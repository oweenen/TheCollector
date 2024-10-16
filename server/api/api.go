package api

import (
	"TheCollectorDG/db"
	"TheCollectorDG/workerManager"
	"TheCollectorDG/workers/tasks"
	"context"
	"encoding/json"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ApiEnv struct {
	Queries       *db.Queries
	Pool          *pgxpool.Pool
	WorkerManager *workerManager.Manager
}

func (env ApiEnv) New() *http.ServeMux {
	router := http.NewServeMux()
	router.HandleFunc("/summoner/{puuid}", env.getSummonerByPuuid)
	router.HandleFunc("/account/{name}/{tag}", env.getSummonerByNameTag)
	router.HandleFunc("/summoner/{puuid}/matches", env.getSummonerMatches)
	router.HandleFunc("/summoner/{puuid}/refresh", env.refreshSummoner)
	return router
}

func (env ApiEnv) getSummonerByPuuid(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	puuid := r.PathValue("puuid")

	summoner, err := env.Queries.GetSummonerByPuuid(context.Background(), puuid)
	if err != nil {
		http.Error(w, "summoner not found", 404)
		return
	}

	bytes, _ := json.Marshal(summoner)
	w.Write(bytes)
}

func (env ApiEnv) getSummonerByNameTag(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	ctx := context.Background()
	name := r.PathValue("name")
	tag := r.PathValue("tag")

	exists, err := env.Queries.SummonerExistsByNameTag(ctx, db.SummonerExistsByNameTagParams{
		Name: name,
		Tag:  tag,
	})
	if err != nil {
		w.WriteHeader(500)
		return
	}

	if !exists {
		done, err := env.WorkerManager.AssignTaskWithDone("americas", tasks.AccountByNameTagTask{
			Name:    name,
			Tag:     tag,
			Cluster: "americas",
			Queries: env.Queries,
		})
		if err != nil {
			w.WriteHeader(500)
			return
		}

		<-done
	}

	summoner, err := env.Queries.GetSummonerByNameTag(ctx, db.GetSummonerByNameTagParams{
		Name: name,
		Tag:  tag,
	})
	if err != nil {
		w.WriteHeader(404)
		return
	}

	bytes, _ := json.Marshal(summoner)
	w.Write(bytes)
}

func (env ApiEnv) getSummonerMatches(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	puuid := r.PathValue("puuid")

	matches, err := env.Queries.SummonerMatchHistory(context.Background(), db.SummonerMatchHistoryParams{
		SummonerPuuid: puuid,
		Limit:         20,
		Offset:        0,
	})
	if err != nil {
		http.Error(w, "summoner not found", 404)
		return
	}

	// prevent returning null when list is empty
	if matches == nil {
		matches = []db.SummonerMatchHistoryRow{}
	}

	bytes, _ := json.Marshal(matches)
	w.Write(bytes)
}

func (env ApiEnv) refreshSummoner(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	puuid := r.PathValue("puuid")

	env.WorkerManager.AssignTask("na1", tasks.SummonerDetailsTask{
		Puuid:   puuid,
		Region:  "na1",
		Queries: env.Queries,
	})

	env.WorkerManager.AssignTask("americas", tasks.AccountByPuuidTask{
		Puuid:   puuid,
		Cluster: "americas",
		Queries: env.Queries,
	})

	env.WorkerManager.AssignTask("americas", tasks.MatchHistoryTask{
		Puuid:   puuid,
		Pool:    env.Pool,
		Cluster: "americas",
		Queries: env.Queries,
	})
}
