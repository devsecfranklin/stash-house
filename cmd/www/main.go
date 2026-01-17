/*
# SPDX-FileCopyrightText: ©2025 franklin <smoooth.y62wj@passmail.net>
#
# SPDX-License-Identifier: MIT
*/

package main

import (
	"fmt"
	"html/template"
	"internal/logging"
	"log"
	"net/http"
)

type (
	Stuff struct {
		token string
	}

	// twitchUser struct { // A simple struct to ht
	OauthToken struct { // OauthToken data structure passed to the template
		ClientID  string
		StateRand string
		Title     string
	}

	Page struct { // Page data structure for a generic page
		Title string
	}
)

var (
	err error
	LayoutDir string = "template/www"
	tmpls     *template.Template
)

func main() {
	fs := http.FileServer(http.Dir("static/www"))
	http.Handle("/static/www/", http.StripPrefix("/static/www/", fs))

	tmpls, err = template.ParseGlob(LayoutDir + "/*.tmpl")
	if err != nil {
		panic(err)
	}

	// Register the handler for all paths
	http.HandleFunc("/", handler)
	http.HandleFunc("/lab", labPageHandler)

	logging.Log_header("Server listening on :8080")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		logging.Log_fatal(fmt.Sprintf("Server failed to start: %v", err))
	}
}

func handler(w http.ResponseWriter, r *http.Request) { // handler for the root path
	log.Println("Serving index page")

	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	page := Page{"www.bitsmasher.net"}

	err := tmpls.ExecuteTemplate(w, "indexPage", page) // Assuming you have an "indexPage" template
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal server error: Could not render index page.", http.StatusInternalServerError)
	}
}

func labPageHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Serving lab page")

	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	page := Page{"www.bitsmasher.net/lab"}

	err := tmpls.ExecuteTemplate(w, "labPage", page)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal server error: Could not render index page.", http.StatusInternalServerError)
	}
}
