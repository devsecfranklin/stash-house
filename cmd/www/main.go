/*
# SPDX-FileCopyrightText: ©2021-2025 franklin <franklin@bitsmasher.net>
#
# SPDX-License-Identifier: MIT
*/

package www

import (
	"fmt"

	"log"
	"net/http"
)

type Answer struct {
	Flag   string
	Puzzle string
	Result string
}

func answerHandler(w http.ResponseWriter, r *http.Request) {

	data := Answer{Flag: r.URL.Query().Get("flag"), Puzzle: r.URL.Query().Get("puzzle")}

	//fmt.Fprintf(w, "Hello from Go App! You requested: %s\n", r.URL.Path)
	//log.Printf("Request received for: %s from %s\n", r.URL.Path, r.RemoteAddr)
	//fmt.Fprintf(w, "For %s you submitted: %s\n", Puzzle, Flag)
	log.Printf("For %s you submitted: %s\n", data.Puzzle, data.Flag)

	if (data.Puzzle == "puzzle1") && (data.Flag == "HAHASTRUGGLE") {
		log.Fprintf(w, "Correct answer!\n")
	} else {
		fmt.Fprintf(w, "NOPE. Guess the struggle continues!\n")
	}
}

func main() {

	// Serve static files
	fs := http.FileServer(http.Dir("./static"))
	router.Handle("/static/*", http.StripPrefix("/static/", fs))
	// Register the handler for all paths
	http.HandleFunc("/", answerHandler)

	// http.Handle("/", http.FileServer(http.Dir("../../web/static/"))) // this is to server static web pages

	log.Println("Server listening on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
