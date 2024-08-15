-- name: InsertPuuid :exec
INSERT INTO tft_summoner (
    puuid
) VALUES (
    $1
) ON CONFLICT (puuid) DO NOTHING;

-- name: UpdateSummoner :exec
UPDATE tft_summoner
SET summoner_id = $2, profile_icon_id = $3, summoner_level = $4
WHERE puuid = $1;

-- name: UpdateAccount :exec
UPDATE tft_summoner
SET name = $2, tag = $3
WHERE puuid = $1;

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

-- name: CreateComp :exec
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

-- name: GetPuuidsWithNullSummonerData :many
SELECT
    puuid
FROM tft_summoner
WHERE summoner_id IS NULL
LIMIT $1;

-- name: GetPuuidsWithNullAccountData :many
SELECT
    puuid
FROM tft_summoner
WHERE name IS NULL OR tag IS NULL
LIMIT $1;

-- name: MatchExists :one
SELECT EXISTS (
    SELECT * FROM tft_match WHERE id = $1
);

-- name: SetMatchesUpdated :exec
UPDATE tft_summoner SET matches_updated = $2
WHERE puuid = $1;
