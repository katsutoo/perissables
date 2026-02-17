# Les Périssables

A stupid, fun, multiplayer tabletop RPG where 4 players pick premade supermarket characters — fruits, vegetables, canned goods, frozen items — and try to survive a randomly selected story set inside a grocery store. Think Dungeons & Dragons but your party is a Banana Rogue, a Canned Beans Paladin, and a Leek Bard fighting a Frozen Pizza Golem in Aisle 7.

Players don't build characters or level up. They pick from a roster, each character comes with preset stats, spells, and a personality. They explore the supermarket, hit story events, make choices, roll dice, and enter turn-based combat. If someone dies, their teammates can loot 1-2 items off the corpse and keep going. You either finish the story or everyone dies trying. When a group finishes a story, they hop back to the lobby and pick a new one.

## Visual Style

The game uses a 2D top-down perspective with pixel art, looking down at the supermarket from a bird's eye view. Players see their character sprites walking through aisles, past shelves, freezer sections, and checkout lanes. The camera follows the party as they navigate the store. Think classic Pokemon or early Zelda — tile-based movement on a grid, character sprites with simple walk animations in four directions. This perspective maps perfectly to a grocery store layout since you're essentially looking down at a floor plan. Combat transitions to a separate screen with the party facing enemies, again Pokemon-style, with turn order, spell selections, and dice roll animations playing out.

## Story System

Stories are data, not code. Each story is a JSON file defining a map, events, branching choices, dice checks, and combat encounters. The game engine is a generic "dungeon master" that can run any story file. The game ships with a handful of built-in stories at launch, with more added over time.

The story format is open source under MIT. The JSON story schema, character preset schema, UI/theme manifest schema, specification, documentation, and example packs live in a dedicated public repository licensed under MIT. The community can freely create, share, modify, and distribute their own custom stories and character data packs. Players can drop community-made pack files into the game and play them immediately. This gives the game infinite replayability without giving away any of the engine or game code. The community becomes the dungeon masters.

## Licensing

The game itself is proprietary with a public codebase. The engine, combat system, rendering, UI, multiplayer server, assets, sprites, sounds, and built-in stories are all commercially licensed under All Rights Reserved. The source code is hosted in a public repository for transparency and educational purposes, but no permission is granted to copy, modify, distribute, or use any part of the code or assets without explicit written permission. People can look, learn, open issues, and submit pull requests, but they cannot legally reuse or redistribute the game.

```
Public repo — les-perissables (All Rights Reserved):
├── COPYRIGHT
├── Game engine
├── Combat system
├── Rendering / UI
├── Multiplayer server
├── Assets (sprites, sounds, music)
└── Built-in stories

Public repo — les-perissables-stories (MIT):
├── LICENSE (MIT)
├── Story JSON schema specification
├── Character preset schema specification
├── Theme/UI manifest schema specification
├── Story validation tools
├── Example community stories
└── Documentation for story creators
```

## Community Content Policy (Explicit)

Community creators are welcome to build and share custom content, with clear boundaries:

Allowed (MIT data repo scope):

- Create/share custom story JSON files.
- Create/share custom character preset JSON files.
- Reuse supported built-in map/theme IDs in story metadata.
- Create/share custom UI skin/theme manifests and original assets.

Not allowed (ARR game repo scope):

- Redistribute or reuse proprietary game code.
- Redistribute built-in proprietary sprites, music, SFX, or other game assets outside normal gameplay use.
- Claim rights to engine/runtime internals under MIT.

Runtime boundary:

- Community packs are data consumed by the proprietary game runtime.
- Playing community packs requires owning/installing the game build.
- The MIT data/spec repository does not grant permission to copy, modify, or redistribute proprietary engine/assets.

Multiplayer compatibility rule:

- All clients in a session must have matching pack `id`, `version`, and checksum before run start.

## Tech Stack

The entire game is written in Go using raylib-go (Go bindings for raylib) for 2D rendering. The UI surface is intentionally minimal — character select screen, the top-down tile map with camera follow, a dialog box for story text and choices, a combat screen with HP bars and spell buttons, dice roll animations, and a tiny item slot per character. No complex menus, no inventory management, no stat screens. raylib-go/raylib is perfect for this scope and handles tile-based 2D rendering with ease.

### Tech Stack Schema

| Layer         | Technology                                        | Purpose                                                           |
| ------------- | ------------------------------------------------- | ----------------------------------------------------------------- |
| Game Client   | Go + raylib-go (raylib bindings)                  | Window, render loop, input, scene transitions, UI drawing         |
| Game Server   | Go + Chi (`github.com/go-chi/chi/v5`)             | Routing, middleware, health endpoints, session orchestration      |
| Networking    | `github.com/gorilla/websocket` over `net/http`    | WebSocket transport for input events and state snapshots          |
| Story Content | JSON files + schema validator                     | Data-driven stories (events, choices, checks, encounters)         |
| Theme System  | Asset manifests + per-theme packs                 | Swap tilesets, ambience, combat backdrops, and UI skins per story |
| Build/Quality | `go test`, `go vet`, `staticcheck`, `govulncheck` | Reliability, code quality, and security baseline                  |
| Release Automation | GoReleaser (`goreleaser`)                    | Versioned Linux/Windows artifacts, checksums, changelog, packaging |
| Distribution  | Steam + Steamworks SDK                            | Native Linux/Windows release, lobbies/invites/achievements        |

