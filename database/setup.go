package database

import (
	"TheCollectorDG/config"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func dataSourceName(config config.MySqlConfig) string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s?tls=true", config.Username, config.Password, config.Host, config.DbName)
}

func Setup(config config.MySqlConfig) {
	var err error
	db, err = sql.Open("mysql", dataSourceName(config))
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("failed to ping: %v", err)
	}
	log.Println("Successfully connected to PlanetScale!")
}
