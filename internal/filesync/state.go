package filesync

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sort"
	"time"
)

// loadState reads the state file or initializes
// a new one
func (m *Manager) loadState() error {
	b, err := os.ReadFile(m.statePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			m.state = &State{
				Files:     []FileEntry{},
				UpdatedAt: time.Now().UTC().Format(time.RFC3339Nano),
			}
			return m.saveState()
		}
		return err
	}

	var st State
	if len(b) > 0 {
		if err := json.Unmarshal(b, &st); err != nil {
			return err
		}
	}
	if st.Files == nil {
		st.Files = []FileEntry{}
	}
	if st.UpdatedAt == "" {
		st.UpdatedAt = time.Now().UTC().Format(time.RFC3339Nano)
	}
	m.state = &st
	m.normalizeState()
	return nil
}

// saveState writes the current state to disk
func (m *Manager) saveState() error {
	if m.state == nil {
		return errors.New("filesync state not initialized")
	}
	m.normalizeState()
	b, err := json.MarshalIndent(m.state, "", "  ")
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(m.statePath), 0o755); err != nil {
		return err
	}
	return os.WriteFile(m.statePath, b, 0o600)
}

// touchStateTime updates the state's UpdatedAt timestamp.
func (m *Manager) touchStateTime() {
	if m.state != nil {
		m.state.UpdatedAt = time.Now().UTC().Format(time.RFC3339Nano)
	}
}

// normalizeState de-duplicates entries and
// ensures rel paths are current
func (m *Manager) normalizeState() {
	if m.state == nil {
		return
	}

	seen := make(map[string]FileEntry)
	for _, f := range m.state.Files {
		abs, err := normalizePath(f.Path)
		if err != nil {
			continue
		}
		rel, ok := m.withinEngagement(abs)
		if !ok {
			continue
		}
		f.Path = abs
		f.RelPath = rel
		seen[abs] = f
	}

	m.state.Files = m.state.Files[:0]
	for _, f := range seen {
		m.state.Files = append(m.state.Files, f)
	}

	sort.Slice(m.state.Files, func(i, j int) bool {
		return m.state.Files[i].Path < m.state.Files[j].Path
	})
}
