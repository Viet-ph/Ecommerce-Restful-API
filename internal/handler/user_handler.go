package handler

import (
	"log"
	"net/http"

	"github.com/Viet-ph/Furniture-Store-Server/internal/helper"
	"github.com/Viet-ph/Furniture-Store-Server/internal/middleware"
	"github.com/Viet-ph/Furniture-Store-Server/internal/model"
	"github.com/Viet-ph/Furniture-Store-Server/internal/service"
)

type UserHandler struct {
	*service.UserService
}

func (u *UserHandler) UserSignUp() http.HandlerFunc {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Username string `json:"username"`
		Location string `json:"location"`
	}

	type response struct {
		Email    string `json:"email"`
		Username string `json:"username"`
		Location string `json:"location"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req, err := helper.Decode[request](r)
		if err != nil {
			log.Printf("Error decoding parameters: %s", err)
			w.WriteHeader(500)
			return
		}

		user, err := u.Create(r.Context(), req.Location, req.Email, req.Password, req.Username)
		if err != nil {
			helper.RespondWithError(w, http.StatusInternalServerError, "Error creating new user: "+err.Error())
			return
		}

		helper.RespondWithJSON(w, http.StatusCreated, response{
			Email:    user.Email,
			Username: user.Username,
			Location: user.Location,
		})
	}
}

func (u *UserHandler) GetUserInfo() http.HandlerFunc {
	type response struct {
		Email     string `json:"email"`
		Username  string `json:"username"`
		Location  string `json:"location"`
		CreatedAt string `json:"created_at"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := r.Context().Value(middleware.ContextUserKey).(model.User)
		if !ok {
			helper.RespondWithError(w, http.StatusInternalServerError, "Context value is not User type")
			return
		}

		helper.RespondWithJSON(w, http.StatusOK,
			response{
				Email:     user.Email,
				Username:  user.Username,
				Location:  user.Location,
				CreatedAt: user.CreatedAt.String(),
			},
		)
	}
}
