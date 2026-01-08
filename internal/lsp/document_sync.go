package lsp

import (
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

// textDocumentDidOpen handles the textDocument/didOpen notification.
func (s *Server) textDocumentDidOpen(ctx *glsp.Context, params *protocol.DidOpenTextDocumentParams) error {
	uri := params.TextDocument.URI
	content := params.TextDocument.Text

	// Store the document and parse it
	doc := s.workspace.OpenDocument(uri, content)

	// Publish diagnostics
	s.publishDiagnostics(ctx, uri, doc)

	return nil
}

// textDocumentDidChange handles the textDocument/didChange notification.
func (s *Server) textDocumentDidChange(ctx *glsp.Context, params *protocol.DidChangeTextDocumentParams) error {
	uri := params.TextDocument.URI

	// With full sync, we get the complete new content
	if len(params.ContentChanges) > 0 {
		// The last change contains the full content in full sync mode
		// ContentChanges is []any in glsp, need to type assert
		lastChange := params.ContentChanges[len(params.ContentChanges)-1]
		if change, ok := lastChange.(protocol.TextDocumentContentChangeEventWhole); ok {
			doc := s.workspace.UpdateDocument(uri, change.Text)
			s.publishDiagnostics(ctx, uri, doc)
		} else if changeMap, ok := lastChange.(map[string]any); ok {
			// Fallback for when it comes as a map
			if text, ok := changeMap["text"].(string); ok {
				doc := s.workspace.UpdateDocument(uri, text)
				s.publishDiagnostics(ctx, uri, doc)
			}
		}
	}

	return nil
}

// textDocumentDidClose handles the textDocument/didClose notification.
func (s *Server) textDocumentDidClose(ctx *glsp.Context, params *protocol.DidCloseTextDocumentParams) error {
	uri := params.TextDocument.URI
	s.workspace.CloseDocument(uri)

	// Clear diagnostics for the closed document
	ctx.Notify(protocol.ServerTextDocumentPublishDiagnostics, protocol.PublishDiagnosticsParams{
		URI:         uri,
		Diagnostics: []protocol.Diagnostic{},
	})

	return nil
}

// textDocumentDidSave handles the textDocument/didSave notification.
func (s *Server) textDocumentDidSave(ctx *glsp.Context, params *protocol.DidSaveTextDocumentParams) error {
	// Re-parse if text is included
	if params.Text != nil {
		uri := params.TextDocument.URI
		doc := s.workspace.UpdateDocument(uri, *params.Text)
		s.publishDiagnostics(ctx, uri, doc)
	}
	return nil
}
