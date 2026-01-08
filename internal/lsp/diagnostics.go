package lsp

import (
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"

	"github.com/andresperezl/gdls/internal/analysis"
	"github.com/andresperezl/gdls/internal/parser"
)

// publishDiagnostics publishes diagnostics for a document.
func (s *Server) publishDiagnostics(ctx *glsp.Context, uri string, doc *analysis.Document) {
	if doc == nil {
		return
	}

	switch doc.Type {
	case analysis.DocumentTypeTSCN:
		s.publishTSCNDiagnostics(ctx, uri, doc)
	case analysis.DocumentTypeGDShader:
		s.publishGDShaderDiagnostics(ctx, uri, doc)
	}
}

// publishTSCNDiagnostics publishes diagnostics for a TSCN document.
func (s *Server) publishTSCNDiagnostics(ctx *glsp.Context, uri string, doc *analysis.Document) {
	if doc.TSCNAST == nil {
		return
	}

	diagnostics := []protocol.Diagnostic{}

	// Add parse errors
	for _, err := range doc.TSCNAST.Errors {
		diagnostics = append(diagnostics, protocol.Diagnostic{
			Range: protocol.Range{
				Start: protocol.Position{
					Line:      uint32(err.Range.Start.Line),
					Character: uint32(err.Range.Start.Column),
				},
				End: protocol.Position{
					Line:      uint32(err.Range.End.Line),
					Character: uint32(err.Range.End.Column),
				},
			},
			Severity: severityPtr(protocol.DiagnosticSeverityError),
			Source:   strPtr("gdls"),
			Message:  err.Message,
		})
	}

	// Check format version
	if doc.TSCNAST.Descriptor != nil && doc.TSCNAST.Descriptor.Format != 3 {
		diagnostics = append(diagnostics, protocol.Diagnostic{
			Range: protocol.Range{
				Start: protocol.Position{
					Line:      uint32(doc.TSCNAST.Descriptor.Range.Start.Line),
					Character: uint32(doc.TSCNAST.Descriptor.Range.Start.Column),
				},
				End: protocol.Position{
					Line:      uint32(doc.TSCNAST.Descriptor.Range.End.Line),
					Character: uint32(doc.TSCNAST.Descriptor.Range.End.Column),
				},
			},
			Severity: severityPtr(protocol.DiagnosticSeverityError),
			Source:   strPtr("gdls"),
			Message:  "Only format=3 (Godot 4.x) is supported",
		})
	}

	// Check for missing resource references
	diagnostics = append(diagnostics, s.checkResourceReferences(doc)...)

	// Check for missing parent nodes
	diagnostics = append(diagnostics, s.checkParentReferences(doc)...)

	// Check for duplicate resource IDs
	diagnostics = append(diagnostics, s.checkDuplicateIDs(doc)...)

	ctx.Notify(protocol.ServerTextDocumentPublishDiagnostics, protocol.PublishDiagnosticsParams{
		URI:         uri,
		Diagnostics: diagnostics,
	})
}

// checkResourceReferences checks for references to non-existent resources.
func (s *Server) checkResourceReferences(doc *analysis.Document) []protocol.Diagnostic {
	diagnostics := []protocol.Diagnostic{}

	// Build maps of valid IDs
	extIDs := make(map[string]bool)
	subIDs := make(map[string]bool)

	for _, ext := range doc.TSCNAST.ExtResources {
		extIDs[ext.ID] = true
	}
	for _, sub := range doc.TSCNAST.SubResources {
		subIDs[sub.ID] = true
	}

	// Check all resource references in the document
	checkValue := func(v parser.Value) {
		if ref, ok := v.(*parser.ResourceRef); ok {
			var valid bool
			switch ref.RefType {
			case "ExtResource":
				valid = extIDs[ref.ID]
			case "SubResource":
				valid = subIDs[ref.ID]
			}
			if !valid {
				diagnostics = append(diagnostics, protocol.Diagnostic{
					Range: protocol.Range{
						Start: protocol.Position{
							Line:      uint32(ref.Range.Start.Line),
							Character: uint32(ref.Range.Start.Column),
						},
						End: protocol.Position{
							Line:      uint32(ref.Range.End.Line),
							Character: uint32(ref.Range.End.Column),
						},
					},
					Severity: severityPtr(protocol.DiagnosticSeverityError),
					Source:   strPtr("gdls"),
					Message:  "Reference to undefined resource: " + ref.ID,
				})
			}
		}
	}

	// Check in sub_resource properties
	for _, sub := range doc.TSCNAST.SubResources {
		for _, prop := range sub.Properties {
			walkValue(prop.Value, checkValue)
		}
	}

	// Check in node properties
	for _, node := range doc.TSCNAST.Nodes {
		if node.Instance != nil {
			checkValue(node.Instance)
		}
		for _, prop := range node.Properties {
			walkValue(prop.Value, checkValue)
		}
	}

	return diagnostics
}

