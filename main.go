package main

import (
	"TheCollectorDG/db"
	"TheCollectorDG/tasks"
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

var conn *pgx.Conn
var queries *db.Queries

func init() {
	var err error

	ctx := context.Background()

	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	log.Println("Env variables loaded")

	conn, err = pgx.Connect(ctx, os.Getenv("DB_URL"))
	if err != nil {
		panic(err)
	}
	log.Println("Db connection successful")

	queries = db.New(conn)
}

func main() {
	queue := make(chan tasks.Task, 1000)

	worker(queue)
}

func worker(queue chan tasks.Task) {
	for {
		select {
		case task := <-queue:
			task.Exec(context.Background())
		default:
			spawnMatchHistoryTasks(queue)
		}
	}
}

func spawnMatchHistoryTasks(queue chan tasks.Task) {
	rows, _ := queries.GetOldestMatchHistories(context.Background(), 1)

	for _, row := range rows {
		queue <- tasks.MatchHistoryTask{
			Cluster: "americas",
			Puuid:   row.Puuid,
			Queue:   queue,
			Conn:    conn,
			Queries: queries,
		}
	}
}
