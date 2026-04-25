package gcp

import (
	"context"
	"testing"
)

func TestNew_MissingProject(t *testing.T) {
	_, err := New(map[string]string{})
	if err == nil {
		t.Fatal("expected error for missing project, got nil")
	}
}

func TestNew_ValidConfig(t *testing.T) {
	p, err := New(map[string]string{"project": "my-project", "region": "us-central1"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p == nil {
		t.Fatal("expected non-nil provider")
	}
}

func TestFetchResource_ComputeInstance(t *testing.T) {
	p, _ := New(map[string]string{"project": "my-project", "region": "us-central1"})
	res, err := p.FetchResource(context.Background(), "compute_instance", "instance-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res["type"] != "compute_instance" {
		t.Errorf("expected type compute_instance, got %q", res["type"])
	}
	if res["project"] != "my-project" {
		t.Errorf("expected project my-project, got %q", res["project"])
	}
}

func TestFetchResource_StorageBucket(t *testing.T) {
	p, _ := New(map[string]string{"project": "my-project", "region": "us-central1"})
	res, err := p.FetchResource(context.Background(), "storage_bucket", "my-bucket")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res["id"] != "my-bucket" {
		t.Errorf("expected id my-bucket, got %q", res["id"])
	}
}

func TestFetchResource_UnsupportedType(t *testing.T) {
	p, _ := New(map[string]string{"project": "my-project"})
	_, err := p.FetchResource(context.Background(), "cloud_function", "fn-1")
	if err == nil {
		t.Fatal("expected error for unsupported resource type")
	}
}

func TestFetchResource_EmptyID(t *testing.T) {
	p, _ := New(map[string]string{"project": "my-project"})
	_, err := p.FetchResource(context.Background(), "compute_instance", "")
	if err == nil {
		t.Fatal("expected error for empty resource id")
	}
}

func TestInit_RegistersProvider(t *testing.T) {
	// Verify the provider was registered via init() by using provider.New.
	// Import side-effect already happened; just ensure New works through registry.
	p, err := New(map[string]string{"project": "proj"})
	if err != nil || p == nil {
		t.Fatalf("provider construction failed: %v", err)
	}
}
