package pilot

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

//StatusHandler is the API interface for the occupancy pilot test
type StatusHandler interface {
	HealthCheck(w http.ResponseWriter, r *http.Request)
	GetLotInfo(w http.ResponseWriter, r *http.Request)
	NewLot(w http.ResponseWriter, r *http.Request)
	GetOccupancyInfo(w http.ResponseWriter, r *http.Request)
	TrackEntry(w http.ResponseWriter, r *http.Request)
	TrackExit(w http.ResponseWriter, r *http.Request)
}

//OccupancyStatusHandler is used to handle HTTP requests related to lots and occupancy
type OccupancyStatusHandler struct {
	lots map[string]*Lot
}

type errMessage struct {
	Message string `json:"message"`
}

//NewStatusHandler returns a new status handler
func NewStatusHandler(lots map[string]*Lot) StatusHandler {
	return &OccupancyStatusHandler{
		lots: lots,
	}
}

func (sh *OccupancyStatusHandler) addLot(l *Lot) error {
	if _, ok := sh.lots[l.ID]; !ok {
		sh.lots[l.ID] = l
		return nil
	}
	return fmt.Errorf("a lot with that ID already exists: %s", sh.lots[l.ID].LotName)
}

//HealthCheck can be used to check the status of the API
func (sh *OccupancyStatusHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
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

//NewLot creates a new lot
func (sh *OccupancyStatusHandler) NewLot(w http.ResponseWriter, r *http.Request) {
	d := r.Body
	q := r.URL.Query()
	var occupancy int
	var err error
	if vals := q.Get("start_occupancy"); len(vals) > 0 {
		occupancy, err = strconv.Atoi(string(vals[0]))
	} else {
		occupancy = 0
	}
	if err == nil {
		var data []byte
		data, err = ioutil.ReadAll(d)
		if err == nil {
			var lot *Lot
			lot, err = decodeLot(data, occupancy)
			if err == nil {
				sh.lots[lot.ID] = lot
				w.WriteHeader(http.StatusOK)
				return
			}
		}
	}
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(errMessage{Message: err.Error()})
}

//GetOccupancyInfo gets occupancy info for the lot given
func (sh *OccupancyStatusHandler) GetOccupancyInfo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if v, ok := vars["lot"]; !ok {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errMessage{Message: "you must include a lot id"})
	} else {
		if l, ok := sh.lots[v]; !ok {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(errMessage{Message: "no lot matching that id found"})
		} else {
			type resp struct {
				OpenSpots   int     `json:"open_spots"`
				Occupied    int     `json:"occupied_spots"`
				PercentFull float64 `json:"percent_full"`
			}
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(
				resp{
					OpenSpots:   l.MaxOccupancy - l.Occupancy(),
					Occupied:    l.Occupancy(),
					PercentFull: l.PercentFull(),
				})
		}
	}
}

//TrackExit tracks an exit from a lot
func (sh *OccupancyStatusHandler) TrackExit(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if v, ok := vars["lot"]; !ok {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errMessage{Message: "you must include a lot id"})
	} else {
		if l, ok := sh.lots[v]; !ok {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(errMessage{Message: fmt.Sprintf("no lot matching id: %s found", v)})
		} else {
			err := l.TrackExit()
			if err != nil {
				w.WriteHeader(http.StatusMethodNotAllowed)
				json.NewEncoder(w).Encode(errMessage{Message: "lot is empty; no exits possible"})
			} else {
				w.WriteHeader(http.StatusOK)
			}
		}
	}
}

//TrackEntry tracks an entry to a lot
func (sh *OccupancyStatusHandler) TrackEntry(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if v, ok := vars["lot"]; !ok {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errMessage{Message: "you must include a lot id"})
	} else {
		if l, ok := sh.lots[v]; !ok {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(errMessage{Message: "no lot matching that id found"})
		} else {
			err := l.TrackEntry()
			if err != nil {
				w.WriteHeader(http.StatusMethodNotAllowed)
				json.NewEncoder(w).Encode(errMessage{Message: "lot is full; no spots available"})
			} else {
				w.WriteHeader(http.StatusOK)
			}
		}
	}
}
