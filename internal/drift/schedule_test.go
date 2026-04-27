package drift

import (
	"testing"
	"time"
)

func TestParseFrequency_Known(t *testing.T) {
	cases := []struct {
		input    string
		want     Frequency
	}{
		{"hourly", FrequencyHourly},
		{"daily", FrequencyDaily},
		{"weekly", FrequencyWeekly},
	}
	for _, tc := range cases {
		got, err := ParseFrequency(tc.input)
		if err != nil {
			t.Errorf("ParseFrequency(%q) unexpected error: %v", tc.input, err)
		}
		if got != tc.want {
			t.Errorf("ParseFrequency(%q) = %v, want %v", tc.input, got, tc.want)
		}
	}
}

func TestParseFrequency_Unknown(t *testing.T) {
	_, err := ParseFrequency("monthly")
	if err == nil {
		t.Error("expected error for unknown frequency, got nil")
	}
}

func TestFrequency_String(t *testing.T) {
	if got := FrequencyDaily.String(); got != "daily" {
		t.Errorf("FrequencyDaily.String() = %q, want \"daily\"", got)
	}
	if got := Frequency(99).String(); got != "unknown" {
		t.Errorf("unknown Frequency.String() = %q, want \"unknown\"", got)
	}
}

func TestFrequency_Interval(t *testing.T) {
	if got := FrequencyHourly.Interval(); got != time.Hour {
		t.Errorf("FrequencyHourly.Interval() = %v, want %v", got, time.Hour)
	}
	if got := FrequencyWeekly.Interval(); got != 7*24*time.Hour {
		t.Errorf("FrequencyWeekly.Interval() = %v, want %v", got, 7*24*time.Hour)
	}
}

func TestSchedule_IsDue_NoLastRun(t *testing.T) {
	start := time.Now().Add(-time.Minute)
	s := &Schedule{Frequency: FrequencyHourly, StartAt: start}
	if !s.IsDue(time.Now()) {
		t.Error("expected schedule to be due when past StartAt and no LastRun")
	}
}

func TestSchedule_IsDue_NotYet(t *testing.T) {
	start := time.Now().Add(time.Hour)
	s := &Schedule{Frequency: FrequencyHourly, StartAt: start}
	if s.IsDue(time.Now()) {
		t.Error("expected schedule not to be due before StartAt")
	}
}

func TestSchedule_IsDue_WithLastRun(t *testing.T) {
	last := time.Now().Add(-2 * time.Hour)
	s := &Schedule{
		Frequency: FrequencyHourly,
		StartAt:   last,
		LastRun:   &last,
	}
	if !s.IsDue(time.Now()) {
		t.Error("expected schedule to be due when interval has elapsed since LastRun")
	}
}

func TestSchedule_NextRun_NoLastRun(t *testing.T) {
	start := time.Now()
	s := &Schedule{Frequency: FrequencyDaily, StartAt: start}
	if got := s.NextRun(); !got.Equal(start) {
		t.Errorf("NextRun() = %v, want %v", got, start)
	}
}

func TestSchedule_NextRun_WithLastRun(t *testing.T) {
	last := time.Now()
	s := &Schedule{
		Frequency: FrequencyDaily,
		StartAt:   last,
		LastRun:   &last,
	}
	want := last.Add(24 * time.Hour)
	if got := s.NextRun(); !got.Equal(want) {
		t.Errorf("NextRun() = %v, want %v", got, want)
	}
}
