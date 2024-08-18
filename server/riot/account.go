package riot

import (
	"fmt"
)

type Account struct {
	Puuid string `json:"puuid"`
	Name  string `json:"gameName"`
	Tag   string `json:"tagLine"`
}

func GetAccountByName(cluster string, name string, tag string) (*Account, error) {
	accountRes := new(Account)
	route := fmt.Sprintf("riot/account/v1/accounts/by-riot-id/%v/%v", name, tag)
	err := getJson(cluster, route, accountRes)
	if err != nil {
		return nil, err
	}

	return accountRes, err
}

func GetAccountByPuuid(cluster string, puuid string) (*Account, error) {
	accountRes := new(Account)
	route := fmt.Sprintf("riot/account/v1/accounts/by-puuid/%v", puuid)
	err := getJson(cluster, route, accountRes)
	if err != nil {
		return nil, err
	}

	return accountRes, err
}
