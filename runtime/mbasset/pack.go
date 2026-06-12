package mbasset

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"moonbasic/runtime"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

type packManifest struct {
	Textures map[string]string `json:"textures"`
	Models   map[string]string `json:"models"`
	Sounds   map[string]string `json:"sounds"`
}

type packState struct {
	baseDir  string
	manifest packManifest
	textures map[string]heap.Handle
	models   map[string]heap.Handle
	sounds   map[string]heap.Handle
}

func manifestLookup(table map[string]string, id, idLower string) string {
	if table == nil {
		return ""
	}
	if rel, ok := table[id]; ok {
		return rel
	}
	if rel, ok := table[idLower]; ok {
		return rel
	}
	for k, v := range table {
		if strings.EqualFold(k, id) {
			return v
		}
	}
	return ""
}

func (m *Module) clearPack(rt *runtime.Runtime) {
	if m.pack == nil {
		return
	}
	reg := runtime.ActiveRegistry()
	if reg != nil && rt != nil {
		freeKind := func(handles map[string]heap.Handle, cmd string) {
			for id, h := range handles {
				if h == 0 {
					continue
				}
				_, _ = reg.Call(cmd, []value.Value{value.FromHandle(int32(h))})
				handles[id] = 0
			}
		}
		freeKind(m.pack.textures, "TEXTURE.FREE")
		freeKind(m.pack.models, "MODEL.FREE")
		freeKind(m.pack.sounds, "FREESOUND")
	}
	m.pack = nil
}

func (m *Module) assetLoadPack(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("ASSET.LOADPACK expects 1 argument (manifest path)")
	}
	path, err := rt.ArgString(args, 0)
	if err != nil {
		return value.Nil, err
	}
	abs := rt.ResolveAssetPath(path)
	b, err := os.ReadFile(abs)
	if err != nil {
		return value.Nil, err
	}
	var man packManifest
	if err := json.Unmarshal(b, &man); err != nil {
		return value.Nil, err
	}
	m.mu.Lock()
	m.clearPack(rt)
	m.pack = &packState{
		baseDir:  filepath.Dir(abs),
		manifest: man,
		textures: make(map[string]heap.Handle),
		models:   make(map[string]heap.Handle),
		sounds:   make(map[string]heap.Handle),
	}
	m.mu.Unlock()
	return value.Nil, nil
}

func (m *Module) packPath(id, kind string) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.pack == nil {
		return "", fmt.Errorf("ASSET: no pack loaded (call ASSET.LOADPACK first)")
	}
	id = strings.TrimSpace(id)
	idLower := strings.ToLower(id)
	var rel string
	switch kind {
	case "texture":
		rel = manifestLookup(m.pack.manifest.Textures, id, idLower)
	case "model":
		rel = manifestLookup(m.pack.manifest.Models, id, idLower)
	case "sound":
		rel = manifestLookup(m.pack.manifest.Sounds, id, idLower)
	default:
		return "", fmt.Errorf("ASSET: internal kind %q", kind)
	}
	if rel == "" {
		return "", fmt.Errorf("ASSET: unknown %s id %q", kind, id)
	}
	if filepath.IsAbs(rel) {
		return rel, nil
	}
	return filepath.Clean(filepath.Join(m.pack.baseDir, filepath.FromSlash(rel))), nil
}

func (m *Module) loadCached(rt *runtime.Runtime, kind, id, loadCmd string) (value.Value, error) {
	path, err := m.packPath(id, kind)
	if err != nil {
		return value.Nil, err
	}

	cacheKey := strings.ToLower(strings.TrimSpace(id))
	m.mu.Lock()
	if m.pack == nil {
		m.mu.Unlock()
		return value.Nil, fmt.Errorf("ASSET: no pack loaded")
	}
	var cache map[string]heap.Handle
	switch kind {
	case "texture":
		cache = m.pack.textures
	case "model":
		cache = m.pack.models
	case "sound":
		cache = m.pack.sounds
	}
	if h, ok := cache[cacheKey]; ok && h != 0 {
		m.mu.Unlock()
		return value.FromHandle(int32(h)), nil
	}
	m.mu.Unlock()

	v, err := assetCall(rt, loadCmd, []value.Value{value.FromStringIndex(rt.Heap.Intern(path))})
	if err != nil {
		return value.Nil, err
	}
	if v.Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("ASSET.%s: load returned non-handle", strings.ToUpper(kind))
	}
	m.mu.Lock()
	if m.pack != nil {
		switch kind {
		case "texture":
			m.pack.textures[cacheKey] = heap.Handle(v.IVal)
		case "model":
			m.pack.models[cacheKey] = heap.Handle(v.IVal)
		case "sound":
			m.pack.sounds[cacheKey] = heap.Handle(v.IVal)
		}
	}
	m.mu.Unlock()
	return v, nil
}

func (m *Module) assetTexture(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("ASSET.TEXTURE expects 1 argument (id$)")
	}
	id, err := rt.ArgString(args, 0)
	if err != nil {
		return value.Nil, err
	}
	return m.loadCached(rt, "texture", id, "TEXTURE.LOAD")
}

func (m *Module) assetModel(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("ASSET.MODEL expects 1 argument (id$)")
	}
	id, err := rt.ArgString(args, 0)
	if err != nil {
		return value.Nil, err
	}
	return m.loadCached(rt, "model", id, "MODEL.LOAD")
}

func (m *Module) assetSound(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("ASSET.SOUND expects 1 argument (id$)")
	}
	id, err := rt.ArgString(args, 0)
	if err != nil {
		return value.Nil, err
	}
	return m.loadCached(rt, "sound", id, "AUDIO.LOADSOUND")
}

func (m *Module) assetUnload(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 0 {
		return value.Nil, fmt.Errorf("ASSET.UNLOAD expects 0 arguments")
	}
	m.mu.Lock()
	m.clearPack(rt)
	m.mu.Unlock()
	return value.Nil, nil
}

func assetCall(rt *runtime.Runtime, name string, args []value.Value) (value.Value, error) {
	reg := runtime.ActiveRegistry()
	if reg == nil {
		return value.Nil, fmt.Errorf("%s: registry not active", name)
	}
	return reg.Call(name, args)
}
