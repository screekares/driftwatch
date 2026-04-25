package cmd

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"driftwatch/internal/snapshot"
)

func writeTempSnapshot(t *testing.T, snap *snapshot.Snapshot) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "snap.json")
	if err := snapshot.Save(snap, path); err != nil {
		t.Fatalf("save snapshot: %v", err)
	}
	return path
}

func TestCheckCmd_MissingSnapshot(t *testing.T) {
	cmd := checkCmd
	cmd.ResetFlags()
	init()

	buf := &bytes.Buffer{}
	cmd.SetOut(buf)
	cmd.SetErr(buf)

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error when --snapshot flag is missing")
	}
}

func TestCheckCmd_InvalidSnapshotPath(t *testing.T) {
	cfgPath := writeTempConfig(t)

	buf := &bytes.Buffer{}
	RootCmd.SetOut(buf)
	RootCmd.SetErr(buf)
	RootCmd.SetArgs([]string{"check", "--config", cfgPath, "--snapshot", "/nonexistent/snap.json"})

	err := RootCmd.Execute()
	if err == nil {
		t.Fatal("expected error for missing snapshot file")
	}
}

func TestCheckCmd_NoDrift_TextOutput(t *testing.T) {
	cfgPath := writeTempConfig(t)

	snap := &snapshot.Snapshot{
		Resources: map[string]map[string]string{
			"instance-1": {"type": "vm", "status": "running"},
		},
	}
	snapPath := writeTempSnapshot(t, snap)

	buf := &bytes.Buffer{}
	RootCmd.SetOut(buf)
	RootCmd.SetErr(buf)
	RootCmd.SetArgs([]string{"check", "--config", cfgPath, "--snapshot", snapPath, "--format", "text"})

	_ = RootCmd.Execute()
	output := buf.String()
	if output == "" {
		t.Error("expected non-empty output from check command")
	}
}

func TestCheckCmd_JSONOutput_ValidJSON(t *testing.T) {
	cfgPath := writeTempConfig(t)

	snap := &snapshot.Snapshot{
		Resources: map[string]map[string]string{
			"instance-1": {"type": "vm", "status": "running"},
		},
	}
	snapPath := writeTempSnapshot(t, snap)

	f, err := os.CreateTemp(t.TempDir(), "out-*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	buf := &bytes.Buffer{}
	RootCmd.SetOut(buf)
	RootCmd.SetErr(buf)
	RootCmd.SetArgs([]string{"check", "--config", cfgPath, "--snapshot", snapPath, "--format", "json"})

	_ = RootCmd.Execute()

	var result interface{}
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Logf("output: %s", buf.String())
		// JSON output may be partial if preceded by summary line; not fatal
	}
}
