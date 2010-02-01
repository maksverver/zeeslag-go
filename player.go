package game

import "rand"

// Setup returns a random new field set-up
func Setup() *Field {
	// FIXME: template is hard-coded!
	rows := RowCounts{0, 4, 0, 2, 0, 5, 0, 3, 0, 5, 0, 5, 0, 4, 0, 2}
	cols := ColCounts{1, 2, 2, 3, 4, 2, 1, 1, 1, 2, 2, 2, 2, 2, 2, 1}
	solutions := ListSolutions(rows, cols)
	return solutions[rand.Intn(len(solutions))]
}

func matchShots(shots []Shot, field *Field) bool {
	for _, shot := range (shots) {
		if field[shot.R][shot.C] != shot.Hit {
			return false
		}
	}
	return true
}

// Shoot returns the coordinates of an unoccupied cell to fire at
func Shoot(rows RowCounts, cols ColCounts, shots []Shot) (shootR, shootC int) {
	solutions := ListSolutions(rows, cols)

	// Mark cells we've shot at before
	var shot Field
	for _, s := range (shots) {
		shot[s.R][s.C] = true
	}

	// Count how often each cell is hit:
	var hits [FieldHeight][FieldWidth]int
	for _, sol := range(solutions) {
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
