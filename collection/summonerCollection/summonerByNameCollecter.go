package summonerCollection

import (
	"TheCollectorDG/database"
	"TheCollectorDG/riot"
	"TheCollectorDG/types"
	"fmt"
	"time"
)

type SummonerByNameCollecter struct {
	RawName string
	Region  string
}

func NewSummonerByNameCollecter(region string, name string) SummonerByNameCollecter {
	return SummonerByNameCollecter{
		RawName: types.ToRawName(name),
		Region:  region,
	}
}

func (c SummonerByNameCollecter) Id() string {
	return c.Region + c.RawName
}

func (c SummonerByNameCollecter) Collect() error {
	// get summoner from riot
	updatedAt := time.Now().Unix()
	summoner, err := riot.GetSummonerByName(c.Region, c.RawName)
	if err != nil {
		fmt.Printf("Error getting summoner %s from riot: %s\n", c.RawName, err)
		return err
	}

	summoner.LastUpdated = updatedAt

	err = database.StoreSummoner(summoner)
	if err != nil {
		fmt.Printf("Error inserting summoner %s into db %s\n", c.RawName, err)
		return err
	}

	fmt.Printf("Collected summoner %v\n", c.RawName)
	return nil
}
