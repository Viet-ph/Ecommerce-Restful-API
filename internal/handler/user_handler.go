package handler

import (
	"log"
	"net/http"

	db "github.com/Viet-ph/Furniture-Store-Server/internal/database"
	"github.com/Viet-ph/Furniture-Store-Server/internal/dto"
	"github.com/Viet-ph/Furniture-Store-Server/internal/helper"
	"github.com/Viet-ph/Furniture-Store-Server/internal/middleware"
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
	return func(w http.ResponseWriter, r *http.Request) {
		log.Print("Hit Sign Up endpoint.")
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

		helper.RespondWithJSON(w, http.StatusCreated, user)
	}
}

func (u *UserHandler) GetPersonalInfo() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := r.Context().Value(middleware.ContextUserKey).(db.User)
		if !ok {
			helper.RespondWithError(w, http.StatusInternalServerError, "Context value is not User type")
			return
		}

		helper.RespondWithJSON(w, http.StatusOK, dto.DbUsertoDto(&user))
	}
}

func (u *UserHandler) ChangePassword() http.HandlerFunc {
	type request struct {
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := r.Context().Value(middleware.ContextUserKey).(db.User)
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

		helper.RespondWithJSON(w, http.StatusOK, "Password change successfully!")
	}
}

func (u *UserHandler) DeleteAccount() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := r.Context().Value(middleware.ContextUserKey).(db.User)
		if !ok {
			helper.RespondWithError(w, http.StatusInternalServerError, "Context value is not User type")
			return
		}

		err := u.DeleteUserById(r.Context(), user.ID)
		if err != nil {
			helper.RespondWithError(w, http.StatusInternalServerError, "Couldn't delete account with given user id: "+user.ID.String())
			return
		}

		helper.RespondWithJSON(w, http.StatusOK, "User with id: "+user.ID.String()+" deleted successfully.")
	}
}
