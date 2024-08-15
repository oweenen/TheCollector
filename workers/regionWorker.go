package workers

import (
	"TheCollectorDG/db"
	"TheCollectorDG/tasks"
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func RegionWorker(pool *pgxpool.Pool, queries *db.Queries, queue chan tasks.Task) {
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
			spawnSummonerDetailsTasks(pool, queries, queue)
		default:
		}
	}
}

func spawnSummonerDetailsTasks(pool *pgxpool.Pool, queries *db.Queries, queue chan tasks.Task) {
	puuids, _ := queries.GetPuuidsWithNullSummonerData(context.Background(), 100)
	for _, puuid := range puuids {
		queue <- tasks.SummonerDetailsTask{
			Puuid:   puuid,
			Region:  "na1",
			Queries: queries,
		}
	}
}
