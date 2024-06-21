package auth

import (
	"errors"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func NewJWT(id int, expiresInSeconds int, secretKey string) (string, error) {

	// Create a JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(time.Second * time.Duration(expiresInSeconds))),
		Subject:   strconv.Itoa(id),
	})

	// Sign the token
	signedJWT, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return signedJWT, nil
}

func ExtractIDFromToken(token string, secretKey string) (string, error) {

	claimsStruct := jwt.RegisteredClaims{}

	// Validate signature of token
	// and retrieve user id if token is valid
	validToken, err := jwt.ParseWithClaims(token, &claimsStruct, func(t *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		return "", err
	}

	idString, err := validToken.Claims.GetSubject()
	if err != nil {
		return "", err
	}

	issuer, err := validToken.Claims.GetIssuer()
	if err != nil {
		return "", err
	}
	if issuer != "chirpy" {
		return "", errors.New("invalid issuer")
	}

	return idString, nil
}
