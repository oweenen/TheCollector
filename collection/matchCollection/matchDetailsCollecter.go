package matchCollection

import (
	"TheCollectorDG/collection/summonerCollection"
	"TheCollectorDG/database"
	"TheCollectorDG/datastore"
	"TheCollectorDG/riot"
	"TheCollectorDG/types"
	"fmt"
	"sync"
	"time"
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
	// fetch match data from riot
	matchRes, err := riot.GetMatchDetails(c.RegionalServer, c.MatchId)
	if err != nil {
		fmt.Printf("Failed to get match details from riot\n\tERROR: %v\n\tCONTEXT: %+v\n", err, c)
		return err
	}

	// create match from matchRes
	match := types.NewMatchFromRiotRes(matchRes)

	// skip storing match if match is not Normal, Ranked, DoubleUp
	if match.QueueId != 1090 && match.QueueId != 1100 && match.QueueId != 1160 {
		return nil
	}

	// extract summoners not stored from match
	summonersNotStored := extractSummonersNotStored(match)

	// queue and await summoners not stored
	err = queueAndAwaitSummonersNotStored(summonersNotStored, c.SummonerCollectionQueue)
	if err != nil {
		fmt.Printf("Aborting match details collection, summoner collection failed\n\tERROR: %v\n\tCONTEXT: %+v\n", err, c)
		return err
	}

	// store match details to s3
	err = datastore.StoreMatch(match)
	if err != nil {
		fmt.Printf("Failed to store match details to s3\nERROR: %v\n\tCONTEXT: %+v\n\tMATCH: %+v\n", err, c, match)
		return err
	}

	// queue stale ranks
	QueueStaleRankUpdates(match, c.SummonerCollectionQueue)

	// create database transaction
	tx, err := database.NewTransaction()
	if err != nil {
		tx.Rollback()
		fmt.Printf("Failed to create database transaction\nERROR: %v\n", err)
		return err
	}

	// store match to database
	err = database.StoreMatch(tx, match)
	if err != nil {
		tx.Rollback()
		fmt.Printf("Failed to store match to database\nERROR: %v\n\tCONTEXT: %+v\n\tMATCH: %+v\n", err, c, match)
		return err
	}

	// store comps to database
	for _, comp := range match.Comps {
		err = database.StoreComp(tx, match.Id, &comp)
		if err != nil {
			tx.Rollback()
			fmt.Printf("Failed to store comp to database\nERROR: %v\n\tCONTEXT: %+v\n\tCOMP: %+v\n", err, c, comp)
			return err
		}
	}

	// store augment if match is ranked and match is from within the past week
	if match.QueueId == 1100 && match.Date > time.Now().UnixMilli()-1000*60*60*24*7 {
		for _, comp := range match.Comps {
			for i, augment := range comp.Augments {
				err := database.StoreAugment(tx, match.Id, comp.SummonerPuuid, match.GameVersion, augment, i, comp.Placement)
				if err != nil {
					tx.Rollback()
					fmt.Printf("Failed to store augment to database\nERROR: %v\n\tCONTEXT: %+v\n\tAUGMENT: %+v\n", err, c, augment)
					return err
				}
			}
		}
	}

	// commit database transaction
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		fmt.Printf("Failed to commit transaction to database\nERROR: %v\n\tCONTEXT: %+v\n", err, c)
		return err
	}

	fmt.Printf("Collected match %s\n", c.MatchId)
	return nil
}

func extractSummonersNotStored(match *types.Match) []string {
	var result = []string{}
	for _, comp := range match.Comps {
		if !database.SummonerIsStored(comp.SummonerPuuid) {
			result = append(result, comp.SummonerPuuid)
		}
	}
	return result
}

func queueAndAwaitSummonersNotStored(summoners []string, summonerCollectionQueue *summonerCollection.RegionalSummonerCollectionQueue) error {
	var wg sync.WaitGroup
	errChan := make(chan error, len(summoners))

	for _, puuid := range summoners {
		wg.Add(1)
		go func(puuid string) {
			defer wg.Done()
			if !database.SummonerIsStored(puuid) {
				err := <-summonerCollectionQueue.QueueSummonerByPuuid(puuid)
				if err != nil {
					errChan <- err
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
