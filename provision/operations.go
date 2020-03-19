package provision

import (
	"fmt"

	"github.com/gofrs/uuid"
	sdk "github.com/lsc-chocos/mainflux/sdk/go"
)

// CreateGroup creates a group of fully connected biparted Things and Channels
func (c *Client) CreateGroup(thingsData interface{}, channelsData interface{}) error {
	groupUUID, err := uuid.NewV4()
	if err != nil {
		return fmt.Errorf("generating groupID failed: %s", err.Error())
	}
	groupID := groupUUID.String()
	c.Groups = append(c.Groups, groupID)

	things := buildThings(thingsData, groupID)
	c.MfxSDK.CreateThings(things, c.UserToken)

	channels := buildChannels(channelsData, groupID)
	c.MfxSDK.CreateChannels(channels, c.UserToken)

	thingIDs, _ := c.GetThingIDs("")
	channelIDs, _ := c.GetChannelIDs("")

	connections := sdk.ConnectionIDs{
		ChannelIDs: channelIDs,
		ThingIDs:   thingIDs,
	}

	err = c.MfxSDK.Connect(connections, c.UserToken)
	if err != nil {
		return err
	}

	return nil
}

// GetThingIDs gets all the thingIDs
func (c *Client) GetThingIDs(channelID string) ([]string, error) {
	var allThings []sdk.Thing
	var thingsPage sdk.ThingsPage
	var err error

	if channelID == "" {
		thingsPage, err = c.MfxSDK.Things(c.UserToken, 0, 100, "")
		if err != nil {
			return nil, fmt.Errorf("Error calling SDK.Things(): %s", err.Error())
		}
	} else {
		thingsPage, err = c.MfxSDK.ThingsByChannel(c.UserToken, channelID, 0, 10)
		if err != nil {
			return nil, fmt.Errorf("Error calling SDK.ThingsByChannel(): %s", err.Error())
		}
	}

	allThings = thingsPage.Things
	thingIDs := make([]string, len(allThings))
	for i, thing := range allThings {
		thingIDs[i] = thing.ID
	}
	return thingIDs, nil
}

// GetChannelIDs gets all the thingIDs
func (c *Client) GetChannelIDs(thingID string) ([]string, error) {
	var allChannels []sdk.Channel
	var channelsPage sdk.ChannelsPage
	var err error

	if thingID == "" {
		channelsPage, err = c.MfxSDK.Channels(c.UserToken, 0, 10, "")
		if err != nil {
			return nil, fmt.Errorf("Error calling SDK.Channels(): %s", err.Error())
		}
	} else {
		channelsPage, err = c.MfxSDK.ChannelsByThing(c.UserToken, thingID, 0, 100)
		if err != nil {
			return nil, fmt.Errorf("Error calling SDK.Channels(): %s", err.Error())
		}
	}

	allChannels = channelsPage.Channels
	channelIDs := make([]string, len(allChannels))
	for i, channel := range allChannels {
		channelIDs[i] = channel.ID
	}
	return channelIDs, nil
}

// RemoveAll removes everything for current user
func (c *Client) RemoveAll() error {
	thingsIDs, _ := c.GetThingIDs("")
	for _, id := range thingsIDs {
		err := c.MfxSDK.DeleteThing(id, c.UserToken)
		if err != nil {
			return err
		}
	}
	channelIDs, _ := c.GetChannelIDs("")
	for _, id := range channelIDs {
		err := c.MfxSDK.DeleteChannel(id, c.UserToken)
		if err != nil {
			return err
		}
	}
	return nil
}

func buildThings(thingsData interface{}, groupID string) []sdk.Thing {
	var things []sdk.Thing
	switch v := thingsData.(type) {
	case int:
		things = make([]sdk.Thing, v)
		for i := 0; i < v; i++ {
			thing := sdk.Thing{
				Name:     fmt.Sprintf("Provision-thing-%d", v),
				Metadata: map[string]interface{}{"group_id": groupID},
			}
			things[i] = thing
		}
	}
	return things
}

func buildChannels(channelsData interface{}, groupID string) []sdk.Channel {
	var channels []sdk.Channel
	switch v := channelsData.(type) {
	case int:
		channels = make([]sdk.Channel, v)
		for i := 0; i < v; i++ {
			channel := sdk.Channel{
				Name:     fmt.Sprintf("Provision-channel-%d", v),
				Metadata: map[string]interface{}{"group_id": groupID},
			}
			channels[i] = channel
		}
	}
	return channels
}
