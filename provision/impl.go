package provision

import (
	"fmt"

	"github.com/gofrs/uuid"
	sdk "github.com/lsc-chocos/mainflux/sdk/go"
)

// SetUser sets user and updates usertoken
func (c *lscProvision) SetUser(user sdk.User) error {
	c.user = user
	return c.UpdateUserToken()
}

// UpdateUserToken updates the user token
func (c *lscProvision) UpdateUserToken() error {
	token, err := c.mfxSDK.CreateToken(c.user)
	if err != nil {
		return err
	}
	c.userToken = token
	return nil
}

func (c *lscProvision) SendMessage(channelID string, senMLMessage string, token string) error {
	return c.mfxSDK.SendMessage(channelID, senMLMessage, token)
}

// CreateGroup creates a group of fully connected biparted Things and Channels
func (c *lscProvision) CreateGroup(thingsData interface{}, channelsData interface{}) ([]string, []string, error) {
	groupUUID, err := uuid.NewV4()
	if err != nil {
		return []string{}, []string{}, fmt.Errorf("generating groupID failed: %s", err.Error())
	}
	groupID := groupUUID.String()
	c.groups = append(c.groups, groupID)

	things := buildThings(thingsData, groupID)
	things, err = c.mfxSDK.CreateThings(things, c.userToken)
	if err != nil {
		return []string{}, []string{}, fmt.Errorf("SDK.CreateThings Failed: %w", err)
	}

	channels := buildChannels(channelsData, groupID)
	channels, err = c.mfxSDK.CreateChannels(channels, c.userToken)
	if err != nil {
		return []string{}, []string{}, fmt.Errorf("SDK.CreateChannels Failed: %w", err)
	}

	thingIDs := make([]string, len(things))
	for i, thing := range things {
		thingIDs[i] = thing.ID
	}
	channelIDs := make([]string, len(channels))
	for i, channel := range channels {
		channelIDs[i] = channel.ID
	}

	connections := sdk.ConnectionIDs{
		ChannelIDs: channelIDs,
		ThingIDs:   thingIDs,
	}

	err = c.mfxSDK.Connect(connections, c.userToken)
	if err != nil {
		return []string{}, []string{}, err
	}

	return thingIDs, channelIDs, nil
}

// GetThing gets a thing from thingID
func (c *lscProvision) GetThing(thingID string) (sdk.Thing, error) {
	return c.mfxSDK.Thing(thingID, c.userToken)
}

// GetThings gets all the things drom channelID
// TODO pagination is not handled
func (c *lscProvision) GetThings(channelID string) (sdk.ThingsPage, error) {
	if channelID == "" {
		return c.mfxSDK.Things(c.userToken, 0, 100, "")
	}
	return c.mfxSDK.ThingsByChannel(c.userToken, channelID, 0, 10)
}

// GetThingIDs gets all the thingIDs
func (c *lscProvision) GetThingIDs(channelID string) ([]string, error) {
	thingsPage, err := c.GetThings(channelID)
	if err != nil {
		return []string{}, fmt.Errorf("Error retireving things")
	}
	allThings := thingsPage.Things
	thingIDs := make([]string, len(allThings))
	for i, thing := range allThings {
		thingIDs[i] = thing.ID
	}
	return thingIDs, nil
}

// GetChannelIDs gets all the thingIDs
func (c *lscProvision) GetChannelIDs(thingID string) ([]string, error) {
	var allChannels []sdk.Channel
	var channelsPage sdk.ChannelsPage
	var err error

	if thingID == "" {
		channelsPage, err = c.mfxSDK.Channels(c.userToken, 0, 10, "")
		if err != nil {
			return nil, fmt.Errorf("Error calling SDK.Channels(): %s", err.Error())
		}
	} else {
		channelsPage, err = c.mfxSDK.ChannelsByThing(c.userToken, thingID, 0, 100)
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
func (c *lscProvision) RemoveAll() error {
	thingsIDs, _ := c.GetThingIDs("")
	for _, id := range thingsIDs {
		err := c.mfxSDK.DeleteThing(id, c.userToken)
		if err != nil {
			return err
		}
	}
	channelIDs, _ := c.GetChannelIDs("")
	for _, id := range channelIDs {
		err := c.mfxSDK.DeleteChannel(id, c.userToken)
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
