package choco

import (
	"fmt"
	"time"

	"github.com/lsc-chocos/choco/state"
	sdk "github.com/lsc-chocos/mainflux/sdk/go"
)

// Build builds the choco using Thing
func (c *lscChoco) Build(thing sdk.Thing, sensors []Sensor, chanIDs []string) {
	c.thing = thing
	c.thingToken = c.thing.Key
	c.status = Status{State: state.CREATED}
	c.sensors = sensors
	c.channelIDs = chanIDs
	for i := range c.sensors {
		c.sensors[i].SetState(state.CREATED)
	}

}

// Run starts the choco
func (c *lscChoco) Run() {
	c.status.State = state.RUNNING
	for i := range c.sensors {
		c.sensors[i].SetState(state.RUNNING)
		go c.sensors[i].Run()
	}
}

// Stop stops the choco
func (c *lscChoco) Stop() {
	c.status.State = state.STANDBY
	for i := range c.sensors {
		c.sensors[i].SetState(state.STANDBY)
	}
}

// Observe returns the snapshots of sensors
func (c *lscChoco) Observe() map[string]SensorData {
	result := map[string]SensorData{}
	for _, sensor := range c.sensors {
		result[sensor.Name()] = sensor.Snapshot()
	}
	return result
}

// SenML returns the snapshot in SenML Strings
func (c *lscChoco) SenML() ([]string, error) {
	data := make([]string, len(c.sensors))
	var err error
	for i, sensor := range c.sensors {
		data[i], err = sensor.SenML()
		if err != nil {
			return nil, err
		}
	}
	return data, nil
}

// Observe prints the status of sensor
func (c *lscChoco) ObserveUntil() {
	for {
		fmt.Printf("Buffer: %+v\n", c.sensors)
		time.Sleep(100 * time.Millisecond)
	}
}

// SendStatus updates choco status to server
func (c *lscChoco) SendStatus() error {
	senMLStrs, err := c.SenML()
	if err != nil {
		return err
	}
	for _, chanID := range c.channelIDs {
		for _, senMLStr := range senMLStrs {
			err = c.client.MfxSDK.SendMessage(chanID, senMLStr, c.thingToken)
			if err != nil {
				return fmt.Errorf("Error sending message: %w", err)
			}
		}
	}
	return nil
}