```text
[Story JSON + Theme Manifest]
             |
             v
      [Go Server Authority]
  (story state, dice, combat, sync)
             |
      WebSocket protocol
             |
             v
      [Go + raylib-go Client]
 (render, UI skin, input, audio)
```

### Locked Networking And HTTPS Decisions

- HTTP router (locked): **Chi** (`github.com/go-chi/chi/v5`).
- WebSocket library (locked): **gorilla/websocket** (`github.com/gorilla/websocket`).
- External transport (locked): **WSS** only in production.
- TLS/HTTPS strategy (locked): terminate TLS at reverse proxy (Caddy/Nginx/Traefik/Cloudflare) and forward to Go server on private HTTP.
- Local dev (locked): `ws://localhost` is acceptable; optional local TLS via `mkcert` when testing production-like setup.
- Message format for MVP (locked): JSON frames first; optimize to binary only if profiling shows a bottleneck.
- Release automation (locked): GoReleaser for tagged release builds.
- Release OS targets (locked): Linux + Windows only for MVP/release v1.
- macOS policy (locked): deferred unless Apple Developer Program + signing/notarization workflow is adopted.

### MVP Security Baseline (No Accounts Edition)

- Server-authoritative gameplay only (movement, dice, combat, story state).
- Strict WebSocket message validation (schema, required fields, type checks).
- Reject unknown message types and invalid state transitions.
- Enforce request/frame size limits to reduce abuse and memory pressure.
- Add per-client/session rate limits for action spam prevention.
- Use reconnect/session tokens with expiry; later bind to Steam identity at integration phase.
- Keep dependency and code scanning in CI (`go vet`, `staticcheck`, `govulncheck`).

### WebSocket Message Envelope (Locked v1)

All client/server WebSocket messages use one envelope shape:

```json
{
  "type": "input",
  "schema_version": 1,
  "session_id": "ses_01HXYZ...",
  "player_id": "ply_01HXYZ...",
  "seq": 42,
  "payload": {}
}
```

Field rules:

- `type`: required string. Allowed v1 values: `join`, `ready`, `input`, `state`, `event`, `error`, `ping`, `pong`, `rejoin`, `resync_request`, `resync_state`.
- `schema_version`: required integer. Must be `1` for v1.
- `session_id`: required string after session join.
- `player_id`: required string after player assignment.
- `seq`: required unsigned monotonic counter per connection (starts at `1`).
- `payload`: required object with message-specific data.

Validation rules:

- Reject unknown `type` values.
- Reject missing required fields.
- Reject stale or duplicated `seq` values per connection.
- Reject messages with mismatched `schema_version`.

### Network Limits Defaults (Locked v1)

| Limit | Default |
|------|---------|
| Max inbound WS frame size | `16 KiB` |
| Max outbound WS frame size | `64 KiB` |
| Input message rate | `20 msg/s` (burst `40`) per client |
| Control message rate (`join`/`rejoin`/`ready`) | `5 msg/s` (burst `10`) per client |
| Heartbeat interval | every `10s` |
| Heartbeat timeout | disconnect after `30s` without valid heartbeat |

These defaults can be tuned later from config, but v1 implementation should start with these values.

### TMX Conventions (Locked v1)

Required map layers:

- `ground` (tile layer, required)
- `collision` (tile layer, required; non-empty tile means blocked)
- `events` (object layer, required)
- `decor` (tile layer, optional)

Object naming and properties:

- Event trigger objects: `evt_<event_id>`
- Spawn objects: `spawn_<spawn_id>`
- Warp objects: `warp_<warp_id>`
- Every event object must define custom property `event_id` matching story node/trigger IDs.
- Grid alignment required for collidable and trigger objects.

### Asset Conventions (Locked v1)

Sprites:

- Character frame size: `16x16`.
- Character sheet row order: row `0` down, row `1` left, row `2` right, row `3` up.
- Frames per direction: `4` (frame `0` idle, frames `1-3` walk cycle).
- Naming pattern: `char_<group>_<name>.png` (example: `char_fruit_banana_rogue.png`).
- Theme tileset naming pattern: `tileset_<theme_id>.png`.

Audio:

- Music format: `.ogg` (loop-friendly, compressed).
- SFX format: `.wav` (PCM, short effects).
- Voice bark format: `.ogg` (short compressed one-liners).
- Target sample rate: `48 kHz`.
- Naming pattern: `music_<theme_id>_<track>.ogg`, `sfx_<category>_<name>.wav`, `voice_<actor>_<line_id>.ogg`.

### Audio And Voice Scope (Locked v1)

Audio channels:

- `ambience`: environment loop while exploring (supermarket/garden/storage_room).
- `music`: exploration/combat tracks.
- `sfx`: interactions and gameplay feedback.
- `voice`: short character/enemy barks.

Required gameplay audio events in MVP:

- Movement and interaction: footsteps, object-use sounds (example: toy tank/cart interactions).
- Dice: unique sounds for `success`, `failure`, `critical_success`, `critical_failure`.
- Combat: battle music transition on encounter start/end.
- Spells: special cast SFX per spell (with fallback to generic magic SFX).
- Enemies/characters: tiny voice barks on selected combat/dialog moments.

Voice line constraints (MVP):

- Keep voice content tiny (short words/sentences only; no full VO system).
- Text bubble/dialog remains primary source of narrative.
- Voice playback is optional flavor and must never block gameplay.

Voice chat policy:

- No in-game voice chat for MVP/release v1.
- Players use Discord or other external voice apps.

Multiplayer audio rule:

- Server sends gameplay events only; clients play local audio from those events (no audio streaming over WebSocket).

