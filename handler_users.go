package main

import (
	"encoding/json"
	"log"
	"net/http"

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
	user, err := cfg.DB.CreateUser(param.Email, string(hashedPassword))
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

	// Authenticate user
	user, err := cfg.DB.AuthenticateUser(param.Email, param.Password)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized access")
	}

	type validResp struct {
		ID    int    `json:"id"`
		Email string `json:"email"`
	}

	// Respond valid response
	respondWithJSON(w, http.StatusOK, validResp{
		ID:    user.ID,
		Email: user.Email,
	})

}
