package handler

import (
	"log"
	"net/http"

	"github.com/Viet-ph/Furniture-Store-Server/internal/helper"
	"github.com/Viet-ph/Furniture-Store-Server/internal/service"
)

type AuthHandler struct {
	*service.AuthService
	*service.UserService
}

func (a *AuthHandler) UserLogin() http.HandlerFunc {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	type response struct {
		Email        string `json:"email"`
		Username     string `json:"username"`
		Location     string `json:"location"`
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
			Email:        user.Email,
			Username:     user.Username,
			Location:     user.Location,
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		})
	}
}
