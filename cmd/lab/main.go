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
    fs := http.FileServer(http.Dir("./static/lab"))
    http.Handle("/static/lab/", http.StripPrefix("/static/lab/", fs))

    tmpls, err = template.ParseGlob(LayoutDir + "/*.tmpl")
    if err != nil {
        logging.Log_info(err.Error())
	panic(err)
    }

    http.HandleFunc("/", handler)
    http.HandleFunc("/lab", labPageHandler)
    http.HandleFunc("/minecraft", minecraftPageHandler)
    http.HandleFunc("/auth", authPageHandler)
    http.HandleFunc("/dst", dst)
    http.HandleFunc("/civilization2", civ2)
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

/*
func formatDate(t time.Time) string {
    return t.Format("2006-01-02")
}
*/
 
func minecraftPageHandler(w http.ResponseWriter, r *http.Request) {
        logging.Log_header("Serving Minecraft page")
        
        w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
        w.Header().Set("Pragma", "no-cache")
        w.Header().Set("Expires", "0")
        
        //username := auth.GetUserName(r)
        page := Page{"minecraft"}
        err := tmpls.ExecuteTemplate(w, "mcPage", page)
        if err != nil {
          // logging.CheckError(err) // http.Error(w, err.Error(), http.StatusInternalServerError)
	  logging.Log_info(err.Error())
	  http.Error(w, "Internal server error: Could not render minecraft page.", http.StatusInternalServerError)
	}
}

func labPageHandler(w http.ResponseWriter, r *http.Request) {
        logging.Log_header("Serving Lab page")
        
        w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
        w.Header().Set("Pragma", "no-cache")
        w.Header().Set("Expires", "0")
        
        //username := auth.GetUserName(r)
        page := Page{"lab"}
        
        err := tmpls.ExecuteTemplate(w, "labPage", page)
        if err != nil {
          logging.Log_info(err.Error())
          http.Error(w, "Internal server error: Could not render lab page.", http.StatusInternalServerError)                                                                                             
        }
}

func authPageHandler(w http.ResponseWriter, r *http.Request) {
        logging.Log_header("Serving auth page")
        
        w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
        w.Header().Set("Pragma", "no-cache")
        w.Header().Set("Expires", "0")
        
        //username := auth.GetUserName(r)
        page := Page{"auth"}
        
        err := tmpls.ExecuteTemplate(w, "labAuthPage", page)
        if err != nil {
          logging.Log_info(err.Error())
          http.Error(w, "Internal server error: Could not render auth page.", http.StatusInternalServerError)                                                                                             
        }
}

func civ2(w http.ResponseWriter, r *http.Request) {
        logging.Log_header("Serving civ2 page")

        w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
        w.Header().Set("Pragma", "no-cache")
        w.Header().Set("Expires", "0")

        page := Page{"Civilization ]["}

        err := tmpls.ExecuteTemplate(w, "civ2Page", page)
        if err != nil {
          logging.Log_info(err.Error())
          http.Error(w, "Internal server error: Could not render civ page.", http.StatusInternalServerError)                                 
        }

}

func dst(w http.ResponseWriter, r *http.Request) {
	logging.Log_header("Serving DST page")

	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	page := Page{"Do Not Starve"}

	err := tmpls.ExecuteTemplate(w, "dstPage", page)
        if err != nil {
          logging.Log_info(err.Error())
          http.Error(w, "Internal server error: Could not render dst page.", http.StatusInternalServerError)                                 
        }

}
