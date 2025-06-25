package auth

import (
	"html/template"
	"log"
	"io"
	"crypto/rand"
	"encoding/base64"
	"net/http"

	"github.com/gorilla/securecookie"
)

var cookieHandler = securecookie.New(securecookie.GenerateRandomKey(64), securecookie.GenerateRandomKey(32))

func SignupPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// Retrieve signup form data.
		username := r.FormValue("username")
		password := r.FormValue("password")

		// Perform signup logic here (e.g., store user data in a database).
		// For simplicity, we'll just print the data for demonstration.
		log.Printf("New user signup: Username - %s, Password - %s\n", username, password)

		// Redirect to a welcome or login page after signup.
		http.Redirect(w, r, "/welcome", http.StatusSeeOther)
		return
	}

	// If not a POST request, serve the signup page template.
	tmpl, err := template.ParseFiles("templates/signup.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

// func LoginHandler(response http.ResponseWriter, request *http.Request) {
// 	name := request.FormValue("name")
// 	pass := request.FormValue("password")
// 	redirectTarget := "/"
// 	if name != "" && pass != "" {
// 		// .. check credentials ..
// 		setSession(name, response)
// 		redirectTarget = "/internal"
// 	}
// 	http.Redirect(response, request, redirectTarget, 302)
// }

func LogoutHandler(response http.ResponseWriter, request *http.Request) {
	clearSession(response)
	http.Redirect(response, request, "/", 302)
}

func SetSession(userName string, response http.ResponseWriter) {
	log.Println("cookies: SetSession()")
	value := map[string]string{
		"name": userName,
	}
	if encoded, err := cookieHandler.Encode("session", value); err == nil {
		cookie := &http.Cookie{
			Name:  "session",
			Value: encoded,
			Path:  "/",
		}
		log.Println("settings cookie....")
		http.SetCookie(response, cookie)
	}
}

func GetUserName(request *http.Request) (userName string) {
	if cookie, err := request.Cookie("session"); err == nil {
		cookieValue := make(map[string]string)
		if err = cookieHandler.Decode("session", cookie.Value, &cookieValue); err == nil {
			userName = cookieValue["name"]
		}
	}
	
	return userName
}

func clearSession(response http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:   "session",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(response, cookie)
}

func State(n int) (string, error) {
    data := make([]byte, n)
    if _, err := io.ReadFull(rand.Reader, data); err != nil {
        return "", err
    }
    return trimStringToFirstXRunes(base64.StdEncoding.EncodeToString(data), 16), nil
}

func trimStringToFirstXRunes(s string, x int) string {
	// Convert the string to a slice of runes
	runes := []rune(s)

	// If x is greater than or equal to the number of runes, return the original string
	if x >= len(runes) {
		return s
	}

	// Return the string created from the first x runes
	return string(runes[:x])
}
