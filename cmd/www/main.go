package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io" // Use io instead of ioutil
	"log"
	"math/rand" // For simple state generation (consider crypto/rand for production)
	"net/http"
	"os"
	"sync"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/twitch"
)

// A simple struct to hold the user data we get from Twitch.
type twitchUser struct {
	Data []struct {
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
	} `json:"data"`
}

var (
	LayoutDir string = "template/www"
	tmpls     *template.Template

	twitchOauthConfig *oauth2.Config

	// Temporary storage for OAuth states.
	// In a real application, use a secure, persistent storage (e.g., database, secure session store)
	// and associate it with a user session. This simple map is for demonstration ONLY.
	oauthStates = struct {
		sync.RWMutex
		m map[string]bool // map[state_string]is_valid
	}{m: make(map[string]bool)}
)

// OauthToken data structure passed to the template
type OauthTokenPageData struct {
	ClientID string
	AuthURL  string // The full authorization URL for the template
	Title    string
}

// Page data structure for a generic page
type Page struct {
	Title string
}

func init() {
	// Initialize random number generator for state generation (non-cryptographic)
	// For production, use crypto/rand for truly secure state generation.
	rand.NewSource(time.Now().UnixNano())

	// Load Twitch credentials from environment variables
	clientID := os.Getenv("TWITCH_CLIENT_ID")
	if clientID == "" {
		log_error("TWITCH_CLIENT_ID environment variable not set. Please set it.")
	}
	clientSecret := os.Getenv("TWITCH_CLIENT_SECRET")
	if clientSecret == "" {
		log_error("TWITCH_CLIENT_SECRET environment variable not set. Please set it.")
	}

	// Define the redirect URI. This MUST match your registered redirect URI
	// in the Twitch Developer Console. For local testing, it's typically localhost.
	redirectURI := "https//www.bitsmasher.net/twitch/callback"

	twitchOauthConfig = &oauth2.Config{
		RedirectURL:  redirectURI, // Use the correct local redirect URI
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scopes:       []string{"user:read:email"}, // Scopes determine what permissions you're asking for.
		Endpoint:     twitch.Endpoint,
	}

	// Load templates once at startup
	var err error
	tmpls, err = template.ParseGlob(LayoutDir + "/*.tmpl")
	if err != nil {
		log_error(fmt.Sprintf("Failed to parse templates: %v", err))
	}
}

// generateRandomState generates a random string for CSRF protection.
// In a real application, consider using crypto/rand for stronger randomness.
func generateRandomState(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))] // rand is safe for simple, non-cryptographic use
	}
	return string(b)
}

// oauthLoginHandler displays the initial OAuth authorization link.
func oauthLoginHandler(w http.ResponseWriter, r *http.Request) {
	log_header("Serving initial Twitch OAuth authorization page.")

	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	// Generate a new unique state string for this request
	state := generateRandomState(32) // Use a slightly longer state for better security

	// Store the state securely (e.g., in a session) before redirecting.
	// For this simple example, we use a map. In production, ensure this is tied to the user's session.
	oauthStates.Lock()
	oauthStates.m[state] = true
	oauthStates.Unlock()
	log_info(fmt.Sprintf("Generated and stored new state: %s", state))

	// Generate the full authorization URL that the user's browser will be redirected to
	authURL := twitchOauthConfig.AuthCodeURL(state)

	// Prepare data for the template
	data := OauthTokenPageData{
		ClientID: twitchOauthConfig.ClientID,
		AuthURL:  authURL, // Pass the full generated URL to the template
		Title:    "Authorize with Twitch",
	}

	// Execute the template to render the page
	err := tmpls.ExecuteTemplate(w, "oauthPage", data)
	if err != nil {
		log_error(fmt.Sprintf("Failed to execute oauthPage template: %v", err))
		http.Error(w, "Internal server error: Could not render page.", http.StatusInternalServerError)
	}
}

