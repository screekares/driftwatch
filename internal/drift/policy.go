package drift

import (
	"fmt"
	"strings"
)

// PolicyAction defines what happens when a policy rule matches.
type PolicyAction string

const (
	PolicyActionWarn PolicyAction = "warn"
	PolicyActionFail PolicyAction = "fail"
	PolicyActionIgnore PolicyAction = "ignore"
)

// PolicyRule represents a single rule that matches drift entries and applies an action.
type PolicyRule struct {
	ResourceType string       `json:"resource_type,omitempty"`
	MinSeverity  Severity     `json:"min_severity,omitempty"`
	Action       PolicyAction `json:"action"`
}

// Policy is a collection of rules evaluated against a drift report.
type Policy struct {
	Rules []PolicyRule `json:"rules"`
}

// PolicyResult holds the outcome of evaluating a policy against a report.
type PolicyResult struct {
	Passed   bool
	Warnings []string
	Failures []string
}

// Evaluate applies the policy rules to the given report and returns a PolicyResult.
func (p *Policy) Evaluate(r *Report) PolicyResult {
	result := PolicyResult{Passed: true}

	for _, entry := range r.Entries {
		for _, rule := range p.Rules {
			if rule.ResourceType != "" && !strings.EqualFold(rule.ResourceType, entry.ResourceType) {
				continue
			}
			if entry.Severity < rule.MinSeverity {
				continue
			}
			msg := fmt.Sprintf("resource %s/%s (severity: %s)", entry.ResourceType, entry.ResourceID, entry.Severity)
			switch rule.Action {
			case PolicyActionFail:
				result.Failures = append(result.Failures, msg)
				result.Passed = false
			case PolicyActionWarn:
				result.Warnings = append(result.Warnings, msg)
			}
		}
	}
	return result
}

// ParseAction converts a string to a PolicyAction.
func ParseAction(s string) (PolicyAction, error) {
	switch strings.ToLower(s) {
	case "warn":
		return PolicyActionWarn, nil
	case "fail":
		return PolicyActionFail, nil
	case "ignore":
		return PolicyActionIgnore, nil
	default:
		return "", fmt.Errorf("unknown policy action %q: must be warn, fail, or ignore", s)
	}
}
