package main

import (
	"./game"
	"flag"
	"fmt"
	"rand"
	"time"
)

func printIndent(indent int) {
	for ; indent > 0; indent-- {
		fmt.Print(" ")
	}
}

func printStrategy(strategy *game.Strategy, indent int) {
	printIndent(indent)
	for i, f := range (strategy.Shots) {
		if i > 0 {
			fmt.Print(" ")
		}
		fmt.Print(game.FormatCoords(game.DecodeCoords(f)))
	}
	fmt.Println()
	if strategy.IfHit != nil {
		printIndent(indent)
		fmt.Println("if hit:")
		printStrategy(strategy.IfHit, indent+4)
	}
	if strategy.IfMiss != nil {
		printIndent(indent)
		fmt.Println("if miss:")
		printStrategy(strategy.IfMiss, indent+4)
	}
}

type caseDepth struct {
	field game.Field
	depth int
}

func calcCases(field *game.Field, strategy *game.Strategy, depth int, results chan *caseDepth) {
	for _, shot := range(strategy.Shots) {
		r,c := game.DecodeCoords(shot)
		field[r][c] = true
	}
	depth += len(strategy.Shots)
	if strategy.IfHit != nil {
		calcCases(field, strategy.IfHit, depth, results)
	} else {
		results <- &caseDepth{*field, depth}
	}
	if strategy.IfMiss != nil {
		r,c := game.DecodeCoords(strategy.Shots[len(strategy.Shots) - 1])
		field[r][c] = false
		calcCases(field, strategy.IfMiss, depth, results)
	}
	for _, shot := range(strategy.Shots) {
		r,c := game.DecodeCoords(shot)
		field[r][c] = false
	}
}

func generateCases(strategy *game.Strategy) chan *caseDepth {
	results := make(chan *caseDepth, 100)
	go func() {
		calcCases(&game.Field{}, strategy, 0, results)
		results <- nil
	}()
	return results
}

func simpleDifficulty(rows *game.RowCounts, cols *game.ColCounts, field *game.Field) float {
	var dif, cnt, tot float
	for r := 0; r < game.FieldHeight; r++ {
		for c := 0; c < game.FieldWidth; c++ {
			if field[r][c] {
				val := float(rows[r] + cols[c])
				if val > dif {
					dif = val
				}
			}
		}
	}
	for r := 0; r < game.FieldHeight; r++ {
		for c := 0; c < game.FieldWidth; c++ {
			if float(rows[r] + cols[c]) == dif {
				tot++
				if field[r][c] {
					cnt++
				}
			}
		}
	}
	return dif + cnt/(tot + 1)
}

func main() {
	// Parse command line arguments:
	setupFlag := flag.Bool("Setup", false, "Generate a starting field")
	shipsFlag := flag.String("Ships", "", "Solve a field described as a list of ships")
	rowsFlag := flag.String("Rows", "", "Solve a field with the given row counts (requires -Cols as well)")
	colsFlag := flag.String("Cols", "", "Solve a field with the given column counts (requires -Rows as well)")
	seedFlag := flag.Int64("Seed", 0, "Random seed (0 to pick at random)")
	shotsFlag := flag.String("Shots", "-", "Specify previous shots, and request the next move")
	flag.FloatVar(&game.TimeOut, "TimeOut", game.TimeOut, "Maximum time to spend on solving")
	flag.Parse()

	// Seed pseudo-random number generator:
	if *seedFlag != 0 {
		rand.Seed(*seedFlag)
	} else {
		rand.Seed(time.Nanoseconds())
	}

	var rows game.RowCounts
	var cols game.ColCounts

	if *shipsFlag != "" {
		// Parse row/column counts from Ships flag:
		if field := game.ParseShips(*shipsFlag); field == nil {
			fmt.Println("Couldn't parse field description:", *shipsFlag)
			return
		} else {
			rows, cols = game.CountShips(field)
		}
	} else if *rowsFlag != "" || *colsFlag != "" {
		// Parse row/column counts from Rows and Cols flags:
		if rowsPtr := game.ParseRows(*rowsFlag); rowsPtr == nil {
			fmt.Println("Couldn't parse row counts:", *rowsFlag)
			return
		} else {
			rows = *rowsPtr
		}
		if colsPtr := game.ParseCols(*colsFlag); colsPtr == nil {
			fmt.Println("Couldn't parse column counts:", *colsFlag)
			return
		} else {
			cols = *colsPtr
		}
	} else if *setupFlag {
		// Set up a random field:
		field := game.Setup()
		fmt.Println("Random setup:", game.FormatShips(field))
		rows, cols = game.CountShips(field)
	} else {
		flag.PrintDefaults()
		return
	}

	if *shotsFlag == "-" {
		// No shots passed; just print number of solutions
		fmt.Println("Rows:", game.FormatCounts(&rows))
		fmt.Println("Cols:", game.FormatCounts(&cols))
		solutions := game.ListSolutions(rows, cols)
		fmt.Println(len(solutions), "solutions found.")
		strategy := game.CreateStrategy(solutions)
		fmt.Println("Expected score:", game.GetExpectedScore(strategy))
		wc := game.GetMaximumScore(strategy)
		fmt.Println("Worst-case score:", wc)
		fmt.Println("Bad cases:")
		ch := generateCases(strategy)
		for cd := <- ch; cd != nil; cd = <- ch {
			if wc - cd.depth <= 3 {
				// first column: optimal search depth, the higher the better
				// second column: max(rows[r]+cols[c]) of cells containing ships, the lower the better
				fmt.Println(cd.depth, simpleDifficulty(&rows, &cols, &cd.field), game.FormatShips(&cd.field))
			}
		}
	} else {
		shots := game.ParseShots(*shotsFlag)
		if shots == nil {
			fmt.Println("Couldn't parse shots:", *shotsFlag)
		} else {
			// Determine best move:
			r, c := game.Shoot(rows, cols, shots)
			fmt.Println(game.FormatCoords(r, c))
		}
	}
}
