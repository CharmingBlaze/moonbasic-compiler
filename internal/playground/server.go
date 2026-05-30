package playground

import (
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"moonbasic/compiler/pipeline"
)

// Server hosts the static playground UI and compile API.
type Server struct {
	WebRoot string
	Addr    string
}

// CompileResult is returned from POST /api/compile.
type CompileResult struct {
	OK       bool     `json:"ok"`
	Errors   []string `json:"errors,omitempty"`
	Function int      `json:"functions,omitempty"`
	Disasm   string   `json:"disasm,omitempty"`
}

// RunResult is returned from POST /api/run.
type RunResult struct {
	OK     bool     `json:"ok"`
	Output string   `json:"output,omitempty"`
	Errors []string `json:"errors,omitempty"`
}

// ListenAndServe starts the playground HTTP server.
func (s *Server) ListenAndServe() error {
	if s.Addr == "" {
		s.Addr = "127.0.0.1:8765"
	}
	root := s.WebRoot
	if root == "" {
		root = defaultWebRoot()
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/api/compile", s.handleCompile)
	mux.HandleFunc("/api/run", s.handleRun)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		serveStatic(w, r, root)
	})
	fmt.Printf("moonBASIC playground: http://%s/\n", s.Addr)
	return http.ListenAndServe(s.Addr, mux)
}

func defaultWebRoot() string {
	if _, err := os.Stat("web/playground"); err == nil {
		return "web/playground"
	}
	// Module root when run from repo subdir.
	return filepath.Join("..", "..", "web", "playground")
}

func serveStatic(w http.ResponseWriter, r *http.Request, root string) {
	path := r.URL.Path
	if path == "/" {
		path = "/index.html"
	}
	path = filepath.Clean(strings.TrimPrefix(path, "/"))
	full := filepath.Join(root, path)
	data, err := os.ReadFile(full)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	if strings.HasSuffix(path, ".html") {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
	} else if strings.HasSuffix(path, ".js") {
		w.Header().Set("Content-Type", "application/javascript")
	} else if strings.HasSuffix(path, ".css") {
		w.Header().Set("Content-Type", "text/css")
	}
	w.Write(data)
}

func (s *Server) handleCompile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "POST only", http.StatusMethodNotAllowed)
		return
	}
	body, err := io.ReadAll(io.LimitReader(r.Body, 512*1024))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var req struct {
		Source string `json:"source"`
	}
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	res := compileSource(req.Source)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(res)
}

func (s *Server) handleRun(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "POST only", http.StatusMethodNotAllowed)
		return
	}
	body, err := io.ReadAll(io.LimitReader(r.Body, 512*1024))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var req struct {
		Source string `json:"source"`
	}
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	res := runSource(req.Source)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(res)
}

func runSource(src string) RunResult {
	src = strings.TrimSpace(src)
	if src == "" {
		return RunResult{OK: false, Errors: []string{"source is empty"}}
	}
	prog, err := pipeline.CompileSource("playground.mb", src)
	if err != nil {
		return RunResult{OK: false, Errors: []string{err.Error()}}
	}
	out, err := pipeline.RunProgramHeadless(prog)
	if err != nil {
		msg := err.Error()
		if out != "" {
			return RunResult{OK: false, Output: out, Errors: []string{msg}}
		}
		return RunResult{OK: false, Errors: []string{msg}}
	}
	return RunResult{OK: true, Output: out}
}

func compileSource(src string) CompileResult {
	src = strings.TrimSpace(src)
	if src == "" {
		return CompileResult{OK: false, Errors: []string{"source is empty"}}
	}
	prog, err := pipeline.CompileSource("playground.mb", src)
	if err != nil {
		return CompileResult{OK: false, Errors: []string{err.Error()}}
	}
	disasm := ""
	if prog.Main != nil {
		d := prog.Main.Disassemble()
		const maxLines = 40
		lines := strings.Split(d, "\n")
		if len(lines) > maxLines {
			lines = lines[:maxLines]
			disasm = strings.Join(lines, "\n") + "\n; …"
		} else {
			disasm = d
		}
	}
	return CompileResult{
		OK:       true,
		Function: len(prog.Functions),
		Disasm:   disasm,
	}
}

// OpenWebRoot returns an fs.FS for embedded/static checks (tests).
func OpenWebRoot(root string) fs.FS {
	return os.DirFS(root)
}
