package drift

import (
	"fmt"
	"io"
	"strings"
)

// Report summarises drift results for multiple resources.
type Report struct {
	Results []Result
}

// HasDrift returns true if any resource has drifted or is missing.
func (r *Report) HasDrift() bool {
	for _, res := range r.Results {
		if res.Status != StatusMatch {
			return true
		}
	}
	return false
}

// Summary counts results by status.
func (r *Report) Summary() (match, drifted, missing int) {
	for _, res := range r.Results {
		switch res.Status {
		case StatusMatch:
			match++
		case StatusDrift:
			drifted++
		case StatusMissing:
			missing++
		}
	}
	return
}

// Write renders a human-readable drift report to w.
func (r *Report) Write(w io.Writer) {
	match, drifted, missing := r.Summary()
	fmt.Fprintf(w, "Drift Report\n%s\n", strings.Repeat("=", 40))
	fmt.Fprintf(w, "  OK: %d  Drifted: %d  Missing: %d\n\n", match, drifted, missing)

	for _, res := range r.Results {
		switch res.Status {
		case StatusMatch:
			fmt.Fprintf(w, "[OK]      %s\n", res.ResourceID)
		case StatusMissing:
			fmt.Fprintf(w, "[MISSING] %s\n", res.ResourceID)
		case StatusDrift:
			fmt.Fprintf(w, "[DRIFT]   %s\n", res.ResourceID)
			for _, d := range res.Diffs {
				fmt.Fprintf(w, "            %s: declared=%q live=%q\n", d.Field, d.Declared, d.Live)
			}
		}
	}
}
