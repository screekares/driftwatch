// Package provider defines the Provider interface and a registry for
// named provider implementations.
package provider

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"sync"
)

// ErrNotFound is returned when a resource cannot be located by the provider.
var ErrNotFound = errors.New("resource not found")

// Provider fetches live resource state from an infrastructure backend.
type Provider interface {
	FetchResource(ctx context.Context, id string) (map[string]string, error)
}

// ProviderFunc is a constructor that returns a ready-to-use Provider.
type ProviderFunc func() (Provider, error)

var (
	mu       sync.RWMutex
	registry = map[string]ProviderFunc{}
)

// Register adds a named provider constructor to the global registry.
// It panics if the same name is registered twice.
func Register(name string, fn ProviderFunc) {
	mu.Lock()
	defer mu.Unlock()
	if _, ok := registry[name]; ok {
		panic(fmt.Sprintf("provider %q already registered", name))
	}
	registry[name] = fn
}

// New creates a provider by name. Returns an error for unknown names.
func New(name string) (Provider, error) {
	mu.RLock()
	fn, ok := registry[name]
	mu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("unknown provider %q (available: %s)", name, availableNames())
	}
	return fn()
}

// availableNames returns a sorted, comma-separated list of registered names.
func availableNames() string {
	mu.RLock()
	defer mu.RUnlock()
	names := make([]string, 0, len(registry))
	for n := range registry {
		names = append(names, n)
	}
	sort.Strings(names)
	result := ""
	for i, n := range names {
		if i > 0 {
			result += ", "
		}
		result += n
	}
	return result
}
