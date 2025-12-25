/*
Copyright © 2025 Lance Security <support@lancesecurity.org>
*/

package engagement

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"

	"github.com/lachlanharrisdev/praetor/internal/events"
)

type Status struct {
	EngagementDir string
	Metadata      Metadata
	NoteCount     int
	LastEvent     *events.Event
}

// LoadStatusFromPath finds the current engagement (from the given path or its parents)
// and returns a summary status for it.
func LoadStatusFromPath(startPath string) (*Status, error) {
	engDir, err := FindEngagementDir(startPath)
	if err != nil {
		return nil, err
	}
	return LoadStatus(engDir)
}

// LoadStatus loads engagement metadata and summarises the events log.
func LoadStatus(engagementDir string) (*Status, error) {
	m, err := ReadMetadata(engagementDir)
	if err != nil {
		return nil, err
	}

	eventsPath := EventsPath(engagementDir)
	last, err := events.GetLastEvent(eventsPath)
	if err != nil {
		return nil, err
	}

	notes, err := CountEventsOfType(eventsPath, "note")
	if err != nil {
		return nil, err
	}

	return &Status{
		EngagementDir: filepath.Clean(engagementDir),
		Metadata:      m,
		NoteCount:     notes,
		LastEvent:     last,
	}, nil
}

// CountEventsOfType counts events in a JSONL log matching the given type.
// Missing logs count as 0.
func CountEventsOfType(eventsPath string, eventType string) (int, error) {
	f, err := os.Open(eventsPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return 0, nil
		}
		return 0, err
	}
	defer f.Close()

	r := bufio.NewReader(f)

	type eventTypeOnly struct {
		Type string `json:"type"`
	}

	count := 0
	for {
		b, err := r.ReadBytes('\n')
		if err != nil && !errors.Is(err, io.EOF) {
			return 0, err
		}
		line := bytes.TrimSpace(b)
		if len(line) == 0 {
			if errors.Is(err, io.EOF) {
				break
			}
			continue
		}
		var e eventTypeOnly
		if err := json.Unmarshal(line, &e); err != nil {
			return 0, err
		}
		if e.Type == eventType {
			count++
		}
		if errors.Is(err, io.EOF) {
			break
		}
	}
	return count, nil
}
