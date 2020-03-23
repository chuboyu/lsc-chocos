package provision

import (
	"testing"
)

func TestProvision(t *testing.T) {
	provConf, _, err := ConfigsFromFile("../configs/config_test.json")
	p, err := NewClient(provConf)
	result, err := p.Version()
	if err != nil {
		t.Errorf("%e", err)
	}
	t.Log(result)
}
