package checklistaliases

import "moonbasic/runtime"

func registerINPUT(r runtime.Registrar) {
	r.Register("INPUT.MOUSEDELTA_X", "input", forward0("INPUT.MOUSEDX"))
	r.Register("INPUT.MOUSEDELTA_Y", "input", forward0("INPUT.MOUSEDY"))
	r.Register("INPUT.GAMEPADBUTTONDOWN", "input", forward("INPUT.JOYBUTTON"))
	r.Register("INPUT.GAMEPADAXIS", "input", forward("INPUT.GETGAMEPADAXISVALUE"))
}
