//go:build cgo || (windows && !cgo)

package mbsave

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"
	
	"moonbasic/runtime"
	"moonbasic/vm/value"
)

var (
	saveData = make(map[string]string)
	saveMu   sync.RWMutex
	savePath = "save.json"
)

type Module struct{}

func NewModule() *Module { return &Module{} }

func (m *Module) Register(r runtime.Registrar) {
	r.Register("SAVE.DATA", "save", m.saveData)
	r.Register("SAVE.GET", "save", m.saveGet)
	r.Register("SAVE.WRITEFILE", "save", m.saveWriteFile)
	r.Register("SAVE.READFILE", "save", m.saveReadFile)
	loadSaveData()
}

func (m *Module) Reset()     {}
func (m *Module) Shutdown() {
	flushSaveData()
}

func loadSaveData() {
	saveMu.Lock()
	defer saveMu.Unlock()
	b, err := os.ReadFile(savePath)
	if err == nil {
		json.Unmarshal(b, &saveData)
	}
}

func flushSaveData() {
	saveMu.RLock()
	defer saveMu.RUnlock()
	b, err := json.MarshalIndent(saveData, "", "  ")
	if err == nil {
		os.WriteFile(savePath, b, 0644)
	}
}

func (m *Module) saveData(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 2 { return value.Nil, fmt.Errorf("SAVE.DATA expects (key$, value$)") }
	k, _ := rt.ArgString(args, 0)
	v, _ := rt.ArgString(args, 1)
	k = strings.ToUpper(strings.TrimSpace(k))

	saveMu.Lock()
	saveData[k] = v
	saveMu.Unlock()
	
	return value.Nil, nil
}

func (m *Module) saveGet(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 1 { return value.Nil, fmt.Errorf("SAVE.GET expects (key$)") }
	k, _ := rt.ArgString(args, 0)
	k = strings.ToUpper(strings.TrimSpace(k))

	saveMu.RLock()
	v, ok := saveData[k]
	saveMu.RUnlock()

	if ok {
		return rt.RetString(v), nil
	}
	return rt.RetString(""), nil
}

func (m *Module) saveWriteFile(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("SAVE.WRITEFILE expects (path$)")
	}
	path, err := rt.ArgString(args, 0)
	if err != nil {
		return value.Nil, err
	}
	saveMu.RLock()
	snap := make(map[string]string, len(saveData))
	for k, v := range saveData {
		snap[k] = v
	}
	saveMu.RUnlock()
	b, err := json.MarshalIndent(snap, "", "  ")
	if err != nil {
		return value.Nil, err
	}
	if err := os.WriteFile(path, b, 0o644); err != nil {
		return value.Nil, err
	}
	return value.Nil, nil
}

func (m *Module) saveReadFile(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("SAVE.READFILE expects (path$)")
	}
	path, err := rt.ArgString(args, 0)
	if err != nil {
		return value.Nil, err
	}
	b, err := os.ReadFile(path)
	if err != nil {
		return value.Nil, err
	}
	next := make(map[string]string)
	if err := json.Unmarshal(b, &next); err != nil {
		return value.Nil, err
	}
	saveMu.Lock()
	saveData = next
	savePath = path
	saveMu.Unlock()
	return value.Nil, nil
}
