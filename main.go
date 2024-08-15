package main

import (
	"TheCollectorDG/db"
	"TheCollectorDG/tasks"
	"TheCollectorDG/workers"
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	ctx := context.Background()

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	log.Println("Env variables loaded")

	pool, err := pgxpool.New(ctx, os.Getenv("DB_URL"))
	if err != nil {
		panic(err)
	}
	log.Println("Db connection successful")

	queries := db.New(pool)

	go workers.ClusterWorker(pool, queries, make(chan tasks.Task, 1000))
	workers.RegionWorker(pool, queries, make(chan tasks.Task, 1000))
}
