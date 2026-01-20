/*
Copyright © 2025 Lachlan Harris <contact@lachlanharris.dev>
*/
package events

import (
	"os"

	events "github.com/lachlanharrisdev/praetor/pkg/events"
)

// Re-export types from pkg/events for convenience
type Event = events.Event
type ProcessedEvents = events.ProcessedEvents

// Re-export functions from pkg/events
var (
	NewEvent     = events.NewEvent
	NewNote      = events.NewNote
	MarshalJSONL = events.MarshalJSONL
)

// AppendEvent appends the event to the given file path
func AppendEvent(path string, event *Event) error {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o600)
	if err != nil {
		return err
	}
	defer f.Close()

	event.Id, err = GetLastEventId(path)
	if err != nil {
		return err
	}
	event.Id++

	if err := EnsureEventHash(path, event); err != nil {
		return err
	}
	b, err := MarshalJSONL(event)
	if err != nil {
		return err
	}
	_, err = f.Write(b)
	return err
}
