package workers

import (
	"TheCollectorDG/db"
	"TheCollectorDG/workerManager"
	"TheCollectorDG/workers/tasks"

	"context"
	"log"
	"time"
)

func (env WorkerEnv) RegionWorker(prioQueue chan workerManager.Task) {
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
			spawnSummonerDetailsTasks(env.Queries, queue)
		default:
		}
	}
}

func spawnSummonerDetailsTasks(queries *db.Queries, queue chan workerManager.Task) int {
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
