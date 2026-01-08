// Package lsp implements the Language Server Protocol handlers for TSCN files.
package lsp

import (
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
	"github.com/tliron/glsp/server"

	"github.com/andresperezl/gdls/internal/analysis"
)

// Server represents the TSCN language server.
type Server struct {
	name      string
	version   string
	handler   protocol.Handler
	server    *server.Server
	workspace *analysis.Workspace
}

// NewServer creates a new TSCN language server.
func NewServer(name, version string) *Server {
	s := &Server{
		name:      name,
		version:   version,
		workspace: analysis.NewWorkspace(),
	}

	s.handler = protocol.Handler{
		Initialize:                     s.initialize,
		Initialized:                    s.initialized,
		Shutdown:                       s.shutdown,
		SetTrace:                       s.setTrace,
		TextDocumentDidOpen:            s.textDocumentDidOpen,
		TextDocumentDidChange:          s.textDocumentDidChange,
		TextDocumentDidClose:           s.textDocumentDidClose,
		TextDocumentDidSave:            s.textDocumentDidSave,
		TextDocumentHover:              s.textDocumentHover,
		TextDocumentDefinition:         s.textDocumentDefinition,
		TextDocumentDocumentSymbol:     s.textDocumentDocumentSymbol,
		TextDocumentCompletion:         s.textDocumentCompletion,
		TextDocumentFoldingRange:       s.textDocumentFoldingRange,
		TextDocumentDocumentLink:       s.textDocumentDocumentLink,
		TextDocumentReferences:         s.textDocumentReferences,
		TextDocumentSemanticTokensFull: s.textDocumentSemanticTokensFull,
	}

	s.server = server.NewServer(&s.handler, name, false)

	return s
}

// RunStdio runs the server using stdio transport.
func (s *Server) RunStdio() error {
	return s.server.RunStdio()
}

// initialize handles the initialize request from the client.
func (s *Server) initialize(ctx *glsp.Context, params *protocol.InitializeParams) (any, error) {
	capabilities := s.handler.CreateServerCapabilities()

	// Configure text document sync - use full sync for simplicity
	sync := protocol.TextDocumentSyncKindFull
	capabilities.TextDocumentSync = &protocol.TextDocumentSyncOptions{
		OpenClose: boolPtr(true),
		Change:    &sync,
		Save: &protocol.SaveOptions{
			IncludeText: boolPtr(true),
		},
	}

	// Enable hover
	capabilities.HoverProvider = &protocol.HoverOptions{}

	// Enable go to definition
	capabilities.DefinitionProvider = &protocol.DefinitionOptions{}

	// Enable document symbols (outline)
	capabilities.DocumentSymbolProvider = &protocol.DocumentSymbolOptions{}

	// Enable completion
	capabilities.CompletionProvider = &protocol.CompletionOptions{
		TriggerCharacters: []string{"\"", "/", "="},
		ResolveProvider:   boolPtr(false),
	}

	// Enable folding ranges
	capabilities.FoldingRangeProvider = &protocol.FoldingRangeOptions{}

	// Enable document links (clickable paths)
	capabilities.DocumentLinkProvider = &protocol.DocumentLinkOptions{
		ResolveProvider: boolPtr(false),
	}

	// Enable find references
	capabilities.ReferencesProvider = &protocol.ReferenceOptions{}

	// Enable semantic tokens
	capabilities.SemanticTokensProvider = &protocol.SemanticTokensOptions{
		Legend: protocol.SemanticTokensLegend{
			TokenTypes: []string{
				"keyword",   // gd_scene, ext_resource, sub_resource, node, connection
				"type",      // Vector3, Transform3D, Color, node types
				"string",    // string literals
				"number",    // numeric literals
				"property",  // property keys
				"function",  // type constructors
				"comment",   // ; comments
				"variable",  // resource IDs
				"parameter", // parameters in headings
			},
			TokenModifiers: []string{
				"declaration",
				"definition",
				"reference",
			},
		},
		Full: boolPtr(true),
	}

	// Store workspace folders if provided
	if params.WorkspaceFolders != nil {
		for _, folder := range params.WorkspaceFolders {
			s.workspace.AddFolder(folder.URI)
		}
	} else if params.RootURI != nil {
		s.workspace.AddFolder(*params.RootURI)
	}

	return protocol.InitializeResult{
		Capabilities: capabilities,
		ServerInfo: &protocol.InitializeResultServerInfo{
			Name:    s.name,
			Version: &s.version,
		},
	}, nil
}

// initialized handles the initialized notification from the client.
func (s *Server) initialized(ctx *glsp.Context, params *protocol.InitializedParams) error {
	return nil
}

// shutdown handles the shutdown request from the client.
func (s *Server) shutdown(ctx *glsp.Context) error {
	protocol.SetTraceValue(protocol.TraceValueOff)
	return nil
}

// setTrace handles the setTrace notification from the client.
func (s *Server) setTrace(ctx *glsp.Context, params *protocol.SetTraceParams) error {
	protocol.SetTraceValue(params.Value)
	return nil
}

// boolPtr is a helper to create a pointer to a bool.
func boolPtr(b bool) *bool {
	return &b
}
