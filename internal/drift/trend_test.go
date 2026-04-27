package drift

import (
	"testing"
	"time"
)

func makeEntries(counts [][2]int) []HistoryEntry {
	base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	entries := make([]HistoryEntry, len(counts))
	for i, c := range counts {
		entries[i] = HistoryEntry{
			Timestamp:    base.Add(time.Duration(i) * time.Hour),
			DriftedCount: c[0],
			MissingCount: c[1],
		}
	}
	return entries
}

func TestAnalyzeTrend_Increasing(t *testing.T) {
	entries := makeEntries([][2]int{{1, 0}, {2, 0}, {3, 1}})
	ts, err := AnalyzeTrend(entries, 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ts.Direction != TrendIncreasing {
		t.Errorf("expected increasing, got %s", ts.Direction)
	}
	if ts.Delta != 3 {
		t.Errorf("expected delta 3, got %d", ts.Delta)
	}
}

func TestAnalyzeTrend_Decreasing(t *testing.T) {
	entries := makeEntries([][2]int{{5, 2}, {3, 1}, {1, 0}})
	ts, err := AnalyzeTrend(entries, 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ts.Direction != TrendDecreasing {
		t.Errorf("expected decreasing, got %s", ts.Direction)
	}
	if ts.Delta != -6 {
		t.Errorf("expected delta -6, got %d", ts.Delta)
	}
}

func TestAnalyzeTrend_Stable(t *testing.T) {
	entries := makeEntries([][2]int{{2, 1}, {1, 2}, {2, 1}})
	ts, err := AnalyzeTrend(entries, 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ts.Direction != TrendStable {
		t.Errorf("expected stable, got %s", ts.Direction)
	}
}

func TestAnalyzeTrend_WindowLargerThanEntries(t *testing.T) {
	entries := makeEntries([][2]int{{1, 0}, {2, 0}})
	ts, err := AnalyzeTrend(entries, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ts.WindowSize != 2 {
		t.Errorf("expected window size 2, got %d", ts.WindowSize)
	}
}

func TestAnalyzeTrend_EmptyEntries(t *testing.T) {
	_, err := AnalyzeTrend([]HistoryEntry{}, 3)
	if err == nil {
		t.Error("expected error for empty entries")
	}
}

func TestAnalyzeTrend_ZeroWindow(t *testing.T) {
	entries := makeEntries([][2]int{{1, 0}})
	_, err := AnalyzeTrend(entries, 0)
	if err == nil {
		t.Error("expected error for zero window")
	}
}

func TestTrendSummary_String(t *testing.T) {
	ts := TrendSummary{
		Direction:  TrendStable,
		WindowSize: 5,
		AvgDrifted: 2.4,
		AvgMissing: 0.6,
		Delta:      0,
	}
	s := ts.String()
	if s == "" {
		t.Error("expected non-empty string")
	}
}
