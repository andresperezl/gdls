// Package e2e provides end-to-end tests for the Godot language server (gdls).
package e2e

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
)

const (
	requestTimeout = 5 * time.Second
)

// =============================================================================
// JSON-RPC Types
// =============================================================================

type jsonrpcRequest struct {
	JSONRPC string `json:"jsonrpc"`
	ID      int    `json:"id"`
	Method  string `json:"method"`
	Params  any    `json:"params,omitempty"`
}

type jsonrpcNotification struct {
	JSONRPC string `json:"jsonrpc"`
	Method  string `json:"method"`
	Params  any    `json:"params,omitempty"`
}

type jsonrpcResponse struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      *int            `json:"id,omitempty"`
	Method  string          `json:"method,omitempty"`
	Result  json.RawMessage `json:"result,omitempty"`
	Params  json.RawMessage `json:"params,omitempty"`
	Error   *jsonrpcError   `json:"error,omitempty"`
}

type jsonrpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// =============================================================================
// LSP Types (minimal, for test verification)
// =============================================================================

type initializeParams struct {
	ProcessID    int                `json:"processId"`
	Capabilities clientCapabilities `json:"capabilities"`
	RootURI      *string            `json:"rootUri"`
	Trace        string             `json:"trace,omitempty"`
}

type clientCapabilities struct {
	TextDocument textDocumentClientCapabilities `json:"textDocument,omitempty"`
}

type textDocumentClientCapabilities struct {
	Synchronization    textDocumentSyncClientCapabilities   `json:"synchronization,omitempty"`
	PublishDiagnostics publishDiagnosticsClientCapabilities `json:"publishDiagnostics,omitempty"`
}

type textDocumentSyncClientCapabilities struct {
	DidSave bool `json:"didSave,omitempty"`
}

type publishDiagnosticsClientCapabilities struct {
	RelatedInformation bool `json:"relatedInformation,omitempty"`
}

type initializeResult struct {
	Capabilities serverCapabilities `json:"capabilities"`
	ServerInfo   *serverInfo        `json:"serverInfo,omitempty"`
}

type serverCapabilities struct {
	TextDocumentSync       any `json:"textDocumentSync,omitempty"`
	HoverProvider          any `json:"hoverProvider,omitempty"`
	DocumentSymbolProvider any `json:"documentSymbolProvider,omitempty"`
	CompletionProvider     any `json:"completionProvider,omitempty"`
	DefinitionProvider     any `json:"definitionProvider,omitempty"`
	ReferencesProvider     any `json:"referencesProvider,omitempty"`
	SemanticTokensProvider any `json:"semanticTokensProvider,omitempty"`
}

type serverInfo struct {
	Name    string  `json:"name"`
	Version *string `json:"version,omitempty"`
}

type textDocumentItem struct {
	URI        string `json:"uri"`
	LanguageID string `json:"languageId"`
	Version    int    `json:"version"`
	Text       string `json:"text"`
}

type textDocumentIdentifier struct {
	URI string `json:"uri"`
}

type didOpenTextDocumentParams struct {
	TextDocument textDocumentItem `json:"textDocument"`
}

type didCloseTextDocumentParams struct {
	TextDocument textDocumentIdentifier `json:"textDocument"`
}

type documentSymbolParams struct {
	TextDocument textDocumentIdentifier `json:"textDocument"`
}

type documentSymbol struct {
	Name           string           `json:"name"`
	Kind           int              `json:"kind"`
	Range          lspRange         `json:"range"`
	SelectionRange lspRange         `json:"selectionRange"`
	Children       []documentSymbol `json:"children,omitempty"`
}

type symbolInformation struct {
	Name     string   `json:"name"`
	Kind     int      `json:"kind"`
	Location location `json:"location"`
}

type location struct {
	URI   string   `json:"uri"`
	Range lspRange `json:"range"`
}

type lspRange struct {
	Start position `json:"start"`
	End   position `json:"end"`
}

type position struct {
	Line      int `json:"line"`
	Character int `json:"character"`
}

type hoverParams struct {
	TextDocument textDocumentIdentifier `json:"textDocument"`
	Position     position               `json:"position"`
}

type hoverResult struct {
	Contents any       `json:"contents"`
	Range    *lspRange `json:"range,omitempty"`
}

type semanticTokensParams struct {
	TextDocument textDocumentIdentifier `json:"textDocument"`
}

type semanticTokensResult struct {
	Data []uint32 `json:"data"`
}

type publishDiagnosticsParams struct {
	URI         string       `json:"uri"`
	Diagnostics []diagnostic `json:"diagnostics"`
}

type diagnostic struct {
	Range    lspRange `json:"range"`
	Message  string   `json:"message"`
	Severity *int     `json:"severity,omitempty"`
}

