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

	requestInterval := float32(config.Riot.RatePeriod) / float32(config.Riot.RateLimit) / config.Riot.RateEfficiency
	requestIntervalDuration := time.Duration(requestInterval) * time.Millisecond
	queueSpacing := time.Duration(requestInterval/float32(len(riot.RiotRegionRoutes)+len(riot.RiotRegionClusters))) * time.Millisecond

	summonerCollectionRouter := make(map[string]*summonerCollection.RegionalSummonerCollectionQueue)
	prioSummonerCollectionRouter := make(map[string]*summonerCollection.RegionalSummonerCollectionQueue)
	for region := range riot.RiotRegionRoutes {
		summonerCollectionQueue := summonerCollection.NewRegionalSummonerCollectionQueue(region)
		summonerCollectionRouter[region] = summonerCollectionQueue
		prioSummonerCollectionQueue := summonerCollection.NewRegionalSummonerCollectionQueue(region)
		prioSummonerCollectionRouter[region] = prioSummonerCollectionQueue
		go summonerCollection.SummonerCollectionLoop(prioSummonerCollectionQueue, summonerCollectionQueue, requestIntervalDuration)
		time.Sleep(queueSpacing)
	}

	matchCollectionRouter := make(map[string]*matchCollection.RegionalMatchCollectionQueue)
	prioMatchCollectionRouter := make(map[string]*matchCollection.RegionalMatchCollectionQueue)
	for regionalServer := range riot.RiotRegionClusters {
		matchCollectionQueue := matchCollection.NewRegionalMatchCollectionQueue(regionalServer, summonerCollectionRouter)
		matchCollectionRouter[regionalServer] = matchCollectionQueue
		prioMatchCollectionQueue := matchCollection.NewRegionalMatchCollectionQueue(regionalServer, prioSummonerCollectionRouter)
		prioMatchCollectionRouter[regionalServer] = prioMatchCollectionQueue
		go matchCollection.MatchCollectionLoop(prioMatchCollectionQueue, matchCollectionQueue, requestIntervalDuration)
		time.Sleep(queueSpacing)
	}

	api.Setup(prioSummonerCollectionRouter, prioMatchCollectionRouter)
	api.Start()
}
