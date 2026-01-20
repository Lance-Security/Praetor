/*
Copyright © 2025 Lachlan Harris <contact@lachlanharris.dev>
*/
package sync

import (
	"os"
	"time"

	"github.com/lachlanharrisdev/praetor/internal/filesync"
	"github.com/lachlanharrisdev/praetor/internal/formats"
	"github.com/spf13/cobra"
)

// ListCmd represents the sync list command
var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "Show all synced files and their statuses",
	Long: `Sync list will show all files that have been marked
	to sync with the event log, and show basic metadata such as
	last sync time & current sync status.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, err := filesync.NewManagerFromCwd()
		if err != nil {
			return err
		}

		files := ctx.Manager.List()
		if len(files) == 0 {
			formats.Info("No files are configured for syncing")
			return nil
		}

		formats.Infof("Tracking %d file(s):", len(files))
		for _, f := range files {
			status := "present"
			if _, err := os.Stat(f.Path); err != nil {
				if os.IsNotExist(err) {
					status = "missing"
				} else {
					status = "error"
				}
			}

			last := "never"
			if f.LastSynced != "" {
				if ts, err := time.Parse(time.RFC3339Nano, f.LastSynced); err == nil {
					last = ts.Local().Format("2006-01-02 15:04")
				} else {
					last = f.LastSynced
				}
			}

			hash := f.LastHash
			if len(hash) > 12 {
				hash = hash[:12]
			}

			formats.Infof("%s | status: %s | last sync: %s | size: %d bytes | sha256: %s", f.DisplayPath(), status, last, f.Size, hash)
		}

		return nil
	},
	Args: cobra.NoArgs,
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
