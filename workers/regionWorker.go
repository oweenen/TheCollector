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

func (env WorkerEnv) RegionWorker(queue chan workerManager.Task) {
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
			spawnSummonerDetailsTasks(env.Pool, env.Queries, queue)
		default:
		}
	}
}

func spawnSummonerDetailsTasks(pool *pgxpool.Pool, queries *db.Queries, queue chan workerManager.Task) int {
	puuids, _ := queries.GetPuuidsWithNullSummonerData(context.Background(), 100)
	for i, puuid := range puuids {
		select {
		case queue <- tasks.SummonerDetailsTask{
			Puuid:   puuid,
			Region:  "na1",
			Queries: queries,
		}:
		default:
			return i + 1
		}
	}

	return len(puuids)
}
