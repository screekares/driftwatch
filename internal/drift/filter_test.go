package drift

import (
	"testing"
	"time"
)

func baseReport() Report {
	return Report{
		Provider:  "mock",
		CreatedAt: time.Now(),
		Results: []Result{
			{
				ResourceID:   "res-1",
				ResourceType: "instance",
				Drifted:      true,
				Labels:       map[string]string{"env": "prod"},
			},
			{
				ResourceID:   "res-2",
				ResourceType: "bucket",
				Drifted:      false,
				Labels:       map[string]string{"env": "staging"},
			},
			{
				ResourceID:   "res-3",
				ResourceType: "instance",
				Drifted:      false,
				Labels:       map[string]string{"env": "prod"},
			},
		},
	}
}

func TestFilter_NoConstraints(t *testing.T) {
	r := baseReport()
	f := Filter{}
	out := f.Apply(r)
	if len(out.Results) != len(r.Results) {
		t.Fatalf("expected %d results, got %d", len(r.Results), len(out.Results))
	}
}

func TestFilter_OnlyDrifted(t *testing.T) {
	f := Filter{OnlyDrifted: true}
	out := f.Apply(baseReport())
	for _, res := range out.Results {
		if !res.Drifted {
			t.Errorf("expected only drifted results, got non-drifted %q", res.ResourceID)
		}
	}
	if len(out.Results) != 1 {
		t.Fatalf("expected 1 drifted result, got %d", len(out.Results))
	}
}

func TestFilter_ByResourceType(t *testing.T) {
	f := Filter{ResourceTypes: []string{"instance"}}
	out := f.Apply(baseReport())
	if len(out.Results) != 2 {
		t.Fatalf("expected 2 instance results, got %d", len(out.Results))
	}
}

func TestFilter_ByLabelSelector(t *testing.T) {
	f := Filter{LabelSelector: map[string]string{"env": "prod"}}
	out := f.Apply(baseReport())
	if len(out.Results) != 2 {
		t.Fatalf("expected 2 prod results, got %d", len(out.Results))
	}
}

func TestFilter_Combined(t *testing.T) {
	f := Filter{
		ResourceTypes: []string{"instance"},
		OnlyDrifted:   true,
		LabelSelector: map[string]string{"env": "prod"},
	}
	out := f.Apply(baseReport())
	if len(out.Results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(out.Results))
	}
	if out.Results[0].ResourceID != "res-1" {
		t.Errorf("unexpected resource %q", out.Results[0].ResourceID)
	}
}

func TestFilter_PreservesMetadata(t *testing.T) {
	r := baseReport()
	f := Filter{}
	out := f.Apply(r)
	if out.Provider != r.Provider {
		t.Errorf("provider mismatch: want %q got %q", r.Provider, out.Provider)
	}
}
