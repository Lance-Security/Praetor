package events

import (
	"os"
	"testing"
	"time"
)

func TestPrepareEvents(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "events_*.jsonl")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())
	tmpfile.Close()

	addTestEvents(t, tmpfile)

	result, err := PrepareEvents(tmpfile.Name())
	if err != nil {
		t.Fatalf("PrepareEvents failed: %v", err)
	}

	if len(result.Events) != 1 {
		t.Errorf("Expected 1 event, got %d", len(result.Events))
	}

	if result.Events[0].Content != "Modified" {
		t.Errorf("Expected 'Modified', got %q", result.Events[0].Content)
	}

	if len(result.Audit) != 2 {
		t.Errorf("Expected 2 audit events, got %d", len(result.Audit))
	}
}

func addTestEvents(t *testing.T, tmpfile *os.File) {
	err := AppendEvent(tmpfile.Name(), NewNote("Original", "session", "/home", "user"))
	if err != nil {
		t.Fatal(err)
	}
	err = AppendEvent(tmpfile.Name(), NewNote("Delete", "session", "/home", "user"))
	if err != nil {
		t.Fatal(err)
	}

	modEvent := &Event{
		Type:      "modify",
		Timestamp: time.Now().UTC().Format(time.RFC3339Nano),
		SessionID: "session",
		Cwd:       "/home",
		User:      "user",
		Content:   "Modified",
		RefId:     1,
	}
	err = AppendEvent(tmpfile.Name(), modEvent)
	if err != nil {
		t.Fatal(err)
	}

	delEvent := &Event{
		Type:      "delete",
		Timestamp: time.Now().UTC().Format(time.RFC3339Nano),
		SessionID: "session",
		Cwd:       "/home",
		User:      "user",
		RefId:     2,
	}
	err = AppendEvent(tmpfile.Name(), delEvent)
	if err != nil {
		t.Fatal(err)
	}
}

func TestFilterEvents(t *testing.T) {
	events := []*Event{
		{Id: 1, Type: "note", Tags: []string{"important"}},
		{Id: 2, Type: "command", Tags: []string{"info"}},
		{Id: 3, Type: "note", Tags: []string{"important"}},
	}

	result := FilterEvents(events, []string{"important"}, nil)
	if len(result) != 2 {
		t.Errorf("Expected 2 events with tag, got %d", len(result))
	}

	result = FilterEvents(events, nil, []string{"note"})
	if len(result) != 2 {
		t.Errorf("Expected 2 note events, got %d", len(result))
	}

	result = FilterEvents(events, nil, nil)
	if len(result) != 3 {
		t.Errorf("Expected 3 events with no filter, got %d", len(result))
	}
}
