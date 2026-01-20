// Packages are shared with the closed-source praetor server
/*
Copyright © 2025 Lachlan Harris <contact@lachlanharris.dev>
*/

package events

// Event represents a single event in the engagement log.
type Event struct {
	Id        int      `json:"id"`
	Type      string   `json:"type"`
	Timestamp string   `json:"timestamp"`
	SessionID string   `json:"session_id"`
	Cwd       string   `json:"cwd"`
	User      string   `json:"user"`
	Content   string   `json:"content"`
	Raw       string   `json:"raw,omitempty"`
	Tags      []string `json:"tags,omitempty"`
	Hash      string   `json:"hash,omitempty"`
	PrevHash  string   `json:"prev_hash,omitempty"`
	RefId     int      `json:"ref_id,omitempty"`
}

// ProcessedEvents contains the result of processing a raw event log.
type ProcessedEvents struct {
	Events []*Event
	Audit  []*Event
}

// Event type constants.
const (
	TypeNote         = "note"
	TypeCommand      = "command"
	TypeResult       = "result"
	TypeModify       = "modify"
	TypeDelete       = "delete"
	TypeFileSnapshot = "file_snapshot"
	TypeError        = "error"
)
