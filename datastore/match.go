package datastore

import "TheCollectorDG/types"

func StoreMatch(match *types.Match) {
	// store comps
	for _, comp := range match.Comps {
		storeComp(match.Id, &comp)
	}
}
