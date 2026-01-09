# Godot Language Server (GDLS) - VS Code Extension

Language support for Godot Text Scene (`.tscn`), External Scene (`.escn`), and Shader (`.gdshader`) files.

## Installation

### From VSIX File

1. Download the `.vsix` file from [GitHub Releases](https://github.com/andresperezl/gdls/releases)
2. In VS Code, open the Command Palette (`Ctrl+Shift+P` / `Cmd+Shift+P`)
3. Run `Extensions: Install from VSIX...`
4. Select the downloaded `.vsix` file

## Features

- **Syntax Highlighting** - TextMate grammar for basic highlighting, enhanced by semantic tokens from the language server
- **Hover Information** - Rich documentation for nodes, resources, properties, and connections
- **Go to Definition** - Navigate to resource definitions, node parents, and external files
- **Document Symbols** - Hierarchical outline view of your scene structure
- **Auto-completion** - Context-aware completions for node types, resource IDs, node paths, and value constructors
- **Diagnostics** - Real-time error detection for parse errors, missing references, and duplicate IDs
- **Folding** - Collapse sub_resource and node blocks
- **Document Links** - Clickable `res://` paths
- **Find References** - Find all usages of ExtResource/SubResource IDs
- **Semantic Tokens** - Enhanced syntax highlighting based on semantic analysis

## Requirements

The extension requires the `gdls` language server binary. Choose one of the following options:

### Option 1: Use Bundled Binary (Recommended)

If the extension includes a bundled binary for your platform, it will be used automatically. No additional setup required.

### Option 2: Install via Go

If you have Go 1.25+ installed:

```bash
go install github.com/andresperezl/gdls/cmd/gdls@latest
```

Ensure `$GOPATH/bin` (typically `$HOME/go/bin`) is in your PATH.

### Option 3: Download from Releases

Download the appropriate binary for your platform from [GitHub Releases](https://github.com/andresperezl/gdls/releases) and either:
- Place it in a directory in your PATH, or
- Configure the path via the `gdls.server.path` setting (see below)

### Option 4: Configure Custom Path

If you have `gdls` installed in a custom location, configure it in VS Code settings:

```json
{
  "gdls.server.path": "/path/to/gdls"
}
```

## Extension Settings

| Setting | Default | Description |
|---------|---------|-------------|
| `gdls.server.path` | `""` | Path to the gdls executable. Leave empty to use bundled binary or PATH. |
| `gdls.server.enabled` | `true` | Enable/disable the Godot language server. |
| `gdls.trace.server` | `"off"` | Trace communication between VS Code and the language server (`off`, `messages`, `verbose`). |

## Commands

| Command | Description |
|---------|-------------|
| `Godot: Restart Language Server` | Restart the language server |
| `Godot: Show Output Channel` | Show the language server output channel |

## Supported File Types

- `.tscn` - Godot Text Scene files (Godot 4.x format)
- `.escn` - External Scene files
- `.gdshader` - Godot Shader files

## Troubleshooting

### Language server not starting

1. Check that `gdls` is installed and accessible:
   ```bash
   gdls --version
   ```
2. If using a custom path, verify the `gdls.server.path` setting points to the correct location
3. Check the output channel (`Godot: Show Output Channel`) for error messages
4. Enable tracing (`gdls.trace.server`: `"verbose"`) for detailed logs

### Features not working

1. Ensure the file has the correct extension (`.tscn`, `.escn`, or `.gdshader`)
2. Try restarting the language server (`Godot: Restart Language Server`)
3. Check for errors in the Problems panel

## License

MIT
