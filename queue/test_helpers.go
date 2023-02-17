package queue

import (
	"github.com/GiladLeef/SharkCrack/mocks"
	"github.com/GiladLeef/SharkCrack/models"
	"testing"
)

var threads = uint16(3)

func setupConfig() models.Config {
	return models.Config{
		Threads: threads,
	}
}

func setupConfigProvider() mocks.MockConfigProvider {
	config := setupConfig()
	return mocks.NewMockConfigProvider(&config)
}

func assertConfigProviderCalled(t *testing.T, m mocks.MockConfigProvider) {
	expected := uint64(1)

	actual := m.GetConfigCalls
	if expected != actual {
		t.Errorf("Expected: %d\nActual: %d\n", expected, actual)
	}
}
