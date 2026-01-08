# Godot Language Server (GDLS)

Language support for Godot Text Scene (`.tscn`), External Scene (`.escn`), and Shader (`.gdshader`) files.

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

The extension requires the `gdls` language server binary. You can:

1. **Use the bundled binary** (if included in the extension)
2. **Install globally** and ensure `gdls` is in your PATH
3. **Configure a custom path** via the `gdls.server.path` setting

### Building the Language Server

```bash
# Clone the gdls repository
git clone https://github.com/andresperezl/gdls.git
cd gdls

# Build for your platform
go build -o gdls ./cmd/gdls

# Optionally, move to a directory in your PATH
mv gdls ~/.local/bin/
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

## Development

```bash
# Install dependencies
cd vscode-extension
npm install

# Compile
npm run compile

# Watch for changes
npm run watch

# Package the extension
npm run package
```

## License

MIT
