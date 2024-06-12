package server

import (
	"net/http"

	"github.com/Viet-ph/Furniture-Store-Server/internal/handler"
	"github.com/Viet-ph/Furniture-Store-Server/internal/middleware"
	"github.com/Viet-ph/Furniture-Store-Server/internal/service"
)

func addRoutes(
	mux *http.ServeMux,
	userService *service.UserService,
	userHandler *handler.UserHandler,
	authHandler *handler.AuthHandler,
) {
	//Protected routes
	middlewareAuth := middleware.NewMiddlewareAuth(userService)
	mux.Handle("GET /api/v1/users", middlewareAuth(userHandler.GetPersonalInfo()))
	mux.Handle("PUT /api/v1/users/password", middlewareAuth(userHandler.ChangePassword()))

	//Unrotected routes
	mux.HandleFunc("GET /api/v1/healthz", handler.Readiness)
	mux.HandleFunc("POST /api/v1/users", userHandler.UserSignUp())
	mux.HandleFunc("POST /api/v1/login", authHandler.UserLogin())
	mux.HandleFunc("POST /api/v1/revoke", authHandler.RevokeRefreshToken())
	mux.HandleFunc("POST /api/v1/refresh", authHandler.RefreshAccessToken())
}
