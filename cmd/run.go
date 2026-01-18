/*
Copyright © 2025 Lachlan Harris <contact@lachlanharris.dev>
*/
package cmd

import (
	"github.com/spf13/cobra"

	"github.com/lachlanharrisdev/praetor/internal/run"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run [--sandbox -s] <command> [args...]",
	Short: "Runs any CLI command, optionally in an isolated Bastion",
	Long: `Run allows you to run any CLI command, optionally in an isolated "Bastion".
	Bastions offer a container-like, isolated environment for executing commands
	securely during penetration testing engagements.`,
	Args:               cobra.MinimumNArgs(1),
	SilenceUsage:       true,
	DisableFlagParsing: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Manually parse sandbox flag from args (slightly messy but it'll do)
		sandbox := false
		var toolArgs []string
		if len(args) > 0 && (args[0] == "-s" || args[0] == "--sandbox") {
			sandbox = true
			toolArgs = args[1:]
		} else {
			toolArgs = args
		}

		if len(toolArgs) == 0 {
			return cmd.Help()
		}

		if sandbox {
			return run.RunInBastion(toolArgs)
		}
		return run.RunCmd(toolArgs)
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
	runCmd.Flags().BoolP("sandbox", "s", false, "Run the command in an isolated Bastion") // non-functional but should add to help message

}
