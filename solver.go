package game

type intField [FieldHeight][FieldWidth]int

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

// placeShips is the solver's workhorse. It takes a partially solved field with
// ships, a field of blocked cells, row and column counts, the ship type to
// place (in kind) and how many have been place of this type (unit), and where
// the last ship was placed (start_r, start_c), and then computes all remaining
// solutions to the grid, which are sent to the results channel.
//
// N.B. this routine should not return before all results from its subproblems
// have been sent to the results channel. Specifically, if the routine spawns
// new goroutines, it should wait for them to finish before returning!
func placeShips(rows []int, cols []int, ships *Field, blocked *intField, kind, unit, start_r, start_c int, results chan *Field, notify chan int) {
	// Check if we have a ship to place next:
	if unit == ShipTypes[kind].Units {
		kind++
		if kind == len(ShipTypes) {
			result := *ships
			results <- &result
			return
		}
		start_r = 0
		start_c = 0
		unit = 0
	}

	// Prepare to spawn child goroutines for solving subproblems in parallel:
	var childNotify chan int
	var children int
	if kind < 1 { // HEURISTIC: spawn children for the toplevel ship only
		childNotify = make(chan int, 20) // expect about 20 valid placements
	}

	// Search over all remaining positions for this type of ship:
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
				if childNotify == nil {
					placeShips(rows, cols, ships, blocked, kind, unit+1, r1, bc2, results, nil)
				} else {
					var rowsCopy RowCounts
					var colsCopy ColCounts
					copy(&rowsCopy, rows)
					copy(&colsCopy, cols)
					shipsCopy := *ships
					blockedCopy := *blocked
					go placeShips(&rowsCopy, &colsCopy, &shipsCopy, &blockedCopy, kind, unit+1, r1, bc2, results, childNotify)
					children++
				}

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

	// Wait for my children to finish, then notify my parent I'm done:
	for ; children > 0; children-- {
		<-childNotify
	}
	if notify != nil {
		notify <- 1
	}
}

// GenerateSolutions writes all solution fields for the given row and column
// counts to a channel (and then a nil value to terminate the list).
func GenerateSolutions(rows RowCounts, cols ColCounts) chan *Field {
	results := make(chan *Field, 1000) // expect a lot of solutions
	go func() {
		placeShips(&rows, &cols, &Field{}, &intField{}, 0, 0, 0, 0, results, nil)
		results <- nil
	}()
	return results
}

// ListSolutions returns a slice with all solutions for the given field counts
func ListSolutions(rows RowCounts, cols ColCounts) (solutions []Field) {
	ch := GenerateSolutions(rows, cols)
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
