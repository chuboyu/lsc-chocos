package provision

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"net/http"

	sdk "github.com/lsc-chocos/mainflux/sdk/go"
)

// Provision controls channels and things
type Provision struct {
	mfxSDK    sdk.SDK
	userToken string
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

// NewProvision creates provision with http
func NewProvision(conf Config) (*Provision, error) {
	sdkConf := sdk.Config{
		BaseURL:           conf.BaseURL,
		ReaderURL:         conf.ReaderURL,
		ReaderPrefix:      conf.ReaderPrefix,
		UsersPrefix:       conf.UsersPrefix,
		ThingsPrefix:      conf.ThingsPrefix,
		HTTPAdapterPrefix: conf.HTTPAdapterPrefix,
		MsgContentType:    sdk.CTJSON,
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
		return &Provision{mfxSDK: mfxSDK, userToken: ""}, nil
	}
	mfxSDK := sdk.NewSDK(sdkConf)
	return &Provision{mfxSDK: mfxSDK, userToken: ""}, nil
}

//Version gets the provision version
func (p *Provision) Version() (string, error) {
	return p.mfxSDK.Version()
}
