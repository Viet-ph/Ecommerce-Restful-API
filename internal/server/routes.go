package server

import (
	"net/http"

	"github.com/Viet-ph/Furniture-Store-Server/internal/handler"
	"github.com/Viet-ph/Furniture-Store-Server/internal/middleware"
)

func addRoutes(
	mux *http.ServeMux,
	userHandler *handler.UserHandler,
	authHandler *handler.AuthHandler,
	productHandler *handler.ProductHandler,
	cartHandler *handler.CartHandler,

) {
	//Protected routes
	authMiddleware := middleware.NewAuthMiddleware(userHandler.UserService)
	cartMiddleware := middleware.NewCartMiddleware(cartHandler.CartService)

	mux.Handle("GET /api/v1/users/profile", authMiddleware(userHandler.GetPersonalInfo()))
	mux.Handle("PUT /api/v1/users/password", authMiddleware(userHandler.ChangePassword()))
	mux.Handle("DELETE /api/v1/users", authMiddleware(userHandler.DeleteAccount()))

	mux.Handle("POST /api/v1/carts/items", authMiddleware(cartMiddleware(cartHandler.AddNewCartItem())))
	mux.Handle("GET /api/v1/carts/items", authMiddleware(cartMiddleware(cartHandler.GetAllCartItems())))
	mux.Handle("PUT /api/v1/carts/items", authMiddleware(cartMiddleware(cartHandler.UpdateCartItem())))
	mux.Handle("DELETE /api/v1/carts/items", authMiddleware(cartMiddleware(cartHandler.RemoveCartItem())))
	mux.Handle("PUT /api/v1/carts", authMiddleware(cartMiddleware(cartHandler.UpdateCart())))
	mux.Handle("DELETE /api/v1/carts/items", authMiddleware(cartMiddleware(cartHandler.ClearCart())))

	//Unprotected routes
	mux.HandleFunc("GET /api/v1/healthz", handler.Readiness)

	mux.HandleFunc("GET /api/v1/products", productHandler.GetProductsWithFilters())
	mux.HandleFunc("GET /api/v1/products/{id}", productHandler.GetProductDetail())
	mux.HandleFunc("DELETE /api/v1/products/{id}", productHandler.DeleteProduct())

	mux.HandleFunc("POST /api/v1/users", userHandler.UserSignUp())
	mux.HandleFunc("POST /api/v1/login", authHandler.UserLogin())
	mux.HandleFunc("POST /api/v1/revoke", authHandler.RevokeRefreshToken())
	mux.HandleFunc("POST /api/v1/refresh", authHandler.RefreshAccessToken())
}
