package database

import (
	"TheCollectorDG/types"
)

func StoreMatch(match *types.Match) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	_, err = tx.Exec(`
		INSERT IGNORE INTO TFT_Match (
			id,
			date,
			game_length,
			game_version,
			queue_id,
			game_type,
			set_name,
			set_number
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		`,
		match.Id,
		match.Date,
		match.GameLength,
		match.GameVersion,
		match.QueueId,
		match.GameType,
		match.SetName,
		match.SetNumber,
	)
	if err != nil {
		return err
	}

	// store comps
	for _, comp := range match.Comps {
		err := storeComp(tx, match.Id, &comp)
		if err != nil {
			return err
		}
	}

	err = tx.Commit()

	return err
}

func MatchIsStored(matchId string) bool {
	var stored bool
	db.QueryRow(`SELECT EXISTS (SELECT 1 FROM TFT_Match WHERE id = ?)`, matchId).Scan(&stored)
	return stored
}
