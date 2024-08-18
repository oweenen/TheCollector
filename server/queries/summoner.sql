-- name: GetSummonerByPuuid :one
SELECT * FROM tft_summoner WHERE puuid = $1;

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
