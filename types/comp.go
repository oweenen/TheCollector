package types

type Comp struct {
	SummonerPuuid     string   `json:"summoner_puuid,omitempty"`
	Placement         int      `json:"placement"`
	LastRound         int      `json:"last_round"`
	Level             int      `json:"level,omitempty"`
	RemainingGold     int      `json:"remaining_gold"`
	PlayersEliminated int      `json:"players_eliminated"`
	PlayerDamageDealt int      `json:"player_damage_dealt"`
	TimeEliminated    float32  `json:"time_eliminated"`
	Companion         int      `json:"companion"`
	Augments          []string `json:"augments"`
	Traits            []Trait  `json:"traits"`
	Units             []Unit   `json:"units"`
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
