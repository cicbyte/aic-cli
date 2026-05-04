# AIC CLI

> 🚀 Powerful Command Line Tool for Managing Claude Skills — Search, Install, Manage, All-in-One Solution

[![Go Report Card](https://goreportcard.com/badge/github.com/cicbyte/aic-cli)](https://goreportcard.com/report/github.com/cicbyte/aic-cli)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.23.2-blue)](https://go.dev/dl/)

**[🇨🇳 中文版](README.md)** | **[🇺🇸 English](README_EN.md)**

<!-- screenshot: Add terminal screenshot or recording here -->

## ✨ Features

- **🔍 Smart Search** — Quickly search skills from AIC server with keyword and category filtering
- **📦 One-Click Install** — Download and automatically install to project or global directory, supports zip packages and URL imports
- **🎨 Beautiful TUI** — Built-in interactive terminal interface for visual browsing and operations
- **🗂️ Category Management** — Organize skills by category to easily find the tools you need
- **🔗 Symbolic Links** — Support for symlink management to save storage space
- **📊 Detailed Information** — Display version, downloads, favorites, and other detailed information
- **🧹 Cleanup Tool** — Automatically clean unused skills to keep your environment tidy
- **⚙️ Flexible Configuration** — Support for project-level and global skills directories with automatic `.claude` directory detection

## 📦 Installation

### Requirements

- Go 1.23.2 or higher

### Install from Source

```bash
# Clone repository
git clone https://github.com/cicbyte/aic-cli.git
cd aic-cli

# Build and install
go build -o aic-cli main.go

# Move to PATH
sudo mv aic-cli /usr/local/bin/
```

### Using Go install

```bash
go install github.com/cicbyte/aic-cli@latest
```

### Download from Release

Visit the [Releases](https://github.com/cicbyte/aic-cli/releases) page to download the binary for your system.

## 🚀 Quick Start

```bash
# Interactive TUI mode (recommended)
aic-cli

# Search skills
aic-cli search claude

# Install a skill
aic-cli add 123

# List installed skills
aic-cli list
```

## 📖 Usage

### `aic-cli`

Launch the interactive TUI interface to visually browse and manage skills.

```bash
aic-cli
```

When run without any arguments, it automatically enters interactive mode with a friendly terminal user interface.

### `aic-cli search [keyword]`

Search for skills with keyword and category filtering.

```bash
aic-cli search [keyword] [options]
```

| Option | Alias | Description | Default |
|---|---|---|---|
| `--category` | `-c` | Filter by category ID | `0` |
| `--page` | `-p` | Page number | `1` |
| `--size` | `-n` | Number per page | `20` |

**Examples:**

```bash
# Search all skills
aic-cli search

# Search by keyword
aic-cli search claude

# Filter by category
aic-cli search --category 2

# Paginate through results
aic-cli search --page 2 --size 10
```

### `aic-cli add <id>`

Download and install a skill locally.

```bash
aic-cli add <skill-id>
```

**Examples:**

```bash
# Install a specific skill
aic-cli add 123

# Install to global directory
aic-cli add 123 --global
```

### `aic-cli list`

List locally installed skills.

```bash
aic-cli list [options]
```

| Option | Alias | Description | Default |
|---|---|---|---|
| `--global` | `-g` | Show global skills directory | `false` |

**Examples:**

```bash
# List current project skills
aic-cli list

# List global skills
aic-cli list -g
```

### `aic-cli remove <name>`

Remove an installed skill.

```bash
aic-cli remove <skill-name>
```

### `aic-cli categories`

List all skill categories.

```bash
aic-cli categories
```

### `aic-cli download <id>`

Download skill zip package to current directory.

```bash
aic-cli download <skill-id>
```

### `aic-cli clean`

Clean unused skills and cache.

```bash
aic-cli clean
```

### `aic-cli import <url>`

Import a skill from URL.

```bash
aic-cli import <url>
```

### All Commands

```bash
aic-cli --help
```

| Command | Description |
|---|---|
| `search` | Search skills |
| `add` | Install skill |
| `remove` | Remove skill |
| `list` | List installed skills |
| `categories` | View categories |
| `download` | Download zip package |
| `clean` | Clean cache |
| `import` | Import from URL |
| `tui` | Interactive interface |

## ⚙️ Configuration

### Directory Structure

AIC CLI supports two skills directories:

- **Project-level directory**: `<project>/.claude/skills/`
  - Automatically detects `.claude` directory in current project
  - Suitable for project-specific skills

- **Global directory**: `~/.aic-cli/skills/`
  - Cross-project shared common skills
  - Access using `--global` or `-g` flag

### Configuration File

Configuration file is located at `~/.aic-cli/config.yaml`:

```yaml
# AIC Server Configuration
aic:
  base_url: "https://api.example.com"  # AIC server address
  timeout: 30

log:
  level: "info"
  path: "~/.aic-cli/logs/aic-cli.log"
```

### Environment Variables

| Variable | Description | Default |
|---|---|---|
| `AIC_CLI_BASE_URL` | AIC server API base URL | — |
| `AIC_CLI_CONFIG` | Configuration file path | `~/.aic-cli/config.yaml` |
| `AIC_CLI_LOG_LEVEL` | Log level | `info` |

## 🏗️ Project Structure

```
aic-cli/
├── cmd/              # CLI command definitions
│   ├── root.go       # Root command
│   ├── search.go     # Search command
│   ├── add.go        # Add command
│   ├── list.go       # List command
│   └── tui.go        # TUI interface
├── internal/         # Internal packages
│   ├── api/          # API client
│   ├── common/       # Common definitions
│   ├── log/          # Logging module
│   └── utils/        # Utility functions
└── main.go           # Program entry
```

## 🛠️ Development

### Build Project

```bash
# Clone repository
git clone https://github.com/cicbyte/aic-cli.git
cd aic-cli

# Download dependencies
go mod download

# Build
go build -o aic-cli main.go
```

### Run Tests

```bash
go test ./...
```

### Contributing

Issues and Pull Requests are welcome! Before submitting a PR, please ensure:

1. Code passes `go test` and `go vet`
2. Add necessary test cases
3. Update relevant documentation

## 📄 License

[MIT](LICENSE) © 2026 Cicbyte

## 🙏 Acknowledgments

- [Cobra](https://github.com/spf13/cobra) — Powerful CLI framework
- [Bubbletea](https://github.com/charmbracelet/bubbletea) — Beautiful TUI framework
- [Lipgloss](https://github.com/charmbracelet/lipgloss) — Terminal styling library

---

**Project URL:** [https://github.com/cicbyte/aic-cli](https://github.com/cicbyte/aic-cli)

**Feedback & Suggestions:** Please submit [Issues](https://github.com/cicbyte/aic-cli/issues)
