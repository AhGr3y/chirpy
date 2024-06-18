package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/ahgr3y/chirpy/internal/database"
)

type apiConfig struct {
	fileserverHits int
	DB             *database.DB
}

func main() {

	// Set up debug flag
	dbg := flag.Bool("debug", false, "Enable debug mode")

	// Parse all flags
	flag.Parse()

	// Implement debug flag logic
	if *dbg { // Flag enabled

		fmt.Println("Debug mode is enabled. Clearing database.json...")

		// Delete database.json
		err := os.Remove("database.json")
		if err != nil {
			log.Fatal(err)
		} else {
			fmt.Println("database.json cleared successfully")
		}
	} else { // Flag disabled
		fmt.Println("Running in normal mode...")
	}

	const rootFilepath = "."
	const port = "8080"

	db, err := database.NewDB("database.json")
	if err != nil {
		log.Fatal(err)
	}

	apiCfg := apiConfig{
		fileserverHits: 0,
		DB:             db,
	}

	// Create a ServeMux
	serveMux := http.NewServeMux()

	// Add a handler for the root path
	// By default, FileServer will look for index.html
	fileserverHandler := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(rootFilepath))))
	serveMux.Handle("/app/*", fileserverHandler)

	// Register handler for healthz endpoint
	serveMux.HandleFunc("/api/healthz", handlerReadinessGet)
	serveMux.HandleFunc("POST /api/healthz", handlerReadinessPost)
	serveMux.HandleFunc("DELETE /api/healthz", handlerReadinessDelete)

	// Register handler for metrics endpoint
	serveMux.HandleFunc("/admin/metrics", apiCfg.handlerGetServerHits)
	serveMux.HandleFunc("POST /admin/metrics", apiCfg.handlerPostServerHits)
	serveMux.HandleFunc("DELETE /admin/metrics", apiCfg.handlerDeleteServerHits)

	// Register handler for reset endpoint
	serveMux.HandleFunc("/api/reset", apiCfg.handlerResetServerHits)

	// Register handler for chirps endpoint
	serveMux.HandleFunc("POST /api/chirps", apiCfg.handlerChirpPost)
	serveMux.HandleFunc("GET /api/chirps", apiCfg.handlerChirpGet)
	serveMux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.handlerChirpGetByID)

	// Register handler for users endpoint
	serveMux.HandleFunc("POST /api/users", apiCfg.handlerUsersPost)

	// Create a pointer to a server
	server := &http.Server{
		Addr:    ":" + port,
		Handler: serveMux,
	}

	// Start the server
	log.Printf("Serving files from %s on port %s...\n", rootFilepath, port)
	// Log errors
	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

}
