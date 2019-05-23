package occupancy_pilot_test

import "errors"

type tracker struct {
	occupancy    int
	maxOccupancy int
}

func (t *tracker) trackEntry() error {
	if t.occupancy < t.maxOccupancy {
		t.occupancy++
		return nil
	}
	return errors.New("max occupancy reached")
}

func (t *tracker) trackExit() error {
	if t.occupancy > 0 {
		t.occupancy--
		return nil
	}
	return errors.New("lot is empty")
}

type coord struct {
	Lat  float64 `json:"latitude"`
	Long float64 `json:"longitude"`
}

//Lot is a parking lot
type Lot struct {
	Location         coord  `json:"location"`
	ID               string `json:"id"`
	LotName          string `json:"name"`
	MaxOccupancy     int    `json:"max_occupancy"`
	occupancyTracker *tracker
}

//NewLot returns a new Lot object
func NewLot(lat, long float64, id, lotName string, maxOccupancy, startOccupancy int) *Lot {
	return &Lot{
		Location:         coord{Lat: lat, Long: long},
		ID:               id,
		LotName:          lotName,
		MaxOccupancy:     maxOccupancy,
		occupancyTracker: &tracker{maxOccupancy: maxOccupancy, occupancy: startOccupancy},
	}
}

//TrackEntry should be called when a car enters the lot
func (l *Lot) TrackEntry() error {
	return l.occupancyTracker.trackEntry()
}

//TrackExit should be called when a car exits the lot
func (l *Lot) TrackExit() error {
	return l.occupancyTracker.trackExit()
}

//Occupancy tells you how many cars are currently in the lot
func (l *Lot) Occupancy() int {
	return l.occupancyTracker.occupancy
}

//PercentFull tells you how full the lot is on a scale of 0 to 1
func (l *Lot) PercentFull() float64 {
	return (float64(l.occupancyTracker.occupancy) / float64(l.MaxOccupancy))
}
