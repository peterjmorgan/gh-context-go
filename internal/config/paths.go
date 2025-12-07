// ABOUTME: Cross-platform path resolution for gh-context configuration
// ABOUTME: Handles XDG on Unix and APPDATA on Windows via go-gh

package config

import (
	"os"
	"path/filepath"

	ghConfig "github.com/cli/go-gh/v2/pkg/config"
)

// ContextDir returns the directory where contexts are stored.
// Uses go-gh's config directory resolution which handles:
// - GH_CONFIG_DIR environment variable
// - XDG_CONFIG_HOME on Unix (~/.config/gh)
// - APPDATA on Windows
func ContextDir() (string, error) {
	configDir := ghConfig.ConfigDir()
	contextDir := filepath.Join(configDir, "contexts")

	// Ensure the directory exists
	if err := os.MkdirAll(contextDir, 0755); err != nil {
		return "", err
	}

	return contextDir, nil
}

// ContextFile returns the full path to a context file by name.
func ContextFile(name string) (string, error) {
	dir, err := ContextDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, name+".ctx"), nil
}

// ActiveFile returns the path to the active context pointer file.
func ActiveFile() (string, error) {
	dir, err := ContextDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "active"), nil
}
