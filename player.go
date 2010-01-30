package game

import "rand"

/*
// Setup returns a random new field set-up
// FIXME: generates weak-ish fields!
func Setup() Field {
	var total int
	for _, t := range(ShipTypes) {
		total += t.Length * t.Units
	}
	for {
		var rows RowCounts
		var cols ColCounts
		for i := 0; i < total; i++ {
			pos := rand.Intn(FieldHeight*FieldWidth)
			cols[pos%FieldWidth]++
			rows[pos/FieldWidth]++
		}
		solutions := ListSolutions(rows, cols)
		if len(solutions) > 100 {
			return solutions[rand.Intn(len(solutions))]
		}
	}
	return Field{}
}
*/

// Setup returns a random new field set-up
func Setup() Field {
	rows := RowCounts{0, 0, 0, 0, 0, 0, 2, 2, 2, 2, 3, 3, 3, 4, 4, 5}
	cols := ColCounts{0, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2}
	solutions := ListSolutions(rows, cols)
	return solutions[rand.Intn(len(solutions))]
}

func matchShots(shots []Shot, field *Field) bool {
	for _, shot := range (shots) {
		if (*field)[shot.R][shot.C] != shot.Hit {
			return false
		}
	}
	return true
}


// Shoot returns the coordinates of an unoccupied cell to fire at
func Shoot(rows RowCounts, cols ColCounts, shots []Shot) (shootR, shootC int) {

	// Generate solutions
	ch := make(chan *Field)
	go GenerateSolutions(rows, cols, ch)

	// Mark cells we've shot at before
	var shot Field
	for _, s := range (shots) {
		shot[s.R][s.C] = true
	}

	// Count how often each cell is hit:
	var hits [FieldHeight][FieldWidth]int
	for sol := <-ch; sol != nil; sol = <-ch {
		if matchShots(shots, sol) {
			for r, row := range (*sol) {
				for c, hit := range (row) {
					if hit {
						hits[r][c]++
					}
				}
			}
		}
	}

	// Select an unfired cell with maximum hit probability, at random:
	var cnt, max int
	for r, row := range (hits) {
		for c, hit := range (row) {
			if !shot[r][c] && hit >= max {
				if hit > max {
					max = hit
					cnt = 1
				} else {
					cnt++
				}
				if rand.Intn(cnt) == 0 {
					shootR, shootC = r, c
				}
			}
		}
	}
	return
}
