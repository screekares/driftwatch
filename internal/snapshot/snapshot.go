// Package snapshot provides functionality to capture and persist
// the current state of live resources for later drift comparison.
package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/driftwatch/internal/provider"
)

// ResourceSnapshot holds the captured state of a single resource.
type ResourceSnapshot struct {
	ID         string            `json:"id"`
	Type       string            `json:"type"`
	Attributes map[string]string `json:"attributes"`
	CapturedAt time.Time         `json:"captured_at"`
}

// Snapshot represents a full point-in-time capture of resources.
type Snapshot struct {
	Provider  string             `json:"provider"`
	CreatedAt time.Time          `json:"created_at"`
	Resources []ResourceSnapshot `json:"resources"`
}

// Capture fetches all declared resource IDs from the provider and
// records their current attributes into a Snapshot.
func Capture(p provider.Provider, resourceIDs []string) (*Snapshot, error) {
	snap := &Snapshot{
		CreatedAt: time.Now().UTC(),
		Resources: make([]ResourceSnapshot, 0, len(resourceIDs)),
	}

	for _, id := range resourceIDs {
		attrs, err := p.FetchResource(id)
		if err != nil {
			return nil, fmt.Errorf("snapshot: fetching resource %q: %w", id, err)
		}
		snap.Resources = append(snap.Resources, ResourceSnapshot{
			ID:         id,
			Attributes: attrs,
			CapturedAt: time.Now().UTC(),
		})
	}
	return snap, nil
}

// Save writes the snapshot as JSON to the given file path.
func Save(snap *Snapshot, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("snapshot: creating file %q: %w", path, err)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(snap); err != nil {
		return fmt.Errorf("snapshot: encoding snapshot: %w", err)
	}
	return nil
}

// Load reads a snapshot from the given JSON file path.
func Load(path string) (*Snapshot, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("snapshot: opening file %q: %w", path, err)
	}
	defer f.Close()

	var snap Snapshot
	if err := json.NewDecoder(f).Decode(&snap); err != nil {
		return nil, fmt.Errorf("snapshot: decoding snapshot: %w", err)
	}
	return &snap, nil
}
