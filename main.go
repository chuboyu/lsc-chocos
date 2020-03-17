package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/lsc-chocos/choco"
	sdk "github.com/lsc-chocos/mainflux/sdk/go"
	"github.com/lsc-chocos/provision"
)

func main() {
	provConf := provision.Config{
		BaseURL:           "https://localhost",
		UsersPrefix:       "",
		ThingsPrefix:      "",
		HTTPAdapterPrefix: "",
		MsgContentType:    sdk.CTJSONSenML,
		TLSVerification:   true,
		CaFilePath:        "ssl/ca.crt",
	}
	ch, err := choco.NewChoco(provConf)
	if err != nil {
		fmt.Printf(err.Error())
		os.Exit(1)
	}
	locSensor := choco.Sensor{
		SensorFunc: choco.SensorFunc(func() choco.SensorData {
			data := choco.SensorData{}
			data["long"] = rand.Float64()
			data["lat"] = rand.Float64()
			return data
		}),
		Period: time.Second,
	}
	ch.Build(sdk.Thing{}, locSensor)
	go func() {
		ch.Run()
		time.Sleep(2 * time.Second)
		ch.Stop()
	}()
	ch.Observe()
}
