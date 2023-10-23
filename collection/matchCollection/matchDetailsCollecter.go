package matchCollection

import (
	"TheCollectorDG/collection/summonerCollection"
	"TheCollectorDG/database"
	"TheCollectorDG/datastore"
	"TheCollectorDG/riot"
	"TheCollectorDG/types"
	"fmt"
	"sync"
)

type MatchDetailsCollecter struct {
	SummonerCollectionQueue *summonerCollection.RegionalSummonerCollectionQueue
	RegionalServer          string
	MatchId                 string
}

func NewMatchDetailsCollecter(regionalServer string, matchId string, summonerCollectionQueue *summonerCollection.RegionalSummonerCollectionQueue) MatchDetailsCollecter {
	return MatchDetailsCollecter{
		SummonerCollectionQueue: summonerCollectionQueue,
		RegionalServer:          regionalServer,
		MatchId:                 matchId,
	}
}

func (c MatchDetailsCollecter) Id() string {
	return c.MatchId
}

func (c MatchDetailsCollecter) Collect() error {
	match, err := riot.GetMatchDetails(c.RegionalServer, c.MatchId)
	if err != nil {
		fmt.Printf("ERROR failed to get match %s from riot: %s\n", c.MatchId, err)
		return err
	}

	if match.QueueId == 1110 || match.QueueId == 1111 {
		return nil
	}

	err = QueueSummonersNotStored(match, c.SummonerCollectionQueue)
	if err != nil {
		return err
	}

	datastore.StoreMatch(match)

	err = database.StoreMatch(match)
	if err != nil {
		return err
	}

	QueueStaleRankUpdates(match, c.SummonerCollectionQueue)

	fmt.Printf("Collected match %s\n", c.MatchId)
	return nil
}

func QueueSummonersNotStored(match *types.Match, summonerCollectionQueue *summonerCollection.RegionalSummonerCollectionQueue) error {
	var wg sync.WaitGroup
	errChan := make(chan error, len(match.Comps))

	for _, comp := range match.Comps {
		puuid := comp.SummonerPuuid
		wg.Add(1)
		go func(puuid string) {
			defer wg.Done()
			if !database.SummonerIsStored(puuid) {
				err := <-summonerCollectionQueue.QueueSummonerByPuuid(puuid)
				if err != nil {
					errChan <- fmt.Errorf("failed to collect summoner %s", puuid)
				}
			}
		}(puuid)
	}

	wg.Wait()
	if len(errChan) > 0 {
		err := <-errChan
		close(errChan)
		return err
	}
	return nil
}

func QueueStaleRankUpdates(match *types.Match, summonerCollectionQueue *summonerCollection.RegionalSummonerCollectionQueue) error {
	// if match is not ranked skip
	if match.QueueId != 1100 {
		return nil
	}

	staleRanks, err := database.GetStaleRankFromMatch(match.Id)
	if err != nil {
		return err
	}

	for _, updateInfo := range staleRanks {
		summonerCollectionQueue.QueueRank(updateInfo.Puuid, updateInfo.SummonerId)
	}

	return nil
}
