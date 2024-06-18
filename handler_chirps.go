package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/ahgr3y/chirpy/internal/database"
)

func (cfg *apiConfig) handlerChirpPost(w http.ResponseWriter, r *http.Request) {

	// To store JSON data from request
	type chirpStructure struct {
		Body string `json:"body"`
	}

	// Parse JSON Chirp to chirpStructure
	decoder := json.NewDecoder(r.Body)
	chirpStruct := chirpStructure{}
	err := decoder.Decode(&chirpStruct)
	if err != nil {
		log.Printf("Error decoding JSON: %s", err)
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	// Validate chirp body
	cleanChirp, err := validateChirp(chirpStruct.Body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
	}

	type validResp struct {
		ID          int    `json:"id"`
		CleanedBody string `json:"body"`
	}

	// Save chirp to database
	chirpObj, err := cfg.DB.CreateChirp(cleanChirp)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong while creating chirp")
	}

	// Respond valid response
	respondWithJSON(w, http.StatusCreated, validResp{
		ID:          chirpObj.ID,
		CleanedBody: cleanChirp,
	})

}

// handlerChirpGet responds with a JSON of all chirps in database in ascending order
func (cfg *apiConfig) handlerChirpGet(w http.ResponseWriter, r *http.Request) {

	// Retrieve chirps from database
	dbChirps, err := cfg.DB.GetChirps()
	if err != nil {
		log.Fatal(err)
	}

	// Create a copy of dbChirps
	chirps := []database.Chirp{}
	for _, dbChirp := range dbChirps {
		chirps = append(chirps, database.Chirp{
			ID:   dbChirp.ID,
			Body: dbChirp.Body,
		})
	}

	// Sort chirps in ascending order of id
	sort.Slice(chirps, func(i, j int) bool {
		return chirps[i].ID < chirps[j].ID
	})

	respondWithJSON(w, http.StatusOK, chirps)
}

// handlerChirpGetByID response with a chirp with the given id.
func (cfg *apiConfig) handlerChirpGetByID(w http.ResponseWriter, r *http.Request) {

	// Get user's requested chirpID from URL path
	stringID := r.PathValue("chirpID")
	requestedID, err := strconv.Atoi(stringID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID")
		return
	}

	// Retrieve chirps from database.
	chirp, err := cfg.DB.GetChirp(requestedID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Failed to retrieve chirp")
		return
	}

	respondWithJSON(w, http.StatusOK, database.Chirp{
		ID:   chirp.ID,
		Body: chirp.Body,
	})
}

func cleanBody(body string, profanities []string) string {

	// Get words from original body
	originalBodyWords := strings.Split(body, " ")

	// Get words from lower case body
	lowerBody := strings.ToLower(body)
	words := strings.Split(lowerBody, " ")

	// Replace profanities with ****
	for i, word := range words {
		for _, profane := range profanities {
			if word == profane {
				originalBodyWords[i] = "****"
				continue
			}
		}
	}

	// Join back words into a sentence
	cleanedBody := strings.Join(originalBodyWords, " ")
	return cleanedBody
}

func validateChirp(body string) (string, error) {

	// Chirp cannot be too long
	const maxChirpLength = 140
	if len(body) > maxChirpLength {
		return "", errors.New("chirp is too long")
	}

	profanities := []string{
		"kerfuffle",
		"sharbert",
		"fornax",
	}

	// Sensor profanities
	cleaned := cleanBody(body, profanities)

	return cleaned, nil

}
