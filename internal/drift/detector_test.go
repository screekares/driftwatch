package drift_test

import (
	"context"
	"errors"
	"testing"

	"github.com/driftwatch/driftwatch/internal/drift"
	_ "github.com/driftwatch/driftwatch/internal/provider/mock"
	"github.com/driftwatch/driftwatch/internal/provider"
)

func newMockDetector(t *testing.T) *drift.Detector {
	t.Helper()
	p, err := provider.New("mock")
	if err != nil {
		t.Fatalf("provider.New: %v", err)
	}
	return drift.New(p)
}

func TestCheck_NoДrift(t *testing.T) {
	d := newMockDetector(t)
	declared := map[string]string{
		"image": "nginx:latest",
		"replicas": "3",
	}
	res, err := d.Check(context.Background(), "service-a", declared)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Status != drift.StatusMatch {
		t.Errorf("expected StatusMatch, got %v", res.Status)
	}
	if len(res.Diffs) != 0 {
		t.Errorf("expected no diffs, got %v", res.Diffs)
	}
}

func TestCheck_Drift(t *testing.T) {
	d := newMockDetector(t)
	declared := map[string]string{
		"image":    "nginx:1.19", // differs from mock live value
		"replicas": "3",
	}
	res, err := d.Check(context.Background(), "service-a", declared)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Status != drift.StatusDrift {
		t.Errorf("expected StatusDrift, got %v", res.Status)
	}
	if len(res.Diffs) == 0 {
		t.Error("expected diffs, got none")
	}
}

func TestCheck_Missing(t *testing.T) {
	d := newMockDetector(t)
	_, err := d.Check(context.Background(), "nonexistent", map[string]string{})
	if err == nil {
		t.Fatal("expected error for missing resource")
	}
	if !errors.Is(err, provider.ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}
