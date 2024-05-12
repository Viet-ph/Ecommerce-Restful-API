package auth

import (
	"errors"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GetAPIKey(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("no authentication info found")
	}
	vals := strings.Split(authHeader, " ")
	if len(vals) != 2 {
		return "", errors.New("malformed auth header")
	}
	if vals[0] != "ApiKey" {
		return "", errors.New("malformed first part of auth header")
	}
	return vals[1], nil
}

func ExtractTokenFromHeader(r *http.Request) string {
	reqToken := r.Header.Get("Authorization")
	splitToken := strings.Split(reqToken, "Bearer ")
	return splitToken[1]
}

func ExtractIdFromToken(token *jwt.Token) (string, error) {
	id, err := token.Claims.GetSubject()
	if err != nil {
		return "", err
	}

	return id, nil
}

func createAccessToken(userId string) (singedTolen string, err error) {
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "Furniture-access",
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
		Subject:   userId,
	})
	return accessToken.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

func CreateRefreshToken(userId string) (singedTolen string, err error) {
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "Furniture-refresh",
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * 60 * time.Hour)),
		Subject:   userId,
	})
	return refreshToken.SignedString([]byte(os.Getenv("JWT_SECRET")))
}
