package database

import (
	"TheCollectorDG/types"
	"fmt"
	"strings"
)

func GetStaleMatchHistory(excludePuuids []string) (*types.UpdateInfo, error) {
	var query string
	if len(excludePuuids) > 0 {
		query = fmt.Sprintf(`
			SELECT
				puuid,
				region,
				matches_last_updated
			FROM Summoner WHERE puuid NOT IN ('%s')
			ORDER BY matches_last_updated LIMIT 1
		`, strings.Join(excludePuuids, "', '"))
	} else {
		query = `
			SELECT
				puuid,
				region,
				matches_last_updated
			FROM Summoner
			ORDER BY matches_last_updated LIMIT 1
		`
	}

	updateInfo := new(types.UpdateInfo)
	row := db.QueryRow(query)
	err := row.Scan(
		&updateInfo.Puuid,
		&updateInfo.Region,
		&updateInfo.MatchesLastUpdated,
	)

	return updateInfo, err
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
