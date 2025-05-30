package main

import (
	"fmt"
	"log"
	"net/http"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello from Go App! You requested: %s\n", r.URL.Path)
	log.Printf("Request received for: %s from %s\n", r.URL.Path, r.RemoteAddr)
}

func main() {
	// Register the handler for all paths
	http.HandleFunc("/", helloHandler)

	// Define the port for the Go application to listen on
	// It's common to use a port like 8080 for internal applications
	port := ":8080"
	fmt.Printf("Go application starting on http://localhost%s\n", port)

	// Start the HTTP server
	err := http.ListenAndServe(port, nil)
	if err!= nil {
		log.Fatalf("Go application failed to start: %v", err)
	}
}
