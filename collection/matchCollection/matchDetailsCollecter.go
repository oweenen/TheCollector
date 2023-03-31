package matchCollection

import (
	"TheCollectorDG/collection/summonerCollection"
	"TheCollectorDG/database"
	"TheCollectorDG/riot"
	"TheCollectorDG/types"
	"fmt"
	"sync"
)

type MatchDetailsCollecter struct {
	SummonerCollectionQueue *summonerCollection.SummonerCollectionQueue
	MatchId                 string
}

func NewMatchDetailsCollecter(matchId string, summonerCollectionQueue *summonerCollection.SummonerCollectionQueue) MatchDetailsCollecter {
	return MatchDetailsCollecter{
		SummonerCollectionQueue: summonerCollectionQueue,
		MatchId:                 matchId,
	}
}

func (c MatchDetailsCollecter) Id() string {
	return c.MatchId
}

func (c MatchDetailsCollecter) Collect() error {
	fmt.Printf("Collecting match %s\n", c.MatchId)

	match, err := riot.GetMatchDetails(c.MatchId)
	if err != nil {
		fmt.Printf("ERROR failed to get match %s from riot: %s\n", c.MatchId, err)
		return err
	}

	err = QueueSummonersNotStored(match, c.SummonerCollectionQueue)
	if err != nil {
		return err
	}

	err = database.StoreMatch(match)

	return err
}

func QueueSummonersNotStored(match *types.Match, summonerCollectionQueue *summonerCollection.SummonerCollectionQueue) error {
	var wg sync.WaitGroup
	errChan := make(chan error)
	region := types.GetMatchIdRegion(match.Id)

	for _, comp := range match.Comps {
		puuid := comp.Summoner.Puuid
		wg.Add(1)
		go func(puuid string) {
			defer wg.Done()
			if !database.SummonerIsStored(puuid) {
				err := <-summonerCollectionQueue.QueueSummonerByPuuid(region, puuid)
				if err != nil {
					errChan <- fmt.Errorf("failed to collect summoner %s %s", region, puuid)
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
