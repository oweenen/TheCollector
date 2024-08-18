package api

import (
	"TheCollectorDG/db"
	"context"
	"encoding/json"
	"net/http"
)

type ApiEnv struct {
	Queries *db.Queries
}

func (env ApiEnv) New() *http.ServeMux {
	router := http.NewServeMux()
	router.HandleFunc("/summoner/{puuid}", env.getSummoner)
	return router
}

func (env ApiEnv) getSummoner(w http.ResponseWriter, r *http.Request) {
	puuid := r.PathValue("puuid")

	summoner, err := env.Queries.GetSummonerByPuuid(context.Background(), puuid)
	if err != nil {
		http.Error(w, "summoner not found", 404)
		return
	}

	bytes, _ := json.Marshal(summoner)
	w.Write(bytes)
}