### Repo Split And Legal Baseline (Locked)

Main repo (`les-perissables`, proprietary / ARR):

- Required day-1 legal files: `COPYRIGHT`, `LICENSE` (All Rights Reserved text), `README.md`.
- Contains engine, client/server code, assets, built-in stories.

Stories repo (`les-perissables-stories`, MIT):

- Required day-1 legal files: `LICENSE` (MIT), `README.md`.
- Contains story schema, validator tooling, examples, creator docs.

### CI Starter Baseline (Locked)

First commit CI must run:

- `go test ./...`
- `go vet ./...`
- `staticcheck ./...`
- `go tool govulncheck ./...`

Workflow target file: `.github/workflows/ci.yml`.

### Save File Paths (Locked v1)

Use standard per-OS user data locations:

- Linux saves/config: `~/.local/share/les-perissables/`
- Windows saves/config: `%AppData%/LesPerissables/`

Path rules:

- Keep all run data under one game root folder per OS.
- Use atomic save writes (`tmp` + rename) in this folder.
- Keep a backup save slot for corruption recovery.

### Privacy And Data Minimization (Locked v1)

Principle: collect the minimum data required to run the game.

- Game runtime: no optional telemetry by default.
- Multiplayer identity: platform/session identifiers required for gameplay and anti-cheat only.
- No in-game voice recording or voice chat storage.
- If crash upload is enabled later, make it opt-in and document exactly what is sent.

### Community Website Strategy (Post-MVP, Separate Repo)

Plan community story sharing as a separate web project repository.

- Recommended repo: `les-perissables-hub` (separate from game runtime repo).
- Purpose: upload/share/discover packs (`story`, `character`, `theme` data).
- Website should host metadata and user-created content; game runtime stays proprietary.
- Apply moderation policy and terms before public submissions open.
- Keep website data minimal (only what is needed for accounts/moderation if enabled).

Why Gorilla as default:

- Team familiarity and cleaner learning path for MVP implementation.
- Documentation/examples are easy to follow for rapid execution.

## Multiplayer

Multiplayer runs over WebSockets. A Go server holds the authoritative game state — map positions, combat turns, dice rolls, story progression. All game logic is server-side so nobody can cheat their rolls. Clients just render what the server tells them.

## Platform

The game ships exclusively on Steam with native builds for both Linux and Windows. Go cross-compiles effortlessly between the two platforms, and raylib-go (with underlying raylib) supports both out of the box. Steam integration via Steamworks SDK handles multiplayer lobbies, friend invites, and achievements. No browser version — just clean native binaries. Priced at 2-3 euros.

macOS note: a Mac build can technically be uploaded to Steam, but a smooth user experience on modern macOS generally requires code signing and notarization, which requires an Apple Developer Program membership.

## Story Engine

The story engine is a simple state machine. It loads a JSON story file, places events on the map, presents choices to players, resolves dice checks against character stats, and triggers combat encounters. It doesn't know anything about frozen aisles or forbidden spices — all the flavor lives in the story data files. This means adding a new story is just writing a new JSON file with no code changes. Community stories go through the same engine as built-in ones — no difference.

## Combat

Combat is turn-based and dice-driven. When an encounter triggers, the screen transitions to a dedicated combat view — party on one side, enemies on the other, Pokemon-style. Each character has preset spells and stats. On their turn a player picks an action, dice are rolled, damage or effects are resolved. Enemies are defined per encounter in the story file. Simple, fast, lethal. No grinding, no progression — just survive or don't.

## Dice Rules (d100)

The game uses a percentile system with a d100 roll.

- Character stats range from **5** to **70** (70 is current max).
- Roll lower than or equal to the relevant stat: **success**.
- Roll higher than the relevant stat: **failure**.
- Roll **000**: **critical success**.
- Roll **100**: **critical failure**.

This applies to story checks and combat checks, unless a specific effect says otherwise.

Engine evaluation order for a check:

1. If roll is `000` (stored internally as `0`), result is critical success.
2. Else if roll is `100`, result is critical failure.
3. Else if roll is lower than or equal to stat, result is success.
4. Else result is failure.

Example story JSON check (d100):

```json
{
  "id": "check_open_backdoor",
  "type": "check",
  "label": "Force the rusty backdoor",
  "check": {
    "stat": "strength",
    "difficulty": 0,
    "dice": "d100",
    "critical_success": 0,
    "critical_failure": 100,
    "success_rule": "roll_lte_stat"
  },
  "outcomes": {
    "critical_success": {
      "goto": "node_backdoor_bursts_open",
      "effects": ["grant_item:crowbar"]
    },
    "success": {
      "goto": "node_backdoor_open"
    },
    "failure": {
      "goto": "node_backdoor_stuck",
      "effects": ["apply_status:embarrassed"]
    },
    "critical_failure": {
      "goto": "node_alarm_triggers",
      "effects": ["start_encounter:security_drones"]
    }
  }
}
```

## Implementation Roadmap (Go + raylib-go)

This roadmap is intentionally split into small phases so we can ship step by step without writing too much code per prompt.

### Delivery Milestones And Time Estimates

These estimates assume one developer, steady execution, and controlled scope.

