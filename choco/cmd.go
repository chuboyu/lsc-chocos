package choco

import (
	"encoding/json"
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

// Observe returns the snapshots of sensors
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

// SendStatus updates choco status to server
func (c *lscChoco) SendStatus() error {
	data := map[string]interface{}{}
	data["name"] = c.thing.Name
	data["id"] = c.thing.ID
	data["ts"] = time.Now().UnixNano()
	data["observe"] = c.Observe()
	result, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("Json Marshal error: %w", err)
	}
	resultStr := string(result)
	for _, chanID := range c.channelIDs {
		fmt.Printf("%s, %s, %s\n", chanID, resultStr, c.thingToken)
		err = c.client.MfxSDK.SendMessage(chanID, resultStr, c.thingToken)
		if err != nil {
			return fmt.Errorf("Error sending message: %w", err)
		}
	}
	return nil
}
