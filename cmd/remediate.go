package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"driftwatch/internal/drift"
	"driftwatch/internal/snapshot"
)

var remediateOutputFormat string

var remediateCmd = &cobra.Command{
	Use:   "remediate",
	Short: "Suggest remediation actions for detected drift",
	Long:  `Loads a snapshot, runs drift detection, and outputs a remediation plan.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		snapshotPath, _ := cmd.Flags().GetString("snapshot")
		if snapshotPath == "" {
			return fmt.Errorf("--snapshot flag is required")
		}

		snap, err := snapshot.Load(snapshotPath)
		if err != nil {
			return fmt.Errorf("loading snapshot: %w", err)
		}

		detector := drift.New(snap)
		report, err := detector.Check(cmd.Context())
		if err != nil {
			return fmt.Errorf("running drift check: %w", err)
		}

		plan := drift.BuildRemediationPlan(report)

		switch remediateOutputFormat {
		case "json":
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(plan)
		default:
			fmt.Println(plan.Summary())
			for _, a := range plan.Actions {
				fmt.Printf("  [%s] %s (%s): %s\n", a.Severity, a.ResourceID, a.ResourceType, a.Suggestion)
			}
		}

		return nil
	},
}

func init() {
	remediateCmd.Flags().String("snapshot", "", "Path to snapshot file (required)")
	remediateCmd.Flags().StringVarP(&remediateOutputFormat, "output", "o", "text", "Output format: text or json")
	rootCmd.AddCommand(remediateCmd)
}
