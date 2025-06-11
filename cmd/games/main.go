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
	"strings"
	"regexp"
)

var LayoutDir string = "template/games"
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
	var err error

	log.Println("Server starting...")

	table_status := Table{name: "users", db: "games", exist: false}
	table_status.exist = check_table_exist(&table_status) 
	if ! table_status.exist {
		table_status.exist = CrudDB(&table_status)
	}
	
	//  --- Make static files available ------------------------------------------
	// http.Handle("/static/", http.StripPrefix("/static/", fs))
	// "/static/images/my_image.jpg" will look for "images/my_image.jpg" in the "static" directory.
	fs := http.FileServer(http.Dir("./static/games"))
	http.Handle("/static/games/", http.StripPrefix("/static/games/", fs))
	//http.Handle("/static", fs)

	index, err = template.ParseGlob(LayoutDir + "/*.tmpl")
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/", handler)
	http.HandleFunc("/civilization2", civ2)
	http.HandleFunc("/dst", dst)
	http.HandleFunc("/linux_games", lg)
	http.HandleFunc("/minecraft", mc)
	http.HandleFunc("/submission", submission)
	http.HandleFunc("/puzzle1", puzzle1)
	http.HandleFunc("/puzzle2", puzzle2)
	http.HandleFunc("/puzzle3", puzzle3)
	http.HandleFunc("/puzzle4", puzzle4)
	http.ListenAndServe(":8080", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	log.Println("Serving index page")

	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	page := Page{Title: "games.bitsmasher.net"}

	err := index.ExecuteTemplate(w, "indexPage", page)
	if err != nil {
		log.Println(err) // http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

func civ2(w http.ResponseWriter, r *http.Request) {
	log.Println("Serving civ2 page")

	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	page := Page{"Civilization ]["}

	err := index.ExecuteTemplate(w, "civ2Page", page)
	if err != nil {
		log.Println(err) // http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

func dst(w http.ResponseWriter, r *http.Request) {
	log.Println("Serving DST page")

	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	page := Page{"Do Not Starve"}

	err := index.ExecuteTemplate(w, "dstPage", page)
	if err != nil {
		log.Println(err) // http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

func lg(w http.ResponseWriter, r *http.Request) {
	log.Println("Serving linux games page")

	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	page := Page{"Linux Games"}

	err := index.ExecuteTemplate(w, "lgPage", page)
	if err != nil {
		log.Println(err) // http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

func mc(w http.ResponseWriter, r *http.Request) {
	log.Println("Serving Minecraft page")

	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	page := Page{"minecraft"}

	err := index.ExecuteTemplate(w, "mcPage", page)
	if err != nil {
		log.Println(err) // http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

func submission(w http.ResponseWriter, r *http.Request) {
	data := Answer{Title: "wrong answer headquarters", Flag: r.URL.Query().Get("flag"), Puzzle: r.URL.Query().Get("puzzle"), Result: false}
	log.Printf("For %s you submitted: %s\n", data.Puzzle, data.Flag)

	// --------------- form input validation ---------------
	if len(strings.TrimSpace(data.Flag)) == 0 {
		log.Println("User submitted blank value.")
		incorrect(w, r)
	}

	if !strings.ContainsAny(data.Flag, `!-_,.@/()=+?`) {
		log.Println("User submitted special character nonsense.")
		incorrect(w, r)
	}

	var emailRegexp = regexp.MustCompile(`^[^@]+@[^@]+\.[^@]+$`)
	if !emailRegexp.Match([]byte(data.Flag)) {
		log.Println("User submitted a possible email address.")
		incorrect(w, r)	
	}

	// --------------- check the submitted flag  ---------------
	if ((data.Puzzle == "puzzle1") && (data.Flag == "HAHASTRUGGLE")) {
		log.Println(w, "Correct answer!")
		log.Println("Serving correct page")
		data.Title = "YOU DID IT"
		data.Result = true
	} else if ((data.Puzzle == "puzzle2-1") && (data.Flag == "DELICIOUSTEARS")) {
		log.Println(w, "Correct answer!")
		log.Println("Serving correct page")
		data.Title = "YOU DID IT"
		data.Result = true
	} else if 	((data.Puzzle == "puzzle2-2") && (data.Flag == "COMPUTERSAREHARD")) {
		log.Println(w, "Correct answer!")
		log.Println("Serving correct page")
		data.Title = "YOU DID IT"
		data.Result = true
	} else if 	((data.Puzzle == "puzzle2-3") && (data.Flag == "ENJOYTHEPAIN")) {
		log.Println(w, "Correct answer!")
		log.Println("Serving correct page")
		data.Title = "YOU DID IT"
		data.Result = true
	} else if 	((data.Puzzle == "puzzle2-4") && (data.Flag == "ANTISOCIALNETWORKING")) {
		log.Println(w, "Correct answer!")
		log.Println("Serving correct page")
		data.Title = "YOU DID IT"
		data.Result = true
	} else {
		data.Result = false
	}

	if data.Result {
		correct(w, r)
	} else {
		incorrect(w, r)
	}
}

func correct(w http.ResponseWriter, r *http.Request) {
	log.Println("Serving CORRECT page")

	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	data := Answer{Title: "CORRECT ANSWER", Flag: r.URL.Query().Get("flag"), Puzzle: r.URL.Query().Get("puzzle"), Result: true}
	log.Println(w, "Correct answer!")
	log.Println("Serving correct page")

	data.Title = "YOU DID IT"
	data.Result = true

	err := index.ExecuteTemplate(w, "correctPage", data)
	if err != nil {
		log.Println(err) // http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

func incorrect(w http.ResponseWriter, r *http.Request) {
	log.Println("Serving WRONG page")
	
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
	
    data := Answer{Title: "WRONG ANSWER", Flag: r.URL.Query().Get("flag"), Puzzle: r.URL.Query().Get("puzzle"), Result: false}
	log.Println(w, "NOPE. Guess the struggle continues!")
	log.Println("Serving incorrect page")

	data.Title = "NOPE. Guess the struggle continues!"
	data.Result = false

	err := index.ExecuteTemplate(w, "incorrectPage", data)
	if err != nil {
		log.Println(err) // http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}


func puzzle1(w http.ResponseWriter, r *http.Request) {
	log.Println("Serving puzzle1 page")

	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	page := Page{Title: "puzzle1"}

	err := index.ExecuteTemplate(w, "puzzle1Page", page)
	if err != nil {
		log.Println(err) // http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

func puzzle2(w http.ResponseWriter, r *http.Request) {
	log.Println("Serving puzzle2 page")

	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	page := Page{Title: "puzzle2"}


	
	err := index.ExecuteTemplate(w, "puzzle2Page", page)
	if err != nil {
		log.Println(err) // http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

func puzzle3(w http.ResponseWriter, r *http.Request) {
	log.Println("Serving puzzle3 page")
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
	page := Page{Title: "puzzle3"}
	err := index.ExecuteTemplate(w, "puzzle3Page", page)
	if err != nil {
		log.Println(err) // http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

func puzzle4(w http.ResponseWriter, r *http.Request) {
	log.Println("Serving puzzle4 page")

	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	page := Page{Title: "puzzle4"}

	err := index.ExecuteTemplate(w, "puzzle4Page", page)
	if err != nil {
		log.Println(err) // http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

func puzzle5(w http.ResponseWriter, r *http.Request) {
	log.Println("Serving puzzle5 page")

	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	page := Page{Title: "puzzle5"}
	
	err := index.ExecuteTemplate(w, "puzzle5Page", page)
	if err != nil {
		log.Println(err) // http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}


/*
package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/", serveTemplate)

	log.Print("Listening on :3000...")
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func serveTemplate(w http.ResponseWriter, r *http.Request) {
	lp := filepath.Join("templates", "layout.html")
	fp := filepath.Join("templates", filepath.Clean(r.URL.Path))

	// Return a 404 if the template doesn't exist
	info, err := os.Stat(fp)
	if err != nil {
		if os.IsNotExist(err) {
			http.NotFound(w, r)
			return
		}
	}

	// Return a 404 if the request is for a directory
	if info.IsDir() {
		http.NotFound(w, r)
		return
	}

	tmpl, err := template.ParseFiles(lp, fp)
	if err != nil {
		// Log the detailed error
		log.Print(err.Error())
		// Return a generic "Internal Server Error" message
		http.Error(w, http.StatusText(500), 500)
		return
	}

	err = tmpl.ExecuteTemplate(w, "layout", nil)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, http.StatusText(500), 500)
	}
}

*/