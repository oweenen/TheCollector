package database

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func SetupConnection() {
	var err error
	db, err = sql.Open("mysql", os.Getenv("DSN"))
	if err != nil {
		log.Fatalf("Failed to connect to PlanetScale: %v\n", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping PlanetScale: %v\n", err)
	}
	log.Println("Successfully connected to PlanetScale!")
}
