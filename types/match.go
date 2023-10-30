package types

import (
	"TheCollectorDG/riot"
	"strings"
)

type Match struct {
	Id          string  `json:"id"`
	Date        int64   `json:"date"`
	GameLength  float64 `json:"game_length"`
	GameVersion string  `json:"game_version"`
	QueueId     int     `json:"queue_id"`
	GameType    string  `json:"game_type"`
	SetName     string  `json:"set_name"`
	SetNumber   int     `json:"set_number"`
	Comps       []Comp  `json:"comps,omitempty"`
}

func NewMatchFromRiotRes(matchRes *riot.RiotMatchRes) *Match {
	match := &Match{
		Id:          matchRes.MetaData.MatchId,
		Date:        matchRes.Info.Date,
		GameLength:  matchRes.Info.Length,
		GameVersion: formatGameVersion(matchRes.Info.Version),
		QueueId:     matchRes.Info.QueueId,
		GameType:    matchRes.Info.GameType,
		SetName:     matchRes.Info.SetName,
		SetNumber:   matchRes.Info.SetNumber,
		Comps:       make([]Comp, len(matchRes.Info.Comps)),
	}

	for i, comp := range matchRes.Info.Comps {
		if len(comp.Augments) == expectedAugments(comp.LastRound)-1 {
			comp.Augments = append([]string{"TFT9_Augment_Legend"}, comp.Augments...)
		}

		match.Comps[i] = Comp{
			SummonerPuuid:     comp.Puuid,
			Placement:         comp.Placement,
			LastRound:         comp.LastRound,
			Level:             comp.Level,
			RemainingGold:     comp.RemainingGold,
			PlayersEliminated: comp.PlayersEliminated,
			PlayerDamageDealt: comp.DamageToPlayers,
			TimeEliminated:    float32(comp.TimeEliminated),
			Companion:         comp.Companion.ItemId,
			Augments:          comp.Augments,
			Units:             make([]Unit, len(comp.Units)),
		}

		for _, trait := range comp.Traits {
			if trait.TierActive == 0 {
				continue
			}
			match.Comps[i].Traits = append(match.Comps[i].Traits, Trait{
				Name:       trait.Name,
				Style:      trait.Style,
				TierActive: trait.TierActive,
				TierMax:    trait.TierMax,
			})
		}

		for j, unit := range comp.Units {
			match.Comps[i].Units[j] = Unit{
				CharactedId: unit.Id,
				Rarity:      unit.Rarity,
				Tier:        unit.Tier,
				Items:       unit.ItemNames,
			}
		}
	}

	return match
}

func GetMatchIdRegion(matchId string) string {
	return strings.ToLower(strings.Split(matchId, "_")[0])
}

func expectedAugments(lastRound int) int {
	if lastRound >= 20 {
		return 3
	} else if lastRound >= 13 {
		return 2
	} else if lastRound >= 5 {
		return 1
	}
	return 0
}

func formatGameVersion(gameVersion string) string {
	parts := strings.Split(gameVersion, " ")
	return parts[len(parts)-1]
}
