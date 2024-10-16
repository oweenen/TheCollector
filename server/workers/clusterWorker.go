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

func (env WorkerEnv) ClusterWorker(prioQueue chan workerManager.Task) {
	backoffTicker := time.NewTicker(BACKOFF_TIME)
	queue := make(chan workerManager.Task, 1000)

	for {
		select {
		case task := <-prioQueue:
			err := task.Exec(context.Background())
			if err != nil {
				log.Println(err.Error())
			}
			continue
		default:
		}

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

	for i, row := range rows {
		select {
		case queue <- tasks.MatchHistoryTask{
			Cluster:      "americas",
			Puuid:        row.Puuid,
			Pool:         pool,
			Queries:      queries,
			MatchesAfter: time.Now().Add(-time.Hour * 24 * 3),
		}:
		default:
			return i + 1
		}
	}

	return len(rows)
}

func spawnAccountDetailsTasks(pool *pgxpool.Pool, queries *db.Queries, queue chan workerManager.Task) int {
	puuids, _ := queries.GetPuuidsWithNullAccountData(context.Background(), 100)
	for i, puuid := range puuids {
		select {
		case queue <- tasks.AccountByPuuidTask{
			Puuid:   puuid,
			Cluster: "americas",
			Queries: queries,
		}:
		default:
			return i + 1
		}
	}

	return len(puuids)
}
