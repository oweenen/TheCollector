package database

import (
	"TheCollectorDG/types"
	"fmt"
	"time"
)

func GetUpdateInfo(puuid string) (*types.UpdateInfo, error) {
	updateInfo := new(types.UpdateInfo)
	row := db.QueryRow(`
		SELECT
			puuid,
			region,
			summoner_id,
			matches_last_updated,
			last_updated,
			rank_last_updated
		FROM Summoner
		WHERE puuid = ?
		`,
		puuid)
	err := row.Scan(
		&updateInfo.Puuid,
		&updateInfo.Region,
		&updateInfo.SummonerId,
		&updateInfo.MatchesLastUpdated,
		&updateInfo.LastUpdated,
		&updateInfo.RankLastUpdated,
	)

	return updateInfo, err
}

func GetStaleMatchHistory(regionCluster string) (*types.UpdateInfo, error) {
	query := fmt.Sprintf(`
		SELECT
			puuid,
			region,
			matches_last_updated
		FROM Summoner WHERE region_cluster = '%s'
		ORDER BY matches_last_updated LIMIT 1
	`, regionCluster)

	updateInfo := new(types.UpdateInfo)
	row := db.QueryRow(query)
	err := row.Scan(
		&updateInfo.Puuid,
		&updateInfo.Region,
		&updateInfo.MatchesLastUpdated,
	)
	return updateInfo, err
}

func GetStaleRankFromMatch(matchId string) ([]*types.UpdateInfo, error) {
	// Threshold is 1 week ago
	staleThreshold := time.Now().Unix() - 60*60*24*7

	rows, err := db.Query(`
	SELECT
    	s.puuid,
    	s.region,
    	s.summoner_id,
    	s.rank_last_updated
	FROM
    	Summoner AS s
	INNER JOIN
	    Comp AS c
	ON
	    c.summoner_puuid = s.puuid
	WHERE
	    c.match_id = ?
	    AND s.rank_last_updated < ?
	`, matchId, staleThreshold)
	if err != nil {
		return nil, err
	}
	var staleRanks []*types.UpdateInfo
	defer rows.Close()
	for rows.Next() {
		updateInfo := new(types.UpdateInfo)

		err = rows.Scan(
			&updateInfo.Puuid,
			&updateInfo.Region,
			&updateInfo.SummonerId,
			&updateInfo.RankLastUpdated,
		)

		if err != nil {
			continue
		}

		staleRanks = append(staleRanks, updateInfo)
	}

	return staleRanks, err
}

func SetLastUpdated(puuid string, updatedAt int64) error {
	_, err := db.Exec(`
		UPDATE Summoner
			SET last_updated = ?
		WHERE puuid = ?
		`,
		updatedAt,
		puuid,
	)
	return err
}

func SetMatchesUpdatedAt(puuid string, updatedAt int64) error {
	_, err := db.Exec(`
		UPDATE Summoner
			SET matches_last_updated = ?
		WHERE puuid = ?
		`,
		updatedAt,
		puuid,
	)
	return err
}

func SetRankUpdatedAt(puuid string, updatedAt int64) error {
	_, err := db.Exec(`
		UPDATE Summoner
			SET rank_last_updated = ?
		WHERE puuid = ?
		`,
		updatedAt,
		puuid,
	)
	return err
}
