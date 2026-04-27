package drift

import (
	"strings"
	"testing"
)

func remediationReport() Report {
	return Report{
		Entries: []ReportEntry{
			{
				ResourceID:   "i-001",
				ResourceType: "ec2_instance",
				Status:       StatusMatch,
			},
			{
				ResourceID:   "i-002",
				ResourceType: "ec2_instance",
				Status:       StatusDrifted,
				Differences: []Difference{
					{Field: "instance_type", Expected: "t3.medium", Actual: "t3.small"},
				},
			},
			{
				ResourceID:   "i-003",
				ResourceType: "s3_bucket",
				Status:       StatusMissing,
			},
		},
	}
}

func TestBuildRemediationPlan_SkipsMatchedEntries(t *testing.T) {
	r := remediationReport()
	plan := BuildRemediationPlan(r)
	for _, a := range plan.Actions {
		if a.ResourceID == "i-001" {
			t.Error("matched resource should not appear in remediation plan")
		}
	}
}

func TestBuildRemediationPlan_MissingResource(t *testing.T) {
	r := remediationReport()
	plan := BuildRemediationPlan(r)
	var found bool
	for _, a := range plan.Actions {
		if a.ResourceID == "i-003" && a.ResourceType == "s3_bucket" {
			found = true
			if a.Severity != SeverityHigh {
				t.Errorf("expected SeverityHigh for missing resource, got %s", a.Severity)
			}
			if !strings.Contains(a.Suggestion, "re-provisioning") {
				t.Errorf("suggestion should mention re-provisioning, got: %s", a.Suggestion)
			}
		}
	}
	if !found {
		t.Error("missing resource i-003 not found in plan")
	}
}

func TestBuildRemediationPlan_DriftedField(t *testing.T) {
	r := remediationReport()
	plan := BuildRemediationPlan(r)
	var found bool
	for _, a := range plan.Actions {
		if a.ResourceID == "i-002" && a.Field == "instance_type" {
			found = true
			if !strings.Contains(a.Suggestion, "instance_type") {
				t.Errorf("suggestion should mention field name, got: %s", a.Suggestion)
			}
			if !strings.Contains(a.Suggestion, "t3.medium") {
				t.Errorf("suggestion should mention expected value, got: %s", a.Suggestion)
			}
		}
	}
	if !found {
		t.Error("drifted field action for i-002 not found in plan")
	}
}

func TestRemediationPlan_HasActions(t *testing.T) {
	empty := RemediationPlan{}
	if empty.HasActions() {
		t.Error("expected HasActions to be false for empty plan")
	}

	plan := BuildRemediationPlan(remediationReport())
	if !plan.HasActions() {
		t.Error("expected HasActions to be true")
	}
}

func TestRemediationPlan_Summary(t *testing.T) {
	empty := RemediationPlan{}
	if empty.Summary() != "No remediation actions required." {
		t.Errorf("unexpected empty summary: %s", empty.Summary())
	}

	plan := BuildRemediationPlan(remediationReport())
	if !strings.Contains(plan.Summary(), "remediation action") {
		t.Errorf("unexpected summary: %s", plan.Summary())
	}
}
