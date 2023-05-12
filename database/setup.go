package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func SetupConnection() {
	var err error
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s)/%s", os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_URL"), os.Getenv("DB_NAME"))
	db, err = sql.Open("mysql", dataSourceName)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("failed to ping: %v", err)
	}
	log.Println("Successfully connected to database!")
}
