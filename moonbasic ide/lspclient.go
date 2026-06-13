package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
)

// LSPItem is a completion suggestion from the language server.
type LSPItem struct {
	Label         string `json:"label"`
	Detail        string `json:"detail,omitempty"`
	InsertText    string `json:"insertText,omitempty"`
	Documentation string `json:"documentation,omitempty"`
}

// LSPClient talks to moonbasic --lsp over stdio (hover + completion).
type LSPClient struct {
	mu     sync.Mutex
	cmd    *exec.Cmd
	stdin  io.WriteCloser
	reader *bufio.Reader
	nextID atomic.Int64
	ready  bool
}

func newLSPClient() *LSPClient {
	return &LSPClient{}
}

func (c *LSPClient) start(exe string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.ready {
		return nil
	}
	if exe == "" {
		exe = "moonbasic"
	}
	cmd := exec.Command(exe, "--lsp")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	if err := cmd.Start(); err != nil {
		return err
	}
	c.cmd = cmd
	c.stdin = stdin
	c.reader = bufio.NewReader(stdout)
	c.ready = true
	return c.requestLocked(map[string]any{
		"jsonrpc": "2.0",
		"id":      c.nextID.Add(1),
		"method":  "initialize",
		"params": map[string]any{
			"processId":    nil,
			"rootUri":      nil,
			"capabilities": map[string]any{},
		},
	})
}

func (c *LSPClient) stop() {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.stdin != nil {
		_ = c.stdin.Close()
	}
	if c.cmd != nil && c.cmd.Process != nil {
		_ = c.cmd.Process.Kill()
	}
	c.ready = false
}

func readLSPFrame(r *bufio.Reader) ([]byte, error) {
	var length int
	for {
		line, err := r.ReadString('\n')
		if err != nil {
			return nil, err
		}
		line = strings.TrimRight(line, "\r\n")
		if line == "" {
			break
		}
		if strings.HasPrefix(strings.ToLower(line), "content-length:") {
			fmt.Sscanf(strings.TrimSpace(line[15:]), "%d", &length)
		}
	}
	if length <= 0 {
		return nil, fmt.Errorf("bad lsp frame")
	}
	buf := make([]byte, length)
	_, err := io.ReadFull(r, buf)
	return buf, err
}

func writeLSPFrame(w io.Writer, payload []byte) error {
	if _, err := fmt.Fprintf(w, "Content-Length: %d\r\n\r\n", len(payload)); err != nil {
		return err
	}
	_, err := w.Write(payload)
	return err
}

func (c *LSPClient) notifyLocked(msg map[string]any) error {
	msg["jsonrpc"] = "2.0"
	raw, _ := json.Marshal(msg)
	return writeLSPFrame(c.stdin, raw)
}

func (c *LSPClient) requestLocked(msg map[string]any) error {
	wantID := msg["id"]
	msg["jsonrpc"] = "2.0"
	raw, _ := json.Marshal(msg)
	if err := writeLSPFrame(c.stdin, raw); err != nil {
		return err
	}
	for {
		body, err := readLSPFrame(c.reader)
		if err != nil {
			return err
		}
		var env struct {
			ID    json.RawMessage `json:"id"`
			Error *struct {
				Message string `json:"message"`
			} `json:"error"`
		}
		if json.Unmarshal(body, &env) != nil || len(env.ID) == 0 {
			continue
		}
		if string(env.ID) != fmt.Sprint(wantID) {
			continue
		}
		if env.Error != nil {
			return fmt.Errorf(env.Error.Message)
		}
		return nil
	}
}

func (c *LSPClient) callLocked(msg map[string]any, result any) error {
	wantID := msg["id"]
	msg["jsonrpc"] = "2.0"
	raw, _ := json.Marshal(msg)
	if err := writeLSPFrame(c.stdin, raw); err != nil {
		return err
	}
	for {
		body, err := readLSPFrame(c.reader)
		if err != nil {
			return err
		}
		var env struct {
			ID     json.RawMessage `json:"id"`
			Result json.RawMessage `json:"result"`
			Error  *struct {
				Message string `json:"message"`
			} `json:"error"`
		}
		if json.Unmarshal(body, &env) != nil || len(env.ID) == 0 {
			continue
		}
		if string(env.ID) != fmt.Sprint(wantID) {
			continue
		}
		if env.Error != nil {
			return fmt.Errorf(env.Error.Message)
		}
		if result != nil && len(env.Result) > 0 {
			return json.Unmarshal(env.Result, result)
		}
		return nil
	}
}

func fileURI(path string) string {
	p := filepath.ToSlash(path)
	if strings.HasPrefix(p, "/") {
		return "file://" + p
	}
	return "file:///" + p
}

func (a *App) ensureLSP() error {
	if a.lsp != nil && a.lsp.ready {
		return nil
	}
	tc := a.resolveToolchain()
	if !tc.Found {
		return fmt.Errorf("moonbasic not found for LSP — set path in Settings")
	}
	a.lsp = newLSPClient()
	if err := a.lsp.start(tc.Moonbasic); err != nil {
		a.lsp = nil
		return err
	}
	_ = a.lsp.notifyLocked(map[string]any{"method": "initialized", "params": map[string]any{}})
	return nil
}

func (a *App) syncDocumentLocked(filePath, content string) error {
	uri := fileURI(filePath)
	_ = a.lsp.notifyLocked(map[string]any{
		"method": "textDocument/didOpen",
		"params": map[string]any{
			"textDocument": map[string]any{"uri": uri, "text": content},
		},
	})
	return a.lsp.notifyLocked(map[string]any{
		"method": "textDocument/didChange",
		"params": map[string]any{
			"textDocument":   map[string]any{"uri": uri, "version": 1},
			"contentChanges": []map[string]any{{"text": content}},
		},
	})
}

// GetLSPHover returns markdown help at a position (0-based line/col).
func (a *App) GetLSPHover(filePath, content string, line, col int) string {
	if err := a.ensureLSP(); err != nil {
		return ""
	}
	a.lsp.mu.Lock()
	defer a.lsp.mu.Unlock()
	_ = a.syncDocumentLocked(filePath, content)
	var result struct {
		Contents struct {
			Value string `json:"value"`
		} `json:"contents"`
	}
	id := a.lsp.nextID.Add(1)
	err := a.lsp.callLocked(map[string]any{
		"id":     id,
		"method": "textDocument/hover",
		"params": map[string]any{
			"textDocument": map[string]any{"uri": fileURI(filePath)},
			"position":     map[string]any{"line": line, "character": col},
		},
	}, &result)
	if err != nil {
		return ""
	}
	return result.Contents.Value
}

// GetLSPCompletion returns completion items at a position (0-based line/col).
func (a *App) GetLSPCompletion(filePath, content string, line, col int) []LSPItem {
	if err := a.ensureLSP(); err != nil {
		return nil
	}
	a.lsp.mu.Lock()
	defer a.lsp.mu.Unlock()
	_ = a.syncDocumentLocked(filePath, content)
	var result struct {
		Items []LSPItem `json:"items"`
	}
	id := a.lsp.nextID.Add(1)
	err := a.lsp.callLocked(map[string]any{
		"id":     id,
		"method": "textDocument/completion",
		"params": map[string]any{
			"textDocument": map[string]any{"uri": fileURI(filePath)},
			"position":     map[string]any{"line": line, "character": col},
		},
	}, &result)
	if err != nil {
		return nil
	}
	return result.Items
}