// =============================================================================
// Test LSP Client
// =============================================================================

type testLSPClient struct {
	t      *testing.T
	cmd    *exec.Cmd
	stdin  io.WriteCloser
	stdout *bufio.Reader
	stderr io.ReadCloser
	nextID int
	mu     sync.Mutex

	// For async notification handling
	notifications chan jsonrpcResponse
	responses     map[int]chan jsonrpcResponse
	responsesMu   sync.Mutex
	done          chan struct{}
	readerStarted bool
}

func newTestLSPClient(t *testing.T) *testLSPClient {
	t.Helper()

	// Find the project root (where go.mod is)
	projectRoot, err := findProjectRoot()
	if err != nil {
		t.Fatalf("failed to find project root: %v", err)
	}

	// Start the server using go run
	cmd := exec.Command("go", "run", "./cmd/gdls")
	cmd.Dir = projectRoot

	stdin, err := cmd.StdinPipe()
	if err != nil {
		t.Fatalf("failed to create stdin pipe: %v", err)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		t.Fatalf("failed to create stdout pipe: %v", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		t.Fatalf("failed to create stderr pipe: %v", err)
	}

	if err := cmd.Start(); err != nil {
		t.Fatalf("failed to start server: %v", err)
	}

	client := &testLSPClient{
		t:             t,
		cmd:           cmd,
		stdin:         stdin,
		stdout:        bufio.NewReader(stdout),
		stderr:        stderr,
		nextID:        1,
		notifications: make(chan jsonrpcResponse, 100),
		responses:     make(map[int]chan jsonrpcResponse),
		done:          make(chan struct{}),
	}

	// Start the async reader
	client.startReader()

	return client
}

func findProjectRoot() (string, error) {
	// Start from the current working directory and go up
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("could not find go.mod in any parent directory")
		}
		dir = parent
	}
}

func (c *testLSPClient) startReader() {
	c.readerStarted = true
	go func() {
		for {
			select {
			case <-c.done:
				return
			default:
				msg, err := c.readMessage()
				if err != nil {
					// Connection closed or error
					return
				}

				var resp jsonrpcResponse
				if err := json.Unmarshal(msg, &resp); err != nil {
					c.t.Logf("failed to unmarshal response: %v", err)
					continue
				}

				if resp.ID != nil {
					// It's a response to a request
					c.responsesMu.Lock()
					if ch, ok := c.responses[*resp.ID]; ok {
						ch <- resp
						delete(c.responses, *resp.ID)
					}
					c.responsesMu.Unlock()
				} else if resp.Method != "" {
					// It's a notification from server
					select {
					case c.notifications <- resp:
					default:
						// Channel full, drop notification
					}
				}
			}
		}
	}()
}

func (c *testLSPClient) writeMessage(msg []byte) error {
	header := fmt.Sprintf("Content-Length: %d\r\n\r\n", len(msg))
	if _, err := c.stdin.Write([]byte(header)); err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}
	if _, err := c.stdin.Write(msg); err != nil {
		return fmt.Errorf("failed to write body: %w", err)
	}
	return nil
}

func (c *testLSPClient) readMessage() ([]byte, error) {
	// Read headers until empty line
	var contentLength int
	for {
		line, err := c.stdout.ReadString('\n')
		if err != nil {
			return nil, fmt.Errorf("failed to read header line: %w", err)
		}
		line = strings.TrimSpace(line)
		if line == "" {
			break // End of headers
		}
		if strings.HasPrefix(line, "Content-Length:") {
			lenStr := strings.TrimSpace(strings.TrimPrefix(line, "Content-Length:"))
			contentLength, err = strconv.Atoi(lenStr)
			if err != nil {
				return nil, fmt.Errorf("invalid Content-Length: %w", err)
			}
		}
	}

	if contentLength == 0 {
		return nil, fmt.Errorf("no Content-Length header")
	}

	// Read exactly contentLength bytes
	body := make([]byte, contentLength)
	if _, err := io.ReadFull(c.stdout, body); err != nil {
		return nil, fmt.Errorf("failed to read body: %w", err)
	}

	return body, nil
}

func (c *testLSPClient) sendRequest(ctx context.Context, method string, params any) (json.RawMessage, error) {
	c.mu.Lock()
	id := c.nextID
	c.nextID++
	c.mu.Unlock()

	req := jsonrpcRequest{
		JSONRPC: "2.0",
		ID:      id,
		Method:  method,
		Params:  params,
	}

	msg, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create response channel
	respCh := make(chan jsonrpcResponse, 1)
	c.responsesMu.Lock()
	c.responses[id] = respCh
	c.responsesMu.Unlock()

	// Send the request
	if err := c.writeMessage(msg); err != nil {
		c.responsesMu.Lock()
		delete(c.responses, id)
		c.responsesMu.Unlock()
		return nil, err
	}

	// Wait for response with timeout
	select {
	case <-ctx.Done():
		c.responsesMu.Lock()
		delete(c.responses, id)
		c.responsesMu.Unlock()
		return nil, ctx.Err()
	case resp := <-respCh:
		if resp.Error != nil {
			return nil, fmt.Errorf("LSP error %d: %s", resp.Error.Code, resp.Error.Message)
		}
		return resp.Result, nil
	}
}

