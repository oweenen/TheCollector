package api

import (
	"TheCollectorDG/database"
	"net/url"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

// summoner/:region/:name
func GetSummoner(c *fiber.Ctx) error {
	region := strings.ToUpper(c.Params("region"))
	name, err := url.QueryUnescape(c.Params("name"))
	if err != nil {
		c.Status(400)
		return nil
	}

	// summoner in db
	summoner, err := database.GetSummoner(region, name)
	if err == nil {
		c.Status(200).JSON(*summoner)
		return nil
	}

	// summoner not in db
	err = <-summonerCollectionQueue.QueueSummonerByName(region, name)
	if err != nil {
		c.SendStatus(404)
		return nil
	}
	summoner, err = database.GetSummoner(region, name)
	if err != nil {
		c.SendStatus(500)
		return nil
	}

	c.Status(200).JSON(*summoner)
	return nil
}

// update/profile/:puuid
func UpdateProfile(c *fiber.Ctx) error {
	puuid := c.Params("puuid")
	updateInfo, err := database.GetUpdateInfo(puuid)
	if err != nil {
		c.SendStatus(404)
		return nil
	}

	err = <-summonerCollectionQueue.QueueSummonerByPuuid(updateInfo.Region, updateInfo.Puuid)
	if err != nil {
		c.SendStatus(500)
		return nil
	}

	err = <-matchCollectionQueue.QueueMatchHistory(updateInfo.Region, updateInfo.Puuid, updateInfo.MatchesLastUpdated)
	if err != nil {
		c.SendStatus(500)
		return nil
	}

	err = database.SetLastUpdated(updateInfo.Puuid, time.Now().Unix())
	if err != nil {
		c.SendStatus(500)
		return nil
	}

	c.SendStatus(200)
	return nil
}

// comps/:puuid
func GetCompHistory(c *fiber.Ctx) error {
	puuid := c.Params("puuid")

	compHistory, err := database.GetRecentComps(puuid, 10)
	if err != nil {
		c.SendStatus(404)
		return nil
	}

	c.Status(200).JSON(compHistory)
	return nil
}
