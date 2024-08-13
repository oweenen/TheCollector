package main

import (
	"TheCollectorDG/db"
	"TheCollectorDG/riot"
	"TheCollectorDG/types"
	"context"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/joho/godotenv"
)

var conn *pgx.Conn
var queries *db.Queries

func init() {
	var err error

	ctx := context.Background()

	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	conn, err = pgx.Connect(ctx, os.Getenv("DB_URL"))
	if err != nil {
		panic(err)
	}

	queries = db.New(conn)
}

func main() {
	ctx := context.Background()

	collectionLoop(ctx)
}

func collectionLoop(ctx context.Context) {
	for {
		oldestMatchHistories, _ := queries.GetOldestMatchHistories(ctx, 1)
		puuid := oldestMatchHistories[0].Puuid

		updatedAt := time.Now()
		matchIds, _ := riot.GetMatchHistory("americas", puuid, oldestMatchHistories[0].MatchesUpdated.Time.UnixMilli())
		matchIdsLen := len(matchIds)

		log.Printf("Got %v matchIds from summoner %v\n", matchIdsLen, puuid)

		for i, matchId := range matchIds {
			if exists, _ := queries.MatchExists(ctx, matchId); exists {
				log.Printf("%v/%v: Skipping matchId %v...\n", i+1, matchIdsLen, matchId)
				continue
			}

			matchDetails, _ := riot.GetMatchDetails("americas", matchId)
			storeMatchDetails(ctx, matchDetails)
			log.Printf("%v/%v: Inserted new match with matchId %v!\n", i+1, matchIdsLen, matchId)
		}

		queries.SetMatchesUpdated(ctx, db.SetMatchesUpdatedParams{
			Puuid: puuid,
			MatchesUpdated: pgtype.Timestamp{
				Time: updatedAt,
			},
		})
	}
}

func storeMatchDetails(ctx context.Context, matchDetails *riot.Match) {
	// insert participants
	for _, puuid := range matchDetails.MetaData.Participants {
		queries.InsertPuuid(ctx, puuid)
	}

	tx, _ := conn.Begin(ctx)
	defer tx.Rollback(ctx)
	qtx := queries.WithTx(tx)

	// create match
	qtx.CreateMatch(ctx, db.CreateMatchParams{
		ID:          matchDetails.MetaData.MatchId,
		DataVersion: matchDetails.MetaData.DataVersion,
		GameVersion: matchDetails.Info.GameVersion,
		QueueID:     matchDetails.Info.QueueId,
		GameType:    matchDetails.Info.GameType,
		SetName:     matchDetails.Info.SetName,
		SetNumber:   matchDetails.Info.SetNumber,
	})

	// create comps
	for _, compDetails := range matchDetails.Info.Comps {
		qtx.CreateComp(ctx, db.CreateCompParams{
			MatchID:       matchDetails.MetaData.MatchId,
			SummonerPuuid: compDetails.Puuid,
			CompData:      types.CompData(compDetails),
		})
	}

	tx.Commit(ctx)
}
