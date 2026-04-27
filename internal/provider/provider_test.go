package provider_test

import (
	"errors"
	"testing"

	"github.com/example/driftwatch/internal/provider"
	_ "github.com/example/driftwatch/internal/provider/mock"
)

func TestNew_KnownProvider(t *testing.T) {
	p, err := provider.New("mock", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p == nil {
		t.Fatal("expected non-nil provider")
	}
}

func TestNew_UnknownProvider(t *testing.T) {
	_, err := provider.New("nonexistent", nil)
	if err == nil {
		t.Fatal("expected error for unknown provider")
	}
}

func TestMockProvider_FetchResource_Found(t *testing.T) {
	p, _ := provider.New("mock", nil)
	res, err := p.FetchResource("instance", "i-001")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res == nil {
		t.Fatal("expected non-nil resource")
	}
	if res.ID != "i-001" {
		t.Errorf("expected ID i-001, got %s", res.ID)
	}
}

func TestMockProvider_FetchResource_NotFound(t *testing.T) {
	p, _ := provider.New("mock", nil)
	_, err := p.FetchResource("instance", "i-999")
	if err == nil {
		t.Fatal("expected error for missing resource")
	}
}

// TestMockProvider_FetchResource_EmptyID verifies that an empty resource ID
// returns an error rather than a nil or zero-value resource.
func TestMockProvider_FetchResource_EmptyID(t *testing.T) {
	p, _ := provider.New("mock", nil)
	_, err := p.FetchResource("instance", "")
	if err == nil {
		t.Fatal("expected error for empty resource ID")
	}
}

func TestRegister_Duplicate(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic on duplicate registration")
		}
	}()
	factory := func(_ map[string]string) (provider.Provider, error) {
		return nil, errors.New("stub")
	}
	provider.Register("mock", factory)
}
