package drift

import (
	"encoding/json"
	"fmt"
	"io"
	"text/tabwriter"
)

// Format represents the output format for drift reports.
type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
)

// Formatter writes a Report to an io.Writer in a specific format.
type Formatter interface {
	Format(r *Report, w io.Writer) error
}

// NewFormatter returns a Formatter for the given format string.
// It returns an error if the format is unsupported.
func NewFormatter(f Format) (Formatter, error) {
	switch f {
	case FormatText:
		return &textFormatter{}, nil
	case FormatJSON:
		return &jsonFormatter{}, nil
	default:
		return nil, fmt.Errorf("unsupported format %q: choose \"text\" or \"json\"", f)
	}
}

// textFormatter renders drift results as a human-readable table.
type textFormatter struct{}

func (t *textFormatter) Format(r *Report, w io.Writer) error {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "RESOURCE\tFIELD\tEXPECTED\tACTUAL\tSTATUS")
	for _, d := range r.Drifts {
		for _, f := range d.Fields {
			fmt.Fprintf(tw, "%s\t%s\t%v\t%v\t%s\n",
				d.ResourceID, f.Field, f.Expected, f.Actual, d.Status)
		}
		if len(d.Fields) == 0 {
			fmt.Fprintf(tw, "%s\t-\t-\t-\t%s\n", d.ResourceID, d.Status)
		}
	}
	return tw.Flush()
}

// jsonFormatter renders drift results as indented JSON.
type jsonFormatter struct{}

func (j *jsonFormatter) Format(r *Report, w io.Writer) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(r)
}
