// ABOUTME: Context struct definition and serialization for gh-context
// ABOUTME: Handles reading/writing context configuration files (KEY=VALUE format)

package config

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

// Context represents a saved GitHub CLI context (account/host configuration).
type Context struct {
	Name         string // Context name (derived from filename, not stored in file)
	Hostname     string // GitHub host (e.g., github.com, github.enterprise.com)
	User         string // GitHub username
	Transport    string // ssh or https
	SSHHostAlias string // Optional SSH host alias for custom SSH configs
}

// validNamePattern defines valid context name characters.
var validNamePattern = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)

// ValidateName checks if a context name contains only valid characters.
func ValidateName(name string) error {
	if !validNamePattern.MatchString(name) {
		return fmt.Errorf("context name '%s' contains invalid characters (use only alphanumeric, hyphens, underscores)", name)
	}
	return nil
}

// Load reads a context from a .ctx file.
func Load(name string) (*Context, error) {
	path, err := ContextFile(name)
	if err != nil {
		return nil, err
	}

	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("context '%s' not found", name)
		}
		return nil, err
	}
	defer file.Close()

	ctx := &Context{Name: name}
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		switch key {
		case "HOSTNAME":
			ctx.Hostname = value
		case "USER":
			ctx.User = value
		case "TRANSPORT":
			ctx.Transport = value
		case "SSH_HOST_ALIAS":
			ctx.SSHHostAlias = value
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return ctx, nil
}

// Save writes a context to a .ctx file.
func (c *Context) Save() error {
	path, err := ContextFile(c.Name)
	if err != nil {
		return err
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	fmt.Fprintf(file, "HOSTNAME=%s\n", c.Hostname)
	fmt.Fprintf(file, "USER=%s\n", c.User)
	fmt.Fprintf(file, "TRANSPORT=%s\n", c.Transport)
	fmt.Fprintf(file, "SSH_HOST_ALIAS=%s\n", c.SSHHostAlias)

	return nil
}

// Exists checks if a context with the given name exists.
func Exists(name string) (bool, error) {
	path, err := ContextFile(name)
	if err != nil {
		return false, err
	}

	_, err = os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// Delete removes a context file.
func Delete(name string) error {
	path, err := ContextFile(name)
	if err != nil {
		return err
	}

	if err := os.Remove(path); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("context '%s' not found", name)
		}
		return err
	}

	// Clear active pointer if it points to this context
	active, _ := GetActive()
	if active == name {
		if err := ClearActive(); err != nil {
			return err
		}
	}

	return nil
}

// String returns a human-readable representation of the context.
func (c *Context) String() string {
	s := fmt.Sprintf("%s@%s, %s", c.User, c.Hostname, c.Transport)
	if c.SSHHostAlias != "" {
		s += fmt.Sprintf(", ssh_host=%s", c.SSHHostAlias)
	}
	return s
}
