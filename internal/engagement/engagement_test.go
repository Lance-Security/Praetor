package engagement

import (
	"os"
	"testing"
)

func TestEnsureEngagement(t *testing.T) {
	tmpDir := t.TempDir()
	engName := "test-eng"

	engDir, err := EnsureEngagement(tmpDir, engName, "")
	if err != nil {
		t.Fatalf("EnsureEngagement failed: %v", err)
	}

	if _, err := os.Stat(engDir); err != nil {
		t.Error("Engagement directory not created")
	}

	praetorDir := PraetorDir(engDir)
	if _, err := os.Stat(praetorDir); err != nil {
		t.Error(".praetor directory not created")
	}
}

func TestReadMetadata(t *testing.T) {
	tmpDir := t.TempDir()
	engDir, err := EnsureEngagement(tmpDir, "test", "")
	if err != nil {
		t.Fatalf("EnsureEngagement failed: %v", err)
	}

	m, err := ReadMetadata(engDir)
	if err != nil {
		t.Fatalf("ReadMetadata failed: %v", err)
	}

	if m.EngagementID == "" {
		t.Error("EngagementID should not be empty")
	}
	if m.Name != "test" {
		t.Errorf("Name = %s, want test", m.Name)
	}
}

func TestTouchLastUsed(t *testing.T) {
	tmpDir := t.TempDir()
	engDir, err := EnsureEngagement(tmpDir, "test", "")
	if err != nil {
		t.Fatalf("EnsureEngagement failed: %v", err)
	}

	err = TouchLastUsed(engDir)
	if err != nil {
		t.Fatalf("TouchLastUsed failed: %v", err)
	}

	m, err := ReadMetadata(engDir)
	if err != nil {
		t.Fatalf("ReadMetadata failed: %v", err)
	}
	if m.LastUsed == "" {
		t.Error("LastUsed should be updated")
	}
}

func TestFindEngagementDir(t *testing.T) {
	tmpDir := t.TempDir()
	engDir, _ := EnsureEngagement(tmpDir, "test", "")

	found, err := FindEngagementDir(engDir)
	if err != nil {
		t.Fatalf("FindEngagementDir failed: %v", err)
	}

	if found != engDir {
		t.Errorf("FindEngagementDir = %s, want %s", found, engDir)
	}
}
