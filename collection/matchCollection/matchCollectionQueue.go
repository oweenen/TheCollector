package matchCollection

import (
	"TheCollectorDG/collection"
	"TheCollectorDG/collection/summonerCollection"
	"TheCollectorDG/types"
	"fmt"
)

type RegionalMatchCollectionQueue struct {
	regionalServer                 string
	matchDetailsCollectionQueue    collection.CollectionQueue
	matchHistoryCollectionQueue    collection.CollectionQueue
	summonerCollectionRegionRouter map[string]*summonerCollection.RegionalSummonerCollectionQueue
}

func NewRegionalMatchCollectionQueue(regionalServer string, summonerCollectionRegionRouter map[string]*summonerCollection.RegionalSummonerCollectionQueue) *RegionalMatchCollectionQueue {
	return &RegionalMatchCollectionQueue{
		regionalServer:                 regionalServer,
		matchDetailsCollectionQueue:    collection.NewCollectionQueue(),
		matchHistoryCollectionQueue:    collection.NewCollectionQueue(),
		summonerCollectionRegionRouter: summonerCollectionRegionRouter,
	}
}

func (cq *RegionalMatchCollectionQueue) QueueMatchDetails(matchId string) chan error {
	region := types.GetMatchIdRegion(matchId)
	summonerCollectionQueue, ok := cq.summonerCollectionRegionRouter[region]
	if !ok {
		errChan := make(chan error)
		errChan <- fmt.Errorf("error routing match details collection for %v no routing value for region %v", matchId, region)
		return errChan
	}
	collecter := NewMatchDetailsCollecter(cq.regionalServer, matchId, summonerCollectionQueue)
	return cq.matchDetailsCollectionQueue.Queue(collecter)
}

func (cq *RegionalMatchCollectionQueue) QueueMatchHistory(puuid string, after int64) chan error {
	collecter := NewMatchHistoryCollecter(cq.regionalServer, puuid, after, cq)
	return cq.matchHistoryCollectionQueue.Queue(collecter)
}
