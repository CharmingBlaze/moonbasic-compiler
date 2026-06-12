// Apidoc writes docs/API_CONSISTENCY.md from compiler/builtinmanifest/commands.json.
// Run from the module root: go run ./tools/apidoc
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type root struct {
	Commands []struct {
		Key         string   `json:"key"`
		Args        []string `json:"args"`
		Returns     string   `json:"returns,omitempty"`
		Namespace   string   `json:"namespace,omitempty"`
		Description string   `json:"description,omitempty"`
	} `json:"commands"`
}

type sig struct {
	key, ret, desc string
	args           []string
}

func main() {
	path := filepath.Join("compiler", "builtinmanifest", "commands.json")
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "read %s: %v (run from repository root)\n", path, err)
		os.Exit(1)
	}
	var r root
	if err := json.Unmarshal(data, &r); err != nil {
		fmt.Fprintf(os.Stderr, "parse: %v\n", err)
		os.Exit(1)
	}
	byNS := make(map[string][]sig)
	for _, c := range r.Commands {
		// Skip legacy suffixed names in the human documentation
		if strings.HasSuffix(c.Key, "$") || strings.HasSuffix(c.Key, "#") ||
			strings.HasSuffix(c.Key, "%") || strings.HasSuffix(c.Key, "?") {
			continue
		}

		parts := strings.SplitN(c.Key, ".", 2)
		ns := parts[0]
		if c.Namespace != "" {
			ns = strings.ToUpper(c.Namespace)
		}
		byNS[ns] = append(byNS[ns], sig{key: c.Key, ret: c.Returns, args: c.Args, desc: c.Description})
	}
	nss := make([]string, 0, len(byNS))
	for ns := range byNS {
		nss = append(nss, ns)
	}
	sort.Strings(nss)

	var b strings.Builder
	b.WriteString("# moonBASIC API consistency\n\n")
	b.WriteString("This document is generated from `compiler/builtinmanifest/commands.json`.\n\n")
	b.WriteString("**Contributor contract:** Treat this file as the authoritative checklist of **registered overloads** (name, arity, and manifest metadata). New builtins belong in **`compiler/builtinmanifest/commands.json`**; refresh this doc after manifest edits so tooling, reviews, and external contributors stay aligned.\n\n")
	b.WriteString("Refresh: `go run ./tools/apidoc` (from the repository root).\n\n")
	b.WriteString("## Related documentation\n\n")
	b.WriteString("- **[ERROR_MESSAGES.md](../ERROR_MESSAGES.md)** — compile-time vs runtime errors, did-you-mean, heap handle hints.\n")
	b.WriteString("- **[ROADMAP.md](../ROADMAP.md)** — phased engineering plan (polish → rendering → 2D → systems → …).\n")
	b.WriteString("- **[COMMAND_AUDIT.md](../COMMAND_AUDIT.md)** — namespace → primary `docs/reference/*.md` file; run **`go run ./tools/cmdaudit`** to verify every manifest namespace maps to an existing reference page (exit code **2** if a namespace is unmapped or a referenced file is missing).\n")
	b.WriteString("- **[reference/API_CONVENTIONS.md](../reference/API_CONVENTIONS.md)** — consistent verbs (`LOAD`, `SETPOS`, `SETSCALE`, …) across object types.\n\n")
	b.WriteString("## Naming conventions\n\n")
	b.WriteString("- **Registry / source form**: `NS.ACTION` in uppercase with a dot (e.g. `CAMERA.SETPOS`).\n")
	b.WriteString("- **Handle methods** (on a handle value): `cam.SetPos` dispatches to `CAMERA.SETPOS`. **`SetPosition`** is an alias for **`SetPos`** where both are registered (same handler).\n")
	b.WriteString("- **Spatial handles** (`Camera3D`, `Body3D`, `Model`, `Sprite`, `Light2D`): use **`SETPOS`** for position. Aliases **`SETPOSITION`** exist for **Camera**, **Model**, **Body3D**, **Sprite**, **Light2D** — same implementation as `SETPOS`.\n")
	b.WriteString("- **3D lights** (`LIGHT.*`): use **`LIGHT.SETDIR`** for the directional sun (normalized). **`LIGHT.SETPOS`** stores point/spot position; **`LIGHT.SETTARGET`** moves the shadow frustum look-at; **`RENDER.SETAMBIENT`** sets PBR ambient tint.\n")
	b.WriteString("- **`MODEL.SETPOS`**: sets the model root transform to a **translation matrix** (replaces prior rotation/scale on that matrix).\n")
	b.WriteString("- **Creation verbs**: prefer **`*.CREATE`** for procedural handles; deprecated **`*.MAKE`** aliases point at the same handlers where registered. **`*.LOAD`** for assets (`SPRITE.LOAD`, `MODEL.LOAD`); materials use `MATERIAL.MAKEDEFAULT` / `MATERIAL.MAKEPBR`.\n")
	b.WriteString("- **Cross-type patterns**: see **[API_CONVENTIONS.md](../reference/API_CONVENTIONS.md)**.\n\n")
	b.WriteString("## Default values (common no-arg `CREATE` paths)\n\n")
	b.WriteString("| Command | Defaults |\n")
	b.WriteString("|----------|----------|\n")
	b.WriteString("| `CAMERA.CREATE` (deprecated `CAMERA.MAKE`) | position (0, 2, 8), target (0, 0, 0), up (0, 1, 0), FOV 45°, perspective |\n")
	b.WriteString("| `LIGHT.CREATE` (deprecated `LIGHT.MAKE`) | kind `directional`, white, intensity 1.0, direction toward normalized (-1,-2,-1) |\n")
	b.WriteString("| `BODY3D.CREATE` (deprecated `BODY3D.MAKE`) | no args → **DYNAMIC** motion type |\n")
	b.WriteString("| `MATERIAL.MAKEDEFAULT` / `MAKEPBR` | see `runtime/mbmodel3d` (material modules) |\n\n")
	b.WriteString("## Debug watch overlay\n\n")
	b.WriteString("`DEBUG.WATCH(label, value)` stores rows; `DEBUG.WATCHCLEAR` clears them. With **CGO** and Raylib, the window pipeline may draw a **top-left overlay** each frame (`runtime/mbdebug/overlay_cgo.go`) when **`DEBUG.ENABLE`** was called or the host enabled **`Registry.DebugMode`** (e.g. **`--debug`**). **`DEBUG.DISABLE`** clears the user override. Without CGO, watches are stored but not drawn.\n\n")
	b.WriteString("## Errors\n\n")
	b.WriteString("- **Compile-time**: unknown `NS.METHOD` → did-you-mean within namespace + manifest listing (see `compiler/semantic/cmdhint.go`).\n")
	b.WriteString("- **Runtime**: VM wraps native errors with **source file and line** when available (`vm/vm.go`). Unknown registry keys → `runtime.FormatUnknownRegistryCommand`.\n\n")
	b.WriteString("## Commands by namespace\n\n")
	for _, ns := range nss {
		entries := byNS[ns]
		sort.Slice(entries, func(i, j int) bool {
			if entries[i].key != entries[j].key {
				return entries[i].key < entries[j].key
			}
			return len(entries[i].args) < len(entries[j].args)
		})
		b.WriteString("### ")
		b.WriteString(ns)
		b.WriteString("\n\n")
		for _, e := range entries {
			b.WriteString("- **`")
			b.WriteString(e.key)
			b.WriteString("`**")
			if len(e.args) > 0 {
				b.WriteString(" - args: ")
				b.WriteString(strings.Join(e.args, ", "))
			} else {
				b.WriteString(" - args: (none)")
			}
			if e.ret != "" {
				b.WriteString(" -> returns ")
				b.WriteString(e.ret)
			}
			if strings.TrimSpace(e.desc) != "" {
				b.WriteString(" — ")
				b.WriteString(strings.TrimSpace(e.desc))
			}
			b.WriteString("\n")
		}
		b.WriteString("\n")
	}
	out := filepath.Join("docs", "API_CONSISTENCY.md")
	if err := os.WriteFile(out, []byte(b.String()), 0o644); err != nil {
		fmt.Fprintf(os.Stderr, "write %s: %v\n", out, err)
		os.Exit(1)
	}
	fmt.Fprintf(os.Stderr, "wrote %s\n", out)

	if err := writeSystemsCommandRegistry(byNS, nss); err != nil {
		fmt.Fprintf(os.Stderr, "systems registry: %v\n", err)
		os.Exit(1)
	}
}

