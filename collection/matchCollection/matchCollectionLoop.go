package matchCollection

import (
	"TheCollectorDG/database"
	"log"
	"time"
)

func MatchCollectionLoop(priorityQueue *MatchCollectionQueue, queue *MatchCollectionQueue, interval time.Duration) {
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

func queueStaleMatchHistory(cq *MatchCollectionQueue) error {
	updateInfo, err := database.GetStaleMatchHistory(cq.matchHistoryCollectionQueue.ListIds())
	if err == nil && updateInfo != nil {
		cq.QueueMatchHistory(updateInfo.Puuid, updateInfo.MatchesLastUpdated)
	}
	return err
}
