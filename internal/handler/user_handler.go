package handler

import (
	"log"
	"net/http"

	"github.com/Viet-ph/Furniture-Store-Server/internal/helper"
	"github.com/Viet-ph/Furniture-Store-Server/internal/middleware"
	"github.com/Viet-ph/Furniture-Store-Server/internal/model"
	"github.com/Viet-ph/Furniture-Store-Server/internal/service"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	*service.UserService
}

func NewUserHandler(u *service.UserService) *UserHandler {
	return &UserHandler{
		UserService: u,
	}
}

func (u *UserHandler) UserSignUp() http.HandlerFunc {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Username string `json:"username"`
		Location string `json:"location"`
	}

	type response struct {
		Id       string `json:"id"`
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
			Id:       user.ID.String(),
			Email:    user.Email,
			Username: user.Username,
			Location: user.Location,
		})
	}
}

func (u *UserHandler) GetPersonalInfo() http.HandlerFunc {
	type response struct {
		Id        string `json:"id"`
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
				Id:        user.ID.String(),
				Email:     user.Email,
				Username:  user.Username,
				Location:  user.Location,
				CreatedAt: user.CreatedAt.String(),
			},
		)
	}
}

func (u *UserHandler) ChangePassword() http.HandlerFunc {
	type request struct {
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := r.Context().Value(middleware.ContextUserKey).(model.User)
		if !ok {
			helper.RespondWithError(w, http.StatusInternalServerError, "Context value is not User type")
			return
		}

		req, err := helper.Decode[request](r)
		if err != nil {
			log.Printf("Error decoding parameters: %s", err)
			w.WriteHeader(500)
			return
		}

		if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword)); err != nil {
			helper.RespondWithError(w, http.StatusBadRequest, "Wrong password")
			return
		}

		err = u.UserService.UpdateUserPassword(r.Context(), user.ID, req.NewPassword)
		if err != nil {
			helper.RespondWithError(w, http.StatusBadRequest, "Error updating new password: "+err.Error())
			return
		}

		helper.RespondWithJSON(w, http.StatusOK, struct{}{})
	}
}
