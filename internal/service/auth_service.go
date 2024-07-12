package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	db "github.com/Viet-ph/Furniture-Store-Server/internal/database"
	"github.com/Viet-ph/Furniture-Store-Server/internal/dto"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	queries *db.Queries
}

const ACCESS_TOKEN_LIFETIME = 30 * time.Minute
const REFRESH_TOKEN_LIFETIME = 24 * 60 * time.Hour

func NewAuthService(q *db.Queries) *AuthService {
	return &AuthService{
		queries: q,
	}
}

func (a *AuthService) Login(context context.Context, email, password string) (user dto.User, signedAccessToken, signedRefreshToken string, err error) {
	dbUser, err := a.queries.GetUserByEmail(context, email)
	if err == sql.ErrNoRows {
		return dto.User{}, "", "", fmt.Errorf("wrong email")
	} else if err != nil {
		return dto.User{}, "", "", fmt.Errorf("unable to get user by email: %v", err)
	}

	if err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(password)); err != nil {
		return dto.User{}, "", "", errors.New("wrong password")
	}

	signedAccessToken, err = signAccessToken(dbUser.ID.String(), ACCESS_TOKEN_LIFETIME)
	if err != nil {
		return dto.User{}, "", "", fmt.Errorf("unable to sign access token: %v", err)
	}

	dbRefreshToken, err := a.queries.GetValidTokenByUserId(context, dbUser.ID)
	if err == sql.ErrNoRows {
		signedRefreshToken, err = signRefreshToken(dbUser.ID.String(), REFRESH_TOKEN_LIFETIME)
		if err != nil {
			return dto.User{}, "", "", fmt.Errorf("unable to sign refresh token: %v", err)
		}

		_, err = a.queries.SaveTokenToDB(context, db.SaveTokenToDBParams{
			ID:        uuid.New(),
			UserID:    dbUser.ID,
			Token:     signedRefreshToken,
			ExpiresAt: time.Now().Add(REFRESH_TOKEN_LIFETIME),
			CreatedAt: time.Now().UTC(),
		})
		if err != nil {
			return dto.User{}, "", "", fmt.Errorf("unable to save refresh token to database: %v", err)
		}
	} else {
		signedRefreshToken = dbRefreshToken.Token
	}

	return dto.DbUsertoDto(&dbUser), signedAccessToken, signedRefreshToken, nil
}

func (a *AuthService) RefreshAccessToken(context context.Context, r *http.Request) (string, error) {
	refreshToken, err := ExtractTokenFromHeader(r)
	if err != nil {
		return "", fmt.Errorf("error getting token from header: %v", err)
	}

	dbRefreshToken, err := a.queries.GetTokenDetail(context, refreshToken)
	if err != nil {
		return "", fmt.Errorf("error getting refresh token info in database, error: %v", err)
	}

	if valid, _ := isRefreshTokenValid(dbRefreshToken); !valid {
		return "", fmt.Errorf("invalid refresh token: %v", err)
	}

	userId, err := ValidateTokenAndExtractId(dbRefreshToken.Token)
	if err != nil {
		return "", fmt.Errorf("unable to get user Id from refresh token: %v", err)
	}

	return signAccessToken(userId, ACCESS_TOKEN_LIFETIME)
}

func (a *AuthService) RevokeRefreshToken(context context.Context, refreshToken string) error {
	if err := a.queries.RevokeToken(context, refreshToken); err != nil {
		return fmt.Errorf("unable to revoke refresh token: %v", err)
	}

	return nil
}

// ExtractTokenFromHeader extracts the Bearer token from the Authorization header of the request.
// It returns the token string and an error, if any.
func ExtractTokenFromHeader(r *http.Request) (string, error) {
	// Get the Authorization header
	reqToken := r.Header.Get("Authorization")
	if reqToken == "" {
		return "", errors.New("authorization header missing")
	}

	// Split the token to separate the "Bearer " part
	splitToken := strings.Split(reqToken, "Bearer ")
	if len(splitToken) != 2 {
		return "", errors.New("invalid authorization header format")
	}

	// Return the token part
	token := splitToken[1]
	if token == "" {
		return "", errors.New("token missing after Bearer")
	}

	return token, nil
}

func ValidateTokenAndExtractId(tokenString string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil || !token.Valid {
		return "", err
	}

	id, err := token.Claims.GetSubject()
	if err != nil {
		return "", err
	}

	return id, nil
}

func signAccessToken(userId string, lifetime time.Duration) (singedToken string, err error) {
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "Furniture-access",
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(lifetime)),
		Subject:   userId,
	})
	return accessToken.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

func signRefreshToken(userId string, lifetime time.Duration) (singedToken string, err error) {
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "Furniture-refresh",
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(lifetime)),
		Subject:   userId,
	})
	return refreshToken.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

func isRefreshTokenValid(dbToken db.RefreshToken) (bool, error) {
	token, err := jwt.ParseWithClaims(dbToken.Token, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil || !token.Valid {
		return false, err
	} else if issuer, _ := token.Claims.GetIssuer(); issuer != "Furniture-refresh" {
		return false, errors.New("invalid refresh token issuer")
	} else if dbToken.Revoked {
		return false, errors.New("refresh token revoked")
	}

	return true, nil
}
