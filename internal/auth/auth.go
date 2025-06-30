package auth

import (
	"crypto/rand"
	"encoding/base64"
	"html/template"
	"io"
	"log"
	"net/http"
	"math/rand" // For simple state generation

	"github.com/gorilla/securecookie"
	"os"

)

type TwitchUser struct { // TwitchUser represents the user information returned by the Twitch API.
	ID              string `json:"id"`
	Login           string `json:"login"`
	DisplayName     string `json:"display_name"`
	Type            string `json:"type"`
	BroadcasterType string `json:"broadcaster_type"`
	Description     string `json:"description"`
	ProfileImageURL string `json:"profile_image_url"`
	OfflineImageURL string `json:"offline_image_url"`
	ViewCount       int    `json:"view_count"`
	Email           string `json:"email"`
	CreatedAt       string `json:"created_at"`
}

var (
	clientID      = os.Getenv("TWITCH_CLIENT_ID")                                        // Get from environment
	clientSecret  = os.Getenv("TWITCH_CLIENT_SECRET")                                    // Get from environment
	redirectURL   = "https://www.bitsmasher.net:8080/twitch/callback"                    // Must match the registered redirect URI
	scopes        = []string{"user:read:email", "bits:read", "channel:read:redemptions", "chat:read"} // Desired scopes
	cookieHandler = securecookie.New(securecookie.GenerateRandomKey(64), securecookie.GenerateRandomKey(32))
)

func GenerateRandomState(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

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
