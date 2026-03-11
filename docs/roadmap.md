# Roadmap

## Delivery Milestones

These estimates assume one developer, controlled scope, and strict phase discipline.

| Milestone | Included phases | Result | Estimated time |
| --- | --- | --- | --- |
| Vertical slice | 00-07 | One map, one story, checks, first combat slice | 3-6 weeks |
| Playable MVP | 00-13 | Core multiplayer, three themes, stable lobby-to-run loop | 2-4 months |
| Release-ready | 14-18 | Reconnect robustness, creator tooling, QA, Steam-ready builds | 4-8 months |
| Community hub | 19-20 | Separate website for sharing and moderating packs | 3-8 weeks |

## How To Use This Roadmap

- Build in order. Each phase should leave the project in a playable, testable, or clearly reviewable state.
- Do not pull later-phase work forward unless an earlier phase is blocked without it.
- Keep exact MVP locks in `docs/mvp_contract.md`; use this file for sequencing and planning.
- Use `docs/progress_tracker.md` as the detailed checkbox board for completion tracking.

## Phase Plan

### Foundation (00-02)

These phases establish scope, repo shape, and the first runnable client/server shell.

### Phase 00 - Scope Freeze

Lock MVP scope, non-goals, and acceptance rules before code starts.

### Phase 01 - Repo Bootstrap

Create the Go workspace, client/server entrypoints, logging, legal files, and CI baseline.

### Phase 02 - Render Loop And Scene Skeleton

Get a stable game window, basic scenes, input abstraction, and audio manager shell.

### Core Gameplay Slice (03-09)

These phases build the first full single-run experience in order: move in the world, load stories, resolve checks, fight, select characters, and survive death states.

### Phase 03 - Tilemap World Prototype

Load one TMX map, walk it, collide with it, and follow the player camera.

### Phase 04 - Story Schema v1

Define story data structures, validation, and one canonical example story.

### Phase 05 - Story Runtime

Execute branching story progression from data only.

### Phase 06 - Dice And Checks

Implement visible d100 checks with correct critical handling and clear feedback.

### Phase 07 - Combat Core

Ship the first playable turn-based combat loop tied to story encounters.

### Phase 08 - Character Presets

Add the premade roster, validation, and character-select flow.

### Phase 09 - Death And Loot

Support lethal runs, corpse looting, and continue-until-wipe flow.

### Replayability And Presentation (10-12)

These phases make the shared engine feel broader without changing its core rules.

### Phase 10 - Theme Packs

Switch environment presentation through data-driven theme loading.

### Phase 11 - UI Variants

Support multiple UI skins while keeping behavior identical.

### Phase 12 - Lobby And Story Rotation

Complete the repeatable lobby -> run -> summary -> lobby loop.

### Multiplayer Authority And Recovery (13-14)

These phases move authority fully server-side, then harden sync and reconnect behavior.

### Phase 13 - Authoritative Multiplayer

Move gameplay authority server-side and validate 2-4 player sessions.

### Phase 14 - Reconnect And Sync Robustness

Add resync, heartbeat, replay protection, and disconnect recovery.

### Tooling, Persistence, And Quality (15-17)

These phases make the project easier to extend, safer to resume, and more stable to ship.

### Phase 15 - Creator Tooling

Build schema validation and pack-check tooling for story creators.

### Phase 16 - Save And Resume

Persist interrupted runs safely with versioned saves and compatibility checks.

### Phase 17 - QA, Balance, Performance

Tune the game, profile hot paths, and harden regression coverage.

### Release And Community Follow-Through (18-20)

These phases cover shipping, then optional community infrastructure after the core game is ready.

### Phase 18 - Steam Packaging

Prepare reproducible release builds, Steamworks integration, and release operations.

### Phase 19 - Community Hub Foundation

Create a separate web repo for pack discovery and download.

### Phase 20 - Moderation And Trust

Add moderation policy, abuse controls, and minimal admin tooling for community content.

## Definition Of Done For Every Phase

- The phase has a small demo or testable proof.
- New logic has tests where it should.
- Errors are handled explicitly.
- Logs are structured.
- Scope for the next phase is clear before coding starts.

## Known Risks

- `Steamworks`: Steam integration is conceptually simple on paper but often messy in practice, especially around auth, lobby behavior, native SDK setup, and release testing.
- `Packaging`: `raylib-go`, native dependencies, and cross-platform release builds are likely to cause more friction than the core game logic.
- `Multiplayer sync`: Story state, combat state, reconnect flow, and deterministic-looking client behavior can become subtle quickly once multiple players act under latency.
- `Creator tooling`: the data-driven model is a strength, but it only pays off if validation and authoring tools arrive early enough.

## Tracking Note

This document keeps the roadmap at the planning level. Use `docs/progress_tracker.md` as the built-in execution checklist, and move day-to-day implementation detail into issues or a project board if the tracker becomes too granular.
