package api

import (
	"TheCollectorDG/database"
	"TheCollectorDG/riot"
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

// summoner/:region/:name
func GetSummonerByName(c *fiber.Ctx) error {
	region := strings.ToLower(c.Params("region"))
	name, err := url.QueryUnescape(c.Params("name"))
	if err != nil {
		c.Status(400)
		return nil
	}

	// summoner in db
	summoner, err := database.GetSummonerByName(region, name)
	if err != nil {
		summonerCQ, ok := summonerCollectionRegionRouter[region]
		if !ok {
			c.SendStatus(400)
			return nil
		}
		err = <-summonerCQ.QueueSummonerByName(name)
		if err != nil {
			c.SendStatus(404)
			return nil
		}
		summoner, err = database.GetSummonerByName(region, name)
		if err != nil {
			c.SendStatus(500)
			return nil
		}
	}

	rank, err := database.GetRank(summoner.Puuid)
	if err == nil {
		summoner.Rank = rank
	}

	c.Status(200).JSON(*summoner)
	return nil
}

// summoner/:puuid
func GetSummonerByPuuid(c *fiber.Ctx) error {
	puuid := c.Params("puuid")

	summoner, err := database.GetSummonerByPuuid(puuid)
	if err != nil {
		c.SendStatus(404)
		return nil
	}

	rank, err := database.GetRank(summoner.Puuid)
	if err == nil {
		summoner.Rank = rank
	}

	c.Status(200).JSON(*summoner)
	return nil
}

// rank/:puuid
func GetSummonerRank(c *fiber.Ctx) error {
	puuid := c.Params("puuid")
	rank, err := database.GetRank(puuid)
	if err != nil {
		c.SendStatus(404)
		return err
	}
	c.Status(200).JSON(rank)
	return nil
}

// update/profile/:puuid
func UpdateProfile(c *fiber.Ctx) error {
	c.Set("Access-Control-Allow-Origin", "*")
	puuid := c.Params("puuid")
	updateInfo, err := database.GetUpdateInfo(puuid)
	if err != nil {
		c.SendStatus(404)
		return nil
	}

	summonerCQ, ok := summonerCollectionRegionRouter[updateInfo.Region]
	if !ok {
		c.SendStatus(400)
		return nil
	}
	matchCQ, ok := matchCollectionRegionRouter[riot.RiotRegionRoutes[updateInfo.Region]]
	if !ok {
		c.SendStatus(500)
		return nil
	}

	err = <-summonerCQ.QueueSummonerByPuuid(updateInfo.Puuid)
	if err != nil {
		c.SendStatus(500)
		return nil
	}

	err = <-summonerCQ.QueueRank(updateInfo.Puuid, updateInfo.SummonerId)
	if err != nil {
		log.Panicln(err)
		c.SendStatus(500)
		return nil
	}

	err = <-matchCQ.QueueMatchHistory(updateInfo.Puuid, updateInfo.MatchesLastUpdated)
	if err != nil {
		c.SendStatus(500)
		return nil
	}

	err = database.SetLastUpdated(updateInfo.Puuid, time.Now().Unix())
	if err != nil {
		c.SendStatus(500)
		return nil
	}

	c.Status(200).SendString("done")
	return nil
}

// matches/:puuid
func GetMatchHistory(c *fiber.Ctx) error {
	puuid := c.Params("puuid")

	matchHistory, err := database.GetRecentMatches(puuid, 10)
	if err != nil {
		c.SendStatus(404)
		return nil
	}

	c.Status(200).JSON(matchHistory)
	return nil
}

// matches/stats/:puuid
func GetSummonerStats(c *fiber.Ctx) error {
	puuid := c.Params("puuid")

	stats, err := database.GetMatchStats(puuid)
	if err != nil {
		c.Status(404).SendString(err.Error())
		fmt.Print(err)
		return err
	}

	c.Status(200).JSON(*stats)
	return nil
}
