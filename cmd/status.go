/*
Copyright © 2025 Lachlan Harris <contact@lachlanharris.dev>
*/
package cmd

import (
	"os"
	"time"

	"github.com/lachlanharrisdev/praetor/internal/engagement"
	"github.com/lachlanharrisdev/praetor/internal/events"
	"github.com/lachlanharrisdev/praetor/internal/output"
	"github.com/lachlanharrisdev/praetor/internal/utils"

	"github.com/spf13/cobra"
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Shows the current status of the engagement",
	Long: `Status shows the current status of the engagement,
	including the engagement name, start time, number of notes
	recorded and the latest activity.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}

		s, err := engagement.LoadStatusFromPath(cwd)
		if err != nil {
			return err
		}

		output.LogPrimary(s.Metadata.Name)

		started := s.Metadata.CreatedAt
		if ts, err := time.Parse(time.RFC3339Nano, s.Metadata.CreatedAt); err == nil {
			started = ts.Local().Format("2006-01-02 15:04")
		}
		output.Indent()
		defer output.Dedent()
		output.Println(utils.Muted("Started:"), utils.Default(started))
		output.Println(utils.Muted("Notes:"), utils.Primary(s.NoteCount))

		if s.LastEvent == nil {
			output.Println(utils.Muted("Latest:"), utils.Muted("(none)"))
			return nil
		}
		output.Println(utils.Muted("Latest:"), events.ShowEventTerminal(*s.LastEvent))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// statusCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// statusCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
