// Package checklistaliases registers beginner-facing API names from the final polish
// checklist as aliases to canonical moonBASIC commands (WINDOW.*, TIME.*, etc.).
package checklistaliases

import "moonbasic/runtime"

// Module registers checklist API aliases.
type Module struct{}

// NewModule creates the alias module.
func NewModule() *Module { return &Module{} }

func (m *Module) Register(r runtime.Registrar) {
	registerAPP(r)
	registerRENDER(r)
	registerACTION(r)
	registerINPUT(r)
	registerJSON(r)
	registerSAVE(r)
	registerTEXT(r)
	registerAUDIO(r)
	registerPICK(r)
	registerBODY(r)
}

func (m *Module) Reset()     {}
func (m *Module) Shutdown() {}
