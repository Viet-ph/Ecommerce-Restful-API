package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Viet-ph/Furniture-Store-Server/internal/handler"
	"github.com/Viet-ph/Furniture-Store-Server/internal/helper"
	"github.com/Viet-ph/Furniture-Store-Server/internal/server"
	"github.com/Viet-ph/Furniture-Store-Server/internal/service"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("Port is not found in the environment")
	}

	queries, err := helper.ConnectDatabase()
	if err != nil {
		log.Fatal("Error conenction to database.")
	}

	log.Print("Database connected.")

	cartService := service.NewCartService(queries)
	userHandler := handler.NewUserHandler(service.NewUserService(queries), cartService)
	authHandler := handler.NewAuthHandler(service.NewAuthService(queries))
	productHandler := handler.NewProductHandler(service.NewProductService(queries))
	cartHandler := handler.NewCartHandler(cartService)

	srv := server.NewServer(
		authHandler,
		userHandler,
		productHandler,
		cartHandler,
	)
	server := &http.Server{
		Addr:    ":" + port,
		Handler: srv,
	}

	log.Printf("Listening on port: %s\n", port)
	log.Fatal(server.ListenAndServe())
}
