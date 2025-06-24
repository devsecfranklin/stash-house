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

var LayoutDir string = "template/www"
var index *template.Template

type OauthToken struct {
	clientID string
	channelName string
	secret string
}

type Page struct {
	Title string
}

func oauthHandler(w http.ResponseWriter, r *http.Request) {
	log_header("Serving index page")

	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
	//data := OauthToken{clientID: r.URL.Query().Get("clientID"), secret: r.URL.Query().Get("secret")}

	page := Page{"oauth"}
	//log.Printf("For %s you submitted: %s\n", data.clientID, data.secret)

	err := index.ExecuteTemplate(w, "oauthPage", page)
	if err != nil {
		log.Println(err) // http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {

	// Serve static files
	fs := http.FileServer(http.Dir("./static"))
	router.Handle("/static/*", http.StripPrefix("/static/", fs))
	// Register the handler for all paths
	http.HandleFunc("/", answerHandler)
	http.HandleFunc("/oauth", oauthHandler)

	// http.Handle("/", http.FileServer(http.Dir("../../web/static/"))) // this is to server static web pages

	log_header("Server listening on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}


const (
	LRED   = "\033[1;31m"
	LGREEN = "\033[1;32m"
	LBLUE  = "\033[1;34m"
	LPURP  = "\033[1;35m"
	NC     = "\033[0m" // No Color
)

func log_header(msg string) {
	fmt.Printf("\n%s# --- %s %s\n", LPURP, msg, NC)
}

func log_info(msg string) {
	fmt.Printf("%s%s%s\n", LBLUE, msg, NC)
}

func log_success(msg string) {
	fmt.Printf("%s%s%s\n", LGREEN, msg, NC)
}

func log_error(msg string) {
	fmt.Printf("%sERROR: %s%s\n", LRED, msg, NC)
	os.Exit(1)
}