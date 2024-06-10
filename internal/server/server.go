package server

import (
	"net/http"

	"github.com/Viet-ph/Furniture-Store-Server/internal/handler"
	"github.com/Viet-ph/Furniture-Store-Server/internal/middleware"
	"github.com/Viet-ph/Furniture-Store-Server/internal/service"
)

func NewServer(
	userService *service.UserService,
	authHandler *handler.AuthHandler,
	userHandler *handler.UserHandler,
) http.Handler {
	mux := http.NewServeMux()
	addRoutes(mux,
		userService,
		userHandler,
		authHandler,
	)
	var handler http.Handler = mux
	handler = middleware.MiddlewareCors(handler)
	return handler
}
