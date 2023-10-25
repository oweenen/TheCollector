package database

import (
	"TheCollectorDG/types"
)

func StoreRank(puuid string, rank *types.Rank) error {
	_, err := db.Exec(`
		INSERT INTO TFT_Rank (
			summoner_puuid,
			tier,
			division,
			lp
		)
		VALUES (?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE
			tier = VALUES(tier),
			division = VALUES(division),
			lp = VALUES(lp)
		`,
		puuid,
		rank.Tier,
		rank.Division,
		rank.Lp,
	)
	return err
}

func GetRank(puuid string) (*types.Rank, error) {
	rank := new(types.Rank)
	row := db.QueryRow(`
		SELECT
			tier,
			division,
			lp
		FROM TFT_Rank WHERE summoner_puuid = ? LIMIT 1
		`,
		puuid,
	)
	err := row.Scan(
		&rank.Tier,
		&rank.Division,
		&rank.Lp,
	)
	return rank, err
}
