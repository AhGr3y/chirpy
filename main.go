package main

import (
	"log"
	"net/http"
)

func main() {

	const port = "8080"

	// Create a handler
	serveMux := http.NewServeMux()

	// Create a server
	server := http.Server{
		Addr:    ":" + port,
		Handler: serveMux,
	}

	// Start the server
	log.Printf("Starting server on port %s...\n", port)
	// Log errors
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
