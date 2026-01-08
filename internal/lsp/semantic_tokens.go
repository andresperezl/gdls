package lsp

import (
	"sort"

	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"

	"github.com/andresperezl/gdls/internal/analysis"
	"github.com/andresperezl/gdls/internal/parser"
)

// Semantic token types (must match the order in server.go SemanticTokensLegend)
const (
	tokenTypeKeyword   = 0 // gd_scene, ext_resource, sub_resource, node, connection
	tokenTypeType      = 1 // Vector3, Transform3D, Color, node types
	tokenTypeString    = 2 // string literals
	tokenTypeNumber    = 3 // numeric literals
	tokenTypeProperty  = 4 // property keys
	tokenTypeFunction  = 5 // type constructors
	tokenTypeComment   = 6 // ; comments
	tokenTypeVariable  = 7 // resource IDs
	tokenTypeParameter = 8 // parameters in headings
)

// semanticToken represents a single semantic token.
type semanticToken struct {
	line      uint32
	startChar uint32
	length    uint32
	tokenType uint32
	modifiers uint32
}

// textDocumentSemanticTokensFull handles the textDocument/semanticTokens/full request.
func (s *Server) textDocumentSemanticTokensFull(ctx *glsp.Context, params *protocol.SemanticTokensParams) (*protocol.SemanticTokens, error) {
	doc := s.workspace.GetDocument(params.TextDocument.URI)
	if doc == nil {
		return nil, nil
	}

	var tokens []semanticToken

	switch doc.Type {
	case analysis.DocumentTypeTSCN:
		if doc.TSCNAST == nil {
			return nil, nil
		}
		tokens = s.collectTSCNSemanticTokens(doc.TSCNAST)
	case analysis.DocumentTypeGDShader:
		// For now, GDShader semantic tokens are handled by VS Code's TextMate grammar
		// We can add more advanced semantic highlighting later
		return nil, nil
	default:
		return nil, nil
	}

	// Sort tokens by position
	sort.Slice(tokens, func(i, j int) bool {
		if tokens[i].line != tokens[j].line {
			return tokens[i].line < tokens[j].line
		}
		return tokens[i].startChar < tokens[j].startChar
	})

	// Convert to delta-encoded format
	data := make([]uint32, 0, len(tokens)*5)
	var prevLine, prevChar uint32

	for _, tok := range tokens {
		deltaLine := tok.line - prevLine
		var deltaChar uint32
		if deltaLine == 0 {
			deltaChar = tok.startChar - prevChar
		} else {
			deltaChar = tok.startChar
		}

		data = append(data, deltaLine, deltaChar, tok.length, tok.tokenType, tok.modifiers)
		prevLine = tok.line
		prevChar = tok.startChar
	}

	return &protocol.SemanticTokens{
		Data: data,
	}, nil
}

// collectTSCNSemanticTokens collects all semantic tokens from a TSCN AST.
func (s *Server) collectTSCNSemanticTokens(ast *parser.Document) []semanticToken {
	tokens := []semanticToken{}

	// Descriptor tokens
	if ast.Descriptor != nil {
		// The section keyword
		tokens = append(tokens, semanticToken{
			line:      uint32(ast.Descriptor.Range.Start.Line),
			startChar: uint32(ast.Descriptor.Range.Start.Column) + 1, // Skip '['
			length:    uint32(len(ast.Descriptor.Type)),
			tokenType: tokenTypeKeyword,
		})
	}

	// External resources
	for _, ext := range ast.ExtResources {
		// Section keyword
		tokens = append(tokens, semanticToken{
			line:      uint32(ext.Range.Start.Line),
			startChar: uint32(ext.Range.Start.Column) + 1,
			length:    12, // "ext_resource"
			tokenType: tokenTypeKeyword,
		})
	}

	// Sub resources
	for _, sub := range ast.SubResources {
		// Section keyword
		tokens = append(tokens, semanticToken{
			line:      uint32(sub.Range.Start.Line),
			startChar: uint32(sub.Range.Start.Column) + 1,
			length:    12, // "sub_resource"
			tokenType: tokenTypeKeyword,
		})

		// Properties
		for _, prop := range sub.Properties {
			tokens = append(tokens, s.tokenizeProperty(prop)...)
		}
	}

	// Nodes
	for _, node := range ast.Nodes {
		// Section keyword
		tokens = append(tokens, semanticToken{
			line:      uint32(node.Range.Start.Line),
			startChar: uint32(node.Range.Start.Column) + 1,
			length:    4, // "node"
			tokenType: tokenTypeKeyword,
		})

		// Properties
		for _, prop := range node.Properties {
			tokens = append(tokens, s.tokenizeProperty(prop)...)
		}
	}

	// Connections
	for _, conn := range ast.Connections {
		tokens = append(tokens, semanticToken{
			line:      uint32(conn.Range.Start.Line),
			startChar: uint32(conn.Range.Start.Column) + 1,
			length:    10, // "connection"
			tokenType: tokenTypeKeyword,
		})
	}

	// Comments
	for _, comment := range ast.Comments {
		tokens = append(tokens, semanticToken{
			line:      uint32(comment.Range.Start.Line),
			startChar: uint32(comment.Range.Start.Column),
			length:    uint32(comment.Range.End.Column - comment.Range.Start.Column),
			tokenType: tokenTypeComment,
		})
	}

	return tokens
}

