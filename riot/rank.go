package riot

import (
	"fmt"
)

type RiotRankRes struct {
	LeagueId     string `json:"leagueId"`
	QueueType    string `json:"queueType"`
	Tier         string `json:"tier"`
	Rank         string `json:"rank"`
	SummonerId   string `json:"summonerId"`
	SummonerName string `json:"summonerName"`
	LeaguePoints int    `json:"leaguePoints"`
	Wins         int    `json:"wins"`
	Losses       int    `json:"losses"`
	Veteran      bool   `json:"veteran"`
	Inactive     bool   `json:"inactive"`
	FreshBlood   bool   `json:"freshBlood"`
	HotStreak    bool   `json:"hotStreak"`
}

func GetRank(region string, summonerId string) (*RiotRankRes, error) {
	var rankRes []RiotRankRes
	url := fmt.Sprintf("https://%v.api.riotgames.com/tft/league/v1/entries/by-summoner/%v?api_key=%v", region, summonerId, key)
	err := getJson(url, &rankRes)
	if err != nil {
		return nil, err
	}

	for _, rank := range rankRes {
		if rank.QueueType == "RANKED_TFT" {
			return &rank, nil
		}
	}
	return nil, nil
}
