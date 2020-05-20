package choco

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/lsc-chocos/choco/state"
)

var _ Sensor = (*lscSensor)(nil)

// Sensor can be runned/stopped, with functions to retrieve SenML format data
type Sensor interface {
	//Name gets the name of sensor
	Name() string

	//UpdateData updates the data with Sensor data
	UpdateData(SensorData) error

	//Run collects/mocks the from actual running data
	Run()

	//SenML returns the current snapshot of in SenML format
	SenML() (string, error)

	//GetState get the state of current sensor
	GetState() (state.State, error)

	//SetState sets the state of current sensor
	SetState(state.State) error

	//Snapshot returns the snapshot of Sensor
	Snapshot() SensorData
}

type lscSensor struct {
	name       string
	buffer     SensorBuffer
	sensorFunc SensorFunc
	period     time.Duration
	state      state.State
	unit       string
}

//NewLSCSensor returns a sensor given params provided
func NewLSCSensor(name string, buffer SensorBuffer,
	sensorFunc SensorFunc, period time.Duration, unit string) (Sensor, error) {
	sensor := &lscSensor{
		name:       name,
		buffer:     buffer,
		sensorFunc: sensorFunc,
		period:     period,
		state:      state.CREATED,
		unit:       unit,
	}
	return sensor, nil
}

func (s *lscSensor) Name() string {
	return s.name
}

// UpdateData updates data to Buffer
func (s *lscSensor) UpdateData(data SensorData) error {
	s.buffer.UpdateData(data)
	return nil
}

// Run regularly updates sensorBuffer
func (s *lscSensor) Run() {
	for s.state == state.RUNNING {
		s.buffer.UpdateData(s.sensorFunc())
		time.Sleep(s.period)
	}
}

// SenML returns the sensor data in senml format
func (s *lscSensor) SenML() (string, error) {
	jsonObj := s.buffer.DumpSenML()
	jsonObj[0]["bn"] = fmt.Sprintf("%s:", s.name)
	jsonObj[0]["bu"] = s.unit
	jsonObj[0]["bver"] = 1

	jsonBytes, err := json.Marshal(jsonObj)
	jsonStr := string(jsonBytes)
	return jsonStr, err
}

func (s *lscSensor) GetState() (state.State, error) {
	return s.state, nil
}

func (s *lscSensor) SetState(state state.State) error {
	s.state = state
	return nil
}

func (s *lscSensor) Snapshot() SensorData {
	return s.buffer.Snapshot()
}
