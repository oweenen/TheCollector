package tasks

import (
	"TheCollectorDG/db"
	"TheCollectorDG/riot"
	"context"
	"fmt"
	"log"
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
		Puuid:         task.Puuid,
		SummonerID:    res.SummonerId,
		ProfileIconID: res.ProfileIconId,
		SummonerLevel: res.SummonerLevel,
	})

	log.Printf("Summoner details collected for %v\n", task.Puuid)

	return err
}
