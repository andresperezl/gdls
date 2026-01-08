package lsp

import (
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"

	"github.com/andresperezl/gdls/internal/analysis"
	"github.com/andresperezl/gdls/internal/parser"
)

// textDocumentReferences handles the textDocument/references request.
func (s *Server) textDocumentReferences(ctx *glsp.Context, params *protocol.ReferenceParams) ([]protocol.Location, error) {
	doc := s.workspace.GetDocument(params.TextDocument.URI)
	if doc == nil || doc.TSCNAST == nil {
		return nil, nil
	}

	line := int(params.Position.Line)
	col := int(params.Position.Character)

	// Find what's at this position
	locations := s.findReferences(doc, params.TextDocument.URI, line, col, params.Context.IncludeDeclaration)

	return locations, nil
}

// findReferences finds all references to the item at the given position.
func (s *Server) findReferences(doc *analysis.Document, uri string, line, col int, includeDeclaration bool) []protocol.Location {
	ast := doc.TSCNAST
	locations := []protocol.Location{}

	// Check if we're on an ext_resource
	for _, ext := range ast.ExtResources {
		if isInRange(ext.Range, line, col) {
			// Find all references to this ext_resource ID
			if includeDeclaration {
				locations = append(locations, protocol.Location{
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
				})
			}
			locations = append(locations, s.findResourceReferences(ast, ext.ID, "ExtResource", uri)...)
			return locations
		}
	}

	// Check if we're on a sub_resource
	for _, sub := range ast.SubResources {
		if isInRange(sub.Range, line, col) {
			if includeDeclaration {
				locations = append(locations, protocol.Location{
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
				})
			}
			locations = append(locations, s.findResourceReferences(ast, sub.ID, "SubResource", uri)...)
			return locations
		}
	}

	// Check if we're on a node
	for _, node := range ast.Nodes {
		if isInRange(node.Range, line, col) {
			// Find the node's path
			var nodePath string
			switch node.Parent {
			case "":
				nodePath = ""
			case ".":
				nodePath = node.Name
			default:
				nodePath = node.Parent + "/" + node.Name
			}

			if includeDeclaration {
				locations = append(locations, protocol.Location{
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
				})
			}
			locations = append(locations, s.findNodeReferences(ast, nodePath, uri)...)
			return locations
		}
	}

	return locations
}

// findResourceReferences finds all references to a resource ID.
func (s *Server) findResourceReferences(ast *parser.Document, id, refType, uri string) []protocol.Location {
	locations := []protocol.Location{}

	var findInValue func(v parser.Value)
	findInValue = func(v parser.Value) {
		if v == nil {
			return
		}

		switch val := v.(type) {
		case *parser.ResourceRef:
			if val.RefType == refType && val.ID == id {
				locations = append(locations, protocol.Location{
					URI: uri,
					Range: protocol.Range{
						Start: protocol.Position{
							Line:      uint32(val.Range.Start.Line),
							Character: uint32(val.Range.Start.Column),
						},
						End: protocol.Position{
							Line:      uint32(val.Range.End.Line),
							Character: uint32(val.Range.End.Column),
						},
					},
				})
			}
		case *parser.ArrayValue:
			for _, elem := range val.Values {
				findInValue(elem)
			}
		case *parser.DictValue:
			for _, entry := range val.Entries {
				findInValue(entry.Value)
			}
		case *parser.TypedValue:
			for _, arg := range val.Arguments {
				findInValue(arg)
			}
		}
	}

	// Search in sub_resource properties
	for _, sub := range ast.SubResources {
		for _, prop := range sub.Properties {
			findInValue(prop.Value)
		}
	}

	// Search in node properties and instances
	for _, node := range ast.Nodes {
		if node.Instance != nil {
			findInValue(node.Instance)
		}
		for _, prop := range node.Properties {
			findInValue(prop.Value)
		}
	}

	return locations
}

// findNodeReferences finds all references to a node path (in parent= and connections).
func (s *Server) findNodeReferences(ast *parser.Document, nodePath, uri string) []protocol.Location {
	locations := []protocol.Location{}

	// Find nodes that have this path as parent
	for _, node := range ast.Nodes {
		if node.Parent == nodePath || (nodePath == "" && node.Parent == ".") {
			// This node references our target as its parent
			locations = append(locations, protocol.Location{
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
			})
		}
	}

	// Find connections that reference this node
	for _, conn := range ast.Connections {
		if conn.From == nodePath || conn.To == nodePath {
			locations = append(locations, protocol.Location{
				URI: uri,
				Range: protocol.Range{
					Start: protocol.Position{
						Line:      uint32(conn.Range.Start.Line),
						Character: uint32(conn.Range.Start.Column),
					},
					End: protocol.Position{
						Line:      uint32(conn.Range.End.Line),
						Character: uint32(conn.Range.End.Column),
					},
				},
			})
		}
	}

	return locations
}
