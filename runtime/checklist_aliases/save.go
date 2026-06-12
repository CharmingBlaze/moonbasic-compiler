package checklistaliases

import "moonbasic/runtime"

func registerSAVE(r runtime.Registrar) {
	r.Register("SAVE.SET", "save", forward("SAVE.DATA"))
	r.Register("SAVE.WRITE", "save", forward("SAVE.WRITEFILE"))
	r.Register("SAVE.READ", "save", forward("SAVE.READFILE"))
}
