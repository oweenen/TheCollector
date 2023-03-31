package riot

import (
	"TheCollectorDG/types"
	"fmt"
)

func GetSummonerByName(region string, name string) (*types.Summoner, error) {
	summonerRes := new(types.RiotSummonerRes)
	url := fmt.Sprintf("https://%v.api.riotgames.com/tft/summoner/v1/summoners/by-name/%v?api_key=%v", region, name, key)
	err := getJson(url, summonerRes)
	if err != nil {
		return nil, err
	}
	summoner := types.NewSummonerFromRiotRes(region, summonerRes)
	return summoner, err
}

func GetSummonerByPuuid(region string, puuid string) (*types.Summoner, error) {
	summonerRes := new(types.RiotSummonerRes)
	url := fmt.Sprintf("https://%v.api.riotgames.com/tft/summoner/v1/summoners/by-puuid/%v?api_key=%v", region, puuid, key)
	err := getJson(url, summonerRes)
	if err != nil {
		return nil, err
	}
	summoner := types.NewSummonerFromRiotRes(region, summonerRes)
	return summoner, err
}
