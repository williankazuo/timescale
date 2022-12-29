package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func InitDB() *sql.DB {
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		log.Fatal("ENV DB_HOST not set")
	}

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		log.Fatal("ENV DB_NAME not set")
	}

	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		log.Fatal("ENV DB_HOST not set")
	}

	dbPassword := os.Getenv("DB_PASSWORD")
	if dbPassword == "" {
		log.Fatal("ENV DB_PASSWORD not set")
	}

	connString := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", dbUser, dbPassword, dbHost, dbName)
	db, err := sql.Open("postgres", connString)
	if err != nil {
		log.Fatal("Error opening connection to dabatase")
	}

	return db
}
