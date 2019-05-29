package pilot_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"github.com/ParkBot/src/mocks"
	"github.com/ParkBot/src/services/occupancy_pilot_test/pilot"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

type errReader int

func (*errReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("test error")
}

func Test_NewStatusHandler_ReturnsStatusHandlerInterface(t *testing.T) {
	sh := pilot.NewStatusHandler(make(map[string]*pilot.Lot))
	assert.Implements(t, (*pilot.StatusHandler)(nil), sh)
}

func Test_HealthCheck_WritesCorrectHeader(t *testing.T) {
	w := mocks.NewMockResponseWriter(false)
	sh := pilot.NewStatusHandler(make(map[string]*pilot.Lot))
	sh.HealthCheck(w, nil)
	assert.Equal(t, http.StatusOK, (w.(*mocks.MockResponseWriter)).StatusHeader)
}

func Test_GetLotInfo_OnError_ReturnsError(t *testing.T) {
	w := mocks.NewMockResponseWriter(false)
	sh := pilot.NewStatusHandler(make(map[string]*pilot.Lot))
	r, _ := http.NewRequest("GET", "/lots", nil)
	sh.GetLotInfo(w, r)
	assert.Equal(t, http.StatusBadRequest, (w.(*mocks.MockResponseWriter)).StatusHeader)
	assert.Equal(t, "{\"message\":\"you must include a lot id\"}\n", string((w.(*mocks.MockResponseWriter)).Message))
}

func Test_GetLotInfo_OnMapError_ReturnsError(t *testing.T) {
	w := mocks.NewMockResponseWriter(false)
	sh := pilot.NewStatusHandler(make(map[string]*pilot.Lot))
	r, _ := http.NewRequest("GET", "/lots", nil)
	r = mux.SetURLVars(r, map[string]string{"lot": "1"})
	sh.GetLotInfo(w, r)
	assert.Equal(t, http.StatusNotFound, (w.(*mocks.MockResponseWriter)).StatusHeader)
	assert.Equal(t, "{\"message\":\"no lot matching that id found\"}\n", string((w.(*mocks.MockResponseWriter)).Message))
}

func Test_GetLotInfo_Normally_ReturnsNoError(t *testing.T) {
	w := mocks.NewMockResponseWriter(false)
	sh := pilot.NewStatusHandler(map[string]*pilot.Lot{"1": pilot.NewLot(0.0, 0.0, "1", "test-lot", 100, 0)})
	r, _ := http.NewRequest("GET", "/lots", nil)
	r = mux.SetURLVars(r, map[string]string{"lot": "1"})
	sh.GetLotInfo(w, r)
	assert.Equal(t, http.StatusOK, (w.(*mocks.MockResponseWriter)).StatusHeader)
	assert.Equal(
		t,
		"{\"location\":{},\"id\":\"1\",\"name\":\"test-lot\",\"max_occupancy\":100}\n",
		string((w.(*mocks.MockResponseWriter)).Message),
	)
}

func Test_GetOccupancyInfo_OnError_ReturnsError(t *testing.T) {
	w := mocks.NewMockResponseWriter(false)
	sh := pilot.NewStatusHandler(make(map[string]*pilot.Lot))
	r, _ := http.NewRequest("GET", "/lots", nil)
	sh.GetOccupancyInfo(w, r)
	assert.Equal(t, http.StatusBadRequest, (w.(*mocks.MockResponseWriter)).StatusHeader)
	assert.Equal(t, "{\"message\":\"you must include a lot id\"}\n", string((w.(*mocks.MockResponseWriter)).Message))
}

func Test_GetOccupancyInfo_OnMapError_ReturnsError(t *testing.T) {
	w := mocks.NewMockResponseWriter(false)
	sh := pilot.NewStatusHandler(make(map[string]*pilot.Lot))
	r, _ := http.NewRequest("GET", "/lots", nil)
	r = mux.SetURLVars(r, map[string]string{"lot": "1"})
	sh.GetOccupancyInfo(w, r)
	assert.Equal(t, http.StatusNotFound, (w.(*mocks.MockResponseWriter)).StatusHeader)
	assert.Equal(t, "{\"message\":\"no lot matching that id found\"}\n", string((w.(*mocks.MockResponseWriter)).Message))
}

func Test_GetOccupancyInfo_Normally_ReturnsNoError(t *testing.T) {
	w := mocks.NewMockResponseWriter(false)
	sh := pilot.NewStatusHandler(map[string]*pilot.Lot{"1": pilot.NewLot(0.0, 0.0, "1", "test-lot", 100, 0)})
	r, _ := http.NewRequest("GET", "/lots", nil)
	r = mux.SetURLVars(r, map[string]string{"lot": "1"})
	sh.GetOccupancyInfo(w, r)
	assert.Equal(t, http.StatusOK, (w.(*mocks.MockResponseWriter)).StatusHeader)
	assert.Equal(
		t,
		"{\"open_spots\":100,\"occupied_spots\":0,\"percent_full\":0}\n",
		string((w.(*mocks.MockResponseWriter)).Message),
	)

}

func Test_NewLot_OnAtoiError_ReturnsError(t *testing.T) {
	w := mocks.NewMockResponseWriter(false)
	sh := pilot.NewStatusHandler(make(map[string]*pilot.Lot))
	r, _ := http.NewRequest("GET", "/lots?start_occupancy=n", nil)
	sh.NewLot(w, r)
	assert.Equal(t, http.StatusBadRequest, (w.(*mocks.MockResponseWriter)).StatusHeader)
	assert.Equal(t, "{\"message\":\"strconv.Atoi: parsing \\\"n\\\": invalid syntax\"}\n", string((w.(*mocks.MockResponseWriter)).Message))
}

func Test_NewLot_OnReadError_ReturnsError(t *testing.T) {
	w := mocks.NewMockResponseWriter(false)
	sh := pilot.NewStatusHandler(make(map[string]*pilot.Lot))
	r, _ := http.NewRequest("GET", "/lots", new(errReader))
	sh.NewLot(w, r)
	assert.Equal(t, http.StatusBadRequest, (w.(*mocks.MockResponseWriter)).StatusHeader)
	assert.Equal(t, "{\"message\":\"test error\"}\n", string((w.(*mocks.MockResponseWriter)).Message))
}

func Test_NewLot_OnDecodeLotError_ReturnsError(t *testing.T) {
	w := mocks.NewMockResponseWriter(false)
	sh := pilot.NewStatusHandler(make(map[string]*pilot.Lot))
	b := []byte{}
	r, _ := http.NewRequest("GET", "/lots", bytes.NewReader(b))
	sh.NewLot(w, r)
	assert.Equal(t, http.StatusBadRequest, (w.(*mocks.MockResponseWriter)).StatusHeader)
	assert.Equal(t, "{\"message\":\"unexpected end of JSON input\"}\n", string((w.(*mocks.MockResponseWriter)).Message))

}

func Test_NewLot_Normally_ReturnsNoError(t *testing.T) {
	w := mocks.NewMockResponseWriter(false)
	sh := pilot.NewStatusHandler(make(map[string]*pilot.Lot))
	lot := pilot.NewLot(0.0, 0.0, "1", "test-lot", 100, 0)
	b, _ := json.Marshal(lot)
	r, _ := http.NewRequest("GET", "/lots/?start_occupancy=15", bytes.NewReader(b))
	sh.NewLot(w, r)
	assert.Equal(t, http.StatusOK, (w.(*mocks.MockResponseWriter)).StatusHeader)
	assert.Empty(t, (w.(*mocks.MockResponseWriter)).Message)
}
