package db

import (
	"TheCollectorDG/types"
	"context"
	"encoding/json"
)

type CreateCompParams struct {
	MatchID       string
	SummonerPuuid string
	CompData      types.CompData
}

func (q *Queries) CreateComp(ctx context.Context, arg CreateCompParams) error {
	compJson, _ := json.Marshal(types.CompData(arg.CompData))
	return q._createComp(ctx, _createCompParams{
		MatchID:       arg.MatchID,
		SummonerPuuid: arg.SummonerPuuid,
		CompData:      compJson,
	})
}
