package snapshot_test

import (
	"os"
	"testing"
	"time"

	"github.com/driftwatch/internal/snapshot"
)

// stubProvider satisfies provider.Provider for testing.
type stubProvider struct {
	resources map[string]map[string]string
}

func (s *stubProvider) FetchResource(id string) (map[string]string, error) {
	attrs, ok := s.resources[id]
	if !ok {
		return nil, nil
	}
	return attrs, nil
}

func newStub() *stubProvider {
	return &stubProvider{
		resources: map[string]map[string]string{
			"res-1": {"env": "prod", "size": "large"},
			"res-2": {"env": "staging", "size": "small"},
		},
	}
}

func TestCapture_AllResources(t *testing.T) {
	p := newStub()
	snap, err := snapshot.Capture(p, []string{"res-1", "res-2"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(snap.Resources) != 2 {
		t.Fatalf("expected 2 resources, got %d", len(snap.Resources))
	}
	if snap.Resources[0].ID != "res-1" {
		t.Errorf("expected res-1, got %s", snap.Resources[0].ID)
	}
}

func TestCapture_Empty(t *testing.T) {
	p := newStub()
	snap, err := snapshot.Capture(p, []string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(snap.Resources) != 0 {
		t.Errorf("expected 0 resources")
	}
}

func TestSaveAndLoad_RoundTrip(t *testing.T) {
	f, err := os.CreateTemp("", "snap-*.json")
	if err != nil {
		t.Fatal(err)
	}
	f.Close()
	defer os.Remove(f.Name())

	orig := &snapshot.Snapshot{
		Provider:  "mock",
		CreatedAt: time.Now().UTC().Truncate(time.Second),
		Resources: []snapshot.ResourceSnapshot{
			{ID: "r1", Attributes: map[string]string{"k": "v"}, CapturedAt: time.Now().UTC().Truncate(time.Second)},
		},
	}

	if err := snapshot.Save(orig, f.Name()); err != nil {
		t.Fatalf("save failed: %v", err)
	}

	loaded, err := snapshot.Load(f.Name())
	if err != nil {
		t.Fatalf("load failed: %v", err)
	}
	if loaded.Provider != orig.Provider {
		t.Errorf("provider mismatch: got %s", loaded.Provider)
	}
	if len(loaded.Resources) != 1 {
		t.Errorf("expected 1 resource, got %d", len(loaded.Resources))
	}
}

func TestLoad_FileNotFound(t *testing.T) {
	_, err := snapshot.Load("/nonexistent/path/snap.json")
	if err == nil {
		t.Error("expected error for missing file")
	}
}
