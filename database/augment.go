package database

import (
	"database/sql"
	"fmt"
)

func StoreAugment(tx *sql.Tx, matchId, summonerPuuid, gameVersion, augment string, taken, placement int) error {
	compHashBin := compHashBin(matchId, summonerPuuid)
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
		compHashBin,
		gameVersion,
		augment,
		taken,
		placement,
	)
	return err
}

type AugmentStats struct {
	AvgPlacement float32 `json:"avg_placement"`
	TimesPlayed  int     `json:"times_played"`
	Frequency    float32 `json:"frequency"`
}

func GetAugmentStats() (map[string]map[string][]*AugmentStats, error) {
	augmentStatsMap := make(map[string]map[string][]*AugmentStats)

	rows, err := db.Query(`
	SELECT
		game_version,
		augment_id,
		-1 as taken,
		AVG(placement) AS avg_placement,
		COUNT(*) AS times_played,
		(COUNT(*) / (
			SELECT
				COUNT(*)
			FROM
				Augment) * 100) AS frequency
	FROM
		Augment
	GROUP BY
		game_version, augment_id
		
	UNION

	SELECT
		game_version,
		augment_id,
		taken,
		AVG(placement) AS avg_placement,
		COUNT(*) AS times_played,
		(COUNT(*) / (
				SELECT
					COUNT(*)
				FROM
					Augment) * 100) AS frequency
	FROM
		Augment
	GROUP BY
		game_version, augment_id, taken
	`)
	if err != nil {
		return augmentStatsMap, err
	}

	defer rows.Close()
	for rows.Next() {
		augmentStats := new(AugmentStats)

		var gameVersion string
		var augmentId string
		var taken int

		err = rows.Scan(
			&gameVersion,
			&augmentId,
			&taken,
			&augmentStats.AvgPlacement,
			&augmentStats.TimesPlayed,
			&augmentStats.Frequency,
		)
		if err != nil {
			fmt.Println(err.Error())
			return augmentStatsMap, err
		}

		patchAugmentsStatsMap, ok := augmentStatsMap[gameVersion]
		if !ok {
			patchAugmentsStatsMap = make(map[string][]*AugmentStats)
			augmentStatsMap[gameVersion] = patchAugmentsStatsMap
		}

		augmentStatsArr, ok := patchAugmentsStatsMap[augmentId]
		if !ok {
			augmentStatsArr = make([]*AugmentStats, 4)
			patchAugmentsStatsMap[augmentId] = augmentStatsArr
		}

		augmentStatsArr[taken+1] = augmentStats
	}

	return augmentStatsMap, nil
}
