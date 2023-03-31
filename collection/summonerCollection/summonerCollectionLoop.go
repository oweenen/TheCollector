package summonerCollection

import "time"

func SummonerCollectionLoop(priorityQueue *SummonerCollectionQueue, queue *SummonerCollectionQueue, interval time.Duration) {
	for range time.Tick(interval) {
		if priorityQueue.collectionQueue.HasNext() {
			go priorityQueue.collectionQueue.CollectNext()
		} else {
			go queue.collectionQueue.CollectNext()
		}
	}
}
