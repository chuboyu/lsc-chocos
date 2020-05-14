package choco

import (
	"math/rand"
	"testing"
	"time"

	sdk "github.com/lsc-chocos/mainflux/sdk/go"
)

func TestChoco(t *testing.T) {
	choco, err := NewChoco(Config{})
	if err != nil {
		t.Errorf("Choco initial failed: %s", err.Error())
	}
	locSensor := Sensor{
		Name: "location",
		SensorFunc: SensorFunc(func() SensorData {
			data := SensorData{}
			data["long"] = rand.Float64()
			data["lat"] = rand.Float64()
			return data
		}),
		Period: time.Second,
		Buffer: NewSensorBuffer(5),
	}
	speedSensor := Sensor{
		Name: "speed",
		SensorFunc: SensorFunc(func() SensorData {
			data := SensorData{}
			data["speed"] = rand.Float64()
			return data
		}),
		Period: 100 * time.Millisecond,
		Buffer: NewSensorBuffer(5),
	}

	sensorList := []Sensor{locSensor, speedSensor}
	choco.Build(sdk.Thing{}, sensorList, []string{})
	choco.Run()
	time.Sleep(time.Second)
	t.Logf("%+v", choco.Observe())
	choco.Stop()
}
