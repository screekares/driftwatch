// Package azure provides a driftwatch provider for Microsoft Azure resources.
package azure

import (
	"fmt"

	"github.com/driftwatch/driftwatch/internal/provider"
)

const providerName = "azure"
unc init() {
	provider.Register(providerName, New)
}

// azureProvider implements provider.Provider for Azure.
type azureProvider struct {
	subscriptionID string
	resourceGroup  string
}

// New constructs an azureProvider from the given config map.
// Required keys: subscription_id, resource_group.
func New(cfg map[string]string) (provider.Provider, error) {
	sub, ok := cfg["subscription_id"]
	if !ok || sub == "" {
		return nil, fmt.Errorf("azure provider: missing required config key 'subscription_id'")
	}
	rg, ok := cfg["resource_group"]
	if !ok || rg == "" {
		return nil, fmt.Errorf("azure provider: missing required config key 'resource_group'")
	}
	return &azureProvider{subscriptionID: sub, resourceGroup: rg}, nil
}

// Name returns the provider identifier.
func (p *azureProvider) Name() string { return providerName }

// FetchResource retrieves a stub resource by type and id.
// Supported types: VirtualMachine, StorageAccount.
func (p *azureProvider) FetchResource(resourceType, id string) (map[string]string, error) {
	switch resourceType {
	case "VirtualMachine":
		return map[string]string{
			"id":             id,
			"type":           "VirtualMachine",
			"subscription":   p.subscriptionID,
			"resource_group": p.resourceGroup,
			"size":           "Standard_D2s_v3",
			"location":       "eastus",
		}, nil
	case "StorageAccount":
		return map[string]string{
			"id":             id,
			"type":           "StorageAccount",
			"subscription":   p.subscriptionID,
			"resource_group": p.resourceGroup,
			"sku":            "Standard_LRS",
			"location":       "eastus",
		}, nil
	default:
		return nil, fmt.Errorf("azure provider: unsupported resource type %q", resourceType)
	}
}
