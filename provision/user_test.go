package provision

import (
	"testing"
)

func TestUser(t *testing.T) {
	var err error

	provConf, user, _ := ConfigsFromFile("../configs/config_test.json")
	p, _ := NewClient(provConf)

	err = p.Initialize()
	if err != nil {
		t.Error("Initialization of Client Failed")
	}

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
