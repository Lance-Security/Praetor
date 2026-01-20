/*
Copyright © 2025 Lachlan Harris <contact@lachlanharris.dev>
*/
package cmd

import (
	"github.com/lachlanharrisdev/praetor/cmd/sync"
	"github.com/lachlanharrisdev/praetor/internal/engagement"
	"github.com/lachlanharrisdev/praetor/internal/filesync"
	"github.com/lachlanharrisdev/praetor/internal/formats"
	"github.com/spf13/cobra"
)

// syncCmd represents the sync command
var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Manually sync all syncable files",
	Long: `Sync will manually synchronize all files that have explicitly
	been marked to sync with Praetor.`,

	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, err := filesync.NewManagerFromCwd()
		if err != nil {
			return err
		}

		if len(tags) > 0 {
			ctx.Manager.SetTags(tags)
		}

		results, err := ctx.Manager.SyncAll()
		if err != nil {
			return err
		}

		if len(results) == 0 {
			formats.Info("No files are configured for syncing")
			return nil
		}

		if _, err := filesync.PrintSyncResults(results); err != nil {
			return err
		}

		return engagement.TouchLastUsed(ctx.EngDir)
	},
}

func init() {
	rootCmd.AddCommand(syncCmd)

	syncCmd.AddCommand(sync.AddCmd)
	syncCmd.AddCommand(sync.ListCmd)
	syncCmd.AddCommand(sync.RemoveCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// syncCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// syncCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
