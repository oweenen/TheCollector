package matchCollection

import (
	"TheCollectorDG/database"
	"log"
	"os"
	"strconv"
	"time"
)

func MatchCollectionLoop(priorityQueue *RegionalMatchCollectionQueue, queue *RegionalMatchCollectionQueue, interval time.Duration) {
	doPassiveCollection, err := strconv.ParseBool(os.Getenv("DO_PASSIVE_COLLECTION"))
	if err != nil {
		doPassiveCollection = false
	}

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

		if doPassiveCollection && queue.matchHistoryCollectionQueue.NumActiveJobs() == 0 {
			err := queueStaleMatchHistory(queue)
			if err != nil {
				log.Println(err)
			}
		}
	}
}

func queueStaleMatchHistory(cq *RegionalMatchCollectionQueue) error {
	updateInfo, err := database.GetStaleMatchHistory(cq.regionalServer)
	if err == nil && updateInfo != nil {
		cq.QueueMatchHistory(updateInfo.Puuid, updateInfo.MatchesLastUpdated)
	}
	return err
}
