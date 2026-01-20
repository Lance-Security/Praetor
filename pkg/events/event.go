// Packages are shared with the closed-source praetor server
/*
Copyright © 2025 Lachlan Harris <contact@lachlanharris.dev>
*/

package events

import (
	"bytes"
	"time"

	"github.com/simonfrey/jsonl"
)

// NewEvent creates a new event with the specified parameters.
// The Id field is left at 0 and should be set when appending to a log.
func NewEvent(eventType, content, timestamp, sessionID, cwd, user, raw string, tags []string) *Event {
	return &Event{
		Type:      eventType,
		Timestamp: timestamp,
		SessionID: sessionID,
		Cwd:       cwd,
		User:      user,
		Content:   content,
		Raw:       raw,
		Tags:      tags,
	}
}

// NewNote creates a new note event with the current timestamp.
func NewNote(content, sessionID, cwd, user string) *Event {
	return NewEvent(TypeNote, content, time.Now().UTC().Format(time.RFC3339Nano), sessionID, cwd, user, "", nil)
}

// MarshalJSONL marshals the event into JSONL format.
func MarshalJSONL(event *Event) ([]byte, error) {
	var buf bytes.Buffer
	if err := jsonl.NewWriter(&buf).Write(event); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
