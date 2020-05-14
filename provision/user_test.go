package provision_test

import (
	"testing"

	sdk "github.com/lsc-chocos/mainflux/sdk/go"
	"github.com/lsc-chocos/provision"
)

func TestUser(t *testing.T) {
	p, err := provision.NewClient(CreateProvisionTestConfig())

	err = p.Initialize()
	if err != nil {
		t.Error("Initialization of Client Failed")
	}

	user := sdk.User{Email: "boyu@test.com", Password: "testtest"}
	err = p.SetUser(user)
	if err != nil {
		t.Errorf("set user failed: %s, %s", user.Email, user.Password)
	}
	if user.Email != p.User.Email {
		t.Errorf("set user email failed: %s", user.Email)
	}
	if user.Password != p.User.Password {
		t.Errorf("set user email failed: %s", user.Password)
	}

	err = p.UpdateUserToken()
	if err != nil {
		t.Errorf("Update user token failed: %s", err.Error())
	}
	if p.UserToken == "" {
		t.Error("User token empty after update")
	}
}
