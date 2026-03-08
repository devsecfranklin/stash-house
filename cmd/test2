package main

import (
	"fmt"
	"net/http"
)

type Page struct {
	Title string
	Body  string
}

func loadPage(title string) (*Page, error) {
	// implement loading (file, DB, etc). Example stub:
	if title == "" {
		return nil, fmt.Errorf("empty title")
	}
	return &Page{Title: title, Body: "Page body for " + title}, nil
}

func configHandler(w http.ResponseWriter, r *http.Request) {
	const prefix = "/config/"
	if len(r.URL.Path) <= len(prefix) {
		http.Error(w, "missing title", http.StatusBadRequest)
		return
	}
	title := r.URL.Path[len(prefix):]

	p, err := loadPage(title)
	if err != nil {
		http.Error(w, "failed to load page: "+err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "<h1>%s</h1><div>%s</div>", p.Title, p.Body)
}

func main() {
	http.HandleFunc("/config/", configHandler)
	http.ListenAndServe(":8080", nil)
}
