package helper

import (
	"database/sql"
	"log"
	"os"

	"github.com/Viet-ph/Furniture-Store-Server/internal/database"
)

func ConnectDatabase() (*database.Queries, error) {
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("Connection string is not found in the environment")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, err
	}
	dbQueries := database.New(db)

	return dbQueries, nil
}