| Milestone | Included Phases | Target Result | Estimated Time |
|-----------|-----------------|---------------|----------------|
| Vertical Slice | Phase 00 -> Phase 07 | 1 map, 1 story, d100 checks, basic combat, no full multiplayer | ~3-6 weeks |
| Playable MVP | Phase 00 -> Phase 13 | Core multiplayer, 3 themes (supermarket/garden/storage room), lobby run loop, stable core flow | ~2-4 months |
| Polished Release-Ready | Phase 14 -> Phase 18 | Reconnect robustness, creator tooling, QA/balance, Steam-ready builds | ~4-8 months |
| Community Web Hub (optional) | Phase 19 -> Phase 20 | Separate website repo for sharing/moderating community packs | ~3-8 weeks |

Schedule guardrails:

- If scope grows, timeline grows linearly.
- If you stay phase-locked (no mid-phase feature additions), timeline stays realistic.
- Finish criteria are the checkbox tasks in the tracker section.

### Important Decisions To Lock Before Coding

1. **Single engine, many themes:** Keep one gameplay engine for all stories. Only change assets, map tilesets, audio, and UI skin per story.
2. **Data-driven first:** Story behavior stays in JSON. Theme and UI variant selection must also come from JSON metadata.
3. **No custom logic per story (at first):** Avoid story-specific code paths until after MVP. This prevents complexity explosion.
4. **Small vertical slices:** Each phase should end with something playable or testable.

### Go Practices We Will Follow (adapted from AGENTS_GO.md)

- Stdlib-first mindset, minimal dependencies, explicit code.
- Constructor injection and clear package boundaries.
- No ignored errors, no panic outside startup, structured logging with `slog`.
- `context.Context` first param for network/server/service operations.
- Focused functions (generally under one screen), extract helpers early.
- Table-driven tests for parser/logic systems.
- CI checks: `go test`, `go vet`, `staticcheck`, `govulncheck`.

### Proposed Project Structure (game-oriented)

```text
cmd/
  client/
    main.go
  server/
    main.go

internal/
  app/                # App wiring and lifecycle
  scene/              # Scene manager (lobby, world, combat)
  render/             # raylib-go draw wrappers, camera, atlas
  input/              # Input mapping
  netcode/            # WebSocket protocol, client sync
  game/               # Runtime game state
  story/              # Story schema, loader, validator
  theme/              # Theme packs (tileset, UI skin, audio)
  combat/             # Turn system and combat rules
  character/          # Preset character definitions
  item/               # Loot and inventory slots
  dice/               # Dice service and RNG
  logx/               # slog setup helpers

assets/
  themes/
    supermarket/
    garden/
    storage_room/
  ui/
  sprites/
  audio/

stories/
  builtin/
  community/

test/
  integration/          # Multi-package integration tests (server/session/story flow)
  graphics/             # raylib-go smoke/visual tests (run with `-tags=graphics`)
  testutil/             # Shared fixtures, builders, test helpers

docs/
  mvp_contract.md       # Locked MVP scope and non-goals
  phase_acceptance.md   # Reusable phase completion template

COPYRIGHT
LICENSE
```

### Phase 00 - Foundation And Scope Freeze

**Goal:** Freeze MVP boundaries and avoid feature creep.

**Build:**
- Write short MVP list and explicit non-goals.
- Lock core runtime constraints:
  - max party size: `4`
  - target simulation/render rate: `60 FPS`
  - world tile size: `16x16`
  - design/output resolutions for MVP: `1280x720` and `1920x1080`
- Define coding rules and done checklist for every phase.

**Done when:**
- Team agrees on one-page MVP contract.
- New ideas are tracked as post-MVP, not merged into active phases.

**Phase 00 output artifacts:**
- `docs/mvp_contract.md`
- `docs/phase_acceptance.md`
- `COPYRIGHT`
- `LICENSE`

### Phase 01 - Repo Bootstrap

**Goal:** Create clean Go workspace for client and server.

**Build:**
- Initialize module, folder structure, and package boundaries.
- Add `justfile` tasks for run/test/lint/security checks.
- Add basic config loading and `slog` initialization.
- Add legal baseline files (`COPYRIGHT`, ARR `LICENSE`) from day 1.
- Add starter CI workflow with minimum checks (`go test`, `go vet`, `staticcheck`, `govulncheck`).

**Done when:**
- `go run ./cmd/client` and `go run ./cmd/server` both start cleanly.
- Baseline CI/local checks pass.

### Phase 02 - Render Loop And Scene Skeleton

**Goal:** Have a stable game window and scene switching.

**Build:**
- raylib-go window boot, fixed update loop, draw loop.
- Scene manager with `LobbyScene`, `WorldScene`, `CombatScene` placeholders.
- Basic keyboard input abstraction.
- Audio manager skeleton with channels (`ambience`, `music`, `sfx`, `voice`) and master/music/sfx/voice volume controls.

**Done when:**
- You can switch scenes with debug keys and keep stable FPS.

### Phase 03 - Tilemap World Prototype

**Goal:** Walk in a top-down map with camera follow.

**Build:**
- Load one tilemap (TMX export adapter for MVP).
- Collision layers and blocked tiles.
- Character movement with 4-direction animation.
- Trigger footsteps and exploration ambience playback in world scene.

**Done when:**
- A player can move around a store map and collide with walls/shelves.

### Phase 04 - Story Schema v1 (Core)

**Goal:** Define and validate story JSON contract.

**Build:**
- Story schema for: metadata, map reference, events, choices, checks, encounters.
- Validation package with clear errors (line/path aware if possible).
- Example story file loaded at startup.

**Done when:**
- Invalid JSON fails with readable errors.
- Valid story loads into in-memory structs.

### Phase 05 - Story Runtime State Machine

**Goal:** Execute story progression from data only.

**Build:**
- Event trigger system (position and scripted trigger IDs).
- Dialogue + choice UI with branching.
- Story state save structure (current node, flags, completed events).

