package provider_test

import (
	"testing"

	"github.com/yourorg/driftwatch/internal/provider"
	_ "github.com/yourorg/driftwatch/internal/provider/mock" // register mock
)

func TestNew_KnownProvider(t *testing.T) {
	p, err := provider.New("mock", nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if p.Name() != "mock" {
		t.Errorf("expected name \"mock\", got %q", p.Name())
	}
}

func TestNew_UnknownProvider(t *testing.T) {
	_, err := provider.New("nonexistent", nil)
	if err == nil {
		t.Fatal("expected error for unknown provider, got nil")
	}
}

func TestMockProvider_FetchResource_Found(t *testing.T) {
	p, _ := provider.New("mock", nil)
	state, err := p.FetchResource("instance", "web-01")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if state["instance_type"] != "t3.micro" {
		t.Errorf("expected instance_type t3.micro, got %v", state["instance_type"])
	}
}

func TestMockProvider_FetchResource_NotFound(t *testing.T) {
	p, _ := provider.New("mock", nil)
	_, err := p.FetchResource("instance", "does-not-exist")
	if err == nil {
		t.Fatal("expected error for missing resource, got nil")
	}
}

func TestRegister_Duplicate(t *testing.T) {
	// Re-registering should silently overwrite — no panic expected.
	provider.Register("mock", func(cfg map[string]string) (provider.Provider, error) {
		return nil, nil
	})
}
