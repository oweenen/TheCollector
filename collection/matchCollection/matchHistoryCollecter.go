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

	// fetch match history from riot
	matchIds, err := riot.GetMatchHistory(c.RegionalServer, c.Puuid, c.After)
	if err != nil {
		fmt.Printf("Failed to fetch match history from riot\nERROR: %v\nDETAILS:%+v\n", err.Error(), c)
		return err
	}

	// filter match ids not stored
	matchIdsNotStored := filterMatchIdsNotStored(matchIds)

	// queue new matches
	queueAndAwaitMatches(matchIdsNotStored, c.MatchCQ)

	// set matches_updated_at
	err = database.SetMatchesUpdatedAt(c.Puuid, updatedAt)
	if err != nil {
		fmt.Printf("Failed to set matches_updated_at\nERROR: %v\nDETAILS: {puuid: %v, updatedAt: %v}\n", err.Error(), c.Puuid, updatedAt)
		return err
	}

	// log collection
	fmt.Printf("Collected match history for summoner %s\n", c.Puuid)

	return nil
}

func filterMatchIdsNotStored(matchIds []string) []string {
	result := []string{}
	for _, matchId := range matchIds {
		if !database.MatchIsStored(matchId) {
			result = append(result, matchId)
		}
	}
	return result
}

func queueAndAwaitMatches(matchIds []string, matchCQ *RegionalMatchCollectionQueue) {
	var wg sync.WaitGroup
	for _, matchId := range matchIds {
		wg.Add(1)
		go func(matchId string) {
			defer wg.Done()
			<-matchCQ.QueueMatchDetails(matchId)
		}(matchId)
	}
	wg.Wait()
}
