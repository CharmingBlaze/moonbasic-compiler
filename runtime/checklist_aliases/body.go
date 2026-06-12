package checklistaliases

import (
	"fmt"

	"moonbasic/runtime"
	"moonbasic/vm/value"
)

func registerBODY(r runtime.Registrar) {
	r.Register("BODY.ADDSTATICBOX", "body", bodyBox("STATIC"))
	r.Register("BODY.ADDDYNAMICBOX", "body", bodyBox("DYNAMIC"))
	r.Register("BODY.ADDSPHERE", "body", bodySphere("DYNAMIC"))
	r.Register("BODY.ADDCAPSULE", "body", bodyCapsule("DYNAMIC"))
	r.Register("BODY.SETMASS", "body", forward("ENTITY.SETMASS"))
	r.Register("BODY.SETFRICTION", "body", forward("PHYSICS.FRICTION"))
	r.Register("BODY.SETBOUNCE", "body", forward("PHYSICS.BOUNCE"))
	r.Register("BODY.APPLYFORCE", "body", forward("PHYSICS.FORCE"))
	r.Register("BODY.APPLYIMPULSE", "body", forward("PHYSICS.IMPULSE"))
}

func bodyBox(motion string) runtime.BuiltinFn {
	return func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 4 {
			return value.Nil, fmt.Errorf("BODY.ADD%sBOX expects (entity, width, height, depth)", motion)
		}
		ent := args[0]
		if _, err := call(rt, "ENTITY.SETSCALE", []value.Value{ent, args[1], args[2], args[3]}); err != nil {
			return value.Nil, err
		}
		motionVal := value.FromStringIndex(rt.Heap.Intern(motion))
		shapeVal := value.FromStringIndex(rt.Heap.Intern("BOX"))
		return call(rt, "ENTITY.ADDPHYSICS", []value.Value{ent, motionVal, shapeVal})
	}
}

func bodySphere(motion string) runtime.BuiltinFn {
	return func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 2 {
			return value.Nil, fmt.Errorf("BODY.ADDSPHERE expects (entity, radius)")
		}
		ent := args[0]
		r := args[1]
		if _, err := call(rt, "ENTITY.SETSCALE", []value.Value{ent, r, r, r}); err != nil {
			return value.Nil, err
		}
		motionVal := value.FromStringIndex(rt.Heap.Intern(motion))
		shapeVal := value.FromStringIndex(rt.Heap.Intern("SPHERE"))
		return call(rt, "ENTITY.ADDPHYSICS", []value.Value{ent, motionVal, shapeVal})
	}
}

func bodyCapsule(motion string) runtime.BuiltinFn {
	return func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 3 {
			return value.Nil, fmt.Errorf("BODY.ADDCAPSULE expects (entity, radius, height)")
		}
		ent := args[0]
		r, h := args[1], args[2]
		if _, err := call(rt, "ENTITY.SETSCALE", []value.Value{ent, r, h, r}); err != nil {
			return value.Nil, err
		}
		motionVal := value.FromStringIndex(rt.Heap.Intern(motion))
		shapeVal := value.FromStringIndex(rt.Heap.Intern("CAPSULE"))
		return call(rt, "ENTITY.ADDPHYSICS", []value.Value{ent, motionVal, shapeVal})
	}
}
