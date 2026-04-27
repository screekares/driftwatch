package drift

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// HistoryEntry records the result of a single drift check run.
type HistoryEntry struct {
	Timestamp   time.Time `json:"timestamp"`
	Provider    string    `json:"provider"`
	TotalChecked int      `json:"total_checked"`
	DriftedCount int      `json:"drifted_count"`
	MissingCount int      `json:"missing_count"`
	Summary     string    `json:"summary"`
}

// History holds a list of past drift check entries.
type History struct {
	Entries []HistoryEntry `json:"entries"`
}

// AppendEntry adds a new entry to the history file at the given path.
// If the file does not exist it is created.
func AppendEntry(path string, entry HistoryEntry) error {
	h, err := LoadHistory(path)
	if err != nil {
		return fmt.Errorf("load history: %w", err)
	}
	h.Entries = append(h.Entries, entry)
	return saveHistory(path, h)
}

// LoadHistory reads a history file from disk. Returns an empty History if the
// file does not exist.
func LoadHistory(path string) (*History, error) {
	h := &History{}
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return h, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read history file: %w", err)
	}
	if err := json.Unmarshal(data, h); err != nil {
		return nil, fmt.Errorf("parse history file: %w", err)
	}
	return h, nil
}

func saveHistory(path string, h *History) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("create history dir: %w", err)
	}
	data, err := json.MarshalIndent(h, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal history: %w", err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("write history file: %w", err)
	}
	return nil
}

// EntryFromReport builds a HistoryEntry from a Report and provider name.
func EntryFromReport(r *Report, provider string) HistoryEntry {
	drifted := 0
	missing := 0
	for _, res := range r.Resources {
		switch res.Status {
		case StatusDrifted:
			drifted++
		case StatusMissing:
			missing++
		}
	}
	return HistoryEntry{
		Timestamp:    time.Now().UTC(),
		Provider:     provider,
		TotalChecked: len(r.Resources),
		DriftedCount: drifted,
		MissingCount: missing,
		Summary:      r.Summary(),
	}
}
