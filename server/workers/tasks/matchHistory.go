package tasks

import (
	"TheCollectorDG/db"
	"TheCollectorDG/riot"
	"TheCollectorDG/types"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type MatchHistoryTask struct {
	Cluster      string
	Puuid        string
	MatchesAfter time.Time
	Pool         *pgxpool.Pool
	Queries      *db.Queries
}

func (task MatchHistoryTask) Id() string {
	return fmt.Sprintf("MatchHistoryTask-%v", task.Puuid)
}

func (task MatchHistoryTask) Exec(ctx context.Context) error {
	updatedAt := time.Now()

	res, err := riot.GetMatchHistory(task.Cluster, task.Puuid, task.MatchesAfter)
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

	for _, matchId := range res {
		if exists, _ := task.Queries.MatchExists(ctx, matchId); exists {
			log.Printf("Skipping match %v...\n", matchId)
			return nil
		}

		res, err := riot.GetMatchDetails(task.Cluster, matchId)
		if err != nil {
			return err
		}

		err = storeMatchDetails(ctx, task.Pool, task.Queries, res)

		log.Printf("Stored match %v!\n", matchId)
	}

	return err
}

func storeMatchDetails(ctx context.Context, pool *pgxpool.Pool, queries *db.Queries, matchDetails *riot.Match) error {
	var err error

	// insert participants
	for _, puuid := range matchDetails.MetaData.Participants {
		err = queries.InsertPuuid(ctx, puuid)
		if err != nil {
			return err
		}
	}

	tx, err := pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	qtx := queries.WithTx(tx)

	// create match
	err = qtx.CreateMatch(ctx, db.CreateMatchParams{
		ID:          matchDetails.MetaData.MatchId,
		DataVersion: matchDetails.MetaData.DataVersion,
		GameVersion: matchDetails.Info.GameVersion,
		QueueID:     matchDetails.Info.QueueId,
		GameType:    matchDetails.Info.GameType,
		SetName:     matchDetails.Info.SetName,
		SetNumber:   matchDetails.Info.SetNumber,
		MatchDate: pgtype.Timestamp{
			Time:  time.UnixMilli(matchDetails.Info.Date),
			Valid: true,
		},
	})
	if err != nil {
		return err
	}

	// create comps
	for _, compDetails := range matchDetails.Info.Comps {
		err = qtx.CreateComp(ctx, db.CreateCompParams{
			MatchID:       matchDetails.MetaData.MatchId,
			SummonerPuuid: compDetails.Puuid,
			CompData:      types.CompData(compDetails),
		})
		if err != nil {
			return err
		}
	}

	err = tx.Commit(ctx)

	return err
}
