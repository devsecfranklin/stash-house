package main

import (
	"encoding/json" // For parsing Twitch token response
	"fmt"
	"html/template"
	"io/ioutil" // For reading response body
	"log"
	"net/http"
	"net/url"   // For URL encoding parameters
	"os"
	"sync"      // For mutex to protect the state map
	"time"      // For state expiration/cleanup (simple example)
	"math/rand" // For simple state generation
)

// In a real application, you'd use a more robust random string generator,
// potentially from a crypto package. This is a simple example.
func generateRandomState(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

var LayoutDir string = "template/www"
var tmpls *template.Template // Renamed from 'index' for clarity when multiple templates are loaded

// Twitch API Credentials (IMPORTANT: Load from environment variables in production!)
var (
	twitchClientID     string
	twitchClientSecret string // Needed for code exchange
	twitchRedirectURI  string = "https://www.bitsmasher.net/oauth" // This MUST match your registered redirect URI
)

// Temporary storage for OAuth states.
// In a real application, use a secure, persistent storage (e.g., database, secure session store)
// and associate it with a user session. This simple map is for demonstration ONLY.
var oauthStates = struct {
	sync.RWMutex
	m map[string]bool // map[state_string]is_valid
}{m: make(map[string]bool)}

// OauthToken data structure passed to the template
type OauthToken struct {
	ClientID string
	StateRand string
	Title     string
}

// Page data structure for a generic page
type Page struct {
	Title string
}

func init() {
	// Seed random number generator for state generation
	rand.Seed(time.Now().UnixNano())

	// Load Twitch credentials from environment variables
	twitchClientID = os.Getenv("TWITCH_CLIENT_ID")
	if twitchClientID == "" {
		log_error("TWITCH_CLIENT_ID environment variable not set. Please set it.")
	}
	twitchClientSecret = os.Getenv("TWITCH_TOKEN")
	if twitchClientSecret == "" {
		log_error("TWITCH_TOKEN environment variable not set. Please set it.")
	}

	// Load templates once at startup
	var err error
	tmpls, err = template.ParseGlob(LayoutDir + "/*.tmpl")
	if err != nil {
		log_error(fmt.Sprintf("Failed to parse templates: %v", err))
	}
}

// oauthHandler handles both the initial request for the OAuth page
// and the callback from Twitch after authorization.
func oauthHandler(w http.ResponseWriter, r *http.Request) {
	log_header("Serving oauth page / Handling OAuth callback")

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
		log_info(fmt.Sprintf("Twitch OAuth Error: %s - %s", twitchError, twitchErrorDescription))
		http.Error(w, fmt.Sprintf("Twitch authorization failed: %s", twitchErrorDescription), http.StatusForbidden)
		return
	}

	if code != "" {
		// This is a callback from Twitch after user authorization
		log_info("Received Twitch OAuth callback with code.")

		// 1. Validate the 'state' parameter to prevent CSRF attacks
		oauthStates.RLock()
		isValidState := oauthStates.m[returnedState]
		oauthStates.RUnlock()

		if !isValidState {
			log_error(fmt.Sprintf("Invalid or missing state parameter from Twitch callback: %s", returnedState))
			http.Error(w, "Invalid OAuth state. Possible CSRF attack or expired request.", http.StatusBadRequest)
			return
		}

		// Remove the state after use to prevent replay attacks (simple example, more robust solution needed for production)
		oauthStates.Lock()
		delete(oauthStates.m, returnedState)
		oauthStates.Unlock()
		log_success("State parameter validated.")

		// 2. Exchange the authorization code for an access token
		log_info("Exchanging authorization code for access token...")
		tokenURL := "https://id.twitch.tv/oauth2/token"
		
		// Use url.Values to properly encode form data
		data := url.Values{}
		data.Set("client_id", twitchClientID)
		data.Set("client_secret", twitchClientSecret)
		data.Set("code", code)
		data.Set("grant_type", "authorization_code")
		data.Set("redirect_uri", twitchRedirectURI) // This must be the same redirect_uri used in the initial request

		resp, err := http.PostForm(tokenURL, data)
		if err != nil {
			log_error(fmt.Sprintf("Failed to exchange code for token: %v", err))
			http.Error(w, "Failed to get access token from Twitch.", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log_error(fmt.Sprintf("Failed to read token response body: %v", err))
			http.Error(w, "Failed to read Twitch token response.", http.StatusInternalServerError)
			return
		}

		if resp.StatusCode != http.StatusOK {
			log_error(fmt.Sprintf("Twitch token exchange failed with status %d: %s", resp.StatusCode, string(body)))
			http.Error(w, fmt.Sprintf("Twitch token exchange failed: %s", string(body)), http.StatusInternalServerError)
			return
		}

		// Parse the token response
		var tokenResponse struct {
			AccessToken  string `json:"access_token"`
			ExpiresIn    int    `json:"expires_in"`
			RefreshToken string `json:"refresh_token"`
			Scope        []string `json:"scope"`
			TokenType    string `json:"token_type"`
		}
		if err := json.Unmarshal(body, &tokenResponse); err != nil {
			log_error(fmt.Sprintf("Failed to unmarshal token response: %v", err))
			http.Error(w, "Failed to parse Twitch token response.", http.StatusInternalServerError)
			return
		}

		log_success("Successfully obtained Twitch access token!")
		log_info(fmt.Sprintf("Access Token: %s...", tokenResponse.AccessToken[:10])) // Log a prefix, not full token
		log_info(fmt.Sprintf("Refresh Token: %s...", tokenResponse.RefreshToken[:10]))
		log_info(fmt.Sprintf("Expires in: %d seconds", tokenResponse.ExpiresIn))
		log_info(fmt.Sprintf("Scopes: %v", tokenResponse.Scope))

		// IMPORTANT: Store tokenResponse.AccessToken and tokenResponse.RefreshToken securely.
		// Associate them with the user's account in your database.
		// For this example, we'll just display a success message.
		fmt.Fprintf(w, "<h1>Twitch Authorization Successful!</h1>")
		fmt.Fprintf(w, "<p>Access Token: <code>%s...</code></p>", tokenResponse.AccessToken[:20])
		fmt.Fprintf(w, "<p>Refresh Token: <code>%s...</code></p>", tokenResponse.RefreshToken[:20])
		fmt.Fprintf(w, "<p>You can now close this window/tab.</p>")
		return

	} else {
		// This is an initial request to display the OAuth authorization link
		log_info("Displaying initial Twitch OAuth authorization page.")

		// Generate a new unique state string for this request
		state := generateRandomState(32) // Generate a random string of 32 characters

		// Store the state securely (e.g., in a session) before redirecting.
		// For this simple example, we use a map. In production, ensure this is tied to the user's session.
		oauthStates.Lock()
		oauthStates.m[state] = true
		oauthStates.Unlock()
		log_info(fmt.Sprintf("Generated and stored new state: %s", state))

		// Prepare data for the template
		data := OauthToken{
			ClientID: twitchClientID,
			StateRand: state,
			Title:     "Authorize with Twitch",
		}

		// Execute the template to render the page
		err := tmpls.ExecuteTemplate(w, "oauthPage", data)
		if err != nil {
			log_error(fmt.Sprintf("Failed to execute oauthPage template: %v", err))
			http.Error(w, "Internal server error: Could not render page.", http.StatusInternalServerError)
		}
	}
}

func main() {
	// Serve static files (ensure your `template/www/static` path is correct)
	fs := http.FileServer(http.Dir("./static/www"))
	http.Handle("/static/www/", http.StripPrefix("/static/www/", fs))

	// Register handlers
	http.HandleFunc("/", handler)
	http.HandleFunc("/oauth", oauthHandler) // This now handles both initial request and callback

	log_header("Server listening on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log_fatal(fmt.Sprintf("Server failed to start: %v", err))
	}
}

// handler for the root path
func handler(w http.ResponseWriter, r *http.Request) {
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

// Custom logging functions (kept as is)
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
	os.Exit(1) // Exit on critical errors during setup
}

func log_fatal(msg string) { // Added for graceful server shutdown logging
	fmt.Printf("%sFATAL: %s%s\n", LRED, msg, NC)
	log.Fatal(msg)
}
