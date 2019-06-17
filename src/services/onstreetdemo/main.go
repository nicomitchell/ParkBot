package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	handler "github.com/ParkBot/src/services/onstreetdemo/handler"
	"github.com/gorilla/mux"
)

func main() {
	h := handler.NewHandler(20)
	router := mux.NewRouter()
	router.HandleFunc("/state", h.GetState).Methods("GET")
	go func() {
		h.GenerateEvent()
		fmt.Println("generated")
		time.Sleep(time.Second * time.Duration(rand.Int()%6+6))
	}()
	log.Fatal(http.ListenAndServe(":8081", router))
}
