package database

import (
	"TheCollectorDG/types"
	"crypto/md5"
	"database/sql"
)

func compHashBin(matchId string, summonerPuuid string) []byte {
	hasher := md5.New()
	hasher.Write([]byte(summonerPuuid + matchId))
	hashBytes := hasher.Sum(nil)

	return hashBytes
}

func storeComp(tx *sql.Tx, matchId string, comp *types.Comp) error {
	hashBytes := compHashBin(matchId, comp.SummonerPuuid)

	_, err := tx.Exec(`
		INSERT IGNORE INTO Comp (
			hash_bin,
			match_id,
			summoner_puuid,
			placement
		)
		VALUES (?, ?, ?, ?)
		`,
		hashBytes,
		matchId,
		comp.SummonerPuuid,
		comp.Placement,
	)
	return err
}