type systemGroup struct {
	title    string
	guide    string
	anchor   string
	namespaces []string
}

func writeSystemsCommandRegistry(byNS map[string][]sig, allNS []string) error {
	groups := []systemGroup{
		{"Core: window, time, render, scene, entity", "01-CORE.md", "core-window-time", []string{"WINDOW", "TIME", "SYSTEM", "RENDER", "SCENE", "ENTITY"}},
		{"Camera and light", "02-CAMERA-LIGHT.md", "camera-light", []string{"CAMERA", "CAMERA2D", "LIGHT", "LIGHT2D"}},
		{"Meshes, models, materials, textures, asset packs", "03-ASSETS.md", "assets", []string{"MESH", "MODEL", "MATERIAL", "TEXTURE", "ASSET", "BBOX", "BSPHERE"}},
		{"Input and actions", "04-INPUT.md", "input-action", []string{"INPUT", "ACTION", "GAMEPAD", "CURSOR", "GESTURE"}},
		{"Physics, bodies, collision, picking", "05-PHYSICS.md", "physics", []string{"PHYSICS", "PHYSICS3D", "BODY", "BODY2D", "BODY3D", "BODYREF", "COLLISION", "PICK", "RAY", "RAY2D"}},
		{"Audio (2D and 3D)", "06-AUDIO.md", "audio", []string{"AUDIO", "AUDIOSTREAM", "SOUND"}},
		{"2D sprites, tilemaps, terrain, particles, animation", "07-2D-WORLD.md", "2d-world", []string{"SPRITE", "TILEMAP", "TERRAIN", "PARTICLE", "ANIM", "WORLD", "CHUNK"}},
		{"UI, fonts, and text", "08-UI-TEXT.md", "ui-text", []string{"GUI", "FONT", "DRAW", "TEXT", "COLOR"}},
		{"Save, files, JSON, math, vectors", "09-DATA.md", "data", []string{"SAVE", "FILE", "JSON", "MATH", "VEC3", "VEC2", "CONFIG", "CSV"}},
		{"Debug, timers", "10-DEBUG-TIMER.md", "debug-timer", []string{"DEBUG", "TIMER"}},
	}

	seen := make(map[string]bool)
	for _, g := range groups {
		for _, ns := range g.namespaces {
			seen[ns] = true
		}
	}

	var b strings.Builder
	b.WriteString("# Game systems — complete command registry\n\n")
	b.WriteString("> Every registered command for the **40 beginner systems** and their related namespaces.\n\n")
	b.WriteString("Generated from `compiler/builtinmanifest/commands.json` (same source as [API_CONSISTENCY.md](../API_CONSISTENCY.md)).\n\n")
	b.WriteString("**How to use this page:**\n\n")
	b.WriteString("- Learn **why** and **when** to call commands in [00-START.md](00-START.md) and the numbered guides (`01-CORE` … `11-TOOLING`).\n")
	b.WriteString("- Look up **arity** and return types here or in [API_CONSISTENCY.md](../API_CONSISTENCY.md).\n")
	b.WriteString("- Deep behavior: [COMMAND_AUDIT.md](../COMMAND_AUDIT.md) → `docs/reference/*.md`.\n")
	b.WriteString("- Validate a script: `moonbasic --check yourgame.mb`.\n")
	b.WriteString("- In-game help: `HELP(\"ENTITY.SETPOS\")`.\n\n")
	b.WriteString("**Case:** Command names are **case-insensitive** in source.\n\n")
	b.WriteString("---\n\n## Table of contents\n\n")
	for _, g := range groups {
		b.WriteString("- [")
		b.WriteString(g.title)
		b.WriteString("](#")
		b.WriteString(g.anchor)
		b.WriteString(")\n")
	}
	b.WriteString("- [Globals and language builtins](#globals)\n")
	b.WriteString("- [All other engine namespaces](#all-other-namespaces)\n\n")
	b.WriteString("---\n\n")

	total := 0
	for _, g := range groups {
		b.WriteString("## ")
		b.WriteString(g.title)
		b.WriteString("\n\n")
		b.WriteString("Guide: [")
		b.WriteString(g.guide)
		b.WriteString("](")
		b.WriteString(g.guide)
		b.WriteString(")\n\n")
		groupCount := 0
		for _, ns := range g.namespaces {
			entries := byNS[ns]
			if len(entries) == 0 {
				continue
			}
			sort.Slice(entries, func(i, j int) bool {
				if entries[i].key != entries[j].key {
					return entries[i].key < entries[j].key
				}
				return len(entries[i].args) < len(entries[j].args)
			})
			b.WriteString("### ")
			b.WriteString(ns)
			b.WriteString("\n\n")
			for _, e := range entries {
				writeSigLine(&b, e)
				groupCount++
			}
			b.WriteString("\n")
		}
		total += groupCount
		b.WriteString("*")
		b.WriteString(fmt.Sprintf("%d", groupCount))
		b.WriteString(" overloads in this section.*\n\n---\n\n")
	}

	b.WriteString("## Globals and language builtins\n\n")
	b.WriteString("Guide: [11-TOOLING.md](11-TOOLING.md) · [LANGUAGE.md](../LANGUAGE.md)\n\n")
	globalNS := []string{"PRINT", "HELP", "IMPORT", "ABS", "SIN", "COS", "RAND", "LEN", "STR", "VAL", "CHR", "ASC"}
	globalCount := 0
	for _, ns := range globalNS {
		entries := byNS[ns]
		if len(entries) == 0 {
			continue
		}
		sort.Slice(entries, func(i, j int) bool { return entries[i].key < entries[j].key })
		b.WriteString("### ")
		b.WriteString(ns)
		b.WriteString("\n\n")
		for _, e := range entries {
			writeSigLine(&b, e)
			globalCount++
		}
		b.WriteString("\n")
	}
	b.WriteString("Language keywords (`IF`, `WHILE`, `FUNCTION`, `IMPORT \"file.mb\"`, …) are documented in [LANGUAGE.md](../LANGUAGE.md).\n\n")
	b.WriteString("CLI (`moonbasic new`, `moonrun`, `moonbasic --check`, …): [11-TOOLING.md](11-TOOLING.md).\n\n---\n\n")

	b.WriteString("## All other engine namespaces\n\n")
	b.WriteString("The full engine registers **")
	b.WriteString(fmt.Sprintf("%d", len(allNS)))
	b.WriteString("** dotted namespaces and thousands of overloads. Everything is listed in [API_CONSISTENCY.md](../API_CONSISTENCY.md).\n\n")
	b.WriteString("Namespace → reference file map: [COMMAND_AUDIT.md](../COMMAND_AUDIT.md).\n\n")
	b.WriteString("| Namespace | Overloads | Primary reference |\n")
	b.WriteString("|-------------|----------:|-------------------|\n")
	// COMMAND_AUDIT has counts; here we only list namespaces not in beginner groups
	otherCount := 0
	for _, ns := range allNS {
		if seen[ns] || ns == "PRINT" || ns == "HELP" {
			continue
		}
		skipGlobal := false
		for _, g := range globalNS {
			if ns == g {
				skipGlobal = true
				break
			}
		}
		if skipGlobal {
			continue
		}
		n := len(byNS[ns])
		otherCount += n
		b.WriteString("| `")
		b.WriteString(ns)
		b.WriteString("` | ")
		b.WriteString(fmt.Sprintf("%d", n))
		b.WriteString(" | [API_CONSISTENCY.md](../API_CONSISTENCY.md#")
		b.WriteString(strings.ToLower(ns))
		b.WriteString(") |\n")
	}
	b.WriteString("\n*Beginner-system overloads above: ")
	b.WriteString(fmt.Sprintf("%d", total+globalCount))
	b.WriteString(" · Other namespaces: ")
	b.WriteString(fmt.Sprintf("%d", otherCount))
	b.WriteString(" · See [API_CONSISTENCY.md](../API_CONSISTENCY.md) for the complete machine-readable list.*\n")

	out := filepath.Join("docs", "systems", "COMMAND_REGISTRY.md")
	if err := os.WriteFile(out, []byte(b.String()), 0o644); err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "wrote %s\n", out)
	return nil
}

func writeSigLine(b *strings.Builder, e sig) {
	b.WriteString("- **`")
	b.WriteString(e.key)
	b.WriteString("`**")
	if len(e.args) > 0 {
		b.WriteString(" — args: ")
		b.WriteString(strings.Join(e.args, ", "))
	} else {
		b.WriteString(" — args: (none)")
	}
	if e.ret != "" {
		b.WriteString(" → ")
		b.WriteString(e.ret)
	}
	if strings.TrimSpace(e.desc) != "" {
		b.WriteString(" — ")
		b.WriteString(strings.TrimSpace(e.desc))
	}
	b.WriteString("\n")
}
