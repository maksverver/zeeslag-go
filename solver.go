package game

// The main solver for the game is implemented here, where solving means to
// find all fields that satisfy a given pair of row and column counts. These
// solutions are used to determine the firing strategy.

import "./util"

// solverstate describes a partial solution, used by placeShips:
type solverState struct {
	rows    RowCounts
	cols    ColCounts
	ships   Field
	blocked [FieldHeight][FieldWidth]int
	results chan *Field
}

// Returns a new copy of a partial solution. The only reason that this is in a
// separate function, is that inlining it into placeShips significantly reduces
// performance. (Probably a compiler bug.)
func copyState(ss *solverState) *solverState {
	copy := *ss
	return &copy
}

// placeShips is the solver's workhorse. It takes a partially solved field with
// ships, a field of blocked cells, row and column counts, the next ship to
// place, and where the last ship was placed (start_r, start_c), and then
// computes all remaining solutions to the grid, which are sent to the results
// channel.
//
// N.B. this routine should not return before all results from its subproblems
// have been sent to the results channel. Specifically, if the routine spawns
// new goroutines, it should wait for them to finish before returning!
func placeShips(ss *solverState, ship, start_r, start_c int, notify chan int) {

	// Check if we need to restart placing ship at the top left corner:
	if ship > 0 && ShipLengths[ship] != ShipLengths[ship-1] {
		start_r = 0
		start_c = 0
	}

	// Prepare to spawn child goroutines for solving subproblems in parallel:
	var childNotify chan int
	var children int
	if ship < 1 { // HEURISTIC: spawn children for the toplevel ship only
		childNotify = make(chan int, 20) // expect about 20 valid placements
	}

	// Search over all remaining positions for this type of ship:
	for dir := 0; dir < 2; dir++ {
		h := dir*(ShipLengths[ship]-1) + 1
		w := (1-dir)*(ShipLengths[ship]-1) + 1

		for r1 := start_r; r1 <= FieldHeight-h; r1++ {
			if ss.rows[r1] < w {
				continue
			}
		loop:
			for c1 := util.Ifc(r1 == start_r, start_c, 0); c1 <= FieldWidth-w; c1++ {

				if ss.cols[c1] < h || ss.blocked[r1][c1] > 0 {
					continue
				}

				// Check if space is available here:
				r2, c2 := r1+h, c1+w
				if c2 > FieldWidth || r2 > FieldHeight {
					continue
				}
				for r := r1; r < r2; r++ {
					if ss.rows[r] < w {
						continue loop
					}
				}
				for c := c1; c < c2; c++ {
					if ss.cols[c] < h {
						continue loop
					}
				}
				for r := r1; r < r2; r++ {
					for c := c1; c < c2; c++ {
						if ss.blocked[r][c] > 0 {
							continue loop
						}
					}
				}

				// Calculate area blocked by the ship:
				br1 := util.Max(0, r1-1)
				bc1 := util.Max(0, c1-1)
				br2 := util.Min(FieldHeight, r2+1)
				bc2 := util.Min(FieldWidth, c2+1)

				// Claim space
				for r := r1; r < r2; r++ {
					ss.rows[r] -= w
				}
				for c := c1; c < c2; c++ {
					ss.cols[c] -= h
				}
				for r := r1; r < r2; r++ {
					for c := c1; c < c2; c++ {
						ss.ships[r][c] = true
					}
				}
				for r := br1; r < br2; r++ {
					for c := bc1; c < bc2; c++ {
						ss.blocked[r][c]++
					}
				}

				// Quick check to see if field is still solvable:
				for r := br1; r < br2; r++ {
					if ss.rows[r] == 1 &&
						(r == 0 || ss.rows[r-1] == 0) &&
						(r == FieldHeight-1 || ss.rows[r+1] == 0) {
						goto unsolvable
					}
				}
				for c := bc1; c < bc2; c++ {
					if ss.cols[c] == 1 &&
						(c == 0 || ss.cols[c-1] == 0) &&
						(c == FieldWidth-1 || ss.cols[c+1] == 0) {
						goto unsolvable
					}
				}

				// Solve recursively
				if ship+1 == len(ShipLengths) {
					result := ss.ships // make a copy
					ss.results <- &result
				} else if childNotify == nil {
					placeShips(ss, ship+1, r1, bc2, nil)
				} else {
					go placeShips(copyState(ss), ship+1, r1, bc2, childNotify)
					children++
				}

				// Return claimed space
			unsolvable:
				for r := r1; r < r2; r++ {
					ss.rows[r] += w
				}
				for c := c1; c < c2; c++ {
					ss.cols[c] += h
				}
				for r := r1; r < r2; r++ {
					for c := c1; c < c2; c++ {
						ss.ships[r][c] = false
					}
				}
				for r := br1; r < br2; r++ {
					for c := bc1; c < bc2; c++ {
						ss.blocked[r][c]--
					}
				}
			}
		}
	}

	for ; children > 0; children-- {
		<-childNotify  // wait for child to finish
	}
	if notify != nil {
		notify <- 1  // notify parent I'm done
	}
}

// GenerateSolutions writes all solution fields for the given row and column
// counts to a channel (and then a nil value to terminate the list).
func GenerateSolutions(rows RowCounts, cols ColCounts) chan *Field {
	results := make(chan *Field, 1000) // expect a lot of solutions
	go func() {
		state := solverState{rows: rows, cols: cols, results: results}
		placeShips(&state, 0, 0, 0, nil)
		results <- nil
	}()
	return results
}

// ListSolutions returns a slice with all solutions for the given field counts
func ListSolutions(rows RowCounts, cols ColCounts) (solutions []*Field) {
	ch := GenerateSolutions(rows, cols)
	for sol := <-ch; sol != nil; sol = <-ch {
		i := len(solutions)
		if i == cap(solutions) {
			tmp := make([]*Field, i, util.Max(2*i, 16))
			copy(tmp, solutions)
			solutions = tmp
		}
		solutions = solutions[0 : i+1]
		solutions[i] = sol
	}
	return
}
