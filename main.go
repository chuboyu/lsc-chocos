package main

import (
	"flag"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/lsc-chocos/choco"
	sdk "github.com/lsc-chocos/mainflux/sdk/go"
	"github.com/lsc-chocos/provision"
	"github.com/lsc-chocos/sim"
	log "github.com/sirupsen/logrus"
)

func initLogger() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
}

func initProvision(provConf provision.Config, user sdk.User, crtFilePath string) *provision.Client {
	p, err := provision.NewClient(provConf, crtFilePath)
	if err != nil {
		exitWithError(err, 1, "Initializtion Provision Failed")
	}
	p.SetUser(user)
	p.UpdateUserToken()
	return p
}

func main() {
	configFilePath := flag.String("f", "config.json", "path of the config file")
	crtFilePath := flag.String("cacert", "mainflux-server.crt", "path of certificate file")
	flag.Parse()

	pConf, user, err := provision.ConfigsFromFile(*configFilePath)
	if err != nil {
		exitWithError(err, 1, "Config Read Failed")
	}
	p := initProvision(pConf, user, *crtFilePath)
	thingIDs, channelIDs, err := p.CreateGroup(1, 1)
	if err != nil {
		exitWithError(err, 1, "Create Group Failed")
	}

	chocoList := []choco.Choco{}
	for _, thingID := range thingIDs {
		ch, err := choco.NewChoco(pConf)
		if err != nil {
			exitWithError(err, 1, "Create Choco Failed")
		}
		thing, err := p.MfxSDK.Thing(thingID, p.UserToken)
		if err != nil {
			exitWithError(err, 1, fmt.Sprintf("retrieve thing failed for thingID: %s\n", thingID))
		}
		ch.Build(thing, sim.SensorsV0(), channelIDs)
		ch.Run()
		chocoList = append(chocoList, ch)
	}
	time.Sleep(time.Second)
	for _, ch := range chocoList {
		go func(cho choco.Choco) {
			for {
				cho.SendStatus()
				time.Sleep(time.Second)
			}
		}(ch)
	}

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
