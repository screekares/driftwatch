package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/driftwatch/internal/config"
	"github.com/driftwatch/internal/provider"
	"github.com/driftwatch/internal/snapshot"
	"github.com/spf13/cobra"
)

var snapshotOutput string

var snapshotCmd = &cobra.Command{
	Use:   "snapshot",
	Short: "Capture the current state of live resources",
	Long:  `Fetches live resource attributes from the configured provider and saves them as a JSON snapshot file.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load(cfgFile)
		if err != nil {
			return fmt.Errorf("loading config: %w", err)
		}

		p, err := provider.New(cfg.Provider)
		if err != nil {
			return fmt.Errorf("initialising provider: %w", err)
		}

		ids := make([]string, 0, len(cfg.Resources))
		for _, r := range cfg.Resources {
			ids = append(ids, r.ID)
		}

		snap, err := snapshot.Capture(p, ids)
		if err != nil {
			return fmt.Errorf("capturing snapshot: %w", err)
		}
		snap.Provider = cfg.Provider

		if snapshotOutput == "" {
			snapshotOutput = filepath.Join("snapshots",
				fmt.Sprintf("%s.json", time.Now().UTC().Format("20060102T150405Z")))
		}

		if err := os.MkdirAll(filepath.Dir(snapshotOutput), 0o755); err != nil {
			return fmt.Errorf("creating snapshot directory: %w", err)
		}

		if err := snapshot.Save(snap, snapshotOutput); err != nil {
			return fmt.Errorf("saving snapshot: %w", err)
		}

		fmt.Fprintf(os.Stdout, "Snapshot saved to %s (%d resources)\n",
			snapshotOutput, len(snap.Resources))
		return nil
	},
}

func init() {
	snapshotCmd.Flags().StringVarP(&snapshotOutput, "output", "o", "",
		"path to write the snapshot JSON file (default: snapshots/<timestamp>.json)")
	rootCmd.AddCommand(snapshotCmd)
}
