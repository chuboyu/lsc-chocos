package choco

import (
	"encoding/json"
	"fmt"
	"os"

	MQTT "github.com/eclipse/paho.mqtt.golang"
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
	mqttClient *MQTT.Client
	status     Status
	sensors    []Sensor
}

// MqttConfig is the config use for paho mqtt client
type MqttConfig struct {
	Broker string `json:"broker"`
}

type fileConfig struct {
	Provision provision.Config `json:"provision"`
	User      sdk.User         `json:"user"`
	Mqtt      MqttConfig       `json:"mqtt"`
}

// ConfigsFromFile creates provision config from file (currently no use)
func ConfigsFromFile(configFilePath string) (provision.Config, sdk.User, error) {
	file, err := os.Open(configFilePath)
	if err != nil {
		return provision.Config{}, sdk.User{}, err
	}

	decoder := json.NewDecoder(file)
	var fileConf fileConfig

	err = decoder.Decode(&fileConf)
	if err != nil {
		return provision.Config{}, sdk.User{}, err
	}

	return fileConf.Provision, fileConf.User, nil
}

// NewChoco re
func NewChoco(conf provision.Config, crtFilePath string) (Choco, error) {
	client, err := provision.NewClient(conf, crtFilePath)
	if err != nil {
		return nil, fmt.Errorf("client initialization failed with config %+v: %s", conf, err.Error())
	}
	return &lscChoco{client: client}, nil
}