func (c *testLSPClient) sendNotification(method string, params any) error {
	notif := jsonrpcNotification{
		JSONRPC: "2.0",
		Method:  method,
		Params:  params,
	}

	msg, err := json.Marshal(notif)
	if err != nil {
		return fmt.Errorf("failed to marshal notification: %w", err)
	}

	return c.writeMessage(msg)
}

func (c *testLSPClient) waitForNotification(ctx context.Context, method string) (json.RawMessage, error) {
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case notif := <-c.notifications:
			if notif.Method == method {
				return notif.Params, nil
			}
			// Not the notification we're looking for, continue waiting
		}
	}
}

func (c *testLSPClient) initialize(ctx context.Context) (*initializeResult, error) {
	params := initializeParams{
		ProcessID: os.Getpid(),
		Capabilities: clientCapabilities{
			TextDocument: textDocumentClientCapabilities{
				Synchronization: textDocumentSyncClientCapabilities{
					DidSave: true,
				},
				PublishDiagnostics: publishDiagnosticsClientCapabilities{
					RelatedInformation: true,
				},
			},
		},
		RootURI: nil,
		Trace:   "off",
	}

	result, err := c.sendRequest(ctx, "initialize", params)
	if err != nil {
		return nil, err
	}

	var initResult initializeResult
	if err := json.Unmarshal(result, &initResult); err != nil {
		return nil, fmt.Errorf("failed to unmarshal initialize result: %w", err)
	}

	// Send initialized notification
	if err := c.sendNotification("initialized", struct{}{}); err != nil {
		return nil, fmt.Errorf("failed to send initialized: %w", err)
	}

	return &initResult, nil
}

func (c *testLSPClient) openDocument(uri, content string) error {
	params := didOpenTextDocumentParams{
		TextDocument: textDocumentItem{
			URI:        uri,
			LanguageID: "tscn",
			Version:    1,
			Text:       content,
		},
	}

	return c.sendNotification("textDocument/didOpen", params)
}

func (c *testLSPClient) closeDocument(uri string) error {
	params := didCloseTextDocumentParams{
		TextDocument: textDocumentIdentifier{
			URI: uri,
		},
	}

	return c.sendNotification("textDocument/didClose", params)
}

func (c *testLSPClient) shutdown(ctx context.Context) error {
	_, err := c.sendRequest(ctx, "shutdown", nil)
	return err
}

func (c *testLSPClient) exit() error {
	return c.sendNotification("exit", nil)
}

func (c *testLSPClient) close() int {
	// Signal reader to stop
	close(c.done)

	// Close stdin to signal EOF to the server
	c.stdin.Close()

	// Wait for the process to exit
	err := c.cmd.Wait()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return exitErr.ExitCode()
		}
		return -1
	}
	return 0
}

// =============================================================================
// Helper Functions
// =============================================================================

func loadTestFile(t *testing.T, filename string) string {
	t.Helper()

	projectRoot, err := findProjectRoot()
	if err != nil {
		t.Fatalf("failed to find project root: %v", err)
	}

	path := filepath.Join(projectRoot, "testdata", filename)
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read test file %s: %v", filename, err)
	}

	return string(content)
}

// =============================================================================
// Tests
// =============================================================================

func TestLSPInitialization(t *testing.T) {
	t.Parallel()

	client := newTestLSPClient(t)
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
		defer cancel()
		client.shutdown(ctx)
		client.exit()
		client.close()
	}()

	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	result, err := client.initialize(ctx)
	if err != nil {
		t.Fatalf("initialize failed: %v", err)
	}

	// Verify server info
	if result.ServerInfo == nil {
		t.Error("expected ServerInfo to be present")
	} else {
		if result.ServerInfo.Name != "gdls" {
			t.Errorf("expected server name 'gdls', got '%s'", result.ServerInfo.Name)
		}
	}

	// Verify capabilities
	if result.Capabilities.TextDocumentSync == nil {
		t.Error("expected TextDocumentSync capability")
	}
	if result.Capabilities.HoverProvider == nil {
		t.Error("expected HoverProvider capability")
	}
	if result.Capabilities.DocumentSymbolProvider == nil {
		t.Error("expected DocumentSymbolProvider capability")
	}
	if result.Capabilities.SemanticTokensProvider == nil {
		t.Error("expected SemanticTokensProvider capability")
	}
}

