package auth

import (
	"errors"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

var ErrNoAuthHeaderIncluded = errors.New("no auth header included in request")

func HashPassword(password string) (string, error) {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

func AuthenticatePassword(hashedPassword string, password string) error {

	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return err
	}

	return nil
}

// ExtractApiKey extracts the api key from headers.
func ExtractApiKey(headers http.Header) (string, error) {

	authHeader := headers.Get("Authorization")

	if authHeader == "" {
		return "", ErrNoAuthHeaderIncluded
	}

	splitAuth := strings.Split(authHeader, " ")
	if len(splitAuth) < 2 || splitAuth[0] != "ApiKey" {
		return "", errors.New("malformed authorization header")
	}

	return splitAuth[1], nil
}
