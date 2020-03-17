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
	Build(sdk.Thing, Sensor)

	// Run starts the Choco
	Run()

	// Stop stops the Choco
	Stop()

	// Oberserve continuously observe the status of Choco
	Observe()

	/*
		// UpdateStatus forces the updating of information
		UpdateStatus() error

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
	channels   []sdk.Channel
	client     *provision.Client
	status     Status
	sensor     Sensor
}

// NewChoco re
func NewChoco(conf provision.Config) (Choco, error) {
	client, err := provision.NewClient(conf)
	if err != nil {
		return nil, fmt.Errorf("client initialization failed with config %+v: %s", conf, err.Error())
	}
	return &lscChoco{client: client}, nil
}
