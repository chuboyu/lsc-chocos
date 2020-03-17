package choco

import (
	"fmt"
	"time"

	"github.com/lsc-chocos/choco/state"
	sdk "github.com/lsc-chocos/mainflux/sdk/go"
)

// Build builds the choco using Thing
func (c *lscChoco) Build(thing sdk.Thing, sensors []Sensor) {
	//c.Thing = thing.(sdk.Thing)
	//c.ThingToken = c.Thing.Key
	c.status = Status{State: state.CREATED}
	c.sensors = sensors
	for i := range c.sensors {
		c.sensors[i].State = state.CREATED
	}

}

// Run starts the choco
func (c *lscChoco) Run() {
	c.status.State = state.RUNNING
	for i := range c.sensors {
		c.sensors[i].State = state.RUNNING
		go c.sensors[i].Run()
	}
}

// Stop stops the choco
func (c *lscChoco) Stop() {
	c.status.State = state.STANDBY
	for i := range c.sensors {
		c.sensors[i].State = state.STANDBY
	}
}

func (c *lscChoco) Observe() map[string]SensorData {
	result := map[string]SensorData{}
	for _, sensor := range c.sensors {
		result[sensor.Name] = sensor.Buffer.Snapshot()
	}
	return result
}

// Observe prints the status of sensor
func (c *lscChoco) ObserveUntil() {
	for {
		fmt.Printf("Buffer: %+v\n", c.sensors)
		time.Sleep(100 * time.Millisecond)
	}
}
