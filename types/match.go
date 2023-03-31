package types

import (
	"strings"
)

type Match struct {
	Id          string  `json:"id,omitempty"`
	Date        int64   `json:"date,omitempty"`
	GameLength  float64 `json:"game_length,omitempty"`
	GameVersion string  `json:"game_version,omitempty"`
	QueueId     int     `json:"queue_id,omitempty"`
	GameType    string  `json:"game_type,omitempty"`
	SetName     string  `json:"set_name,omitempty"`
	SetNumber   int     `json:"set_number,omitempty"`
	Comps       []Comp  `json:"comps,omitempty"`
}

type Comp struct {
	Match             *Match    `json:"match,omitempty"`
	Summoner          *Summoner `json:"summoner,omitempty"`
	Placement         int       `json:"placement,omitempty"`
	LastRound         int       `json:"last_round,omitempty"`
	Level             int       `json:"level,omitempty"`
	RemainingGold     int       `json:"remaining_gold,omitempty"`
	PlayersEliminated int       `json:"players_eliminated,omitempty"`
	PlayerDamageDealt int       `json:"player_damage_dealt,omitempty"`
	TimeEliminated    float32   `json:"time_eliminated,omitempty"`
	Companion         Companion `json:"companion,omitempty"`
	Augments          []string  `json:"augments,omitempty"`
	Traits            []Trait   `json:"traits,omitempty"`
	Units             []Unit    `json:"units,omitempty"`
}

type Companion struct {
	ContentId string `json:"content_id,omitempty"`
	ItemId    int    `json:"item_id,omitempty"`
	SkinId    int    `json:"skin_id,omitempty"`
	Species   string `json:"species,omitempty"`
}

type Trait struct {
	Name       string `json:"name"`
	Style      int    `json:"style"`
	TierActive int    `json:"tier_active"`
	TierMax    int    `json:"tier_max"`
}

type Unit struct {
	CharactedId string   `json:"character_id"`
	Rarity      int      `json:"rarity"`
	Tier        int      `json:"tier"`
	Items       []string `json:"items"`
}

func GetMatchIdRegion(matchId string) string {
	return strings.Split(matchId, "_")[0]
}

func NewMatchFromRiotRes(matchRes *RiotMatchRes) *Match {
	match := &Match{
		Id:          matchRes.MetaData.MatchId,
		Date:        matchRes.Info.Date,
		GameLength:  matchRes.Info.Length,
		GameVersion: matchRes.Info.Version,
		QueueId:     matchRes.Info.QueueId,
		GameType:    matchRes.Info.GameType,
		SetName:     matchRes.Info.SetName,
		SetNumber:   matchRes.Info.SetNumber,
		Comps:       make([]Comp, len(matchRes.Info.Comps)),
	}

	for i, comp := range matchRes.Info.Comps {
		match.Comps[i] = Comp{
			Summoner: &Summoner{
				Puuid: comp.Puuid,
			},
			Placement:         comp.Placement,
			LastRound:         comp.LastRound,
			Level:             comp.Level,
			RemainingGold:     comp.RemainingGold,
			PlayersEliminated: comp.PlayersEliminated,
			PlayerDamageDealt: comp.DamageToPlayers,
			TimeEliminated:    float32(comp.TimeEliminated),
			Companion: Companion{
				ContentId: comp.Companion.ContentId,
				ItemId:    comp.Companion.ItemId,
				SkinId:    comp.Companion.SkinId,
				Species:   comp.Companion.Species,
			},
			Augments: comp.Augments,
			Traits:   make([]Trait, len(comp.Traits)),
			Units:    make([]Unit, len(comp.Units)),
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
