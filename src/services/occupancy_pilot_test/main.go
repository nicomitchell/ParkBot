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
	router.HandleFunc("/lots/{lot}", sh.GetLotInfo).Methods("GET")
	router.HandleFunc("/lots/new", sh.NewLot).Methods("POST")
	router.HandleFunc("/lots/{lot}/occupancy", sh.GetOccupancyInfo).Methods("GET")

	log.Fatal(http.ListenAndServe(":8081", router))
}
