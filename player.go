package game

import "rand"
import "sync"
import "time"

var TimeOut float = 30

var solutionsCache = make(map[string][]*Field)       // caches known solutions
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
					default:
						break notifyWaiters
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
		case <-ticker.C:
			solutions = nil // timer expired!
		}
		ticker.Stop()
	}
	return solutions
}

// Setup returns a random new field set-up
func Setup() *Field {
	// FIXME: template is hard-coded!
	// FIXME: I do have harder templates.
	rows := RowCounts{2, 0, 4, 0, 2, 0, 3, 0, 7, 0, 5, 0, 3, 0, 4, 0}
	cols := ColCounts{1, 1, 2, 2, 2, 3, 2, 2, 2, 2, 2, 2, 2, 2, 2, 1}
	solutions := ListSolutions(rows, cols)
	return solutions[rand.Intn(len(solutions))]
}

func filterShots(solutions []*Field, shots []Shot) []*Field {
	count := 0
	filtered := make([]*Field, len(solutions))
loop:
	for _, solution := range (solutions) {
		for _, shot := range (shots) {
			if solution[shot.R][shot.C] != shot.Hit {
				continue loop
			}
		}
		filtered[count] = solution
		count++
	}
	return filtered[0:count]
}

// SimpleShoot fires at a cell with a maximum probability of hitting, estimating
// this probability as rows[r] + cols[c]. This algorithm is simplistic, but very
// fast.
func SimpleShoot(rows RowCounts, cols ColCounts, shots []Shot) (shootR, shootC int) {
	var shot Field
	for _, s := range (shots) {
		shot[s.R][s.C] = true
	}
	var maxHit, hitCount int
	for r := 0; r < FieldHeight; r++ {
		for c := 0; c < FieldWidth; c++ {
			if !shot[r][c] && rows[r] > 0 && cols[c] > 0 {
				hit := rows[r] + cols[c]
				if hit > maxHit {
					maxHit = hit
					hitCount = 0
				}
				if hit == maxHit {
					hitCount++
					if rand.Intn(hitCount) == 0 {
						shootR, shootC = r, c
					}
				}
			}
		}
	}
	return
}

// Counts how often cell r,c is hit in the given solution set:
func CountHits(solutions []*Field, r, c int) int {
	var count int
	for _, solution := range (solutions) {
		if solution[r][c] {
			count++
		}
	}
	return count
}

// Shoot returns the coordinates of an unoccupied cell to fire at
func Shoot(rows RowCounts, cols ColCounts, shots []Shot) (shootR, shootC int) {
	// Mark cells we've shot at before
	var shot Field
	for _, s := range (shots) {
		shot[s.R][s.C] = true
	}

	// Find all solutions:
	solutions := getSolutions(rows, cols, int64(TimeOut*1e9))
	if solutions == nil {
		// Solver timed out; use a less sophisticated algorithm:
		return SimpleShoot(rows, cols, shots)
	}
	solutions = filterShots(solutions, shots)

	// Count how often each (unfired) cell is hit:
	var hits [FieldHeight][FieldWidth]int
	{
		var children int
		notify := make(chan struct{}, FieldHeight*FieldWidth)
		for r := 0; r < FieldHeight; r++ {
			for c := 0; c < FieldWidth; c++ {
				if !shot[r][c] && rows[r] > 0 && cols[c] > 0 {
					children++
					go func(r, c int) {
						hits[r][c] = CountHits(solutions, r, c)
						notify <- struct{}{}
					}(r, c)
				}
			}
		}
		for children > 0 {
			<-notify
			children--
		}
	}

	// Select an unfired cell with maximum hit probability, at random:
	var cnt, max int
	for r, row := range (hits) {
		for c, hit := range (row) {
			if hit > max {
				max = hit
				cnt = 0
			}
			if hit == max {
				cnt++
				if rand.Intn(cnt) == 0 {
					shootR, shootC = r, c
				}
			}
		}
	}
	return
}
