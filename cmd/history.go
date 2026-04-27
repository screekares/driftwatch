package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"

	"github.com/driftwatch/internal/drift"
)

// historyCmd represents the history command, which displays past drift check
// results recorded in the local history file.
var historyCmd = &cobra.Command{
	Use:   "history",
	Short: "Show past drift check results",
	Long: `Display a chronological list of previous drift checks stored in the
local history file. Each entry shows the timestamp, total resources checked,
number of drifted resources, missing resources, and overall severity.

Use --format json to emit machine-readable output.`,
	RunE: runHistory,
}

var (
	hisoryFile   string
	historyLimit int
	historyFmt   string
)

func init() {
	rootCmd.AddCommand(historyCmd)

	historyCmd.Flags().StringVar(&hisoryFile, "history-file", "drift-history.json",
		"Path to the history file")
	historyCmd.Flags().IntVar(&historyLimit, "limit", 20,
		"Maximum number of entries to display (0 = all)")
	historyCmd.Flags().StringVar(&historyFmt, "format", "text",
		"Output format: text or json")
}

func runHistory(cmd *cobra.Command, _ []string) error {
	entries, err := drift.LoadHistory(hisoryFile)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Fprintln(cmd.OutOrStdout(), "No history found. Run 'driftwatch check' to generate entries.")
			return nil
		}
		return fmt.Errorf("loading history: %w", err)
	}

	if len(entries) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "History is empty.")
		return nil
	}

	// Apply limit (most-recent first).
	if historyLimit > 0 && len(entries) > historyLimit {
		entries = entries[len(entries)-historyLimit:]
	}

	// Reverse so newest appears at the top.
	for i, j := 0, len(entries)-1; i < j; i, j = i+1, j-1 {
		entries[i], entries[j] = entries[j], entries[i]
	}

	switch historyFmt {
	case "json":
		return printHistoryJSON(cmd, entries)
	default:
		return printHistoryText(cmd, entries)
	}
}

func printHistoryText(cmd *cobra.Command, entries []drift.HistoryEntry) error {
	w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "TIMESTAMP\tTOTAL\tDRIFTED\tMISSING\tSEVERITY")
	for _, e := range entries {
		ts := e.Timestamp.Format(time.RFC3339)
		fmt.Fprintf(w, "%s\t%d\t%d\t%d\t%s\n",
			ts, e.TotalResources, e.DriftedCount, e.MissingCount, e.Severity)
	}
	return w.Flush()
}

func printHistoryJSON(cmd *cobra.Command, entries []drift.HistoryEntry) error {
	enc := json.NewEncoder(cmd.OutOrStdout())
	enc.SetIndent("", "  ")
	return enc.Encode(entries)
}
