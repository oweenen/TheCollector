package types

type CompData struct {
	Augments  []string `json:"augments"`
	Companion struct {
		ContentId string `json:"contentId"`
		ItemId    int    `json:"itemId"`
		SkinId    int    `json:"skinId"`
		Species   string `json:"species"`
	} `json:"companion"`
	RemainingGold     int     `json:"goldLeft"`
	LastRound         int     `json:"lastRound"`
	Level             int     `json:"level"`
	Placement         int     `json:"placement"`
	PlayersEliminated int     `json:"playersEliminated"`
	Puuid             string  `json:"puuid"`
	TimeEliminated    float64 `json:"timeEliminated"`
	DamageToPlayers   int     `json:"totalDamageToPlayers"`
	Traits            []struct {
		Name       string `json:"name"`
		NumUnits   int    `json:"numUnits"`
		Style      int    `json:"style"`
		TierActive int    `json:"tierRurrent"`
		TierMax    int    `json:"tierTotal"`
	} `json:"traits"`
	Units []struct {
		CharacterId string   `json:"characterId"`
		ItemNames   []string `json:"itemNames"`
		Rarity      int      `json:"rarity"`
		Tier        int      `json:"tier"`
	} `json:"units"`
}
