package drift

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Baseline represents a saved drift report used as a reference point
// for suppressing known/accepted drift in future checks.
type Baseline struct {
	CreatedAt  time.Time         `json:"created_at"`
	Provider   string            `json:"provider"`
	Ignored    map[string]bool   `json:"ignored"`    // key: "resourceType/resourceID"
	Annotation string            `json:"annotation,omitempty"`
}

// NewBaseline creates a Baseline from a Report, marking all drifted
// and missing resources as ignored.
func NewBaseline(r Report, annotation string) Baseline {
	ignored := make(map[string]bool, len(r.Entries))
	for _, e := range r.Entries {
		if e.Status == StatusDrifted || e.Status == StatusMissing {
			key := baselineKey(e.ResourceType, e.ResourceID)
			ignored[key] = true
		}
	}
	return Baseline{
		CreatedAt:  time.Now().UTC(),
		Provider:   r.Provider,
		Ignored:    ignored,
		Annotation: annotation,
	}
}

// SaveBaseline writes a Baseline to the given file path as JSON.
func SaveBaseline(path string, b Baseline) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("baseline: create %q: %w", path, err)
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(b); err != nil {
		return fmt.Errorf("baseline: encode: %w", err)
	}
	return nil
}

// LoadBaseline reads a Baseline from the given file path.
func LoadBaseline(path string) (Baseline, error) {
	f, err := os.Open(path)
	if err != nil {
		return Baseline{}, fmt.Errorf("baseline: open %q: %w", path, err)
	}
	defer f.Close()
	var b Baseline
	if err := json.NewDecoder(f).Decode(&b); err != nil {
		return Baseline{}, fmt.Errorf("baseline: decode: %w", err)
	}
	return b, nil
}

// ApplyBaseline removes entries from the Report that are present in the
// Baseline's ignored set, returning the filtered Report.
func ApplyBaseline(r Report, b Baseline) Report {
	filtered := make([]ReportEntry, 0, len(r.Entries))
	for _, e := range r.Entries {
		key := baselineKey(e.ResourceType, e.ResourceID)
		if !b.Ignored[key] {
			filtered = append(filtered, e)
		}
	}
	r.Entries = filtered
	return r
}

func baselineKey(resourceType, resourceID string) string {
	return resourceType + "/" + resourceID
}
