// ABOUTME: Context store operations for listing and managing saved contexts
// ABOUTME: Handles active context pointer and context enumeration

package config

import (
	"os"
	"path/filepath"
	"strings"
)

// List returns all saved context names.
func List() ([]string, error) {
	dir, err := ContextDir()
	if err != nil {
		return nil, err
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, err
	}

	var contexts []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if strings.HasSuffix(name, ".ctx") {
			contexts = append(contexts, strings.TrimSuffix(name, ".ctx"))
		}
	}

	return contexts, nil
}

// ListContexts returns all saved contexts with their full configuration.
func ListContexts() ([]*Context, error) {
	names, err := List()
	if err != nil {
		return nil, err
	}

	var contexts []*Context
	for _, name := range names {
		ctx, err := Load(name)
		if err != nil {
			continue // Skip contexts that fail to load
		}
		contexts = append(contexts, ctx)
	}

	return contexts, nil
}

// GetActive returns the name of the currently active context.
// Returns empty string if no context is active.
func GetActive() (string, error) {
	path, err := ActiveFile()
	if err != nil {
		return "", err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", err
	}

	return strings.TrimSpace(string(data)), nil
}

// SetActive sets the active context pointer.
func SetActive(name string) error {
	path, err := ActiveFile()
	if err != nil {
		return err
	}

	// Ensure parent directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	return os.WriteFile(path, []byte(name+"\n"), 0644)
}

// ClearActive removes the active context pointer.
func ClearActive() error {
	path, err := ActiveFile()
	if err != nil {
		return err
	}

	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}
