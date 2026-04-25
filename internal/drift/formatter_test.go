package drift

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestNewFormatter_Text(t *testing.T) {
	f, err := NewFormatter(FormatText)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := f.(*textFormatter); !ok {
		t.Fatalf("expected *textFormatter, got %T", f)
	}
}

func TestNewFormatter_JSON(t *testing.T) {
	f, err := NewFormatter(FormatJSON)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := f.(*jsonFormatter); !ok {
		t.Fatalf("expected *jsonFormatter, got %T", f)
	}
}

func TestNewFormatter_Unknown(t *testing.T) {
	_, err := NewFormatter("xml")
	if err == nil {
		t.Fatal("expected error for unknown format, got nil")
	}
}

func TestTextFormatter_ContainsHeaders(t *testing.T) {
	r := makeReport()
	f, _ := NewFormatter(FormatText)
	var buf bytes.Buffer
	if err := f.Format(r, &buf); err != nil {
		t.Fatalf("Format error: %v", err)
	}
	output := buf.String()
	for _, header := range []string{"RESOURCE", "FIELD", "EXPECTED", "ACTUAL", "STATUS"} {
		if !strings.Contains(output, header) {
			t.Errorf("output missing header %q", header)
		}
	}
}

func TestJSONFormatter_ValidJSON(t *testing.T) {
	r := makeReport()
	f, _ := NewFormatter(FormatJSON)
	var buf bytes.Buffer
	if err := f.Format(r, &buf); err != nil {
		t.Fatalf("Format error: %v", err)
	}
	var out map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
}

func TestTextFormatter_NoDrift_MissingStatus(t *testing.T) {
	r := &Report{
		Drifts: []DriftResult{
			{ResourceID: "svc-empty", Status: StatusMissing},
		},
	}
	f, _ := NewFormatter(FormatText)
	var buf bytes.Buffer
	if err := f.Format(r, &buf); err != nil {
		t.Fatalf("Format error: %v", err)
	}
	if !strings.Contains(buf.String(), "svc-empty") {
		t.Error("expected resource ID in output")
	}
}
