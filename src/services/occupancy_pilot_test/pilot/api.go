package pilot

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

//StatusHandler is the API interface for the occupancy pilot test
type StatusHandler interface {
	GetLotInfo(w http.ResponseWriter, r *http.Request)
	NewLot(w http.ResponseWriter, r *http.Request)
	GetOccupancyInfo(w http.ResponseWriter, r *http.Request)
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
		data := []byte{}
		_, err = d.Read(data)
		if err == nil {
			lot, err := decodeLot(data, occupancy)
			if err == nil {
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(lot)
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
