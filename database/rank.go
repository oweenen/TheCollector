package database

import (
	"TheCollectorDG/types"
)

func StoreRank(puuid string, rank *types.Rank) error {
	_, err := db.Exec(`
		INSERT INTO TFT_Rank (
			summoner_puuid,
			type,
			tier,
			division,
			lp,
			raw_lp
		)
		VALUES (?, ?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE
			tier = VALUES(tier),
			division = VALUES(division),
			lp = VALUES(lp),
			raw_lp = VALUES(raw_lp)
		`,
		puuid,
		rank.Type,
		rank.Tier,
		rank.Division,
		rank.Lp,
		rank.RawLp,
	)
	return err
}

func GetRank(puuid string) (*types.Rank, error) {
	rank := new(types.Rank)
	row := db.QueryRow(`
		SELECT
			type,
			tier,
			division,
			lp,
			raw_lp
		FROM TFT_Rank WHERE summoner_puuid = ? LIMIT 1
		`,
		puuid,
	)
	err := row.Scan(
		&rank.Type,
		&rank.Tier,
		&rank.Division,
		&rank.Lp,
		&rank.RawLp,
	)
	return rank, err
}
