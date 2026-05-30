//go:build cgo || (windows && !cgo)

package mbentity

import (
	"unsafe"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func pickAnimFrameWithTime(e *ent, anim rl.ModelAnimation, animTime float32, animMode int32) int32 {
	if anim.FrameCount <= 0 {
		return 0
	}
	ext := e.getExt()
	savedTime, savedMode := ext.animTime, ext.animMode
	ext.animTime = animTime
	ext.animMode = animMode
	frame := pickAnimFrame(e, anim)
	ext.animTime = savedTime
	ext.animMode = savedMode
	return frame
}

func lerpMatrix(a, b rl.Matrix, t float32) rl.Matrix {
	if t <= 0 {
		return a
	}
	if t >= 1 {
		return b
	}
	pa := (*[16]float32)(unsafe.Pointer(&a))
	pb := (*[16]float32)(unsafe.Pointer(&b))
	var out rl.Matrix
	po := (*[16]float32)(unsafe.Pointer(&out))
	for i := 0; i < 16; i++ {
		po[i] = pa[i] + (pb[i]-pa[i])*t
	}
	return out
}

func (m *Module) updateEntityAnimation(e *ent, ext *entExt, dt float32) {
	if e == nil || !e.hasRLModel || ext == nil || len(ext.modelAnims) == 0 {
		return
	}
	ai := ext.animIndex
	if ai < 0 || int(ai) >= len(ext.modelAnims) {
		ai = 0
	}
	toAnim := ext.modelAnims[ai]
	if toAnim.FrameCount <= 0 {
		return
	}
	if ext.animSpeed != 0 {
		ext.animTime += dt * ext.animSpeed * 30
	}

	if ext.animBlendFrom >= 0 && ext.animBlendDur > 0 {
		ext.animBlendT += dt / ext.animBlendDur
		if ext.animBlendT > 1 {
			ext.animBlendT = 1
		}
		fromIdx := ext.animBlendFrom
		if fromIdx < 0 || int(fromIdx) >= len(ext.modelAnims) {
			ext.animBlendFrom = -1
		} else {
			fromAnim := ext.modelAnims[fromIdx]
			fromFrame := pickAnimFrameWithTime(e, fromAnim, ext.animBlendFromTime, ext.animMode)
			toFrame := pickAnimFrame(e, toAnim)
			m.applyBlendedPose(e, fromAnim, toAnim, fromFrame, toFrame, ext.animBlendT)
			if ext.animBlendT >= 1 {
				ext.animBlendFrom = -1
				ext.animBlendT = 0
				ext.animBlendDur = 0
			}
			return
		}
	}

	frame := pickAnimFrame(e, toAnim)
	rl.UpdateModelAnimation(e.rlModel, toAnim, frame)
	rl.UpdateModelAnimationBones(e.rlModel, toAnim, frame)
}

func (m *Module) applyBlendedPose(e *ent, fromAnim, toAnim rl.ModelAnimation, fromFrame, toFrame int32, t float32) {
	meshes := e.rlModel.GetMeshes()
	if len(meshes) == 0 || meshes[0].BoneMatrices == nil || meshes[0].BoneCount <= 0 {
		return
	}
	// Sample "from" pose into scratch matrices.
	rl.UpdateModelAnimation(e.rlModel, fromAnim, fromFrame)
	rl.UpdateModelAnimationBones(e.rlModel, fromAnim, fromFrame)
	ext := e.getExt()
	bones := unsafe.Slice(meshes[0].BoneMatrices, meshes[0].BoneCount)
	n := int(meshes[0].BoneCount)
	if cap(ext.blendBoneScratch) < n {
		ext.blendBoneScratch = make([]rl.Matrix, n)
	}
	copy(ext.blendBoneScratch[:n], bones[:n])

	// Sample "to" pose and lerp bone matrices.
	rl.UpdateModelAnimation(e.rlModel, toAnim, toFrame)
	rl.UpdateModelAnimationBones(e.rlModel, toAnim, toFrame)
	for i := 0; i < n; i++ {
		bones[i] = lerpMatrix(ext.blendBoneScratch[i], bones[i], t)
	}
}

func (m *Module) beginAnimCrossfade(ext *entExt, toIndex int32, durationSec float32) {
	if ext == nil {
		return
	}
	if durationSec <= 0 {
		durationSec = 0.1
	}
	ext.animBlendFrom = ext.animIndex
	ext.animBlendFromTime = ext.animTime
	ext.animIndex = toIndex
	ext.animTime = 0
	ext.animBlendT = 0
	ext.animBlendDur = durationSec
}
