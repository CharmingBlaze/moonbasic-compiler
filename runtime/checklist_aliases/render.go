package checklistaliases

import (
	"fmt"

	"moonbasic/runtime"
	"moonbasic/vm/value"
)

func registerRENDER(r runtime.Registrar) {
	r.Register("RENDER.SETBACKGROUND", "render", forward("RENDER.CLEAR"))
	r.Register("RENDER.BEGIN", "render", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) == 1 {
			return call(rt, "RENDER.BEGIN3D", args)
		}
		if len(args) != 0 {
			return value.Nil, fmt.Errorf("RENDER.BEGIN expects 0 or 1 arguments (optional camera handle)")
		}
		active, err := call(rt, "CAMERA.GETACTIVE", nil)
		if err != nil {
			return value.Nil, err
		}
		if active.Kind == value.KindHandle && active.IVal == 0 {
			return value.Nil, fmt.Errorf("RENDER.BEGIN: no active camera (call CAMERA.SETACTIVE first)")
		}
		return call(rt, "RENDER.BEGIN3D", []value.Value{active})
	})
	r.Register("RENDER.END", "render", forward0("RENDER.END3D"))
}
