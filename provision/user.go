package provision

import sdk "github.com/lsc-chocos/mainflux/sdk/go"

// Initialize currently does nothing
func (c *Client) Initialize() error {
	return nil
}

// SetUser sets user and updates usertoken
func (c *Client) SetUser(user sdk.User) error {
	c.User = user
	return c.UpdateUserToken()
}

// UpdateUserToken updates the user token
func (c *Client) UpdateUserToken() error {
	token, err := c.MfxSDK.CreateToken(c.User)
	if err != nil {
		return err
	}
	c.UserToken = token
	return nil
}
