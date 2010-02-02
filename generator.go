package main

// A tool to generate random fields.

import (
	"./game"
	"./util"
	"fmt"
	"malloc"
	"rand"
	"runtime"
	"time"
)

const concurrency = 12
const minDifficulty = 40000

// Generates a random field by placing each ship at a random location.
func GenerateField(rng *rand.Rand) (field game.Field) {
	var blocked [game.FieldHeight][game.FieldWidth]bool
	for ship := 0; ship < len(game.ShipLengths); ship++ {
		length := game.ShipLengths[ship]
	retry:
		for {
			r1 := rng.Intn(game.FieldHeight)
			c1 := rng.Intn(game.FieldWidth)
			dir := rng.Intn(2)
			r2 := r1 + (length-1)*dir + 1
			c2 := c1 + (length-1)*(1-dir) + 1
			if r2 > game.FieldHeight || c2 > game.FieldWidth {
				continue
			}
			for r := r1; r < r2; r++ {
				for c := c1; c < c2; c++ {
					if blocked[r][c] {
						continue retry
					}
				}
			}
			// Place ship
			for r := r1; r < r2; r++ {
				for c := c1; c < c2; c++ {
					field[r][c] = true
				}
			}
			// Mark blocked fields
			br1 := util.Max(0, r1-1)
			bc1 := util.Max(0, c1-1)
			br2 := util.Min(game.FieldHeight, r2+1)
			bc2 := util.Min(game.FieldWidth, c2+1)
			for r := br1; r < br2; r++ {
				for c := bc1; c < bc2; c++ {
					blocked[r][c] = true
				}
			}
			break
		}
	}
	return
}

// Continuously generates fields:
func generate() {
	rng := rand.New(rand.NewSource(rand.Int63()))
	for {
		field := GenerateField(rng)
		rows, cols := game.CountShips(&field)
		difficulty := 0
		ch := game.GenerateSolutions(rows, cols)
		for sol := <-ch; sol != nil; sol = <-ch {
			difficulty++
		}
		if difficulty >= minDifficulty {
			fmt.Println(difficulty, game.FormatCounts(&rows), game.FormatCounts(&cols), game.FormatShips(&field))
			malloc.GC()
		}
	}
}

func main() {
	runtime.GOMAXPROCS(concurrency)
	rand.Seed(time.Nanoseconds())
	for i := 0; i < concurrency; i++ {
		go generate()
	}
	<-make(chan int) // block forever
}
