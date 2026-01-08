package lsp

import (
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"

	"github.com/andresperezl/gdls/internal/analysis"
	"github.com/andresperezl/gdls/internal/parser"
)

// textDocumentDefinition handles the textDocument/definition request.
func (s *Server) textDocumentDefinition(ctx *glsp.Context, params *protocol.DefinitionParams) (any, error) {
	doc := s.workspace.GetDocument(params.TextDocument.URI)
	if doc == nil || doc.TSCNAST == nil {
		return nil, nil
	}

	line := int(params.Position.Line)
	col := int(params.Position.Character)

	// Find what's at this position and where it's defined
	location := s.findDefinition(doc, params.TextDocument.URI, line, col)
	if location == nil {
		return nil, nil
	}

	return location, nil
}

// findDefinition finds the definition location for the item at the given position.
func (s *Server) findDefinition(doc *analysis.Document, uri string, line, col int) *protocol.Location {
	ast := doc.TSCNAST

	// Check if we're on a resource reference in a property value
	for _, sub := range ast.SubResources {
		for _, prop := range sub.Properties {
			if loc := s.findDefinitionInValue(prop.Value, ast, uri, line, col); loc != nil {
				return loc
			}
		}
	}

	for _, node := range ast.Nodes {
		// Check instance references
		if node.Instance != nil {
			if loc := s.findDefinitionInValue(node.Instance, ast, uri, line, col); loc != nil {
				return loc
			}
		}
		// Check property values
		for _, prop := range node.Properties {
			if loc := s.findDefinitionInValue(prop.Value, ast, uri, line, col); loc != nil {
				return loc
			}
		}
	}

	// Check if we're on an ext_resource path
	for _, ext := range ast.ExtResources {
		if isInRange(ext.Range, line, col) {
			// Return location to the file itself
			return s.resolveResourcePath(ext.Path, uri)
		}
	}

	// Check if we're on a node with a parent reference
	for _, node := range ast.Nodes {
		if isInRange(node.Range, line, col) && node.Parent != "" && node.Parent != "." {
			// Find the parent node
			return s.findNodeByPath(ast, node.Parent, uri)
		}
	}

	return nil
}

// findDefinitionInValue finds definition for a resource reference in a value.
func (s *Server) findDefinitionInValue(v parser.Value, ast *parser.Document, uri string, line, col int) *protocol.Location {
	if v == nil {
		return nil
	}

	switch val := v.(type) {
	case *parser.ResourceRef:
		if isInRange(val.Range, line, col) {
			switch val.RefType {
			case "ExtResource":
				// Find the ext_resource definition
				for _, ext := range ast.ExtResources {
					if ext.ID == val.ID {
						return &protocol.Location{
							URI: uri,
							Range: protocol.Range{
								Start: protocol.Position{
									Line:      uint32(ext.Range.Start.Line),
									Character: uint32(ext.Range.Start.Column),
								},
								End: protocol.Position{
									Line:      uint32(ext.Range.End.Line),
									Character: uint32(ext.Range.End.Column),
								},
							},
						}
					}
				}
			case "SubResource":
				// Find the sub_resource definition
				for _, sub := range ast.SubResources {
					if sub.ID == val.ID {
						return &protocol.Location{
							URI: uri,
							Range: protocol.Range{
								Start: protocol.Position{
									Line:      uint32(sub.Range.Start.Line),
									Character: uint32(sub.Range.Start.Column),
								},
								End: protocol.Position{
									Line:      uint32(sub.Range.End.Line),
									Character: uint32(sub.Range.End.Column),
								},
							},
						}
					}
				}
			}
		}

	case *parser.ArrayValue:
		for _, elem := range val.Values {
			if loc := s.findDefinitionInValue(elem, ast, uri, line, col); loc != nil {
				return loc
			}
		}

	case *parser.DictValue:
		for _, entry := range val.Entries {
			if loc := s.findDefinitionInValue(entry.Value, ast, uri, line, col); loc != nil {
				return loc
			}
		}

	case *parser.TypedValue:
		for _, arg := range val.Arguments {
			if loc := s.findDefinitionInValue(arg, ast, uri, line, col); loc != nil {
				return loc
			}
		}
	}

	return nil
}

