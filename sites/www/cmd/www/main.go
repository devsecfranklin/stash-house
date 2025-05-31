/*
# SPDX-FileCopyrightText: ©2021-2025 franklin <franklin@bitsmasher.net>
#
# SPDX-License-Identifier: MIT
*/

package main

import (
	"fmt"
	"log"
	"net/http"
)

func answerHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello from Go App! You requested: %s\n", r.URL.Path)
	log.Printf("Request received for: %s from %s\n", r.URL.Path, r.RemoteAddr)
	flag := r.URL.Query().Get("flag")
	puzzle := r.URL.Query().Get("puzzle")
	fmt.Fprintf(w, "For %s you submitted: %s\n", puzzle, flag)
	log.Printf("For %s you submitted: %s\n", puzzle, flag)

	if(puzzle == "puzzle1") && (flag == "HAHASTRUGGLE") {
        fmt.Fprintf(w, "Correct answer!\n")
	} else {
		fmt.Fprintf(w, "NOPE. Guess the struggle continues!\n")
    }
}

func main() {

	// Register the handler for all paths
	http.HandleFunc("/", answerHandler)

	// http.Handle("/", http.FileServer(http.Dir("../../web/static/"))) // this is to server static web pages

	log.Println("Server listening on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}

