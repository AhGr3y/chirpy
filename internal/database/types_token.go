package database

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/ahgr3y/chirpy/internal/auth"
)

type RefreshToken struct {
	ID        int       `json:"id"`
	Token     string    `json:"refresh_token"`
	ExpiresAt time.Time `json:"expires_at"`
}

// CreateRefreshToken generates a refresh token
// and stores it it database
func (db *DB) CreateRefreshToken(id int) (RefreshToken, error) {

	// Generate a refresh token
	token, err := GenerateRefreshToken(id)
	if err != nil {
		return RefreshToken{}, err
	}

	// Save token to database
	err = db.SaveTokenToDB(token)
	if err != nil {
		return RefreshToken{}, err
	}

	return token, nil
}

// GenerateRefreshToken generates a refresh token
func GenerateRefreshToken(id int) (RefreshToken, error) {

	// Generate 32 bytes of random data in a slice
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return RefreshToken{}, err
	}

	// Convert bytes to hex string
	bytesHexString := hex.EncodeToString(bytes)

	// RefreshToken that expires in 60 days
	token := RefreshToken{
		ID:        id,
		Token:     bytesHexString,
		ExpiresAt: time.Now().Add(time.Second * 5184000),
	}

	return token, nil
}

func (db *DB) SaveTokenToDB(token RefreshToken) error {

	// Load database
	dbStructure, err := db.loadDB()
	if err != nil {
		return err
	}

	// Add/Update token to dbStructure
	dbStructure.RefreshTokens[token.ID] = token

	// Update database
	err = db.writeDB(dbStructure)
	if err != nil {
		return err
	}

	return nil
}

// RenewJWT checks the validity of refreshToken.
// If valid, returns a new JWT.
func (db *DB) RenewJWT(refreshToken string, secretKey string) (string, error) {

	// Validate refreshToken.
	id, err := db.validateRefreshToken(refreshToken)
	if err != nil {
		return "", err
	}

	// Create JWT.
	token, err := auth.NewJWT(id, secretKey)
	if err != nil {
		return "", err
	}

	return token, nil
}

// validateRefreshToken looks up refreshToken in the database.
// Returns an error message if it doesn't exist, or has expired.
// Otherwise, return the user id of the user that corresponds to refreshToken.
func (db *DB) validateRefreshToken(refreshToken string) (int, error) {

	// Load database.
	dbStructure, err := db.loadDB()
	if err != nil {
		return 0, err
	}

	// Check if refreshToken exist in dbStructure.
	// Check if refreshToken expired.
	dbTokens := dbStructure.RefreshTokens
	for _, dbToken := range dbTokens {
		tokenExist := false
		tokenNotExpired := false

		if dbToken.Token == refreshToken {
			tokenExist = true
		}
		if time.Now().Before(dbToken.ExpiresAt) {
			tokenNotExpired = true
		}

		if tokenExist && tokenNotExpired {
			return dbToken.ID, nil
		} else if !tokenExist {
			return 0, errors.New("refresh token does not exist")
		} else {
			return 0, errors.New("refresh token expired")
		}
	}

	return 0, errors.New("error validating refresh token")
}

// RevokeRefreshToken revokes the RefreshToken associated with
// refreshToken from the database.
func (db *DB) RevokeRefreshToken(refreshToken string) error {

	// Load database.
	dbStructure, err := db.loadDB()
	if err != nil {
		return err
	}

	// Revoke the associated RefreshToken
	dbTokens := dbStructure.RefreshTokens
	for _, dbToken := range dbTokens {
		if dbToken.Token == refreshToken {
			delete(dbTokens, dbToken.ID)
		}
	}

	// Update database
	err = db.writeDB(dbStructure)
	if err != nil {
		return err
	}

	return nil
}
