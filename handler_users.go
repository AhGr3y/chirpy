package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/ahgr3y/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerUsersPost(w http.ResponseWriter, r *http.Request) {

	// To store JSON data from request
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// Parse JSON to parameters
	decoder := json.NewDecoder(r.Body)
	param := parameters{}
	err := decoder.Decode(&param)
	if err != nil {
		log.Printf("Error decoding JSON: %s", err)
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	// Hash the password
	hashedPassword, err := auth.HashPassword(param.Password)
	if err != nil {
		log.Printf("Error hashing password: %s", err)
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	// Save chirp to database
	user, err := cfg.DB.CreateUser(param.Email, hashedPassword)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong while creating user")
	}

	type validResp struct {
		ID    int    `json:"id"`
		Email string `json:"email"`
	}

	// Respond valid response
	respondWithJSON(w, http.StatusCreated, validResp{
		ID:    user.ID,
		Email: user.Email,
	})

}

func (cfg *apiConfig) handlerUsersLogin(w http.ResponseWriter, r *http.Request) {

	// To store JSON data from request
	type parameters struct {
		Email            string `json:"email"`
		Password         string `json:"password"`
		ExpiresInSeconds int    `json:"expires_in_seconds"`
	}

	// Parse JSON to parameters
	decoder := json.NewDecoder(r.Body)
	param := parameters{}
	err := decoder.Decode(&param)
	if err != nil {
		log.Printf("Error decoding JSON: %s", err)
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	// Authenticate user
	user, err := cfg.DB.AuthenticateUser(param.Email, param.Password)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized access")
	}

	// Create a signedJWT
	signedJWT, err := auth.NewJWT(user.ID, validateExpiry(param.ExpiresInSeconds), cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
	}

	type validResp struct {
		ID    int    `json:"id"`
		Email string `json:"email"`
		Token string `json:"token"`
	}

	// Respond valid response
	respondWithJSON(w, http.StatusOK, validResp{
		ID:    user.ID,
		Email: user.Email,
		Token: signedJWT,
	})

}

// validateExpiry sets valid values for expiresInSeconds
func validateExpiry(expiresInSeconds int) int {

	// Set expiry to 24 hours if not assigned
	if expiresInSeconds == 0 {
		return 86400
	}

	// Expiry cannot exceed 24 hours
	if expiresInSeconds > 86400 {
		return 86400
	}

	return expiresInSeconds
}

// handlerUpdateUser updates user details with parameters from request
func (cfg *apiConfig) handlerUpdateUser(w http.ResponseWriter, r *http.Request) {

	// To store JSON data from request
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// Parse JSON to parameters
	decoder := json.NewDecoder(r.Body)
	param := parameters{}
	err := decoder.Decode(&param)
	if err != nil {
		log.Printf("Error decoding JSON: %s", err)
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

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

	// Convert idString to its int equivalent
	id, err := strconv.Atoi(idString)
	if err != nil {
		log.Printf("Error converting int to string: %s", err)
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	// Hash the password
	hashedPassword, err := auth.HashPassword(param.Password)
	if err != nil {
		log.Printf("Error hashing password: %s", err)
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	// Update user email and password
	user, err := cfg.DB.UpdateUser(id, param.Email, hashedPassword)
	if err != nil {
		log.Printf("Error updating user: %s", err)
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	type validResp struct {
		ID    int    `json:"id"`
		Email string `json:"email"`
	}

	respondWithJSON(w, http.StatusOK, validResp{
		ID:    user.ID,
		Email: user.Email,
	})
}
