package checklistaliases

import "moonbasic/runtime"

func registerAUDIO(r runtime.Registrar) {
	r.Register("AUDIO3D.LOAD", "audio3d", forward("AUDIO.LOADSOUND"))
	r.Register("AUDIO3D.PLAYAT", "audio3d", forward("SOUND.PLAY3D"))
	r.Register("AUDIO3D.ATTACH", "audio3d", forward("SOUND.ATTACH"))
	r.Register("AUDIO3D.SETLISTENER", "audio3d", forward("AUDIO.LISTENERCAMERA"))
	r.Register("AUDIO3D.SETRANGE", "audio3d", forward("AUDIO.SETSOUNDVOLUME"))
	r.Register("AUDIO.PLAYMUSIC", "audio", forward("AUDIO.PLAY"))
	r.Register("AUDIO.PLAYSOUND", "audio", forward("AUDIO.PLAY"))
	r.Register("AUDIO.STOPSOUND", "audio", forward("AUDIO.STOP"))
	r.Register("AUDIO.STOPMUSIC", "audio", forward("AUDIO.STOP"))
	r.Register("AUDIO.SETVOLUME", "audio", forward("AUDIO.SETSOUNDVOLUME"))
}
