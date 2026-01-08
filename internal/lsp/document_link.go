package lsp

import (
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

// textDocumentDocumentLink handles the textDocument/documentLink request.
func (s *Server) textDocumentDocumentLink(ctx *glsp.Context, params *protocol.DocumentLinkParams) ([]protocol.DocumentLink, error) {
	doc := s.workspace.GetDocument(params.TextDocument.URI)
	if doc == nil || doc.TSCNAST == nil {
		return nil, nil
	}

	links := []protocol.DocumentLink{}

	// Add links for external resource paths
	for _, ext := range doc.TSCNAST.ExtResources {
		if ext.Path == "" {
			continue
		}

		// Resolve the path
		location := s.resolveResourcePath(ext.Path, params.TextDocument.URI)
		if location == nil {
			continue
		}

		// Use PathRange if available, otherwise fall back to full range
		linkRange := ext.PathRange
		if linkRange.Start.Line == 0 && linkRange.Start.Column == 0 && linkRange.End.Line == 0 && linkRange.End.Column == 0 {
			// PathRange not set, use full range as fallback
			linkRange = ext.Range
		}

		links = append(links, protocol.DocumentLink{
			Range: protocol.Range{
				Start: protocol.Position{
					Line:      uint32(linkRange.Start.Line),
					Character: uint32(linkRange.Start.Column),
				},
				End: protocol.Position{
					Line:      uint32(linkRange.End.Line),
					Character: uint32(linkRange.End.Column),
				},
			},
			Target:  strPtr(location.URI),
			Tooltip: strPtr("Open " + ext.Path),
		})
	}

	return links, nil
}
