package onstreetdemo

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

//Handler is the interface for the demo api
type Handler interface {
	GenerateEvent()
	GetState(http.ResponseWriter, *http.Request)
}

type handler struct {
	numSpots int
	spots    map[string]bool
	sync.Mutex
}

//NewHandler returns a new Handler with the appropriate number of spots
func NewHandler(spots int) Handler {
	h := &handler{
		numSpots: spots,
		spots:    make(map[string]bool),
	}
	for i := 0; i < spots; i++ {
		h.spots[fmt.Sprintf("L00%2d", i+1)] = false
	}
	return h
}

//GenerateEvent may generate a change in the state of one of the spots
func (h *handler) GenerateEvent() {
	h.Lock()
	if rand.Int()%10 <= 3 {
		toggle := rand.Int() % h.numSpots
		h.spots[fmt.Sprintf("L00%2d", toggle+1)] = !h.spots[fmt.Sprintf("L00%2d", toggle+1)]
	}
	h.Unlock()
}

type state struct {
	ID       string `json:"id"`
	Occupied bool   `json:"occupied"`
	Time     string `json:"time"`
}

//GetState returns the current state of the lot
func (h *handler) GetState(w http.ResponseWriter, r *http.Request) {
	current := make([]state, h.numSpots)
	h.Lock()
	for s, o := range h.spots {
		spot := state{ID: s, Occupied: o, Time: time.Now().String()}
		current = append(current, spot)
	}
	h.Unlock()
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(current)
}
