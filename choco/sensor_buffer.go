package choco

import (
	"container/ring"
	"time"
)

// sensorBufferNode is the node in sensor buffer ring
type sensorBufferNode struct {
	SensorData SensorData
	Metadata   map[string]interface{}
}

var _ SensorBuffer = (*lscSensorBuffer)(nil)

//SensorBuffer stores a window of SensorData
type SensorBuffer interface {
	UpdateData(data SensorData) error
	Snapshot() (SensorData, error)
	DumpSenML() ([]map[string]interface{}, error)
}

// SensorBuffer is a ring buffer storing sensor data
type lscSensorBuffer struct {
	ringBuffer *ring.Ring
}

// NewSensorBuffer returns an empty sensor buffer
func NewSensorBuffer(length int) (SensorBuffer, error) {
	r := ring.New(length)
	for i := 0; i < length; i++ {
		r.Value = &sensorBufferNode{
			SensorData: SensorData{},
			Metadata:   map[string]interface{}{},
		}
		r = r.Next()
	}
	return &lscSensorBuffer{ringBuffer: r}, nil
}

// UpdateData updates sensor data to buffer
func (s *lscSensorBuffer) UpdateData(data SensorData) error {
	s.ringBuffer = s.ringBuffer.Next()
	node := s.ringBuffer.Value.(*sensorBufferNode)
	node.SensorData = data
	node.Metadata["ts"] = float64(time.Now().UnixNano()) / float64(1e9)
	return nil
}

// Snapshot returns the latest observed numbers
func (s *lscSensorBuffer) Snapshot() (SensorData, error) {
	return s.ringBuffer.Value.(*sensorBufferNode).SensorData, nil
}

// DumpSenML dumps all the data in the buffer into senml json
func (s *lscSensorBuffer) DumpSenML() ([]map[string]interface{}, error) {
	obj := []map[string]interface{}{}
	for i := 0; i < s.ringBuffer.Len(); i++ {
		node := s.ringBuffer.Value.(*sensorBufferNode)
		data := node.SensorData
		meta := node.Metadata
		for k, v := range data {
			obj = append(obj, map[string]interface{}{
				"n": k,
				"v": v,
				"t": meta["ts"],
			})
		}
		s.ringBuffer = s.ringBuffer.Next()
	}
	return obj, nil
}
