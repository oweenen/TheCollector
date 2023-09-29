package database

import (
	"TheCollectorDG/riot"
	"TheCollectorDG/types"
	"strings"
)

func StoreSummoner(summoner *types.Summoner) error {
	_, err := db.Exec(`
		INSERT INTO Summoner (
			puuid,
			region,
			region_cluster,
			summoner_id,
			raw_name,
			display_name,
			profile_icon_id,
			summoner_level
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE
			region = VALUES(region),
			region_cluster = VALUES(region_cluster),
			summoner_id = VALUES(summoner_id),
			raw_name = VALUES(raw_name),
			display_name = VALUES(display_name),
			profile_icon_id = VALUES(profile_icon_id),
			summoner_level = VALUES(summoner_level)
		`,
		summoner.Puuid,
		summoner.Region,
		riot.RiotRegionRoutes[summoner.Region],
		summoner.SummonerId,
		types.ToRawName(summoner.Name),
		summoner.Name,
		summoner.ProfileIconId,
		summoner.SummonerLevel,
	)
	return err
}

func GetSummoner(region string, name string) (*types.Summoner, error) {
	summoner := new(types.Summoner)
	row := db.QueryRow(`
		SELECT
			puuid,
			region,
			summoner_id,
			display_name,
			profile_icon_id,
			summoner_level,
			last_updated
		FROM Summoner WHERE raw_name = ? AND region = ? LIMIT 1
		`,
		types.ToRawName(name),
		strings.ToLower(region),
	)
	err := row.Scan(
		&summoner.Puuid,
		&summoner.Region,
		&summoner.SummonerId,
		&summoner.Name,
		&summoner.ProfileIconId,
		&summoner.SummonerLevel,
		&summoner.LastUpdated,
	)
	return summoner, err
}

func SummonerIsStored(puuid string) bool {
	var stored bool
	db.QueryRow(`SELECT EXISTS (SELECT 1 FROM Summoner WHERE puuid = ?)`, puuid).Scan(&stored)
	return stored
}
