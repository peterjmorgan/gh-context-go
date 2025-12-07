# gh-context

A kubectx-style context switcher for GitHub CLI - manage multiple GitHub accounts and hosts easily.

## Installation

```bash
gh extension install pmorgan/gh-context
```

## Features

- **Cross-platform**: Native binaries for Windows, macOS, and Linux
- **Multi-account**: Switch between personal, work, and enterprise GitHub accounts
- **Repository binding**: Automatically apply contexts when entering specific repositories
- **Shell integration**: Auto-apply contexts on `cd` (bash, zsh, PowerShell, fish)
- **No runtime dependencies**: Single binary, no shell interpreters required

## Quick Start

```bash
# Create a context from your current session
gh context new --from-current --name personal

# Create a context with explicit parameters
gh context new --hostname github.enterprise.com --user myuser --name work

# Switch contexts
gh context use work

# List all contexts
gh context list

# Bind a repository to a context
gh context bind work

# Show current context and repo binding
gh context current
```

## Commands

| Command | Description |
|---------|-------------|
| `list` | List all contexts with active indicator |
| `current` | Show active context and repo-bound context |
| `new` | Create a new context |
| `use <name>` | Switch to a context |
| `delete <name>` | Remove a saved context |
| `bind <name>` | Bind current repository to a context |
| `unbind` | Remove repository binding |
| `apply` | Apply the repo's bound context |
| `shell-hook [shell]` | Print shell integration code |
| `auth-status` | Show authentication status for all contexts |

## Shell Integration

Add automatic context switching when entering repositories:

### Bash
```bash
gh context shell-hook bash >> ~/.bashrc
source ~/.bashrc
```

### Zsh
```bash
gh context shell-hook zsh >> ~/.zshrc
source ~/.zshrc
```

### PowerShell
```powershell
gh context shell-hook powershell >> $PROFILE
. $PROFILE
```

### Fish
```fish
gh context shell-hook fish >> ~/.config/fish/config.fish
source ~/.config/fish/config.fish
```

## Context File Format

Contexts are stored in `~/.config/gh/contexts/` (or `%APPDATA%\gh\contexts` on Windows):

```
HOSTNAME=github.com
USER=myuser
TRANSPORT=ssh
SSH_HOST_ALIAS=
```

## How It Works

1. **Context Storage**: Each context is saved as a `.ctx` file with hostname, user, transport, and optional SSH host alias
2. **Active Context**: A pointer file tracks which context is currently active
3. **Repository Binding**: A `.ghcontext` file in the repo root stores the bound context name
4. **Authentication**: Uses `gh auth switch` to change the active GitHub CLI user
5. **Shell Hooks**: Monitor directory changes and auto-apply bound contexts

## Building from Source

```bash
go build -o gh-context .
```

## Release

Binaries are automatically built for all platforms when a version tag is pushed:

```bash
git tag v1.0.0
git push origin v1.0.0
```

The `gh-extension-precompile` action builds binaries for:
- `darwin-amd64`, `darwin-arm64` (macOS)
- `linux-amd64`, `linux-arm64`, `linux-386`
- `windows-amd64`, `windows-386`, `windows-arm64`

## License

MIT
