//go:build cgo || (windows && !cgo)

package mbaudio

import "moonbasic/runtime"

func (m *Module) registerJamHooks(reg runtime.Registrar) {
	registerJamAudio(reg, m)
}
