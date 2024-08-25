package main

import (
	"TheCollectorDG/api"
	"TheCollectorDG/db"
	"TheCollectorDG/workerManager"
	"TheCollectorDG/workers"
	"net/http"

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

	workerEnv := workers.WorkerEnv{
		Pool:    pool,
		Queries: queries,
	}

	workerManager := workerManager.New()
	workerManager.AddWorker("na1", workerEnv.RegionWorker)
	workerManager.AddWorker("americas", workerEnv.ClusterWorker)

	apiEnv := api.ApiEnv{
		WorkerManager: workerManager,
		Queries:       queries,
	}
	http.ListenAndServe(":8080", apiEnv.New())
}
