package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Viet-ph/Furniture-Store-Server/internal/server"
	"github.com/Viet-ph/Furniture-Store-Server/internal/service"
	"github.com/Viet-ph/Furniture-Store-Server/internal/utils"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("Port is not found in the environment")
	}

	queries, err := utils.ConnectDatabase()
	if err != nil {
		log.Fatal("Error conenction to database.")
	}

	srv := server.NewServer(
		service.NewUserService(queries),
		service.NewFeedService(queries),
		service.NewPostService(queries),
		service.NewFeedFollowService(queries),
	)
	server := &http.Server{
		Addr:    ":" + port,
		Handler: srv,
	}

	log.Printf("Listening on port: %s\n", port)
	log.Fatal(server.ListenAndServe())
}
