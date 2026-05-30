package vm

import (
	"sort"
	"time"

	"moonbasic/lineprof"
)

var _ lineprof.LineProfiler = (*ProfileRecorder)(nil)

// ProfileRecorder counts VM instructions and wall time per 1-based source line and per user function.
type ProfileRecorder struct {
	LineHits  map[int]uint64
	LineNanos map[int]uint64
	FuncNanos map[string]uint64
	lastLine  int
	lastAt    time.Time
	seen      bool
	funcStack []funcProfileFrame
}

type funcProfileFrame struct {
	name string
	at   time.Time
}

// NewProfileRecorder returns an empty recorder.
func NewProfileRecorder() *ProfileRecorder {
	return &ProfileRecorder{
		LineHits:  make(map[int]uint64),
		LineNanos: make(map[int]uint64),
		FuncNanos: make(map[string]uint64),
	}
}

// EnterFunc marks entry into a user function for wall-time attribution.
func (p *ProfileRecorder) EnterFunc(name string) {
	if p == nil {
		return
	}
	p.flushElapsed(time.Now())
	p.funcStack = append(p.funcStack, funcProfileFrame{name: name, at: time.Now()})
}

// LeaveFunc marks return from a user function.
func (p *ProfileRecorder) LeaveFunc() {
	if p == nil || len(p.funcStack) == 0 {
		return
	}
	now := time.Now()
	top := p.funcStack[len(p.funcStack)-1]
	p.funcStack = p.funcStack[:len(p.funcStack)-1]
	if elapsed := now.Sub(top.at).Nanoseconds(); elapsed > 0 {
		if p.FuncNanos == nil {
			p.FuncNanos = make(map[string]uint64)
		}
		p.FuncNanos[top.name] += uint64(elapsed)
	}
	p.lastAt = now
}

func (p *ProfileRecorder) flushElapsed(now time.Time) {
	if !p.seen {
		p.lastAt = now
		p.seen = true
		return
	}
	if elapsed := now.Sub(p.lastAt).Nanoseconds(); elapsed > 0 {
		if p.lastLine >= 1 {
			if p.LineNanos == nil {
				p.LineNanos = make(map[int]uint64)
			}
			p.LineNanos[p.lastLine] += uint64(elapsed)
		}
		if len(p.funcStack) > 0 {
			name := p.funcStack[len(p.funcStack)-1].name
			if p.FuncNanos == nil {
				p.FuncNanos = make(map[string]uint64)
			}
			p.FuncNanos[name] += uint64(elapsed)
		}
	}
	p.lastAt = now
}

// RecordLine increments instruction count and attributes elapsed wall time to the previous line.
func (p *ProfileRecorder) RecordLine(line int) {
	if p == nil || line < 1 {
		return
	}
	now := time.Now()
	p.flushElapsed(now)
	p.lastLine = line
	if p.LineHits == nil {
		p.LineHits = make(map[int]uint64)
	}
	p.LineHits[line]++
}

// TopProfileLines returns up to n source lines with the highest hit counts (stable tie-break by line).
func TopProfileLines(p *ProfileRecorder, n int) []struct {
	Line  int
	Count uint64
	Nanos uint64
} {
	if p == nil || len(p.LineHits) == 0 {
		return nil
	}
	type pair struct {
		line  int
		count uint64
	}
	var ps []pair
	for ln, c := range p.LineHits {
		ps = append(ps, pair{ln, c})
	}
	sort.Slice(ps, func(i, j int) bool {
		if ps[i].count == ps[j].count {
			return ps[i].line < ps[j].line
		}
		return ps[i].count > ps[j].count
	})
	if n > len(ps) {
		n = len(ps)
	}
	out := make([]struct {
		Line  int
		Count uint64
		Nanos uint64
	}, n)
	for i := 0; i < n; i++ {
		out[i].Line = ps[i].line
		out[i].Count = ps[i].count
		if p.LineNanos != nil {
			out[i].Nanos = p.LineNanos[ps[i].line]
		}
	}
	return out
}

// TopProfileFuncs returns up to n user functions with highest wall-time totals.
func TopProfileFuncs(p *ProfileRecorder, n int) []struct {
	Name  string
	Nanos uint64
} {
	if p == nil || len(p.FuncNanos) == 0 {
		return nil
	}
	type pair struct {
		name  string
		nanos uint64
	}
	var ps []pair
	for name, ns := range p.FuncNanos {
		ps = append(ps, pair{name, ns})
	}
	sort.Slice(ps, func(i, j int) bool {
		if ps[i].nanos == ps[j].nanos {
			return ps[i].name < ps[j].name
		}
		return ps[i].nanos > ps[j].nanos
	})
	if n > len(ps) {
		n = len(ps)
	}
	out := make([]struct {
		Name  string
		Nanos uint64
	}, n)
	for i := 0; i < n; i++ {
		out[i].Name = ps[i].name
		out[i].Nanos = ps[i].nanos
	}
	return out
}

func (v *VM) profileEnterFunc(name string) {
	if rec, ok := v.Profiler.(*ProfileRecorder); ok {
		rec.EnterFunc(name)
	}
}

func (v *VM) profileLeaveFunc() {
	if rec, ok := v.Profiler.(*ProfileRecorder); ok {
		rec.LeaveFunc()
	}
}
