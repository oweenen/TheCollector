package riot

import (
	"fmt"
)

type RiotSummonerRes struct {
	Puuid         string `json:"puuid"`
	Id            string `json:"id"`
	AccountId     string `json:"accountId"`
	Name          string `json:"name"`
	ProfileIconId int    `json:"profileIconId"`
	RevisionDate  int    `json:"revisionDate"`
	SummonerLevel int    `json:"summonerLevel"`
}

func GetSummonerByName(region string, name string) (*RiotSummonerRes, error) {
	summonerRes := new(RiotSummonerRes)
	url := fmt.Sprintf("https://%v.api.riotgames.com/tft/summoner/v1/summoners/by-name/%v?api_key=%v", region, name, key)
	err := getJson(url, summonerRes)
	if err != nil {
		return nil, err
	}

	return summonerRes, err
}

func GetSummonerByPuuid(region string, puuid string) (*RiotSummonerRes, error) {
	summonerRes := new(RiotSummonerRes)
	url := fmt.Sprintf("https://%v.api.riotgames.com/tft/summoner/v1/summoners/by-puuid/%v?api_key=%v", region, puuid, key)
	err := getJson(url, summonerRes)
	if err != nil {
		return nil, err
	}

	return summonerRes, err
}
