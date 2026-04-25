package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"driftwatch/internal/config"
	"driftwatch/internal/drift"
	"driftwatch/internal/provider"
	"driftwatch/internal/snapshot"
)

var (
	checkConfigFile string
	checkFormat     string
	checkSnapshot   string
)

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Detect drift between live infrastructure and a saved snapshot",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load(checkConfigFile)
		if err != nil {
			return fmt.Errorf("loading config: %w", err)
		}

		p, err := provider.New(cfg.Provider)
		if err != nil {
			return fmt.Errorf("initialising provider: %w", err)
		}

		snap, err := snapshot.Load(checkSnapshot)
		if err != nil {
			return fmt.Errorf("loading snapshot: %w", err)
		}

		detector := drift.New(p)
		report, err := detector.Check(cmd.Context(), snap.Resources)
		if err != nil {
			return fmt.Errorf("running drift check: %w", err)
		}

		fmt.Fprintf(cmd.OutOrStdout(), "%s\n", report.Summary())

		formatter, err := drift.NewFormatter(checkFormat)
		if err != nil {
			return fmt.Errorf("creating formatter: %w", err)
		}

		if err := report.Write(cmd.OutOrStdout(), formatter); err != nil {
			return fmt.Errorf("writing report: %w", err)
		}

		if report.HasDrift() {
			os.Exit(1)
		}
		return nil
	},
}

func init() {
	checkCmd.Flags().StringVarP(&checkConfigFile, "config", "c", "driftwatch.yaml", "path to config file")
	checkCmd.Flags().StringVarP(&checkFormat, "format", "f", "text", "output format: text or json")
	checkCmd.Flags().StringVarP(&checkSnapshot, "snapshot", "s", "", "path to snapshot file (required)")
	_ = checkCmd.MarkFlagRequired("snapshot")
	RootCmd.AddCommand(checkCmd)
}
