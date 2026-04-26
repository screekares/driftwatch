package aws_test

import (
	"testing"

	// Blank import triggers init() registration
	_ "github.com/example/driftwatch/internal/provider/aws"
	"github.com/example/driftwatch/internal/provider"
)

func TestAWSProvider_RegisteredViaInit(t *testing.T) {
	names := provider.AvailableNames()
	for _, n := range names {
		if n == "aws" {
			return
		}
	}
	t.Error("expected 'aws' to be registered via init(), but it was not found")
}

func TestAWSProvider_NewViaRegistry(t *testing.T) {
	p, err := provider.New("aws", map[string]string{
		"region": "us-east-1",
	})
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if p == nil {
		t.Fatal("expected non-nil provider")
	}
}

func TestAWSProvider_NewViaRegistry_MissingConfig(t *testing.T) {
	_, err := provider.New("aws", map[string]string{})
	if err == nil {
		t.Fatal("expected error for missing region config, got nil")
	}
}
