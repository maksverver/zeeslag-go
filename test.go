package main

import (
	"./game"
	"flag"
	"fmt"
	"rand"
	"time"
)

func main() {
	// Parse command line arguments:
	setupFlag := flag.Bool("Setup", false, "Generate a starting field")
	shipsFlag := flag.String("Ships", "", "Solve a field described as a list of ships")
	rowsFlag := flag.String("Rows", "", "Solve a field with the given row counts (requires -Cols as well)")
	colsFlag := flag.String("Cols", "", "Solve a field with the given column counts (requires -Rows as well)")
	seedFlag := flag.Int64("Seed", 0, "Random seed (0 to pick at random)")
	shotsFlag := flag.String("Shots", "-", "Specify previous shots, and request the next move")
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
		fmt.Println("Random setup: " + game.FormatShips(field))
		rows, cols = game.CountShips(field)
	} else {
		flag.PrintDefaults()
		return
	}
	
	if (*shotsFlag == "-") {
		// No shots passed; just print number of solutions
		fmt.Println("Rows:", game.FormatCounts(&rows))
		fmt.Println("Cols:", game.FormatCounts(&cols))
		ch := game.GenerateSolutions(rows, cols)
		cnt := 0
		for sol := <-ch; sol != nil; sol = <-ch {
			//fmt.Print("Solution:\n", sol)
			cnt++
		}
		fmt.Println(cnt, "solutions found.")
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
