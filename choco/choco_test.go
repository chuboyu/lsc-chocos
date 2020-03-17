package choco

import (
	"math/rand"
	"testing"
	"time"

	sdk "github.com/lsc-chocos/mainflux/sdk/go"
	"github.com/lsc-chocos/provision"
)

func TestChoco(t *testing.T) {
	provConf := provision.Config{
		BaseURL:           "https://localhost",
		UsersPrefix:       "",
		ThingsPrefix:      "",
		HTTPAdapterPrefix: "",
		MsgContentType:    sdk.CTJSONSenML,
		TLSVerification:   true,
		CaFilePath:        "../ssl/ca.crt",
	}
	choco, err := NewChoco(provConf)
	if err != nil {
		t.Errorf("Choco initial failed: %s", err.Error())
	}
	locSensor := Sensor{
		Name: "location",
		SensorFunc: SensorFunc(func() SensorData {
			data := SensorData{}
			data["long"] = rand.Float64()
			data["lat"] = rand.Float64()
			return data
		}),
		Period: time.Second,
	}
	speedSensor := Sensor{
		Name: "speed",
		SensorFunc: SensorFunc(func() SensorData {
			data := SensorData{}
			data["speed"] = rand.Float64()
			return data
		}),
		Period: time.Millisecond,
	}

	sensorList := []Sensor{locSensor, speedSensor}
	choco.Build(sdk.Thing{}, sensorList)
	choco.Run()
	time.Sleep(time.Second)
	t.Logf("%+v", choco.Observe())
	choco.Stop()
}
