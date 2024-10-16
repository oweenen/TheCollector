CREATE TABLE tft_match (
	id VARCHAR PRIMARY KEY NOT NULL,
	data_version VARCHAR NOT NULL,
	game_version VARCHAR NOT NULL,
	queue_id INT NOT NULL,
	game_type VARCHAR NOT NULL,
	set_name VARCHAR NOT NULL,
	set_number INT NOT NULL,
	match_date TIMESTAMP NOT NULL
);

CREATE TABLE tft_summoner (
	puuid VARCHAR PRIMARY KEY NOT NULL,
	name VARCHAR,
	tag VARCHAR,
	summoner_id VARCHAR,
	profile_icon_id INT,
	summoner_level INT,
	full_update_timestamp TIMESTAMP,
	background_update_timestamp TIMESTAMP,
	skip_account BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE TABLE tft_comp (
	match_id VARCHAR NOT NULL REFERENCES tft_match,
	summoner_puuid VARCHAR NOT NULL REFERENCES tft_summoner,
	comp_data JSONB NOT NULL,
	PRIMARY KEY(match_id, summoner_puuid)
);
