package drift

import (
	"fmt"
	"time"
)

// Frequency represents how often a drift check should run.
type Frequency int

const (
	FrequencyHourly  Frequency = iota
	FrequencyDaily
	FrequencyWeekly
)

// Schedule defines when and how often drift checks should be performed.
type Schedule struct {
	Frequency Frequency    `json:"frequency"`
	StartAt   time.Time    `json:"start_at"`
	LastRun   *time.Time   `json:"last_run,omitempty"`
}

// ParseFrequency converts a string like "hourly", "daily", or "weekly"
// into a Frequency value.
func ParseFrequency(s string) (Frequency, error) {
	switch s {
	case "hourly":
		return FrequencyHourly, nil
	case "daily":
		return FrequencyDaily, nil
	case "weekly":
		return FrequencyWeekly, nil
	default:
		return 0, fmt.Errorf("unknown frequency %q: must be hourly, daily, or weekly", s)
	}
}

// String returns the string representation of a Frequency.
func (f Frequency) String() string {
	switch f {
	case FrequencyHourly:
		return "hourly"
	case FrequencyDaily:
		return "daily"
	case FrequencyWeekly:
		return "weekly"
	default:
		return "unknown"
	}
}

// Interval returns the time.Duration corresponding to the Frequency.
func (f Frequency) Interval() time.Duration {
	switch f {
	case FrequencyHourly:
		return time.Hour
	case FrequencyDaily:
		return 24 * time.Hour
	case FrequencyWeekly:
		return 7 * 24 * time.Hour
	default:
		return 0
	}
}

// IsDue reports whether the schedule's next run time has been reached
// relative to the provided now time.
func (s *Schedule) IsDue(now time.Time) bool {
	if s.LastRun == nil {
		return !now.Before(s.StartAt)
	}
	return now.Sub(*s.LastRun) >= s.Frequency.Interval()
}

// NextRun returns the time at which the next check is due.
func (s *Schedule) NextRun() time.Time {
	if s.LastRun == nil {
		return s.StartAt
	}
	return s.LastRun.Add(s.Frequency.Interval())
}
