package datastore

import "TheCollectorDG/types"

func StoreMatch(match *types.Match) error {
	// store comps
	for _, comp := range match.Comps {
		err := storeComp(match.Id, &comp)
		if err != nil {
			return err
		}
	}

	return nil
}
