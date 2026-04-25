package drift

import (
	"context"
	"fmt"

	"github.com/driftwatch/driftwatch/internal/provider"
)

// Status represents the drift status of a single resource.
type Status int

const (
	StatusMatch  Status = iota // declared == live
	StatusDrift               // live differs from declared
	StatusMissing             // resource not found in provider
)

// Result holds the drift analysis for one resource.
type Result struct {
	ResourceID string
	Status     Status
	Declared   map[string]string
	Live       map[string]string
	Diffs      []FieldDiff
}

// FieldDiff describes a single field that differs.
type FieldDiff struct {
	Field    string
	Declared string
	Live     string
}

// Detector compares declared resources against live provider state.
type Detector struct {
	provider provider.Provider
}

// New returns a Detector backed by the given provider.
func New(p provider.Provider) *Detector {
	return &Detector{provider: p}
}

// Check fetches the live state of resourceID and compares it to declared.
func (d *Detector) Check(ctx context.Context, resourceID string, declared map[string]string) (Result, error) {
	live, err := d.provider.FetchResource(ctx, resourceID)
	if err != nil {
		return Result{ResourceID: resourceID, Status: StatusMissing}, fmt.Errorf("fetch %q: %w", resourceID, err)
	}

	diffs := compare(declared, live)
	status := StatusMatch
	if len(diffs) > 0 {
		status = StatusDrift
	}

	return Result{
		ResourceID: resourceID,
		Status:     status,
		Declared:   declared,
		Live:       live,
		Diffs:      diffs,
	}, nil
}

// compare returns field-level differences between declared and live maps.
func compare(declared, live map[string]string) []FieldDiff {
	seen := make(map[string]bool)
	var diffs []FieldDiff

	for k, dv := range declared {
		seen[k] = true
		if lv, ok := live[k]; !ok || lv != dv {
			diffs = append(diffs, FieldDiff{Field: k, Declared: dv, Live: lv})
		}
	}
	for k, lv := range live {
		if !seen[k] {
			diffs = append(diffs, FieldDiff{Field: k, Declared: "", Live: lv})
		}
	}
	return diffs
}
