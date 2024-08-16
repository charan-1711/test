package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

var DB *sql.DB

func ConnectToDB() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading env file")
	}

	DB, err := sql.Open(
		"postgres",
		fmt.Sprintf(
			"host=%s user=%s dbname=%s password=%s sslmode=disable",
			os.Getenv("DB_HOST"),
			os.Getenv("DB_USER"),
			os.Getenv("DB_NAME"),
			os.Getenv("DB_PASSWORD"),
		),
	)

	if err != nil {
		log.Fatal("Error connecting to database")
	}

	if err = DB.Ping(); err != nil {
		log.Fatal("Error connecting to database")
	}

	log.Println("Successfully connected to database")
}
