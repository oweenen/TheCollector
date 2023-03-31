package config

import (
	"github.com/spf13/viper"
)

type Neo4jConfig struct {
	Uri      string `mapstructure:"uri"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

type MySqlConfig struct {
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Host     string `mapstructure:"host"`
	DbName   string `mapstructure:"dbname"`
}

type RiotConfig struct {
	Key            string  `mapstructure:"key"`
	RateLimit      int     `mapstructure:"rate_limit"`
	RatePeriod     int     `mapstructure:"rate_period"`
	RateEfficiency float32 `mapstructure:"rate_efficiency"`
	MatchesAfter   int64   `mapstructure:"matches_after"`
}

type Config struct {
	Neo4jConfig Neo4jConfig `mapstructure:"neo4j"`
	MySqlConfig `mapstructure:"mysql"`
	Riot        RiotConfig `mapstructure:"riot"`
}

func LoadConfig() (Config, error) {
	var config Config
	vp := viper.New()
	vp.SetConfigFile("./config/config.json")

	err := vp.ReadInConfig()
	if err != nil {
		return config, err
	}

	err = vp.Unmarshal(&config)

	return config, err
}
