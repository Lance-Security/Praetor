/*
Copyright © 2025 Lachlan Harris <contact@lachlanharris.dev>
*/
package cmd

import (
	"github.com/lachlanharrisdev/praetor/internal/output"
	"github.com/spf13/cobra"

	"github.com/lachlanharrisdev/praetor/internal/version"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Shows current version info",
	Long:  `Version shows the current version information of the application.`,
	Run: func(cmd *cobra.Command, args []string) {
		output.LogSuccessf("Praetor `pt` version %s (commit: %s, date: %s)", version.Version, version.Commit, version.Date)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// versionCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// versionCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
