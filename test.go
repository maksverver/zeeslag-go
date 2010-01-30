package main

import "./game"
import "./solver"
import "./io"
import "fmt"

func main() {

	var rows game.RowCounts
	var cols game.ColCounts

	const desc = "5HJ11.4HK16.4HF14.3HC8.3HI4.3HE1.2HO2.2HA5.2HF16.2HK6" // 6346 solutions
	//const desc = "5HL2.4HF16.4HD13.3HD5.3HC11.3HM10.2VD8.2HJ9.2HA15.2HH3"  // 1966 solutions

	if field := io.ParseShips(desc); field == nil {
		fmt.Println("Couldn't parse field description:", desc)
	} else {
		fmt.Print("Parsed field:\n", field)
		rows, cols = game.CountShips(field)
	}

/*
	const rowsDesc = "0.0.0.0.0.0.2.2.2.2.3.3.3.4.4.5"
	const colsDesc = "2.2.2.2.2.2.2.2.2.2.2.2.2.2.2.0"

	if rowsPtr := io.ParseRows(rowsDesc); rowsPtr == nil {
		fmt.Println("Couldn't parse row counts:", rowsDesc)
	} else {
		if colsPtr := io.ParseCols(colsDesc); colsPtr == nil {
			fmt.Println("Couldn't parse column counts:", rowsDesc)
		} else {
			rows = *rowsPtr
			cols = *colsPtr
		}
	}
*/

	fmt.Println("Rows:", rows)
	fmt.Println("Cols:", cols)
	ch := make(chan *game.Field)
	go solver.GenerateSolutions(rows, cols, ch)
	cnt := 0
	for sol := <-ch; sol != nil; sol = <-ch {
		//fmt.Print("Solution:\n", sol)
		cnt++
	}
	fmt.Println(cnt, "solutions found.")
}
