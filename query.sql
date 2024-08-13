-- name: InsertPuuid :exec
INSERT INTO tft_summoner (
    puuid
) VALUES (
    $1
) ON CONFLICT (puuid) DO NOTHING;

-- name: CreateMatch :exec
INSERT INTO tft_match (
    id,
	data_version,
	game_version,
	queue_id,
	game_type,
	set_name,
	set_number
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
);

-- name: _createComp :exec
INSERT INTO tft_comp (
    match_id,
	summoner_puuid,
	comp_data
) VALUES (
    $1, $2, $3
);

-- name: GetOldestMatchHistories :many
SELECT
    puuid,
    matches_updated
FROM tft_summoner
ORDER BY matches_updated ASC NULLS FIRST
LIMIT $1;

-- name: MatchExists :one
SELECT EXISTS (
    SELECT * FROM tft_match WHERE id = $1
);

-- name: SetMatchesUpdated :exec
UPDATE tft_summoner SET matches_updated = $2
WHERE puuid = $1;
