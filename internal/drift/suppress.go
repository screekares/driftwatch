package drift

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Suppression represents a time-bounded rule that silences drift for a specific resource.
type Suppression struct {
	ResourceType string    `json:"resource_type"`
	ResourceID   string    `json:"resource_id"`
	Reason       string    `json:"reason"`
	ExpiresAt    time.Time `json:"expires_at"`
	CreatedAt    time.Time `json:"created_at"`
}

// SuppressionList holds all active suppressions.
type SuppressionList struct {
	Suppressions []Suppression `json:"suppressions"`
}

// IsActive reports whether the suppression has not yet expired.
func (s Suppression) IsActive(now time.Time) bool {
	return now.Before(s.ExpiresAt)
}

// LoadSuppressions reads a suppression file from disk.
// Returns an empty list if the file does not exist.
func LoadSuppressions(path string) (SuppressionList, error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return SuppressionList{}, nil
	}
	if err != nil {
		return SuppressionList{}, fmt.Errorf("read suppressions: %w", err)
	}
	var sl SuppressionList
	if err := json.Unmarshal(data, &sl); err != nil {
		return SuppressionList{}, fmt.Errorf("parse suppressions: %w", err)
	}
	return sl, nil
}

// SaveSuppressions writes the suppression list to disk.
func SaveSuppressions(path string, sl SuppressionList) error {
	data, err := json.MarshalIndent(sl, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal suppressions: %w", err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("write suppressions: %w", err)
	}
	return nil
}

// ApplySuppressions removes report entries that match an active suppression.
func ApplySuppressions(r Report, sl SuppressionList, now time.Time) Report {
	active := make([]Suppression, 0, len(sl.Suppressions))
	for _, s := range sl.Suppressions {
		if s.IsActive(now) {
			active = append(active, s)
		}
	}
	if len(active) == 0 {
		return r
	}

	filtered := make([]Entry, 0, len(r.Entries))
	for _, e := range r.Entries {
		suppressed := false
		for _, s := range active {
			if s.ResourceType == e.ResourceType &&
				(s.ResourceID == "*" || s.ResourceID == e.ResourceID) {
				suppressed = true
				break
			}
		}
		if !suppressed {
			filtered = append(filtered, e)
		}
	}
	r.Entries = filtered
	return r
}
