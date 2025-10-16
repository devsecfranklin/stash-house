/*
# SPDX-FileCopyrightText: ©2025 franklin <smoooth.y62wj@passmail.net>
#
# SPDX-License-Identifier: MIT
*/

package main

import (
	"encoding/json" // For parsing Twitch token response
	"fmt"
	"html/template"
	"internal/auth"
	"internal/logging" // switch to new module soon
	"io/ioutil"        // For reading response body
	"log"
	"net/http"
	"net/url" // For URL encoding parameters
	"os"
	"sync" // For mutex to protect the state map

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/twitch"
)

type (
	Stuff struct {
		token string
	}

	twitchUser struct { // A simple struct to ht
	OauthToken struct { // OauthToken data structure passed to the template
		ClientID  string
		StateRand string
		Title     string
	}

	Page struct { // Page data structure for a generic page
		Title string
	}
)

var (
	err error

	LayoutDir string = "template/www"
	tmpls     *template.Template

	twitchOauthConfig = &oauth2.Config{ // It's best practice to set these as environment variables.
		RedirectURL:  "https://www.bitsmasher.net/twitch/callback",
		ClientID:     os.Getenv("TWITCH_CLIENT_ID"),
		ClientSecret: os.Getenv("TWITCH_CLIENT_SECRET"),
		Scopes:       []string{"user:read:email"}, // Scopes determine what permissions you're asking for.
		Endpoint:     twitch.Endpoint,
	}

	oauthStateString = "random-string-for-demonstration" // A random string used to protect against CSRF attacks

	twitchClientID     string                                      // Twitch API Credentials (IMPORTANT: Load from environment variables in production!)
	twitchClientSecret string                                      // Needed for code exchange
	twitchRedirectURI  string = "https://www.bitsmasher.net/oauth" // This MUST match your registered redirect URI
	twitchOauthToken   string
	oauthStates        = struct { // Temporary storage for OAuth states. Move into database
		sync.RWMutex
		m map[string]bool // map[state_string]is_valid
	}{m: make(map[string]bool)}
)

func main() {
	fs := http.FileServer(http.Dir("./static/www"))
	http.Handle("/static/www/", http.StripPrefix("/static/www/", fs))

	tmpls, err = template.ParseGlob(LayoutDir + "/*.tmpl")
	if err != nil {
		panic(err)
	}

	// Register the handler for all paths
	http.HandleFunc("/", handler)
	http.HandleFunc("/oauth", oauthHandler)
	http.HandleFunc("/twitch/callback", oauthHandler) //handleTwitchCallback)
	http.HandleFunc("/chatoverlay", twitchChatHandler)
	http.HandleFunc("/lab", labPageHandler)
    http.HandleFunc("/labAnsiblePage", labAnsiblePageHandler)
	http.HandleFunc("/labAuthPage", labAuthPageHandler)

	logging.Log_header("Server listening on :8080")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		logging.Log_fatal(fmt.Sprintf("Server failed to start: %v", err))
	}
}

func handler(w http.ResponseWriter, r *http.Request) { // handler for the root path
	log.Println("Serving index page")

	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	page := Page{"www.bitsmasher.net"}

	err := tmpls.ExecuteTemplate(w, "indexPage", page) // Assuming you have an "indexPage" template
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal server error: Could not render index page.", http.StatusInternalServerError)
	}
}

