package onstreetdemo

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

type Handler struct {
	numSpots int
	spots    map[int]bool
	sync.Mutex
}

//NewHandler returns a new Handler with the appropriate number of spots
func NewHandler(spots int) *Handler {
	h := &Handler{
		numSpots: spots,
		spots:    make(map[int]bool),
	}
	for i := 0; i < spots; i++ {
		h.spots[i] = false
	}
	return h
}

func (h *Handler) GenerateEvent() {
	h.Lock()
	if rand.Int()%10 <= 3 {
		toggle := rand.Int() % h.numSpots
		h.spots[toggle] = !h.spots[toggle]
	}
	h.Unlock()
}

type state struct {
	Occupied   []int  `json:"occupied"`
	Unoccupied []int  `json:"unoccupied"`
	Time       string `json:"time"`
}

func (h *Handler) GetState(w http.ResponseWriter, r *http.Request) {
	current := state{Occupied: []int{}, Unoccupied: []int{}, Time: time.Now().String()}
	h.Lock()
	for s, o := range h.spots {
		if o {
			current.Occupied = append(current.Occupied, s)
		} else {
			current.Unoccupied = append(current.Unoccupied, s)
		}
	}
	h.Unlock()
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(current)
}
