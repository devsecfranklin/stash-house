package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"

	"os"
	"strings"
	"testing"
	"text/template"

	"github.com/stretchr/testify/assert"
	// "github.com/stretchr/testify/require"

	// You might need to mock or stub these internal packages if they have side effects
	// or external dependencies that are hard to control in unit tests.
	// For simplicity, we're assuming they are well-behaved or can be bypassed for now.
	"internal/auth"
	"internal/logging"
)

// Mock the logging functions to prevent actual logging during tests
func init() {
	logging.Log_header = func(msg string) {}
	logging.Log_info = func(msg string) {}
	logging.Log_error = func(msg string) {}
	logging.Log_success = func(msg string) {}
	logging.Log_fatal = func(msg string) { panic(msg) } // Panic on fatal to catch it in tests
	logging.CheckError = func(err error) {
		if err != nil {
			panic(err)
		}
	}

	// Initialize templates for tests
	var err error
	tmpls, err = template.ParseGlob("template/www/*.tmpl") // Adjust path if necessary for your test environment
	if err != nil {
		// If templates are not found, tests that rely on them will fail.
		// This might require creating dummy template files for testing purposes.
		// For now, just log and allow tests to proceed, they might panic later.
		logging.Log_error("Failed to parse templates for tests: " + err.Error())
	}
}

// Helper to create a mock HTTP request
func newRequest(method, path string, body string) *http.Request {
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return req
}

// TestHandler tests the handler for the root path
func TestHandler(t *testing.T) {
	// Create a response recorder to capture the handler's output
	rr := httptest.NewRecorder()
	req := newRequest("GET", "/", "")

	// Call the handler function
	handler(rr, req)

	// Assert the status code
	assert.Equal(t, http.StatusOK, rr.Code, "Expected status OK for root path")

	// Assert that certain headers are set
	assert.Equal(t, "no-cache, no-store, must-revalidate", rr.Header().Get("Cache-Control"), "Cache-Control header mismatch")
	assert.Equal(t, "no-cache", rr.Header().Get("Pragma"), "Pragma header mismatch")
	assert.Equal(t, "0", rr.Header().Get("Expires"), "Expires header mismatch")

	// Assert that the body contains expected content (assuming indexPage renders something identifiable)
	assert.Contains(t, rr.Body.String(), "www.bitsmasher.net", "Expected body to contain site title")
}

// TestOauthHandlerInitialRequest tests the initial request to /oauth
func TestOauthHandlerInitialRequest(t *testing.T) {
	rr := httptest.NewRecorder()
	req := newRequest("GET", "/oauth", "")

	// Clear oauthStates for a clean test
	oauthStates.Lock()
	oauthStates.m = make(map[string]bool)
	oauthStates.Unlock()

	oauthHandler(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "Expected status OK for initial OAuth request")
	assert.Contains(t, rr.Body.String(), "Authorize with Twitch", "Expected body to contain OAuth authorization message")

	// Verify a state was generated and stored
	oauthStates.RLock()
	assert.Len(t, oauthStates.m, 1, "Expected one state to be stored")
	oauthStates.RUnlock()
}

// TestOauthHandlerCallbackSuccess tests a successful OAuth callback
func TestOauthHandlerCallbackSuccess(t *testing.T) {
	// Mock environment variables for Twitch credentials
	os.Setenv("TWITCH_CLIENT_ID", "mock_client_id")
	os.Setenv("TWITCH_CLIENT_SECRET", "mock_client_secret")
	defer func() {
		os.Unsetenv("TWITCH_CLIENT_ID")
		os.Unsetenv("TWITCH_CLIENT_SECRET")
	}()

	mockState := auth.GenerateRandomState(16)
	oauthStates.Lock()
	oauthStates.m[mockState] = true
	oauthStates.Unlock()

	// Mock HTTP client for token exchange
	// We need to create a custom HTTP client that always returns a successful mock response
	// This is a simplified mock for demonstration; a more robust solution would involve
	// creating a mock HTTP server or using a library like `http.DefaultClient.Transport`.
	oldHTTPClient := http.DefaultClient
	defer func() { http.DefaultClient = oldHTTPClient }() // Restore original client after test

	http.DefaultClient = &http.Client{
		Transport: RoundTripFunc(func(req *http.Request) (*http.Response, error) {
			assert.Contains(t, req.URL.String(), "https://id.twitch.tv/oauth2/token", "Token exchange URL mismatch")
			assert.Equal(t, "POST", req.Method, "Token exchange method mismatch")

			// Simulate successful token response
			tokenJSON := `{"access_token":"mock_access_token", "expires_in":3600, "refresh_token":"mock_refresh_token", "scope":["user:read:email"], "token_type":"bearer"}`
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(bytes.NewBufferString(tokenJSON)),
				Header:     make(http.Header),
			}, nil
		}),
	}

	rr := httptest.NewRecorder()
	// Simulate Twitch callback with code and state
	req := newRequest("GET", "/twitch/callback?code=mock_code&state="+mockState, "")

	oauthHandler(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "Expected status OK for successful OAuth callback")
	assert.Contains(t, rr.Body.String(), "mock_access_token", "Expected access token in success page")

	// Verify state was removed
	oauthStates.RLock()
	_, found := oauthStates.m[mockState]
	assert.False(t, found, "Expected state to be removed after use")
	oauthStates.RUnlock()
}

