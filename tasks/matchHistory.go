package tasks

import (
	"TheCollectorDG/db"
	"TheCollectorDG/riot"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type MatchHistoryTask struct {
	Cluster string
	Puuid   string
	Queue   chan Task
	Conn    *pgx.Conn
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

	task.Queries.SetMatchesUpdated(ctx, db.SetMatchesUpdatedParams{
		Puuid: task.Puuid,
		MatchesUpdated: pgtype.Timestamp{
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
			Conn:    task.Conn,
			Queries: task.Queries,
		}
	}

	return err
}
