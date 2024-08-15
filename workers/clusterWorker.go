package workers

import (
	"TheCollectorDG/db"
	"TheCollectorDG/tasks"
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func ClusterWorker(pool *pgxpool.Pool, queries *db.Queries, queue chan tasks.Task) {
	for {
		select {
		case task := <-queue:
			err := task.Exec(context.Background())
			if err != nil {
				log.Println(err.Error())
			}
		default:
			spawnMatchHistoryTasks(pool, queries, queue)
		}
	}
}

func spawnMatchHistoryTasks(pool *pgxpool.Pool, queries *db.Queries, queue chan tasks.Task) {
	rows, _ := queries.GetOldestMatchHistories(context.Background(), 1)

	for _, row := range rows {
		queue <- tasks.MatchHistoryTask{
			Cluster: "americas",
			Puuid:   row.Puuid,
			Queue:   queue,
			Pool:    pool,
			Queries: queries,
		}
	}
}
