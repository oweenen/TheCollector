package riot

import (
	"fmt"
)

type RiotMatchRes struct {
	MetaData struct {
		DataVersion  string   `json:"data_version"`
		MatchId      string   `json:"match_id"`
		Participants []string `json:"participants"`
	} `json:"metadata"`
	Info struct {
		Date      int64   `json:"game_datetime"`
		Length    float64 `json:"game_length"`
		Version   string  `json:"game_version"`
		QueueId   int     `json:"queue_id"`
		GameType  string  `json:"tft_game_type"`
		SetName   string  `json:"tft_set_core_name"`
		SetNumber int     `json:"tft_set_number"`
		Comps     []struct {
			Augments  []string `json:"augments"`
			Companion struct {
				ContentId string `json:"content_ID"`
				ItemId    int    `json:"item_ID"`
				SkinId    int    `json:"skin_ID"`
				Species   string `json:"species"`
			} `json:"companion"`
			RemainingGold     int     `json:"gold_left"`
			LastRound         int     `json:"last_round"`
			Level             int     `json:"level"`
			Placement         int     `json:"placement"`
			PlayersEliminated int     `json:"players_eliminated"`
			Puuid             string  `json:"puuid" Match:"Summoner"`
			TimeEliminated    float64 `json:"time_eliminated"`
			DamageToPlayers   int     `json:"total_damage_to_players"`
			Traits            []struct {
				Name       string `json:"name"`
				NumUnits   int    `json:"num_units"`
				Style      int    `json:"style"`
				TierActive int    `json:"tier_current"`
				TierMax    int    `json:"tier_total"`
			} `json:"traits"`
			Units []struct {
				Id        string   `json:"character_id"`
				ItemNames []string `json:"itemNames"`
				Rarity    int      `json:"rarity"`
				Tier      int      `json:"tier"`
			} `json:"units"`
		} `json:"participants"`
	} `json:"info"`
}

func GetMatchDetails(regionalServer, matchId string) (*RiotMatchRes, error) {
	matchRes := new(RiotMatchRes)
	url := fmt.Sprintf("https://%v.api.riotgames.com/tft/match/v1/matches/%v?api_key=%v", regionalServer, matchId, key)
	err := getJson(url, matchRes)
	if err != nil {
		return nil, err
	}

	return matchRes, err
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
