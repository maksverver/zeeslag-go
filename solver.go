package solver

import "./game"

const Height, Width = game.FieldHeight, game.FieldWidth

// min returns the minimum value of its two arguments
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// max returns the maximum value of its two arguments
func max(a, b int) int {
	if a < b {
		return b
	}
	return a
}

// ifc returns the a or b argument, depending on whether c is true or false
func ifc(c bool, a, b int) int {
	if c {
		return a
	}
	return b
}

// GenerateSolutions writes all solution fields for the given row and column
// counts to a channel (and then a nil value to terminate the list).
func GenerateSolutions(rows game.RowCounts, cols game.ColCounts, ch chan *game.Field) {
	var ships game.Field
	var blocked [Height][Width]int
	var placeShips func(kind, unit, start_r, start_c int)

	placeShips = func(kind, unit, start_r, start_c int) {
		// Check if we have a ship to place next:
		if unit == game.ShipTypes[kind].Units {
			kind++
			if kind == len(game.ShipTypes) {
				result := ships
				ch <- &result
				return
			}
			start_r = 0
			start_c = 0
			unit = 0
		}

		for dir := 0; dir < 2; dir++ {
			h := dir*(game.ShipTypes[kind].Length-1) + 1
			w := (1-dir)*(game.ShipTypes[kind].Length-1) + 1

			for r1 := start_r; r1 <= Height-h; r1++ {
				if rows[r1] < w {
					continue
				}

				c1 := 0
				if r1 == start_r {
					c1 = start_c
				}
			loop:
				for ; c1 <= Width-w; c1++ {

					if cols[c1] < h || blocked[r1][c1] > 0 {
						continue
					}

					// Check if space is available here:
					r2, c2 := r1+h, c1+w
					if c2 > Width || r2 > Height {
						continue
					}
					for r := r1; r < r2; r++ {
						if rows[r] < w {
							continue loop
						}
					}
					for c := c1; c < c2; c++ {
						if cols[c] < h {
							continue loop
						}
					}
					for r := r1; r < r2; r++ {
						for c := c1; c < c2; c++ {
							if blocked[r][c] > 0 {
								continue loop
							}
						}
					}

					// Calculate area blocked by the ship:
					br1 := max(0, r1-1)
					bc1 := max(0, c1-1)
					br2 := min(Height, r2+1)
					bc2 := min(Width, c2+1)

					// Claim space
					for r := r1; r < r2; r++ {
						rows[r] -= w
					}
					for c := c1; c < c2; c++ {
						cols[c] -= h
					}
					for r := r1; r < r2; r++ {
						for c := c1; c < c2; c++ {
							ships[r][c] = true
						}
					}
					for r := br1; r < br2; r++ {
						for c := bc1; c < bc2; c++ {
							blocked[r][c]++
						}
					}

					// Quick check to see if field is still solvable:
					for r := br1; r < br2; r++ {
						if rows[r] == 1 && (r == 0 || rows[r-1] == 0) &&
							(r == Height-1 || rows[r+1] == 0) {
							goto unsolvable
						}
					}
					for c := bc1; c < bc2; c++ {
						if cols[c] == 1 && (c == 0 || cols[c-1] == 0) &&
							(c == Width-1 || cols[c+1] == 0) {
							goto unsolvable
						}
					}

					// Solve recursively
					placeShips(kind, unit+1, r1, bc2)

				unsolvable:

					// Return claimed space
					for r := r1; r < r2; r++ {
						rows[r] += w
					}
					for c := c1; c < c2; c++ {
						cols[c] += h
					}
					for r := r1; r < r2; r++ {
						for c := c1; c < c2; c++ {
							ships[r][c] = false
						}
					}
					for r := br1; r < br2; r++ {
						for c := bc1; c < bc2; c++ {
							blocked[r][c]--
						}
					}
				}
			}
		}
	}

	// Generate all solutions, then write a nil to signal end of the list:
	placeShips(0, 0, 0, 0)
	ch <- nil
}
