/*
# SPDX-FileCopyrightText: ©2025 franklin <smoooth.y62wj@passmail.net>
#
# SPDX-License-Identifier: MIT
*/
package main

import (
    "html/template"
    "net/http"
    "internal/logging"
	"fmt"
)

var (
	err error

	LayoutDir string = "template/lab"
	tmpls     *template.Template
)

type (
	Page struct { // Page data structure for a generic page
		Title string
	}
)

func main() {
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/*", http.StripPrefix("/static/", fs))

	tmpls, err = template.ParseGlob(LayoutDir + "/*.tmpl")
	if err != nil {
		logging.Log_info(err.Error())
		panic(err)
	}

	http.HandleFunc("/", handler)
        // http.HandleFunc("/lab", labHandler)
        http.HandleFunc("/minecraft", minecraftPageHandler)


	logging.Log_header("Server listening on :8080")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		logging.Log_fatal(fmt.Sprintf("Server failed to start: %v", err))
	}
}

func handler(w http.ResponseWriter, r *http.Request) { // handler for the root path
	logging.Log_info("Serving index page")

	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	page := Page{"lab.bitsmasher.net"}

	err := tmpls.ExecuteTemplate(w, "indexPage", page) // Assuming you have an "indexPage" template
	if err != nil {
		logging.Log_info(err.Error())
		http.Error(w, "Internal server error: Could not render index page.", http.StatusInternalServerError)
	}
}

        
func minecraftPageHandler(w http.ResponseWriter, r *http.Request) {
        logging.Log_header("Serving Minecraft page")
        
        w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
        w.Header().Set("Pragma", "no-cache")
        w.Header().Set("Expires", "0")
        
        //username := auth.GetUserName(r)
        page := Page{"minecraft"}
        
        err := index.ExecuteTemplate(w, "mcPage", page)
        if err != nil {
                // logging.CheckError(err) // http.Error(w, err.Error(), http.StatusInternalServerError)
		logging.Log_info(err.Error())
		http.Error(w, "Internal server error: Could not render lab auth page.", http.StatusInternalServerError)
	}
}
