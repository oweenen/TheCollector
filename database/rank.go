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
		FROM TFT_Rank WHERE summoner_puuid = ?
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

type RankStats struct {
	TotalGames       int     `json:"total_games"`
	AveragePlacement float32 `json:"average_placement"`
	Top4Rate         float32 `json:"top_4_rate"`
}

func GetRankStats(puuid string) (*RankStats, error) {
	rankStats := new(RankStats)
	row := db.QueryRow(`
		SELECT
			played,
			avg_placement,
			top_4_rate
		FROM TFT_Rank WHERE summoner_puuid = ?
	`, puuid)
	err := row.Scan(
		&rankStats.TotalGames,
		&rankStats.AveragePlacement,
		&rankStats.Top4Rate,
	)
	return rankStats, err
}

func UpdateRankStats(puuid string) error {
	_, err := db.Exec(`
	UPDATE
	TFT_Rank AS r
	JOIN (
		SELECT
			summoner_puuid,
			COUNT(*) AS played,
			AVG(placement) AS avg_placement,
			(SUM(CASE WHEN placement <= 4 THEN 1 ELSE 0 END) / COUNT(*)) * 100 AS top_4_rate
		FROM
			Comp
			JOIN TFT_Match ON Comp.match_id = TFT_Match.id
		WHERE
			Comp.summoner_puuid = ?
			AND TFT_Match.queue_id = 1100) AS rank_stats 
	ON r.summoner_puuid = rank_stats.summoner_puuid 
	SET 
		r.played = rank_stats.played,
		r.avg_placement = rank_stats.avg_placement,
		r.top_4_rate = rank_stats.top_4_rate
	`,
		puuid,
	)

	return err
}