func TestLSPDoesNotCrashOnSimpleTSCN(t *testing.T) {
	t.Parallel()
	testFileDoesNotCrash(t, "simple.tscn")
}

func TestLSPDoesNotCrashOnComplexTSCN(t *testing.T) {
	t.Parallel()
	testFileDoesNotCrash(t, "complex.tscn")
}

func TestLSPDoesNotCrashOnErrorsTSCN(t *testing.T) {
	t.Parallel()
	testFileDoesNotCrash(t, "with_errors.tscn")
}

func testFileDoesNotCrash(t *testing.T, filename string) {
	t.Helper()

	client := newTestLSPClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	// Initialize
	_, err := client.initialize(ctx)
	if err != nil {
		t.Fatalf("initialize failed: %v", err)
	}

	// Open document
	content := loadTestFile(t, filename)
	uri := "file:///test/" + filename
	if err := client.openDocument(uri, content); err != nil {
		t.Fatalf("failed to open document: %v", err)
	}

	// Give server time to process
	time.Sleep(100 * time.Millisecond)

	// Shutdown gracefully
	if err := client.shutdown(ctx); err != nil {
		t.Fatalf("shutdown failed: %v", err)
	}

	// Exit
	if err := client.exit(); err != nil {
		t.Fatalf("exit notification failed: %v", err)
	}

	// Verify exit code is 0
	exitCode := client.close()
	if exitCode != 0 {
		t.Errorf("expected exit code 0, got %d", exitCode)
	}
}

func TestLSPDocumentSymbols(t *testing.T) {
	t.Parallel()

	client := newTestLSPClient(t)
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
		defer cancel()
		client.shutdown(ctx)
		client.exit()
		client.close()
	}()

	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	_, err := client.initialize(ctx)
	if err != nil {
		t.Fatalf("initialize failed: %v", err)
	}

	// Open simple.tscn (has Player and CollisionShape3D nodes)
	content := loadTestFile(t, "simple.tscn")
	uri := "file:///test/simple.tscn"
	if err := client.openDocument(uri, content); err != nil {
		t.Fatalf("failed to open document: %v", err)
	}

	// Give server time to process
	time.Sleep(50 * time.Millisecond)

	// Request document symbols
	params := documentSymbolParams{
		TextDocument: textDocumentIdentifier{URI: uri},
	}
	result, err := client.sendRequest(ctx, "textDocument/documentSymbol", params)
	if err != nil {
		t.Fatalf("documentSymbol request failed: %v", err)
	}

	// Try to parse as DocumentSymbol[] first
	var symbols []documentSymbol
	if err := json.Unmarshal(result, &symbols); err != nil {
		// Try SymbolInformation[] (older format)
		var symbolInfos []symbolInformation
		if err2 := json.Unmarshal(result, &symbolInfos); err2 != nil {
			t.Fatalf("failed to unmarshal symbols: %v (also tried: %v)", err, err2)
		}
		// Convert to check
		if len(symbolInfos) < 2 {
			t.Errorf("expected at least 2 symbols (Player, CollisionShape3D), got %d", len(symbolInfos))
		}

		// Find Player node
		found := false
		for _, sym := range symbolInfos {
			if sym.Name == "Player" {
				found = true
				break
			}
		}
		if !found {
			t.Error("expected to find 'Player' symbol")
		}
		return
	}

	// Verify we got symbols for the nodes
	if len(symbols) < 2 {
		t.Errorf("expected at least 2 symbols (Player, CollisionShape3D), got %d", len(symbols))
	}

	// Find Player node
	found := false
	for _, sym := range symbols {
		if sym.Name == "Player" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected to find 'Player' symbol")
	}
}

func TestLSPHover(t *testing.T) {
	t.Parallel()

	client := newTestLSPClient(t)
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
		defer cancel()
		client.shutdown(ctx)
		client.exit()
		client.close()
	}()

	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	_, err := client.initialize(ctx)
	if err != nil {
		t.Fatalf("initialize failed: %v", err)
	}

	content := loadTestFile(t, "simple.tscn")
	uri := "file:///test/simple.tscn"
	if err := client.openDocument(uri, content); err != nil {
		t.Fatalf("failed to open document: %v", err)
	}

	// Give server time to process
	time.Sleep(50 * time.Millisecond)

	// Hover over line 5 (the [node name="Player"...] line) - 0-indexed, so line 4
	params := hoverParams{
		TextDocument: textDocumentIdentifier{URI: uri},
		Position:     position{Line: 4, Character: 10},
	}
	result, err := client.sendRequest(ctx, "textDocument/hover", params)
	if err != nil {
		t.Fatalf("hover request failed: %v", err)
	}

	// Result may be null if no hover info at that position, but should not error
	if string(result) != "null" {
		var hover hoverResult
		if err := json.Unmarshal(result, &hover); err != nil {
			t.Fatalf("failed to unmarshal hover result: %v", err)
		}
		// Verify contents is present
		if hover.Contents == nil {
			t.Error("expected hover contents to be present")
		}
	}
	// If null, that's OK - just means no hover info at that position
}

