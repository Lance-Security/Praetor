package filesync

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/lachlanharrisdev/praetor/internal/engagement"
)

func setupTestEngagement(t *testing.T) string {
	t.Helper()
	tmpDir := t.TempDir()
	engDir, err := engagement.EnsureEngagement(tmpDir, "test-filesync", "")
	if err != nil {
		t.Fatalf("failed to create test engagement: %v", err)
	}
	return engDir
}

// TestAddRemoveListLifecycle tests the full lifecycle of adding,
// listing and removing synced files
func TestAddRemoveListLifecycle(t *testing.T) {
	engDir := setupTestEngagement(t)

	testFile := filepath.Join(engDir, "notes.txt")
	if err := os.WriteFile(testFile, []byte("test content"), 0o644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	mgr, err := NewManager(engDir)
	if err != nil {
		t.Fatalf("NewManager failed: %v", err)
	}

	if len(mgr.List()) != 0 {
		t.Error("expected empty list initially")
	}

	entry, err := mgr.Add(testFile)
	if err != nil {
		t.Fatalf("Add failed: %v", err)
	}
	if entry.RelPath != "notes.txt" {
		t.Errorf("RelPath = %q, want %q", entry.RelPath, "notes.txt")
	}
	if len(mgr.List()) != 1 {
		t.Fatalf("List returned %d files, want 1", len(mgr.List()))
	}

	// Adding same file again returns existing entry
	if _, err := mgr.Add(testFile); err != nil || len(mgr.List()) != 1 {
		t.Error("duplicate Add should not create new entry")
	}

	removed, ok, err := mgr.Remove(testFile)
	if err != nil || !ok || removed.RelPath != "notes.txt" {
		t.Error("Remove failed or returned wrong entry")
	}
	if len(mgr.List()) != 0 {
		t.Error("expected empty list after Remove")
	}

	_, ok, _ = mgr.Remove(testFile)
	if ok {
		t.Error("Remove should return ok=false for non-existent file")
	}
}

// TestSyncAllDetectsChanges tests various cases for the modification
// (or lack of modification) of synced files
func TestSyncAllDetectsChanges(t *testing.T) {
	engDir := setupTestEngagement(t)

	testFile := filepath.Join(engDir, "data.txt")
	if err := os.WriteFile(testFile, []byte("initial"), 0o644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	mgr, err := NewManager(engDir)
	if err != nil {
		t.Fatalf("NewManager failed: %v", err)
	}
	if _, err := mgr.Add(testFile); err != nil {
		t.Fatalf("Add failed: %v", err)
	}

	// first sync with change (new file)
	results, err := mgr.SyncAll()
	if err != nil || len(results) != 1 || !results[0].Changed {
		t.Error("first SyncAll should detect change")
	}

	// second sync without change
	results, err = mgr.SyncAll()
	if err != nil || results[0].Changed {
		t.Error("SyncAll without changes should not detect change")
	}

	// modify and sync again
	if err := os.WriteFile(testFile, []byte("modified"), 0o644); err != nil {
		t.Fatalf("failed to modify file: %v", err)
	}
	results, err = mgr.SyncAll()
	if err != nil || !results[0].Changed {
		t.Error("SyncAll after modification should detect change")
	}
}
