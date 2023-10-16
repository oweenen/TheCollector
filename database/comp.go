package database

import (
	"TheCollectorDG/types"
	"crypto/md5"
	"database/sql"
)

func storeComp(tx *sql.Tx, matchId string, comp *types.Comp) error {
	hasher := md5.New()
	hasher.Write([]byte(comp.SummonerPuuid + matchId))
	hashBytes := hasher.Sum(nil)

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
