package azure_test

import (
	"testing"

	"github.com/driftwatch/driftwatch/internal/provider"
	_ "github.com/driftwatch/driftwatch/internal/provider/azure" // trigger init
)

func TestAzureProvider_RegisteredViaInit(t *testing.T) {
	names := provider.AvailableNames()
	for _, n := range names {
		if n == "azure" {
			return
		}
	}
	t.Error("azure provider not found in registered providers")
}

func TestAzureProvider_NewViaRegistry(t *testing.T) {
	p, err := provider.New("azure", map[string]string{
		"subscription_id": "sub-reg-test",
		"resource_group":  "rg-reg-test",
	})
	if err != nil {
		t.Fatalf("provider.New(azure) error: %v", err)
	}
	if p.Name() != "azure" {
		t.Errorf("expected name azure, got %q", p.Name())
	}
}

func TestAzureProvider_NewViaRegistry_MissingConfig(t *testing.T) {
	_, err := provider.New("azure", map[string]string{})
	if err == nil {
		t.Fatal("expected error for missing config, got nil")
	}
}

func TestAzureProvider_NewViaRegistry_MissingResourceGroup(t *testing.T) {
	// Ensure that omitting resource_group alone also produces an error,
	// not just a fully empty config map.
	_, err := provider.New("azure", map[string]string{
		"subscription_id": "sub-reg-test",
	})
	if err == nil {
		t.Fatal("expected error when resource_group is missing, got nil")
	}
}
