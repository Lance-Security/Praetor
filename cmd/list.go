/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/lachlanharrisdev/praetor/internal/engagement"
	"github.com/lachlanharrisdev/praetor/internal/events"
	"github.com/lachlanharrisdev/praetor/internal/formats"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list [<event_count>]",
	Short: "Shows the most recent events",
	Long: `List shows the last 10 events and their metadata by
	default, or the number supplied as an argument`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}
		engDir, err := engagement.FindEngagementDir(cwd)
		if err != nil {
			return err
		}

		n := 10
		if len(args) > 0 {
			var err error
			n, err = strconv.Atoi(args[0])
			if err != nil {
				return err
			}
		}
		e, err := events.GetLastNEvents(engagement.EventsPath(engDir), n)
		if err != nil {
			return err
		}

		msgs := make([]formats.Message, 0, len(e))
		for _, ev := range e {
			msgs = append(msgs, formats.Message{Level: formats.LevelInfo, Event: ev})
		}

		out, err := formats.RenderMessages(formats.FormatTerminal, msgs, formats.Options{Format: formats.FormatTerminal})
		if err != nil {
			return err
		}

		fmt.Print(out)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
