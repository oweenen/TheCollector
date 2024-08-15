package riot

import (
	"fmt"
)

type Match struct {
	MetaData struct {
		DataVersion  string   `json:"data_version"`
		MatchId      string   `json:"match_id"`
		Participants []string `json:"participants"`
	} `json:"metadata"`
	Info struct {
		Date        int64   `json:"game_datetime"`
		Length      float64 `json:"game_length"`
		GameVersion string  `json:"game_version"`
		QueueId     int     `json:"queue_id"`
		GameType    string  `json:"tft_game_type"`
		SetName     string  `json:"tft_set_core_name"`
		SetNumber   int     `json:"tft_set_number"`
		Comps       []Comp  `json:"participants"`
	} `json:"info"`
}

type Comp struct {
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
	Puuid             string  `json:"puuid"`
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
		CharacterId string   `json:"character_id"`
		ItemNames   []string `json:"itemNames"`
		Rarity      int      `json:"rarity"`
		Tier        int      `json:"tier"`
	} `json:"units"`
}

func GetMatchDetails(cluster, matchId string) (*Match, error) {
	matchRes := new(Match)
	route := fmt.Sprintf("tft/match/v1/matches/%v", matchId)
	err := getJson(cluster, route, matchRes)
	if err != nil {
		return nil, err
	}

	return matchRes, err
}

func GetMatchHistory(cluster string, puuid string) ([]string, error) {
	var history []string
	count := 200

	route := fmt.Sprintf("tft/match/v1/matches/by-puuid/%v/ids?count=%v", puuid, count)
	err := getJson(cluster, route, &history)
	return history, err
}
