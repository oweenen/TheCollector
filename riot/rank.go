package riot

import (
	"TheCollectorDG/types"
	"fmt"
)

func GetRank(region string, summonerId string) (*types.Rank, error) {
	var rankRes []types.RiotRankRes
	url := fmt.Sprintf("https://%v.api.riotgames.com/tft/league/v1/entries/by-summoner/%v?api_key=%v", region, summonerId, key)
	err := getJson(url, &rankRes)
	if err != nil {
		return nil, err
	}

	for _, r := range rankRes {
		if r.QueueType == "RANKED_TFT" {
			return types.NewRankFromRiotRes(r), nil
		}
	}
	return nil, nil
}
