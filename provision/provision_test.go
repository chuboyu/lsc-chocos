package provision

import (
	"testing"

	"github.com/lsc-chocos/choco"
)

func TestProvision(t *testing.T) {
	provConf, _, err := choco.ConfigsFromFile("../configs/config_test.json")
	p, err := NewClient(provConf, "../ssl/mainflux-server.crt")
	result, err := p.Version()
	if err != nil {
		t.Errorf("%e", err)
	}
	t.Log(result)
}
