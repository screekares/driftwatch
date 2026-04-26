package gcp_test

import (
	"testing"

	"github.com/driftwatch/driftwatch/internal/provider/gcp"
)

func TestNew_MissingProject(t *testing.T) {
	_, err := gcp.New(map[string]string{})
	if err == nil {
		t.Fatal("expected error for missing project, got nil")
	}
}

func TestNew_ValidConfig(t *testing.T) {
	p, err := gcp.New(map[string]string{"project": "my-project", "region": "us-central1"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p == nil {
		t.Fatal("expected non-nil provider")
	}
}

func TestFetchResource_ComputeInstance(t *testing.T) {
	p, _ := gcp.New(map[string]string{"project": "proj-1", "region": "us-east1"})
	res, err := p.FetchResource("compute_instance", "vm-abc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res["id"] != "vm-abc" {
		t.Errorf("expected id=vm-abc, got %s", res["id"])
	}
	if res["status"] != "RUNNING" {
		t.Errorf("expected status=RUNNING, got %s", res["status"])
	}
}

func TestFetchResource_StorageBucket(t *testing.T) {
	p, _ := gcp.New(map[string]string{"project": "proj-1", "region": "us-east1"})
	res, err := p.FetchResource("storage_bucket", "my-bucket")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res["storage_class"] != "STANDARD" {
		t.Errorf("expected storage_class=STANDARD, got %s", res["storage_class"])
	}
}

func TestFetchResource_UnsupportedType(t *testing.T) {
	p, _ := gcp.New(map[string]string{"project": "proj-1"})
	_, err := p.FetchResource("unknown_resource", "id-1")
	if err == nil {
		t.Fatal("expected error for unsupported resource type, got nil")
	}
}

func TestFetchResource_EmptyID(t *testing.T) {
	p, _ := gcp.New(map[string]string{"project": "proj-1"})
	_, err := p.FetchResource("compute_instance", "")
	if err == nil {
		t.Fatal("expected error for empty id, got nil")
	}
}
