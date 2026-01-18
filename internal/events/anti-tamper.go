/*
Copyright © 2025 Lachlan Harris <contact@lachlanharris.dev>
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
	"os"
)

type hashMaterial struct {
	Id        int      `json:"id"`
	Type      string   `json:"type"`
	Timestamp string   `json:"timestamp"`
	SessionID string   `json:"session_id"`
	Cwd       string   `json:"cwd"`
	User      string   `json:"user"`
	Content   string   `json:"content"`
	Raw       string   `json:"raw"`
	Tags      []string `json:"tags"`
	PrevHash  string   `json:"prev_hash"`
	RefId     int      `json:"ref_id"`
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
		Id:        event.Id,
		Type:      event.Type,
		Timestamp: event.Timestamp,
		SessionID: event.SessionID,
		Cwd:       event.Cwd,
		User:      event.User,
		Content:   event.Content,
		Raw:       event.Raw,
		Tags:      tags,
		PrevHash:  event.PrevHash,
		RefId:     event.RefId,
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

	scanner := bufio.NewScanner(f)
	return verifyLogLines(scanner)
}

// verifyLogLines processes each line of a log and validates the hash chain.
func verifyLogLines(scanner *bufio.Scanner) error {
	var lastHash string
	lineNumber := 0

	for scanner.Scan() {
		lineNumber++
		line := scanner.Bytes()

		// Skip empty lines
		if len(bytes.TrimSpace(line)) == 0 {
			continue
		}

		var event Event
		if err := json.Unmarshal(line, &event); err != nil {
			return fmt.Errorf("line %d: %w", lineNumber, err)
		}

		if err := verifyEventChain(&event, lastHash, lineNumber); err != nil {
			return err
		}

		// Update lastHash if this event has a hash
		if event.Hash != "" {
			lastHash = event.Hash
		}
	}

	return scanner.Err()
}

// verifyEventChain validates a single event against the chain state
// it checks:
// - if event has no hash but chain has started (missing hash)
// - if event starts chain with unexpected prev_hash
// - if event's prev_hash doesnt match last hash (broken chain)
// - if event's hash itself is valid
func verifyEventChain(event *Event, lastHash string, lineNumber int) error {
	if event.Hash == "" {
		if lastHash != "" {
			return fmt.Errorf("line %d: hash chain broken (missing hash)", lineNumber)
		}
		return nil
	}

	if lastHash == "" {
		// attempted to start a new chain
		if event.PrevHash != "" {
			return fmt.Errorf("line %d: events log appears modified (unexpected prev_hash)", lineNumber)
		}
	} else {
		// attempted to continue existing chain
		if event.PrevHash != lastHash {
			return fmt.Errorf("line %d: events log appears modified (broken hash chain)", lineNumber)
		}
	}

	if err := VerifyEvent(event); err != nil {
		return fmt.Errorf("line %d: %w", lineNumber, err)
	}

	return nil
}

// readLastNonEmptyLines reads the last n non empty lines from a file.
// uses an expanding window approach starting from a small buffer and doubling
// until it has enough lines or reaches the file start
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

	size, err := getFileSize(f)
	if err != nil || size == 0 {
		return nil, err
	}

	return readLinesFromEnd(f, size, n)
}

func getFileSize(f *os.File) (int64, error) {
	fi, err := f.Stat()
	if err != nil {
		return 0, err
	}
	return fi.Size(), nil
}

// readLinesFromEnd reads the last n non empty lines by expanding the read window
// from the end of the file until enough lines are found or the file start is reached
func readLinesFromEnd(f *os.File, size int64, n int) ([][]byte, error) {
	const (
		initialTail = int64(64 * 1024)
		maxTail     = int64(1024 * 1024)
	)

	for tail := initialTail; ; tail *= 2 {
		// Don't read past the start of file
		if tail > size {
			tail = size
		}

		start := size - tail
		lines, err := readAndParseLines(f, start, tail, n)
		if err != nil {
			return nil, err
		}

		if len(lines) >= n || start == 0 || tail >= maxTail {
			return lines, nil
		}
	}
}

// readAndParseLines reads bytes from the file and extracts non-empty lines
// returns up to n lines from the buffer, in reverse order
func readAndParseLines(f *os.File, start int64, tail int64, n int) ([][]byte, error) {
	buf := make([]byte, tail)
	if _, err := f.ReadAt(buf, start); err != nil {
		return nil, err
	}

	trimmed := bytes.TrimSpace(buf)
	if len(trimmed) == 0 {
		return nil, nil
	}

	return extractNonEmptyLines(trimmed, n)
}

// extractNonEmptyLines splits the buffer by newlines and extracts up to n
// non-empty lines in reverse order
func extractNonEmptyLines(buf []byte, n int) ([][]byte, error) {
	parts := bytes.Split(buf, []byte("\n"))
	lines := make([][]byte, 0, n)

	// Iterate from the end of parts backward to maintain reverse order
	for i := len(parts) - 1; i >= 0 && len(lines) < n; i-- {
		p := bytes.TrimSpace(parts[i])
		if len(p) > 0 {
			lines = append(lines, p)
		}
	}

	return lines, nil
}