func TestLSPSemanticTokens(t *testing.T) {
	t.Parallel()

	client := newTestLSPClient(t)
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
		defer cancel()
		client.shutdown(ctx)
		client.exit()
		client.close()
	}()

	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	_, err := client.initialize(ctx)
	if err != nil {
		t.Fatalf("initialize failed: %v", err)
	}

	content := loadTestFile(t, "complex.tscn")
	uri := "file:///test/complex.tscn"
	if err := client.openDocument(uri, content); err != nil {
		t.Fatalf("failed to open document: %v", err)
	}

	// Give server time to process
	time.Sleep(50 * time.Millisecond)

	// Request semantic tokens
	params := semanticTokensParams{
		TextDocument: textDocumentIdentifier{URI: uri},
	}
	result, err := client.sendRequest(ctx, "textDocument/semanticTokens/full", params)
	if err != nil {
		t.Fatalf("semanticTokens request failed: %v", err)
	}

	if string(result) == "null" {
		t.Error("expected semantic tokens result, got null")
		return
	}

	var tokens semanticTokensResult
	if err := json.Unmarshal(result, &tokens); err != nil {
		t.Fatalf("failed to unmarshal semantic tokens: %v", err)
	}

	if len(tokens.Data) == 0 {
		t.Error("expected non-empty semantic tokens data for TSCN content")
	}
}

func TestLSPDiagnosticsOnErrors(t *testing.T) {
	t.Parallel()

	client := newTestLSPClient(t)
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
		defer cancel()
		client.shutdown(ctx)
		client.exit()
		client.close()
	}()

	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	_, err := client.initialize(ctx)
	if err != nil {
		t.Fatalf("initialize failed: %v", err)
	}

	// Open file with errors
	content := loadTestFile(t, "with_errors.tscn")
	uri := "file:///test/with_errors.tscn"
	if err := client.openDocument(uri, content); err != nil {
		t.Fatalf("failed to open document: %v", err)
	}

	// Wait for publishDiagnostics notification
	notifCtx, notifCancel := context.WithTimeout(ctx, 2*time.Second)
	defer notifCancel()

	params, err := client.waitForNotification(notifCtx, "textDocument/publishDiagnostics")
	if err != nil {
		t.Fatalf("failed to receive diagnostics notification: %v", err)
	}

	var diagParams publishDiagnosticsParams
	if err := json.Unmarshal(params, &diagParams); err != nil {
		t.Fatalf("failed to unmarshal diagnostics: %v", err)
	}

	if diagParams.URI != uri {
		t.Errorf("expected diagnostics for URI %s, got %s", uri, diagParams.URI)
	}

	if len(diagParams.Diagnostics) == 0 {
		t.Error("expected diagnostics for file with errors, got none")
	}
}

func TestLSPShutdownGracefully(t *testing.T) {
	t.Parallel()

	client := newTestLSPClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	// Full lifecycle
	_, err := client.initialize(ctx)
	if err != nil {
		t.Fatalf("initialize failed: %v", err)
	}

	content := loadTestFile(t, "complex.tscn")
	uri := "file:///test/complex.tscn"
	if err := client.openDocument(uri, content); err != nil {
		t.Fatalf("failed to open document: %v", err)
	}

	// Give server time to process
	time.Sleep(50 * time.Millisecond)

	// Close document
	if err := client.closeDocument(uri); err != nil {
		t.Fatalf("failed to close document: %v", err)
	}

	// Shutdown
	if err := client.shutdown(ctx); err != nil {
		t.Fatalf("shutdown failed: %v", err)
	}

	// Exit
	if err := client.exit(); err != nil {
		t.Fatalf("exit notification failed: %v", err)
	}

	// Wait and verify exit code
	exitCode := client.close()
	if exitCode != 0 {
		t.Errorf("expected exit code 0, got %d", exitCode)
	}
}

