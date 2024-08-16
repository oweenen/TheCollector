// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: query.sql

package db

import (
	"context"

	"TheCollectorDG/types"
	"github.com/jackc/pgx/v5/pgtype"
)

const createComp = `-- name: CreateComp :exec
INSERT INTO tft_comp (
    match_id,
	summoner_puuid,
	comp_data
) VALUES (
    $1, $2, $3
)
`

type CreateCompParams struct {
	MatchID       string
	SummonerPuuid string
	CompData      types.CompData
}

func (q *Queries) CreateComp(ctx context.Context, arg CreateCompParams) error {
	_, err := q.db.Exec(ctx, createComp, arg.MatchID, arg.SummonerPuuid, arg.CompData)
	return err
}

const createMatch = `-- name: CreateMatch :exec
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
)
`

type CreateMatchParams struct {
	ID          string
	DataVersion string
	GameVersion string
	QueueID     int32
	GameType    string
	SetName     string
	SetNumber   int32
}

func (q *Queries) CreateMatch(ctx context.Context, arg CreateMatchParams) error {
	_, err := q.db.Exec(ctx, createMatch,
		arg.ID,
		arg.DataVersion,
		arg.GameVersion,
		arg.QueueID,
		arg.GameType,
		arg.SetName,
		arg.SetNumber,
	)
	return err
}

const getOldestMatchHistories = `-- name: GetOldestMatchHistories :many
SELECT
    puuid,
    matches_updated
FROM tft_summoner
ORDER BY matches_updated ASC NULLS FIRST
LIMIT $1
`

type GetOldestMatchHistoriesRow struct {
	Puuid          string
	MatchesUpdated pgtype.Timestamp
}

func (q *Queries) GetOldestMatchHistories(ctx context.Context, limit int32) ([]GetOldestMatchHistoriesRow, error) {
	rows, err := q.db.Query(ctx, getOldestMatchHistories, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetOldestMatchHistoriesRow
	for rows.Next() {
		var i GetOldestMatchHistoriesRow
		if err := rows.Scan(&i.Puuid, &i.MatchesUpdated); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getPuuidsWithNullAccountData = `-- name: GetPuuidsWithNullAccountData :many
SELECT
    puuid
FROM tft_summoner
WHERE name IS NULL OR tag IS NULL
LIMIT $1
`

func (q *Queries) GetPuuidsWithNullAccountData(ctx context.Context, limit int32) ([]string, error) {
	rows, err := q.db.Query(ctx, getPuuidsWithNullAccountData, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []string
	for rows.Next() {
		var puuid string
		if err := rows.Scan(&puuid); err != nil {
			return nil, err
		}
		items = append(items, puuid)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getPuuidsWithNullSummonerData = `-- name: GetPuuidsWithNullSummonerData :many
SELECT
    puuid
FROM tft_summoner
WHERE summoner_id IS NULL
LIMIT $1
`

func (q *Queries) GetPuuidsWithNullSummonerData(ctx context.Context, limit int32) ([]string, error) {
	rows, err := q.db.Query(ctx, getPuuidsWithNullSummonerData, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []string
	for rows.Next() {
		var puuid string
		if err := rows.Scan(&puuid); err != nil {
			return nil, err
		}
		items = append(items, puuid)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const insertPuuid = `-- name: InsertPuuid :exec
INSERT INTO tft_summoner (
    puuid
) VALUES (
    $1
) ON CONFLICT (puuid) DO NOTHING
`

func (q *Queries) InsertPuuid(ctx context.Context, puuid string) error {
	_, err := q.db.Exec(ctx, insertPuuid, puuid)
	return err
}

const matchExists = `-- name: MatchExists :one
SELECT EXISTS (
    SELECT id, data_version, game_version, queue_id, game_type, set_name, set_number FROM tft_match WHERE id = $1
)
`

func (q *Queries) MatchExists(ctx context.Context, id string) (bool, error) {
	row := q.db.QueryRow(ctx, matchExists, id)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const setMatchesUpdated = `-- name: SetMatchesUpdated :exec
UPDATE tft_summoner SET matches_updated = $2
WHERE puuid = $1
`

type SetMatchesUpdatedParams struct {
	Puuid          string
	MatchesUpdated pgtype.Timestamp
}

func (q *Queries) SetMatchesUpdated(ctx context.Context, arg SetMatchesUpdatedParams) error {
	_, err := q.db.Exec(ctx, setMatchesUpdated, arg.Puuid, arg.MatchesUpdated)
	return err
}

const updateAccount = `-- name: UpdateAccount :exec
UPDATE tft_summoner
SET name = $2, tag = $3
WHERE puuid = $1
`

type UpdateAccountParams struct {
	Puuid string
	Name  *string
	Tag   *string
}

func (q *Queries) UpdateAccount(ctx context.Context, arg UpdateAccountParams) error {
	_, err := q.db.Exec(ctx, updateAccount, arg.Puuid, arg.Name, arg.Tag)
	return err
}

const updateSummoner = `-- name: UpdateSummoner :exec
UPDATE tft_summoner
SET summoner_id = $2, profile_icon_id = $3, summoner_level = $4
WHERE puuid = $1
`

type UpdateSummonerParams struct {
	Puuid         string
	SummonerID    *string
	ProfileIconID *int32
	SummonerLevel *int32
}

func (q *Queries) UpdateSummoner(ctx context.Context, arg UpdateSummonerParams) error {
	_, err := q.db.Exec(ctx, updateSummoner,
		arg.Puuid,
		arg.SummonerID,
		arg.ProfileIconID,
		arg.SummonerLevel,
	)
	return err
}
