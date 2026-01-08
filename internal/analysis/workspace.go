// Package analysis provides document management and semantic analysis for TSCN and GDShader files.
package analysis

import (
	"strings"
	"sync"

	"github.com/andresperezl/gdls/internal/gdshader"
	"github.com/andresperezl/gdls/internal/parser"
)

// DocumentType represents the type of document.
type DocumentType int

const (
	DocumentTypeTSCN DocumentType = iota
	DocumentTypeGDShader
	DocumentTypeUnknown
)

// Workspace manages all open documents and workspace folders.
type Workspace struct {
	mu        sync.RWMutex
	documents map[string]*Document
	folders   []string
}

// Document represents an open document with its parsed AST.
type Document struct {
	URI        string
	Content    string
	Type       DocumentType
	TSCNAST    *parser.Document          // For TSCN/ESCN files
	ShaderAST  *gdshader.ShaderDocument  // For GDShader files
	ShaderErrs []*gdshader.SemanticError // Semantic errors from GDShader analysis
	Version    int
}

// NewWorkspace creates a new workspace.
func NewWorkspace() *Workspace {
	return &Workspace{
		documents: make(map[string]*Document),
		folders:   []string{},
	}
}

// AddFolder adds a workspace folder.
func (w *Workspace) AddFolder(uri string) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.folders = append(w.folders, uri)
}

// GetFolders returns all workspace folders.
func (w *Workspace) GetFolders() []string {
	w.mu.RLock()
	defer w.mu.RUnlock()
	// Return a copy to avoid race conditions
	folders := make([]string, len(w.folders))
	copy(folders, w.folders)
	return folders
}

// GetDocumentType determines the document type from URI.
func GetDocumentType(uri string) DocumentType {
	lowerURI := strings.ToLower(uri)
	if strings.HasSuffix(lowerURI, ".tscn") || strings.HasSuffix(lowerURI, ".escn") {
		return DocumentTypeTSCN
	}
	if strings.HasSuffix(lowerURI, ".gdshader") || strings.HasSuffix(lowerURI, ".gdshaderinc") {
		return DocumentTypeGDShader
	}
	return DocumentTypeUnknown
}

// OpenDocument opens a document and parses it.
func (w *Workspace) OpenDocument(uri, content string) *Document {
	w.mu.Lock()
	defer w.mu.Unlock()

	doc := parseDocument(uri, content)
	doc.Version = 1
	w.documents[uri] = doc
	return doc
}

// UpdateDocument updates a document's content and re-parses it.
func (w *Workspace) UpdateDocument(uri, content string) *Document {
	w.mu.Lock()
	defer w.mu.Unlock()

	existingDoc, exists := w.documents[uri]
	doc := parseDocument(uri, content)

	if exists {
		doc.Version = existingDoc.Version + 1
	} else {
		doc.Version = 1
	}
	w.documents[uri] = doc
	return doc
}

// parseDocument parses a document based on its type.
func parseDocument(uri, content string) *Document {
	docType := GetDocumentType(uri)
	doc := &Document{
		URI:     uri,
		Content: content,
		Type:    docType,
	}

	switch docType {
	case DocumentTypeTSCN:
		doc.TSCNAST = parser.Parse(content)
	case DocumentTypeGDShader:
		p := gdshader.NewParser(content)
		doc.ShaderAST = p.Parse()
		// Run semantic analysis
		if doc.ShaderAST != nil {
			analyzer := gdshader.NewAnalyzer(doc.ShaderAST)
			doc.ShaderErrs = analyzer.Analyze()
		}
	}

	return doc
}

// CloseDocument closes a document.
func (w *Workspace) CloseDocument(uri string) {
	w.mu.Lock()
	defer w.mu.Unlock()
	delete(w.documents, uri)
}

// GetDocument returns a document by URI.
func (w *Workspace) GetDocument(uri string) *Document {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.documents[uri]
}

// GetAllDocuments returns all open documents.
func (w *Workspace) GetAllDocuments() []*Document {
	w.mu.RLock()
	defer w.mu.RUnlock()

	docs := make([]*Document, 0, len(w.documents))
	for _, doc := range w.documents {
		docs = append(docs, doc)
	}
	return docs
}

// PositionToOffset converts a line/character position to a byte offset.
func (d *Document) PositionToOffset(line, character uint32) int {
	lines := strings.Split(d.Content, "\n")
	offset := 0

	for i := uint32(0); i < line && int(i) < len(lines); i++ {
		offset += len(lines[i]) + 1 // +1 for newline
	}

	if int(line) < len(lines) {
		lineContent := lines[line]
		if int(character) <= len(lineContent) {
			offset += int(character)
		} else {
			offset += len(lineContent)
		}
	}

	return offset
}

// OffsetToPosition converts a byte offset to a line/character position.
func (d *Document) OffsetToPosition(offset int) (line, character uint32) {
	content := d.Content
	if offset > len(content) {
		offset = len(content)
	}

	line = 0
	lineStart := 0

	for i := 0; i < offset; i++ {
		if content[i] == '\n' {
			line++
			lineStart = i + 1
		}
	}

	character = uint32(offset - lineStart)
	return
}
