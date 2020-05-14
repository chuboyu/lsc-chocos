package provision_test

import (
	"testing"

	sdk "github.com/lsc-chocos/mainflux/sdk/go"
	"github.com/lsc-chocos/provision"
)

func TestCreatGroup(t *testing.T) {
	c, err := provision.NewClient(CreateProvisionTestConfig())
	c.SetUser(sdk.User{Email: "boyu@test.com", Password: "testtest"})
	if err != nil {
		t.Error("Initialization Failed", err.Error())
	}
	t.Log("Initialization Done")

	testCases := []struct {
		numThings   int
		numChannels int
	}{
		{numThings: 1, numChannels: 1},
		{numThings: 10, numChannels: 10},
	}

	for _, testCase := range testCases {
		t.Logf("Test Case: %+v", testCase)
		err := c.RemoveAll()
		if err != nil {
			t.Error("Error when calling Client.RemoveAll()")
		}
		thingIDs, channelIDs, err := c.CreateGroup(testCase.numThings, testCase.numChannels)
		if err != nil {
			t.Errorf("Error Creating Group")
		}
		if len(thingIDs) != testCase.numThings {
			t.Errorf("Number of things incorrect: expect %d, exist %d", testCase.numThings, len(thingIDs))
		}
		if len(channelIDs) != testCase.numChannels {
			t.Errorf("Number of channels incorrect: expect %d, exist %d", testCase.numChannels, len(channelIDs))
		}

		for _, thingID := range thingIDs {
			thingChannelIDs, err := c.GetChannelIDs(thingID)
			if err != nil {
				t.Errorf("Error when retrieving channels with thing %s: %s", thingID, err.Error())
			}
			if len(thingChannelIDs) != testCase.numChannels {
				t.Errorf("error # of channels for thing %s: expect %d, exist %d", thingID, testCase.numChannels, len(thingChannelIDs))
			}
		}

		for _, channelID := range channelIDs {
			channelThingIDs, err := c.GetThingIDs(channelID)
			if err != nil {
				t.Errorf("Error when retrieving things with channel %s: %s", channelID, err.Error())
			}
			if len(channelThingIDs) != testCase.numThings {
				t.Errorf("error # of things with channel %s: expect %d, exist %d", channelID, testCase.numThings, len(channelThingIDs))
			}
		}
		err = c.RemoveAll()
		if err != nil {
			t.Error("Error when calling Client.RemoveAll()")
		}
		thingIDs, _ = c.GetThingIDs("")
		if len(thingIDs) != 0 {
			t.Error("c.RemoveAll() Failed")
		}
		channelIDs, _ = c.GetChannelIDs("")
		if len(channelIDs) != 0 {
			t.Error("c.RemoveAll() Failed")
		}
	}
}