**Done when:**
- One full mini-story can be completed without hardcoded event logic.

### Phase 06 - Dice And Checks

**Goal:** Resolve checks consistently and visibly.

**Build:**
- Dice service (`d100`, modifiers, and critical handling for `000`/`100`).
- Character stat check resolution from story nodes.
- Lightweight dice roll animation and result log line.
- Dice result SFX mapping (`success`, `failure`, `critical_success`, `critical_failure`).

**Done when:**
- Choice checks branch correctly based on roll + stat.

### Phase 07 - Combat Core v1

**Goal:** First playable turn-based combat slice.

**Build:**
- Turn order, action selection, damage resolution.
- Enemy definitions loaded from story encounter data.
- Simple combat UI: HP bars, action buttons, combat log.
- Combat music transition (exploration -> battle -> exploration) and spell cast SFX hooks.
- Tiny character/enemy combat voice bark hooks with cooldown (text remains primary).

**Done when:**
- Entering encounter from story launches combat and returns outcome.

### Phase 08 - Character Roster And Presets

**Goal:** Add premade funny characters cleanly.

**Build:**
- Character data files: stats, spells, traits, sprite refs.
- Character select scene and party lock-in.
- Server-safe validation for allowed roster choices.
- Add optional character voice bark references and text fallback keys.

**Done when:**
- Players can pick from presets and start a run.

### Phase 09 - Death, Loot, And Run Continuation

**Goal:** Support lethal runs without ending instantly.

**Build:**
- Character death state and removal from active turns.
- Corpse loot rule (1-2 items) and tiny inventory slots.
- Continue-until-all-dead or story-complete flow.

**Done when:**
- A run survives one character death and can still finish.

### Phase 10 - Theme Packs (Supermarket, Garden, Storage Room)

**Goal:** Enable different environments per story with minimal code changes.

**Build:**
- Add `theme_id` to story metadata.
- Theme pack format: tileset, props, ambience loop, exploration music, combat music, combat backdrop, UI skin id.
- Asset loader that swaps theme resources at story start.
- Add third environment pack: `storage_room`.
- Add per-theme audio manifest + fallback rules for missing assets.

**Done when:**
- Story A loads supermarket visuals; Story B loads garden visuals; Story C loads storage room visuals.

### Phase 11 - UI Variants Per Story

**Goal:** Make UI look different per story while logic stays shared.

**Build:**
- Add `ui_variant` (or link from theme pack).
- Skinning system for frames, buttons, fonts, colors, dialog portraits.
- Keep same UI behaviors and controls across variants.

**Done when:**
- At least 2 distinct UI skins render correctly with the same game logic.

### Phase 12 - Lobby And Story Rotation

**Goal:** Finish run, return to lobby, pick another story fast.

**Build:**
- Lobby flow with story list and preview metadata.
- End-of-run summary and reset hooks.
- Quick restart loop for party replayability.

**Done when:**
- Players can complete story, return to lobby, and launch another theme/story.

### Phase 13 - Multiplayer Authoritative Server

**Goal:** Move game logic authority server-side.

**Build:**
- WebSocket protocol (messages for input, state snapshots, events).
- Server-authoritative state for movement, story progression, combat turns, dice.
- Client as renderer/input sender only.
- Strict message validation and unknown-message rejection.
- Frame size limits and per-client rate limiting.
- Session/rejoin token lifecycle (issue, validate, expire).

**Done when:**
- 2-4 players can play one run with synchronized state.

### Phase 14 - Sync Robustness And Reconnect

**Goal:** Make multiplayer resilient.

**Build:**
- Sequence numbers / tick IDs for state updates.
- Reconnect flow (rejoin session, state resync).
- Basic anti-cheat sanity checks (invalid move/input rejection).
- Ping/pong heartbeat with idle timeout handling.
- Duplicate/replay protection for stale or repeated input frames.

**Done when:**
- A dropped client can reconnect without corrupting session state.

### Phase 15 - Content Tooling For Story Creators

**Goal:** Make story creation easy and safe.

**Build:**
- CLI validator for story + character + theme references.
- Pack manifest support (`pack_id`, `version`, checksum) for multiplayer compatibility checks.
- Docs for schema, examples, and common mistakes.
- Optional dry-run mode to simulate event graph reachability.

**Done when:**
- A non-programmer can create and validate a story/character pack.

### Phase 16 - Save System And Session Persistence

**Goal:** Support resume for interrupted sessions.

**Build:**
- Save current run state (party, flags, map positions, encounter state).
- Load/resume flow with compatibility checks on story version.
- Corruption-safe save writes (temp + atomic rename).
- Use locked OS-specific save roots for file placement.

**Done when:**
- Mid-run quit and resume works reliably.

### Phase 17 - QA, Balance, And Performance

**Goal:** Make it fun, stable, and smooth.

**Build:**
- Playtest pass on pacing, fairness, and "fun/stupid" tone.
- Profile update/render hot paths; optimize draw calls and allocations.
- Regression tests for story parser, branching, combat math.

**Done when:**
- Stable FPS on target hardware and no blocker bugs in core loop.

### Phase 18 - Steam Packaging And Release Readiness

**Goal:** Prepare shipping builds for Linux + Windows.

**Build:**
- GoReleaser config for reproducible Linux/Windows artifacts (with CGO/raylib-go + raylib build matrix strategy).
- Build scripts, assets packaging, version stamping.
- Steamworks integration pass (lobbies, invites, achievements as scoped).
- Bind multiplayer identity to Steam auth/session tickets for production trust.
- Release checklist, crash log collection, hotfix protocol.

