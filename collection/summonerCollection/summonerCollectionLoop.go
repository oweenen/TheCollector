package summonerCollection

import "time"

func SummonerCollectionLoop(priorityQueue *RegionalSummonerCollectionQueue, queue *RegionalSummonerCollectionQueue, interval time.Duration) {
	for range time.Tick(interval) {
		if priorityQueue.summonerCollectionQueue.HasNext() {
			go priorityQueue.summonerCollectionQueue.CollectNext()
		} else if queue.summonerCollectionQueue.HasNext() {
			go queue.summonerCollectionQueue.CollectNext()
		}
	}
}
