// ABOUTME: Git repository operations for gh-context
// ABOUTME: Handles repo root detection and .ghcontext file management

package git

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const ghContextFile = ".ghcontext"

// RepoRoot returns the root directory of the current git repository.
// Returns empty string if not in a git repository.
func RepoRoot() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	output, err := cmd.Output()
	if err != nil {
		// Not in a git repository
		return "", nil
	}
	return strings.TrimSpace(string(output)), nil
}

// GetBinding reads the context name from .ghcontext in the repo root.
// Returns empty string if no binding exists.
func GetBinding() (string, error) {
	root, err := RepoRoot()
	if err != nil {
		return "", err
	}
	if root == "" {
		return "", nil
	}

	bindingPath := filepath.Join(root, ghContextFile)
	data, err := os.ReadFile(bindingPath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", err
	}

	return strings.TrimSpace(string(data)), nil
}

// SetBinding writes a context name to .ghcontext in the repo root.
func SetBinding(contextName string) error {
	root, err := RepoRoot()
	if err != nil {
		return err
	}
	if root == "" {
		return fmt.Errorf("not inside a Git repository")
	}

	bindingPath := filepath.Join(root, ghContextFile)
	return os.WriteFile(bindingPath, []byte(contextName+"\n"), 0644)
}

// RemoveBinding deletes the .ghcontext file from the repo root.
func RemoveBinding() error {
	root, err := RepoRoot()
	if err != nil {
		return err
	}
	if root == "" {
		return fmt.Errorf("not inside a Git repository")
	}

	bindingPath := filepath.Join(root, ghContextFile)
	if err := os.Remove(bindingPath); err != nil {
		if os.IsNotExist(err) {
			return nil // Already gone, not an error
		}
		return err
	}
	return nil
}

// HasBinding checks if the current repo has a .ghcontext file.
func HasBinding() (bool, error) {
	root, err := RepoRoot()
	if err != nil {
		return false, err
	}
	if root == "" {
		return false, nil
	}

	bindingPath := filepath.Join(root, ghContextFile)
	_, err = os.Stat(bindingPath)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// BindingPath returns the full path to .ghcontext in the current repo.
// Returns empty string if not in a git repository.
func BindingPath() (string, error) {
	root, err := RepoRoot()
	if err != nil {
		return "", err
	}
	if root == "" {
		return "", nil
	}
	return filepath.Join(root, ghContextFile), nil
}
