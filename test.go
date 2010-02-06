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
		fmt.Println("Worst-case score:", game.GetMaximumScore(strategy))
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
