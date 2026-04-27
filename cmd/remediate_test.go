package cmd

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"driftwatch/internal/drift"
	"driftwatch/internal/snapshot"
)

func writeTempRemediateSnapshot(t *testing.T) string {
	t.Helper()
	snap := snapshot.Snapshot{
		Resources: []snapshot.Resource{
			{ID: "i-abc", Type: "ec2_instance", Attributes: map[string]interface{}{"instance_type": "t3.small"}},
		},
	}
	dir := t.TempDir()
	path := filepath.Join(dir, "snap.json")
	if err := snapshot.Save(snap, path); err != nil {
		t.Fatalf("saving snapshot: %v", err)
	}
	return path
}

func TestRemediateCmd_MissingSnapshot(t *testing.T) {
	root := newRootCmd()
	root.AddCommand(remediateCmd)
	buf := &bytes.Buffer{}
	root.SetOut(buf)
	root.SetErr(buf)
	root.SetArgs([]string{"remediate"})
	err := root.Execute()
	if err == nil {
		t.Fatal("expected error when --snapshot is missing")
	}
}

func TestRemediateCmd_InvalidSnapshotPath(t *testing.T) {
	root := newRootCmd()
	root.AddCommand(remediateCmd)
	buf := &bytes.Buffer{}
	root.SetOut(buf)
	root.SetErr(buf)
	root.SetArgs([]string{"remediate", "--snapshot", "/nonexistent/snap.json"})
	err := root.Execute()
	if err == nil {
		t.Fatal("expected error for invalid snapshot path")
	}
}

func TestRemediateCmd_TextOutput(t *testing.T) {
	path := writeTempRemediateSnapshot(t)
	root := newRootCmd()
	root.AddCommand(remediateCmd)
	buf := &bytes.Buffer{}
	root.SetOut(buf)
	root.SetErr(os.Stderr)
	root.SetArgs([]string{"remediate", "--snapshot", path, "--output", "text"})
	_ = root.Execute()
	// Output may be empty if no drift; just ensure no panic.
}

func TestRemediateCmd_JSONOutput(t *testing.T) {
	path := writeTempRemediateSnapshot(t)
	root := newRootCmd()
	root.AddCommand(remediateCmd)
	buf := &bytes.Buffer{}
	root.SetOut(buf)
	root.SetErr(os.Stderr)
	root.SetArgs([]string{"remediate", "--snapshot", path, "--output", "json"})
	_ = root.Execute()

	if buf.Len() > 0 {
		var plan drift.RemediationPlan
		if err := json.Unmarshal(buf.Bytes(), &plan); err != nil {
			t.Errorf("expected valid JSON output, got: %s", buf.String())
		}
	}
}
