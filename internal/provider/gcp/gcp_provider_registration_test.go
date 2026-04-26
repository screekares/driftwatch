package gcp_test

import (
	"testing"

	"github.com/driftwatch/driftwatch/internal/provider"
	_ "github.com/driftwatch/driftwatch/internal/provider/gcp"
)

func TestGCPProvider_RegisteredViaInit(t *testing.T) {
	names := provider.AvailableNames()
	for _, name := range names {
		if name == "gcp" {
			return
		}
	}
	t.Error("expected 'gcp' to be registered, but it was not found")
}

func TestGCPProvider_NewViaRegistry(t *testing.T) {
	p, err := provider.New("gcp", map[string]string{"project": "test-project", "region": "us-west1"})
	if err != nil {
		t.Fatalf("unexpected error creating gcp provider via registry: %v", err)
	}
	if p == nil {
		t.Fatal("expected non-nil provider from registry")
	}
}

func TestGCPProvider_NewViaRegistry_MissingConfig(t *testing.T) {
	_, err := provider.New("gcp", map[string]string{})
	if err == nil {
		t.Fatal("expected error when creating gcp provider without required config, got nil")
	}
}
