package main

import (
	"./game"
	"http"
	"flag"
	"fmt"
	"log"
	"io"
	"strconv"
	"time"
)

func PlayerServer(conn *http.Conn, request *http.Request) {
	if request.ParseForm() != nil {
		conn.WriteHeader(http.StatusInternalServerError)
		return
	}
	nsBegin := time.Nanoseconds()
	var response string
	var succeeded bool
	if action, ok := request.Form["Action"]; !ok {
		response = "no Action parameter supplied"
	} else {
		switch action[0] {
		case "Ships":
			field := game.Setup()
			response = game.FormatShips(&field)
			succeeded = true
		case "Fire":
			if rows, ok := request.Form["Rows"]; !ok {
				response = "no Rows parameter supplied"
			} else if rows := game.ParseRows(rows[0]); rows == nil {
				response = "invalid row count data"
			} else if cols, ok := request.Form["Cols"]; !ok {
				response = "no Cols parameter supplied"
			} else if cols := game.ParseCols(cols[0]); cols == nil {
				response = "invalid column count data"
			} else if shots, ok := request.Form["Shots"]; !ok {
				response = "no Shots parameter supplied"
			} else if shots := game.ParseShots(shots[0]); shots == nil {
				response = "invalid shot data"
			} else {
				r, c := game.Shoot(*rows, *cols, shots)
				response = game.FormatCoords(r, c)
				succeeded = true
			}
		case "Finished":
			succeeded = true
		default:
			response = "unknown Action value supplied"
		}
	}

	// Log request details:
	{
		var parameters string
		for name, value := range (request.Form) {
			if len(name) > 0 {
				if parameters != "" {
					parameters += " "
				}
				parameters += (name + "=" + value[0])
			}
		}
		delay := fmt.Sprintf("\t(%.3fs)", float64(time.Nanoseconds()-nsBegin)/1e9)
		log.Stdout(
			"\t("+conn.RemoteAddr+")",
			"\t"+parameters,
			"\t"+fmt.Sprintf("%v", succeeded),
			"\t"+response,
			"\t"+delay)
	}


	// Write response to client:
	conn.SetHeader("Content-Type", "text/plain")
	if !succeeded {
		response = "ERROR: " + response + "!"
	}
	io.WriteString(conn, response)
}

func main() {
	// Parse command line arguments:
	host := flag.String("h", "", "hostname to bind")
	port := flag.Int("p", 14000, "port to bind")
	path := flag.String("r", "/player", "root path for player")
	flag.Parse()
	addr := *host + ":" + strconv.Itoa(*port)

	// Start an HTTP server with a player handler:
	http.Handle(*path, http.HandlerFunc(PlayerServer))
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Stderr("Could not serve on address " + addr + ": " + err.String())
	}
}
