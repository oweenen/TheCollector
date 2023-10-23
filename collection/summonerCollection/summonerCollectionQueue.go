package summonerCollection

import (
	"TheCollectorDG/collection"
)

type RegionalSummonerCollectionQueue struct {
	region                  string
	summonerCollectionQueue collection.CollectionQueue
}

func NewRegionalSummonerCollectionQueue(region string) *RegionalSummonerCollectionQueue {
	return &RegionalSummonerCollectionQueue{
		region:                  region,
		summonerCollectionQueue: collection.NewCollectionQueue(),
	}
}

func (cq *RegionalSummonerCollectionQueue) QueueSummonerByName(name string) chan error {
	collecter := NewSummonerByNameCollecter(cq.region, name)
	return cq.summonerCollectionQueue.Queue(collecter)
}

func (cq *RegionalSummonerCollectionQueue) QueueSummonerByPuuid(puuid string) chan error {
	collecter := NewSummonerByPuuidCollecter(cq.region, puuid)
	return cq.summonerCollectionQueue.Queue(collecter)
}

func (cq *RegionalSummonerCollectionQueue) QueueRank(puuid string, summonerId string) chan error {
	collecter := NewRankCollecter(cq.region, puuid, summonerId)
	return cq.summonerCollectionQueue.Queue(collecter)
}