// handleTwitchCallback handles the callback from Twitch after user authorization.
func handleTwitchCallback(w http.ResponseWriter, r *http.Request) {
	log_header("Handling Twitch OAuth callback")

	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	// Extract parameters from the query string
	code := r.URL.Query().Get("code")
	returnedState := r.URL.Query().Get("state")
	twitchError := r.URL.Query().Get("error")
	twitchErrorDescription := r.URL.Query().Get("error_description")

	if twitchError != "" {
		log_info(fmt.Sprintf("Twitch OAuth Error: %s - %s", twitchError, twitchErrorDescription))
		http.Error(w, fmt.Sprintf("Twitch authorization failed: %s", twitchErrorDescription), http.StatusForbidden)
		return
	}

	if code == "" {
		log_error("No authorization code received in Twitch callback.")
		http.Error(w, "Missing authorization code.", http.StatusBadRequest)
		return
	}

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
	// Use the oauth2.Config's Exchange method, which correctly handles the POST request internally.
	token, err := twitchOauthConfig.Exchange(r.Context(), code) // Use r.Context() for better practice
	if err != nil {
		log_error(fmt.Sprintf("oauthConf.Exchange() failed: %v", err))
		http.Error(w, fmt.Sprintf("Failed to get access token from Twitch: %v", err), http.StatusInternalServerError)
		return
	}

	log_success("Successfully obtained Twitch access token!")
	log_info(fmt.Sprintf("Access Token: %s...", token.AccessToken[:10])) // Log a prefix, not full token
	log_info(fmt.Sprintf("Refresh Token: %s...", token.RefreshToken[:10]))
	log_info(fmt.Sprintf("Expires in: %d seconds", int(time.Until(token.Expiry).Seconds()))) // More accurate expiry calc
	log_info(fmt.Sprintf("Scopes: %v", token.Extra("scope"))) // Scopes might be in Extra

	// 3. Use the access token to get user data from the Twitch API
	log_info("Fetching user data from Twitch API...")
	client := twitchOauthConfig.Client(r.Context(), token) // Use r.Context()

	// The Twitch Helix API requires the Client-ID header, even for bearer token requests.
	req, err := http.NewRequest("GET", "https://api.twitch.tv/helix/users", nil)
	if err != nil {
		log_error(fmt.Sprintf("Failed to create Twitch API request: %v", err))
		http.Error(w, "Internal server error: Could not create Twitch API request.", http.StatusInternalServerError)
		return
	}
	req.Header.Set("Authorization", "Bearer "+token.AccessToken)
	req.Header.Set("Client-Id", twitchOauthConfig.ClientID)

	resp, err := client.Do(req)
	if err != nil {
		log_error(fmt.Sprintf("Failed to make Twitch API request: %v", err))
		http.Error(w, "Internal server error: Could not fetch user data from Twitch.", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body) // Read body for error logging
		log_error(fmt.Sprintf("Twitch API request failed with status %d: %s", resp.StatusCode, string(bodyBytes)))
		http.Error(w, fmt.Sprintf("Twitch API request failed: %s", string(bodyBytes)), http.StatusInternalServerError)
		return
	}

	// 4. Decode the user data and render the success page
	var user twitchUser
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		log_error(fmt.Sprintf("Failed to decode Twitch user data: %v", err))
		http.Error(w, "Internal server error: Could not decode Twitch user data.", http.StatusInternalServerError)
		return
	}

	// Create a data structure to pass to the template
	templateData := struct {
		User  twitchUser
		Token *oauth2.Token
	}{
		User:  user,
		Token: token,
	}

	// IMPORTANT: In a real app, store user.Data[0].ID, user.Data[0].Login, token.AccessToken, token.RefreshToken securely in your DB.

	// Assuming you have a 'success.html' template
	err = tmpls.ExecuteTemplate(w, "successPage", templateData) // Assuming "successPage" is defined in your templates
	if err != nil {
		log_error(fmt.Sprintf("Failed to execute successPage template: %v", err))
		http.Error(w, "Internal server error: Could not render success page.", http.StatusInternalServerError)
		return // Ensure to return after http.Error
	}
}

func main() {
	// --- Make static files available ------------------------------------------
	fs := http.FileServer(http.Dir("./static/www"))
	http.Handle("/static/www/", http.StripPrefix("/static/www/", fs))

	// Register handlers
	http.HandleFunc("/", handler)                  // Your main index page
	http.HandleFunc("/oauth", oauthLoginHandler)   // Handles initiating the OAuth flow
	http.HandleFunc("/twitch/callback", handleTwitchCallback) // Handles Twitch's redirect back

	log_header("Server listening on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log_fatal(fmt.Sprintf("Server failed to start: %v", err))
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
	// In a web server, you might not want to os.Exit(1) on every error,
	// especially in handlers, as it would crash the server.
	// For setup errors in init, it's more acceptable.
	// For handler errors, typically just log and return http.Error.
	// Keeping it for now as per your original intent in init() where it was used.
}

func log_fatal(msg string) { // Added for graceful server shutdown logging
	fmt.Printf("%sFATAL: %s%s\n", LRED, msg, NC)
	log.Fatal(msg)
}
