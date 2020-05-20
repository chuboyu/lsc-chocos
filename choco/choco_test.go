package choco_test

import (
	"testing"

	"github.com/lsc-chocos/choco"
	"github.com/lsc-chocos/choco/mocks"
	sdk "github.com/lsc-chocos/mainflux/sdk/go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestChoco(t *testing.T) {
	ch, err := choco.NewChoco(choco.Config{})
	assert.Equal(t, err, nil, "")
	mockSensorList := make([]choco.Sensor, 5)
	for i := range mockSensorList {
		ms := &mocks.Sensor{}
		ms.On("SetState", mock.AnythingOfType("state.State")).Return(nil)
		ms.On("Run").Return()
		mockSensorList[i] = ms
	}
	ch.Build(sdk.Thing{}, mockSensorList, []string{})
	ch.Run()
	ch.Stop()
}
