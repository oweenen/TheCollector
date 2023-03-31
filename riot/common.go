package riot

import (
	"encoding/json"
	"errors"
	"net/http"
)

func getJson(url string, target interface{}) error {
	res, err := http.Get(url)

	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return errors.New(res.Status)
	}

	return json.NewDecoder(res.Body).Decode(target)
}
