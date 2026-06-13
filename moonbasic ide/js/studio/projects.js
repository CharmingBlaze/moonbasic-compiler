/**
 * moonBASIC project templates
 */

export const PROJECT_TYPES = [
  { id: 'blank', label: 'Blank project', hint: 'main.mb + assets folder' },
  { id: 'hello', label: 'Hello window', hint: 'Simple APP.OPEN demo' },
  { id: 'spin_cube', label: '3D spinning cube', hint: 'ENTITY + CAMERA + grid' },
  { id: 'pong', label: '2D Pong', hint: 'DRAW rectangles and input' }
];

const templates = {
  blank: {
    'main.mb': `; moonBASIC project
APP.OPEN(800, 600, "My Game")
APP.SETFPS(60)

WHILE NOT (INPUT.KEYDOWN(KEY_ESCAPE) OR APP.SHOULDCLOSE())
    RENDER.CLEAR(18, 22, 32)
    DRAW.TEXT("Edit main.mb and press F5 to run", 120, 280, 20, 220, 230, 245, 255)
    RENDER.FRAME()
WEND
APP.CLOSE()
`,
    'assets/.gitkeep': ''
  },
  hello: {
    'main.mb': `APP.OPEN(640, 480, "Hello moonBASIC")
APP.SETFPS(60)
WHILE NOT (INPUT.KEYDOWN(KEY_ESCAPE) OR APP.SHOULDCLOSE())
    RENDER.CLEAR(20, 24, 36)
    DRAW.TEXT("Hello, moonBASIC!", 180, 220, 28, 255, 255, 255, 255)
    RENDER.FRAME()
WEND
APP.CLOSE()
`,
    'assets/.gitkeep': ''
  },
  spin_cube: {
    'main.mb': `SetMSAA(0)
APP.OPEN(800, 600, "Spinning cube")
APP.SETFPS(60)
cam = CAMERA.CREATE().fov(55)
cube = ENTITY.CREATECUBE(1, 1, 1).scale(1.4, 1.4, 1.4).pos(0, 1, 0).col(255, 150, 70)
ang = 0.0
WHILE NOT (INPUT.KEYDOWN(KEY_ESCAPE) OR APP.SHOULDCLOSE())
    RENDER.CLEAR(38, 42, 58)
    ang = ang + APP.DELTA() * 0.4
    cube.rot(0, ang, 0)
    CAMERA.SETORBIT(cam, 0, 1, 0, 0.55, 0.32, 6.5)
    cam.Begin()
        Draw3D.Grid(14, 1.0)
        ENTITY.DRAWALL()
    cam.End()
    RENDER.FRAME()
WEND
cube.free()
cam.free()
APP.CLOSE()
`,
    'assets/.gitkeep': ''
  },
  pong: {
    'main.mb': `APP.OPEN(960, 540, "Pong")
APP.SETFPS(60)
py1 = 200.0
py2 = 200.0
bx = 480.0
by = 270.0
bvx = 280.0
bvy = 140.0
WHILE NOT (INPUT.KEYDOWN(KEY_ESCAPE) OR APP.SHOULDCLOSE())
    dt = APP.DELTA()
    IF INPUT.KEYDOWN(KEY_W) THEN py1 = py1 - 280 * dt
    IF INPUT.KEYDOWN(KEY_S) THEN py1 = py1 + 280 * dt
    IF INPUT.KEYDOWN(KEY_I) THEN py2 = py2 - 280 * dt
    IF INPUT.KEYDOWN(KEY_K) THEN py2 = py2 + 280 * dt
    bx = bx + bvx * dt
    by = by + bvy * dt
    RENDER.CLEAR(12, 14, 22)
    DRAW.RECTANGLE(30, INT(py1), 16, 80, 200, 230, 255, 255)
    DRAW.RECTANGLE(914, INT(py2), 16, 80, 255, 170, 110, 255)
    DRAW.RECTANGLE(INT(bx) - 8, INT(by) - 8, 16, 16, 255, 255, 255, 255)
    RENDER.FRAME()
WEND
APP.CLOSE()
`,
    'assets/.gitkeep': ''
  }
};

export function buildProject(typeId) {
  const files = templates[typeId] || templates.blank;
  return {
    title: PROJECT_TYPES.find(t => t.id === typeId)?.label || 'Project',
    folderName: 'my-moonbasic-game',
    files: Object.entries(files).map(([name, content]) => ({ name, content }))
  };
}

export function flattenProjectFiles(files) {
  return files;
}
