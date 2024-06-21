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
		log.Printf("Error creating user: %s", err)
		respondWithError(w, http.StatusInternalServerError, "Something went wrong while creating user")
		return
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
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// Parse JSON request body to parameters
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
		log.Printf("Error authenticating user: %s", err)
		respondWithError(w, http.StatusUnauthorized, "Unauthorized access")
		return
	}

	// Create a signedJWT
	signedJWT, err := auth.NewJWT(user.ID, cfg.jwtSecret)
	if err != nil {
		log.Printf("Error creating JWT: %s", err)
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	// Create RefreshToken
	refreshToken, err := cfg.DB.CreateRefreshToken(user.ID)
	if err != nil {
		log.Printf("Error generating refresh token: %s", err)
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	type validResp struct {
		ID           int    `json:"id"`
		Email        string `json:"email"`
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}

	// Respond valid response
	respondWithJSON(w, http.StatusOK, validResp{
		ID:           user.ID,
		Email:        user.Email,
		Token:        signedJWT,
		RefreshToken: refreshToken.Token,
	})

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
	user, err := cfg.DB.UpdateUserEmailPassword(id, param.Email, hashedPassword)
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

func (cfg *apiConfig) handlerRefreshToken(w http.ResponseWriter, r *http.Request) {

	// Extract token from request header
	token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")

	// Renew JWT
	token, err := cfg.DB.RenewJWT(token, cfg.jwtSecret)
	if err != nil {
		log.Printf("Error renewing JWT: %s", err)
		respondWithError(w, http.StatusUnauthorized, "Token doesn't exist or expired")
		return
	}

	type validResp struct {
		Token string `json:"token"`
	}

	respondWithJSON(w, http.StatusOK, validResp{
		Token: token,
	})
}

// handlerRevokeRefreshToken deletes the refresh token
// (associated with the refresh token extracted from the request header)
// from the database.
func (cfg *apiConfig) handlerRevokeRefreshToken(w http.ResponseWriter, r *http.Request) {

	// Extract token from request header
	token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")

	// Revoke refresh token
	err := cfg.DB.RevokeRefreshToken(token)
	if err != nil {
		log.Printf("Error revoking refresh token: %s", err)
		respondWithError(w, http.StatusBadRequest, "Something went wrong")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
