package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func respondWithError(w http.ResponseWriter, code int, errorMsg string) {

	if code > 499 {
		log.Printf("Unexpected error: %d - %s", code, errorMsg)
	}

	type errorResponse struct {
		Error string `json:"error"`
	}

	respondWithJSON(w, code, errorResponse{
		Error: errorMsg,
	})
}

func respondWithJSON(w http.ResponseWriter, code int, respBody interface{}) {

	// Parse responseVal struct to JSON
	dat, err := json.Marshal(respBody)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(dat)
}
