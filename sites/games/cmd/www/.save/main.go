package main

import (
    "fmt"
    "html/template"
    "log"
    "net/http"
    "regexp"
    "strconv"
    "strings"
)

// FormData struct to hold parsed and validated data, and errors/success messages for template rendering
type FormData struct {
    Name    string
    Email   string
    Message string
    Rating  string   // Stored as string to easily re-populate HTML form
    Category string
    Errors string // Slice to collect validation errors
    Success string   // Success message
}

// isValidEmail checks for a basic email format using a regular expression.
// Note: A more robust regex or dedicated library should be used in production.
func isValidEmail(email string) bool {
    emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
    return emailRegex.MatchString(email)
}

func feedbackHandler(w http.ResponseWriter, r *http.Request) {
    tmpl, err := template.ParseFiles("feedback.html")
    if err!= nil {
        http.Error(w, "Unable to load template", http.StatusInternalServerError)
        return
    }

    data := FormData{} // Initialize with empty data for GET request or initial display

    if r.Method == http.MethodPost {
        err := r.ParseForm() // Parse all form data from the request body
        if err!= nil {
            data.Errors = append(data.Errors, "Failed to parse form data.")
            tmpl.Execute(w, data) // Re-render form with error
            return
        }

        // Retrieve and populate data from the form. Trim spaces for cleaner input.
        data.Name = strings.TrimSpace(r.FormValue("name"))
        data.Email = strings.TrimSpace(r.FormValue("email"))
        data.Message = strings.TrimSpace(r.FormValue("message"))
        data.Rating = r.FormValue("rating")     // Radio button: will be empty string if not selected
        data.Category = r.FormValue("category") // Dropdown: will be empty if default "-- Select --" is chosen

        // --- Server-side Validation ---
        // Validate Name
        if data.Name == "" {
            data.Errors = append(data.Errors, "Name is required.")
        }

        // Validate Email
        if data.Email == "" {
            data.Errors = append(data.Errors, "Email is required.")
        } else if!isValidEmail(data.Email) {
            data.Errors = append(data.Errors, "Invalid email format.")
        }

        // Validate Message
        if data.Message == "" {
            data.Errors = append(data.Errors, "Message is required.")
        } else if len(data.Message) < 10 {
            data.Errors = append(data.Errors, "Message must be at least 10 characters long.")
        }

        // Validate Rating (radio button)
        validRatings := map[string]bool{"1": true, "2": true, "3": true, "4": true, "5": true}
        if data.Rating == "" {
            data.Errors = append(data.Errors, "Please select an overall rating.")
        } else if _, ok := validRatings;!ok {
            data.Errors = append(data.Errors, "Invalid rating value.")
        } else {
            // Optional: convert rating to int for further processing, but keep string for re-display
            _, err := strconv.Atoi(data.Rating)
            if err!= nil {
                data.Errors = append(data.Errors, "Rating is not a valid number.") // Should not happen if validRatings check passes
            }
        }

        // Validate Category (dropdown)
        validCategories := map[string]bool{"product": true, "support": true, "bug": true, "other": true}
        if data.Category == "" {
            data.Errors = append(data.Errors, "Please select a feedback category.")
        } else if _, ok := validCategories[data.Category];!ok {
            data.Errors = append(data.Errors, "Invalid feedback category.")
        }

        if len(data.Errors) > 0 {
            // If there are errors, re-render the form with errors and previous inputs
            tmpl.Execute(w, data)
            return
        }

        // --- If validation passes, process the data ---
        log.Printf("Received Feedback:\nName: %s\nEmail: %s\nMessage: %s\nRating: %s\nCategory: %s\n",
            data.Name, data.Email, data.Message, data.Rating, data.Category)

        // In a real application, this is where data would be saved to a database,
        // sent as an email, or integrated with other services.
        data.Success = "Thank you for your feedback! It has been successfully submitted."
        // Clear form data on success for a fresh display, but retain success message
        tmpl.Execute(w, FormData{Success: data.Success})
        return
    }

    // For GET requests, simply render the empty form
    tmpl.Execute(w, data)
}

func main() {
    http.HandleFunc("/feedback", feedbackHandler)

    port := ":8080"
    fmt.Printf("Server starting on port %s...\n", port)
    err := http.ListenAndServe(port, nil)
    if err!= nil {
        log.Fatalf("Server failed to start: %v", err)
    }
}