**Done when:**
- RC build is reproducible and publish-ready.

### Phase 19 - Community Web Hub Foundation (Separate Repo)

**Goal:** Launch a minimal website where players can publish/discover packs.

**Build:**
- Create separate repo for web hub (`les-perissables-hub`).
- Implement pack listing pages (story/character/theme packs).
- Implement pack metadata schema (title, author, version, checksum, tags).
- Provide safe download workflow for data packs only (no executables).
- Link game docs and validator usage for creators.

**Done when:**
- Players can browse and download validated community packs from the web hub.

### Phase 20 - Community Moderation And Trust

**Goal:** Keep shared content safe, lawful, and low-abuse.

**Build:**
- Publish moderation guidelines and submission rules.
- Add report/flag flow for inappropriate or infringing content.
- Add minimal admin moderation tools (hide/unlist/remove).
- Add abuse controls (upload limits, basic anti-spam, checksum verification).
- Publish a simple privacy page for website data handling.

**Done when:**
- Community submissions can be moderated with clear policy and low operational risk.

### Definition Of Done For Every Phase

- A small demo or test proves the phase works end-to-end.
- No ignored errors, no hidden panics, logs are structured.
- New logic has tests (table-driven where relevant).
- Scope for the next phase is explicit before coding starts.

### Notes On Complexity (your question about different environments/UI)

This is **very possible** and not too tricky if we keep these guardrails:

1. One engine and one rule system.
2. Per-story theme and UI are data-driven asset swaps.
3. No story-specific custom UI behavior in MVP.

If we respect those three rules, supermarket, garden, and storage room stories can feel different visually while staying simple to build and maintain.

## Clickable Progress Tracker

Use this checklist as your execution board. Tick boxes as you finish each item.

### Pre-Implementation Lock Checklist
- [x] LOCK-1 Map format selected: `TMX`
- [x] LOCK-2 WebSocket package selected: `github.com/gorilla/websocket`
- [x] LOCK-3 HTTP router selected: `github.com/go-chi/chi/v5`
- [x] LOCK-4 Production transport fixed: `wss://` with reverse-proxy TLS termination
- [x] LOCK-5 WebSocket payload format fixed for MVP: JSON
- [x] LOCK-6 Story/theme schema versioning fixed: `schema_version` integer, start at `1`, reject unsupported major versions
- [x] LOCK-7 Save-file versioning/migration fixed: `save_version` integer, support current + previous version with explicit migrators
- [x] LOCK-8 Release targets fixed: Linux + Windows only; macOS deferred (Apple signing/notarization cost)
- [x] LOCK-9 WS envelope fixed: `type`, `schema_version`, `session_id`, `player_id`, `seq`, `payload`
- [x] LOCK-10 Network limits fixed: frame caps, per-client rate limits, heartbeat interval/timeout
- [x] LOCK-11 TMX conventions fixed: layer names + object naming/property rules
- [x] LOCK-12 Asset conventions fixed: sprite sheet/frame order/naming + audio formats
- [x] LOCK-13 Repo split/legal baseline fixed: ARR main repo + MIT stories repo with legal files day 1
- [x] LOCK-14 CI baseline fixed: `go test`, `go vet`, `staticcheck`, `govulncheck`
- [x] LOCK-15 Audio scope fixed: ambience/music/sfx/voice channels + gameplay audio events
- [x] LOCK-16 Voice chat policy fixed: no in-game voice chat (external apps only)
- [x] LOCK-17 Community content/licensing boundary fixed: MIT data packs + ARR runtime/assets
- [x] LOCK-18 Save path conventions fixed: Linux `~/.local/share/les-perissables/`, Windows `%AppData%/LesPerissables/`
- [x] LOCK-19 Privacy baseline fixed: minimum data, no default telemetry, opt-in crash upload if added later
- [x] LOCK-20 Community website plan fixed: separate repo (`les-perissables-hub`) in post-MVP phases

### Phase 00 - Foundation And Scope Freeze
- [x] 00.1 Write `docs/mvp_contract.md` with goals, non-goals, and "not in MVP" list
- [x] 00.2 Freeze core constraints: party size (`4`), tile size (`16x16`), target FPS (`60`), target resolutions (`1280x720`, `1920x1080`)
- [x] 00.3 Freeze dice/check rules (`d100`, stat range `5-70`, `000` crit success, `100` crit fail)
- [x] 00.4 Freeze networking scope for MVP (hosted server, no peer-to-peer)
- [x] 00.5 Create `docs/phase_acceptance.md` template for all future phases
- [x] 00.6 Phase 00 complete

### Phase 01 - Repo Bootstrap
- [ ] 01.1 Initialize module and root folders (`cmd/`, `internal/`, `assets/`, `stories/`, `docs/`)
- [ ] 01.2 Create `cmd/client/main.go` and `cmd/server/main.go` with startup wiring only
- [ ] 01.3 Add logging bootstrap (`internal/logx`) using `slog` structured logs
- [ ] 01.4 Add `justfile` commands: `run-client`, `run-server`, `test`, `lint`, `security-scan`
- [ ] 01.5 Add baseline checks (`go test ./...`, `go vet ./...`, `staticcheck`, `govulncheck`)
- [ ] 01.6 Add legal files (`COPYRIGHT`, ARR `LICENSE`) in main repo scaffold
- [ ] 01.7 Add starter CI workflow at `.github/workflows/ci.yml` with locked baseline checks
- [ ] 01.8 Phase 01 complete

