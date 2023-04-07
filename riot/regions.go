package riot

var RiotRegionRoutes = map[string]string{
	"na1":  "americas",
	"br1":  "americas",
	"la1":  "americas",
	"la2":  "americas",
	"kr":   "asia",
	"jp1":  "asia",
	"eun1": "europe",
	"euw1": "europe",
	"tr1":  "europe",
	"ru":   "europe",
	"ph2":  "sea",
	"sg2":  "sea",
	"th2":  "sea",
	"tw2":  "sea",
	"vn2":  "sea",
	"oc1":  "sea",
}

var RiotRegionClusters = map[string][]string{
	"americas": {"na1", "br1", "la1", "la2"},
	"asia":     {"kr", "jp1"},
	"europe":   {"eun1", "euw1", "tr1", "ru"},
	"sea":      {"ph2", "sg2", "th2", "tw2", "vn2", "oc1"},
}
