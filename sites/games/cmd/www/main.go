/*
# SPDX-FileCopyrightText: © 2020-2025 franklin <franklin@bitsmasher.net>
#
# SPDX-License-Identifier: MIT
*/
package main

import (
	"html/template"
	"log"
	"net/http"
)

var LayoutDir string = "web/template"
var index *template.Template

type Page struct {
	Title string
}

type Answer struct {
	Title  string
	CSS    string
	Flag   string
	Puzzle string
	Result bool
}

func main() {
	log.Println("Server starting...")
	var err error
	index, err = template.ParseGlob(LayoutDir + "/*.tmpl")
	if err != nil {
		panic(err)
	}
	http.HandleFunc("/", handler)
	http.HandleFunc("/dst/", dst)
	http.HandleFunc("/minecraft/", mc)
	http.HandleFunc("/submission", submission)
	http.HandleFunc("/puzzle1", puzzle1)
	http.HandleFunc("/puzzle2", puzzle2)
	http.HandleFunc("/puzzle3", puzzle3)


    //  --- Make static files available ------------------------------------------
	// http.Handle("/static/", http.StripPrefix("/static/", fs)) // "/static/images/my_image.jpg" will look for "images/my_image.jpg" in the "static" directory.
    fs := http.FileServer(http.Dir("./web/static"))
    http.Handle("/web/static/", http.StripPrefix("/web/static/", fs))
    if err != nil {
        panic(err)
    }
	http.Handle("/static", fs)
	
	http.ListenAndServe(":8080", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	log.Println("Serving index page")
	page := Page{Title: "games.bitsmasher.net"}
	index.ExecuteTemplate(w, "indexPage", page)
}

func dst(w http.ResponseWriter, r *http.Request) {
	log.Println("Serving DST page")
	page := Page{"Do Not Starve"}
	index.ExecuteTemplate(w, "dstPage", page)
}

func mc(w http.ResponseWriter, r *http.Request) {
	log.Println("Serving Minecraft page")
	page := Page{"minecraft"}
	index.ExecuteTemplate(w, "mcPage", page)
}

func submission(w http.ResponseWriter, r *http.Request) {
	data := Answer{Title: "wrong answer headquarters", Flag: r.URL.Query().Get("flag"), Puzzle: r.URL.Query().Get("puzzle"), Result: false}
	log.Printf("For %s you submitted: %s\n", data.Puzzle, data.Flag)

	if (data.Puzzle == "puzzle1") && (data.Flag == "HAHASTRUGGLE") {
		log.Println(w, "Correct answer!")
	} else {
		log.Println(w, "NOPE. Guess the struggle continues!")
	}

	log.Println("Serving submission page")
	index.ExecuteTemplate(w, "submissionPage", data)
}

func puzzle1(w http.ResponseWriter, r *http.Request) {
	log.Println("Serving puzzle1 page")
	page := Page{Title: "puzzle1"}
	index.ExecuteTemplate(w, "puzzle1Page", page)
}

func puzzle2(w http.ResponseWriter, r *http.Request) {
	log.Println("Serving puzzle2 page")
	page := Page{Title: "puzzle2"}
	index.ExecuteTemplate(w, "puzzle2Page", page)
}

func puzzle3(w http.ResponseWriter, r *http.Request) {
	log.Println("Serving puzzle3 page")
	page := Page{Title: "puzzle3"}
	index.ExecuteTemplate(w, "puzzle3Page", page)
}