### Phase 02 - Render Loop And Scene Skeleton
- [ ] 02.1 Create fixed-timestep game loop (`update` + `draw`) in client app layer
- [ ] 02.2 Add scene interface (`Enter`, `Update`, `Draw`, `Exit`)
- [ ] 02.3 Implement placeholder scenes: lobby, world, combat
- [ ] 02.4 Add scene router and debug scene switching keys
- [ ] 02.5 Add framerate/debug overlay to confirm stable loop behavior
- [ ] 02.6 Add audio manager with channel buses (`ambience`, `music`, `sfx`, `voice`) and volume controls
- [ ] 02.7 Phase 02 complete

### Phase 03 - Tilemap World Prototype
- [x] 03.1 Map format locked for MVP: `TMX`
- [ ] 03.2 Load and render one sample supermarket map with layers
- [ ] 03.3 Implement blocked tile collision and movement rejection
- [ ] 03.4 Implement 4-direction player movement and sprite facing/animation
- [ ] 03.5 Add camera follow and map bounds clamp
- [ ] 03.6 Add footsteps + exploration ambience playback in world scene
- [ ] 03.7 Phase 03 complete

### Phase 04 - Story Schema v1 (Core)
- [ ] 04.1 Define schema files for story metadata, events, choices, checks, encounters
- [ ] 04.2 Implement Go structs in `internal/story` matching schema fields
- [ ] 04.3 Implement loader with strict validation and helpful path-based errors
- [ ] 04.4 Add table-driven tests for valid/invalid story files
- [ ] 04.5 Add one canonical example story in `stories/builtin/`
- [ ] 04.6 Phase 04 complete

### Phase 05 - Story Runtime State Machine
- [ ] 05.1 Implement runtime state struct (current node, flags, completed events)
- [ ] 05.2 Implement trigger resolver (tile position + trigger ID)
- [ ] 05.3 Implement dialogue panel and branching choice handling
- [ ] 05.4 Implement side effects (set/unset flags, start encounter)
- [ ] 05.5 Add tests for branching paths and unreachable-node detection
- [ ] 05.6 Phase 05 complete

### Phase 06 - Dice And Checks
- [ ] 06.1 Implement d100 roll generator with internal roll range `0..100`
- [ ] 06.2 Implement resolver ordering: `0` critical success, `100` critical failure, then normal success/fail
- [ ] 06.3 Implement stat-based checks (`roll <= stat`), stat cap enforcement (`5..70`)
- [ ] 06.4 Add UI formatting: display internal `0` as `000`
- [ ] 06.5 Add table-driven tests covering boundaries (`0`, `1`, stat, stat+1, `100`)
- [ ] 06.6 Add dice result SFX mapping (`success`, `failure`, `critical_success`, `critical_failure`)
- [ ] 06.7 Phase 06 complete

### Phase 07 - Combat Core v1
- [ ] 07.1 Implement combat state machine (start, player turn, enemy turn, resolution)
- [ ] 07.2 Implement action set (`attack`, `spell`, `item`, `pass`) with costs/effects
- [ ] 07.3 Implement hit/check resolution using d100 rules
- [ ] 07.4 Load encounter definitions from story data only (no hardcoded fights)
- [ ] 07.5 Return structured outcome to story runtime (win/loss/rewards)
- [ ] 07.6 Add combat music transition hooks (enter/exit combat)
- [ ] 07.7 Add spell cast SFX and tiny enemy/character combat voice bark hooks with cooldown
- [ ] 07.8 Phase 07 complete

### Phase 08 - Character Roster And Presets
- [ ] 08.1 Define character preset schema (stats, spells, traits, sprite IDs)
- [ ] 08.2 Implement character loader and validation (`stat <= 70`)
- [ ] 08.3 Build character select UI with lock-in flow
- [ ] 08.4 Enforce duplicate/invalid pick rules server-side
- [ ] 08.5 Add sample roster pack (fruit/vegetable/canned/frozen archetypes)
- [ ] 08.6 Add optional voice bark references + text fallback keys in character data
- [ ] 08.7 Phase 08 complete

### Phase 09 - Death, Loot, And Run Continuation
- [ ] 09.1 Add character death state and remove dead actors from turn queue
- [ ] 09.2 Add corpse interaction with loot limit (`1-2` items)
- [ ] 09.3 Add tiny inventory slot constraints and item transfer rules
- [ ] 09.4 Implement run continuation logic (alive party continues)
- [ ] 09.5 Implement wipe logic (all dead -> run failed summary)
- [ ] 09.6 Phase 09 complete

### Phase 10 - Theme Packs (Supermarket, Garden, Storage Room)
- [ ] 10.1 Add `theme_id` to story metadata schema and validator
- [ ] 10.2 Define theme manifest format (tiles, props, ambience loop, exploration music, combat music, combat backdrop, default skin)
- [ ] 10.3 Implement theme asset loader/unloader with missing-asset fallbacks
- [ ] 10.4 Create `assets/themes/supermarket`, `assets/themes/garden`, and `assets/themes/storage_room`
- [ ] 10.5 Verify story switch applies full environment swap for all 3 themes without logic changes
- [ ] 10.6 Add per-theme audio manifest + fallback rules (including shared default SFX)
- [ ] 10.7 Phase 10 complete