// findNodeByPath finds a node by its path and returns its location.
func (s *Server) findNodeByPath(ast *parser.Document, path, uri string) *protocol.Location {
	// Build node path map
	pathToNode := make(map[string]*parser.Node)
	for _, node := range ast.Nodes {
		var nodePath string
		switch node.Parent {
		case "":
			nodePath = ""
		case ".":
			nodePath = node.Name
		default:
			nodePath = node.Parent + "/" + node.Name
		}
		pathToNode[nodePath] = node
	}

	// Find the target node
	if node, ok := pathToNode[path]; ok {
		return &protocol.Location{
			URI: uri,
			Range: protocol.Range{
				Start: protocol.Position{
					Line:      uint32(node.Range.Start.Line),
					Character: uint32(node.Range.Start.Column),
				},
				End: protocol.Position{
					Line:      uint32(node.Range.End.Line),
					Character: uint32(node.Range.End.Column),
				},
			},
		}
	}

	return nil
}

// resolveResourcePath resolves a resource path to an actual file URI.
func (s *Server) resolveResourcePath(resPath, currentURI string) *protocol.Location {
	// Handle res:// paths
	if strings.HasPrefix(resPath, "res://") {
		relPath := strings.TrimPrefix(resPath, "res://")

		// Find the project root (where project.godot is located)
		projectRoot := s.findProjectRoot(currentURI)
		if projectRoot == "" {
			return nil
		}

		targetPath := filepath.Join(projectRoot, relPath)
		targetURI := pathToURI(targetPath)

		return &protocol.Location{
			URI: targetURI,
			Range: protocol.Range{
				Start: protocol.Position{Line: 0, Character: 0},
				End:   protocol.Position{Line: 0, Character: 0},
			},
		}
	}

	// Handle relative paths
	currentPath := uriToPath(currentURI)
	if currentPath == "" {
		return nil
	}
	currentDir := filepath.Dir(currentPath)
	targetPath := filepath.Join(currentDir, resPath)
	targetURI := pathToURI(targetPath)

	return &protocol.Location{
		URI: targetURI,
		Range: protocol.Range{
			Start: protocol.Position{Line: 0, Character: 0},
			End:   protocol.Position{Line: 0, Character: 0},
		},
	}
}

// findProjectRoot finds the Godot project root by looking for project.godot file.
func (s *Server) findProjectRoot(currentURI string) string {
	// First, try to use workspace folders
	for _, folder := range s.workspace.GetFolders() {
		folderPath := uriToPath(folder)
		if folderPath == "" {
			continue
		}

		// Check if project.godot exists in this folder
		if fileExists(filepath.Join(folderPath, "project.godot")) {
			return folderPath
		}
	}

	// Fallback: walk up from the current file's directory to find project.godot
	currentPath := uriToPath(currentURI)
	if currentPath == "" {
		return ""
	}

	dir := filepath.Dir(currentPath)
	for dir != "/" && dir != "." && dir != "" {
		if fileExists(filepath.Join(dir, "project.godot")) {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break // Reached root
		}
		dir = parent
	}

	return ""
}

// uriToPath converts a file:// URI to a filesystem path, handling Windows paths correctly.
func uriToPath(uri string) string {
	parsedURI, err := url.Parse(uri)
	if err != nil {
		return ""
	}

	path := parsedURI.Path

	// On Windows, file:///C:/path becomes /C:/path after parsing
	// We need to remove the leading slash for Windows paths
	if len(path) >= 3 && path[0] == '/' && path[2] == ':' {
		// Looks like /C:/... - remove leading slash
		path = path[1:]
	}

	return path
}

// pathToURI converts a filesystem path to a file:// URI, handling Windows paths correctly.
func pathToURI(path string) string {
	// On Windows, paths look like C:\Users\... or C:/Users/...
	// We need to convert to file:///C:/Users/...
	if len(path) >= 2 && path[1] == ':' {
		// Windows absolute path - add extra slash
		return "file:///" + filepath.ToSlash(path)
	}
	// Unix path or relative path
	return "file://" + path
}

// fileExists checks if a file exists at the given path.
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
