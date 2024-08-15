package riot

import (
	"fmt"
)

type RiotSummonerRes struct {
	Puuid         string `json:"puuid"`
	SummonerId    string `json:"id"`
	AccountId     string `json:"accountId"`
	ProfileIconId int    `json:"profileIconId"`
	RevisionDate  int    `json:"revisionDate"`
	SummonerLevel int    `json:"summonerLevel"`
}

func GetSummonerByPuuid(region string, puuid string) (*RiotSummonerRes, error) {
	summonerRes := new(RiotSummonerRes)
	route := fmt.Sprintf("tft/summoner/v1/summoners/by-puuid/%v", puuid)
	err := getJson(region, route, summonerRes)
	if err != nil {
		return nil, err
	}

	return summonerRes, err
}
