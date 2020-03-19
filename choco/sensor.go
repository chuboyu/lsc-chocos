package choco

import (
	"container/ring"
	"encoding/json"
	"fmt"
	"time"

	"github.com/lsc-chocos/choco/state"
)

// SensorData has stored attributes
type SensorData map[string]float64

// SensorFunc type defines function returning the latest data
type SensorFunc func() SensorData

// SensorBuffer is a ring buffer storing sensor data
type SensorBuffer struct {
	ringBuffer *ring.Ring
}

// NewSensorBuffer returns an empty sensor buffer
func NewSensorBuffer(length int) *SensorBuffer {
	r := ring.New(length)
	for i := 0; i < length; i++ {
		r.Value = SensorData{}
		r = r.Next()
	}
	return &SensorBuffer{ringBuffer: r}
}

// UpdateData updates sensor data to buffer
func (s *SensorBuffer) UpdateData(data SensorData) {
	s.ringBuffer = s.ringBuffer.Next()
	s.ringBuffer.Value = data
}

// Snapshot returns the latest observed numbers
func (s *SensorBuffer) Snapshot() SensorData {
	return s.ringBuffer.Value.(SensorData)
}

// DumpSenML dumps all the data in the buffer into senml json
func (s *SensorBuffer) DumpSenML() []map[string]interface{} {
	obj := []map[string]interface{}{}
	for i := 0; i < s.ringBuffer.Len(); i++ {
		data := s.ringBuffer.Value.(SensorData)
		for k, v := range data {
			obj = append(obj, map[string]interface{}{
				"n": k,
				"v": v,
			})
		}
		s.ringBuffer = s.ringBuffer.Next()
	}
	return obj
}

// Sensor consist of
type Sensor struct {
	Name       string
	Buffer     *SensorBuffer
	SensorFunc SensorFunc
	Period     time.Duration
	State      state.State
	Unit       string
}

// UpdateData updates data to Buffer
func (s *Sensor) UpdateData(data SensorData) {
	fmt.Printf("%+v\n", data)
	s.Buffer.UpdateData(data)
}

// Run regularly updates sensorBuffer
func (s *Sensor) Run() {
	for s.State == state.RUNNING {
		s.Buffer.UpdateData(s.SensorFunc())
		time.Sleep(s.Period)
	}
}

// SenML returns the sensor data in senml format
func (s *Sensor) SenML() (string, error) {
	jsonObj := s.Buffer.DumpSenML()
	jsonObj[0]["bn"] = s.Name
	jsonObj[0]["bu"] = s.Unit
	jsonObj[0]["bt"] = float64(time.Now().UnixNano()) / float64(1e9)
	jsonObj[0]["bver"] = 1

	jsonBytes, err := json.Marshal(jsonObj)
	jsonStr := string(jsonBytes)
	return jsonStr, err
}
