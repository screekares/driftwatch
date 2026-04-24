package provider

import "fmt"

// Provider defines the interface for fetching live infrastructure state.
type Provider interface {
	// Name returns the provider identifier (e.g. "aws", "gcp").
	Name() string
	// FetchResource retrieves the live state of a named resource.
	FetchResource(resourceType, resourceID string) (map[string]interface{}, error)
}

// Registry holds registered provider factories.
var registry = map[string]func(cfg map[string]string) (Provider, error){}

// Register adds a provider factory under the given name.
func Register(name string, factory func(cfg map[string]string) (Provider, error)) {
	registry[name] = factory
}

// New instantiates a provider by name using the supplied configuration.
func New(name string, cfg map[string]string) (Provider, error) {
	factory, ok := registry[name]
	if !ok {
		return nil, fmt.Errorf("unknown provider %q; available: %v", name, availableNames())
	}
	return factory(cfg)
}

// availableNames returns a slice of registered provider names.
func availableNames() []string {
	names := make([]string, 0, len(registry))
	for k := range registry {
		names = append(names, k)
	}
	return names
}
