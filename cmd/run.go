/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"

	"github.com/lachlanharrisdev/praetor/internal/bastions"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Runs any CLI command in an isolated Bastion",
	Long: `Run allows you to run any CLI command within an isolated "Bastion", offering
	a container-like, isolated environment for executing commands securely during
	penetration testing engagements.`,
	DisableFlagParsing: true,
	Args:               cobra.MinimumNArgs(1),
	SilenceUsage:       true,
	RunE: func(cmd *cobra.Command, args []string) error {
		return bastions.RunInBastion(args)
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// runCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
