package types

import "strings"

type Summoner struct {
	Puuid         string `json:"puuid"`
	Region        string `json:"region"`
	SummonerId    string `json:"summoner_id"`
	Name          string `json:"name"`
	ProfileIconId int    `json:"profileIconId"`
	SummonerLevel int    `json:"summonerLevel"`
	LastUpdated   int64  `json:"lastUpdated"`
}

func ToRawName(displayName string) string {
	return strings.ToLower(strings.ReplaceAll(displayName, " ", ""))
}

func NewSummonerFromRiotRes(region string, summonerRes *RiotSummonerRes) *Summoner {
	return &Summoner{
		Puuid:         summonerRes.Puuid,
		Region:        region,
		SummonerId:    summonerRes.Id,
		Name:          summonerRes.Name,
		ProfileIconId: summonerRes.ProfileIconId,
		SummonerLevel: summonerRes.SummonerLevel,
	}
}
