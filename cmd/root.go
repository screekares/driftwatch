package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	configFile string
	profile    string
)

// rootCmd is the base command for the driftwatch CLI.
var rootCmd = &cobra.Command{
	Use:   "driftwatch",
	Short: "Detect configuration drift between deployed services and IaC definitions",
	Long: `driftwatch compares the live state of your cloud infrastructure
against your declared infrastructure-as-code definitions and reports
any configuration drift it finds.`,
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(
		&configFile, "config", "c", "",
		"path to driftwatch config file (default: .driftwatch.json)",
	)
	rootCmd.PersistentFlags().StringVarP(
		&profile, "profile", "p", "",
		"profile name to use from config file",
	)
}
