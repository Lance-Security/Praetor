/*
Copyright © 2025 Lance Security <support@lancesecurity.org>
*/

package events

import (
	"encoding/json"
	"os"
)

// GetLastEvent reads and returns the last event from the given events log file as an event struct
func GetLastEvent(path string) (*Event, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	stat, err := f.Stat()
	if err != nil {
		return nil, err
	}
	if stat.Size() == 0 {
		return nil, nil
	}

	events, err := readLastNonEmptyLines(path, 1)
	if err != nil {
		return nil, err
	}
	if len(events) == 0 {
		return nil, nil
	}

	var lastEvent Event
	if err := json.Unmarshal(events[0], &lastEvent); err != nil {
		return nil, err
	}
	return &lastEvent, nil
}
