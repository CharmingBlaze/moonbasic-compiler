/** moonBASIC built-in examples */

export const EXAMPLES = {
  spin_cube: {
    title: 'Spinning 3D Cube',
    category: '3D',
    code: `; moonBASIC — spinning cube (3D)
; Requires full runtime: moonrun main.mb

SetMSAA(0)
APP.OPEN(800, 600, "moonBASIC — spinning cube")
APP.SETFPS(60)

cam = CAMERA.CREATE().fov(55)
cube = ENTITY.CREATECUBE(1, 1, 1).scale(1.4, 1.4, 1.4).pos(0.0, 1.0, 0.0).col(255, 150, 70)

camYaw = 0.55
camPitch = 0.32
cdist = 6.5
ang = 0.0

WHILE NOT (INPUT.KEYDOWN(KEY_ESCAPE) OR APP.SHOULDCLOSE())
    RENDER.CLEAR(38, 42, 58)
    dt = APP.DELTA()
    ang = ang + dt * 0.4
    cube.rot(ang * 0.35, ang, ang * 0.2)
    CAMERA.SETORBIT(cam, 0.0, 1.0, 0.0, camYaw, camPitch, cdist)
    cam.Begin()
        Draw3D.Grid(14, 1.0)
        ENTITY.DRAWALL()
    cam.End()
    Draw.Text("Spinning cube   ESC quit", 12, 10, 18, 235, 240, 255, 255)
    RENDER.FRAME()
WEND

cube.free()
cam.free()
APP.CLOSE()
`
  },
  pong: {
    title: '2D Pong',
    category: '2D',
    code: `; moonBASIC — two-player Pong
APP.OPEN(960, 540, "moonBASIC Pong")
APP.SETFPS(60)

pw = 16
ph = 80
phh = 40
pyMin = 40
pyMax = 420
py1 = 200.0
py2 = 200.0
bx = 480.0
by = 270.0
bvx = 280.0
bvy = 140.0
bCap = 480.0
p1s = 0
p2s = 0

WHILE NOT (INPUT.KEYDOWN(KEY_ESCAPE) OR APP.SHOULDCLOSE())
    dt = APP.DELTA()
    ps = 280.0 * dt
    IF INPUT.KEYDOWN(KEY_W) THEN py1 = py1 - ps
    IF INPUT.KEYDOWN(KEY_S) THEN py1 = py1 + ps
    IF INPUT.KEYDOWN(KEY_I) THEN py2 = py2 - ps
    IF INPUT.KEYDOWN(KEY_K) THEN py2 = py2 + ps
    py1 = MATH.CLAMP(py1, pyMin, pyMax)
    py2 = MATH.CLAMP(py2, pyMin, pyMax)
    bx = bx + bvx * dt
    by = by + bvy * dt
    IF by < 24.0 OR by > 516.0 THEN bvy = 0.0 - bvy
    RENDER.CLEAR(12, 14, 22)
    DRAW.RECTANGLE(30, INT(py1), pw, ph, 200, 230, 255, 255)
    DRAW.RECTANGLE(914, INT(py2), pw, ph, 255, 170, 110, 255)
    DRAW.RECTANGLE(INT(bx) - 8, INT(by) - 8, 16, 16, 255, 255, 255, 255)
    DRAW.TEXT("P1: " + STR(p1s) + "    P2: " + STR(p2s), 360, 10, 24, 245, 245, 250, 255)
    RENDER.FRAME()
WEND
APP.CLOSE()
`
  },
  gui_basics: {
    title: 'GUI Basics',
    category: 'UI',
    code: `; moonBASIC — raygui window, label, button
APP.OPEN(800, 600, "GUI basics")
APP.SETFPS(60)
GUI.ENABLE()

status = "Press OK or close the panel"
WHILE NOT (INPUT.KEYDOWN(KEY_ESCAPE) OR APP.SHOULDCLOSE())
    RENDER.CLEAR(28, 28, 36)
    IF GUI.WINDOWBOX(40, 40, 520, 320, "Demo Panel") THEN
        status = "Panel closed"
    ENDIF
    GUI.LABEL(56, 72, 460, 22, "GUI.LABEL — immediate-mode UI")
    IF GUI.BUTTON(56, 100, 120, 28, "OK") THEN
        status = "OK pressed"
    ENDIF
    GUI.LABEL(56, 140, 460, 22, status)
    DRAW.TEXT("ESC to quit", 16, 580, 16, 180, 180, 190, 255)
    RENDER.FRAME()
WEND
APP.CLOSE()
`
  },
  hello: {
    title: 'Hello moonBASIC',
    category: 'Starter',
    code: `; Your first moonBASIC program
APP.OPEN(640, 480, "Hello moonBASIC")
APP.SETFPS(60)

WHILE NOT (INPUT.KEYDOWN(KEY_ESCAPE) OR APP.SHOULDCLOSE())
    RENDER.CLEAR(20, 24, 36)
    DRAW.TEXT("Hello, moonBASIC!", 180, 220, 28, 255, 255, 255, 255)
    DRAW.TEXT("Press ESC to quit", 220, 260, 18, 180, 190, 210, 255)
    RENDER.FRAME()
WEND
APP.CLOSE()
`
  }
};

export const EXAMPLE_CATEGORIES = ['Starter', '2D', '3D', 'UI'];

export const SNIPPETS = {
  game_loop: {
    title: 'Game loop',
    code: `WHILE NOT (INPUT.KEYDOWN(KEY_ESCAPE) OR APP.SHOULDCLOSE())
    dt = APP.DELTA()
    RENDER.CLEAR(20, 24, 36)
    ; update & draw here
    RENDER.FRAME()
WEND
`
  },
  function_block: {
    title: 'Function',
    code: `FUNCTION Add(a, b)
    RETURN a + b
ENDFUNCTION
`
  },
  entity_3d: {
    title: '3D entity + camera',
    code: `cam = CAMERA.CREATE().fov(55).pos(0, 4, -8)
cube = ENTITY.CREATECUBE(1, 1, 1).pos(0, 1, 0).col(200, 120, 80, 255)
`
  }
};
