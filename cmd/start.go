/*
Copyright © 2025 Lachlan Harris <contact@lachlanharris.dev>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/lachlanharrisdev/praetor/internal/config"
	"github.com/lachlanharrisdev/praetor/internal/engagement"
	"github.com/lachlanharrisdev/praetor/internal/formats"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start <name>",
	Short: "Initialize a new penetration testing engagement environment",
	Long: `Start initializes a new penetration testing engagement environment
	by creating a new folder structure and configuration, and changes into the
	new directory.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}
		engDir, err := engagement.EnsureEngagement(cfg.EngagementRoot, args[0], cfg.TemplateDir)
		if err != nil {
			return err
		}
		if err = os.Chdir(engDir); err != nil {
			return fmt.Errorf("Failed to change directory to %s: %w", engDir, err)
		}
		formats.Success(engDir)
		return nil
	},
	Aliases: []string{"engage", "init"},
	Args:    cobra.MinimumNArgs(1), // <name> argument is required
}

func init() {
	rootCmd.AddCommand(startCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// startCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// startCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
