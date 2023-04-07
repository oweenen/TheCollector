package matchCollection

import (
	"TheCollectorDG/database"
	"TheCollectorDG/riot"
	"fmt"
	"sync"
	"time"
)

type MatchHistoryCollecter struct {
	MatchCQ        *RegionalMatchCollectionQueue
	RegionalServer string
	Puuid          string
	After          int64
}

func NewMatchHistoryCollecter(regionalServer string, puuid string, after int64, matchCQ *RegionalMatchCollectionQueue) MatchHistoryCollecter {
	return MatchHistoryCollecter{
		MatchCQ:        matchCQ,
		RegionalServer: regionalServer,
		Puuid:          puuid,
		After:          after,
	}
}

func (c MatchHistoryCollecter) Id() string {
	return c.Puuid
}

func (c MatchHistoryCollecter) Collect() error {
	updatedAt := time.Now().Unix()
	history, err := riot.GetMatchHistory(c.RegionalServer, c.Puuid, c.After)
	if err != nil {
		return err
	}

	collectMatches(history, c.MatchCQ)

	err = database.SetMatchesUpdatedAt(c.Puuid, updatedAt)
	fmt.Printf("Collected match history for summoner %s\n", c.Puuid)
	return err
}

func collectMatches(matchIds []string, matchCQ *RegionalMatchCollectionQueue) error {
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
