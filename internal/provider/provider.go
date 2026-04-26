// Package provider defines the interface and registry for cloud providers.
package provider

import (
	"fmt"
	"sort"
	"sync"
)

// Resource represents a single cloud resource fetched from a provider.
type Resource struct {
	ID     string            `json:"id"`
	Type   string            `json:"type"`
	Region string            `json:"region,omitempty"`
	Labels map[string]string `json:"labels,omitempty"`
	Attrs  map[string]string `json:"attrs,omitempty"`
}

// Provider is the interface all cloud provider implementations must satisfy.
type Provider interface {
	FetchResource(resourceType, id string) (*Resource, error)
}

// FactoryFunc constructs a Provider from a config map.
type FactoryFunc func(cfg map[string]string) (Provider, error)

var (
	mu       sync.RWMutex
	registry = map[string]FactoryFunc{}
)

// Register adds a named provider factory to the global registry.
// It panics if the same name is registered twice.
func Register(name string, fn FactoryFunc) {
	mu.Lock()
	defer mu.Unlock()
	if _, exists := registry[name]; exists {
		panic(fmt.Sprintf("provider: duplicate registration for %q", name))
	}
	registry[name] = fn
}

// New creates a provider by name using the global registry.
func New(name string, cfg map[string]string) (Provider, error) {
	mu.RLock()
	fn, ok := registry[name]
	mu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("provider: unknown provider %q (available: %v)", name, AvailableNames())
	}
	return fn(cfg)
}

// AvailableNames returns a sorted list of registered provider names.
func AvailableNames() []string {
	mu.RLock()
	defer mu.RUnlock()
	names := make([]string, 0, len(registry))
	for k := range registry {
		names = append(names, k)
	}
	sort.Strings(names)
	return names
}

// availableNames is kept for backward compatibility within the package.
func availableNames() []string { return AvailableNames() }
