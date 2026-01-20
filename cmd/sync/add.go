/*
Copyright © 2025 Lachlan Harris <contact@lachlanharris.dev>
*/
package sync

import (
	"path/filepath"

	"github.com/lachlanharrisdev/praetor/internal/engagement"
	"github.com/lachlanharrisdev/praetor/internal/filesync"
	"github.com/lachlanharrisdev/praetor/internal/formats"
	"github.com/spf13/cobra"
)

// AddCmd represents the sync add command
var AddCmd = &cobra.Command{
	Use:   "add <file>",
	Short: "Add a new file to automatically sync",
	Long: `Sync add will let Praetor watch a new file for changes and
	automatically append file-related events to the event log, as well
	as adding an event showing the file is now monitored.`,
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
		ctx.Manager.SetTags(tagVals)

		entry, err := ctx.Manager.Add(target)
		if err != nil {
			return err
		}

		formats.Successf("Added file for sync: %s", entry.DisplayPath())

		results, err := ctx.Manager.SyncAll()
		if err != nil {
			return err
		}

		if _, err := filesync.PrintSyncResults(results); err != nil {
			return err
		}

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
