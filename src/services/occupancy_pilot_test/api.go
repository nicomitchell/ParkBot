package occupancy_pilot_test

import (
	"fmt"
)

//OccupancyStatusHandler is used to handle HTTP requests
type OccupancyStatusHandler struct {
	lots map[string]*Lot
}

func (sh *OccupancyStatusHandler) addLot(l *Lot) error {
	if _, ok := sh.lots[l.ID]; !ok {
		sh.lots[l.ID] = l
		return nil
	}
	return fmt.Errorf("a lot with that ID already exists: %s", sh.lots[l.ID].LotName)
}
