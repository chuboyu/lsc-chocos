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

func initProvision(chocoConf choco.Config) *provision.Client {
	p, err := provision.NewClient(chocoConf.Provision)
	if err != nil {
		exitWithError(err, 1, "Initializtion Provision Failed")
	}
	p.SetUser(chocoConf.User)
	p.UpdateUserToken()
	return p
}

func main() {
	initLogger()

	configFilePath := flag.String("f", "config.json", "path of the config file")
	crtFilePath := flag.String("cacert", "mainflux-server.crt", "path of certificate file")
	flag.Parse()

	chocoConf, err := choco.ConfigsFromFile(*configFilePath)
	if err != nil {
		exitWithError(err, 1, "Config Read Failed")
	}
	chocoConf.Provision.CaCert, err = ioutil.ReadFile(*crtFilePath)
	if err != nil {
		exitWithError(err, 1, "Cert Read Failed")
	}

	p := initProvision(chocoConf)
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
		thing, err := p.MfxSDK.Thing(thingID, p.UserToken)
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
