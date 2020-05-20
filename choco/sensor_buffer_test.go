package choco_test

import (
	"testing"
	"time"

	"github.com/lsc-chocos/choco"
	"github.com/stretchr/testify/assert"
)

const (
	BUFFER_SIZE int = 5
	TEST_SIZE   int = 10
)

func CreateTestSensorBuffer() (choco.SensorBuffer, error) {
	return choco.NewSensorBuffer(BUFFER_SIZE)
}
func CreateTestSensorFunc() choco.SensorFunc {
	return choco.SensorFunc(func() choco.SensorData {
		data := choco.SensorData{}
		data["testCol"] = float64(time.Now().UnixNano() / int64(time.Millisecond))
		return data
	})
}
func CreateTestSensorData(i float64) choco.SensorData {
	return choco.SensorData{
		"testCol0": i,
		"testCol1": i + 1,
	}
}

func TestSensorBuffer(t *testing.T) {
	tsb, err := CreateTestSensorBuffer()
	assert.Equal(t, err, nil)
	out, err := tsb.Snapshot()
	assert.Equal(t, err, nil)
	assert.Equal(t, out, choco.SensorData{})

	sensorFunc := CreateTestSensorFunc()

	for i := 0; i < TEST_SIZE; i++ {
		sensorData := sensorFunc()
		err = tsb.UpdateData(sensorData)
		assert.Equal(t, err, nil)
		sout, err := tsb.Snapshot()
		assert.Equal(t, err, nil)
		assert.Equal(t, sout, sensorData)

		senmlOut, err := tsb.DumpSenML()
		assert.Equal(t, err, nil)
		if i+1 < BUFFER_SIZE {
			assert.Equal(t, len(senmlOut), i+1)
		} else {
			assert.Equal(t, len(senmlOut), BUFFER_SIZE)
		}
	}
	for j := 0; j < TEST_SIZE; j++ {
		sensorData := CreateTestSensorData(float64(j))
		tsb.UpdateData(sensorData)
		sout, _ := tsb.Snapshot()
		assert.Equal(t, sout, sensorData)
	}
	//dumpedSenML, _ := tsb.DumpSenML()
	//TODO to best tested
}
