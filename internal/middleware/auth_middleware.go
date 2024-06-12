package middleware

import (
	"context"
	"net/http"

	"github.com/Viet-ph/Furniture-Store-Server/internal/helper"
	"github.com/Viet-ph/Furniture-Store-Server/internal/service"
	"github.com/google/uuid"
)

type ContextKey string

const ContextUserKey ContextKey = "user"

func NewMiddlewareAuth(userService *service.UserService) func(http.Handler) http.Handler {
	//This will take the dependencies and return a authentication middleware that accepts only a single handler.
	//By doing this, will clean up the middleware function arguments and create closure to outter deps.
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenString, err := service.ExtractTokenFromHeader(r)
			if err != nil {
				helper.RespondWithError(w, http.StatusUnauthorized, "Cannot get token from header: "+err.Error())
			}

			userId, err := service.ValidateTokenAndExtractId(tokenString)
			if err != nil {
				helper.RespondWithError(w, http.StatusUnauthorized, "Unauthorized Access: "+err.Error())
				return
			}

			user, err := userService.GetUserById(r.Context(), uuid.MustParse(userId))
			if err != nil {
				helper.RespondWithError(w, http.StatusUnauthorized, "Unauthorized Access: "+err.Error())
				return
			}

			ctx := context.WithValue(r.Context(), ContextUserKey, user)
			handler.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
