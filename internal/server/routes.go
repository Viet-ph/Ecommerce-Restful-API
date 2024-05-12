package server

import (
	"net/http"
)

func addRoutes(
	mux *http.ServeMux,
) {
	//Protected routes
	//middlewareAuth := middleware.NewMiddlewareAuth(userService)

}
