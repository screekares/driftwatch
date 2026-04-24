// Package mock provides a test/stub provider for use in unit tests and
// local development without requiring real cloud credentials.
package mock

import (
	"fmt"

	"github.com/yourorg/driftwatch/internal/provider"
)

func init() {
	provider.Register("mock", func(cfg map[string]string) (provider.Provider, error) {
		return &MockProvider{resources: defaultResources()}, nil
	})
}

// MockProvider returns hard-coded resource state.
type MockProvider struct {
	resources map[string]map[string]interface{}
}

func (m *MockProvider) Name() string { return "mock" }

// FetchResource returns a fake resource state keyed by "type/id".
func (m *MockProvider) FetchResource(resourceType, resourceID string) (map[string]interface{}, error) {
	key := resourceType + "/" + resourceID
	res, ok := m.resources[key]
	if !ok {
		return nil, fmt.Errorf("mock: resource %q not found", key)
	}
	return res, nil
}

// SetResource allows tests to inject custom resource state.
func (m *MockProvider) SetResource(resourceType, resourceID string, state map[string]interface{}) {
	m.resources[resourceType+"/"+resourceID] = state
}

func defaultResources() map[string]map[string]interface{} {
	return map[string]map[string]interface{}{
		"instance/web-01": {
			"instance_type": "t3.micro",
			"ami":           "ami-0abcdef1234567890",
			"region":        "us-east-1",
		},
		"bucket/assets": {
			"versioning": true,
			"region":     "us-east-1",
		},
	}
}
