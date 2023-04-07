package types

type Rank struct {
	Type     string `json:"type"`
	Tier     string `json:"tier"`
	Division string `json:"division"`
	Lp       int    `json:"lp"`
}

func NewRankFromRiotRes(rankRes RiotRankRes) *Rank {
	return &Rank{
		Type:     rankRes.QueueType,
		Tier:     rankRes.Tier,
		Division: rankRes.Rank,
		Lp:       rankRes.LeaguePoints,
	}
}
