# Godot Language Server (GDLS)

[![CI](https://github.com/andresperezl/gdls/actions/workflows/ci.yml/badge.svg)](https://github.com/andresperezl/gdls/actions/workflows/ci.yml)
[![Go Version](https://img.shields.io/github/go-mod/go-version/andresperezl/gdls)](https://go.dev/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Release](https://img.shields.io/github/v/release/andresperezl/gdls)](https://github.com/andresperezl/gdls/releases)

A Language Server Protocol (LSP) implementation for Godot Engine files, written in Go. Provides IDE features for Text Scene (`.tscn`), External Scene (`.escn`), and Shader (`.gdshader`) files.

## Features

- **Syntax Highlighting** - Semantic tokens for enhanced highlighting beyond TextMate grammars
- **Hover Information** - Rich documentation for nodes, resources, properties, and connections
- **Go to Definition** - Navigate to resource definitions, node parents, and external files
- **Document Symbols** - Hierarchical outline view of your scene structure
- **Auto-completion** - Context-aware completions for node types, resource IDs, node paths, and value constructors
- **Diagnostics** - Real-time error detection for parse errors, missing references, and duplicate IDs
- **Folding** - Collapse sub_resource and node blocks
- **Document Links** - Clickable `res://` paths
- **Find References** - Find all usages of ExtResource/SubResource IDs

## Installation

### Option 1: Go Install (Recommended)

If you have Go 1.25+ installed:

```bash
go install github.com/andresperezl/gdls/cmd/gdls@latest
```

Ensure `$GOPATH/bin` (typically `$HOME/go/bin`) is in your PATH.

### Option 2: Download Binary

Download pre-built binaries for your platform from [GitHub Releases](https://github.com/andresperezl/gdls/releases).

Available platforms:
- Linux (amd64, arm64)
- macOS (amd64, arm64)
- Windows (amd64)

### Option 3: Build from Source

```bash
git clone https://github.com/andresperezl/gdls.git
cd gdls
go build -o gdls ./cmd/gdls

# Optionally, move to a directory in your PATH
mv gdls ~/.local/bin/
```

## Usage

The server communicates via stdio using the Language Server Protocol:

```bash
gdls [options]
```

**Options:**
- `-v`, `--version` - Print version information
- `-h`, `--help` - Print help message

## Editor Integration

### VS Code

A VS Code extension is available in the [`vscode-extension`](./vscode-extension) directory.

**Installation:**
1. Download the `.vsix` file from [GitHub Releases](https://github.com/andresperezl/gdls/releases)
2. In VS Code, open the Command Palette (`Ctrl+Shift+P` / `Cmd+Shift+P`)
3. Run `Extensions: Install from VSIX...`
4. Select the downloaded `.vsix` file

The extension can use a bundled binary or find `gdls` in your PATH.

### Other Editors

Any editor with LSP support can use GDLS. Configure your LSP client to run `gdls` via stdio.

**Neovim (with nvim-lspconfig):**

```lua
vim.api.nvim_create_autocmd("FileType", {
  pattern = { "tscn", "escn", "gdshader" },
  callback = function()
    vim.lsp.start({
      name = "gdls",
      cmd = { "gdls" },
    })
  end,
})
```

**Helix** (`~/.config/helix/languages.toml`):

```toml
[[language]]
name = "tscn"
scope = "source.tscn"
file-types = ["tscn", "escn"]
language-servers = ["gdls"]

[language-server.gdls]
command = "gdls"
```

## Supported File Types

| Extension | Description |
|-----------|-------------|
| `.tscn` | Godot Text Scene files (Godot 4.x format) |
| `.escn` | External Scene files |
| `.gdshader` | Godot Shader files |
| `.gdshaderinc` | Godot Shader include files |

## Development

Requires Go 1.25+ and [Task](https://taskfile.dev/) for build automation.

```bash
# Build for current platform
task build

# Run tests
task test

# Run linter
task lint

# Build for all platforms
task build:all

# See all available tasks
task --list
```

## Contributing

Contributions are welcome! Please:

1. Fork the repository
2. Create a feature branch
3. Use [semantic commits](https://www.conventionalcommits.org/)
4. Submit a pull request

## License

[MIT](./LICENSE)
