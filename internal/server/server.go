package server

import (
	"net/http"

	"github.com/Viet-ph/Furniture-Store-Server/internal/middleware"
)

func NewServer() http.Handler {
	mux := http.NewServeMux()
	addRoutes(mux)
	var handler http.Handler = mux
	handler = middleware.MiddlewareCors(handler)
	return handler
}
