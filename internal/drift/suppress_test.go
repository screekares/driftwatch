package drift

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

var (
	now    = time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)
	future = now.Add(24 * time.Hour)
	past   = now.Add(-24 * time.Hour)
)

func suppressionReport() Report {
	return Report{
		Entries: []Entry{
			{ResourceType: "ec2_instance", ResourceID: "i-aaa", Status: StatusDrifted},
			{ResourceType: "s3_bucket", ResourceID: "bucket-1", Status: StatusDrifted},
			{ResourceType: "ec2_instance", ResourceID: "i-bbb", Status: StatusMissing},
		},
	}
}

func TestIsActive_NotExpired(t *testing.T) {
	s := Suppression{ExpiresAt: future}
	if !s.IsActive(now) {
		t.Error("expected suppression to be active")
	}
}

func TestIsActive_Expired(t *testing.T) {
	s := Suppression{ExpiresAt: past}
	if s.IsActive(now) {
		t.Error("expected suppression to be inactive")
	}
}

func TestApplySuppressions_NoActive(t *testing.T) {
	sl := SuppressionList{
		Suppressions: []Suppression{
			{ResourceType: "ec2_instance", ResourceID: "i-aaa", ExpiresAt: past},
		},
	}
	r := ApplySuppressions(suppressionReport(), sl, now)
	if len(r.Entries) != 3 {
		t.Errorf("expected 3 entries, got %d", len(r.Entries))
	}
}

func TestApplySuppressions_SpecificResource(t *testing.T) {
	sl := SuppressionList{
		Suppressions: []Suppression{
			{ResourceType: "ec2_instance", ResourceID: "i-aaa", ExpiresAt: future},
		},
	}
	r := ApplySuppressions(suppressionReport(), sl, now)
	if len(r.Entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(r.Entries))
	}
}

func TestApplySuppressions_Wildcard(t *testing.T) {
	sl := SuppressionList{
		Suppressions: []Suppression{
			{ResourceType: "ec2_instance", ResourceID: "*", ExpiresAt: future},
		},
	}
	r := ApplySuppressions(suppressionReport(), sl, now)
	if len(r.Entries) != 1 {
		t.Errorf("expected 1 entry, got %d", len(r.Entries))
	}
	if r.Entries[0].ResourceType != "s3_bucket" {
		t.Errorf("unexpected resource type: %s", r.Entries[0].ResourceType)
	}
}

func TestSaveAndLoadSuppressions_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "suppressions.json")

	sl := SuppressionList{
		Suppressions: []Suppression{
			{ResourceType: "s3_bucket", ResourceID: "*", Reason: "maintenance", ExpiresAt: future, CreatedAt: now},
		},
	}
	if err := SaveSuppressions(path, sl); err != nil {
		t.Fatalf("save: %v", err)
	}
	loaded, err := LoadSuppressions(path)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if len(loaded.Suppressions) != 1 {
		t.Errorf("expected 1 suppression, got %d", len(loaded.Suppressions))
	}
}

func TestLoadSuppressions_FileNotFound(t *testing.T) {
	sl, err := LoadSuppressions("/nonexistent/path/suppressions.json")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(sl.Suppressions) != 0 {
		t.Errorf("expected empty list, got %d", len(sl.Suppressions))
	}
}

func TestLoadSuppressions_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	_ = os.WriteFile(path, []byte("not-json"), 0o644)
	_, err := LoadSuppressions(path)
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestSaveSuppressions_WritesValidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "s.json")
	sl := SuppressionList{Suppressions: []Suppression{}}
	if err := SaveSuppressions(path, sl); err != nil {
		t.Fatalf("save: %v", err)
	}
	data, _ := os.ReadFile(path)
	var out SuppressionList
	if err := json.Unmarshal(data, &out); err != nil {
		t.Errorf("invalid JSON written: %v", err)
	}
}
