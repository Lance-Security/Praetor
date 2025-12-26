/*
Copyright © 2025 Lance Security <support@lancesecurity.org>
*/
package engagement

import (
	"encoding/json"
	"errors"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/lance-security/praetor/internal/version"
)

type Metadata struct {
	EngagementID string `json:"engagement_id"`
	Name         string `json:"name"`
	CreatedAt    string `json:"created_at"`
	ToolVersion  string `json:"tool_version"`
	LastUsed     string `json:"last_used"`
}

// utility functions for directories

// Dir returns the path to the engagement directory given
// the root and engagement name
func Dir(root, name string) string {
	return filepath.Join(root, name)
}

// PraetorDir returns the path to the .praetor directory
// inside the engagement directory
func PraetorDir(engagementDir string) string {
	return filepath.Join(engagementDir, ".praetor")
}

// MetadataPath returns the path to the metadata.json file
// inside the .praetor directory
func MetadataPath(engagementDir string) string {
	return filepath.Join(PraetorDir(engagementDir), "metadata.json")
}

// EventsPath returns the path to the events.jsonl file
// inside the .praetor directory
func EventsPath(engagementDir string) string {
	return filepath.Join(PraetorDir(engagementDir), "events.jsonl")
}

// EnsureEngagement ensures that an engagement directory
// exists at the given root with the given name.
// If the directory does not exist, it is created.
// If a templateDir is provided and exists, its contents
// are copied into the new engagement directory.
func EnsureEngagement(root, name, templateDir string) (string, error) {
	root = filepath.Clean(root)
	name = strings.TrimSpace(name)
	if name == "" {
		return "", errors.New("empty engagement name")
	}

	engagementDir := Dir(root, name)
	if err := os.MkdirAll(engagementDir, 0o755); err != nil {
		return "", err
	}

	if fi, err := os.Stat(templateDir); err == nil && fi.IsDir() {
		if err := copyDir(templateDir, engagementDir); err != nil {
			return "", err
		}
	}

	if err := EnsurePraetorFiles(engagementDir, name); err != nil {
		return "", err
	}

	if err := TouchLastUsed(engagementDir); err != nil {
		return "", err
	}

	return engagementDir, nil
}

// EnsurePraetorFiles ensures that the .praetor directory
// and necessary files exist inside the given engagement directory.
func EnsurePraetorFiles(engagementDir, name string) error {
	if err := os.MkdirAll(PraetorDir(engagementDir), 0o755); err != nil {
		return err
	}

	metaPath := MetadataPath(engagementDir)
	if _, err := os.Stat(metaPath); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return err
		}
		now := time.Now().UTC().Format(time.RFC3339Nano)
		m := Metadata{
			EngagementID: uuid.NewString(),
			Name:         name,
			CreatedAt:    now,
			ToolVersion:  version.Version,
			LastUsed:     now,
		}
		if err := writeJSONFile(metaPath, &m, 0o600); err != nil {
			return err
		}
	}

	eventsPath := EventsPath(engagementDir)
	f, err := os.OpenFile(eventsPath, os.O_CREATE, 0o600)
	if err != nil {
		return err
	}
	return f.Close()
}

// ReadMetadata reads the engagement metadata from the metadata.json file
// inside the .praetor directory
func ReadMetadata(engagementDir string) (Metadata, error) {
	b, err := os.ReadFile(MetadataPath(engagementDir))
	if err != nil {
		return Metadata{}, err
	}
	var m Metadata
	if err := json.Unmarshal(b, &m); err != nil {
		return Metadata{}, err
	}
	return m, nil
}

// TouchLastUsed updates the LastUsed field in the metadata.json file
// to the current time
func TouchLastUsed(engagementDir string) error {
	m, err := ReadMetadata(engagementDir)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}
	m.LastUsed = time.Now().UTC().Format(time.RFC3339Nano)
	if m.ToolVersion == "" {
		m.ToolVersion = version.Version
	}
	return writeJSONFile(MetadataPath(engagementDir), &m, 0o600)
}

// FindEngagementDir searches for the engagement directory
// by looking for a .praetor directory in the current
// and parent directories
func FindEngagementDir(start string) (string, error) {
	dir := filepath.Clean(start)
	for {
		p := PraetorDir(dir)
		if fi, err := os.Stat(p); err == nil && fi.IsDir() {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", errors.New("not in engagement")
		}
		dir = parent
	}
}

// writeJSONFile writes the given value as JSON
// to the specified path with the given permissions
func writeJSONFile(path string, v any, perm fs.FileMode) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, b, perm); err != nil {
		return err
	}
	if err = os.Rename(tmp, path); err != nil {
		_ = os.Remove(tmp) // cleanup temporary file
		return err
	}
	return nil
}

// copyDir copies the contents of the src directory
// to the dst directory, excluding any .praetor directories
func copyDir(src, dst string) error {
	return filepath.WalkDir(src, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		if rel == "." {
			return nil
		}
		if rel == ".praetor" || strings.HasPrefix(rel, ".praetor"+string(filepath.Separator)) {
			return nil
		}
		target := filepath.Join(dst, rel)
		if d.IsDir() {
			return os.MkdirAll(target, 0o755)
		}
		if !d.Type().IsRegular() {
			return nil
		}
		// check for existing files
		if _, statErr := os.Stat(target); statErr == nil {
			return nil
		} else if !os.IsNotExist(statErr) {
			return statErr
		}
		return copyFile(path, target)
	})
}

// copyFile copies a file from src to dst
func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}

	out, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
	if err != nil {
		return err
	}

	if _, err := io.Copy(out, in); err != nil {
		_ = out.Close()
		return err
	}
	return out.Close()
}
