package middleware

import (
	"context"
	"net/http"

	db "github.com/Viet-ph/Furniture-Store-Server/internal/database"
	"github.com/Viet-ph/Furniture-Store-Server/internal/helper"
	"github.com/Viet-ph/Furniture-Store-Server/internal/service"
)

const ContextCartKey ContextKey = "userCart"

func NewCartMiddleware(cartService *service.CartService) func(http.Handler) http.Handler {
	//This will take the dependencies and return a authentication middleware that accepts only a single handler.
	//By doing this, will clean up the middleware function arguments and create closure to outter deps.
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, ok := r.Context().Value(ContextUserKey).(db.User)
			if !ok {
				helper.RespondWithError(w, http.StatusInternalServerError, "Context value is not User type")
				return
			}

			cart, err := cartService.GetCartId(r.Context(), user.ID)
			if err != nil {
				helper.RespondWithError(w, http.StatusInternalServerError, "Error getting cart Id, "+err.Error())
				return
			}

			ctx := context.WithValue(r.Context(), ContextCartKey, cart)
			handler.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
