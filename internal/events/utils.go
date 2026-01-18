/*
Copyright © 2025 Lachlan Harris <contact@lachlanharris.dev>
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

// GetUser attempts to get the current terminal user
func GetUser() string {
	user := os.Getenv("USER")
	if user == "" {
		user = os.Getenv("USERNAME")
	}
	if user == "" {
		user = os.Getenv("LOGNAME")
	}
	return user
}

// GetLastEventId returns the ID from the last event in the events log file
// If there are no events or no valid ID, it returns 0
func GetLastEventId(path string) (int, error) {
	lastEvent, err := GetLastEvent(path)
	if err != nil {
		return 0, err
	}
	if lastEvent == nil {
		return 0, nil
	}
	if lastEvent.Id < 0 {
		return 0, nil
	}
	return lastEvent.Id, nil
}

// GetEventById retrieves an event by its ID from the events log file
func GetEventById(path string, id int) (*Event, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	decoder := json.NewDecoder(f)
	for {
		var event Event
		if err := decoder.Decode(&event); err != nil {
			if err.Error() == "EOF" {
				break
			}
			return nil, err
		}
		if event.Id == id {
			return &event, nil
		}
	}
	return nil, nil
}

// GetLastNEvents retrieves the last N events from the events log file
func GetLastNEvents(path string, n int) ([]*Event, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	decoder := json.NewDecoder(f)
	var events []*Event

	for {
		var event Event
		if err := decoder.Decode(&event); err != nil {
			if err.Error() == "EOF" {
				break
			}
			return nil, err
		}
		eventCopy := new(Event)
		*eventCopy = event
		events = append(events, eventCopy)
	}

	if len(events) < n {
		return events, nil
	}
	return events[len(events)-n:], nil
}

// GetAllEvents retrieves all events as an array from the event log file
func GetAllEvents(path string) ([]*Event, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	decoder := json.NewDecoder(f)
	var events []*Event

	for {
		var event Event
		if err := decoder.Decode(&event); err != nil {
			if err.Error() == "EOF" {
				break
			}
			return nil, err
		}
		eventCopy := new(Event)
		*eventCopy = event
		events = append(events, eventCopy)
	}

	return events, nil
}
