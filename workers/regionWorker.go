package workers

import (
	"TheCollectorDG/db"
	"TheCollectorDG/tasks"
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func RegionWorker(pool *pgxpool.Pool, queries *db.Queries, queue chan tasks.Task) {
	for {
		select {
		case task := <-queue:
			err := task.Exec(context.Background())
			if err != nil {
				log.Println(err.Error())
			}
		default:
			spawnSummonerDetailsTasks(pool, queries, queue)
		}
	}
}

func spawnSummonerDetailsTasks(pool *pgxpool.Pool, queries *db.Queries, queue chan tasks.Task) {
	puuids, _ := queries.GetPuuidsWithNullSummoner(context.Background(), 5)
	for _, puuid := range puuids {
		queue <- tasks.SummonerDetailsTask{
			Puuid:   puuid,
			Region:  "na1",
			Queries: queries,
		}
	}
}
