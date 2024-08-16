package tasks

import (
	"TheCollectorDG/db"
	"TheCollectorDG/riot"
	"TheCollectorDG/workerManager"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type MatchHistoryTask struct {
	Cluster string
	Puuid   string
	Queue   chan workerManager.Task
	Pool    *pgxpool.Pool
	Queries *db.Queries
}

func (task MatchHistoryTask) Id() string {
	return fmt.Sprintf("MatchHistoryTask-%v", task.Puuid)
}

func (task MatchHistoryTask) Exec(ctx context.Context) error {
	updatedAt := time.Now()

	res, err := riot.GetMatchHistory(task.Cluster, task.Puuid)
	if err != nil {
		return err
	}

	task.Queries.SetBackgroundUpdateTimestamp(ctx, db.SetBackgroundUpdateTimestampParams{
		Puuid: task.Puuid,
		BackgroundUpdateTimestamp: pgtype.Timestamp{
			Time:  updatedAt,
			Valid: true,
		},
	})

	log.Printf("Got %v matchIds from summoner %v\n", len(res), task.Puuid)

	// queue match details
	for _, matchId := range res {
		task.Queue <- MatchDetailsTask{
			Cluster: task.Cluster,
			MatchId: matchId,
			Pool:    task.Pool,
			Queries: task.Queries,
		}
	}

	return err
}
