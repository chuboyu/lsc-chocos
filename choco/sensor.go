package choco

import (
	"fmt"
	"time"

	"github.com/lsc-chocos/choco/state"
)

// SensorData has stored attributes
type SensorData map[string]float64

// SensorFunc type defines function returning the latest data
type SensorFunc func() SensorData

// SensorBuffer store sensor data
type sensorBuffer struct {
	Data SensorData
}

// UpdateData updates sensor data to buffer
func (s *sensorBuffer) UpdateData(data SensorData) {
	if s.Data == nil {
		s.Data = data
	} else {
		for k, v := range data {
			s.Data[k] = v
		}
	}
}

// Snapshot returns the latest observed numbers
func (s *sensorBuffer) Snapshot() SensorData {
	return s.Data
}

// Sensor consist of
type Sensor struct {
	Name       string
	Buffer     sensorBuffer
	SensorFunc SensorFunc
	Period     time.Duration
	State      state.State
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
