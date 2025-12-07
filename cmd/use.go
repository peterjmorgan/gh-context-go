// ABOUTME: Use command for gh-context - switches to a saved context
// ABOUTME: Sets active context and tests/prompts for authentication

package cmd

import (
	"fmt"

	"github.com/pmorgan/gh-context/internal/auth"
	"github.com/pmorgan/gh-context/internal/config"
	"github.com/spf13/cobra"
)

var useCmd = &cobra.Command{
	Use:   "use <name>",
	Short: "Switch to context (fast by default, prompts to auth if needed)",
	Long: `Switch to a saved context. Sets the context immediately and tests authentication.
If authentication is not configured, provides instructions to set it up.`,
	Args: cobra.ExactArgs(1),
	RunE: runUse,
}

func runUse(cmd *cobra.Command, args []string) error {
	name := args[0]

	// Load context to verify it exists
	ctx, loadErr := config.Load(name)
	if loadErr != nil {
		// Context not found - show available contexts
		contexts, listErr := config.List()
		if listErr == nil && len(contexts) > 0 {
			printErr("Context '%s' not found", name)
			printInfo("Available contexts: %v", contexts)
		}
		return loadErr
	}

	// Set context immediately (fast by default)
	if err := config.SetActive(name); err != nil {
		return err
	}

	printOk("Switched to context '%s' (%s@%s)", name, ctx.User, ctx.Hostname)

	// Test if authentication works
	printInfo("Testing authentication...")
	authenticated, testErr := auth.TestAuth(ctx.Hostname, ctx.User)
	if testErr == nil && authenticated {
		printOk("Authentication verified")
		return nil
	}

	// Authentication failed - prompt user to fix it
	printErr("Authentication required for %s@%s", ctx.User, ctx.Hostname)
	fmt.Println()
	printInfo("Your context has been set, but authentication is needed.")
	printInfo("Please authenticate and your context will work automatically:")
	fmt.Println()
	printInfo("  gh auth login --hostname %s --username %s --scopes repo,read:org", ctx.Hostname, ctx.User)
	fmt.Println()
	printInfo("After authentication, all gh commands will use the correct account.")

	return nil
}
