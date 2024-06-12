package main

import (
	"log"
	"net/http"
)

func main() {

	const rootFilepath = "."
	const port = "8080"

	// Create a ServeMux
	serveMux := http.NewServeMux()

	// Add a handler for the root path
	// By default, FileServer will look for index.html
	serveMux.Handle("/", http.FileServer(http.Dir(rootFilepath)))

	// Create a pointer to a server
	server := &http.Server{
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
