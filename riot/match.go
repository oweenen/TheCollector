package riot

import (
	"TheCollectorDG/types"
	"fmt"
	"strings"
)

func getRoutingServer(region string) (string, bool) {
	region = strings.ToLower(region)

	switch region {
	case "na1", "br1", "la1", "la2":
		return "americas", true
	case "kr", "jp1":
		return "asia", true
	case "eun1", "euw1", "tr1", "ru":
		return "europe", true
	case "ph2", "sg2", "th2", "tw2", "vn2", "oc1":
		return "sea", true
	default:
		return "unknown", false
	}
}

func GetMatchDetails(matchId string) (*types.Match, error) {
	matchRes := new(types.RiotMatchRes)
	region := types.GetMatchIdRegion(matchId)
	server, ok := getRoutingServer(region)
	if !ok {
		return nil, fmt.Errorf("invalid region %v", region)
	}
	url := fmt.Sprintf("https://%v.api.riotgames.com/tft/match/v1/matches/%v?api_key=%v", server, matchId, key)
	err := getJson(url, matchRes)
	if err != nil {
		return nil, err
	}
	match := types.NewMatchFromRiotRes(matchRes)
	return match, err
}

func GetMatchHistory(region string, puuid string, after int64) ([]string, error) {
	var history []string
	server, ok := getRoutingServer(region)
	if !ok {
		return nil, fmt.Errorf("invalid region %v", region)
	}
	count := 200
	if matchesAfter > after {
		after = matchesAfter
	}

	url := fmt.Sprintf("https://%v.api.riotgames.com/tft/match/v1/matches/by-puuid/%v/ids?startTime=%v&count=%v&api_key=%v", server, puuid, after, count, key)
	err := getJson(url, &history)
	return history, err
}
