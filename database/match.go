package database

import (
	"TheCollectorDG/types"
)

func GetRecentMatches(puuid string, count int) ([]*types.Match, error) {
	rows, err := db.Query(`
	SELECT
    	TFT_Match.id,
    	TFT_Match.date,
    	TFT_Match.game_length,
    	TFT_Match.game_version,
    	TFT_Match.queue_id,
    	TFT_Match.game_type,
    	TFT_Match.set_name,
    	TFT_Match.set_number
	FROM TFT_Match
	JOIN Comp ON TFT_Match.id = Comp.match_id
	WHERE Comp.summoner_puuid = ?
	ORDER BY TFT_Match.date DESC
	LIMIT ?
		`,
		puuid,
		count,
	)
	if err != nil {
		return nil, err
	}

	var matches []*types.Match
	defer rows.Close()
	for rows.Next() {
		match := new(types.Match)

		err = rows.Scan(
			&match.Id,
			&match.Date,
			&match.GameLength,
			&match.GameVersion,
			&match.QueueId,
			&match.GameType,
			&match.SetName,
			&match.SetNumber,
		)
		if err != nil {
			continue
		}

		matches = append(matches, match)
	}

	return matches, nil
}

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

	// store augments if queue is ranked
	if match.QueueId == 1100 {
		for _, comp := range match.Comps {
			compHashBin := compHashBin(match.Id, comp.SummonerPuuid)
			for i, augment := range comp.Augments {
				var taken int
				if len(comp.Augments) == 3 {
					// all augments present
					taken = i
				} else if len(comp.Augments) == 2 && comp.LastRound >= 20 {
					// 1 augment missing, but all augments chosen (2-1 is legend augment)
					taken = i + 1
				} else {
					// rare edge cases taken is unknown
					taken = -1
				}

				err := StoreAugment(tx, compHashBin, augment, taken)
				if err != nil {
					return err
				}
			}
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
