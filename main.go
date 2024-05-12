package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Viet-ph/Furniture-Store-Server/internal/helper"
	"github.com/Viet-ph/Furniture-Store-Server/internal/server"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("Port is not found in the environment")
	}

	_, err := helper.ConnectDatabase()
	if err != nil {
		log.Fatal("Error conenction to database.")
	}

	srv := server.NewServer()
	server := &http.Server{
		Addr:    ":" + port,
		Handler: srv,
	}

	log.Printf("Listening on port: %s\n", port)
	log.Fatal(server.ListenAndServe())
}
