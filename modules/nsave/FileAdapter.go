package nsave

import (
	"autoworld/domain"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type iFileAdapter interface {
	Write(path string, snap domain.SaveSnapshot) error
	Read(path string) (domain.SaveSnapshot, error)
	Exists(path string) bool
	Delete(path string) error
	ListAll() []string
}

type jsonFileAdapter struct {
	saveDir string
}

func newJsonFileAdapter(saveDir string) *jsonFileAdapter {
	if saveDir == "" {
		saveDir = "./saves"
	}
	return &jsonFileAdapter{saveDir: saveDir}
}

func (a *jsonFileAdapter) getPath(path string) string {
	if path == "" {
		return filepath.Join(a.saveDir, "default.json")
	}
	return path
}

func (a *jsonFileAdapter) Write(path string, snap domain.SaveSnapshot) error {
	p := a.getPath(path)
	// Create directory if not exists
	dir := filepath.Dir(p)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create save directory: %w", err)
	}

	data, err := json.MarshalIndent(snap, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal save data: %w", err)
	}

	if err := os.WriteFile(p, data, 0644); err != nil {
		return fmt.Errorf("failed to write save file: %w", err)
	}

	return nil
}

func (a *jsonFileAdapter) Read(path string) (domain.SaveSnapshot, error) {
	var snap domain.SaveSnapshot
	p := a.getPath(path)

	data, err := os.ReadFile(p)
	if err != nil {
		return snap, fmt.Errorf("failed to read save file: %w", err)
	}

	if err := json.Unmarshal(data, &snap); err != nil {
		return snap, fmt.Errorf("failed to unmarshal save data: %w", err)
	}

	return snap, nil
}

func (a *jsonFileAdapter) Exists(path string) bool {
	p := a.getPath(path)
	_, err := os.Stat(p)
	return err == nil
}

func (a *jsonFileAdapter) Delete(path string) error {
	p := a.getPath(path)
	if !a.Exists(path) {
		return nil
	}
	return os.Remove(p)
}

func (a *jsonFileAdapter) ListAll() []string {
	var slots []string

	entries, err := os.ReadDir(a.saveDir)
	if err != nil {
		return slots
	}

	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".json") {
			name := strings.TrimSuffix(entry.Name(), ".json")
			slots = append(slots, name)
		}
	}

	return slots
}
