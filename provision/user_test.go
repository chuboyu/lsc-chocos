package provision

import (
	"testing"

	sdk "github.com/lsc-chocos/mainflux/sdk/go"
)

func TestUser(t *testing.T) {
	provConf := Config{
		BaseURL:           "https://localhost",
		UsersPrefix:       "",
		ThingsPrefix:      "",
		HTTPAdapterPrefix: "",
		MsgContentType:    sdk.CTJSONSenML,
		TLSVerification:   true,
		CaFilePath:        "../ssl/ca.crt",
	}
	p, _ := NewClient(provConf)

	var err error

	err = p.Initialize()
	if err != nil {
		t.Error("Initialization of Client Failed")
	}

	email := "boyu@test.com"
	password := "testtest"
	user := sdk.User{Email: email, Password: password}
	err = p.SetUser(user)
	if err != nil {
		t.Errorf("set user failed: %s, %s", email, password)
	}
	if email != p.User.Email {
		t.Errorf("set user email failed: %s", email)
	}
	if password != p.User.Password {
		t.Errorf("set user email failed: %s", password)
	}

	err = p.UpdateUserToken()
	if err != nil {
		t.Errorf("Update user token failed: %s", err.Error())
	}
	if p.UserToken == "" {
		t.Error("User token empty after update")
	}
}
