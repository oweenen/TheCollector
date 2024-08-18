package riot

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

var NotFoundError = errors.New("404 Not Found")

func getJson(server string, route string, target interface{}) error {
	url := fmt.Sprintf("https://%v.api.riotgames.com/%v", server, route)

	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("X-Riot-Token", key)

	limiter, ok := limiters[server]
	if !ok {
		return fmt.Errorf("no limiter defined %v", server)
	}
	<-limiter.C

	res, err := client.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()
	if res.StatusCode == http.StatusNotFound {
		return NotFoundError
	}
	if res.StatusCode != http.StatusOK {
		return errors.New(res.Status)
	}

	return json.NewDecoder(res.Body).Decode(target)
}
