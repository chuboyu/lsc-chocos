package choco

import (
	"fmt"
	"time"

	"github.com/lsc-chocos/choco/state"
	sdk "github.com/lsc-chocos/mainflux/sdk/go"
)

// Build builds the choco using Thing
func (c *lscChoco) Build(thing sdk.Thing, sensor Sensor) {
	//c.Thing = thing.(sdk.Thing)
	//c.ThingToken = c.Thing.Key
	c.status = Status{State: state.CREATED}
	c.sensor = sensor
}

// Run starts the choco
func (c *lscChoco) Run() {
	c.status.State = state.RUNNING
	c.sensor.State = state.RUNNING
	go c.sensor.Run()
}

// Stop stops the choco
func (c *lscChoco) Stop() {
	c.status.State = state.STANDBY
	c.sensor.State = state.STANDBY
}

// Observe prints the status of sensor
func (c *lscChoco) Observe() {
	for {
		fmt.Printf("Buffer: %+v\n", c.sensor.Buffer)
		time.Sleep(100 * time.Millisecond)
	}
}
