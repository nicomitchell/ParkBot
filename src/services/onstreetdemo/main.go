package main

import (
	"log"
	"net/http"
	"time"

	handler "github.com/ParkBot/src/services/onstreetdemo/handler"
	"github.com/gorilla/mux"
)

func main() {
	h := handler.NewHandler(10)
	router := mux.NewRouter()
	router.HandleFunc("/state", h.GetState).Methods("GET")
	go func() {
		time.Sleep(10 * time.Second)
		h.GenerateEvent()
	}()
	log.Fatal(http.ListenAndServe(":8081", router))
}
