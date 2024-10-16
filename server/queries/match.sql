-- name: CreateMatch :exec
INSERT INTO tft_match (
    id,
	data_version,
	game_version,
	queue_id,
	game_type,
	set_name,
	set_number,
	match_date
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
);

-- name: CreateComp :exec
INSERT INTO tft_comp (
    match_id,
	summoner_puuid,
	comp_data
) VALUES (
    $1, $2, $3
);

-- name: MatchExists :one
SELECT EXISTS (
    SELECT * FROM tft_match WHERE id = $1
);

-- name: SummonerMatchHistory :many
SELECT
	match_id,
	comp_data,
	game_version,
	queue_id,
	game_type,
	set_number,
	match_date
FROM
	tft_comp
	JOIN tft_match ON tft_comp.match_id = tft_match.id
WHERE
	tft_comp.summoner_puuid = $1
ORDER BY
	tft_match.match_date
LIMIT $2
OFFSET $3;
