package types

import "strings"

type Summoner struct {
	Puuid         string `json:"puuid,omitempty"`
	Region        string `json:"region,omitempty"`
	SummonerId    string `json:"summoner_id,omitempty"`
	Name          string `json:"name,omitempty"`
	ProfileIconId int    `json:"profileIconId,omitempty"`
	SummonerLevel int    `json:"summonerLevel,omitempty"`
	LastUpdated   int64  `json:"lastUpdated,omitempty"`
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
