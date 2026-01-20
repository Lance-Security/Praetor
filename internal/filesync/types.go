package filesync

import (
	"path/filepath"

	"github.com/lachlanharrisdev/praetor/internal/events"
)

// manager-specific types

type Manager struct {
	EngagementDir string
	statePath     string
	eventsPath    string
	sessionID     string
	user          string
	cwd           string
	tags          []string
	state         *State
}

type ManagerContext struct {
	Manager *Manager
	EngDir  string
	Cwd     string
}

// general types

const (
	maxSyncSizeBytes = 5 * 1024 * 1024 // 5mb
)

// State keeps the list of tracked files
type State struct {
	Files     []FileEntry `json:"files"`
	UpdatedAt string      `json:"updated_at"`
}

type FileEntry struct {
	Path       string `json:"path"`        // absolute path
	RelPath    string `json:"rel_path"`    // relative to engagement dir when possible
	LastHash   string `json:"last_hash"`   // sha256 of last snapshot
	LastSynced string `json:"last_synced"` // RFC3339Nano timestamp
	Size       int64  `json:"size"`        // size at last snapshot
}

// SyncResult summarizes one sync attempt.
type SyncResult struct {
	Entry   FileEntry
	Changed bool
	Event   *events.Event
	Reason  string
	Err     error
}

// utility functions

// DisplayPath prefers the relative path if present
func (f *FileEntry) DisplayPath() string {
	if f.RelPath != "" {
		return f.RelPath
	}
	return f.Path
}

// normalizePath returns a cleaned absolute path
func normalizePath(path string) (string, error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	return filepath.Clean(abs), nil
}
