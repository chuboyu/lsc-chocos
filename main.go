package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
	"time"

	"github.com/lsc-chocos/choco"
	sdk "github.com/lsc-chocos/mainflux/sdk/go"
	"github.com/lsc-chocos/provision"
	"github.com/lsc-chocos/sim"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

func initLogger() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
}

func initProvision(provConf provision.Config, user sdk.User) provision.Provision {
	p, err := provision.NewProvision(provConf)
	if err != nil {
		exitWithError(err, 1, "Initializtion Provision Failed")
	}
	p.SetUser(user)
	p.UpdateUserToken()
	return p
}

func main() {
	initLogger()

	configFilePath := flag.String("f", "config.json", "path of the config file")
	crtFilePath := flag.String("cacert", "mainflux-server.crt", "path of certificate file")
	flag.Parse()

	file, err := os.Open(*configFilePath)
	chocoConf, err := choco.ParseJSONConfig(file)
	if err != nil {
		exitWithError(err, 1, "Config Read Failed")
	}
	chocoConf.Provision.CaCert, err = ioutil.ReadFile(*crtFilePath)
	if err != nil {
		exitWithError(err, 1, "Cert Read Failed")
	}

	p := initProvision(chocoConf.Provision, chocoConf.User)
	thingIDs, channelIDs, err := p.CreateGroup(1, 1)
	if err != nil {
		exitWithError(err, 1, "Create Group Failed")
	}

	chocoList := []choco.Choco{}
	for _, thingID := range thingIDs {
		ch, err := choco.NewChoco(chocoConf)
		if err != nil {
			exitWithError(err, 1, "Create Choco Failed")
		}
		thing, err := p.GetThing(thingID)
		if err != nil {
			exitWithError(err, 1, fmt.Sprintf("retrieve thing failed for thingID: %s\n", thingID))
		}
		simSensors, err := sim.SensorsV0()
		if err != nil {
			exitWithError(err, 1, "Simulation Sensor create failed")
		}
		ch.Build(thing, simSensors, channelIDs)
		chocoList = append(chocoList, ch)
	}

	for _, ch := range chocoList {
		ch.Run()
	}

	time.Sleep(time.Second)

	sendStatuses := func(ctx context.Context, chocoList []choco.Choco) error {
		g, ctx := errgroup.WithContext(ctx)
		for _, ch := range chocoList {
			g.Go(func() error {
				for {
					err := ch.SendStatus()
					if err != nil {
						log.WithFields(log.Fields{
							"Error": err.Error(),
						}).Fatal("error sending messages")
					}
					time.Sleep(time.Second)
				}
			})
		}
		return nil
	}

	err = sendStatuses(context.Background(), chocoList)

	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()
}

func exitWithError(err error, code int, msg string) {
	log.WithFields(log.Fields{
		"Error": err.Error(),
		"code":  code,
	}).Fatal(msg)
	os.Exit(1)
}
