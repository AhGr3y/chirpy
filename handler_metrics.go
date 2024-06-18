package main

import (
	"fmt"
	"net/http"
)

func (cfg *apiConfig) handlerGetServerHits(w http.ResponseWriter, r *http.Request) {
	body := []byte(fmt.Sprintf(`<html>

<body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
</body>

</html>`, cfg.fileserverHits))
	w.Header().Add("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

// Converts next to a Handler that increments fileserverHits
func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits++
		next.ServeHTTP(w, r)
	})
}

// Prevent POST request to metrics route
func (cfg *apiConfig) handlerPostServerHits(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusMethodNotAllowed)
}

// Prevent DELETE request to metrics route
func (cfg *apiConfig) handlerDeleteServerHits(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusMethodNotAllowed)
}
