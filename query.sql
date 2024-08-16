-- name: InsertPuuid :exec
INSERT INTO tft_summoner (
    puuid
) VALUES (
    $1
) ON CONFLICT (puuid) DO NOTHING;

-- name: UpdateSummoner :exec
UPDATE tft_summoner
SET summoner_id = @summoner_id::VARCHAR,
    profile_icon_id = @profile_icon_id::INT,
    summoner_level = @summoner_level::INT
WHERE puuid = $1;

-- name: UpdateAccount :exec
UPDATE tft_summoner
SET name = @name::VARCHAR,
    tag = @tag::VARCHAR
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
    background_update_timestamp
FROM tft_summoner
ORDER BY background_update_timestamp ASC NULLS FIRST
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

-- name: SetBackgroundUpdateTimestamp :exec
UPDATE tft_summoner
SET background_update_timestamp = @background_update_timestamp::TIMESTAMP
WHERE puuid = $1;

-- name: GetSummonerByPuuid :one
SELECT * FROM tft_summoner WHERE puuid = $1;
