package checklistaliases

import "moonbasic/runtime"

func registerPICK(r runtime.Registrar) {
	r.Register("PICK.MOUSE", "pick", forward("PICK.SCREENCAST"))
	r.Register("PICK.RAY", "pick", forward("PICK.CAST"))
	r.Register("PICK.DISTANCE", "pick", forward("PICK.DIST"))
}
