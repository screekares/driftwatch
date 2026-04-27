package drift

// Severity represents the importance level of a detected drift.
type Severity int

const (
	// SeverityLow indicates a minor configuration difference (e.g. tag changes).
	SeverityLow Severity = iota
	// SeverityMedium indicates a notable configuration difference (e.g. size changes).
	SeverityMedium
	// SeverityHigh indicates a critical configuration difference (e.g. missing resource).
	SeverityHigh
)

// String returns a human-readable label for the severity level.
func (s Severity) String() string {
	switch s {
	case SeverityLow:
		return "low"
	case SeverityMedium:
		return "medium"
	case SeverityHigh:
		return "high"
	default:
		return "unknown"
	}
}

// ParseSeverity converts a string label into a Severity value.
// Returns SeverityLow and false if the label is not recognised.
func ParseSeverity(s string) (Severity, bool) {
	switch s {
	case "low":
		return SeverityLow, true
	case "medium":
		return SeverityMedium, true
	case "high":
		return SeverityHigh, true
	default:
		return SeverityLow, false
	}
}

// ClassifyDrift assigns a Severity to a DriftResult based on simple heuristics:
//   - Missing resources are High.
//   - Resources with more than two differing fields are Medium.
//   - All other drifted resources are Low.
func ClassifyDrift(result DriftResult) Severity {
	if result.Missing {
		return SeverityHigh
	}
	if len(result.Differences) > 2 {
		return SeverityMedium
	}
	return SeverityLow
}
