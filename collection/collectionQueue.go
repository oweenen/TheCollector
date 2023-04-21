package collection

import (
	"sync"
)

type CollectionQueue struct {
	mu                    sync.Mutex
	collectionBucketQueue []*CollectionBucket
	collectionBuckets     map[string]*CollectionBucket
}

func NewCollectionQueue() CollectionQueue {
	return CollectionQueue{
		collectionBuckets: make(map[string]*CollectionBucket),
	}
}

func (cq *CollectionQueue) Queue(collecter Collecter) chan error {
	cq.mu.Lock()
	defer cq.mu.Unlock()

	collectionBucket, exists := cq.collectionBuckets[collecter.Id()]
	if !exists {
		collectionBucket = &CollectionBucket{
			Collecter: collecter,
		}
		cq.collectionBucketQueue = append(cq.collectionBucketQueue, collectionBucket)
		cq.collectionBuckets[collectionBucket.Collecter.Id()] = collectionBucket
	}

	collectionResChannel := make(chan error, 1)
	collectionBucket.Listeners = append(collectionBucket.Listeners, collectionResChannel)

	return collectionResChannel
}

func (cq *CollectionQueue) CollectNext() {
	collectionBucket := cq.NextInQueue()
	if collectionBucket == nil {
		return
	}

	collectionRes := collectionBucket.Collecter.Collect()

	cq.mu.Lock()
	defer cq.mu.Unlock()
	collectionBucket.NotifyListeners(collectionRes)
	delete(cq.collectionBuckets, collectionBucket.Collecter.Id())
}

func (cq *CollectionQueue) NextInQueue() *CollectionBucket {
	cq.mu.Lock()
	defer cq.mu.Unlock()

	if len(cq.collectionBucketQueue) == 0 {
		return nil
	}

	collectionBucket := cq.collectionBucketQueue[0]
	cq.collectionBucketQueue = cq.collectionBucketQueue[1:]

	return collectionBucket
}

func (cq *CollectionQueue) HasNext() bool {
	cq.mu.Lock()
	defer cq.mu.Unlock()
	return len(cq.collectionBucketQueue) > 0
}

func (cq *CollectionQueue) NumActiveJobs() int {
	cq.mu.Lock()
	defer cq.mu.Unlock()
	return len(cq.collectionBuckets)
}

func (cq *CollectionQueue) ListIds() []string {
	cq.mu.Lock()
	defer cq.mu.Unlock()

	ids := make([]string, len(cq.collectionBuckets))
	index := 0
	for id := range cq.collectionBuckets {
		ids[index] = id
		index++
	}

	return ids
}
