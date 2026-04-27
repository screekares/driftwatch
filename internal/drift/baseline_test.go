package drift

import (
	"os"
	"path/filepath"
	"testing"
)

func baselineReport() Report {
	return Report{
		Provider: "mock",
		Entries: []ReportEntry{
			{ResourceType: "instance", ResourceID: "i-aaa", Status: StatusDrifted},
			{ResourceType: "bucket", ResourceID: "b-bbb", Status: StatusMissing},
			{ResourceType: "instance", ResourceID: "i-ccc", Status: StatusOK},
		},
	}
}

func TestNewBaseline_IgnoresDriftedAndMissing(t *testing.T) {
	r := baselineReport()
	b := NewBaseline(r, "initial baseline")

	if b.Provider != "mock" {
		t.Errorf("expected provider mock, got %s", b.Provider)
	}
	if b.Annotation != "initial baseline" {
		t.Errorf("unexpected annotation: %s", b.Annotation)
	}
	if !b.Ignored["instance/i-aaa"] {
		t.Error("expected instance/i-aaa to be ignored")
	}
	if !b.Ignored["bucket/b-bbb"] {
		t.Error("expected bucket/b-bbb to be ignored")
	}
	if b.Ignored["instance/i-ccc"] {
		t.Error("expected instance/i-ccc NOT to be ignored (status OK)")
	}
}

func TestSaveAndLoadBaseline_RoundTrip(t *testing.T) {
	r := baselineReport()
	b := NewBaseline(r, "round-trip test")

	dir := t.TempDir()
	path := filepath.Join(dir, "baseline.json")

	if err := SaveBaseline(path, b); err != nil {
		t.Fatalf("SaveBaseline: %v", err)
	}

	loaded, err := LoadBaseline(path)
	if err != nil {
		t.Fatalf("LoadBaseline: %v", err)
	}

	if loaded.Provider != b.Provider {
		t.Errorf("provider mismatch: got %s", loaded.Provider)
	}
	if loaded.Annotation != b.Annotation {
		t.Errorf("annotation mismatch: got %q, want %q", loaded.Annotation, b.Annotation)
	}
	if len(loaded.Ignored) != len(b.Ignored) {
		t.Errorf("ignored count mismatch: got %d, want %d", len(loaded.Ignored), len(b.Ignored))
	}
	// Verify individual ignored keys survived the round-trip.
	for key := range b.Ignored {
		if !loaded.Ignored[key] {
			t.Errorf("ignored key %q missing after round-trip", key)
		}
	}
}

func TestLoadBaseline_FileNotFound(t *testing.T) {
	_, err := LoadBaseline("/nonexistent/baseline.json")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestApplyBaseline_FiltersIgnoredEntries(t *testing.T) {
	r := baselineReport()
	b := NewBaseline(r, "")

	result := ApplyBaseline(r, b)

	if len(result.Entries) != 1 {
		t.Fatalf("expected 1 entry after baseline, got %d", len(result.Entries))
	}
	if result.Entries[0].ResourceID != "i-ccc" {
		t.Errorf("expected i-ccc to remain, got %s", result.Entries[0].ResourceID)
	}
}

func TestApplyBaseline_EmptyBaseline_NoFilter(t *testing.T) {
	r := baselineReport()
	b := Baseline{Ignored: map[string]bool{}}

	result := ApplyBaseline(r, b)
	if len(result.Entries) != len(r.Entries) {
		t.Errorf("expected all entries preserved, got %d", len(result.Entries))
	}
}

func TestSaveBaseline_InvalidPath(t *testing.T) {
	b := Baseline{Ignored: map[string]bool{}}
	err := SaveBaseline(filepath.Join(os.DevNull, "subdir", "baseline.json"), b)
	if err == nil {
		t.Fatal("expected error for invalid path")
	}
}
