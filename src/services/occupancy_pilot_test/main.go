package main

import (
	"log"
	"net/http"

	"github.com/ParkBot/src/services/occupancy_pilot_test/pilot"
	"github.com/gorilla/mux"
)

func main() {
	sh := pilot.NewStatusHandler(make(map[string]*pilot.Lot))
	router := mux.NewRouter()
	router.HandleFunc("/healthcheck", sh.HealthCheck).Methods("GET")
	router.HandleFunc("/lots/{lot}", sh.GetLotInfo).Methods("GET")
	router.HandleFunc("/lots/new", sh.NewLot).Methods("POST")
	router.HandleFunc("/lots/{lot}/occupancy", sh.GetOccupancyInfo).Methods("GET")
	router.HandleFunc("/lots/{lot}/entry", sh.TrackEntry).Methods("POST")
	router.HandleFunc("/lots/{lot}/exit", sh.TrackExit).Methods("POST")
	log.Fatal(http.ListenAndServe(":8081", router))
}
