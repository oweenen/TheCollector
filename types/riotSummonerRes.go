package types

type RiotSummonerRes struct {
	Puuid         string `json:"puuid"`
	Id            string `json:"id"`
	AccountId     string `json:"accountId"`
	Name          string `json:"name"`
	ProfileIconId int    `json:"profileIconId"`
	RevisionDate  int    `json:"revisionDate"`
	SummonerLevel int    `json:"summonerLevel"`
}
