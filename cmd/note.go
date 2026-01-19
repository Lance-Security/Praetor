/*
Copyright © 2025 Lachlan Harris <contact@lachlanharris.dev>
*/
package cmd

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/lachlanharrisdev/praetor/internal/engagement"
	"github.com/lachlanharrisdev/praetor/internal/events"
	"github.com/lachlanharrisdev/praetor/internal/formats"
)

// noteCmd represents the note command
var noteCmd = &cobra.Command{
	Use:   "note <comment...>",
	Short: "Add a note to the current engagement log",
	Long: `Note adds a note entry to the current engagement's event log.
	Everything after "note" is treated as the note content, including quotations.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}
		engDir, err := engagement.FindEngagementDir(cwd)
		if err != nil {
			return err
		}
		m, err := engagement.ReadMetadata(engDir)
		if err != nil {
			return err
		}
		user := events.GetUser()
		content := strings.Join(args, " ")

		n := events.NewNote(
			content,
			m.EngagementID,
			filepath.Clean(cwd),
			user,
		)

		if len(tags) > 0 {
			n.Tags = append(n.Tags, tags...)
		}

		if err := events.AppendEvent(engagement.EventsPath(engDir), n); err != nil {
			return err
		}
		formats.Success("Successfully added note")
		return engagement.TouchLastUsed(engDir)
	},
	Args:    cobra.MinimumNArgs(1),
	Aliases: []string{"log"},
}

func init() {
	rootCmd.AddCommand(noteCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// noteCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// noteCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
