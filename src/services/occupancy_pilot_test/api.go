package occupancy_pilot_test

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

//OccupancyStatusHandler is used to handle HTTP requests related to lots and occupancy
type OccupancyStatusHandler struct {
	lots map[string]*Lot
}

type errMessage struct {
	Message string `json:"message"`
}

func (sh *OccupancyStatusHandler) addLot(l *Lot) error {
	if _, ok := sh.lots[l.ID]; !ok {
		sh.lots[l.ID] = l
		return nil
	}
	return fmt.Errorf("a lot with that ID already exists: %s", sh.lots[l.ID].LotName)
}

//GetLotInfo returns all the stored information about a parking lot
func (sh *OccupancyStatusHandler) GetLotInfo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if v, ok := vars["lot"]; !ok {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errMessage{Message: "you must include a lot id"})
	} else {
		if l, ok := sh.lots[v]; !ok {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(errMessage{Message: "no lot matching that id found"})
		} else {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(l)
		}
	}
}
