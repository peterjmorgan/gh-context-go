// ABOUTME: New command for gh-context - creates a new saved context
// ABOUTME: Supports --from-current to capture current session or explicit parameters

package cmd

import (
	"fmt"
	"os"

	"github.com/pmorgan/gh-context/internal/auth"
	"github.com/pmorgan/gh-context/internal/config"
	"github.com/spf13/cobra"
)

var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Create a new context",
	Long: `Create a new context from the current session or with explicit parameters.

Examples:
  gh context new --from-current --name work
  gh context new --hostname github.com --user myuser --transport ssh --name personal
  gh context new --hostname github.enterprise.com --user myuser --transport https --ssh-host gh-work --name enterprise`,
	RunE: runNew,
}

var (
	newName        string
	newFromCurrent bool
	newHostname    string
	newUser        string
	newTransport   string
	newSSHHost     string
)

func init() {
	newCmd.Flags().StringVar(&newName, "name", "", "Context name (required)")
	newCmd.Flags().BoolVar(&newFromCurrent, "from-current", false, "Create context from current gh session")
	newCmd.Flags().StringVar(&newHostname, "hostname", "", "GitHub hostname (e.g., github.com)")
	newCmd.Flags().StringVar(&newUser, "user", "", "GitHub username")
	newCmd.Flags().StringVar(&newTransport, "transport", "ssh", "Transport protocol (ssh or https)")
	newCmd.Flags().StringVar(&newSSHHost, "ssh-host", "", "SSH host alias for custom SSH configs")

	newCmd.MarkFlagRequired("name")
}

func runNew(cmd *cobra.Command, args []string) error {
	// Validate context name
	if err := config.ValidateName(newName); err != nil {
		return err
	}

	// Check if context already exists
	exists, err := config.Exists(newName)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("context '%s' already exists", newName)
	}

	var hostname, user string

	if newFromCurrent {
		// Get from current session
		hostname = newHostname
		if hostname == "" {
			hostname = os.Getenv("GH_HOST")
		}
		if hostname == "" {
			hostname = "github.com"
		}

		// Get current user from API
		currentUser, authErr := auth.GetCurrentUserFromSession(hostname)
		if authErr != nil {
			printErr("Could not detect current user on '%s'", hostname)
			printInfo("Make sure you're logged in: gh auth login --hostname %s", hostname)
			return fmt.Errorf("authentication required")
		}

		user = newUser
		if user == "" {
			user = currentUser
		}
	} else {
		// Explicit parameters required
		if newHostname == "" || newUser == "" {
			return fmt.Errorf("provide either --from-current or both --hostname and --user")
		}
		hostname = newHostname
		user = newUser
	}

	// Validate transport
	switch newTransport {
	case "ssh", "https":
		// Valid
	default:
		return fmt.Errorf("transport must be 'ssh' or 'https', got: %s", newTransport)
	}

	// Create and save context
	ctx := &config.Context{
		Name:         newName,
		Hostname:     hostname,
		User:         user,
		Transport:    newTransport,
		SSHHostAlias: newSSHHost,
	}

	if err := ctx.Save(); err != nil {
		return err
	}

	sshInfo := ""
	if newSSHHost != "" {
		sshInfo = fmt.Sprintf(", ssh_host=%s", newSSHHost)
	}

	printOk("Created context '%s' → %s@%s (%s%s)", newName, user, hostname, newTransport, sshInfo)
	return nil
}
