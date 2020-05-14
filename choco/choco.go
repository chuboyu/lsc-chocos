package choco

import (
	"encoding/json"
	"fmt"
	"io"
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

// Config holds the config read from file
type Config struct {
	Provision provision.Config `json:"provision"`
	User      sdk.User         `json:"user"`
	Mqtt      MqttConfig       `json:"mqtt"`
}

// ConfigsFromFile creates provision config from file
func ConfigsFromFile(configFilePath string) (Config, error) {
	file, err := os.Open(configFilePath)
	if err != nil {
		return Config{}, err
	}
	return ParseJSONConfig(file)
}

// ParseJSONConfig parses the the json from io.Reader
func ParseJSONConfig(r io.Reader) (Config, error) {
	decoder := json.NewDecoder(r)
	var chocoConf Config
	err := decoder.Decode(&chocoConf)
	if err != nil {
		return Config{}, err
	}
	return chocoConf, nil
}

// NewChoco re
func NewChoco(conf Config) (Choco, error) {
	client, err := provision.NewClient(conf.Provision)
	if err != nil {
		return nil, fmt.Errorf("client initialization failed with config %+v: %s", conf, err.Error())
	}
	return &lscChoco{client: client}, nil
}
