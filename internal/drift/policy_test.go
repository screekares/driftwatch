package drift

import (
	"testing"
)

func policyReport() *Report {
	return &Report{
		Entries: []Entry{
			{ResourceID: "i-001", ResourceType: "ec2_instance", Status: StatusDrifted, Severity: SeverityHigh},
			{ResourceID: "bucket-a", ResourceType: "s3_bucket", Status: StatusDrifted, Severity: SeverityLow},
			{ResourceID: "i-002", ResourceType: "ec2_instance", Status: StatusMissing, Severity: SeverityCritical},
		},
	}
}

func TestPolicy_NoRules(t *testing.T) {
	p := &Policy{}
	res := p.Evaluate(policyReport())
	if !res.Passed {
		t.Error("expected policy with no rules to pass")
	}
	if len(res.Failures) != 0 || len(res.Warnings) != 0 {
		t.Error("expected no failures or warnings")
	}
}

func TestPolicy_FailOnHighSeverity(t *testing.T) {
	p := &Policy{
		Rules: []PolicyRule{
			{MinSeverity: SeverityHigh, Action: PolicyActionFail},
		},
	}
	res := p.Evaluate(policyReport())
	if res.Passed {
		t.Error("expected policy to fail")
	}
	if len(res.Failures) != 2 {
		t.Errorf("expected 2 failures (high + critical), got %d", len(res.Failures))
	}
}

func TestPolicy_WarnOnResourceType(t *testing.T) {
	p := &Policy{
		Rules: []PolicyRule{
			{ResourceType: "s3_bucket", MinSeverity: SeverityLow, Action: PolicyActionWarn},
		},
	}
	res := p.Evaluate(policyReport())
	if !res.Passed {
		t.Error("warn-only rule should not fail the policy")
	}
	if len(res.Warnings) != 1 {
		t.Errorf("expected 1 warning, got %d", len(res.Warnings))
	}
}

func TestPolicy_IgnoreAction(t *testing.T) {
	p := &Policy{
		Rules: []PolicyRule{
			{ResourceType: "ec2_instance", MinSeverity: SeverityLow, Action: PolicyActionIgnore},
		},
	}
	res := p.Evaluate(policyReport())
	if !res.Passed {
		t.Error("ignore action should not fail the policy")
	}
	if len(res.Warnings) != 0 || len(res.Failures) != 0 {
		t.Error("ignore action should produce no warnings or failures")
	}
}

func TestParseAction_Valid(t *testing.T) {
	cases := []struct {
		input    string
		expected PolicyAction
	}{
		{"warn", PolicyActionWarn},
		{"fail", PolicyActionFail},
		{"ignore", PolicyActionIgnore},
		{"WARN", PolicyActionWarn},
	}
	for _, tc := range cases {
		got, err := ParseAction(tc.input)
		if err != nil {
			t.Errorf("unexpected error for %q: %v", tc.input, err)
		}
		if got != tc.expected {
			t.Errorf("ParseAction(%q) = %q, want %q", tc.input, got, tc.expected)
		}
	}
}

func TestParseAction_Invalid(t *testing.T) {
	_, err := ParseAction("block")
	if err == nil {
		t.Error("expected error for unknown action")
	}
}
