package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/ahgr3y/chirpy/internal/auth"
	"github.com/ahgr3y/chirpy/internal/database"
)

// handlerPostChirp stores the chirp in the request body
// and saves it to the database.
// Ensures that only authenticated user can post chirps.
func (cfg *apiConfig) handlerPostChirp(w http.ResponseWriter, r *http.Request) {

	// Extract token from request header
	token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")

	// Validate signature of token
	// and retrieve user id if token is valid
	idString, err := auth.ExtractIDFromToken(token, cfg.jwtSecret)
	if err != nil {
		log.Printf("Error extracting id from token: %s", err)
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Convert idString to int type
	userID, err := strconv.Atoi(string(idString))
	if err != nil {
		log.Printf("Error converting idString to int type: %s", err)
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	// To store JSON data from request
	type chirpStructure struct {
		Body string `json:"body"`
	}

	// Parse JSON Chirp to chirpStructure
	decoder := json.NewDecoder(r.Body)
	chirpStruct := chirpStructure{}
	err = decoder.Decode(&chirpStruct)
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
		AuthorID    int    `json:"author_id"`
		ID          int    `json:"id"`
		CleanedBody string `json:"body"`
	}

	// Save chirp to database
	chirpObj, err := cfg.DB.CreateChirp(userID, cleanChirp)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong while creating chirp")
	}

	// Respond valid response
	respondWithJSON(w, http.StatusCreated, validResp{
		AuthorID:    userID,
		ID:          chirpObj.ID,
		CleanedBody: cleanChirp,
	})

}

// handlerGetChirps responds with a JSON of all chirps in database in ascending order.
// If author_id query parameter is provided, handlerGetChirps will respond with
// all chirps created by author_id.
func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {

	// Retrieve chirps from database
	dbChirps, err := cfg.DB.GetChirps()
	if err != nil {
		log.Fatal(err)
	}

	// Create a copy of dbChirps
	chirps := []database.Chirp{}

	// Check if request parameter contains author_id
	idString := r.URL.Query().Get("author_id")

	// Retrieve chirps by id.
	if idString != "" {

		authorID, err := strconv.Atoi(idString)
		if err != nil {
			log.Printf("Error converting string to int: %s", err)
			respondWithError(w, http.StatusInternalServerError, "Something went wrong")
			return
		}

		chirps, err = cfg.DB.GetChirpsByID(authorID)
		if err != nil {
			log.Printf("Error retrieving chirps: %s", err)
			return
		}
	} else {
		for _, dbChirp := range dbChirps {
			chirps = append(chirps, database.Chirp{
				AuthorID: dbChirp.AuthorID,
				ID:       dbChirp.ID,
				Body:     dbChirp.Body,
			})
		}
	}

	// Check if request parameter contains sort
	sortBy := r.URL.Query().Get("sort")
	if sortBy == "" {
		sortBy = "asc"
	}

	database.SortChirpsByID(chirps, sortBy)

	respondWithJSON(w, http.StatusOK, chirps)
}

// handlerChirpGetByID response with a chirp with the given id.
func (cfg *apiConfig) handlerChirpGetByID(w http.ResponseWriter, r *http.Request) {

	// Get user's requested chirpID from URL path
	stringID := r.PathValue("chirpID")
	requestedID, err := strconv.Atoi(stringID)
	if err != nil {
		log.Printf("Error converting stringID to int: %s", err)
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

// handlerDeleteChirpByID deletes a Chirp in database with
// the associated ID in the request URL.
// Ensures that only authenticated and authorized user can delete chirp.
func (cfg *apiConfig) handlerDeleteChirpByID(w http.ResponseWriter, r *http.Request) {

	// Extract token from request header
	token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")

	// Validate signature of token
	// and retrieve user id if token is valid
	idString, err := auth.ExtractIDFromToken(token, cfg.jwtSecret)
	if err != nil {
		log.Printf("Error extracting id from token: %s", err)
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Convert idString to int type
	userID, err := strconv.Atoi(string(idString))
	if err != nil {
		log.Printf("Error converting idString to int type: %s", err)
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	// Get user's requested chirpID from URL path
	stringID := r.PathValue("chirpID")
	chirpID, err := strconv.Atoi(stringID)
	if err != nil {
		log.Printf("Error converting stringID to int: %s", err)
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID")
		return
	}

	// Delete chirp.
	err = cfg.DB.DeleteChirp(userID, chirpID)
	if err != nil {
		log.Printf("Error deleting chirp: %s", err)
		respondWithError(w, http.StatusForbidden, "Unauthorized to delete chirp")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
