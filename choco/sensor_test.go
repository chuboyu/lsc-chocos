package choco_test

import (
	"testing"
	"time"

	"github.com/lsc-chocos/choco"
	"github.com/lsc-chocos/choco/mocks"
	"github.com/lsc-chocos/choco/state"
	"github.com/magiconair/properties/assert"
	"github.com/stretchr/testify/mock"
)

func CreateTestSensor() (choco.Sensor, error) {
	msb := &mocks.SensorBuffer{}
	msb.On("UpdateData", mock.AnythingOfType("choco.SensorData")).Return(nil)
	msb.On("Snapshot").Return(choco.SensorData{}, nil)
	msb.On("DumpSenML").Return(
		[]map[string]interface{}{
			map[string]interface{}{
				"bn": "urn:dev:ow:10e2073a0108006:", "bt": 1.276020076001e+09,
				"bu": "A", "bver": 5, "n": "voltage", "u": "V", "v": 120.1,
			},
			map[string]interface{}{"n": "current", "t": -5, "v": 1.2},
			map[string]interface{}{"n": "current", "t": -4, "v": 1.3},
			map[string]interface{}{"n": "current", "t": -3, "v": 1.4},
			map[string]interface{}{"n": "current", "t": -2, "v": 1.5},
			map[string]interface{}{"n": "current", "t": -1, "v": 1.6},
			map[string]interface{}{"n": "current", "v": 1.7},
		}, nil)
	return choco.NewLSCSensor(
		"test_sensor", msb,
		choco.SensorFunc(func() choco.SensorData {
			return choco.SensorData{}
		}),
		time.Second,
		"testunit",
	)
}
func TestSensor(t *testing.T) {
	testSensor, err := CreateTestSensor()
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, testSensor.Name(), "test_sensor")
	assert.Equal(t, testSensor.GetState(), state.CREATED)
	testSensor.SetState(state.RUNNING)
	assert.Equal(t, testSensor.GetState(), state.RUNNING)

	sensorData := choco.SensorData{}
	assert.Equal(t, testSensor.UpdateData(sensorData), nil)

	senml, err := testSensor.SenML()
	assert.Equal(t, err, nil)
	senmlCheck := "[{\"bn\":\"test_sensor:\",\"bt\":1276020076.001,\"bu\":\"testunit\",\"bver\":1,\"n\":\"voltage\",\"u\":\"V\",\"v\":120.1},{\"n\":\"current\",\"t\":-5,\"v\":1.2},{\"n\":\"current\",\"t\":-4,\"v\":1.3},{\"n\":\"current\",\"t\":-3,\"v\":1.4},{\"n\":\"current\",\"t\":-2,\"v\":1.5},{\"n\":\"current\",\"t\":-1,\"v\":1.6},{\"n\":\"current\",\"v\":1.7}]"
	assert.Equal(t, senml, senmlCheck)
	sd, err := testSensor.Snapshot()
	assert.Equal(t, err, nil)
	assert.Equal(t, sd, choco.SensorData{})
}
