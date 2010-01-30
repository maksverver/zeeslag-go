package game

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
func GenerateSolutions(rows RowCounts, cols ColCounts, ch chan *Field) {
	var ships Field
	var blocked [FieldHeight][FieldWidth]int
	var placeShips func(kind, unit, start_r, start_c int)

	placeShips = func(kind, unit, start_r, start_c int) {
		// Check if we have a ship to place next:
		if unit == ShipTypes[kind].Units {
			kind++
			if kind == len(ShipTypes) {
				result := ships
				ch <- &result
				return
			}
			start_r = 0
			start_c = 0
			unit = 0
		}

		for dir := 0; dir < 2; dir++ {
			h := dir*(ShipTypes[kind].Length-1) + 1
			w := (1-dir)*(ShipTypes[kind].Length-1) + 1

			for r1 := start_r; r1 <= FieldHeight-h; r1++ {
				if rows[r1] < w {
					continue
				}

				c1 := 0
				if r1 == start_r {
					c1 = start_c
				}
			loop:
				for ; c1 <= FieldWidth-w; c1++ {

					if cols[c1] < h || blocked[r1][c1] > 0 {
						continue
					}

					// Check if space is available here:
					r2, c2 := r1+h, c1+w
					if c2 > FieldWidth || r2 > FieldHeight {
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
					br2 := min(FieldHeight, r2+1)
					bc2 := min(FieldWidth, c2+1)

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
						if rows[r] == 1 &&
							(r == 0 || rows[r-1] == 0) &&
							(r == FieldHeight-1 || rows[r+1] == 0) {
							goto unsolvable
						}
					}
					for c := bc1; c < bc2; c++ {
						if cols[c] == 1 &&
							(c == 0 || cols[c-1] == 0) &&
							(c == FieldWidth-1 || cols[c+1] == 0) {
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


// Returns a slice with all solutions for the given row/column counts
func ListSolutions(rows RowCounts, cols ColCounts) (solutions []Field) {
	ch := make(chan *Field)
	go GenerateSolutions(rows, cols, ch)
	for sol := <-ch; sol != nil; sol = <-ch {
		i := len(solutions)
		if i == cap(solutions) {
			tmp := make([]Field, i, max(2*i, 16))
			copy(tmp, solutions)
			solutions = tmp
		}
		solutions = solutions[0 : i+1]
		solutions[i] = *sol
	}
	return
}