// RoundTripFunc is a type that implements http.RoundTripper for mocking HTTP client
type RoundTripFunc func(req *http.Request) (*http.Response, error)

func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

// TestOauthHandlerCallbackError tests OAuth callback with an error from Twitch
func TestOauthHandlerCallbackError(t *testing.T) {
	rr := httptest.NewRecorder()
	// Simulate Twitch callback with an error
	req := newRequest("GET", "/twitch/callback?error=access_denied&error_description=User+denied+access", "")

	oauthHandler(rr, req)

	assert.Equal(t, http.StatusForbidden, rr.Code, "Expected status Forbidden for OAuth error")
	assert.Contains(t, rr.Body.String(), "User denied access", "Expected error description in response")
}

// TestOauthHandlerMissingClientID tests oauthHandler when TWITCH_CLIENT_ID is missing
func TestOauthHandlerMissingClientID(t *testing.T) {
	os.Unsetenv("TWITCH_CLIENT_ID") // Ensure it's unset for this test
	os.Setenv("TWITCH_CLIENT_SECRET", "mock_secret")
	defer func() {
		os.Unsetenv("TWITCH_CLIENT_ID")
		os.Unsetenv("TWITCH_CLIENT_SECRET")
	}()

	rr := httptest.NewRecorder()
	req := newRequest("GET", "/oauth", "")

	oauthHandler(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "Expected status OK even if client ID is missing initially (page renders)")
	// Check for a logged error, as the `logging` package is mocked, we can't directly assert on logs.
	// In a real scenario, you might inject a mock logger to capture output.
	// For now, we assume the `logging.Log_error` call works.
}

// TestTwitchChatHandler tests the twitchChatHandler
func TestTwitchChatHandler(t *testing.T) {
	rr := httptest.NewRecorder()
	req := newRequest("GET", "/chatoverlay", "")

	// Set a mock twitchOauthToken
	twitchOauthToken = "mock_twitch_oauth_token"
	defer func() { twitchOauthToken = "" }() // Clean up

	twitchChatHandler(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "Expected status OK for chat overlay")
	// Assert that the token is passed to the template (assuming chatPage renders it)
	assert.Contains(t, rr.Body.String(), "mock_twitch_oauth_token", "Expected chat page to contain twitch token")
}

// Mock of internal/auth.GenerateRandomState for deterministic testing
func (a *auth.Auth) GenerateRandomState(n int) string {
	return strings.Repeat("a", n) // Return a predictable string
}

// Override the global auth.GenerateRandomState with our mock
func init() {
	auth.GenerateRandomState = (&auth.Auth{}).GenerateRandomState
}


/*
The following tests would require more extensive mocking of network calls and potentially the `template` package.
These are often better handled with integration tests or by making the `main` function more testable by
passing dependencies instead of relying on global variables.

func TestMainFunction(t *testing.T) {
	// Testing main() is tricky because it calls http.ListenAndServe, which blocks.
	// You typically test the components called by main() (e.g., handlers) separately.
	// For integration tests, you might run main() in a goroutine and then make HTTP requests to it.
}

func TestStaticFileServing(t *testing.T) {
	// This would involve setting up a test server and making requests to /static/www/
	// and asserting the content.
}
*/

