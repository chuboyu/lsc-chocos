package provision

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net/http"

	sdk "github.com/lsc-chocos/mainflux/sdk/go"
)

// Provision interface contains operation needed for provisioning
type Provision interface {
	// Initialize sets up the provision
	Initialize() error

	// Set User sets the user
	SetUser(user sdk.User) error

	// UpdateUserToken updates the user token
	UpdateUserToken() error

	// CreateGroup creates a group of Things and Channels
	CreateGroup(thingsData interface{}, channelsData interface{}) error

	// GetAllThingIDs retrieves all the thing ids for a given channel
	GetThingIDs(channelID string) ([]string, error)

	// GetAllChannelIDs retrieves all the channel ids for a given thing
	GetChannelIDs(thingID string) ([]string, error)

	// RemoveAll deletes everything corresponding to user from mainflux
	RemoveAll()
}

// Client controls channels and things
type Client struct {
	MfxSDK    sdk.SDK
	UserToken string
	User      sdk.User
	Things    []sdk.Thing
	Channels  []sdk.Channel
	Groups    []string
}

//Config contains provision configs
type Config struct {
	BaseURL           string
	ReaderURL         string
	ReaderPrefix      string
	UsersPrefix       string
	ThingsPrefix      string
	HTTPAdapterPrefix string
	MsgContentType    sdk.ContentType
	TLSVerification   bool
	CaFilePath        string
}

// NewClient creates provision with http
func NewClient(conf Config) (*Client, error) {
	sdkConf := sdk.Config{
		BaseURL:           conf.BaseURL,
		ReaderURL:         conf.ReaderURL,
		ReaderPrefix:      conf.ReaderPrefix,
		UsersPrefix:       conf.UsersPrefix,
		ThingsPrefix:      conf.ThingsPrefix,
		HTTPAdapterPrefix: conf.HTTPAdapterPrefix,
		MsgContentType:    conf.MsgContentType,
		TLSVerification:   conf.TLSVerification,
	}
	if conf.TLSVerification {
		caCert, err := ioutil.ReadFile(conf.CaFilePath)
		if err != nil {
			return nil, err
		}
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)
		client := &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					RootCAs: caCertPool,
				},
			},
		}
		mfxSDK := sdk.NewSDKWithClient(sdkConf, client)
		return &Client{MfxSDK: mfxSDK, UserToken: ""}, nil
	}
	mfxSDK := sdk.NewSDK(sdkConf)
	return &Client{MfxSDK: mfxSDK, UserToken: ""}, nil
}

//Version gets the provision version
func (p *Client) Version() (string, error) {
	version, err := p.MfxSDK.Version()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("mainflux version: %s", version), nil
}