func TestLSPDoesNotCrashOnMalformedTSCN(t *testing.T) {
	t.Parallel()

	client := newTestLSPClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	// Initialize
	_, err := client.initialize(ctx)
	if err != nil {
		t.Fatalf("initialize failed: %v", err)
	}

	// Malformed TSCN content that previously caused infinite loop and memory exhaustion
	malformedContent := `[gd_scene format=3]

[sub_resource type="TestResource" id="TestResource_abc"]
data = {
"key1" "missing_colon_causes_parse_error",
"key2": 123,
"nested": {
"also" "broken"
}
}

[node name="Root" type="Node"]
values = [1, 2, @@@@, 3]
position = Vector3(1.0, @@@@, 3.0)
`

	uri := "file:///test/malformed.tscn"
	if err := client.openDocument(uri, malformedContent); err != nil {
		t.Fatalf("failed to open document: %v", err)
	}

	// Give server time to process
	time.Sleep(100 * time.Millisecond)

	// Server should still be responsive - try to get document symbols
	params := documentSymbolParams{
		TextDocument: textDocumentIdentifier{URI: uri},
	}
	_, err = client.sendRequest(ctx, "textDocument/documentSymbol", params)
	if err != nil {
		t.Fatalf("documentSymbol request failed after malformed input: %v", err)
	}

	// Shutdown gracefully
	if err := client.shutdown(ctx); err != nil {
		t.Fatalf("shutdown failed: %v", err)
	}

	// Exit
	if err := client.exit(); err != nil {
		t.Fatalf("exit notification failed: %v", err)
	}

	// Verify exit code is 0
	exitCode := client.close()
	if exitCode != 0 {
		t.Errorf("expected exit code 0, got %d", exitCode)
	}
}

func TestLSPWithRealGodotProject(t *testing.T) {
	t.Parallel()

	client := newTestLSPClient(t)
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
		defer cancel()
		client.shutdown(ctx)
		client.exit()
		client.close()
	}()

	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout*2) // Longer timeout for large files
	defer cancel()

	_, err := client.initialize(ctx)
	if err != nil {
		t.Fatalf("initialize failed: %v", err)
	}

	// Load a real complex TSCN file from the hextracer project
	projectRoot, err := findProjectRoot()
	if err != nil {
		t.Fatalf("failed to find project root: %v", err)
	}

	mainMenuPath := filepath.Join(projectRoot, "test", "hextracer", "scenes", "main_menu.tscn")
	content, err := os.ReadFile(mainMenuPath)
	if err != nil {
		t.Skipf("hextracer test project not available: %v", err)
	}

	uri := "file:///hextracer/scenes/main_menu.tscn"
	if err := client.openDocument(uri, string(content)); err != nil {
		t.Fatalf("failed to open document: %v", err)
	}

	// Give server time to process large file
	time.Sleep(200 * time.Millisecond)

	// Test document symbols on a real file
	symbolParams := documentSymbolParams{
		TextDocument: textDocumentIdentifier{URI: uri},
	}
	symbolResult, err := client.sendRequest(ctx, "textDocument/documentSymbol", symbolParams)
	if err != nil {
		t.Fatalf("documentSymbol request failed: %v", err)
	}
	t.Logf("Document symbols response length: %d bytes", len(symbolResult))

	// Test hover on a real file
	hoverParams := hoverParams{
		TextDocument: textDocumentIdentifier{URI: uri},
		Position:     position{Line: 57, Character: 10}, // [node name="MainMenu"...] line
	}
	_, err = client.sendRequest(ctx, "textDocument/hover", hoverParams)
	if err != nil {
		t.Fatalf("hover request failed: %v", err)
	}

	// Test semantic tokens
	tokenParams := semanticTokensParams{
		TextDocument: textDocumentIdentifier{URI: uri},
	}
	_, err = client.sendRequest(ctx, "textDocument/semanticTokens/full", tokenParams)
	if err != nil {
		t.Fatalf("semanticTokens request failed: %v", err)
	}

	t.Log("Successfully processed real Godot project file")
}

func TestLSPWithMultipleRealFiles(t *testing.T) {
	t.Parallel()

	client := newTestLSPClient(t)
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
		defer cancel()
		client.shutdown(ctx)
		client.exit()
		client.close()
	}()

	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout*3)
	defer cancel()

	_, err := client.initialize(ctx)
	if err != nil {
		t.Fatalf("initialize failed: %v", err)
	}

	projectRoot, err := findProjectRoot()
	if err != nil {
		t.Fatalf("failed to find project root: %v", err)
	}

	scenesDir := filepath.Join(projectRoot, "test", "hextracer", "scenes")
	entries, err := os.ReadDir(scenesDir)
	if err != nil {
		t.Skipf("hextracer test project not available: %v", err)
	}

	filesOpened := 0
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".tscn") {
			filePath := filepath.Join(scenesDir, entry.Name())
			content, err := os.ReadFile(filePath)
			if err != nil {
				t.Logf("Warning: could not read %s: %v", entry.Name(), err)
				continue
			}

			uri := "file:///hextracer/scenes/" + entry.Name()
			if err := client.openDocument(uri, string(content)); err != nil {
				t.Fatalf("failed to open %s: %v", entry.Name(), err)
			}
			filesOpened++

			// Request document symbols for each file
			params := documentSymbolParams{
				TextDocument: textDocumentIdentifier{URI: uri},
			}
			_, err = client.sendRequest(ctx, "textDocument/documentSymbol", params)
			if err != nil {
				t.Fatalf("documentSymbol request failed for %s: %v", entry.Name(), err)
			}
		}
	}

	t.Logf("Successfully processed %d TSCN files from hextracer project", filesOpened)
}

