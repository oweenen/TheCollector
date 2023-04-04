package matchCollection

import (
	"TheCollectorDG/collection"
	"TheCollectorDG/collection/summonerCollection"
)

type MatchCollectionQueue struct {
	matchDetailsCollectionQueue collection.CollectionQueue
	matchHistoryCollectionQueue collection.CollectionQueue
	summonerCollectionQueue     *summonerCollection.SummonerCollectionQueue
}

func NewMatchCollectionQueue(summonerCollectionQueue *summonerCollection.SummonerCollectionQueue) *MatchCollectionQueue {
	return &MatchCollectionQueue{
		matchDetailsCollectionQueue: collection.NewCollectionQueue(),
		matchHistoryCollectionQueue: collection.NewCollectionQueue(),
		summonerCollectionQueue:     summonerCollectionQueue,
	}
}

func (cq *MatchCollectionQueue) QueueMatchDetails(matchId string) chan error {
	collecter := NewMatchDetailsCollecter(matchId, cq.summonerCollectionQueue)
	return cq.matchDetailsCollectionQueue.Queue(collecter)
}

func (cq *MatchCollectionQueue) QueueMatchHistory(region string, puuid string, after int64) chan error {
	collecter := NewMatchHistoryCollecter(region, puuid, after, cq, cq.summonerCollectionQueue)
	return cq.matchHistoryCollectionQueue.Queue(collecter)
}
