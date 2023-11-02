package api

import (
	"TheCollectorDG/database"
	"TheCollectorDG/riot"
	"TheCollectorDG/stats"
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

// summoner/:region/:name
func GetSummonerByName(c *fiber.Ctx) error {
	c.Set("Access-Control-Allow-Origin", "*")
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

	c.Status(200).JSON(*summoner)
	return nil
}

// summoner/:puuid
func GetSummonerByPuuid(c *fiber.Ctx) error {
	c.Set("Access-Control-Allow-Origin", "*")
	puuid := c.Params("puuid")

	summoner, err := database.GetSummonerByPuuid(puuid)
	if err != nil {
		c.SendStatus(404)
		return nil
	}

	c.Status(200).JSON(*summoner)
	return nil
}

// rank/:puuid
func GetSummonerRank(c *fiber.Ctx) error {
	c.Set("Access-Control-Allow-Origin", "*")
	puuid := c.Params("puuid")
	rank, err := database.GetRank(puuid)
	if err != nil {
		c.SendStatus(404)
		return err
	}
	c.Status(200).JSON(rank)
	return nil
}

// rank/stats/:puuid
func GetRankStats(c *fiber.Ctx) error {
	c.Set("Access-Control-Allow-Origin", "*")
	puuid := c.Params("puuid")

	stats, err := database.GetRankStats(puuid)
	if err != nil {
		c.Status(404).SendString(err.Error())
		fmt.Print(err)
		return err
	}

	c.Status(200).JSON(*stats)
	return nil
}

// update/profile/:puuid
func UpdateProfile(c *fiber.Ctx) error {
	c.Set("Access-Control-Allow-Origin", "*")
	puuid := c.Params("puuid")
	updateInfo, err := database.GetUpdateInfo(puuid)
	if err != nil {
		c.SendStatus(404)
		return err
	}

	summonerCQ, _ := summonerCollectionRegionRouter[updateInfo.Region]
	matchCQ, _ := matchCollectionRegionRouter[riot.RiotRegionRoutes[updateInfo.Region]]

	err = <-summonerCQ.QueueSummonerByPuuid(updateInfo.Puuid)
	if err != nil {
		c.SendStatus(500)
		return err
	}

	err = <-summonerCQ.QueueRank(updateInfo.Puuid, updateInfo.SummonerId)
	if err != nil {
		log.Panicln(err)
		c.SendStatus(500)
		return err
	}

	err = <-matchCQ.QueueMatchHistory(updateInfo.Puuid, updateInfo.MatchesLastUpdated)
	if err != nil {
		c.SendStatus(500)
		return err
	}

	err = database.UpdateRankStats(updateInfo.Puuid)
	if err != nil {
		c.SendStatus(500)
		return err
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
	c.Set("Access-Control-Allow-Origin", "*")
	puuid := c.Params("puuid")

	matchHistory, err := database.GetRecentMatches(puuid, 10)
	if err != nil {
		c.SendStatus(404)
		return nil
	}

	c.Status(200).JSON(matchHistory)
	return nil
}

// matches/participants/:match_id
func GetMatchParticipants(c *fiber.Ctx) error {
	c.Set("Access-Control-Allow-Origin", "*")
	matchId := c.Params("match_id")

	participants, err := database.GetMatchParticipants(matchId)
	if err != nil {
		c.Status(404).SendString(err.Error())
		fmt.Print(err)
		return err
	}

	c.Status(200).JSON(participants)
	return nil
}

// augment/stats
func GetAugmentStats(c *fiber.Ctx) error {
	c.Set("Access-Control-Allow-Origin", "*")

	c.Status(200).JSON(stats.LatestPage)
	return nil
}
