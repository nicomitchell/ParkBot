package onstreetdemo

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"sync"
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
		if i < 9 {
			h.spots[fmt.Sprintf("L000%d", i+1)] = false
		} else {
			h.spots[fmt.Sprintf("L00%d", i+1)] = false

		}
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

type coord struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type state struct {
	ID       string  `json:"id"`
	Occupied bool    `json:"state"`
	Price    float64 `json:"price"`
	Location coord 	 `json:"location"`
}

//GetState returns the current state of the lot
func (h *handler) GetState(w http.ResponseWriter, r *http.Request) {
	locationMap := map[string]coord{
		"L0001" : coord{X:38.219635, Y:-85.763125},
		"L0002" : coord{X:38.219685, Y:-85.763116}, 
		"L0003" : coord{X:38.219733, Y:-85.763108}, 
		"L0004" : coord{X:38.219805, Y:-85.763092}, 
		"L0005" : coord{X:38.219894, Y:-85.763082}, 
		"L0006" : coord{X:38.219975, Y:-85.763068}, 
		"L0007" : coord{X:38.220070, Y:-85.763051}, 
		"L0008" : coord{X:38.220159, Y:-85.763036}, 
		"L0009" : coord{X:38.220263, Y:-85.763027}, 
		"L0010" : coord{X:38.220782, Y:-85.762562}, 
		"L0011" : coord{X:38.220776, Y:-85.762504}, 
		"L0012" : coord{X:38.220771, Y:-85.762426}, 
		"L0013" : coord{X:38.220764, Y:-85.762363}, 
		"L0014" : coord{X:38.220757, Y:-85.762302}, 
		"L0015" : coord{X:38.220742, Y:-85.762167},
		"L0016" : coord{X:38.220734, Y:-85.762097},
		"L0017" : coord{X:38.220729, Y:-85.762034},
		"L0018" : coord{X:38.220721, Y:-85.761962},
		"L0019" : coord{X:38.220715, Y:-85.761899},
		"L0020" : coord{X:38.220702, Y:-85.761812},
	}
	current := []state{}
	h.Lock()
	for s, o := range h.spots {
		spot := state{ID: s, Occupied: o, Price: 2.50,Location:locationMap[s]}
		current = append(current, spot)
	}
	h.Unlock()
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(current)
}
