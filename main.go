package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/ahgr3y/chirpy/internal/database"
	"github.com/joho/godotenv"
)

type apiConfig struct {
	fileserverHits int
	DB             *database.DB
	jwtSecret      string
}

func main() {

	// Load environment variables
	// by default, gotdotenv will look for a file named .env
	// in the current directory
	godotenv.Load()

	// load the JWT
	jwtSecret := os.Getenv("JWT_SECRET")

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
		jwtSecret:      jwtSecret,
	}

	// Create a ServeMux
	serveMux := http.NewServeMux()

	// Add a handler for the root path
	// By default, FileServer will look for index.html
	fileserverHandler := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(rootFilepath))))
	serveMux.Handle("/app/*", fileserverHandler)

	// Register handler for checking server readiness
	serveMux.HandleFunc("/api/healthz", handlerReadinessGet)
	serveMux.HandleFunc("POST /api/healthz", handlerReadinessPost)
	serveMux.HandleFunc("DELETE /api/healthz", handlerReadinessDelete)

	// Register handler to manage api metrics
	serveMux.HandleFunc("/admin/metrics", apiCfg.handlerGetServerHits)
	serveMux.HandleFunc("POST /admin/metrics", apiCfg.handlerPostServerHits)
	serveMux.HandleFunc("DELETE /admin/metrics", apiCfg.handlerDeleteServerHits)
	serveMux.HandleFunc("/api/reset", apiCfg.handlerResetServerHits)

	// Register handler to manage chirps
	serveMux.HandleFunc("POST /api/chirps", apiCfg.handlerChirpPost)
	serveMux.HandleFunc("GET /api/chirps", apiCfg.handlerChirpGet)
	serveMux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.handlerChirpGetByID)

	// Register handler to manage users
	serveMux.HandleFunc("POST /api/users", apiCfg.handlerUsersPost)
	serveMux.HandleFunc("PUT /api/users", apiCfg.handlerUpdateUser)
	serveMux.HandleFunc("POST /api/login", apiCfg.handlerUsersLogin)

	// Register handler to manage user tokens
	serveMux.HandleFunc("POST /api/refresh", apiCfg.handlerRefreshToken)
	serveMux.HandleFunc("POST /api/revoke", apiCfg.handlerRevokeRefreshToken)

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
