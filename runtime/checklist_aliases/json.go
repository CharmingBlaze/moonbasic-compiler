package checklistaliases

import "moonbasic/runtime"

func registerJSON(r runtime.Registrar) {
	r.Register("JSON.STRINGIFY", "json", forward("JSON.TOSTRING"))
	r.Register("JSON.GET", "json", forward("JSON.QUERY"))
	r.Register("JSON.SET", "json", forward("JSON.SETSTRING"))
}
