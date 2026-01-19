/*
Copyright © 2025 Lachlan Harris <contact@lachlanharris.dev>
*/
package cmd

import (
	"os"

	"github.com/lachlanharrisdev/praetor/internal/config"
	"github.com/lachlanharrisdev/praetor/internal/formats"
	"github.com/lachlanharrisdev/praetor/internal/utils"
	"github.com/spf13/cobra"
)

var tags []string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "pt",
	Short: "Praetor revolutionizes state management for penetration testing engagements",
	Long: `Praetor is a powerful tool designed to automatically log and organise note-taking
in penetration testing engagements, as well as offering a minimal suite of tools
to isolate, secure and manage the engagement environment.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		rawFormat, err := cmd.Flags().GetString("format")
		if err != nil {
			return err
		}
		format := formats.ParseFormat(rawFormat)

		// temporary until config is loaded
		utils.ConfigureTerminal(false, false)

		formats.SetDefault(formats.NewEmitter(formats.Options{
			Format:       format,
			Writer:       os.Stdout,
			UseTimestamp: false,
		}))

		cfg, err := config.Load()
		if err != nil {
			formats.Errorf("Failed to load configuration: %v", err)
			return err
		}

		if format == formats.FormatTerminal {
			utils.ConfigureTerminal(cfg.UseColour, cfg.UseBold)
		}

		return nil
	},
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// Repeatable, works on all subcommands, supports both --tag and -t.
	rootCmd.PersistentFlags().StringArrayVarP(&tags, "tag", "t", nil, "add an optional tag to the created event, if applicable (repeatable)")

	// Output format
	rootCmd.PersistentFlags().StringP("format", "f", "terminal", "set the output format (e.g., terminal, json)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
