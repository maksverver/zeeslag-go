package game

import "rand"
import "sync"

var solutionsCache = make(map[string][]*Field)
var solutionsCacheMutex sync.Mutex

// Setup returns a random new field set-up
func Setup() *Field {
	// FIXME: template is hard-coded!
	// FIXME: I do have harder templates.
	rows := RowCounts{2,0,4,0,2,0,3,0,7,0,5,0,3,0,4,0}
	cols := ColCounts{1,1,2,2,2,3,2,2,2,2,2,2,2,2,2,1}
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

func cacheKey(rows RowCounts, cols ColCounts) string {
	return FormatCounts(&rows) + "/" + FormatCounts(&cols)
}

func PurgeCache(rows RowCounts, cols ColCounts) {
	solutionsCacheMutex.Lock()
	solutionsCache[cacheKey(rows, cols)] = nil, false
	solutionsCacheMutex.Unlock()
}

func getSolutions(rows RowCounts, cols ColCounts) []*Field {
	// Lock the mutex, so we generate only one set of solutions at a time.
	key := cacheKey(rows, cols)
	solutionsCacheMutex.Lock()
	solutions, found := solutionsCache[key]
	if !found {
		// Not found in cache; generate from scratch:
		solutions = ListSolutions(rows, cols)
		solutionsCache[key] = solutions
	}
	solutionsCacheMutex.Unlock()
	return solutions
}

// Shoot returns the coordinates of an unoccupied cell to fire at
func Shoot(rows RowCounts, cols ColCounts, shots []Shot) (shootR, shootC int) {
	// Mark cells we've shot at before
	var shot Field
	for _, s := range (shots) {
		shot[s.R][s.C] = true
	}

	// Count how often each cell is hit:
	var hits [FieldHeight][FieldWidth]int
	for _, sol := range (getSolutions(rows, cols)) {
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
