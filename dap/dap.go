// Package dap implements a minimal Debug Adapter Protocol server over stdio for moonBASIC.
package dap

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sync"

	"moonbasic/vm"
)

type message struct {
	Seq        int             `json:"seq"`
	Type       string          `json:"type"`
	Event      string          `json:"event,omitempty"`
	Command    string          `json:"command,omitempty"`
	Arguments  json.RawMessage `json:"arguments,omitempty"`
	RequestSeq int             `json:"request_seq,omitempty"`
	Success    bool            `json:"success,omitempty"`
	Body       any             `json:"body,omitempty"`
	Message    string          `json:"message,omitempty"`
}

type session struct {
	mu            sync.Mutex
	seq           int
	out           io.Writer
	program       string
	breaks        map[int]bool
	vm            *vm.VM
	paused  bool
	running bool
	runDone chan struct{}
}

// ServeStdio runs a DAP session on stdin/stdout (attach to VS Code debug adapter).
func ServeStdio() error {
	s := &session{out: os.Stdout, breaks: make(map[int]bool)}
	r := bufio.NewReader(os.Stdin)
	for {
		headers, err := readHeaders(r)
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		n := headers["Content-Length"]
		if n <= 0 {
			continue
		}
		buf := make([]byte, n)
		if _, err := io.ReadFull(r, buf); err != nil {
			return err
		}
		var req message
		if err := json.Unmarshal(buf, &req); err != nil {
			return err
		}
		if req.Type == "request" {
			if err := s.handleRequest(&req); err != nil {
				return err
			}
		}
	}
}

func readHeaders(r *bufio.Reader) (map[string]int, error) {
	h := map[string]int{}
	for {
		line, err := r.ReadString('\n')
		if err != nil {
			return nil, err
		}
		if line == "\r\n" || line == "\n" {
			break
		}
		var k string
		var v int
		if _, err := fmt.Sscanf(line, "%[^:]: %d", &k, &v); err == nil {
			h[k] = v
		}
	}
	return h, nil
}

func (s *session) send(ev message) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.seq++
	ev.Seq = s.seq
	b, err := json.Marshal(ev)
	if err != nil {
		return err
	}
	if _, err := fmt.Fprintf(s.out, "Content-Length: %d\r\n\r\n", len(b)); err != nil {
		return err
	}
	_, err = s.out.Write(b)
	return err
}

func (s *session) respond(req *message, body any) error {
	return s.send(message{Type: "response", RequestSeq: req.Seq, Command: req.Command, Success: true, Body: body})
}

func (s *session) respondFail(req *message, err error) error {
	return s.send(message{
		Type: "response", RequestSeq: req.Seq, Command: req.Command,
		Success: false, Message: err.Error(),
	})
}

func (s *session) breakLines() []int {
	lines := make([]int, 0, len(s.breaks))
	for ln := range s.breaks {
		lines = append(lines, ln)
	}
	return lines
}

func (s *session) onPaused(v *vm.VM, reason string) {
	s.mu.Lock()
	s.vm = v
	s.paused = true
	s.mu.Unlock()
	_ = s.send(message{
		Type:  "event",
		Event: "stopped",
		Body: map[string]any{
			"reason":            reason,
			"threadId":          1,
			"allThreadsStopped": true,
		},
	})
}

func (s *session) clearPaused() {
	s.mu.Lock()
	s.paused = false
	s.mu.Unlock()
}

func (s *session) handleRequest(req *message) error {
	switch req.Command {
	case "initialize":
		if err := s.respond(req, map[string]any{
			"supportsConfigurationDoneRequest": true,
		}); err != nil {
			return err
		}
		return s.send(message{Type: "event", Event: "initialized"})
	case "configurationDone":
		return s.respond(req, map[string]any{})
	case "launch":
		var args struct {
			Program string `json:"program"`
		}
		_ = json.Unmarshal(req.Arguments, &args)
		if args.Program != "" {
			s.program = args.Program
		}
		if err := s.respond(req, map[string]any{}); err != nil {
			return err
		}
		s.beginDebugRun()
		return nil
	case "continue":
		s.clearPaused()
		s.mu.Lock()
		v := s.vm
		s.mu.Unlock()
		if v != nil {
			v.DebugContinue()
		}
		return s.respond(req, map[string]any{"allThreadsContinued": true})
	case "next", "stepIn", "stepOut":
		s.clearPaused()
		s.mu.Lock()
		v := s.vm
		s.mu.Unlock()
		if v != nil {
			v.DebugStepOnce()
		}
		return s.respond(req, map[string]any{})
	case "disconnect", "terminate":
		s.mu.Lock()
		v := s.vm
		s.mu.Unlock()
		if v != nil {
			v.DebugContinue()
		}
		return s.respond(req, map[string]any{})
	case "setBreakpoints":
		var args struct {
			Breakpoints []struct {
				Line int `json:"line"`
			} `json:"breakpoints"`
		}
		_ = json.Unmarshal(req.Arguments, &args)
		s.breaks = make(map[int]bool)
		out := make([]map[string]any, 0, len(args.Breakpoints))
		for _, bp := range args.Breakpoints {
			ln := bp.Line + 1
			s.breaks[ln] = true
			out = append(out, map[string]any{"id": ln, "verified": true, "line": bp.Line})
		}
		s.mu.Lock()
		v := s.vm
		s.mu.Unlock()
		if v != nil {
			v.SetBreakLines(s.breakLines())
		}
		return s.respond(req, map[string]any{"breakpoints": out})
	case "threads":
		return s.respond(req, map[string]any{"threads": []map[string]any{{"id": 1, "name": "main"}}})
	case "stackTrace":
		line := 0
		path := s.program
		s.mu.Lock()
		if s.vm != nil {
			if ln := s.vm.CurrentLine(); ln > 0 {
				line = ln - 1
			}
			if p := s.vm.CurrentSourcePath(); p != "" {
				path = p
			}
		}
		s.mu.Unlock()
		return s.respond(req, map[string]any{
			"stackFrames": []map[string]any{{
				"id": 1, "name": "main", "line": line, "column": 0,
				"source": map[string]any{"name": path, "path": path},
			}},
			"totalFrames": 1,
		})
	case "scopes":
		return s.respond(req, map[string]any{
			"scopes": []map[string]any{{"name": "Globals", "variablesReference": 1, "expensive": false}},
		})
	case "variables":
		var args struct {
			VariablesReference int `json:"variablesReference"`
		}
		_ = json.Unmarshal(req.Arguments, &args)
		vars := []map[string]any{}
		if args.VariablesReference == 1 {
			s.mu.Lock()
			v := s.vm
			s.mu.Unlock()
			if v != nil {
				for name, val := range v.DebugGlobals() {
					vars = append(vars, map[string]any{
						"name":               name,
						"value":              val,
						"type":               "string",
						"variablesReference": 0,
					})
				}
			}
		}
		return s.respond(req, map[string]any{"variables": vars})
	default:
		return s.respond(req, map[string]any{})
	}
}

func (s *session) sendTerminated(err error) {
	body := map[string]any{"restart": false}
	if err != nil {
		_ = s.send(message{Type: "event", Event: "output", Body: map[string]any{
			"category": "stderr",
			"output":   err.Error() + "\n",
		}})
	}
	_ = s.send(message{Type: "event", Event: "terminated", Body: body})
}
