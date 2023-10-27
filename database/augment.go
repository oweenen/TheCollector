package database

import "database/sql"

func StoreAugment(tx *sql.Tx, compHash []byte, augment string, taken int, placement int) error {
	_, err := tx.Exec(`
		INSERT INTO Augment (
			comp_hash_bin,
			augment_id,
			taken,
			placement
		)
		VALUES (?, ?, ?, ?)
		`,
		compHash,
		augment,
		taken,
		placement,
	)
	return err
}
