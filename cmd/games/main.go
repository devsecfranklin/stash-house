/*
# SPDX-FileCopyrightText: © 2020-2025 franklin <franklin@bitsmasher.net>
#
# SPDX-License-Identifier: MIT
*/
package main

import (
	"html/template"
	"net/http"
	"internal/auth"
	"internal/database"
	"internal/logging"
)

var LayoutDir string = "template/games"
var index *template.Template

type Page struct {
	Title string
	Username string
}

type Answer struct {
	Title  string
	CSS    string
	Flag   string
	Puzzle string
	Result bool
}

type User struct {
	id     int
	name   string
	email  string
	password string
}

func main() {
	var err error
	var table_status bool
	var db_name, table_name string

	logging.Log_header("Server starting...")

	db_name = "games"
	table_name = "users"
	table_status = database.CheckDB(table_name, db_name)
	if !table_status {
		database.CrudDB()
	}

	//  --- Make static files available ------------------------------------------
	fs := http.FileServer(http.Dir("./static/games"))
	http.Handle("/static/games/", http.StripPrefix("/static/games/", fs))

	index, err = template.ParseGlob(LayoutDir + "/*.tmpl")
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/", handler)
	http.HandleFunc("/civilization2", civ2)
	http.HandleFunc("/dst", dst)
	http.HandleFunc("/failure", failure)
	http.HandleFunc("/linux_games", lg)
	http.HandleFunc("/login", LoginPage)
	http.HandleFunc("/minecraft", mc)
	http.HandleFunc("/mcadmin", mcAdmin)
	http.HandleFunc("/mchistory", mcHist)
	http.HandleFunc("/scoreboard", scoreboard)
	http.HandleFunc("/submission", submission)
	//http.HandleFunc("/puzzle0", puzzle0)
	http.HandleFunc("/puzzle1", puzzle1)
	http.HandleFunc("/puzzle2", puzzle2)
	http.HandleFunc("/puzzle3", puzzle3)
	//http.HandleFunc("/puzzle4", puzzle4)
	//http.HandleFunc("/puzzle5", puzzle5)
	http.HandleFunc("/welcome", WelcomePage)

	logging.Log_info("Use browser to access this host on :8080/")
	http.ListenAndServe(":8080", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
  logging.Log_header("Running main handler")

	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	//var cookieStatus bool = cookies.getCookie(w, r)
	//if !cookieStatus {cookies.setCookie(w, r) }

	username := auth.GetUserName(r)
	page := Page{"games.bitsmasher.net", username}

	err := index.ExecuteTemplate(w, "indexPage", page)
	if err != nil {
	  logging.CheckError(err) // http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func LoginPage(w http.ResponseWriter, r *http.Request) {
    logging.Log_header("Serving login page")

	var username string

	redirectTemplate := "LoginPage"

	if r.Method == http.MethodPost {
		username = r.FormValue("usr")
		password := r.FormValue("pwd")

		err := database.ReadUserDB("users", "games", username, password) 
		if err != nil {
			logging.CheckError(err) // http.Error(w, err.Error(), http.StatusInternalServerError)
			redirectTemplate = "FailurePage"
		} else {	
			auth.SetSession(username, w) // set session cookie
			redirectTemplate = "WelcomePage"
		}
	}
	
        //page := Page{"games.bitsmasher.net", username}
	err := index.ExecuteTemplate(w, redirectTemplate, LoginPage)
	
	if err != nil {
		logging.CheckError(err) // http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func WelcomePage(w http.ResponseWriter, r *http.Request) {
	//logging.Log_header(w, "Welcome, you have successfully logged in!")

	username := auth.GetUserName(r)

	logging.Log_info("Serving welcome page for " + username )

	page := Page{"games.bitsmasher.net", username}

	err := index.ExecuteTemplate(w, "WelcomePage", page)

	if err != nil {
		logging.CheckError(err) // http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

func civ2(w http.ResponseWriter, r *http.Request) {
	logging.Log_header("Serving civ2 page")

	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	username := auth.GetUserName(r)
	page := Page{"Civilization ][", username}

	err := index.ExecuteTemplate(w, "civ2Page", page)
	if err != nil {
		logging.CheckError(err) // http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

func dst(w http.ResponseWriter, r *http.Request) {
	logging.Log_header("Serving DST page")

	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

    username := auth.GetUserName(r)
	page := Page{"Do Not Starve", username}

	err := index.ExecuteTemplate(w, "dstPage", page)
	if err != nil {
		logging.CheckError(err) // http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

func failure(w http.ResponseWriter, r *http.Request) {
	logging.Log_header("Serving failure page")

	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	username := auth.GetUserName(r)
	page := Page{"failure", username}

	err := index.ExecuteTemplate(w, "FailurePage", page)
	if err != nil {
		logging.CheckError(err) // http.Error(w, err.Error(), http.StatusInternalServerError)
	}	
}

func lg(w http.ResponseWriter, r *http.Request) {
	logging.Log_header("Serving linux games page")

	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	username := auth.GetUserName(r)
	page := Page{"Linux Games", username}

	err := index.ExecuteTemplate(w, "lgPage", page)
	if err != nil {
		logging.CheckError(err) // http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

func mc(w http.ResponseWriter, r *http.Request) {
	logging.Log_header("Serving minecraft page")
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	username := auth.GetUserName(r)
	page := Page{"minecraft.bitsmasher.net", username}

	err := index.ExecuteTemplate(w, "mcPage", page)
	if err != nil {
		logging.CheckError(err) // http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func mcAdmin(w http.ResponseWriter, r *http.Request) {
	logging.Log_header("Serving minecraft admin page")
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
	username := auth.GetUserName(r)
	page := Page{"minecraft admin", username}
	err := index.ExecuteTemplate(w, "mcAdmin", page)
	if err != nil {
		logging.CheckError(err) // http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func mcHist(w http.ResponseWriter, r *http.Request) {
	logging.Log_header("Serving minecraft history page")
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
	username := auth.GetUserName(r)
	page := Page{"minecraft history", username}
	err := index.ExecuteTemplate(w, "mcHistory", page)
	if err != nil {
		logging.CheckError(err) // http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func scoreboard(w http.ResponseWriter, r *http.Request) {
	logging.Log_header("Serving scoreboard games page")

	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	username := auth.GetUserName(r)
	page := Page{"minecraft", username}

	err := index.ExecuteTemplate(w, "scoreboardPage", page)
	if err != nil {
		logging.CheckError(err) // http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func submission(w http.ResponseWriter, r *http.Request) {
	data := Answer{Title: "wrong answer headquarters", Flag: r.URL.Query().Get("flag"), Puzzle: r.URL.Query().Get("puzzle"), Result: false}
	msg := "For " + data.Puzzle + " you submitted: " + data.Flag
	logging.Log_header(msg)

	// --------------- check the submitted flag  ---------------
	if (data.Puzzle == "puzzle1") && (data.Flag == "HAHASTRUGGLE") {
		logging.Log_info("Correct answer!")
		logging.Log_success("Serving correct page")
		data.Title = "YOU DID IT"
		data.Result = true
	} else if (data.Puzzle == "puzzle2-1") && (data.Flag == "DELICIOUSTEARS") {
		logging.Log_info("Correct answer!")
		logging.Log_success("Serving correct page")
		data.Title = "YOU DID IT"
		data.Result = true
	} else if (data.Puzzle == "puzzle2-2") && (data.Flag == "COMPUTERSAREHARD") {
		logging.Log_info("Correct answer!")
		logging.Log_success("Serving correct page")
		data.Title = "YOU DID IT"
		data.Result = true
	} else if (data.Puzzle == "puzzle2-3") && (data.Flag == "ENJOYTHEPAIN") {
		logging.Log_info("Correct answer!")
		logging.Log_success("Serving correct page")
		data.Title = "YOU DID IT"
		data.Result = true
	} else if (data.Puzzle == "puzzle2-4") && (data.Flag == "ANTISOCIALNETWORKING") {
		logging.Log_info("Correct answer!")
		logging.Log_success("Serving correct page")
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
	logging.Log_header("Serving CORRECT page")

	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	data := Answer{Title: "CORRECT ANSWER", Flag: r.URL.Query().Get("flag"), Puzzle: r.URL.Query().Get("puzzle"), Result: true}
	logging.Log_info("Correct answer!")
	logging.Log_success("Serving correct page")

	data.Title = "YOU DID IT"
	data.Result = true

	err := index.ExecuteTemplate(w, "correctPage", data)
	if err != nil {
		logging.CheckError(err) // http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

func incorrect(w http.ResponseWriter, r *http.Request) {
	logging.Log_header("Serving WRONG page")

	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	data := Answer{Title: "WRONG ANSWER", Flag: r.URL.Query().Get("flag"), Puzzle: r.URL.Query().Get("puzzle"), Result: false}
	logging.Log_info("NOPE. Guess the struggle continues!")
	logging.Log_info("Serving incorrect page")

	data.Title = "NOPE. Guess the struggle continues!"
	data.Result = false

	err := index.ExecuteTemplate(w, "incorrectPage", data)
	if err != nil {
		logging.CheckError(err) // http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

func puzzle1(w http.ResponseWriter, r *http.Request) {
	logging.Log_header("Serving puzzle1 page")

	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	username := auth.GetUserName(r)
	page := Page{"puzzle1", username}

	err := index.ExecuteTemplate(w, "puzzle1Page", page)
	if err != nil {
		logging.CheckError(err) // http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

func puzzle2(w http.ResponseWriter, r *http.Request) {
	logging.Log_header("Serving puzzle2 page")

	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	username := auth.GetUserName(r)
	page := Page{"puzzle2", username}

	err := index.ExecuteTemplate(w, "puzzle2Page", page)
	if err != nil {
		logging.CheckError(err) // http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

func puzzle3(w http.ResponseWriter, r *http.Request) {
	logging.Log_header("Serving puzzle3 page")

	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	username := auth.GetUserName(r)
	page := Page{"puzzle3", username}

	err := index.ExecuteTemplate(w, "puzzle3Page", page)

	if err != nil {
		logging.CheckError(err) // http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

/*
func puzzle4(w http.ResponseWriter, r *http.Request) {
	logging.Log_header("Serving puzzle4 page")

	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	username := auth.GetUserName(r)
	page := Page{"puzzle4", username}

	err := index.ExecuteTemplate(w, "puzzle4Page", page)
	if err != nil {
		logging.CheckError(err) // http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

func puzzle5(w http.ResponseWriter, r *http.Request) {
	logging.Log_header("Serving puzzle5 page")

	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	username := auth.GetUserName(r)
	page := Page{"puzzle5", username}

	err := index.ExecuteTemplate(w, "puzzle5Page", page)
	if err != nil {
		logging.CheckError(err) // http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}
*/
