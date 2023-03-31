package collection

type CollectionBucket struct {
	Collecter Collecter
	Listeners []chan error
}

func (collectionBucket *CollectionBucket) NotifyListeners(res error) {
	for _, listener := range collectionBucket.Listeners {
		listener <- res
	}
}
