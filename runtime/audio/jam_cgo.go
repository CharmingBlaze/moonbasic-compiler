//go:build cgo || (windows && !cgo)

package mbaudio

import (
	"fmt"
	"math"
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"

	"moonbasic/runtime"
	"moonbasic/vm/value"
)

type jamTone struct {
	freqHz   float64
	duration float64
	volume   float64
}

var jamSounds = map[string][]jamTone{
	"jump":    {{440, 0.08, 0.5}, {660, 0.1, 0.45}},
	"hit":     {{120, 0.06, 0.7}, {80, 0.08, 0.5}},
	"coin":    {{880, 0.05, 0.5}, {1175, 0.12, 0.45}},
	"shoot":   {{920, 0.04, 0.55}},
	"powerup": {{523, 0.06, 0.45}, {659, 0.06, 0.45}, {784, 0.1, 0.4}},
	"explode": {{200, 0.05, 0.6}, {100, 0.12, 0.55}, {60, 0.15, 0.4}},
	"select":  {{660, 0.04, 0.4}},
	"error":   {{180, 0.14, 0.55}},
}

func synthJamSound(tones []jamTone) rl.Wave {
	const sampleRate uint32 = 44100
	var pcm []byte
	for _, t := range tones {
		n := int(float64(sampleRate)*t.duration + 0.5)
		if n < 1 {
			continue
		}
		for i := 0; i < n; i++ {
			phase := float64(i) / float64(sampleRate)
			amp := t.volume
			if i > n*3/4 {
				amp *= float64(n-i) / float64(n / 4)
			}
			s := math.Sin(2 * math.Pi * t.freqHz * phase)
			v := int16(s * amp * 32767)
			pcm = append(pcm, byte(v&0xff), byte(v>>8))
		}
	}
	return rl.NewWave(uint32(len(pcm)/2), sampleRate, 16, 1, pcm)
}

func (m *Module) soundBuiltin(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("SOUND.BUILTIN expects 1 string name")
	}
	name, err := rt.ArgString(args, 0)
	if err != nil {
		return value.Nil, err
	}
	name = strings.ToLower(strings.TrimSpace(name))
	tones, ok := jamSounds[name]
	if !ok {
		keys := make([]string, 0, len(jamSounds))
		for k := range jamSounds {
			keys = append(keys, k)
		}
		return value.Nil, fmt.Errorf("SOUND.BUILTIN: unknown %q (try: %s)", name, strings.Join(keys, ", "))
	}
	initAudioOnce()
	w := synthJamSound(tones)
	snd := rl.LoadSoundFromWave(w)
	rl.UnloadWave(w)
	id, err := m.h.Alloc(&soundObj{snd: snd, gain: 1, pan: 0})
	if err != nil {
		rl.UnloadSound(snd)
		return value.Nil, err
	}
	return value.FromHandle(id), nil
}

func registerJamAudio(reg runtime.Registrar, m *Module) {
	reg.Register("SOUND.BUILTIN", "audio", m.soundBuiltin)
}
