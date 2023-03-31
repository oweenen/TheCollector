package summonerCollection

import (
	"TheCollectorDG/database"
	"TheCollectorDG/riot"
	"fmt"
)

type SummonerByPuuidCollecter struct {
	Puuid  string
	Region string
}

func NewSummonerByPuuidCollecter(region string, puuid string) SummonerByPuuidCollecter {
	return SummonerByPuuidCollecter{
		Puuid:  puuid,
		Region: region,
	}
}

func (c SummonerByPuuidCollecter) Id() string {
	return c.Puuid
}

func (c SummonerByPuuidCollecter) Collect() error {
	fmt.Printf("Collecting summoner %v\n", c.Puuid)

	summoner, err := riot.GetSummonerByPuuid(c.Region, c.Puuid)
	if err != nil {
		fmt.Printf("Error getting summoner %s from riot: %s\n", c.Puuid, err)
		return err
	}

	err = database.StoreSummoner(summoner)
	if err != nil {
		fmt.Printf("Error inserting summoner %s into db %s\n", c.Puuid, err)
		return err
	}

	return nil
}
