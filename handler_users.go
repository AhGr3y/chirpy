package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/ahgr3y/chirpy/internal/database"
)

func (cfg *apiConfig) handlerUsersPost(w http.ResponseWriter, r *http.Request) {

	// To store JSON data from request
	type parameters struct {
		Email string `json:"email"`
	}

	// Parse JSON Chirp to parameters
	decoder := json.NewDecoder(r.Body)
	param := parameters{}
	err := decoder.Decode(&param)
	if err != nil {
		log.Printf("Error decoding JSON: %s", err)
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	// Save chirp to database
	userObj, err := cfg.DB.CreateUser(param.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong while creating user")
	}

	// Respond valid response
	respondWithJSON(w, http.StatusCreated, database.User{
		ID:    userObj.ID,
		Email: userObj.Email,
	})

}
