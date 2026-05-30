// One-off helper: annotate commands.json with stub messages for unavailable builtins.
package main

import (
	"encoding/json"
	"fmt"
	"os"
)

func main() {
	path := "compiler/builtinmanifest/commands.json"
	if len(os.Args) > 1 {
		path = os.Args[1]
	}
	data, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	var root map[string]any
	if err := json.Unmarshal(data, &root); err != nil {
		panic(err)
	}
	cmds, ok := root["commands"].([]any)
	if !ok {
		panic("no commands array")
	}

	stubs := map[string]string{
		"PHYSICS3D.DEBUGDRAW":         "Jolt debug wireframe overlay is not available in this release",
		"PHYSICS.SPHERECAST":          "Use PHYSICS3D.RAYCAST for ray queries",
		"PHYSICS.BOXCAST":             "Use PHYSICS3D.RAYCAST for ray queries",
		"PHYSICS.ENABLE":              "Per-body physics enable is not available in this release",
		"PHYSICS.DISABLE":             "Per-body physics disable is not available in this release",
		"JOINT3D.FIXED":               "Use JOINT3D.HINGE where possible; fixed joints are not available yet",
		"JOINT3D.SLIDER":              "Slider joints are not available in this release",
		"JOINT3D.CONE":                "Cone twist joints are not available in this release",
		"PHYSICS2D.ONCOLLISION":       "2D collision callbacks are not available in this release",
		"PHYSICS2D.PROCESSCOLLISIONS":  "2D collision event processing is not available in this release",
		"WORLD.SETREFLECTION":         "Reflection probes are not available in this release",
		"GAME.BURSTSPAWN":             "Use PARTICLE.* APIs directly; burst spawn bridge is not available yet",
		"GAME.SPRITETILEBRIDGE":       "Use collision math and sprite bounds in script",
		"CREATESPRITE3D":              "Use ENTITY.CREATECUBE + ENTITY.TEXTURE",
		"ENTITY.INSTANCE":             "Use MODEL.MAKEINSTANCED for GPU instancing",
		"LEVEL.OPTIMIZE":              "Use MODEL.MAKEINSTANCED for draw batches",
		"LEVEL.APPLYPHYSICS":          "Map Blender extras to BODY3D.* + COMMIT manually (see PHYSICS3D.md)",
		"MESH.GENERATELOD":            "Use MODEL.LOADLOD for file-based LOD",
		"MESH.GENERATELODCHAIN":       "Use MODEL.LOADLOD for file-based LOD chains",
		"MODEL.GETCHILD":              "Use ENTITY.GETCHILD for entity hierarchy",
		"CREATEMIRROR":                "Planar reflections are not available in this release",
		"FITMESH":                     "Procedural mesh fitting is not available in this release",
		"FLIPMESH":                    "Procedural mesh flip is not available in this release",
		"UPDATENORMALS":               "Procedural mesh normal updates are not available in this release",
		"ADDSURFACE":                  "Use MESH.MAKECUSTOM / ENTITY.LOADMESH (see MESH.md)",
	}

	n := 0
	for _, item := range cmds {
		obj, ok := item.(map[string]any)
		if !ok {
			continue
		}
		key, _ := obj["key"].(string)
		stub, ok := stubs[key]
		if !ok {
			continue
		}
		obj["stub"] = stub
		n++
	}
	out, err := json.MarshalIndent(root, "", "  ")
	if err != nil {
		panic(err)
	}
	out = append(out, '\n')
	if err := os.WriteFile(path, out, 0644); err != nil {
		panic(err)
	}
	fmt.Printf("annotated %d commands in %s\n", n, path)
}
