package mbasset

import (
	"fmt"
	"sync"

	"moonbasic/runtime"
	"moonbasic/vm/value"
)

// Module registers ASSET.* path helpers and manifest pack lookup.
type Module struct {
	mu   sync.Mutex
	pack *packState
}

func NewModule() *Module { return &Module{} }

func (m *Module) Register(reg runtime.Registrar) {
	reg.Register("ASSET.PATH", "asset", m.assetPath)
	reg.Register("ASSET.RESOLVE", "asset", m.assetResolve)
	reg.Register("ASSET.LOADPACK", "asset", m.assetLoadPack)
	reg.Register("ASSET.TEXTURE", "asset", m.assetTexture)
	reg.Register("ASSET.MODEL", "asset", m.assetModel)
	reg.Register("ASSET.SOUND", "asset", m.assetSound)
	reg.Register("ASSET.UNLOAD", "asset", m.assetUnload)
}

func (m *Module) Reset() {
	m.mu.Lock()
	m.pack = nil
	m.mu.Unlock()
}

func (m *Module) Shutdown() {
	m.mu.Lock()
	m.pack = nil
	m.mu.Unlock()
}

func (m *Module) assetPath(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("ASSET.PATH expects 1 string argument (relative folder, e.g. assets/)")
	}
	p, err := rt.ArgString(args, 0)
	if err != nil {
		return value.Nil, err
	}
	rt.SetAssetBase(p)
	return value.FromStringIndex(rt.Heap.Intern(p)), nil
}

func (m *Module) assetResolve(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("ASSET.RESOLVE expects 1 string argument (relative asset path)")
	}
	p, err := rt.ArgString(args, 0)
	if err != nil {
		return value.Nil, err
	}
	abs := rt.ResolveAssetPath(p)
	return value.FromStringIndex(rt.Heap.Intern(abs)), nil
}
