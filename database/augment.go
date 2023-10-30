package database

import "database/sql"

func StoreAugment(tx *sql.Tx, compHash []byte, gameVersion string, augment string, taken int, placement int) error {
	_, err := tx.Exec(`
		INSERT INTO Augment (
			comp_hash_bin,
			game_version,
			augment_id,
			taken,
			placement
		)
		VALUES (?, ?, ?, ?, ?)
		`,
		compHash,
		gameVersion,
		augment,
		taken,
		placement,
	)
	return err
}
