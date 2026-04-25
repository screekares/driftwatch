package drift_test

import (
	"strings"
	"testing"

	"github.com/driftwatch/driftwatch/internal/drift"
)

func makeReport(statuses ...drift.Status) *drift.Report {
	var results []drift.Result
	for i, s := range statuses {
		var diffs []drift.FieldDiff
		if s == drift.StatusDrift {
			diffs = []drift.FieldDiff{{Field: "image", Declared: "a", Live: "b"}}
		}
		results = append(results, drift.Result{
			ResourceID: fmt.Sprintf("res-%d", i),
			Status:     s,
			Diffs:      diffs,
		})
	}
	return &drift.Report{Results: results}
}

func TestReport_HasDrift_False(t *testing.T) {
	r := makeReport(drift.StatusMatch, drift.StatusMatch)
	if r.HasDrift() {
		t.Error("expected no drift")
	}
}

func TestReport_HasDrift_True(t *testing.T) {
	r := makeReport(drift.StatusMatch, drift.StatusDrift)
	if !r.HasDrift() {
		t.Error("expected drift")
	}
}

func TestReport_Summary(t *testing.T) {
	r := makeReport(drift.StatusMatch, drift.StatusDrift, drift.StatusMissing, drift.StatusMatch)
	m, d, miss := r.Summary()
	if m != 2 || d != 1 || miss != 1 {
		t.Errorf("unexpected summary: match=%d drift=%d missing=%d", m, d, miss)
	}
}

func TestReport_Write_ContainsLabels(t *testing.T) {
	r := makeReport(drift.StatusMatch, drift.StatusDrift, drift.StatusMissing)
	var sb strings.Builder
	r.Write(&sb)
	out := sb.String()
	for _, want := range []string{"[OK]", "[DRIFT]", "[MISSING]", "image"} {
		if !strings.Contains(out, want) {
			t.Errorf("output missing %q", want)
		}
	}
}
