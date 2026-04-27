package drift

import (
	"fmt"
	"time"
)

// TrendDirection indicates whether drift is increasing, decreasing, or stable.
type TrendDirection string

const (
	TrendIncreasing TrendDirection = "increasing"
	TrendDecreasing TrendDirection = "decreasing"
	TrendStable     TrendDirection = "stable"
)

// TrendSummary holds a computed trend over a window of history entries.
type TrendSummary struct {
	Direction    TrendDirection `json:"direction"`
	WindowSize   int            `json:"window_size"`
	FirstChecked time.Time      `json:"first_checked"`
	LastChecked  time.Time      `json:"last_checked"`
	AvgDrifted   float64        `json:"avg_drifted"`
	AvgMissing   float64        `json:"avg_missing"`
	Delta        int            `json:"delta"` // last minus first total drift count
}

// String returns a human-readable representation of the trend summary.
func (t TrendSummary) String() string {
	return fmt.Sprintf(
		"Trend: %s over %d checks (avg drifted: %.1f, avg missing: %.1f, delta: %+d)",
		t.Direction, t.WindowSize, t.AvgDrifted, t.AvgMissing, t.Delta,
	)
}

// AnalyzeTrend computes a TrendSummary from a slice of HistoryEntry values.
// It uses up to the last `window` entries. If fewer entries exist, all are used.
func AnalyzeTrend(entries []HistoryEntry, window int) (TrendSummary, error) {
	if len(entries) == 0 {
		return TrendSummary{}, fmt.Errorf("no history entries provided")
	}
	if window <= 0 {
		return TrendSummary{}, fmt.Errorf("window must be greater than zero")
	}

	start := len(entries) - window
	if start < 0 {
		start = 0
	}
	slice := entries[start:]

	var totalDrifted, totalMissing int
	for _, e := range slice {
		totalDrifted += e.DriftedCount
		totalMissing += e.MissingCount
	}

	n := len(slice)
	firstTotal := slice[0].DriftedCount + slice[0].MissingCount
	lastTotal := slice[n-1].DriftedCount + slice[n-1].MissingCount
	delta := lastTotal - firstTotal

	var direction TrendDirection
	switch {
	case delta > 0:
		direction = TrendIncreasing
	case delta < 0:
		direction = TrendDecreasing
	default:
		direction = TrendStable
	}

	return TrendSummary{
		Direction:    direction,
		WindowSize:   n,
		FirstChecked: slice[0].Timestamp,
		LastChecked:  slice[n-1].Timestamp,
		AvgDrifted:   float64(totalDrifted) / float64(n),
		AvgMissing:   float64(totalMissing) / float64(n),
		Delta:        delta,
	}, nil
}
