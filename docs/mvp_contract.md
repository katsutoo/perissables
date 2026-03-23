# Les Perissables MVP Contract

Status: Locked for pre-repo planning
Owner: Project team
Updated: 2026-03-23

This document is the single source of truth for locked MVP scope and implementation decisions.

## Product Objective

Ship a funny, fast, multiplayer pixel-art RPG where players pick premade food characters, run story-driven adventures, and survive dice/combat events using one shared engine.

## MVP Scope (In)

1. Native desktop client (Rust + `raylib`/`raylib-rs`) and authoritative Rust server (`axum` + `tokio`).
2. Story engine driven by JSON (no story-specific hardcoded logic).
3. d100 dice system with locked rule set:
   - stat range `5..70`
   - success when `roll <= stat`
   - critical success on `000` (internal `0`)
   - critical failure on `100`
   - this `000`/`100` asymmetry is intentional for flavor even though it creates `101` discrete outcomes (`000` through `100`)
4. Basic turn-based combat integrated with story encounters.
5. Preset character roster (no character build system, no leveling).
6. Run flow: lobby -> story -> end summary -> lobby.
7. Three environment themes:
   - supermarket
   - garden
   - storage_room
8. Core multiplayer for 2-4 players with server authority.
9. Event-driven audio system with channels for ambience, music, SFX, and tiny voice barks.

## Explicit Non-Goals (Not In MVP)

1. Skill trees, progression, talent builds, class customization.
2. Advanced inventory management or deep equipment systems.
3. Procedural map generation.
4. Voice chat.
5. Mod scripting with arbitrary code execution.
6. Story-specific custom UI logic (visual skinning only in MVP).
7. Mobile, browser, or console release targets.

## Locked Runtime Constraints

- Max party size: `4`
- Target rate: `60 FPS`
- Tile size: `16x16`
- Primary resolutions for MVP QA: `1280x720`, `1920x1080`

## Locked Save Paths (v1)

- Linux: `~/.local/share/les-perissables/`
- Windows: `%AppData%/LesPerissables/`

## Locked Architecture Decisions

- HTTP/router stack: `axum` + `tower`
- Async runtime: `tokio`
- WebSocket transport: `axum` WebSockets
- Production transport: `wss://` (TLS terminated by reverse proxy)
- Local development transport: `ws://localhost`
- WebSocket payload for MVP: JSON
- Map format for MVP: TMX
- Release automation: GoReleaser for build/package generation only; it does not change licensing or grant public binary distribution rights
- Distribution policy: production desktop binaries ship through Steam depots, not public release pages
- Release targets: Linux + Windows only
- macOS policy: deferred (requires Apple Developer Program for signing/notarization workflow)
- Story schema versioning: `schema_version` integer, start at `1`, reject unsupported major versions
- Save schema versioning: `save_version` integer, support current + previous version with explicit migrators

## Locked Protocol And Limits

### WebSocket Envelope v1

All messages must use:

```json
{
  "type": "input",
  "schema_version": 1,
  "session_id": "ses_...",
  "player_id": "ply_...",
  "seq": 1,
  "payload": {}
}
```

Required fields are mandatory. Unknown message types, missing fields, stale/duplicate sequence numbers, and unsupported schema versions are rejected.

Allowed `type` values for v1:

- `join`
- `ready`
- `input`
- `state`
- `event`
- `error`
- `ping`
- `pong`
- `rejoin`
- `resync_request`
- `resync_state`

Sequence rule:

- `seq` is a per-connection monotonic counter and starts at `1`.

### Network Limits Defaults

- Max inbound WS frame size: `16 KiB`
- Max outbound WS frame size: `64 KiB`
- Input rate limit: `20 msg/s` (burst `40`) per client
- Control rate limit: `5 msg/s` (burst `10`) per client
- Heartbeat: every `10s`, disconnect after `30s` timeout

## Locked Content Conventions

### TMX

- Required layers: `ground`, `collision`, `events`
- Optional layer: `decor`
- Object naming: `evt_<id>`, `spawn_<id>`, `warp_<id>`
- Required event object property: `event_id`

### Assets

- Character frames: `16x16`
- Direction row order: down, left, right, up
- Frames per direction: `4` (`0` idle, `1-3` walk)
- Character naming pattern: `char_<group>_<name>.png`
- Theme tileset naming pattern: `tileset_<theme_id>.png`
- Music: `.ogg` (`music_<theme_id>_<track>.ogg`)
- SFX: `.wav` (`sfx_<category>_<name>.wav`)
- Voice bark clips: `.ogg` (`voice_<actor>_<line_id>.ogg`)
- Target sample rate: `48 kHz`

### Audio And Voice Scope

- Gameplay audio includes footsteps, object-use sounds, dice result sounds, spell SFX, ambience, and combat transitions.
- Voice is tiny and optional (short barks/short lines only); text remains primary.
- No in-game voice chat (players use external voice apps like Discord).
- Multiplayer audio is event-driven: server emits gameplay events, clients play local audio.

Audio channels for v1:

- `ambience`
- `music`
- `sfx`
- `voice`

## Locked Repo/Legal/CI Baseline

- Main repo (`les-perissables`): add `COPYRIGHT`, ARR `LICENSE`, `README.md` on day 1
- Stories repo (`les-perissables-stories`): add MIT `LICENSE`, `README.md` on day 1
- Minimum CI checks on first commit:
  - `cargo fmt --all --check`
  - `cargo clippy --all-targets --all-features -- -D warnings`
  - `cargo test --all-features`
  - `cargo audit`

## Privacy And Data Minimization

- No optional telemetry by default in the game client/server.
- No in-game voice chat and no voice recording storage.
- Use only gameplay-required platform/session identifiers.
- If crash upload is added, it must be opt-in and clearly documented.

## Community Website (Post-MVP)

- Build website as separate repo: `les-perissables-hub`.
- Purpose: host/discover community data packs, not runtime binaries.
- Website should have moderation policy, abuse controls, and a simple privacy page.

## Community Content And Licensing Boundary

The project has a dual model:

- `les-perissables` (main game): proprietary, All Rights Reserved
- `les-perissables-stories` (data/specs): MIT

Community creators may:

- Create/share story JSON packs
- Create/share character preset JSON packs
- Create/share UI/theme manifests and original assets

Community creators may not:

- Redistribute proprietary game code
- Redistribute proprietary built-in game assets outside the shipped game

Important runtime rule:

- Community packs are data consumed by the proprietary runtime; they do not grant rights to the engine.
- Users still need the game install/ownership to run packs.
- Multiplayer sessions require matching `pack_id`, `version`, and checksum across players.
- Pack manifest support must include `pack_id`, `version`, and checksum validation in tooling/runtime checks.

## Milestone Exit Criteria

### Vertical Slice Exit

- One playable map + one playable story
- d100 checks working
- basic combat working
- no full multiplayer requirement yet

### MVP Exit

- Authoritative multiplayer stable (2-4 players)
- lobby loop is complete
- three themes selectable through story metadata
- core run loop stable for repeated sessions

## Change Control

- Any feature outside this contract is added to post-MVP backlog.
- Scope changes only happen between phases, never inside an active phase.
- If scope grows, timeline updates must be acknowledged before coding continues.
- Pricing and other business-facing release notes belong in `docs/release_notes.md` unless they become locked MVP constraints.
