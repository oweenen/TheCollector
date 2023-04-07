package matchCollection

import (
	"TheCollectorDG/database"
	"TheCollectorDG/riot"
	"log"
	"time"
)

func MatchCollectionLoop(priorityQueue *RegionalMatchCollectionQueue, queue *RegionalMatchCollectionQueue, interval time.Duration) {
	for range time.Tick(interval) {
		if priorityQueue.matchDetailsCollectionQueue.HasNext() {
			go priorityQueue.matchDetailsCollectionQueue.CollectNext()
		} else if priorityQueue.matchHistoryCollectionQueue.HasNext() {
			go priorityQueue.matchHistoryCollectionQueue.CollectNext()
		} else if queue.matchDetailsCollectionQueue.HasNext() {
			go queue.matchDetailsCollectionQueue.CollectNext()
		} else if queue.matchHistoryCollectionQueue.HasNext() {
			go queue.matchHistoryCollectionQueue.CollectNext()
		}

		if !queue.matchHistoryCollectionQueue.HasNext() {
			err := queueStaleMatchHistory(queue)
			if err != nil {
				log.Println(err)
			}
		}
	}
}

func queueStaleMatchHistory(cq *RegionalMatchCollectionQueue) error {
	updateInfo, err := database.GetStaleMatchHistory(riot.RiotRegionClusters[cq.regionalServer], cq.matchHistoryCollectionQueue.ListIds())
	if err == nil && updateInfo != nil {
		cq.QueueMatchHistory(updateInfo.Puuid, updateInfo.MatchesLastUpdated)
	}
	return err
}
