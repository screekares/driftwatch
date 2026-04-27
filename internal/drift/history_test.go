package drift

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestLoadHistory_FileNotFound(t *testing.T) {
	h, err := LoadHistory("/nonexistent/path/history.json")
	if err != nil {
		t.Fatalf("expected no error for missing file, got %v", err)
	}
	if len(h.Entries) != 0 {
		t.Errorf("expected empty entries, got %d", len(h.Entries))
	}
}

func TestAppendEntry_CreatesFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "history.json")

	entry := HistoryEntry{
		Timestamp:    time.Now().UTC(),
		Provider:     "aws",
		TotalChecked: 5,
		DriftedCount: 2,
		MissingCount: 1,
		Summary:      "3 resources checked",
	}

	if err := AppendEntry(path, entry); err != nil {
		t.Fatalf("AppendEntry failed: %v", err)
	}

	if _, err := os.Stat(path); err != nil {
		t.Fatalf("expected file to exist: %v", err)
	}
}

func TestAppendEntry_AccumulatesEntries(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "history.json")

	for i := 0; i < 3; i++ {
		entry := HistoryEntry{
			Timestamp: time.Now().UTC(),
			Provider:  "gcp",
			Summary:   "run",
		}
		if err := AppendEntry(path, entry); err != nil {
			t.Fatalf("AppendEntry iteration %d failed: %v", i, err)
		}
	}

	h, err := LoadHistory(path)
	if err != nil {
		t.Fatalf("LoadHistory failed: %v", err)
	}
	if len(h.Entries) != 3 {
		t.Errorf("expected 3 entries, got %d", len(h.Entries))
	}
}

func TestEntryFromReport_Counts(t *testing.T) {
	r := makeReport()
	entry := EntryFromReport(r, "aws")

	if entry.Provider != "aws" {
		t.Errorf("expected provider aws, got %s", entry.Provider)
	}
	if entry.TotalChecked != len(r.Resources) {
		t.Errorf("expected TotalChecked %d, got %d", len(r.Resources), entry.TotalChecked)
	}
	if entry.Timestamp.IsZero() {
		t.Error("expected non-zero timestamp")
	}
}

func TestLoadHistory_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	_ = os.WriteFile(path, []byte("not-json"), 0o644)

	_, err := LoadHistory(path)
	if err == nil {
		t.Error("expected error for invalid JSON, got nil")
	}
}
