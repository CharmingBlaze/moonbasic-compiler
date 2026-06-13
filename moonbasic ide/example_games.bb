; ============================================================
; Example 1: Rotating Cube World
; ============================================================
Graphics3D 800, 600

Global cam = CreateCamera(0, 5, -15)
Global light = CreateLight(1, 0, -1, 0)
LightColor light, 255, 255, 200

Global cube = CreateCube()
PositionEntity cube, 0, 0, 0
EntityColor cube, 0, 150, 255

Global ground = CreatePlane()
ScaleEntity ground, 20, 1, 20
PositionEntity ground, 0, -1, 0
EntityColor ground, 80, 200, 80

FogMode 1
FogColor 20, 20, 40
FogRange 30, 100

MainLoop:
  TurnEntity cube, 0.5, 1.0, 0.3
  RenderWorld
  FlipBuffers
Goto MainLoop


; ============================================================
; Example 2: First Person Shooter Movement
; ============================================================
; (Uncomment to use)

; Graphics3D 800, 600

; Global cam = CreateCamera(0, 2, 0)
; Global light = CreateLight(2, 0, 10, 0)
; LightColor light, 255, 240, 200

; ; Build a simple arena
; Global floor = CreateCube()
; ScaleEntity floor, 20, 0.5, 20
; PositionEntity floor, 0, -0.25, 0
; EntityColor floor, 100, 80, 60

; Global wall1 = CreateCube()
; ScaleEntity wall1, 20, 4, 0.5
; PositionEntity wall1, 0, 2, 10
; EntityColor wall1, 150, 120, 100

; Global wall2 = CreateCube()
; ScaleEntity wall2, 20, 4, 0.5
; PositionEntity wall2, 0, 2, -10
; EntityColor wall2, 150, 120, 100

; Global speed# = 0.1
; Global rotspd# = 1.5

; MainLoop:
;   If KeyDown(17) Then TranslateEntity cam, 0, 0, speed   ; W
;   If KeyDown(31) Then TranslateEntity cam, 0, 0, -speed  ; S
;   If KeyDown(30) Then TranslateEntity cam, -speed, 0, 0  ; A
;   If KeyDown(32) Then TranslateEntity cam, speed, 0, 0   ; D
;   If KeyDown(75) Then TurnEntity cam, 0, rotspd, 0       ; Left
;   If KeyDown(77) Then TurnEntity cam, 0, -rotspd, 0      ; Right
;   If KeyDown(1) Then End
;   RenderWorld
;   FlipBuffers
; Goto MainLoop


; ============================================================
; Example 3: Solar System
; ============================================================
; (Uncomment to use)

; Graphics3D 800, 600

; Global cam = CreateCamera(0, 30, -80)
; PointCamera cam, 0, 0, 0
; Global sun_light = CreateLight(2, 0, 0, 0)
; LightColor sun_light, 255, 220, 100

; ; Sun
; Global sun = CreateSphere(32)
; ScaleEntity sun, 8, 8, 8
; EntityColor sun, 255, 200, 0
; EntityFX sun, 1  ; Fullbright

; ; Earth
; Global earth = CreateSphere(24)
; ScaleEntity earth, 3, 3, 3
; EntityColor earth, 30, 100, 220

; ; Moon
; Global moon = CreateSphere(16)
; ScaleEntity moon, 1, 1, 1
; EntityColor moon, 200, 200, 200

; Global angle# = 0
; Global moonAngle# = 0

; MainLoop:
;   angle = angle + 0.5
;   moonAngle = moonAngle + 2.0

;   PositionEntity earth, Sin(angle) * 25, 0, Cos(angle) * 25
;   PositionEntity moon, EntityX(earth) + Sin(moonAngle)*5, 0, EntityZ(earth) + Cos(moonAngle)*5

;   TurnEntity sun, 0.1, 0.2, 0
;   TurnEntity earth, 0, 1, 0

;   RenderWorld
;   FlipBuffers
; Goto MainLoop
