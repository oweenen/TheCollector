package tasks

import (
	"TheCollectorDG/db"
	"TheCollectorDG/riot"
	"TheCollectorDG/types"
	"context"
	"log"

	"github.com/jackc/pgx/v5"
)

type MatchDetailsTask struct {
	Cluster string
	MatchId string
	Conn    *pgx.Conn
	Queries *db.Queries
}

func (task MatchDetailsTask) Exec(ctx context.Context) error {
	if exists, _ := task.Queries.MatchExists(ctx, task.MatchId); exists {
		log.Printf("Skipping match %v...\n", task.MatchId)
		return nil
	}

	res, err := riot.GetMatchDetails(task.Cluster, task.MatchId)
	if err != nil {
		return err
	}

	err = storeMatchDetails(ctx, task.Conn, task.Queries, res)

	log.Printf("Stored match %v!\n", task.MatchId)

	return err
}

func (task MatchDetailsTask) Id() string {
	return task.MatchId
}

func storeMatchDetails(ctx context.Context, conn *pgx.Conn, queries *db.Queries, matchDetails *riot.Match) error {
	var err error

	// insert participants
	for _, puuid := range matchDetails.MetaData.Participants {
		err = queries.InsertPuuid(ctx, puuid)
		if err != nil {
			return err
		}
	}

	tx, err := conn.Begin(ctx)
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