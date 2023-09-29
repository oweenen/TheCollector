package database

import (
	"TheCollectorDG/types"
	"database/sql"
	"encoding/json"
)

func GetRecentComps(puuid string, count int) ([]*types.Comp, error) {
	rows, err := db.Query(`
	SELECT
    	TFT_Match.id,
    	TFT_Match.date,
    	TFT_Match.game_length,
    	TFT_Match.game_version,
    	TFT_Match.queue_id,
    	TFT_Match.game_type,
    	TFT_Match.set_name,
    	TFT_Match.set_number,
    	Comp.placement,
    	Comp.last_round,
    	Comp.level,
    	Comp.remaining_gold,
    	Comp.players_eliminated,
    	Comp.player_damage_dealt,
    	Comp.time_eliminated,
    	Comp.companion,
    	Comp.augments,
    	Comp.traits,
    	Comp.units
	FROM TFT_Match
	JOIN (
	    SELECT
	        c.id AS comp_id,
	        c.match_id,
	        c.placement,
	        c.last_round,
	        c.level,
	        c.remaining_gold,
	        c.players_eliminated,
	        c.player_damage_dealt,
	        c.time_eliminated,
	        c.companion,
	        c.augments,
	        c.traits,
	        c.units
	    FROM Comp c
	    WHERE c.summoner_puuid = ?
	    ORDER BY c.date DESC
	    LIMIT ?
	) AS SelectedComps
	ON TFT_Match.id = SelectedComps.match_id
	ORDER BY TFT_Match.date DESC;
		`,
		puuid,
		count,
	)
	if err != nil {
		return nil, err
	}

	var comps []*types.Comp
	defer rows.Close()
	for rows.Next() {
		comp := new(types.Comp)
		comp.Match = new(types.Match)
		var augmentJson, traitJson, unitJson []byte

		err = rows.Scan(
			&comp.Match.Id,
			&comp.Match.Date,
			&comp.Match.GameLength,
			&comp.Match.GameVersion,
			&comp.Match.QueueId,
			&comp.Match.GameType,
			&comp.Match.SetName,
			&comp.Match.SetNumber,
			&comp.Placement,
			&comp.LastRound,
			&comp.Level,
			&comp.RemainingGold,
			&comp.PlayersEliminated,
			&comp.PlayerDamageDealt,
			&comp.TimeEliminated,
			&comp.Companion,
			&augmentJson,
			&traitJson,
			&unitJson,
		)
		if err != nil {
			return nil, err
		}

		json.Unmarshal(augmentJson, &comp.Augments)
		json.Unmarshal(traitJson, &comp.Traits)
		json.Unmarshal(unitJson, &comp.Units)

		comps = append(comps, comp)
	}

	for _, comp := range comps {
		participants, err := getParticipants(comp.Match.Id)
		if err != nil {
			return nil, err
		}
		comp.Match.Participants = participants
	}

	return comps, nil
}

func getParticipants(matchId string) ([]types.Participant, error) {
	rows, err := db.Query(`
		SELECT
			Summoner.region,
			Summoner.display_name
		FROM Summoner JOIN Comp
		ON Summoner.puuid = Comp.summoner_puuid
		WHERE Comp.match_id = ?
		`,
		matchId,
	)
	if err != nil {
		return nil, err
	}

	var participants []types.Participant
	defer rows.Close()
	for rows.Next() {
		var participant types.Participant
		err = rows.Scan(
			&participant.Region,
			&participant.Name,
		)
		if err != nil {
			return nil, err
		}
		participants = append(participants, participant)
	}

	return participants, nil
}

func storeComp(tx *sql.Tx, matchId string, comp *types.Comp) error {
	companionJson, _ := json.Marshal(comp.Companion)
	augmentJson, _ := json.Marshal(comp.Augments)
	traitJson, _ := json.Marshal(comp.Traits)
	unitJson, _ := json.Marshal(comp.Units)
	_, err := tx.Exec(`
		INSERT IGNORE INTO Comp (
			match_id,
			summoner_puuid,
			placement,
			last_round,
			level,
			remaining_gold,
			players_eliminated,
			player_damage_dealt,
			time_eliminated,
			companion,
			augments,
			traits,
			units
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		`,
		matchId,
		comp.Summoner.Puuid,
		comp.Placement,
		comp.LastRound,
		comp.Level,
		comp.RemainingGold,
		comp.PlayersEliminated,
		comp.PlayerDamageDealt,
		comp.TimeEliminated,
		string(companionJson),
		string(augmentJson),
		string(traitJson),
		string(unitJson),
	)
	return err
}
