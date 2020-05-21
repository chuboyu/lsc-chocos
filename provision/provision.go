package provision

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"time"

	sdk "github.com/lsc-chocos/mainflux/sdk/go"
)

var _ Provision = (*lscProvision)(nil)

// Provision interface contains operation needed for provisioning
type Provision interface {

	// Set User sets the user
	SetUser(user sdk.User) error

	// UpdateUserToken updates the user token
	UpdateUserToken() error

	// CreateGroup creates a group of Things and Channels
	CreateGroup(thingsData interface{}, channelsData interface{}) ([]string, []string, error)

	// GetThing retrieves
	GetThing(thingID string) (sdk.Thing, error)

	// GetThings retrieves
	GetThings(channelID string) (sdk.ThingsPage, error)

	// GetThingIDs retrieves all the thing ids for a given channel
	GetThingIDs(channelID string) ([]string, error)

	// GetChannelIDs retrieves all the channel ids for a given thing
	GetChannelIDs(thingID string) ([]string, error)

	// RemoveAll deletes everything corresponding to user from mainflux
	RemoveAll() error

	// Version returns the version of mainflux server
	Version() (string, error)

	// SendMessage sends message to mainflux server
	SendMessage(channelID string, senMLMessage string, token string) error
}

// client controls channels and things
type lscProvision struct {
	mfxSDK    sdk.SDK
	userToken string
	user      sdk.User
	things    []sdk.Thing
	channels  []sdk.Channel
	groups    []string
}

//Config contains provision configs
type Config struct {
	SDKConf sdk.Config `json:"sdk"`
	CaCert  []byte
}

// NewProvision creates provision with http
func NewProvision(conf Config) (Provision, error) {
	var mfxSDK sdk.SDK
	if conf.SDKConf.TLSVerification {
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(conf.CaCert)
		client := &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					RootCAs: caCertPool,
				},
			},
			Timeout: 5 * time.Second,
		}
		mfxSDK = sdk.NewSDKWithClient(conf.SDKConf, client)
	} else {
		mfxSDK = sdk.NewSDK(conf.SDKConf)
	}
	return NewProvisionWithSDK(mfxSDK)
}

// NewProvisionWithSDK sets the SDK
func NewProvisionWithSDK(mfxSDK sdk.SDK) (Provision, error) {
	return &lscProvision{mfxSDK: mfxSDK}, nil
}

//Version gets the provision version
func (p *lscProvision) Version() (string, error) {
	version, err := p.mfxSDK.Version()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("mainflux version: %s", version), nil
}
