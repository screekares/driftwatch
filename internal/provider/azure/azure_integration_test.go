//go:build integration
// +build integration

// Package azure contains integration tests that exercise the provider
// against a real Azure subscription. Run with:
//
//	AZURE_SUBSCRIPTION_ID=<id> AZURE_RESOURCE_GROUP=<rg> \
//	  go test -tags integration ./internal/provider/azure/...
package azure

import (
	"os"
	"testing"
)

func TestIntegration_FetchVirtualMachine(t *testing.T) {
	sub := os.Getenv("AZURE_SUBSCRIPTION_ID")
	rg := os.Getenv("AZURE_RESOURCE_GROUP")
	if sub == "" || rg == "" {
		t.Skip("AZURE_SUBSCRIPTION_ID or AZURE_RESOURCE_GROUP not set")
	}

	p, err := New(map[string]string{
		"subscription_id": sub,
		"resource_group":  rg,
	})
	if err != nil {
		t.Fatalf("failed to create provider: %v", err)
	}

	// Assumes at least one VM named "test-vm" exists in the resource group.
	res, err := p.FetchResource("VirtualMachine", "test-vm")
	if err != nil {
		t.Fatalf("FetchResource error: %v", err)
	}
	if res["id"] != "test-vm" {
		t.Errorf("unexpected id: %q", res["id"])
	}
}
