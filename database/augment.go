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

type AugmentStatsPage struct {
	GameVersion string         `json:"game_version"`
	Augments    []AugmentStats `json:"augments"`
}

type AugmentStats struct {
	AugmentId  string             `json:"augment_id"`
	Overall    *AugmentStageStats `json:"overall"`
	FirstPick  *AugmentStageStats `json:"first_pick"`
	SecondPick *AugmentStageStats `json:"second_pick"`
	ThirdPick  *AugmentStageStats `json:"third_pick"`
}

type AugmentStageStats struct {
	AvgPlacement float32 `json:"avg_placement"`
	TimesPlayed  int     `json:"times_played"`
	Frequency    float32 `json:"frequency"`
}

func GetAugmentStats() (map[string]AugmentStatsPage, error) {
	augmentMapsPages := make(map[string]map[string][]*AugmentStageStats)

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
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		augmentStageStats := new(AugmentStageStats)

		var gameVersion string
		var augmentId string
		var taken int

		err = rows.Scan(
			&gameVersion,
			&augmentId,
			&taken,
			&augmentStageStats.AvgPlacement,
			&augmentStageStats.TimesPlayed,
			&augmentStageStats.Frequency,
		)
		if err != nil {
			fmt.Println(err.Error())
			return nil, err
		}

		augmentMap, ok := augmentMapsPages[gameVersion]
		if !ok {
			augmentMap = make(map[string][]*AugmentStageStats)
			augmentMapsPages[gameVersion] = augmentMap
		}

		augmentStageStatsArr, ok := augmentMap[augmentId]
		if !ok {
			augmentStageStatsArr = make([]*AugmentStageStats, 4)
			augmentMap[augmentId] = augmentStageStatsArr
		}

		augmentStageStatsArr[taken+1] = augmentStageStats
	}

	augmentStatsPages := make(map[string]AugmentStatsPage)
	for gameVersion, augmentStatsMap := range augmentMapsPages {
		augmentStatsPage := AugmentStatsPage{
			GameVersion: gameVersion,
			Augments:    []AugmentStats{},
		}

		for augmentId, stageStatsArr := range augmentStatsMap {
			augmentStatsPage.Augments = append(augmentStatsPage.Augments, AugmentStats{
				AugmentId:  augmentId,
				Overall:    stageStatsArr[0],
				FirstPick:  stageStatsArr[1],
				SecondPick: stageStatsArr[2],
				ThirdPick:  stageStatsArr[3],
			})
		}

		augmentStatsPages[gameVersion] = augmentStatsPage
	}

	return augmentStatsPages, nil
}
