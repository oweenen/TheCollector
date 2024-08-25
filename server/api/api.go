package api

import (
	"TheCollectorDG/db"
	"TheCollectorDG/workerManager"
	"TheCollectorDG/workers/tasks"
	"context"
	"encoding/json"
	"net/http"
)

type ApiEnv struct {
	Queries       *db.Queries
	WorkerManager *workerManager.Manager
}

func (env ApiEnv) New() *http.ServeMux {
	router := http.NewServeMux()
	router.HandleFunc("/summoner/byPuuid", env.getSummonerByPuuid)
	router.HandleFunc("/summoner/byNameTag", env.getSummonerByNameTag)
	return router
}

func (env ApiEnv) getSummonerByPuuid(w http.ResponseWriter, r *http.Request) {
	puuid := r.URL.Query().Get("puuid")

	summoner, err := env.Queries.GetSummonerByPuuid(context.Background(), puuid)
	if err != nil {
		http.Error(w, "summoner not found", 404)
		return
	}

	bytes, _ := json.Marshal(summoner)
	w.Write(bytes)
}

func (env ApiEnv) getSummonerByNameTag(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	name := r.URL.Query().Get("name")
	tag := r.URL.Query().Get("tag")

	exists, err := env.Queries.SummonerExistsByNameTag(ctx, db.SummonerExistsByNameTagParams{
		Name: name,
		Tag:  tag,
	})
	if err != nil {
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
			return
		}

		<-done
	}

	summoner, err := env.Queries.GetSummonerByNameTag(ctx, db.GetSummonerByNameTagParams{
		Name: name,
		Tag:  tag,
	})
	if err != nil {
		return
	}

	bytes, _ := json.Marshal(summoner)
	w.Write(bytes)
}
