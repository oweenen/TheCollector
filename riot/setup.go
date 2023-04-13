package riot

import (
	"log"
	"os"
	"strconv"
)

var key string
var matchesAfter int64

func Setup() {
	var err error
	key = os.Getenv("RIOT_KEY")
	matchesAfter, err = strconv.ParseInt(os.Getenv("MATCHES_AFTER"), 10, 64)
	if err != nil {
		log.Fatalln("error parsing MATCHES_AFTER")
	}
}
