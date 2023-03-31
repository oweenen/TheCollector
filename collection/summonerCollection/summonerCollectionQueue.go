package summonerCollection

import (
	"TheCollectorDG/collection"
)

type SummonerCollectionQueue struct {
	collectionQueue collection.CollectionQueue
}

func NewSummonerCollectionQueue() *SummonerCollectionQueue {
	return &SummonerCollectionQueue{
		collectionQueue: collection.NewCollectionQueue(),
	}
}

func (cq *SummonerCollectionQueue) QueueSummonerByName(region string, name string) chan error {
	collecter := NewSummonerByNameCollecter(region, name)
	return cq.collectionQueue.Queue(collecter)
}

func (cq *SummonerCollectionQueue) QueueSummonerByPuuid(region string, puuid string) chan error {
	collecter := NewSummonerByPuuidCollecter(region, puuid)
	return cq.collectionQueue.Queue(collecter)
}
