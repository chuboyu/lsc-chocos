package provision

import (
	"testing"

	sdk "github.com/lsc-chocos/mainflux/sdk/go"
)

func TestProvision(t *testing.T) {
	provConf := Config{
		BaseURL:           "https://localhost",
		UsersPrefix:       "",
		ThingsPrefix:      "",
		HTTPAdapterPrefix: "",
		MsgContentType:    sdk.CTJSON,
		TLSVerification:   true,
		CaFilePath:        "../ssl/ca.crt",
	}
	p, err := NewClient(provConf)
	result, err := p.Version()
	if err != nil {
		t.Errorf("%e", err)
	}
	t.Log(result)
}
