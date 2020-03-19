package choco

import (
	"fmt"

	"github.com/lsc-chocos/choco/state"
	sdk "github.com/lsc-chocos/mainflux/sdk/go"
	"github.com/lsc-chocos/provision"
)

var _ Choco = (*lscChoco)(nil)

// Choco is an interface that continuously shooting statuses
type Choco interface {
	// Build builds up the Choco
	Build(sdk.Thing, []Sensor, []string)

	// Run starts the Choco
	Run()

	// Stop stops the Choco
	Stop()

	// Observe returns the snapshots of sensors
	Observe() map[string]SensorData

	// SenML returns the current snapshot in SenML string
	SenML() ([]string, error)

	// ObserveUntil continuously observe the status of Choco
	ObserveUntil()

	// SendStatus updates choco status to server
	SendStatus() error

	/*
		UpdateUntil()

		// TearDown destroys the choco
		TearDown()
	*/
}

// Status represent the current status of Choco
type Status struct {
	State state.State
}

// LscChoco is a thing in the field
type lscChoco struct {
	thing      sdk.Thing
	thingToken string
	channelIDs []string
	client     *provision.Client
	status     Status
	sensors    []Sensor
}

// NewChoco re
func NewChoco(conf provision.Config) (Choco, error) {
	client, err := provision.NewClient(conf)
	if err != nil {
		return nil, fmt.Errorf("client initialization failed with config %+v: %s", conf, err.Error())
	}
	return &lscChoco{client: client}, nil
}
