package riot

var key string
var matchesAfter int64

func Setup(riotKey string, after int64) {
	key = riotKey
	matchesAfter = after
}
