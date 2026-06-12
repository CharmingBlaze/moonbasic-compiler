package checklistaliases

import (
	"fmt"

	"moonbasic/runtime"
	"moonbasic/vm/value"
)

func registerTEXT(r runtime.Registrar) {
	r.Register("TEXT.DRAW", "text", forward("DRAW.TEXT"))
	r.Register("TEXT.DRAWFONT", "text", forward("DRAW.TEXTFONT"))
	r.Register("TEXT.SIZE", "text", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 1 {
			return value.Nil, fmt.Errorf("TEXT.SIZE expects 1 argument (text$)")
		}
		size := value.FromInt(24)
		return call(rt, "DRAW.TEXTWIDTH", []value.Value{args[0], size})
	})
}
