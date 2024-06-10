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
	"github.com/Viet-ph/Furniture-Store-Server/internal/model"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	*UserService
	Queries *db.Queries
}

const ACCESS_TOKEN_LIFETIME = 30 * time.Minute
const REFRESH_TOKEN_LIFETIME = 24 * 60 * time.Hour

func NewAuthService(userSv *UserService, q *db.Queries) *AuthService {
	return &AuthService{
		UserService: userSv,
		Queries:     q,
	}
}

func (a *AuthService) Login(context context.Context, email, password string) (user model.User, signedAccessToken, signedRefreshToken string, err error) {
	user, err = a.GetUserByEmail(context, email)
	if err == sql.ErrNoRows {
		return model.User{}, "", "", fmt.Errorf("wrong email")
	} else if err != nil {
		return model.User{}, "", "", err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return model.User{}, "", "", errors.New("wrong password")
	}

	signedAccessToken, err = signAccessToken(user.ID.String(), ACCESS_TOKEN_LIFETIME)
	if err != nil {
		return model.User{}, "", "", err
	}

	dbRefreshToken, err := a.Queries.GetValidTokenByUserId(context, user.ID)
	if err == sql.ErrNoRows {
		signedRefreshToken, err = signRefreshToken(user.ID.String(), REFRESH_TOKEN_LIFETIME)
		if err != nil {
			return model.User{}, "", "", err
		}

		_, err = a.Queries.SaveTokenToDB(context, db.SaveTokenToDBParams{
			ID:        uuid.New(),
			UserID:    user.ID,
			Token:     signedRefreshToken,
			ExpiresAt: time.Now().Add(REFRESH_TOKEN_LIFETIME),
			CreatedAt: time.Now().UTC(),
		})
		if err != nil {
			return model.User{}, "", "", err
		}
	} else {
		signedRefreshToken = dbRefreshToken.Token
	}

	return user, signedAccessToken, signedRefreshToken, nil
}

func (a *AuthService) RefreshAccessToken(context context.Context, refreshToken string) (string, error) {
	dbRefreshToken, err := a.Queries.GetTokenDetail(context, refreshToken)
	if err != nil {
		return "", fmt.Errorf("refresh access token failed, error: %v", err)
	}

	if valid, _ := isRefreshTokenValid(dbRefreshToken); !valid {
		return "", fmt.Errorf("invalid efresh token: %v", err)
	}

	userId, err := ValidateTokenAndExtractId(dbRefreshToken.Token)
	if err != nil {
		return "", fmt.Errorf("unable to get user Id from refresh token: %v", err)
	}

	return signAccessToken(userId, ACCESS_TOKEN_LIFETIME)
}

func ExtractTokenFromHeader(r *http.Request) string {
	reqToken := r.Header.Get("Authorization")
	splitToken := strings.Split(reqToken, "Bearer ")
	return splitToken[1]
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
