package handlecall

import (
	"strings"

	"moonbasic/vm/heap"
)

// TagFromNamespace maps a runtime namespace (e.g. "CAMERA", "MODEL") to the heap
// type tag used for handle-method dispatch. Returns 0, false when unknown.
func TagFromNamespace(ns string) (uint16, bool) {
	switch strings.ToUpper(strings.TrimSpace(ns)) {
	case "CAMERA":
		return heap.TagCamera, true
	case "CAMERA2D":
		return heap.TagCamera2D, true
	case "RENDERTARGET", "RENDERTEXTURE":
		return heap.TagRenderTexture, true
	case "INSTANCE", "INSTANCEDMODEL":
		return heap.TagInstancedModel, true
	case "TRANSFORM", "MAT4", "MATRIX":
		return heap.TagMatrix, true
	case "BODY3D", "PHYSICS3D":
		return heap.TagPhysicsBody, true
	case "PHYSICS2D", "BODY2D":
		return heap.TagBody2D, true
	case "PEER":
		return heap.TagPeer, true
	case "AUDIO":
		// Sound/Music share AUDIO.*; default to Sound for CREATE-shaped returns.
		return heap.TagSound, true
	case "AUDIOSTREAM":
		return heap.TagAudioStream, true
	case "WAVE":
		return heap.TagWave, true
	case "TEXTURE":
		return heap.TagTexture, true
	case "IMAGE":
		return heap.TagImage, true
	case "MESH":
		return heap.TagMesh, true
	case "MATERIAL":
		return heap.TagMaterial, true
	case "SHADER":
		return heap.TagShader, true
	case "SKY":
		return heap.TagSky, true
	case "CLOUD":
		return heap.TagCloud, true
	case "WEATHER":
		return heap.TagWeather, true
	case "DECAL":
		return heap.TagDecal, true
	case "FONT":
		return heap.TagFont, true
	case "SPRITE":
		return heap.TagSprite, true
	case "LIGHT":
		return heap.TagLight, true
	case "LIGHT2D":
		return heap.TagLight2D, true
	case "PARTICLE":
		return heap.TagParticle, true
	case "WATER":
		return heap.TagWater, true
	case "TERRAIN":
		return heap.TagTerrain, true
	case "TILEMAP":
		return heap.TagTilemap, true
	case "ATLAS":
		return heap.TagAtlas, true
	case "CHARCONTROLLER", "CHARACTER", "CHARACTERREF":
		return heap.TagCharController, true
	case "ENTITY":
		return heap.TagEntityRef, true
	case "SHAPE", "SHAPEREF":
		return heap.TagShape, true
	case "KINEMATIC", "KINEMATICBODY":
		return heap.TagKinematicBody, true
	case "STATICBODY":
		return heap.TagStaticBody, true
	case "TRIGGER", "TRIGGERBODY":
		return heap.TagTriggerBody, true
	case "MODEL", "LODMODEL":
		return heap.TagModel, true
	case "FILE":
		return heap.TagFile, true
	case "ARRAY":
		return heap.TagArray, true
	case "NAV":
		return heap.TagNav, true
	case "NAVAGENT":
		return heap.TagNavAgent, true
	case "PATH":
		return heap.TagPath, true
	case "SCATTER":
		return heap.TagScatterSet, true
	case "PROP":
		return heap.TagProp, true
	case "TWEEN":
		return heap.TagTween, true
	case "BIOME":
		return heap.TagBiome, true
	case "NOISE":
		return heap.TagNoise, true
	case "TABLE":
		return heap.TagTable, true
	case "POOL":
		return heap.TagPool, true
	case "JSON":
		return heap.TagJSON, true
	case "CSV":
		return heap.TagCSV, true
	case "DB":
		return heap.TagDB, true
	case "ROWS":
		return heap.TagDBRows, true
	case "RNG":
		return heap.TagRng, true
	case "MEM":
		return heap.TagMem, true
	case "LOBBY":
		return heap.TagLobby, true
	case "NETPACKET", "PACKET":
		return heap.TagNetPacket, true
	case "HOST":
		return heap.TagHost, true
	case "EVENT":
		return heap.TagEvent, true
	case "SPRITEGROUP":
		return heap.TagSpriteGroup, true
	case "SPRITELAYER":
		return heap.TagSpriteLayer, true
	case "SPRITEBATCH":
		return heap.TagSpriteBatch, true
	case "SPRITEUI":
		return heap.TagSpriteUI, true
	case "PARTICLE2D":
		return heap.TagParticle2D, true
	case "QUAT", "QUATERNION":
		return heap.TagQuaternion, true
	case "COLOR":
		return heap.TagColor, true
	case "VEC2":
		return heap.TagVec2, true
	case "VEC3":
		return heap.TagVec3, true
	case "PLAYER2D":
		return heap.TagPlayer2D, true
	case "TIMER", "GAMETIMER":
		return heap.TagGameTimer, true
	case "STOPWATCH":
		return heap.TagGameStopwatch, true
	case "BTREE":
		return heap.TagBTree, true
	case "STEER", "STEERGROUP":
		return heap.TagSteerGroup, true
	case "COMPUTE", "COMPUTESHADER":
		return heap.TagComputeShader, true
	case "SHADERBUFFER":
		return heap.TagShaderBuffer, true
	case "JOINT2D":
		return heap.TagJoint2D, true
	case "DRAW3D", "DRAWPRIM3D":
		return heap.TagDrawPrim3D, true
	case "DRAW2D", "DRAWPRIM2D":
		return heap.TagDrawPrim2D, true
	case "TEXTDRAW":
		return heap.TagTextDraw, true
	case "MOVER":
		return heap.TagMoverFacade, true
	case "GRID", "TACTICALGRID":
		return heap.TagTacticalGrid, true
	default:
		return 0, false
	}
}

// SplitRegistryKey splits "NS.METHOD" into namespace and method.
func SplitRegistryKey(key string) (ns, method string, ok bool) {
	key = strings.ToUpper(strings.TrimSpace(key))
	i := strings.IndexByte(key, '.')
	if i <= 0 || i == len(key)-1 {
		return "", "", false
	}
	return key[:i], key[i+1:], true
}
