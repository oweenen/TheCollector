package tasks

import (
	"TheCollectorDG/db"
	"TheCollectorDG/riot"
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgtype"
)

type SummonerDetailsTask struct {
	Puuid   string
	Region  string
	Queries *db.Queries
}

func (task SummonerDetailsTask) Id() string {
	return fmt.Sprintf("SummonerDetailsTask-%v", task.Puuid)
}

func (task SummonerDetailsTask) Exec(ctx context.Context) error {
	res, err := riot.GetSummonerByPuuid(task.Region, task.Puuid)
	if err != nil {
		return err
	}

	err = task.Queries.UpdateSummoner(ctx, db.UpdateSummonerParams{
		Puuid: task.Puuid,
		SummonerID: pgtype.Text{
			String: res.SummonerId,
			Valid:  true,
		},
		ProfileIconID: pgtype.Int4{
			Int32: int32(res.ProfileIconId),
			Valid: true,
		},
		SummonerLevel: pgtype.Int4{
			Int32: int32(res.SummonerLevel),
			Valid: true,
		},
	})

	log.Printf("Summoner details collected for %v\n", task.Puuid)

	return err
}