// LSP Definition types for testing
type definitionParams struct {
	TextDocument textDocumentIdentifier `json:"textDocument"`
	Position     position               `json:"position"`
}

type locationResult struct {
	URI   string   `json:"uri"`
	Range lspRange `json:"range"`
}

func TestLSPGoToDefinitionResPath(t *testing.T) {
	t.Parallel()

	client := newTestLSPClient(t)
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
		defer cancel()
		client.shutdown(ctx)
		client.exit()
		client.close()
	}()

	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout*2)
	defer cancel()

	// Initialize with workspace folder pointing to hextracer project
	projectRoot, err := findProjectRoot()
	if err != nil {
		t.Fatalf("failed to find project root: %v", err)
	}

	hextracerRoot := filepath.Join(projectRoot, "test", "hextracer")

	// Check if the hextracer project exists
	if _, err := os.Stat(filepath.Join(hextracerRoot, "project.godot")); os.IsNotExist(err) {
		t.Skip("hextracer test project not available")
	}

	// Custom initialization with workspace folder
	initParams := initializeParams{
		ProcessID: os.Getpid(),
		Capabilities: clientCapabilities{
			TextDocument: textDocumentClientCapabilities{
				Synchronization: textDocumentSyncClientCapabilities{
					DidSave: true,
				},
			},
		},
		RootURI: stringPtr("file://" + hextracerRoot),
	}

	result, err := client.sendRequest(ctx, "initialize", initParams)
	if err != nil {
		t.Fatalf("initialize failed: %v", err)
	}
	t.Logf("Initialize result: %s", string(result))

	// Send initialized notification
	if err := client.sendNotification("initialized", struct{}{}); err != nil {
		t.Fatalf("failed to send initialized: %v", err)
	}

	// Load the main_menu.tscn file which has res:// paths
	mainMenuPath := filepath.Join(hextracerRoot, "scenes", "main_menu.tscn")
	content, err := os.ReadFile(mainMenuPath)
	if err != nil {
		t.Fatalf("failed to read main_menu.tscn: %v", err)
	}

	// Use the real file path as URI
	uri := "file://" + mainMenuPath
	if err := client.openDocument(uri, string(content)); err != nil {
		t.Fatalf("failed to open document: %v", err)
	}

	// Give server time to process
	time.Sleep(100 * time.Millisecond)

	// Find line with ext_resource that has res:// path
	// Line 3: [ext_resource type="Script" path="res://Scripts/MainMenu.cs" id="1_mainmenu"]
	params := definitionParams{
		TextDocument: textDocumentIdentifier{URI: uri},
		Position:     position{Line: 2, Character: 45}, // Around the path
	}

	defResult, err := client.sendRequest(ctx, "textDocument/definition", params)
	if err != nil {
		t.Fatalf("definition request failed: %v", err)
	}

	t.Logf("Definition result: %s", string(defResult))

	// The result should be a location pointing to the correct file
	if string(defResult) != "null" {
		var loc locationResult
		if err := json.Unmarshal(defResult, &loc); err != nil {
			// It might be an array
			var locs []locationResult
			if err := json.Unmarshal(defResult, &locs); err != nil {
				t.Logf("Could not parse definition result: %v", err)
			} else if len(locs) > 0 {
				loc = locs[0]
			}
		}

		if loc.URI != "" {
			// Verify the URI points to the project root, not relative to the file
			if !strings.Contains(loc.URI, "hextracer") {
				t.Errorf("Definition URI should be within hextracer project: %s", loc.URI)
			}
			// Should NOT contain scenes/Scripts (which would be wrong relative path)
			if strings.Contains(loc.URI, "scenes/Scripts") {
				t.Errorf("Definition URI incorrectly appended to file directory: %s", loc.URI)
			}
			t.Logf("Definition resolved to: %s", loc.URI)
		}
	}
}

func stringPtr(s string) *string {
	return &s
}

// Document link types for testing
type documentLinkParams struct {
	TextDocument textDocumentIdentifier `json:"textDocument"`
}

type documentLinkResult struct {
	Range   lspRange `json:"range"`
	Target  string   `json:"target"`
	Tooltip string   `json:"tooltip,omitempty"`
}

