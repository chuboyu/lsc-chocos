package main

import (
	"context"
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
	"golang.org/x/sync/errgroup"
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
		ch, err := choco.NewChoco(pConf, *crtFilePath)
		if err != nil {
			exitWithError(err, 1, "Create Choco Failed")
		}
		thing, err := p.MfxSDK.Thing(thingID, p.UserToken)
		if err != nil {
			exitWithError(err, 1, fmt.Sprintf("retrieve thing failed for thingID: %s\n", thingID))
		}
		ch.Build(thing, sim.SensorsV0(), channelIDs)
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
					ch.SendStatus()
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
