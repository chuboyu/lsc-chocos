package main

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/lsc-chocos/choco"
	sdk "github.com/lsc-chocos/mainflux/sdk/go"
	"github.com/lsc-chocos/provision"
	"github.com/lsc-chocos/sim"
)

func initProvision(provConf provision.Config, user sdk.User) *provision.Client {
	p, _ := provision.NewClient(provConf)
	p.SetUser(user)
	p.UpdateUserToken()
	return p
}

func main() {
	pConf, user, _ := provision.ConfigsFromFile("./configs/config_dev.json")
	p := initProvision(pConf, user)
	thingIDs, channelIDs, err := p.CreateGroup(1, 1)
	if err != nil {
		fmt.Printf(err.Error())
		os.Exit(1)
	}

	chocoList := []choco.Choco{}
	for _, thingID := range thingIDs {
		ch, err := choco.NewChoco(pConf)
		if err != nil {
			fmt.Printf(err.Error())
			os.Exit(1)
		}
		thing, err := p.MfxSDK.Thing(thingID, p.UserToken)
		if err != nil {
			fmt.Printf(err.Error())
			os.Exit(1)
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
