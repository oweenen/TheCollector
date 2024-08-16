package workers

import (
	"TheCollectorDG/db"
	"TheCollectorDG/workerManager"
	"TheCollectorDG/workers/tasks"
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func (env WorkerEnv) ClusterWorker(queue chan workerManager.Task) {
	backoffTicker := time.NewTicker(BACKOFF_TIME)

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
			numSpawned := spawnAccountDetailsTasks(env.Pool, env.Queries, queue)
			if numSpawned == 0 {
				spawnMatchHistoryTasks(env.Pool, env.Queries, queue)
			}
		default:
		}
	}
}

func spawnMatchHistoryTasks(pool *pgxpool.Pool, queries *db.Queries, queue chan workerManager.Task) int {
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

	return len(rows)
}

func spawnAccountDetailsTasks(pool *pgxpool.Pool, queries *db.Queries, queue chan workerManager.Task) int {
	puuids, _ := queries.GetPuuidsWithNullAccountData(context.Background(), 100)
	for _, puuid := range puuids {
		queue <- tasks.AccountDetailsTask{
			Puuid:   puuid,
			Cluster: "americas",
			Queries: queries,
		}
	}

	return len(puuids)
}
