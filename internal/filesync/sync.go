package filesync

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/lachlanharrisdev/praetor/internal/events"
	"github.com/lachlanharrisdev/praetor/internal/formats"
)

// PrintSyncResults prints sync results and returns
// the count of changed files
func PrintSyncResults(results []SyncResult) (int, error) {
	changed := 0
	for _, r := range results {
		if r.Err != nil {
			return changed, r.Err
		}
		if r.Changed {
			formats.EmitEvent(r.Event)
			changed++
			continue
		}
		reason := r.Reason
		if reason == "" {
			reason = "unchanged"
		}
		formats.Infof("Skipped %s (%s)", r.Entry.DisplayPath(), reason)
	}

	if changed == 0 {
		formats.Info("No file changes detected")
	} else {
		formats.Successf("Synced %d file(s)", changed)
	}
	return changed, nil
}

// SyncAll iterates tracked files and appends file_snapshot
// events only for changed files
func (m *Manager) SyncAll() ([]SyncResult, error) {
	results := make([]SyncResult, len(m.state.Files))
	anyChanged := false

	for i := range m.state.Files {
		res := m.syncIndex(i)
		results[i] = res
		if res.Err != nil {
			return results, res.Err
		}
		if res.Changed {
			anyChanged = true
		}
	}

	if anyChanged {
		m.touchStateTime()
		if err := m.saveState(); err != nil {
			return results, err
		}
	}

	return results, nil
}

// syncIndex performs sync for a single indexed entry
func (m *Manager) syncIndex(i int) SyncResult {
	entry := m.state.Files[i]

	info, err := os.Stat(entry.Path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return SyncResult{Entry: entry, Changed: false, Reason: "missing"}
		}
		return SyncResult{Entry: entry, Err: err}
	}

	if !info.Mode().IsRegular() {
		return SyncResult{Entry: entry, Changed: false, Reason: "not_regular"}
	}

	if info.Size() > maxSyncSizeBytes {
		return SyncResult{Entry: entry, Changed: false, Reason: "too_large"}
	}

	data, err := os.ReadFile(entry.Path)
	if err != nil {
		return SyncResult{Entry: entry, Err: err}
	}

	hash := sha256.Sum256(data)
	hashHex := hex.EncodeToString(hash[:])

	if entry.LastHash == hashHex {
		return SyncResult{Entry: entry, Changed: false, Reason: "unchanged"}
	}

	now := time.Now().UTC().Format(time.RFC3339Nano)
	content := fmt.Sprintf(
		"Snapshot of %s (sha256=%s, size=%d bytes)",
		entry.DisplayPath(),
		hashHex,
		info.Size(),
	)

	ev := events.NewEvent(
		"file_snapshot",
		content,
		now,
		m.sessionID,
		m.cwd,
		m.user,
		string(data),
		m.tags,
	)

	if err := events.AppendEvent(m.eventsPath, ev); err != nil {
		return SyncResult{Entry: entry, Err: err}
	}

	m.state.Files[i].LastHash = hashHex
	m.state.Files[i].LastSynced = now
	m.state.Files[i].Size = info.Size()

	return SyncResult{Entry: m.state.Files[i], Changed: true, Event: ev, Reason: "updated", Err: nil}
}
