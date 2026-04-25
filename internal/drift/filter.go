package drift

import "strings"

// Filter holds criteria for narrowing drift results.
type Filter struct {
	// ResourceTypes limits results to the given resource types.
	// An empty slice means no filtering by type.
	ResourceTypes []string

	// OnlyDrifted, when true, excludes resources that have no drift.
	OnlyDrifted bool

	// LabelSelector filters resources whose labels contain all provided
	// key=value pairs. An empty map means no label filtering.
	LabelSelector map[string]string
}

// Apply returns a new Report containing only the Results that satisfy f.
func (f Filter) Apply(r Report) Report {
	filtered := make([]Result, 0, len(r.Results))

	for _, res := range r.Results {
		if !f.matchType(res) {
			continue
		}
		if f.OnlyDrifted && !res.Drifted {
			continue
		}
		if !f.matchLabels(res) {
			continue
		}
		filtered = append(filtered, res)
	}

	return Report{
		Provider:  r.Provider,
		Results:   filtered,
		CreatedAt: r.CreatedAt,
	}
}

func (f Filter) matchType(res Result) bool {
	if len(f.ResourceTypes) == 0 {
		return true
	}
	for _, t := range f.ResourceTypes {
		if strings.EqualFold(res.ResourceType, t) {
			return true
		}
	}
	return false
}

func (f Filter) matchLabels(res Result) bool {
	if len(f.LabelSelector) == 0 {
		return true
	}
	for k, v := range f.LabelSelector {
		got, ok := res.Labels[k]
		if !ok || got != v {
			return false
		}
	}
	return true
}
