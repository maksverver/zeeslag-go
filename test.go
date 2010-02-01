package main

import (
	"./game"
	"fmt"
)

func main() {

	var rows game.RowCounts
	var cols game.ColCounts

	const desc = "5HJ11.4HK16.4HF14.3HC8.3HI4.3HE1.2HO2.2HA5.2HF16.2HK6" // 6346 solutions
	//const desc = "5HL2.4HF16.4HD13.3HD5.3HC11.3HM10.2VD8.2HJ9.2HA15.2HH3"  // 1966 solutions

	/*
		if field := game.ParseShips(desc); field == nil {
			fmt.Println("Couldn't parse field description:", desc)
		} else {
			fmt.Print("Parsed field:\n", field)
			rows, cols = game.CountShips(field)
		}
	*/

	/* HARD. Took 6s, or 6.2s after adding some more parameters. 190k solutions? */
	const rowsDesc = "1.2.2.2.2.2.2.3.2.2.2.2.1.2.2.1"
	const colsDesc = "2.4.0.5.0.4.0.5.0.7.0.3.0.0.0.0"

	/* MEDIUM. Takes about 220ms. 11,124 solutions
	const rowsDesc = "2.0.1.3.1.2.4.0.6.0.3.0.0.5.0.3"
	const colsDesc = "4.1.1.2.1.1.2.2.2.2.3.2.2.2.2.1"
	*/

	if rowsPtr := game.ParseRows(rowsDesc); rowsPtr == nil {
		fmt.Println("Couldn't parse row counts:", rowsDesc)
	} else {
		if colsPtr := game.ParseCols(colsDesc); colsPtr == nil {
			fmt.Println("Couldn't parse column counts:", rowsDesc)
		} else {
			rows = *rowsPtr
			cols = *colsPtr
		}
	}

	//fmt.Println("Shots", game.ParseShots("WM15.SE15.SD15"))

	fmt.Println("Rows:", rows)
	fmt.Println("Cols:", cols)
	ch := game.GenerateSolutions(rows, cols)
	cnt := 0
	for sol := <-ch; sol != nil; sol = <-ch {
		//fmt.Print("Solution:\n", sol)
		cnt++
	}
	fmt.Println(cnt, "solutions found.")
}
