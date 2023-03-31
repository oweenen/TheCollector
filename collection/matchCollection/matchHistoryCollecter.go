package matchCollection

import (
	"TheCollectorDG/collection/summonerCollection"
	"TheCollectorDG/database"
	"TheCollectorDG/riot"
	"fmt"
	"sync"
	"time"
)

type MatchHistoryCollecter struct {
	MatchCQ    *MatchCollectionQueue
	SummonerCQ *summonerCollection.SummonerCollectionQueue
	Puuid      string
	After      int64
}

func NewMatchHistoryCollecter(puuid string, after int64, matchCQ *MatchCollectionQueue, summonerCQ *summonerCollection.SummonerCollectionQueue) MatchHistoryCollecter {
	return MatchHistoryCollecter{
		MatchCQ:    matchCQ,
		SummonerCQ: summonerCQ,
		Puuid:      puuid,
		After:      after,
	}
}

func (c MatchHistoryCollecter) Id() string {
	return c.Puuid
}

func (c MatchHistoryCollecter) Collect() error {
	fmt.Printf("Collecting match history for summoner %s\n", c.Puuid)

	updatedAt := time.Now().Unix()
	history, err := riot.GetMatchHistory(c.Puuid, c.After)
	if err != nil {
		return err
	}

	collectMatches(history, c.MatchCQ, c.SummonerCQ)

	err = database.SetMatchesUpdatedAt(c.Puuid, updatedAt)
	return err
}

func collectMatches(matchIds []string, matchCQ *MatchCollectionQueue, summonerCQ *summonerCollection.SummonerCollectionQueue) error {
	var wg sync.WaitGroup
	errChan := make(chan error)
	for _, matchId := range matchIds {
		wg.Add(1)
		go func(matchId string) {
			defer wg.Done()
			if !database.MatchIsStored(matchId) {
				if err := <-matchCQ.QueueMatchDetails(matchId); err != nil {
					errChan <- fmt.Errorf("error collecting match %s: %s", matchId, err)
				}
			}
		}(matchId)
	}
	wg.Wait()
	if len(errChan) > 0 {
		err := <-errChan
		close(errChan)
		return err
	}
	return nil
}
