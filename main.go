package main

import (
	"log"
	"net/http"
)

type apiConfig struct {
	fileserverHits int
}

func main() {

	const rootFilepath = "."
	const port = "8080"

	apiCfg := apiConfig{
		fileserverHits: 0,
	}

	// Create a ServeMux
	serveMux := http.NewServeMux()

	// Add a handler for the root path
	// By default, FileServer will look for index.html
	fileserverHandler := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(rootFilepath))))
	serveMux.Handle("/app/*", fileserverHandler)

	// Register handler for readiness endpoint
	serveMux.HandleFunc("/api/healthz", handlerReadinessGet)
	serveMux.HandleFunc("POST /api/healthz", handlerReadinessPost)
	serveMux.HandleFunc("DELETE /api/healthz", handlerReadinessDelete)

	// Register handler for metrics endpoint for checking server hits
	serveMux.HandleFunc("/admin/metrics", apiCfg.handlerGetServerHits)
	serveMux.HandleFunc("POST /admin/metrics", apiCfg.handlerPostServerHits)
	serveMux.HandleFunc("DELETE /admin/metrics", apiCfg.handlerDeleteServerHits)
	serveMux.HandleFunc("/api/reset", apiCfg.handlerResetServerHits)

	// Create a pointer to a server
	server := &http.Server{
		Addr:    ":" + port,
		Handler: serveMux,
	}

	// Start the server
	log.Printf("Serving files from %s on port %s...\n", rootFilepath, port)
	// Log errors
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

}
