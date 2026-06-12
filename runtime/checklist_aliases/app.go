package checklistaliases

import "moonbasic/runtime"

func registerAPP(r runtime.Registrar) {
	r.Register("APP.OPEN", "app", forward("WINDOW.OPEN"))
	r.Register("APP.CLOSE", "app", forward0("WINDOW.CLOSE"))
	r.Register("APP.SHOULDCLOSE", "app", forward0("WINDOW.SHOULDCLOSE"))
	r.Register("APP.SETFPS", "app", forward("WINDOW.SETFPS"))
	r.Register("APP.WIDTH", "app", forward0("WINDOW.WIDTH"))
	r.Register("APP.HEIGHT", "app", forward0("WINDOW.HEIGHT"))
	r.Register("APP.GETFPS", "app", forward0("TIME.GETFPS"))
	r.Register("APP.TIME", "app", forward0("TIME.GET"))
	r.Register("APP.DELTA", "app", forward0("TIME.DELTA"))
	r.Register("APP.VERSION", "app", forward0("SYSTEM.VERSION"))
}
