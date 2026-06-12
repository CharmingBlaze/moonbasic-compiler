"""Final polish checklist aliases and new commands."""

ENTRIES = [
    # APP system (aliases)
    ("APP.OPEN", ["int", "int", "string"], None, "init"),
    ("APP.CLOSE", [], None, "any"),
    ("APP.SHOULDCLOSE", [], "bool", "any"),
    ("APP.SETFPS", ["int"], None, "any"),
    ("APP.GETFPS", [], "float", "any"),
    ("APP.WIDTH", [], "int", "any"),
    ("APP.HEIGHT", [], "int", "any"),
    ("APP.TIME", [], "float", "any"),
    ("APP.DELTA", [], "float", "any"),
    ("APP.VERSION", [], "string", "any"),
    # RENDER polish names
    ("RENDER.BEGIN", ["handle"], None, "rendering"),
    ("RENDER.END", [], None, "rendering"),
    ("RENDER.SETBACKGROUND", ["int", "int", "int"], None, "rendering"),
    # ACTION checklist names
    ("ACTION.BINDKEY", ["string", "int"], None, "input"),
    ("ACTION.BINDGAMEPAD", ["string", "int", "int"], None, "input"),
    ("ACTION.HIT", ["string"], "bool", "input"),
    # INPUT polish names
    ("INPUT.MOUSEDELTA_X", [], "float", "input"),
    ("INPUT.MOUSEDELTA_Y", [], "float", "input"),
    ("INPUT.GAMEPADBUTTONDOWN", ["int", "int"], "bool", "input"),
    ("INPUT.GAMEPADAXIS", ["int", "int"], "float", "input"),
    # JSON polish names
    ("JSON.STRINGIFY", ["handle"], "string", "any"),
    ("JSON.GET", ["handle", "string"], "any", "any"),
    ("JSON.SET", ["handle", "string", "string"], None, "any"),
    # SAVE polish names
    ("SAVE.SET", ["string", "string"], None, "any"),
    ("SAVE.WRITE", ["string"], None, "any"),
    ("SAVE.READ", ["string"], None, "any"),
    ("SAVE.WRITEFILE", ["string"], None, "any"),
    ("SAVE.READFILE", ["string"], None, "any"),
    # TEXT system
    ("TEXT.DRAW", ["string", "int", "int"], None, "rendering"),
    ("TEXT.DRAWFONT", ["handle", "string", "int", "int"], None, "rendering"),
    ("TEXT.SIZE", ["string"], "int", "rendering"),
    # AUDIO3D
    ("AUDIO3D.LOAD", ["string"], "handle", "audio"),
    ("AUDIO3D.PLAYAT", ["handle", "float", "float", "float", "float"], None, "audio"),
    ("AUDIO3D.ATTACH", ["handle", "handle"], None, "audio"),
    ("AUDIO3D.SETLISTENER", ["handle"], None, "audio"),
    ("AUDIO3D.SETRANGE", ["handle", "float"], None, "audio"),
    ("AUDIO.PLAYSOUND", ["handle"], None, "audio"),
    ("AUDIO.PLAYMUSIC", ["handle"], None, "audio"),
    ("AUDIO.STOPSOUND", ["handle"], None, "audio"),
    ("AUDIO.STOPMUSIC", ["handle"], None, "audio"),
    ("AUDIO.SETVOLUME", ["handle", "float"], None, "audio"),
    # PICK polish names
    ("PICK.MOUSE", ["handle"], "bool", "any"),
    ("PICK.RAY", ["float", "float", "float", "float", "float", "float"], "bool", "any"),
    ("PICK.DISTANCE", [], "float", "any"),
    # BODY beginner API
    ("BODY.ADDSTATICBOX", ["handle", "float", "float", "float"], None, "physics3d"),
    ("BODY.ADDDYNAMICBOX", ["handle", "float", "float", "float"], None, "physics3d"),
    ("BODY.ADDSPHERE", ["handle", "float"], None, "physics3d"),
    ("BODY.ADDCAPSULE", ["handle", "float", "float"], None, "physics3d"),
    ("BODY.SETMASS", ["handle", "float"], None, "physics3d"),
    ("BODY.SETFRICTION", ["handle", "float"], None, "physics3d"),
    ("BODY.SETBOUNCE", ["handle", "float"], None, "physics3d"),
    ("BODY.APPLYFORCE", ["handle", "float", "float", "float"], None, "physics3d"),
    ("BODY.APPLYIMPULSE", ["handle", "float", "float", "float"], None, "physics3d"),
    # ASSET pack
    ("ASSET.LOADPACK", ["string"], None, "asset"),
    ("ASSET.TEXTURE", ["string"], "handle", "asset"),
    ("ASSET.MODEL", ["string"], "handle", "asset"),
    ("ASSET.SOUND", ["string"], "handle", "asset"),
    ("ASSET.UNLOAD", [], None, "asset"),
    # FILE polish
    ("FILE.DELETE", ["string"], "bool", "any"),
    ("FILE.READTEXT", ["string"], "string", "any"),
    ("FILE.WRITETEXT", ["string", "string"], None, "any"),
    # TIMER callbacks
    ("TIMER.AFTER", ["float", "string"], "int", "any"),
    ("TIMER.EVERY", ["float", "string"], "int", "any"),
    ("TIMER.CANCEL", ["int"], None, "any"),
]
