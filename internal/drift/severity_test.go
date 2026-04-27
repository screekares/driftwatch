package drift

import (
	"testing"
)

func TestSeverity_String(t *testing.T) {
	cases := []struct {
		sev  Severity
		want string
	}{
		{SeverityLow, "low"},
		{SeverityMedium, "medium"},
		{SeverityHigh, "high"},
		{Severity(99), "unknown"},
	}
	for _, tc := range cases {
		if got := tc.sev.String(); got != tc.want {
			t.Errorf("Severity(%d).String() = %q, want %q", tc.sev, got, tc.want)
		}
	}
}

func TestParseSeverity_Known(t *testing.T) {
	cases := []struct {
		input string
		want  Severity
	}{
		{"low", SeverityLow},
		{"medium", SeverityMedium},
		{"high", SeverityHigh},
	}
	for _, tc := range cases {
		got, ok := ParseSeverity(tc.input)
		if !ok {
			t.Errorf("ParseSeverity(%q) returned ok=false", tc.input)
		}
		if got != tc.want {
			t.Errorf("ParseSeverity(%q) = %v, want %v", tc.input, got, tc.want)
		}
	}
}

func TestParseSeverity_Unknown(t *testing.T) {
	got, ok := ParseSeverity("critical")
	if ok {
		t.Errorf("ParseSeverity(\"critical\") expected ok=false, got true")
	}
	if got != SeverityLow {
		t.Errorf("ParseSeverity(\"critical\") default = %v, want SeverityLow", got)
	}
}

func TestClassifyDrift_Missing(t *testing.T) {
	result := DriftResult{ResourceID: "res-1", Missing: true}
	if got := ClassifyDrift(result); got != SeverityHigh {
		t.Errorf("ClassifyDrift(missing) = %v, want high", got)
	}
}

func TestClassifyDrift_ManyDifferences(t *testing.T) {
	result := DriftResult{
		ResourceID:  "res-2",
		Differences: []string{"field1", "field2", "field3"},
	}
	if got := ClassifyDrift(result); got != SeverityMedium {
		t.Errorf("ClassifyDrift(3 diffs) = %v, want medium", got)
	}
}

func TestClassifyDrift_FewDifferences(t *testing.T) {
	result := DriftResult{
		ResourceID:  "res-3",
		Differences: []string{"tag"},
	}
	if got := ClassifyDrift(result); got != SeverityLow {
		t.Errorf("ClassifyDrift(1 diff) = %v, want low", got)
	}
}
