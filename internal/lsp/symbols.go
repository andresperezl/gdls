package lsp

import (
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"

	"github.com/andresperezl/gdls/internal/analysis"
	"github.com/andresperezl/gdls/internal/parser"
)

// textDocumentDocumentSymbol handles the textDocument/documentSymbol request.
func (s *Server) textDocumentDocumentSymbol(ctx *glsp.Context, params *protocol.DocumentSymbolParams) (any, error) {
	doc := s.workspace.GetDocument(params.TextDocument.URI)
	if doc == nil {
		return nil, nil
	}

	switch doc.Type {
	case analysis.DocumentTypeTSCN:
		return s.tscnDocumentSymbols(doc)
	case analysis.DocumentTypeGDShader:
		return s.gdshaderDocumentSymbols(doc)
	default:
		return nil, nil
	}
}

// tscnDocumentSymbols returns document symbols for a TSCN document.
func (s *Server) tscnDocumentSymbols(doc *analysis.Document) (any, error) {
	if doc.TSCNAST == nil {
		return nil, nil
	}

	symbols := []protocol.DocumentSymbol{}

	// Add nodes as a hierarchical tree
	nodeSymbols := s.buildNodeTree(doc.TSCNAST.Nodes)
	symbols = append(symbols, nodeSymbols...)

	// Add external resources
	if len(doc.TSCNAST.ExtResources) > 0 {
		extChildren := []protocol.DocumentSymbol{}
		for _, ext := range doc.TSCNAST.ExtResources {
			extChildren = append(extChildren, protocol.DocumentSymbol{
				Name:   ext.ID,
				Detail: strPtr(ext.Type + " - " + ext.Path),
				Kind:   protocol.SymbolKindFile,
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
				SelectionRange: protocol.Range{
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
		symbols = append(symbols, protocol.DocumentSymbol{
			Name:           "External Resources",
			Kind:           protocol.SymbolKindNamespace,
			Range:          extChildren[0].Range,
			SelectionRange: extChildren[0].Range,
			Children:       extChildren,
		})
	}

	// Add sub resources
	if len(doc.TSCNAST.SubResources) > 0 {
		subChildren := []protocol.DocumentSymbol{}
		for _, sub := range doc.TSCNAST.SubResources {
			subChildren = append(subChildren, protocol.DocumentSymbol{
				Name:   sub.ID,
				Detail: strPtr(sub.Type),
				Kind:   protocol.SymbolKindObject,
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
				SelectionRange: protocol.Range{
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
		symbols = append(symbols, protocol.DocumentSymbol{
			Name:           "Sub Resources",
			Kind:           protocol.SymbolKindNamespace,
			Range:          subChildren[0].Range,
			SelectionRange: subChildren[0].Range,
			Children:       subChildren,
		})
	}

	// Add connections
	if len(doc.TSCNAST.Connections) > 0 {
		connChildren := []protocol.DocumentSymbol{}
		for _, conn := range doc.TSCNAST.Connections {
			connChildren = append(connChildren, protocol.DocumentSymbol{
				Name:   conn.Signal + " -> " + conn.Method,
				Detail: strPtr(conn.From + " → " + conn.To),
				Kind:   protocol.SymbolKindEvent,
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
				SelectionRange: protocol.Range{
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
		symbols = append(symbols, protocol.DocumentSymbol{
			Name:           "Connections",
			Kind:           protocol.SymbolKindNamespace,
			Range:          connChildren[0].Range,
			SelectionRange: connChildren[0].Range,
			Children:       connChildren,
		})
	}

	return symbols, nil
}

// buildNodeTree builds a hierarchical tree of node symbols.
func (s *Server) buildNodeTree(nodes []*parser.Node) []protocol.DocumentSymbol {
	if len(nodes) == 0 {
		return nil
	}

	// Build a map of node paths to their children
	type nodeInfo struct {
		node     *parser.Node
		children []*parser.Node
	}

	pathToNode := make(map[string]*nodeInfo)
	var rootNode *parser.Node

	// First pass: create nodeInfo for each node
	for _, node := range nodes {
		var path string
		switch node.Parent {
		case "":
			path = ""
			rootNode = node
		case ".":
			path = node.Name
		default:
			path = node.Parent + "/" + node.Name
		}
		pathToNode[path] = &nodeInfo{node: node, children: []*parser.Node{}}
	}

	// Second pass: link children to parents
	for _, node := range nodes {
		if node.Parent == "" {
			continue // Root has no parent
		}
		parentPath := node.Parent
		if parentPath == "." {
			parentPath = ""
		}
		if parent, ok := pathToNode[parentPath]; ok {
			parent.children = append(parent.children, node)
		}
	}

	// Build symbols recursively
	var buildSymbol func(node *parser.Node, path string) protocol.DocumentSymbol
	buildSymbol = func(node *parser.Node, path string) protocol.DocumentSymbol {
		detail := node.Type
		if detail == "" {
			detail = "(instance)"
		}

		sym := protocol.DocumentSymbol{
			Name:   node.Name,
			Detail: strPtr(detail),
			Kind:   getNodeSymbolKind(node.Type),
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
			SelectionRange: protocol.Range{
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

		// Add children
		if info, ok := pathToNode[path]; ok {
			for _, child := range info.children {
				var childPath string
				if path == "" {
					childPath = child.Name
				} else {
					childPath = path + "/" + child.Name
				}
				sym.Children = append(sym.Children, buildSymbol(child, childPath))
			}
		}

		return sym
	}

	if rootNode == nil {
		return nil
	}

	return []protocol.DocumentSymbol{buildSymbol(rootNode, "")}
}

// getNodeSymbolKind returns the appropriate symbol kind for a node type.
func getNodeSymbolKind(nodeType string) protocol.SymbolKind {
	switch nodeType {
	case "Node", "Node2D", "Node3D":
		return protocol.SymbolKindClass
	case "Control", "Label", "Button", "TextEdit", "LineEdit":
		return protocol.SymbolKindInterface
	case "Camera2D", "Camera3D":
		return protocol.SymbolKindFunction
	case "Sprite2D", "Sprite3D", "MeshInstance3D":
		return protocol.SymbolKindField
	case "CollisionShape2D", "CollisionShape3D":
		return protocol.SymbolKindStruct
	case "AnimationPlayer", "AnimationTree":
		return protocol.SymbolKindMethod
	case "Timer":
		return protocol.SymbolKindEvent
	case "":
		return protocol.SymbolKindModule // Instance
	default:
		return protocol.SymbolKindClass
	}
}

// gdshaderDocumentSymbols returns document symbols for a GDShader document.
func (s *Server) gdshaderDocumentSymbols(doc *analysis.Document) (any, error) {
	if doc.ShaderAST == nil {
		return nil, nil
	}

	symbols := []protocol.DocumentSymbol{}
	ast := doc.ShaderAST

	// Add uniforms
	for _, uniform := range ast.Uniforms {
		symbols = append(symbols, protocol.DocumentSymbol{
			Name:   uniform.Name,
			Detail: strPtr("uniform " + uniform.Type.Name),
			Kind:   protocol.SymbolKindVariable,
			Range: protocol.Range{
				Start: protocol.Position{
					Line:      uint32(uniform.Range.Start.Line),
					Character: uint32(uniform.Range.Start.Column),
				},
				End: protocol.Position{
					Line:      uint32(uniform.Range.End.Line),
					Character: uint32(uniform.Range.End.Column),
				},
			},
			SelectionRange: protocol.Range{
				Start: protocol.Position{
					Line:      uint32(uniform.Range.Start.Line),
					Character: uint32(uniform.Range.Start.Column),
				},
				End: protocol.Position{
					Line:      uint32(uniform.Range.End.Line),
					Character: uint32(uniform.Range.End.Column),
				},
			},
		})
	}

	// Add varyings
	for _, varying := range ast.Varyings {
		symbols = append(symbols, protocol.DocumentSymbol{
			Name:   varying.Name,
			Detail: strPtr("varying " + varying.Type.Name),
			Kind:   protocol.SymbolKindVariable,
			Range: protocol.Range{
				Start: protocol.Position{
					Line:      uint32(varying.Range.Start.Line),
					Character: uint32(varying.Range.Start.Column),
				},
				End: protocol.Position{
					Line:      uint32(varying.Range.End.Line),
					Character: uint32(varying.Range.End.Column),
				},
			},
			SelectionRange: protocol.Range{
				Start: protocol.Position{
					Line:      uint32(varying.Range.Start.Line),
					Character: uint32(varying.Range.Start.Column),
				},
				End: protocol.Position{
					Line:      uint32(varying.Range.End.Line),
					Character: uint32(varying.Range.End.Column),
				},
			},
		})
	}

	// Add constants
	for _, constant := range ast.Constants {
		symbols = append(symbols, protocol.DocumentSymbol{
			Name:   constant.Name,
			Detail: strPtr("const " + constant.Type.Name),
			Kind:   protocol.SymbolKindConstant,
			Range: protocol.Range{
				Start: protocol.Position{
					Line:      uint32(constant.Range.Start.Line),
					Character: uint32(constant.Range.Start.Column),
				},
				End: protocol.Position{
					Line:      uint32(constant.Range.End.Line),
					Character: uint32(constant.Range.End.Column),
				},
			},
			SelectionRange: protocol.Range{
				Start: protocol.Position{
					Line:      uint32(constant.Range.Start.Line),
					Character: uint32(constant.Range.Start.Column),
				},
				End: protocol.Position{
					Line:      uint32(constant.Range.End.Line),
					Character: uint32(constant.Range.End.Column),
				},
			},
		})
	}

	// Add structs
	for _, st := range ast.Structs {
		structSym := protocol.DocumentSymbol{
			Name:   st.Name,
			Detail: strPtr("struct"),
			Kind:   protocol.SymbolKindStruct,
			Range: protocol.Range{
				Start: protocol.Position{
					Line:      uint32(st.Range.Start.Line),
					Character: uint32(st.Range.Start.Column),
				},
				End: protocol.Position{
					Line:      uint32(st.Range.End.Line),
					Character: uint32(st.Range.End.Column),
				},
			},
			SelectionRange: protocol.Range{
				Start: protocol.Position{
					Line:      uint32(st.Range.Start.Line),
					Character: uint32(st.Range.Start.Column),
				},
				End: protocol.Position{
					Line:      uint32(st.Range.End.Line),
					Character: uint32(st.Range.End.Column),
				},
			},
		}

		// Add struct members as children
		for _, member := range st.Members {
			structSym.Children = append(structSym.Children, protocol.DocumentSymbol{
				Name:   member.Name,
				Detail: strPtr(member.Type.Name),
				Kind:   protocol.SymbolKindField,
				Range: protocol.Range{
					Start: protocol.Position{
						Line:      uint32(member.Range.Start.Line),
						Character: uint32(member.Range.Start.Column),
					},
					End: protocol.Position{
						Line:      uint32(member.Range.End.Line),
						Character: uint32(member.Range.End.Column),
					},
				},
				SelectionRange: protocol.Range{
					Start: protocol.Position{
						Line:      uint32(member.Range.Start.Line),
						Character: uint32(member.Range.Start.Column),
					},
					End: protocol.Position{
						Line:      uint32(member.Range.End.Line),
						Character: uint32(member.Range.End.Column),
					},
				},
			})
		}

		symbols = append(symbols, structSym)
	}

	// Add functions
	for _, fn := range ast.Functions {
		funcSym := protocol.DocumentSymbol{
			Name:   fn.Name,
			Detail: strPtr(fn.ReturnType.Name),
			Kind:   protocol.SymbolKindFunction,
			Range: protocol.Range{
				Start: protocol.Position{
					Line:      uint32(fn.Range.Start.Line),
					Character: uint32(fn.Range.Start.Column),
				},
				End: protocol.Position{
					Line:      uint32(fn.Range.End.Line),
					Character: uint32(fn.Range.End.Column),
				},
			},
			SelectionRange: protocol.Range{
				Start: protocol.Position{
					Line:      uint32(fn.Range.Start.Line),
					Character: uint32(fn.Range.Start.Column),
				},
				End: protocol.Position{
					Line:      uint32(fn.Range.End.Line),
					Character: uint32(fn.Range.End.Column),
				},
			},
		}

		symbols = append(symbols, funcSym)
	}

	return symbols, nil
}
