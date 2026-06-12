package checklistaliases

import "moonbasic/runtime"

func registerACTION(r runtime.Registrar) {
	r.Register("ACTION.BINDKEY", "action", forward("ACTION.MAPKEY"))
	r.Register("ACTION.BINDGAMEPAD", "action", forward("ACTION.MAPJOY"))
	r.Register("ACTION.BINDMOUSE", "action", forward("ACTION.MAPMOUSE"))
	r.Register("ACTION.HIT", "action", forward("ACTION.PRESSED"))
}
