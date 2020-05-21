package sim

import (
	"math/rand"
	"time"

	"github.com/lsc-chocos/choco"
)

// SensorsV0 returns list of location and speed sensors for simulation
func SensorsV0() ([]choco.Sensor, error) {
	locBuffer, err := choco.NewSensorBuffer(5)
	if err != nil {
		return []choco.Sensor{}, nil
	}
	locSensor, err := choco.NewLSCSensor(
		"location",
		locBuffer,
		choco.SensorFunc(func() choco.SensorData {
			data := choco.SensorData{}
			data["long"] = rand.Float64()
			data["lat"] = rand.Float64()
			return data
		}),
		time.Second,
		"deg",
	)
	if err != nil {
		return []choco.Sensor{}, nil
	}
	spdBuffer, err := choco.NewSensorBuffer(5)
	if err != nil {
		return []choco.Sensor{}, nil
	}

	speedSensor, err := choco.NewLSCSensor(
		"speed",
		spdBuffer,
		choco.SensorFunc(func() choco.SensorData {
			data := choco.SensorData{}
			data["speed"] = rand.Float64()
			return data
		}),
		100*time.Millisecond,
		"m/s",
	)
	if err != nil {
		return []choco.Sensor{}, nil
	}
	return []choco.Sensor{locSensor, speedSensor}, nil
}