// tokenizeProperty tokenizes a property and its value.
func (s *Server) tokenizeProperty(prop *parser.Property) []semanticToken {
	tokens := []semanticToken{}

	// Property key
	tokens = append(tokens, semanticToken{
		line:      uint32(prop.KeyRange.Start.Line),
		startChar: uint32(prop.KeyRange.Start.Column),
		length:    uint32(prop.KeyRange.End.Column - prop.KeyRange.Start.Column),
		tokenType: tokenTypeProperty,
	})

	// Property value
	tokens = append(tokens, s.tokenizeValue(prop.Value)...)

	return tokens
}

// tokenizeValue tokenizes a value.
func (s *Server) tokenizeValue(v parser.Value) []semanticToken {
	if v == nil {
		return nil
	}

	tokens := []semanticToken{}

	switch val := v.(type) {
	case *parser.StringValue:
		tokens = append(tokens, semanticToken{
			line:      uint32(val.Range.Start.Line),
			startChar: uint32(val.Range.Start.Column),
			length:    uint32(val.Range.End.Column - val.Range.Start.Column),
			tokenType: tokenTypeString,
		})

	case *parser.NumberValue:
		tokens = append(tokens, semanticToken{
			line:      uint32(val.Range.Start.Line),
			startChar: uint32(val.Range.Start.Column),
			length:    uint32(val.Range.End.Column - val.Range.Start.Column),
			tokenType: tokenTypeNumber,
		})

	case *parser.TypedValue:
		// Type name
		tokens = append(tokens, semanticToken{
			line:      uint32(val.TypeRange.Start.Line),
			startChar: uint32(val.TypeRange.Start.Column),
			length:    uint32(val.TypeRange.End.Column - val.TypeRange.Start.Column),
			tokenType: tokenTypeFunction,
		})
		// Arguments
		for _, arg := range val.Arguments {
			tokens = append(tokens, s.tokenizeValue(arg)...)
		}

	case *parser.ResourceRef:
		// Ref type (ExtResource/SubResource)
		tokens = append(tokens, semanticToken{
			line:      uint32(val.Range.Start.Line),
			startChar: uint32(val.Range.Start.Column),
			length:    uint32(len(val.RefType)),
			tokenType: tokenTypeFunction,
		})
		// ID
		tokens = append(tokens, semanticToken{
			line:      uint32(val.IDRange.Start.Line),
			startChar: uint32(val.IDRange.Start.Column),
			length:    uint32(val.IDRange.End.Column - val.IDRange.Start.Column),
			tokenType: tokenTypeVariable,
		})

	case *parser.ArrayValue:
		for _, elem := range val.Values {
			tokens = append(tokens, s.tokenizeValue(elem)...)
		}

	case *parser.DictValue:
		for _, entry := range val.Entries {
			tokens = append(tokens, s.tokenizeValue(entry.Key)...)
			tokens = append(tokens, s.tokenizeValue(entry.Value)...)
		}
	}

	return tokens
}
