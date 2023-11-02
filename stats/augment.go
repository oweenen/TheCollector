package stats

import (
	"TheCollectorDG/database"
	"fmt"
	"sort"
	"time"
)

var StatPages map[string]database.AugmentStatsPage
var LatestPage database.AugmentStatsPage

func AugmentStatsRefreshLoop(interval time.Duration) {
	updateAugmentStats()
	for range time.Tick(interval) {
		updateAugmentStats()
	}
}

func updateAugmentStats() {
	newStats, err := database.GetAugmentStats()
	if err != nil {
		fmt.Printf("Failed to refresh augment stats\n\tERROR: %v\n", err.Error())
		return
	}

	// get latest game version
	gameVersions := make([]string, 0, len(newStats))
	for k := range newStats {
		gameVersions = append(gameVersions, k)
	}
	sort.Sort(sort.Reverse(sort.StringSlice(gameVersions)))
	latestGameVersion := gameVersions[0]

	LatestPage = newStats[latestGameVersion]
	StatPages = newStats

	fmt.Println("Successfully updated augment stats!")
}
