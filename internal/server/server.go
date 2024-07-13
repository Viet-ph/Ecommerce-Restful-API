package server

import (
	"net/http"

	"github.com/Viet-ph/Furniture-Store-Server/internal/handler"
	"github.com/Viet-ph/Furniture-Store-Server/internal/middleware"
)

func NewServer(
	authHandler *handler.AuthHandler,
	userHandler *handler.UserHandler,
	productHandler *handler.ProductHandler,
	cartHandler *handler.CartHandler,
) http.Handler {
	mux := http.NewServeMux()
	addRoutes(mux,
		userHandler,
		authHandler,
		productHandler,
		cartHandler,
	)
	var handler http.Handler = mux
	handler = middleware.MiddlewareCors(handler)
	return handler
}