// checkParentReferences checks for references to non-existent parent nodes.
func (s *Server) checkParentReferences(doc *analysis.Document) []protocol.Diagnostic {
	diagnostics := []protocol.Diagnostic{}

	// Build a map of valid node paths
	nodePaths := make(map[string]bool)
	var rootName string

	for _, node := range doc.TSCNAST.Nodes {
		if node.Parent == "" {
			// Root node
			rootName = node.Name
			nodePaths[""] = true
			nodePaths["."] = true
		} else {
			// Build the full path
			var fullPath string
			if node.Parent == "." {
				fullPath = node.Name
			} else {
				fullPath = node.Parent + "/" + node.Name
			}
			nodePaths[fullPath] = true
		}
	}

	// Check parent references
	for _, node := range doc.TSCNAST.Nodes {
		if node.Parent == "" || node.Parent == "." {
			continue
		}

		// The parent path should exist (without including the root name)
		if !nodePaths[node.Parent] && node.Parent != rootName {
			diagnostics = append(diagnostics, protocol.Diagnostic{
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
				Severity: severityPtr(protocol.DiagnosticSeverityWarning),
				Source:   strPtr("gdls"),
				Message:  "Parent node not found: " + node.Parent,
			})
		}
	}

	return diagnostics
}

// checkDuplicateIDs checks for duplicate resource IDs.
func (s *Server) checkDuplicateIDs(doc *analysis.Document) []protocol.Diagnostic {
	diagnostics := []protocol.Diagnostic{}

	// Check external resources
	extIDs := make(map[string]*parser.ExtResource)
	for _, ext := range doc.TSCNAST.ExtResources {
		if existing, ok := extIDs[ext.ID]; ok {
			diagnostics = append(diagnostics, protocol.Diagnostic{
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
				Severity: severityPtr(protocol.DiagnosticSeverityError),
				Source:   strPtr("gdls"),
				Message:  "Duplicate external resource ID: " + ext.ID,
			})
			_ = existing
		} else {
			extIDs[ext.ID] = ext
		}
	}

	// Check sub resources
	subIDs := make(map[string]*parser.SubResource)
	for _, sub := range doc.TSCNAST.SubResources {
		if existing, ok := subIDs[sub.ID]; ok {
			diagnostics = append(diagnostics, protocol.Diagnostic{
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
				Severity: severityPtr(protocol.DiagnosticSeverityError),
				Source:   strPtr("gdls"),
				Message:  "Duplicate sub-resource ID: " + sub.ID,
			})
			_ = existing
		} else {
			subIDs[sub.ID] = sub
		}
	}

	return diagnostics
}

// walkValue recursively walks a value and calls the callback for each value.
func walkValue(v parser.Value, cb func(parser.Value)) {
	if v == nil {
		return
	}
	cb(v)

	switch val := v.(type) {
	case *parser.ArrayValue:
		for _, elem := range val.Values {
			walkValue(elem, cb)
		}
	case *parser.DictValue:
		for _, entry := range val.Entries {
			walkValue(entry.Key, cb)
			walkValue(entry.Value, cb)
		}
	case *parser.TypedValue:
		for _, arg := range val.Arguments {
			walkValue(arg, cb)
		}
	}
}

// publishGDShaderDiagnostics publishes diagnostics for a GDShader document.
func (s *Server) publishGDShaderDiagnostics(ctx *glsp.Context, uri string, doc *analysis.Document) {
	diagnostics := []protocol.Diagnostic{}

	// Add parse errors from the shader AST
	if doc.ShaderAST != nil {
		for _, err := range doc.ShaderAST.Errors {
			diagnostics = append(diagnostics, protocol.Diagnostic{
				Range: protocol.Range{
					Start: protocol.Position{
						Line:      uint32(err.Range.Start.Line),
						Character: uint32(err.Range.Start.Column),
					},
					End: protocol.Position{
						Line:      uint32(err.Range.End.Line),
						Character: uint32(err.Range.End.Column),
					},
				},
				Severity: severityPtr(protocol.DiagnosticSeverityError),
				Source:   strPtr("gdls"),
				Message:  err.Message,
			})
		}
	}

	// Add semantic errors
	for _, err := range doc.ShaderErrs {
		diagnostics = append(diagnostics, protocol.Diagnostic{
			Range: protocol.Range{
				Start: protocol.Position{
					Line:      uint32(err.Range.Start.Line),
					Character: uint32(err.Range.Start.Column),
				},
				End: protocol.Position{
					Line:      uint32(err.Range.End.Line),
					Character: uint32(err.Range.End.Column),
				},
			},
			Severity: severityPtr(protocol.DiagnosticSeverityError),
			Source:   strPtr("gdls"),
			Message:  err.Message,
		})
	}

	ctx.Notify(protocol.ServerTextDocumentPublishDiagnostics, protocol.PublishDiagnosticsParams{
		URI:         uri,
		Diagnostics: diagnostics,
	})
}

func severityPtr(s protocol.DiagnosticSeverity) *protocol.DiagnosticSeverity {
	return &s
}

func strPtr(s string) *string {
	return &s
}