func TestURIPathConversion(t *testing.T) {
	// Test that Windows-style URIs are handled correctly
	// This is a unit test embedded in e2e for convenience

	testCases := []struct {
		name     string
		uri      string
		wantPath string
	}{
		{
			name:     "Unix path",
			uri:      "file:///home/user/project/file.tscn",
			wantPath: "/home/user/project/file.tscn",
		},
		{
			name:     "Windows path with encoded colon",
			uri:      "file:///c%3A/Users/user/project/file.tscn",
			wantPath: "c:/Users/user/project/file.tscn",
		},
		{
			name:     "Windows path with literal colon",
			uri:      "file:///C:/Users/user/project/file.tscn",
			wantPath: "C:/Users/user/project/file.tscn",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Parse URI to simulate what the server does
			parsedURI, err := url.Parse(tc.uri)
			if err != nil {
				t.Fatalf("failed to parse URI: %v", err)
			}

			path := parsedURI.Path

			// Apply Windows path fix
			if len(path) >= 3 && path[0] == '/' && path[2] == ':' {
				path = path[1:]
			}

			if path != tc.wantPath {
				t.Errorf("got path %q, want %q", path, tc.wantPath)
			}
		})
	}
}

func TestLSPDocumentLinkRange(t *testing.T) {
	t.Parallel()

	client := newTestLSPClient(t)
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
		defer cancel()
		client.shutdown(ctx)
		client.exit()
		client.close()
	}()

	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout*2)
	defer cancel()

	// Initialize with workspace folder
	projectRoot, err := findProjectRoot()
	if err != nil {
		t.Fatalf("failed to find project root: %v", err)
	}

	hextracerRoot := filepath.Join(projectRoot, "test", "hextracer")
	if _, err := os.Stat(filepath.Join(hextracerRoot, "project.godot")); os.IsNotExist(err) {
		t.Skip("hextracer test project not available")
	}

	initParams := initializeParams{
		ProcessID: os.Getpid(),
		Capabilities: clientCapabilities{
			TextDocument: textDocumentClientCapabilities{
				Synchronization: textDocumentSyncClientCapabilities{
					DidSave: true,
				},
			},
		},
		RootURI: stringPtr("file://" + hextracerRoot),
	}

	_, err = client.sendRequest(ctx, "initialize", initParams)
	if err != nil {
		t.Fatalf("initialize failed: %v", err)
	}

	if err := client.sendNotification("initialized", struct{}{}); err != nil {
		t.Fatalf("failed to send initialized: %v", err)
	}

	// Use a simple test content with a known path
	content := `[gd_scene format=3]
[ext_resource type="Script" path="res://Scripts/Test.gd" id="1_test"]
[node name="Root" type="Node"]
`

	uri := "file://" + filepath.Join(hextracerRoot, "test.tscn")
	if err := client.openDocument(uri, content); err != nil {
		t.Fatalf("failed to open document: %v", err)
	}

	time.Sleep(100 * time.Millisecond)

	// Request document links
	params := documentLinkParams{
		TextDocument: textDocumentIdentifier{URI: uri},
	}

	result, err := client.sendRequest(ctx, "textDocument/documentLink", params)
	if err != nil {
		t.Fatalf("documentLink request failed: %v", err)
	}

	t.Logf("Document link result: %s", string(result))

	if string(result) == "null" || string(result) == "[]" {
		t.Log("No document links returned (may be expected if path doesn't exist)")
		return
	}

	var links []documentLinkResult
	if err := json.Unmarshal(result, &links); err != nil {
		t.Fatalf("failed to unmarshal document links: %v", err)
	}

	if len(links) == 0 {
		t.Log("No document links returned")
		return
	}

	// Verify the first link's range covers the full path including extension
	link := links[0]
	t.Logf("Link range: line %d, col %d-%d", link.Range.Start.Line, link.Range.Start.Character, link.Range.End.Character)
	t.Logf("Link target: %s", link.Target)

	// Line 1: [ext_resource type="Script" path="res://Scripts/Test.gd" id="1_test"]
	// The path "res://Scripts/Test.gd" should be fully highlighted (including .gd)
	// The quoted string starts around column 35 and ends at column 57
	expectedPath := `"res://Scripts/Test.gd"`
	rangeLength := link.Range.End.Character - link.Range.Start.Character

	// The range should cover the entire quoted path string
	if rangeLength != len(expectedPath) {
		t.Errorf("expected range length %d (for %s), got %d", len(expectedPath), expectedPath, rangeLength)
	}

	// Verify target includes the extension
	if !strings.HasSuffix(link.Target, ".gd") {
		t.Errorf("expected target to end with .gd, got %s", link.Target)
	}
}
