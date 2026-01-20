package events

import (
	events "github.com/lachlanharrisdev/praetor/pkg/events"
)

// PrepareEvents takes a raw `events.jsonl` and sequentially applies modify
// and delete operations, and returns ProcessedEvents with the tamper-proofing
// removed.
func PrepareEvents(path string) (*ProcessedEvents, error) {
	rawEvents, err := GetAllEvents(path)
	if err != nil {
		return nil, err
	}

	return events.ProcessEvents(rawEvents), nil
}

// Re-export FilterEvents from pkg/events for convenience
var FilterEvents = events.FilterEvents
