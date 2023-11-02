package main

import (
	"TheCollectorDG/api"
	"TheCollectorDG/collection/matchCollection"
	"TheCollectorDG/collection/summonerCollection"
	"TheCollectorDG/database"
	"TheCollectorDG/datastore"
	"TheCollectorDG/riot"
	"TheCollectorDG/stats"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: Failed to load .env file")
	}

	database.SetupConnection()
	datastore.SetupConnection()
	riot.Setup()

	rateLimit, err := strconv.ParseFloat(os.Getenv("RIOT_RATE_LIMIT"), 32)
	if err != nil {
		rateLimit = 100
	}
	ratePeriod, err := strconv.ParseFloat(os.Getenv("RIOT_RATE_PERIOD"), 32)
	if err != nil {
		ratePeriod = 120000
	}
	rateEfficiency, err := strconv.ParseFloat(os.Getenv("RIOT_RATE_EFFICIENCY"), 32)
	if err != nil {
		rateEfficiency = 0.95
	}
	requestInterval := ratePeriod / rateLimit / rateEfficiency
	requestIntervalDuration := time.Duration(requestInterval) * time.Millisecond
	queueSpacing := time.Duration(requestInterval/float64(len(riot.RiotRegionRoutes)+len(riot.RiotRegionClusters))) * time.Millisecond

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

	go stats.AugmentStatsRefreshLoop(time.Hour)

	api.Setup(prioSummonerCollectionRouter, prioMatchCollectionRouter)
	api.Start()
}
