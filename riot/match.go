package riot

import (
	"TheCollectorDG/types"
	"fmt"
)

func GetMatchDetails(regionalServer, matchId string) (*types.Match, error) {
	matchRes := new(types.RiotMatchRes)
	url := fmt.Sprintf("https://%v.api.riotgames.com/tft/match/v1/matches/%v?api_key=%v", regionalServer, matchId, key)
	err := getJson(url, matchRes)
	if err != nil {
		return nil, err
	}
	match := types.NewMatchFromRiotRes(matchRes)
	return match, err
}

func GetMatchHistory(regionalServer string, puuid string, after int64) ([]string, error) {
	var history []string
	count := 200
	if matchesAfter > after {
		after = matchesAfter
	}

	url := fmt.Sprintf("https://%v.api.riotgames.com/tft/match/v1/matches/by-puuid/%v/ids?startTime=%v&count=%v&api_key=%v", regionalServer, puuid, after, count, key)
	err := getJson(url, &history)
	return history, err
}
