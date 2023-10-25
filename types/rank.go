package types

import "TheCollectorDG/riot"

type Rank struct {
	Tier     string `json:"tier"`
	Division string `json:"division"`
	Lp       int    `json:"lp"`
}

func NewRankFromRiotRes(rankRes *riot.RiotRankRes) *Rank {
	return &Rank{
		Tier:     rankRes.Tier,
		Division: rankRes.Rank,
		Lp:       rankRes.LeaguePoints,
	}
}

var tierToLp = map[string]int{
	"IRON":        0,
	"BRONZE":      400,
	"SILVER":      800,
	"GOLD":        1200,
	"PLATINUM":    1600,
	"DIAMOND":     2000,
	"MASTER":      2400,
	"GRANDMASTER": 2400,
	"CHALLENGER":  2400,
}

var divisionToLp = map[string]int{
	"I":   300,
	"II":  200,
	"III": 100,
	"IV":  0,
}

func calcRawLp(tier, division string, lp int) int {
	rawLp := lp

	tierLp, ok := tierToLp[tier]
	if ok {
		rawLp += tierLp
	}

	if tierLp < 2400 {
		if divisionLp, ok := divisionToLp[division]; ok {
			rawLp += divisionLp
		}
	}

	return rawLp
}
