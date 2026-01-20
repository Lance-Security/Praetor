package filesync

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/lachlanharrisdev/praetor/internal/engagement"
	"github.com/lachlanharrisdev/praetor/internal/events"
)

// NewManagerFromCwd creates a manager by finding the
// engagement directory from the cwd
func NewManagerFromCwd() (*ManagerContext, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	engDir, err := engagement.FindEngagementDir(cwd)
	if err != nil {
		return nil, err
	}
	mgr, err := NewManager(engDir)
	if err != nil {
		return nil, err
	}
	return &ManagerContext{Manager: mgr, EngDir: engDir, Cwd: cwd}, nil
}

// NewManager loads or initializes state for an engagement directory.
func NewManager(engDir string) (*Manager, error) {
	engDir = filepath.Clean(engDir)

	meta, err := engagement.ReadMetadata(engDir)
	if err != nil {
		return nil, err
	}

	statePath := filepath.Join(engagement.PraetorDir(engDir), "filesync.json")

	wd, err := os.Getwd()
	if err != nil {
		wd = engDir
	}

	m := &Manager{
		EngagementDir: engDir,
		statePath:     statePath,
		eventsPath:    engagement.EventsPath(engDir),
		sessionID:     meta.EngagementID,
		user:          events.GetUser(),
		cwd:           filepath.Clean(wd),
	}

	if err := m.loadState(); err != nil {
		return nil, err
	}

	return m, nil
}

// Add registers a file for syncing. If already present, it is left unchanged.
func (m *Manager) Add(path string) (*FileEntry, error) {
	abs, err := normalizePath(path)
	if err != nil {
		return nil, err
	}

	rel, ok := m.withinEngagement(abs)
	if !ok {
		return nil, fmt.Errorf("file must reside within engagement directory: %s", m.EngagementDir)
	}

	info, err := os.Stat(abs)
	if err != nil {
		return nil, err
	}
	if !info.Mode().IsRegular() {
		return nil, fmt.Errorf("path is not a regular file: %s", abs)
	}

	if existing := m.find(abs); existing != nil {
		return existing, nil
	}

	entry := FileEntry{
		Path:    abs,
		RelPath: rel,
	}

	m.state.Files = append(m.state.Files, entry)
	m.normalizeState()
	m.touchStateTime()
	return &entry, m.saveState()
}

// Remove unregisters a file. Returns removed entry and whether it was removed.
func (m *Manager) Remove(path string) (*FileEntry, bool, error) {
	abs, err := normalizePath(path)
	if err != nil {
		return nil, false, err
	}

	idx := -1
	for i, f := range m.state.Files {
		if f.Path == abs {
			idx = i
			break
		}
	}
	if idx == -1 {
		return nil, false, nil
	}

	removed := m.state.Files[idx]
	m.state.Files = append(m.state.Files[:idx], m.state.Files[idx+1:]...)
	m.normalizeState()
	m.touchStateTime()
	return &removed, true, m.saveState()
}

// List returns a snapshot of tracked files.
func (m *Manager) List() []FileEntry {
	out := make([]FileEntry, len(m.state.Files))
	copy(out, m.state.Files)
	return out
}

// SetTags sets tags applied to emitted file_snapshot events.
func (m *Manager) SetTags(tags []string) {
	m.tags = append([]string(nil), tags...)
}

// find returns the entry for the absolute path, or nil.
func (m *Manager) find(absPath string) *FileEntry {
	for i := range m.state.Files {
		if m.state.Files[i].Path == absPath {
			return &m.state.Files[i]
		}
	}
	return nil
}

// withinEngagement returns rel path and true if inside engagement (excluding .praetor).
func (m *Manager) withinEngagement(absPath string) (string, bool) {
	rel, err := filepath.Rel(m.EngagementDir, absPath)
	if err != nil || rel == "" || rel == "." {
		return "", false
	}
	if strings.HasPrefix(rel, "..") || strings.HasPrefix(rel, ".praetor") {
		return "", false
	}
	return rel, true
}
