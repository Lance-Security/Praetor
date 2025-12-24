/*
Copyright © 2025 Lance Security <support@lancesecurity.org>
*/
package events

import (
	"bufio"
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
)

type hashMaterial struct {
	Type      string   `json:"type"`
	Timestamp string   `json:"timestamp"`
	SessionID string   `json:"session_id"`
	Cwd       string   `json:"cwd"`
	User      string   `json:"user"`
	Content   string   `json:"content"`
	Raw       string   `json:"raw"`
	Tags      []string `json:"tags"`
	PrevHash  string   `json:"prev_hash"`
}

// GetPreviousHash returns the last event's stored hash (if present).
// This intentionally does not validate or verify hashes; verification is handled
// separately via VerifyLog.
func GetPreviousHash(path string) (string, error) {
	last, err := readLastNonEmptyLines(path, 1)
	if err != nil {
		return "", err
	}
	if len(last) == 0 {
		return "", nil
	}

	var lastEvent Event
	if err := json.Unmarshal(last[0], &lastEvent); err != nil {
		// if the last line can't be parsed, start a new chain.
		return "", nil
	}
	return lastEvent.Hash, nil
}

// EnsureEventHash sets PrevHash and Hash if they are not already set.
//
// this does not verify existing events; it only chains to the last stored hash
// (if any) so later verification can detect tampering.
func EnsureEventHash(path string, event *Event) error {
	if event == nil {
		return errors.New("nil event")
	}
	if event.Hash != "" {
		return nil
	}

	if event.PrevHash == "" {
		prev, err := GetPreviousHash(path)
		if err != nil {
			return err
		}
		event.PrevHash = prev
	}
	h, err := ComputeEventHash(event)
	if err != nil {
		return err
	}
	event.Hash = h
	return nil
}

// ComputeEventHash computes the event hash without including Hash itself.
func ComputeEventHash(event *Event) (string, error) {
	if event == nil {
		return "", errors.New("nil event")
	}
	tags := event.Tags
	if tags == nil {
		tags = []string{}
	}

	m := hashMaterial{
		Type:      event.Type,
		Timestamp: event.Timestamp,
		SessionID: event.SessionID,
		Cwd:       event.Cwd,
		User:      event.User,
		Content:   event.Content,
		Raw:       event.Raw,
		Tags:      tags,
		PrevHash:  event.PrevHash,
	}
	b, err := json.Marshal(&m)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:]), nil
}

// VerifyEvent recomputes the hash and compares it with the stored value.
func VerifyEvent(event *Event) error {
	if event == nil {
		return errors.New("nil event")
	}
	if event.Hash == "" {
		return nil
	}
	expected, err := ComputeEventHash(event)
	if err != nil {
		return err
	}
	if expected != event.Hash {
		return errors.New("events log appears modified (hash mismatch)")
	}
	return nil
}

// VerifyLog validates the entire hash chain in the given JSONL file.
// It tolerates older logs that contain un-hashed events by treating the chain
// as starting at the first event that includes a hash.
func VerifyLog(path string) error {
	f, err := os.Open(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}
	defer f.Close()

	br := bufio.NewReader(f)
	var lastHash string
	lineNumber := 0

	for {
		b, err := br.ReadBytes('\n')
		if err != nil && !errors.Is(err, io.EOF) {
			return err
		}
		if len(b) == 0 && errors.Is(err, io.EOF) {
			break
		}
		lineNumber++
		line := bytes.TrimSpace(b)
		if len(line) == 0 {
			if errors.Is(err, io.EOF) {
				break
			}
			continue
		}

		var event Event
		if uerr := json.Unmarshal(line, &event); uerr != nil {
			return fmt.Errorf("line %d: %w", lineNumber, uerr)
		}

		if event.Hash == "" {
			if lastHash != "" {
				return fmt.Errorf("line %d: %w", lineNumber, errors.New("hash chain broken (missing hash)"))
			}
		} else {
			if lastHash == "" && event.PrevHash != "" {
				return fmt.Errorf("line %d: %w", lineNumber, errors.New("events log appears modified (unexpected prev_hash)"))
			}
			if lastHash != "" && event.PrevHash != lastHash {
				return fmt.Errorf("line %d: %w", lineNumber, errors.New("events log appears modified (broken hash chain)"))
			}
			if verr := VerifyEvent(&event); verr != nil {
				return fmt.Errorf("line %d: %w", lineNumber, verr)
			}
			lastHash = event.Hash
		}

		if errors.Is(err, io.EOF) {
			break
		}
	}
	return nil
}

func readLastNonEmptyLines(path string, n int) ([][]byte, error) {
	if n <= 0 {
		return nil, nil
	}
	f, err := os.Open(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, err
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		return nil, err
	}
	size := fi.Size()
	if size == 0 {
		return nil, nil
	}

	const (
		initialTail = int64(64 * 1024)
		maxTail     = int64(1024 * 1024)
	)

	for tail := initialTail; ; tail *= 2 {
		if tail > size {
			tail = size
		}
		start := size - tail
		buf := make([]byte, tail)
		if _, err := f.ReadAt(buf, start); err != nil {
			return nil, err
		}
		trimmed := bytes.TrimSpace(buf)
		if len(trimmed) == 0 {
			return nil, nil
		}

		parts := bytes.Split(trimmed, []byte("\n"))
		lines := make([][]byte, 0, n)
		for i := len(parts) - 1; i >= 0 && len(lines) < n; i-- {
			p := bytes.TrimSpace(parts[i])
			if len(p) == 0 {
				continue
			}
			lines = append(lines, p)
		}
		if len(lines) >= n || start == 0 || tail >= maxTail {
			return lines, nil
		}
	}
}
