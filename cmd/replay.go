/*
Copyright © 2025 Lachlan Harris <contact@lachlanharris.dev>
*/

package cmd

import (
	"fmt"
	"os"

	"github.com/lachlanharrisdev/praetor/internal/engagement"
	"github.com/lachlanharrisdev/praetor/internal/events"
	"github.com/lachlanharrisdev/praetor/internal/formats"
	"github.com/spf13/cobra"
)

// replayCmd represents the replay command
var replayCmd = &cobra.Command{
	Use:   "replay",
	Short: "Shows a timeline of all events in the event log",
	Long: `Replay displays all events from the current engagement's event log
in the desired format (terminal or JSON). Events are displayed in
chronological order with timestamps, types, and content.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}

		engDir, err := engagement.FindEngagementDir(cwd)
		if err != nil {
			return err
		}

		processed, err := events.PrepareEvents(engagement.EventsPath(engDir))
		if err != nil {
			return err
		}

		rawFormat, err := cmd.Flags().GetString("format")
		if err != nil {
			return err
		}

		format := formats.ParseFormat(rawFormat)
		output, err := formats.Render(format, processed)
		if err != nil {
			return err
		}

		fmt.Print(output)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(replayCmd)

	replayCmd.Flags().StringP("format", "f", "terminal", "Output format (terminal, json)")
}
