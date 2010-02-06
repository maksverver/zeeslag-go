package game

import "rand"
import "sync"
import "time"

var TimeOut float = 30

var solutionsCache = make(map[string][]*Field) // caches known solutions
var solutionsNotify = make(map[string]chan []*Field) // notifies waiters
var solutionsCacheMutex sync.Mutex

func getCacheKey(rows RowCounts, cols ColCounts) string {
	return FormatCounts(&rows) + "/" + FormatCounts(&cols)
}

func PurgeCache(rows RowCounts, cols ColCounts) {
	solutionsCacheMutex.Lock()
	solutionsCache[getCacheKey(rows, cols)] = nil, false
	solutionsCacheMutex.Unlock()
}

func getSolutions(rows RowCounts, cols ColCounts, maxWaitNs int64) []*Field {
	key := getCacheKey(rows, cols)
	solutionsCacheMutex.Lock()
	solutions, found := solutionsCache[key]
	if found {
		solutionsCacheMutex.Unlock()
	} else {
		notify, found := solutionsNotify[key]
		if !found {
			notify = make(chan []*Field)
			solutionsNotify[key] = notify
			go func() {
				solutions := ListSolutions(rows, cols)
				solutionsCacheMutex.Lock()
				solutionsCache[key] = solutions
				solutionsNotify[key] = nil, false
				solutionsCacheMutex.Unlock()
			notifyWaiters:
				for {
					select {
					case notify <- solutions:
					default: break notifyWaiters
					}
				}
			}()
		}
		solutionsCacheMutex.Unlock()
		ticker := time.NewTicker(maxWaitNs)
		select {
		/* FIXME: race conditions here! Since we have unlocked the mutex
		   before selecting, the solver code above may forget to notify us! */
		case solutions = <-notify: // solution found!
		case <-ticker.C: solutions = nil // timer expired!
		}
		ticker.Stop()
	}
	return solutions
}

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

// Shoot returns the coordinates of an unoccupied cell to fire at
func Shoot(rows RowCounts, cols ColCounts, shots []Shot) (shootR, shootC int) {
	// Mark cells we've shot at before
	var shot Field
	for _, s := range (shots) {
		shot[s.R][s.C] = true
	}

	// Count how often each cell is hit:
	solutions := getSolutions(rows, cols, int64(TimeOut*1e9))
	if solutions == nil {
		// TODO: use heuristic here!
		return 0, 0
	}

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
