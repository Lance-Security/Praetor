/*
Copyright © 2025 Lachlan Harris <contact@lachlanharris.dev>
*/
package cmd

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/spf13/cobra"

	"github.com/lachlanharrisdev/praetor/internal/engagement"
	"github.com/lachlanharrisdev/praetor/internal/events"
	"github.com/lachlanharrisdev/praetor/internal/output"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete [<event_id>]",
	Short: "Delete the last event or an event by its ID",
	Long: `Delete allows you to append a delete request to the event log,
	either for the most recent event or for a specific event by its ID.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return deleteEvent(args)
	},
}

// deleteEvent appends a delete event to the log for the specified event ID
func deleteEvent(args []string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	engDir, err := engagement.FindEngagementDir(cwd)
	if err != nil {
		return err
	}

	// get ID from args or get last event
	eventId, err := findIdByArgs(args, engDir)
	if err != nil {
		return err
	}

	dEvent, err := events.GetEventById(engagement.EventsPath(engDir), eventId)
	if err != nil {
		return err
	}
	if dEvent == nil || eventId == 0 {
		output.LogWarningf("Event ID %d not found\n", eventId)
		return fmt.Errorf("event ID %d not found", eventId)
	}

	m, err := engagement.ReadMetadata(engDir)
	if err != nil {
		return err
	}
	user := events.GetUser()

	// create a new event with operation "delete"
	deleteEvent := events.Event{
		Type:      "delete",
		Timestamp: time.Now().UTC().Format(time.RFC3339Nano),
		SessionID: m.EngagementID,
		Cwd:       cwd,
		User:      user,
		Content:   fmt.Sprintf("Delete event ID %d", eventId),
		RefId:     eventId,
		Tags:      tags,
	}

	// append the delete event to the log
	err = events.AppendEvent(engagement.EventsPath(engDir), &deleteEvent)
	if err != nil {
		return err
	}

	output.LogSuccessf("Successfully appended delete event for ID %d\n", eventId)
	return nil
}

func findIdByArgs(args []string, engDir string) (int, error) {
	if len(args) == 0 {
		return events.GetLastEventId(engagement.EventsPath(engDir))
	} else {
		// args[0] is a string, convert to int
		id, err := strconv.Atoi(args[0])
		if err != nil {
			return 0, fmt.Errorf("invalid event ID: %s", args[0])
		}
		return id, nil
	}
}

func init() {
	rootCmd.AddCommand(deleteCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deleteCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deleteCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