### Phase 11 - UI Variants Per Story
- [ ] 11.1 Add `ui_variant` support in schema or derive it from theme manifest
- [ ] 11.2 Define skin tokens (frame sprites, button sprites, font refs, color tokens)
- [ ] 11.3 Refactor UI draw code to read tokens instead of hardcoded values
- [ ] 11.4 Implement at least two complete skins and fallback behavior
- [ ] 11.5 Verify identical controls/behavior across all skins
- [ ] 11.6 Phase 11 complete

### Phase 12 - Lobby And Story Rotation
- [ ] 12.1 Build lobby scene that lists available built-in stories
- [ ] 12.2 Show story preview metadata (title, estimated duration, theme)
- [ ] 12.3 Add end-of-run summary screen (result, deaths, key events)
- [ ] 12.4 Add reset hooks that clear transient run state safely
- [ ] 12.5 Confirm loop: finish run -> lobby -> pick next story -> launch
- [ ] 12.6 Phase 12 complete

### Phase 13 - Multiplayer Authoritative Server
- [ ] 13.1 Define network protocol messages (join, ready, input, state, event, error)
- [ ] 13.2 Implement server session lifecycle and lobby-to-run transition
- [ ] 13.3 Move all authority server-side (movement, story state, combat, dice)
- [ ] 13.4 Implement client intent messages only (never trust client outcomes)
- [ ] 13.5 Add strict payload validation + unknown-message rejection
- [ ] 13.6 Enforce max frame size and per-client/session rate limits
- [ ] 13.7 Implement session/rejoin token issue/validate/expiry flow
- [ ] 13.8 Validate 2-4 player synchronization in one local test session
- [ ] 13.9 Emit semantic audio cue events only (clients play local audio; no streamed audio payloads)
- [ ] 13.10 Phase 13 complete

### Phase 14 - Sync Robustness And Reconnect
- [ ] 14.1 Add sequence numbers/tick IDs to state updates
- [ ] 14.2 Add reconnect handshake (`session_id`, `player_id`, auth/rejoin token)
- [ ] 14.3 Implement full state resync on reconnect before inputs are accepted
- [ ] 14.4 Add ping/pong heartbeat and idle timeout disconnect rules
- [ ] 14.5 Add server-side sanity checks for illegal movement/actions
- [ ] 14.6 Add duplicate/replay input protection using sequence validation
- [ ] 14.7 Add chaos tests (disconnect/reconnect/late packets) for stability
- [ ] 14.8 Phase 14 complete

### Phase 15 - Content Tooling For Story Creators
- [ ] 15.1 Create `cmd/storycheck` CLI for story/character/theme schema and reference validation
- [ ] 15.2 Validate cross-file references (story -> character -> theme -> assets)
- [ ] 15.3 Add pack manifest checks (`pack_id`, `version`, checksum)
- [ ] 15.4 Add `--dry-run` graph walk for branching reachability and dead ends
- [ ] 15.5 Write creator docs with minimal "first story" tutorial + licensing boundaries
- [ ] 15.6 Add example packs and common error troubleshooting section
- [ ] 15.7 Phase 15 complete

### Phase 16 - Save System And Session Persistence
- [ ] 16.1 Define save schema with version field and migration strategy
- [ ] 16.2 Serialize run state (party, flags, position, current node, encounter)
- [ ] 16.3 Implement atomic save writes (temp file + fsync + rename)
- [ ] 16.4 Use locked per-OS save roots (Linux `~/.local/share/les-perissables/`, Windows `%AppData%/LesPerissables/`)
- [ ] 16.5 Implement load/resume validation (story ID/version compatibility)
- [ ] 16.6 Add corruption handling (backup slot + user-facing recovery message)
- [ ] 16.7 Phase 16 complete

### Phase 17 - QA, Balance, And Performance
- [ ] 17.1 Create playtest checklist for pacing, difficulty, and comedy tone
- [ ] 17.2 Collect balancing data (success rates per stat/check type)
- [ ] 17.3 Profile client/server hotspots and reduce avoidable allocations
- [ ] 17.4 Add regression tests for parser, branching, combat math, and dice edges
- [ ] 17.5 Close blocker/critical bugs and verify no regressions
- [ ] 17.6 Phase 17 complete

### Phase 18 - Steam Packaging And Release Readiness
- [ ] 18.1 Add GoReleaser config for Linux/Windows release automation
- [ ] 18.2 Add reproducible Linux/Windows build scripts and version stamping
- [ ] 18.3 Package runtime assets and verify path handling in release builds
- [ ] 18.4 Integrate scoped Steamworks features (lobbies/invites/achievements)
- [ ] 18.5 Bind production multiplayer identity to Steam auth/session tickets
- [ ] 18.6 Add crash log/reporting path and hotfix playbook
- [ ] 18.7 Run release candidate smoke tests on both target OSes
- [ ] 18.8 Phase 18 complete

### Phase 19 - Community Web Hub Foundation (Separate Repo)
- [ ] 19.1 Create separate repository `les-perissables-hub`
- [ ] 19.2 Add pack listing pages and metadata model (`pack_id`, `version`, checksum, tags)
- [ ] 19.3 Add safe pack download flow (data packs only)
- [ ] 19.4 Add creator docs links + validator integration guidance
- [ ] 19.5 Phase 19 complete

### Phase 20 - Community Moderation And Trust
- [ ] 20.1 Publish moderation policy and submission rules
- [ ] 20.2 Add report/flag flow and basic admin moderation actions
- [ ] 20.3 Add anti-spam/abuse controls and upload limits
- [ ] 20.4 Publish website privacy policy (minimal data handling)
- [ ] 20.5 Phase 20 complete
