package provision_test

import (
	"testing"

	sdk "github.com/lsc-chocos/mainflux/sdk/go"
	"github.com/lsc-chocos/provision"
	"github.com/lsc-chocos/provision/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func CreateTestProvision() (provision.Provision, error) {
	mockSDK := &mocks.SDK{}
	mockSDK.On("Version").Return("testversion", nil)
	mockSDK.On("CreateToken", mock.AnythingOfType("sdk.User")).Return("testtoken", nil)
	mockSDK.On("CreateThings", mock.AnythingOfType("[]sdk.Thing"),
		"testtoken").Return(
		[]sdk.Thing{sdk.Thing{ID: "thing_1"}},
		nil,
	)
	mockSDK.On("CreateChannels", mock.AnythingOfType("[]sdk.Channel"),
		"testtoken").Return(
		[]sdk.Channel{sdk.Channel{ID: "channel_1"}},
		nil,
	)
	mockSDK.On("Connect", mock.AnythingOfType("sdk.ConnectionIDs"),
		"testtoken").Return(nil)

	return provision.NewProvisionWithSDK(mockSDK)
}
func TestProvision(t *testing.T) {
	prov, err := CreateTestProvision()
	assert.Equal(t, nil, err)

	version, err := prov.Version()
	assert.Equal(t, nil, err)
	assert.Equal(t, "mainflux version: testversion", version)

	assert.Equal(t, nil, prov.SetUser(sdk.User{}))
	assert.Equal(t, nil, prov.UpdateUserToken())

	thingIDs, channelIDs, err := prov.CreateGroup(5, 5)
	assert.Equal(t, nil, err)
	assert.Equal(t, []string{"thing_1"}, thingIDs)
	assert.Equal(t, []string{"channel_1"}, channelIDs)
}

func CreateIntegrationTestConfig() provision.Config {
	config := provision.Config{}
	config.SDKConf = sdk.Config{
		BaseURL:           "https://ec2-18-179-20-188.ap-northeast-1.compute.amazonaws.com",
		ReaderURL:         "",
		ReaderPrefix:      "",
		UsersPrefix:       "",
		ThingsPrefix:      "",
		HTTPAdapterPrefix: "http",
		MsgContentType:    sdk.CTJSONSenML,
		TLSVerification:   true,
	}

	config.CaCert = []byte(`-----BEGIN CERTIFICATE-----
MIIEXTCCA0UCCQCqmJ7f2TN1WDANBgkqhkiG9w0BAQsFADBXMRIwEAYDVQQDDAls
b2NhbGhvc3QxETAPBgNVBAoMCE1haW5mbHV4MQwwCgYDVQQLDANJb1QxIDAeBgkq
hkiG9w0BCQEWEWluZm9AbWFpbmZsdXguY29tMB4XDTIwMDMyMDA1NTMyMVoXDTIy
MTIxNTA1NTMyMVowgYkxPzA9BgNVBAMMNmVjMi0xOC0xNzktMjAtMTg4LmFwLW5v
cnRoZWFzdC0xLmNvbXB1dGUuYW1hem9uYXdzLmNvbTERMA8GA1UECgwITWFpbmZs
dXgxETAPBgNVBAsMCG1haW5mbHV4MSAwHgYJKoZIhvcNAQkBFhFpbmZvQG1haW5m
bHV4LmNvbTCCAiIwDQYJKoZIhvcNAQEBBQADggIPADCCAgoCggIBALP3NpuauSjF
XHCVotl+3FsBeZZqSOGFouPNyxaueCtVwM9BAUzEigsBZI5vj20zbPXAJSfZDaYw
Wmf1tahWXGAe9qukQzY/TsWXNtORSMUSWSP74RseWY/UOqfrk+NuBl03m7gubryW
lZfli45gUhYlr/D46p2LuC8NloYDmZfKJgTg63eEXD8leJ04PLn4Ej9kvOc7wLBH
9eulyLAlH2Q6NfZ/qXqvdmkfJkzvVsHC99/nbSpi60NdnJf/V1dZ3f+dxpMsitl4
mEHJ+qjisfe4ZlSwz0qXi2r0q3qRUM0UtO63xTX/LYPvRe8Kao6wvsk4Z6wJbZFK
g0keHYxTw3jNoJTIoUJHbrYSad92WsCgc4zLYi13KJBxjQHTrYEzit/0EKqANLIm
Wj9zOFfTObAmREGCvRky7Gvil2DPi1Yp9C8xKsRIqNOVuIIB9gFPAVC2XK1SVli9
T10ZW182rGrkPGtDSgsDEFLyV9gU3xITTo5l7B+ec8cKW/Zd033j9LgFQRdsZjN1
1WBfT8Oq0F0TJa4n338lBFGN5FvvPbOV+F3WQNtmwuTeV7ANhhY4aE2AmGlF+nIy
0DYf4PI7vloAupVbwLNWhx3i9KJiFfKHPqGXkMfxX3sg+XR7Tt6f+cPOLAPDwjML
ajjbI0cFYWauzgd4Wzv6tNAPf3audyadAgMBAAEwDQYJKoZIhvcNAQELBQADggEB
AJxwvECDHvg+GUWlcAR/0wd7EOBBEGEoprQLo+fveZFrItiAWbd4Tj4V/gPAht7y
zqr0HAQUvFOPTfLLjExW9sWiolipCPmrz/hsc3ONdPH52PlBx4X2YmYFAJiNbeWl
Amuiy/BDUb+dcXOyAsYaNgErE5pP/wKF2rblHWrPCfp6NqkvpsA4GbZa6D+HNKmH
xDsceqiOFgfpz3/NFkIDeZ0DlZ7w8fMWrEonq4FmWzKle/xmpRYliQcETcVQadU+
V0kMnv8yE4VrVxWoE5gM9AYdVriGPrFzrEUDt/NBzNHxRk3i5js4awUTMG8h+SHl
5HIadPW7iZCIn1agFv+TQmA=
-----END CERTIFICATE-----
`)
	return config
}
