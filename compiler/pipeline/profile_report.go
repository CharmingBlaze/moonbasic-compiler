package pipeline

import (
	"fmt"
	"html/template"
	"io"
	"os"
	"sort"

	"moonbasic/vm"
)

// PrintProfileReport writes a text summary of top hot lines to w.
func PrintProfileReport(rec *vm.ProfileRecorder, sourceFile string, w io.Writer, topN int) {
	if rec == nil || len(rec.LineHits) == 0 {
		fmt.Fprintln(w, "profile: no samples recorded")
		return
	}
	if topN <= 0 {
		topN = 10
	}
	lines := vm.TopProfileLines(rec, topN)
	fmt.Fprintf(w, "profile: top %d lines by instruction count (%s)\n", len(lines), sourceFile)
	for _, e := range lines {
		ms := float64(e.Nanos) / 1e6
		fmt.Fprintf(w, "  line %4d  %10d ops  %8.2f ms\n", e.Line, e.Count, ms)
	}
	funcs := vm.TopProfileFuncs(rec, topN)
	if len(funcs) > 0 {
		fmt.Fprintf(w, "profile: top %d functions by wall time\n", len(funcs))
		for _, f := range funcs {
			fmt.Fprintf(w, "  %-24s  %8.2f ms\n", f.Name, float64(f.Nanos)/1e6)
		}
	}
}

// WriteProfileHTML writes a simple HTML table of per-line instruction counts and wall time.
func WriteProfileHTML(rec *vm.ProfileRecorder, sourceFile, outPath string) error {
	if rec == nil {
		return fmt.Errorf("profile: nil recorder")
	}
	type row struct {
		Line  int
		Count uint64
		MS    float64
		Bar   int
	}
	var rows []row
	var max uint64
	for ln, c := range rec.LineHits {
		if c > max {
			max = c
		}
		ms := 0.0
		if rec.LineNanos != nil {
			ms = float64(rec.LineNanos[ln]) / 1e6
		}
		rows = append(rows, row{Line: ln, Count: c, MS: ms})
	}
	sort.Slice(rows, func(i, j int) bool {
		if rows[i].Count == rows[j].Count {
			return rows[i].Line < rows[j].Line
		}
		return rows[i].Count > rows[j].Count
	})
	for i := range rows {
		if max > 0 {
			rows[i].Bar = int(rows[i].Count * 100 / max)
		}
	}
	const page = `<!DOCTYPE html>
<html><head><meta charset="utf-8"><title>moonBASIC profile</title>
<style>
body{font-family:system-ui,sans-serif;margin:2rem;background:#111;color:#eee}
h1{font-size:1.2rem} table{border-collapse:collapse;width:100%;max-width:56rem}
td,th{padding:.35rem .5rem;text-align:left;border-bottom:1px solid #333}
.bar{background:#3a7;height:1rem;border-radius:2px}
.num{text-align:right;font-variant-numeric:tabular-nums}
</style></head><body>
<h1>Profile: {{.Source}}</h1>
<p>Instruction counts and attributed wall time per source line.</p>
<table><tr><th>Line</th><th class="num">Ops</th><th class="num">Time (ms)</th><th></th></tr>
{{range .Rows}}<tr><td>{{.Line}}</td><td class="num">{{.Count}}</td><td class="num">{{printf "%.2f" .MS}}</td><td><div class="bar" style="width:{{.Bar}}%"></div></td></tr>
{{end}}</table></body></html>`
	tmpl, err := template.New("profile").Parse(page)
	if err != nil {
		return err
	}
	f, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer f.Close()
	return tmpl.Execute(f, map[string]interface{}{
		"Source": sourceFile,
		"Rows":   rows,
	})
}

// WriteProfileFlameHTML writes a function-level wall-time bar chart (simple flame-style summary).
func WriteProfileFlameHTML(rec *vm.ProfileRecorder, sourceFile, outPath string) error {
	if rec == nil {
		return fmt.Errorf("profile: nil recorder")
	}
	funcs := vm.TopProfileFuncs(rec, 50)
	if len(funcs) == 0 {
		return fmt.Errorf("profile: no function samples recorded")
	}
	type row struct {
		Name string
		MS   float64
		Bar  int
	}
	var rows []row
	var max uint64
	for _, f := range funcs {
		if f.Nanos > max {
			max = f.Nanos
		}
		rows = append(rows, row{Name: f.Name, MS: float64(f.Nanos) / 1e6})
	}
	for i := range rows {
		if max > 0 {
			rows[i].Bar = int(funcs[i].Nanos * 100 / max)
		}
	}
	const page = `<!DOCTYPE html>
<html><head><meta charset="utf-8"><title>moonBASIC profile — functions</title>
<style>
body{font-family:system-ui,sans-serif;margin:2rem;background:#111;color:#eee}
h1{font-size:1.2rem} .row{display:flex;align-items:center;gap:.75rem;margin:.35rem 0;max-width:56rem}
.bar{background:#e85;height:1.25rem;border-radius:2px;min-width:2px}
.name{min-width:14rem;font-family:ui-monospace,Consolas,monospace}
.ms{color:#9aa3b5;min-width:5rem;text-align:right}
</style></head><body>
<h1>Function profile: {{.Source}}</h1>
<p>Wall time per user FUNCTION (nested calls included).</p>
{{range .Rows}}<div class="row"><span class="name">{{.Name}}</span><div class="bar" style="width:{{.Bar}}%"></div><span class="ms">{{printf "%.2f" .MS}} ms</span></div>
{{end}}</body></html>`
	tmpl, err := template.New("flame").Parse(page)
	if err != nil {
		return err
	}
	f, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer f.Close()
	return tmpl.Execute(f, map[string]interface{}{
		"Source": sourceFile,
		"Rows":   rows,
	})
}
