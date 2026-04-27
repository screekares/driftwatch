package drift

import (
	"fmt"
	"strings"
)

// RemediationAction represents a suggested fix for a drifted resource.
type RemediationAction struct {
	ResourceID   string
	ResourceType string
	Field        string
	Expected     interface{}
	Actual       interface{}
	Suggestion   string
	Severity     Severity
}

// RemediationPlan holds all suggested actions for a given report.
type RemediationPlan struct {
	Actions []RemediationAction
}

// HasActions returns true when there is at least one remediation action.
func (p *RemediationPlan) HasActions() bool {
	return len(p.Actions) > 0
}

// Summary returns a human-readable overview of the plan.
func (p *RemediationPlan) Summary() string {
	if !p.HasActions() {
		return "No remediation actions required."
	}
	return fmt.Sprintf("%d remediation action(s) suggested.", len(p.Actions))
}

// BuildRemediationPlan generates a RemediationPlan from a Report.
func BuildRemediationPlan(r Report) RemediationPlan {
	var actions []RemediationAction

	for _, entry := range r.Entries {
		if entry.Status == StatusMatch {
			continue
		}

		if entry.Status == StatusMissing {
			actions = append(actions, RemediationAction{
				ResourceID:   entry.ResourceID,
				ResourceType: entry.ResourceType,
				Suggestion:   fmt.Sprintf("Resource %q of type %q is missing from live state; consider re-provisioning.", entry.ResourceID, entry.ResourceType),
				Severity:     SeverityHigh,
			})
			continue
		}

		for _, diff := range entry.Differences {
			actions = append(actions, RemediationAction{
				ResourceID:   entry.ResourceID,
				ResourceType: entry.ResourceType,
				Field:        diff.Field,
				Expected:     diff.Expected,
				Actual:       diff.Actual,
				Suggestion:   buildSuggestion(entry.ResourceType, diff.Field, diff.Expected),
				Severity:     ClassifyDrift(entry),
			})
		}
	}

	return RemediationPlan{Actions: actions}
}

func buildSuggestion(resourceType, field string, expected interface{}) string {
	return fmt.Sprintf(
		"Update field %q on %s to match declared value %q.",
		field,
		strings.ToLower(resourceType),
		fmt.Sprintf("%v", expected),
	)
}
