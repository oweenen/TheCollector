package api

import (
	"TheCollectorDG/collection/matchCollection"
	"TheCollectorDG/collection/summonerCollection"

	"github.com/gofiber/fiber/v2"
)

var summonerCollectionQueue *summonerCollection.SummonerCollectionQueue
var matchCollectionQueue *matchCollection.MatchCollectionQueue

func Setup(summonerCQ *summonerCollection.SummonerCollectionQueue, matchCQ *matchCollection.MatchCollectionQueue) {
	summonerCollectionQueue = summonerCQ
	matchCollectionQueue = matchCQ
}

func Start() {
	app := fiber.New()

	app.Get("summoner/:region/:name", GetSummoner)
	app.Get("comps/:puuid", GetCompHistory)
	app.Get("update/profile/:puuid", UpdateProfile)

	app.Listen(":9090")
}
