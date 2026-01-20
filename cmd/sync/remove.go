/*
Copyright © 2025 Lachlan Harris <contact@lachlanharris.dev>
*/
package sync

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/lachlanharrisdev/praetor/internal/engagement"
	"github.com/lachlanharrisdev/praetor/internal/events"
	"github.com/lachlanharrisdev/praetor/internal/filesync"
	"github.com/lachlanharrisdev/praetor/internal/formats"
	"github.com/spf13/cobra"
)

// RemoveCmd represents the sync remove command
var RemoveCmd = &cobra.Command{
	Use:   "remove <file>",
	Short: "Stops praetor from automatically syncing a file",
	Long: `Sync remove will stop the file monitoring process for a specified
	file, and add an entry to the event log showing that the file is no longer
	monitored.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, err := filesync.NewManagerFromCwd()
		if err != nil {
			return err
		}

		target := args[0]
		if !filepath.IsAbs(target) {
			target = filepath.Join(ctx.Cwd, target)
		}

		tagVals, err := cmd.Flags().GetStringArray("tag")
		if err != nil {
			return err
		}

		removed, ok, err := ctx.Manager.Remove(target)
		if err != nil {
			return err
		}
		if !ok {
			formats.Warnf("File is not configured for sync: %s", target)
			return nil
		}

		meta, err := engagement.ReadMetadata(ctx.EngDir)
		if err != nil {
			return err
		}

		ev := events.NewEvent(
			"file_snapshot",
			fmt.Sprintf("Stopped syncing file %s", removed.DisplayPath()),
			time.Now().UTC().Format(time.RFC3339Nano),
			meta.EngagementID,
			filepath.Clean(ctx.Cwd),
			events.GetUser(),
			"",
			tagVals,
		)

		if err := events.AppendEvent(engagement.EventsPath(ctx.EngDir), ev); err != nil {
			return err
		}

		formats.Successf("Removed file from sync: %s", removed.DisplayPath())
		return engagement.TouchLastUsed(ctx.EngDir)
	},
	Args: cobra.ExactArgs(1),
}

func init() {
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
