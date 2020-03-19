package sim

import (
	"math/rand"
	"time"

	"github.com/lsc-chocos/choco"
)

// SensorsV0 returns list of location and speed sensors for simulation
func SensorsV0() []choco.Sensor {
	locSensor := choco.Sensor{
		Name: "location",
		SensorFunc: choco.SensorFunc(func() choco.SensorData {
			data := choco.SensorData{}
			data["long"] = rand.Float64()
			data["lat"] = rand.Float64()
			return data
		}),
		Unit:   "deg",
		Period: time.Second,
		Buffer: choco.NewSensorBuffer(5),
	}
	speedSensor := choco.Sensor{
		Name: "speed",
		SensorFunc: choco.SensorFunc(func() choco.SensorData {
			data := choco.SensorData{}
			data["speed"] = rand.Float64()
			return data
		}),
		Unit:   "m/s",
		Period: 100 * time.Millisecond,
		Buffer: choco.NewSensorBuffer(5),
	}
	return []choco.Sensor{locSensor, speedSensor}
}
