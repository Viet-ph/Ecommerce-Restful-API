package handler

import (
	"log"
	"net/http"

	"github.com/Viet-ph/Furniture-Store-Server/internal/helper"
	"github.com/Viet-ph/Furniture-Store-Server/internal/model"
	"github.com/Viet-ph/Furniture-Store-Server/internal/service"
)

type UserHandler struct {
	*service.UserService
}

func (u *UserHandler) CreateUser() http.HandlerFunc {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
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

		userCredential, err := model.NewUserCredential(req.Email, req.Password)
		if err != nil {
			helper.RespondWithError(w, http.StatusInternalServerError, "Error creating new user credential")
			return
		}

		newUser, err := u.Create(r.Context(),
			*userCredential,
			req.Location,
			req.Username)
		if err != nil {
			helper.RespondWithError(w, http.StatusInternalServerError, "Error creating new user")
			return
		}

		helper.RespondWithJSON(w, http.StatusCreated, newUser)
	}
}
