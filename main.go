package main

import (
	"TheCollectorDG/api"
	"TheCollectorDG/collection/matchCollection"
	"TheCollectorDG/collection/summonerCollection"
	"TheCollectorDG/config"
	"TheCollectorDG/database"
	"TheCollectorDG/riot"
	"log"
	"time"
)

func main() {
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatalln(err)
	}

	riot.Setup(config.Riot.Key, config.Riot.MatchesAfter)
	database.Setup(config.MySqlConfig)

	prioSummonerCQ := summonerCollection.NewSummonerCollectionQueue()
	prioMatchCQ := matchCollection.NewMatchCollectionQueue(prioSummonerCQ)
	summonerCQ := summonerCollection.NewSummonerCollectionQueue()
	matchCq := matchCollection.NewMatchCollectionQueue(summonerCQ)

	requestInterval := time.Duration(float32(config.Riot.RatePeriod)/float32(config.Riot.RateLimit)/config.Riot.RateEfficiency) * time.Millisecond

	go summonerCollection.SummonerCollectionLoop(prioSummonerCQ, summonerCQ, requestInterval)
	time.Sleep(requestInterval / 2)
	go matchCollection.MatchCollectionLoop(prioMatchCQ, matchCq, requestInterval)

	api.Setup(prioSummonerCQ, prioMatchCQ)
	api.Start()
}
