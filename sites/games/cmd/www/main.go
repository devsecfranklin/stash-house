/*
# SPDX-FileCopyrightText: ©2021-2025 franklin <franklin@bitsmasher.net>
#
# SPDX-License-Identifier: MIT
*/
package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
)

type Answer struct {
	Flag   string
	Puzzle string
	Result string
}

func answerHandler(w http.ResponseWriter, r *http.Request) {

	data := Answer{Flag: r.URL.Query().Get("flag"), Puzzle: r.URL.Query().Get("puzzle")}

	fmt.Fprintf(w, "Hello from Go App! You requested: %s\n", r.URL.Path)
	log.Printf("Request received for: %s from %s\n", r.URL.Path, r.RemoteAddr)
	//fmt.Printf(w, "For %s you submitted: %s\n", data.Puzzle, data.Flag)
	log.Printf("For %s you submitted: %s\n", data.Puzzle, data.Flag)

	if (data.Puzzle == "puzzle1") && (data.Flag == "HAHASTRUGGLE") {
		log.Println(w, "Correct answer!")
	} else {
		fmt.Fprintf(w, "NOPE. Guess the struggle continues!\n")
	}
}

func main() {

	// http.Handle("/", http.FileServer(http.Dir("../../web/static/"))) // this is to server static web pages

	/*
		http.HandleFunc("/", answerHandler) // Register the handler for all paths
		log.Println("Server listening on :8080")
		err := http.ListenAndServe(":8080", nil)
		if err != nil {
			log.Fatal(err)
		}
	*/

	// Parse the HTML template
	tmpl, err := template.ParseFiles("index.html")
	if err != nil {
		panic(err)
	}

	// Execute the template with the data and print to standard output
	err = tmpl.Execute(os.Stdout, data)
	if err != nil {
		panic(err)
	}
}
