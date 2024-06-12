package handler

import (
	"log"
	"net/http"

	"github.com/Viet-ph/Furniture-Store-Server/internal/dto"
	"github.com/Viet-ph/Furniture-Store-Server/internal/helper"
	"github.com/Viet-ph/Furniture-Store-Server/internal/service"
)

type AuthHandler struct {
	*service.AuthService
	*service.UserService
}

func NewAuthHandler(a *service.AuthService, u *service.UserService) *AuthHandler {
	return &AuthHandler{
		AuthService: a,
		UserService: u,
	}
}

func (a *AuthHandler) UserLogin() http.HandlerFunc {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	type response struct {
		dto.User
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req, err := helper.Decode[request](r)
		if err != nil {
			log.Printf("Error decoding parameters: %s", err)
			w.WriteHeader(500)
			return
		}

		user, accessToken, refreshToken, err := a.Login(r.Context(), req.Email, req.Password)
		if err != nil {
			helper.RespondWithError(w, http.StatusUnauthorized, err.Error())
			return
		}

		helper.RespondWithJSON(w, http.StatusOK, response{
			User:         dto.DbUsertoDto(&user),
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		})
	}
}

func (a *AuthHandler) RefreshAccessToken() http.HandlerFunc {
	type request struct {
		RefreshToken string `json:"refresh_token"`
	}

	type response struct {
		AccessToken string `json:"access_token"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req, err := helper.Decode[request](r)
		if err != nil {
			log.Printf("Error decoding parameters: %s", err)
			w.WriteHeader(500)
			return
		}

		newAccessToken, err := a.AuthService.RefreshAccessToken(r.Context(), req.RefreshToken)
		if err != nil {
			helper.RespondWithError(w, http.StatusUnauthorized, err.Error())
			return
		}

		helper.RespondWithJSON(w, http.StatusAccepted, response{
			AccessToken: newAccessToken,
		})
	}
}

func (a *AuthHandler) RevokeRefreshToken() http.HandlerFunc {
	type request struct {
		RefreshToken string `json:"refresh_token"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req, err := helper.Decode[request](r)
		if err != nil {
			log.Printf("Error decoding parameters: %s", err)
			w.WriteHeader(500)
			return
		}

		err = a.AuthService.RevokeRefreshToken(r.Context(), req.RefreshToken)
		if err != nil {
			helper.RespondWithError(w, http.StatusUnauthorized, "Error revoking token: "+err.Error())
		}
	}
}
