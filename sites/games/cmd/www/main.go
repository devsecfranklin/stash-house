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
	"sync"
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

var (
    tpl *template.Template
    loadOnce sync.Once
    tplErr error
)

func loadTemplate() {
    loadOnce.Do(func() {
        tpl, tplErr = template.ParseGlob("./web/template/*.tmpl")
        if tplErr != nil {
            log.Fatalf("Error loading template: %v", tplErr)
        }
    })
}

func main() {
	var err error

	log.Println("Server starting...")

	//  --- Make static files available ------------------------------------------
	// http.Handle("/static/", http.StripPrefix("/static/", fs))
	// "/static/images/my_image.jpg" will look for "images/my_image.jpg" in the "static" directory.
	fs := http.FileServer(http.Dir("./web/static"))
	http.Handle("/web/static/", http.StripPrefix("/web/static/", fs))
	//http.Handle("/static", fs)

	loadTemplate()
    if tplErr != nil {
        http.Error(w, "Error loading template", http.StatusInternalServerError)
        return
    }

	index, err = template.ParseGlob(LayoutDir + "/*.tmpl")
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/civilization2", civ2)
	http.HandleFunc("/dst", dst)
	http.HandleFunc("/linux_games", lg)
	http.HandleFunc("/minecraft", mc)
	http.HandleFunc("/submission", submission)
	http.HandleFunc("/puzzle1", puzzle1)
	http.HandleFunc("/puzzle2", puzzle2)
	http.HandleFunc("/puzzle3", puzzle3)
	http.HandleFunc("/puzzle4", puzzle4)
	http.HandleFunc("/", handler)

	http.ListenAndServe(":8080", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	log.Println("Serving index page")

	// w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	// w.Header().Set("Pragma", "no-cache")
	// w.Header().Set("Expires", "0")

	loadTemplate()
    if tplErr != nil {
        http.Error(w, "Error loading template", http.StatusInternalServerError)
        return
    }
	page := Page{Title: "games.bitsmasher.net"}
	index.ExecuteTemplate(w, "indexPage", page)
}

func civ2(w http.ResponseWriter, r *http.Request) {
	var tmpl = template.Must(template.ParseFiles("./web/template/civilization2.tmpl"))

	log.Println("Serving civ2 page")

	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	loadTemplate()
    if tplErr != nil {
        http.Error(w, "Error loading template", http.StatusInternalServerError)
        return
    }
	page := Page{"Civilization ]["}

	err := tmpl.ExecuteTemplate(w, "civ2Page", page)
	if err != nil {
        log.Println(err)
		//http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func dst(w http.ResponseWriter, r *http.Request) {
	log.Println("Serving DST page")
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
	page := Page{"Do Not Starve"}
	index.ExecuteTemplate(w, "dstPage", page)
}

func lg(w http.ResponseWriter, r *http.Request) {
	log.Println("Serving linux games page")
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
	page := Page{"Linux Games"}
	index.ExecuteTemplate(w, "lgPage", page)
}

func mc(w http.ResponseWriter, r *http.Request) {
	log.Println("Serving Minecraft page")
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
	page := Page{"minecraft"}
	index.ExecuteTemplate(w, "mcPage", page)
}

func submission(w http.ResponseWriter, r *http.Request) {
	data := Answer{Title: "wrong answer headquarters", Flag: r.URL.Query().Get("flag"), Puzzle: r.URL.Query().Get("puzzle"), Result: false}
    var tmpl = template.Must(template.ParseFiles("./web/template/submission.tmpl"))

	// log.Printf("For %s you submitted: %s\n", data.Puzzle, data.Flag)

	// w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	// w.Header().Set("Pragma", "no-cache")
	// w.Header().Set("Expires", "0")

	loadTemplate()
    if tplErr != nil {
        http.Error(w, "Error loading template", http.StatusInternalServerError)
        return
    }


	if ((data.Puzzle == "puzzle1") && (data.Flag == "HAHASTRUGGLE")) || ((data.Puzzle == "puzzle2-1") && (data.Flag == "DELICIOUSTEARS")) || ((data.Puzzle == "puzzle2-2") && (data.Flag == "COMPUTERSAREHARD")) || ((data.Puzzle == "puzzle3") && (data.Flag == "ENJOYTHEPAIN")) {
		log.Println(w, "Correct answer!")
		log.Println("Serving correct page")

		data.Title = "YOU DID IT"
		data.Result = true

		tmpl = template.Must(template.ParseFiles("./web/template/correct.tmpl"))
		// index.ExecuteTemplate(w, "correctPage", data)
		err := tmpl.ExecuteTemplate(w, "correctPage", data)
		if err != nil {
			log.Println(err)
			//http.Error(w, err.Error(), http.StatusInternalServerError)
		}

	} else {
		
		log.Println(w, "NOPE. Guess the struggle continues!")
		log.Println("Serving submission page")
		data.Title = "NOPE. Guess the struggle continues!"

		//index.ExecuteTemplate(w, "submissionPage", data)
		err := tmpl.ExecuteTemplate(w, "submissionPage", data)
		if err != nil {
			log.Println(err)
			//http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
	
}

func puzzle1(w http.ResponseWriter, r *http.Request) {
	log.Println("Serving puzzle1 page")
	// w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	// w.Header().Set("Pragma", "no-cache")
	// w.Header().Set("Expires", "0")
	page := Page{Title: "puzzle1"}
	index.ExecuteTemplate(w, "puzzle1Page", page)
}

func puzzle2(w http.ResponseWriter, r *http.Request) {
	log.Println("Serving puzzle2 page")

	var tmpl = template.Must(template.ParseFiles("./web/template/puzzle2.tmpl"))

	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	page := Page{Title: "puzzle2"}

	err := tmpl.ExecuteTemplate(w, "puzzle2Page", page)
	if err != nil {
        log.Println(err)
		//http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func puzzle3(w http.ResponseWriter, r *http.Request) {
	log.Println("Serving puzzle3 page")
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
	page := Page{Title: "puzzle3"}
	index.ExecuteTemplate(w, "puzzle3Page", page)
}

func puzzle4(w http.ResponseWriter, r *http.Request) {
	log.Println("Serving puzzle4 page")
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
	page := Page{Title: "puzzle4"}
	index.ExecuteTemplate(w, "puzzle4Page", page)
}
