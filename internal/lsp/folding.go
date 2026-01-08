package lsp

import (
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

// textDocumentFoldingRange handles the textDocument/foldingRange request.
func (s *Server) textDocumentFoldingRange(ctx *glsp.Context, params *protocol.FoldingRangeParams) ([]protocol.FoldingRange, error) {
	doc := s.workspace.GetDocument(params.TextDocument.URI)
	if doc == nil || doc.TSCNAST == nil {
		return nil, nil
	}

	ranges := []protocol.FoldingRange{}
	regionKind := string(protocol.FoldingRangeKindRegion)

	// Add folding ranges for sub_resources with properties
	for _, sub := range doc.TSCNAST.SubResources {
		if len(sub.Properties) > 0 {
			startLine := uint32(sub.Range.Start.Line)
			endLine := uint32(sub.Range.End.Line)
			if endLine > startLine {
				ranges = append(ranges, protocol.FoldingRange{
					StartLine: startLine,
					EndLine:   endLine,
					Kind:      &regionKind,
				})
			}
		}
	}

	// Add folding ranges for nodes with properties
	for _, node := range doc.TSCNAST.Nodes {
		if len(node.Properties) > 0 {
			startLine := uint32(node.Range.Start.Line)
			endLine := uint32(node.Range.End.Line)
			if endLine > startLine {
				ranges = append(ranges, protocol.FoldingRange{
					StartLine: startLine,
					EndLine:   endLine,
					Kind:      &regionKind,
				})
			}
		}
	}

	return ranges, nil
}
