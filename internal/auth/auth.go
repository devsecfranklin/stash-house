package auth

import (
	"crypto/rand"
	"encoding/base64"
	"html/template"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/securecookie"
	"os"

	"context"
	"encoding/json"
	"fmt"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/twitch"

)

// Twitch configuration. You'll need to get these from the Twitch Developer Console.
// Replace with your actual values.
var (
	clientID      = os.Getenv("TWITCH_CLIENT_ID")                                        // Get from environment
	clientSecret  = os.Getenv("TWITCH_CLIENT_SECRET")                                    // Get from environment
	redirectURL   = "https://www.bitsmasher.net:8080/twitch/callback"                    // Must match the registered redirect URI
	scopes        = []string{"user:read:email", "bits:read", "channel:read:redemptions"} // Desired scopes
	cookieHandler = securecookie.New(securecookie.GenerateRandomKey(64), securecookie.GenerateRandomKey(32))
)

func requestTwitchOauth2() {
	if clientID == "" || clientSecret == "" {
		log.Fatal("TWITCH_CLIENT_ID and TWITCH_CLIENT_SECRET environment variables must be set.")
	}

	// OAuth2 configuration for Twitch
	conf := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Scopes:       scopes,
		Endpoint:     twitch.Endpoint,
	}

	// State is used to prevent CSRF attacks.  Generate a random state string and store it (e.g., in a session).
	state := "your-random-state" // Replace with a securely generated random string

	// Route to handle the authorization request
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Construct the authorization URL
		authURL := conf.AuthCodeURL(state, oauth2.AccessTypeOffline)
		fmt.Printf("Authorization URL: %s\n", authURL) //For testing
		// Redirect the user to the Twitch authorization page
		http.Redirect(w, r, authURL, http.StatusFound)
	})

	// Route to handle the callback from Twitch after authorization
	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		// Check for errors in the callback
		if err := r.URL.Query().Get("error"); err != "" {
			fmt.Fprintf(w, "Error from Twitch: %s\n", err)
			return
		}

		// Verify the state parameter (crucial for security)
		if r.URL.Query().Get("state") != state {
			fmt.Fprintln(w, "Invalid state parameter")
			return
		}

		// Get the authorization code from the callback
		code := r.URL.Query().Get("code")

		// Exchange the authorization code for an access token
		token, err := conf.Exchange(context.Background(), code)
		if err != nil {
			fmt.Fprintf(w, "Error exchanging code for token: %s\n", err)
			return
		}

		// Now you have the access token.  Store it securely (e.g., in a database).
		// You can use the token to make API calls to Twitch.
		fmt.Fprintf(w, "Access Token: %s\n", token.AccessToken)
		fmt.Fprintf(w, "Refresh Token: %s\n", token.RefreshToken)
		fmt.Fprintf(w, "Token Type: %s\n", token.TokenType)

		// Example: Get user information
		user, err := getUserInfo(token.AccessToken)
		if err != nil {
			fmt.Fprintf(w, "Error getting user info: %s\n", err)
			return
		}
		fmt.Fprintf(w, "User Info: %+v\n", user)

	})

	// Start the HTTP server
	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// TwitchUser represents the user information returned by the Twitch API.
type TwitchUser struct {
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

// getUserInfo retrieves user information from the Twitch API.
func getUserInfo(accessToken string) (*TwitchUser, error) {
	// Twitch API endpoint for user information.
	apiURL := "https://api.twitch.tv/helix/users"

	// Create a new HTTP client.
	client := &http.Client{}

	// Create a new GET request.
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	// Add the required headers (Authorization and Client-ID).
	req.Header.Add("Authorization", "Bearer "+accessToken)
	req.Header.Add("Client-Id", clientID) // Use the global clientID

	// Send the request.
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	// Check the response status code.
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("received non-200 OK status code: %d, body: %s", resp.StatusCode, string(body))
	}

	// Decode the JSON response.
	var data struct {
		Data []TwitchUser `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("error decoding JSON: %w", err)
	}

	// Check if any user data was returned.
	if len(data.Data) == 0 {
		return nil, fmt.Errorf("no user data found")
	}

	// Return the first user in the data slice.
	return &data.Data[0], nil
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
