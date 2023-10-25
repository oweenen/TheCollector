package summonerCollection

import (
	"TheCollectorDG/database"
	"TheCollectorDG/riot"
	"TheCollectorDG/types"
	"fmt"
	"time"
)

type RankCollecter struct {
	Region     string
	SummonerId string
	Puuid      string
}

func NewRankCollecter(region string, puuid string, summonerId string) RankCollecter {
	return RankCollecter{
		Region:     region,
		SummonerId: summonerId,
		Puuid:      puuid,
	}
}

func (c RankCollecter) Id() string {
	return c.Puuid
}

func (c RankCollecter) Collect() error {
	rankRes, err := riot.GetRank(c.Region, c.SummonerId)
	if err != nil {
		fmt.Printf("Error getting summoner %s from riot: %s\n", c.SummonerId, err)
		return err
	}
	rank := types.NewRankFromRiotRes(rankRes)

	if rank != nil {
		err = database.StoreRank(c.Puuid, rank)
		if err != nil {
			fmt.Printf("Error inserting rank %s into db %s\n", c.SummonerId, err)
			return err
		}
	}

	err = database.SetRankUpdatedAt(c.Puuid, time.Now().Unix())
	if err != nil {
		fmt.Printf("Error updating rank_last_updated for %s in db %s\n", c.Puuid, err)
		return err
	}

	fmt.Printf("Collected rank for summoner %v\n", c.Puuid)
	return nil
}
