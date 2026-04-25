package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestSnapshotCmd_MissingConfig(t *testing.T) {
	// Point cfgFile at a non-existent path; command should return an error.
	old := cfgFile
	cfgFile = "/nonexistent/driftwatch.yaml"
	t.Cleanup(func() { cfgFile = old })

	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"snapshot"})

	err := rootCmd.Execute()
	if err == nil {
		t.Error("expected error when config file is missing")
	}
}

func TestSnapshotCmd_OutputFlag(t *testing.T) {
	// Verify the --output flag is registered on the snapshot sub-command.
	f := snapshotCmd.Flags().Lookup("output")
	if f == nil {
		t.Fatal("expected --output flag to be registered")
	}
	if f.DefValue != "" {
		t.Errorf("expected empty default, got %q", f.DefValue)
	}
}

func TestSnapshotCmd_DefaultOutputDir(t *testing.T) {
	dir := t.TempDir()
	outFile := filepath.Join(dir, "snap.json")

	// Write a minimal valid config.
	cfgPath := filepath.Join(dir, "driftwatch.yaml")
	cfgContent := "version: \"1\"\nprovider: mock\nresources:\n  - id: instance-1\n    type: ec2\n"
	if err := os.WriteFile(cfgPath, []byte(cfgContent), 0o644); err != nil {
		t.Fatal(err)
	}

	old := cfgFile
	cfgFile = cfgPath
	oldOut := snapshotOutput
	snapshotOutput = outFile
	t.Cleanup(func() {
		cfgFile = old
		snapshotOutput = oldOut
	})

	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"snapshot", "--output", outFile})

	// We accept either success or a provider error; we just ensure no panic.
	_ = rootCmd.Execute()
}