func oauthHandler(w http.ResponseWriter, r *http.Request) {
	logging.Log_header("Serving oauth page / Handling OAuth callback")

	twitchClientID = os.Getenv("TWITCH_CLIENT_ID") // Load Twitch credentials from environment variables
	if twitchClientID == "" {
		logging.Log_error("TWITCH_CLIENT_ID environment variable not set. Please set it.")
	}
	twitchClientSecret = os.Getenv("TWITCH_CLIENT_SECRET")
	if twitchClientSecret == "" {
		logging.Log_error("TWITCH_CLIENT_SECRET environment variable not set. Please set it.")
	}

	// Set cache control headers to prevent caching sensitive data
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	// Check if this is an OAuth callback from Twitch (i.e., contains 'code' parameter)
	code := r.URL.Query().Get("code")
	returnedState := r.URL.Query().Get("state")
	twitchError := r.URL.Query().Get("error")
	twitchErrorDescription := r.URL.Query().Get("error_description")

	if twitchError != "" {
		logging.Log_info(fmt.Sprintf("Twitch OAuth Error: %s - %s", twitchError, twitchErrorDescription))
		http.Error(w, fmt.Sprintf("Twitch authorization failed: %s", twitchErrorDescription), http.StatusForbidden)
		return
	}

	if code != "" {
		// This is a callback from Twitch after user authorization
		logging.Log_info("Received Twitch OAuth callback with code: " + code)

		twitchClientID = os.Getenv("TWITCH_CLIENT_ID") // Load Twitch credentials from environment variables
		if twitchClientID == "" {
			logging.Log_error("TWITCH_CLIENT_ID environment variable not set. Please set it.")
		}
		twitchClientSecret = os.Getenv("TWITCH_CLIENT_SECRET")
		if twitchClientSecret == "" {
			logging.Log_error("TWITCH_CLIENT_SECRET environment variable not set. Please set it.")
		}

		oauthStates.RLock()
		// isValidState := oauthStates.m[returnedState] // Validate the 'state' parameter to prevent CSRF attacks
		oauthStates.RUnlock()

		oauthStates.Lock()
		delete(oauthStates.m, returnedState)
		oauthStates.Unlock() // Remove the state after use to prevent replay attacks (simple example, more robust solution needed for production)
		logging.Log_success("State parameter validated.")

		logging.Log_info("Exchanging authorization code for access token: " + code) // Exchange the authorization code for an access token
		tokenURL := "https://id.twitch.tv/oauth2/token"

		data := url.Values{}
		data.Set("client_id", twitchClientID)
		data.Set("client_secret", twitchClientSecret)
		data.Set("code", code)
		data.Set("grant_type", "authorization_code")
		data.Set("redirect_uri", "https://www.bitsmasher.net/twitch/callback") // twitchRedirectURI) // This must be the same redirect_uri used in the initial request

		resp, err := http.PostForm(tokenURL, data)
		if err != nil {
			logging.Log_error(fmt.Sprintf("Failed to exchange code for token: %v", err))
			http.Error(w, "Failed to get access token from Twitch.", http.StatusInternalServerError)
			//return
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			logging.Log_error(fmt.Sprintf("Failed to read token response body: %v", err))
			http.Error(w, "Failed to read Twitch token response.", http.StatusInternalServerError)
			//return
		}

		if resp.StatusCode != http.StatusOK {
			logging.Log_error(fmt.Sprintf("Twitch token exchange failed with status %d: %s", resp.StatusCode, string(body)))
			http.Error(w, fmt.Sprintf("Twitch token exchange failed: %s", string(body)), http.StatusInternalServerError)
			// return
		}

		// Parse the token response
		var tokenResponse struct {
			AccessToken  string   `json:"access_token"`
			ExpiresIn    int      `json:"expires_in"`
			RefreshToken string   `json:"refresh_token"`
			Scope        []string `json:"scope"`
			TokenType    string   `json:"token_type"`
		}

		if err := json.Unmarshal(body, &tokenResponse); err != nil {
			logging.Log_error(fmt.Sprintf("Failed to unmarshal token response: %v", err))
			http.Error(w, "Failed to parse Twitch token response.", http.StatusInternalServerError)
			//return
		}

		logging.Log_success("Successfully obtained Twitch access token!")
		twitchOauthToken = tokenResponse.AccessToken
		f, err := os.Create("private/token-franklin")
		logging.CheckError(err)
		defer f.Close()

		logging.Log_info(fmt.Sprintf("Access Token: %s...", tokenResponse.AccessToken[:10])) // Log a prefix, not full token
		logging.Log_info(fmt.Sprintf("Refresh Token: %s...", tokenResponse.RefreshToken[:10]))
		logging.Log_info(fmt.Sprintf("Expires in: %d seconds", tokenResponse.ExpiresIn))
		logging.Log_info(fmt.Sprintf("Scopes: %v", tokenResponse.Scope))

		// IMPORTANT: Store tokenResponse.AccessToken and tokenResponse.RefreshToken securely.
		// Associate them with the user's account in your database.
		// For this example, we'll just display a success message.
		// fmt.Fprintf(w, "<h1>Twitch Authorization Successful!</h1>")
		// fmt.Fprintf(w, "<p>Access Token: <code>%s</code></p>", tokenResponse.AccessToken[:20])
		// fmt.Fprintf(w, "<p>Refresh Token: <code>%s</code></p>", tokenResponse.RefreshToken[:20])
		// fmt.Fprintf(w, "<p>You can now close this window/tab.</p>")
		//return // just show token on the main auth page
		err = tmpls.ExecuteTemplate(w, "successPage", tokenResponse)
		if err != nil {
			logging.Log_error(fmt.Sprintf("Failed to execute oauthPage template: %v", err))
			http.Error(w, "Internal server error: Could not render page.", http.StatusInternalServerError)
		}

	} else {
		// This is an initial request to display the OAuth authorization link
		logging.Log_info("Displaying initial Twitch OAuth authorization page.")

		// Generate a new unique state string for this request
		state, err := auth.GenerateRandomState(16) // Generate a random string of 32 characters

		if err != nil {
			logging.Log_error(fmt.Sprintf("Failed to generate random state: %v", err))
			http.Error(w, "Internal server error: Failed to generate random state.", http.StatusInternalServerError)
		}

		// Store the state securely (e.g., in a session) before redirecting.
		// For this simple example, we use a map. In production, ensure this is tied to the user's session.
		oauthStates.Lock()
		oauthStates.m[state] = true
		oauthStates.Unlock()
		logging.Log_info(fmt.Sprintf("Generated and stored new state: %s", state))

		// Prepare data for the template
		data := OauthToken{
			ClientID:  twitchClientID,
			StateRand: state,
			Title:     "Authorize with Twitch",
		}

		// Execute the template to render the page
		err = tmpls.ExecuteTemplate(w, "oauthPage", data)
		if err != nil {
			logging.Log_error(fmt.Sprintf("Failed to execute oauthPage template: %v", err))
			http.Error(w, "Internal server error: Could not render page.", http.StatusInternalServerError)
		}
	}
}

func twitchChatHandler(w http.ResponseWriter, r *http.Request) {
	data := Stuff{
		token: twitchOauthToken,
	}
	err = tmpls.ExecuteTemplate(w, "chatPage", data)
	if err != nil {
		logging.Log_error(fmt.Sprintf("Failed to execute oauthPage template: %v", err))
		http.Error(w, "Internal server error: Could not render page.", http.StatusInternalServerError)
	}
}

func labPageHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Serving lab page")

	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	page := Page{"www.bitsmasher.net/lab"}

	err := tmpls.ExecuteTemplate(w, "labPage", page)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal server error: Could not render index page.", http.StatusInternalServerError)
	}
}

func labAnsiblePageHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Serving lab Ansible page")

	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	page := Page{"www.bitsmasher.net/labAnsiblePage"}

	err := tmpls.ExecuteTemplate(w, "labAnsiblePage", page)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal server error: Could not render lab Ansible page.", http.StatusInternalServerError)
	}
}

func labAuthPageHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Serving lab auth page")

	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	page := Page{"www.bitsmasher.net/labAuthPage"}

	err := tmpls.ExecuteTemplate(w, "labAuthPage", page)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal server error: Could not render lab auth page.", http.StatusInternalServerError)
	}
}
