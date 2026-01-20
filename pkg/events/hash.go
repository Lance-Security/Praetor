// Packages are shared with the closed-source praetor server
/*
Copyright © 2025 Lachlan Harris <contact@lachlanharris.dev>
*/

package events

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
)

// hashMaterial contains the fields used for computing an event's hash.
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

// ComputeEventHash computes the SHA-256 hash of an event.
// The hash is computed over all fields except the Hash field itself.
func ComputeEventHash(event *Event) (string, error) {
	if event == nil {
		return "", errors.New("nil event")
	}
	tags := event.Tags
	if tags == nil {
		tags = []string{}
	}
	b, err := json.Marshal(&hashMaterial{
		Id: event.Id, Type: event.Type, Timestamp: event.Timestamp,
		SessionID: event.SessionID, Cwd: event.Cwd, User: event.User,
		Content: event.Content, Raw: event.Raw, Tags: tags,
		PrevHash: event.PrevHash, RefId: event.RefId,
	})
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:]), nil
}

// VerifyEventHash recomputes the hash and compares it with the stored value.
// Returns nil if the event has no hash or if the hash is valid.
func VerifyEventHash(event *Event) error {
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
		return errors.New("hash mismatch: event appears to be modified")
	}
	return nil
}

// SetEventHash computes and sets the Hash field on an event.
// If prevHash is provided, it will also set the PrevHash field.
func SetEventHash(event *Event, prevHash string) error {
	if event == nil {
		return errors.New("nil event")
	}
	if event.Hash != "" {
		return nil
	}
	event.PrevHash = prevHash
	h, err := ComputeEventHash(event)
	if err != nil {
		return err
	}
	event.Hash = h
	return nil
}
