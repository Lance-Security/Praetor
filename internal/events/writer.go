/*
Copyright © 2025 Lance Security <support@lancesecurity.org>
*/
package events

import (
	"bytes"
	"os"
	"time"

	"github.com/simonfrey/jsonl"
)

// Event represents a single event in the engagement log
type Event struct {
	Id        int      `json:"id"`             // auto-incrementing ID
	Type      string   `json:"type"`           // "note" | "command" | "result" | "modify" | "delete"
	Timestamp string   `json:"timestamp"`      // RFC3339Nano format
	SessionID string   `json:"session_id"`     // Engagement session ID
	Cwd       string   `json:"cwd"`            // Current working directory
	User      string   `json:"user"`           // user who performed the action
	Content   string   `json:"content"`        // Main content of the event
	Raw       string   `json:"raw,omitempty"`  // Optional raw data (e.g command output)
	Tags      []string `json:"tags,omitempty"` // optional tags associated with the event (not implemented)

	// tamper protection
	Hash     string `json:"hash,omitempty"`      // Optional hash for tamper protection
	PrevHash string `json:"prev_hash,omitempty"` // Optional previous hash for tamper protection

	// modification / deletion
	RefId int `json:"ref_id,omitempty"` // Optional reference to another event ID
}

// utility function

// MarshalJSONL marshals the event into JSONL format
func MarshalJSONL(event *Event) ([]byte, error) {
	buff := bytes.Buffer{}
	w := jsonl.NewWriter(&buff)
	if err := w.Write(event); err != nil {
		return nil, err
	}
	return buff.Bytes(), nil
}

// ! potential race condition between readin gprevious hash
// ! and appending new event
// could lock the event log to prevent this but for now
// its out of scope and only occurs with multiple processes
// running simultaneously (not supported)

// AppendEvent appends the event to the given file path
func AppendEvent(path string, event *Event) error {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o600) // lock file early so other functions are deterministic
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

// NewEvent creates a new event
func NewEvent(eventType, content, timestamp, sessionID, cwd, user string, raw string, tags []string) *Event {
	return &Event{
		Id:        0, // to be set when appending to log
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

// NewNote creates a new note event and automatically
// handles standard fields
func NewNote(content, sessionID, cwd, user string) *Event {
	return NewEvent(
		"note",
		content,
		time.Now().UTC().Format(time.RFC3339Nano),
		sessionID,
		cwd,
		user,
		"",
		nil,
	)
}
