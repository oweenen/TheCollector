package workers

import (
	"TheCollectorDG/db"
	"TheCollectorDG/tasks"
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func ClusterWorker(pool *pgxpool.Pool, queries *db.Queries, queue chan tasks.Task) {
	backoffTicker := time.NewTicker(time.Second * 10)

	for {
		select {
		case task := <-queue:
			err := task.Exec(context.Background())
			if err != nil {
				log.Println(err.Error())
			}
			continue
		default:
		}

		select {
		case <-backoffTicker.C:
			spawnAccountDetailsTasks(pool, queries, queue)
			spawnMatchHistoryTasks(pool, queries, queue)
		default:
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

func spawnAccountDetailsTasks(pool *pgxpool.Pool, queries *db.Queries, queue chan tasks.Task) {
	puuids, _ := queries.GetPuuidsWithNullAccountData(context.Background(), 100)
	for _, puuid := range puuids {
		queue <- tasks.AccountDetailsTask{
			Puuid:   puuid,
			Cluster: "americas",
			Queries: queries,
		}
	}
}
