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

-- name: MatchExists :one
SELECT EXISTS (
    SELECT * FROM tft_match WHERE id = $1
);
