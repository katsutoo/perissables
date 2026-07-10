# Roadmap

## Delivery Milestones

These are aggressive planning estimates for one developer with controlled scope and strict phase discipline. Add `30-50%` buffer for calendar planning, especially around Steamworks, native packaging, multiplayer sync, and UGC moderation.

Core-game estimates are cumulative from Phase 00. The separate community-hub estimate is additional after or alongside the core release.

| Milestone | Included phases | Result | Estimate basis |
| --- | --- | --- | --- |
| Vertical slice | 00-07 | One map, one story, checks, first combat slice | Cumulative `3-6 weeks` |
| Playable MVP | 00-13 | Core multiplayer, three themes, stable lobby-to-run loop; reconnect hardening and durable run persistence are not complete yet | Cumulative `2-4 months` |
| Release-ready | 00-18 | Reconnect robustness, creator tooling, QA, Steam-ready builds, hosted server | Cumulative `4-8 months` |
| Community hub | 19-20 | Landing page, then accounts + pack sharing/likes/comments + moderation (separate repo) | Additional `3-8 weeks` |

## How To Use This Roadmap

- Build in order. Each phase should leave the project in a playable, testable, or clearly reviewable state.
- Do not pull later-phase work forward unless an earlier phase is blocked without it.
- Keep exact MVP locks in `docs/mvp-contract.md`; use this file for sequencing and planning. Verification methods live in the dedicated test, QA, security, and benchmark plans.
- Use `docs/progress-tracker.md` as the detailed checkbox board for completion tracking.

## Phase Plan

### Foundation (00-02)

These phases establish scope, repo shape, and the first runnable client/server shell.

### Phase 00 - Scope Freeze

Lock MVP scope, non-goals, and completion criteria before code starts.

### Phase 01 - Repo Bootstrap

Create all five locked Rust workspace members, pinned toolchain/lockfile, startup wiring, logging, legal baseline, native target compile checks, CI, and the `les-perissables-stories` repo needed by Phase 04.

### Phase 02 - Render Loop And Scene Skeleton

Get a stable game window, basic scenes, input abstraction, and audio manager shell.

### Core Gameplay Slice (03-09)

These phases build the first full single-run experience in order: move in the world, load stories, resolve checks, fight, select characters, and survive death states. Phases 05-07 may use one hardcoded debug character so the vertical slice can exercise checks and combat before Phase 08 adds the real data-driven roster; that debug path must not survive Phase 08.

### Phase 03 - Tilemap World Prototype

Load one TMX map, walk it, collide with it, and follow the player camera.

### Phase 04 - Story Schema v1

Implement the already-locked complete content schema v1, portable pack validation, conformance vectors, and canonical MIT creator examples before runtime behavior depends on them.

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

Move gameplay authority server-side using the already deterministic `game_core`; validate scripted 2-4 player convergence and all admission boundaries.

### Phase 14 - Reconnect And Sync Robustness

Add resync, heartbeat, replay protection, and disconnect recovery.

### Tooling, Persistence, And Quality (15-17)

These phases make the project easier to extend, safer to resume, and more stable to ship.

### Phase 15 - Creator Tooling

Build the `storycheck` CLI and cross-file pack checks for story creators in the MIT `les-perissables-stories` repo, on top of the schema/validation crate already built in Phase 04.

### Phase 16 - Save And Resume

Persist runs server-side with versioned session snapshots and compatibility checks so restarts and deploys do not destroy them.

### Phase 17 - QA, Balance, Performance

Provision production-equivalent staging, execute the Phase 17 pre-release QA matrix and benchmark plans against production-profile core artifacts, and profile only measured hot paths. Regression tests already belong to the phases introducing behavior; final Steam/package QA follows in Phase 18.

### Release And Community Follow-Through (18-20)

These phases cover shipping, then the community hub (a committed follow-on built soon after the core game is ready).

### Phase 18 - Steam Packaging And Production Server

Prepare and sign reproducible release builds, integrate the locked Steamworks scope, deploy the production game server, and rerun the complete QA/release matrix against the final candidate. Production binaries are distributed through Steam rather than public release pages.

### Phase 19 - Community Hub Foundation

Create the separate proprietary `les-perissables-hub` as a Rust `axum` + `maud` + `htmx` app on Railway with Railway Postgres via `sqlx`. Stage 1 is a content-only landing page. Phase 19 may implement accounts and Tier 1 UGC only behind a private/admin gate, with privacy/deletion terms, baseline abuse limits, complete validation, and immutable private storage already active. The public Stage 2 launch remains blocked on Phase 20 moderation and trust gates. Validation uses the pinned MIT crate from `les-perissables-stories`; files use the immutable R2 publication flow locked in the contract.

### Phase 20 - Moderation And Trust

Because the hub hosts user-generated content (shared packs, likes, comments), moderation and privacy are the public Stage 2 launch gate, not a later add-on: publish a moderation policy and privacy policy, add report/flag plus admin delete/ban actions, anti-spam/upload limits, and account/content deletion. Once asset/format validation is also in place, enable Tier 2 (original-asset) packs.

## Definition Of Done For Every Phase

- Phase 00 is a documentation-only scope gate: CI/build/artifact clauses below are not applicable until Phase 01 creates the workspace. Every later phase satisfies every applicable clause and records any explicit non-applicability.
- Phase-specific acceptance scenarios and expected state/events are written before implementation and linked from the tracker.
- Behavior introduced by the phase has deterministic public-interface tests for happy, boundary, invalid, and operational-error paths; each new test is observed failing for the intended reason before the fix.
- The pinned CI checks pass on a clean checkout, and supported feature/target combinations compile.
- Required manual QA uses the exact release artifact and records environment, steps, expected/actual results, side effects, cleanup, and evidence under `docs/qa-plan.md`.
- Security-sensitive boundaries have threat-model coverage and bounded negative tests under `docs/security-model.md`; no discovered credential is used.
- Performance-sensitive changes preserve correctness and run the smallest applicable release-mode measurement under `docs/benchmark-plan.md`; raw output and environment metadata are retained.
- Operational errors are explicitly propagated and stable externally; programmer invariants are deliberate; required structured logs contain no secrets or personal data.
- No blocker/critical defect or unexplained intermittent test remains. The phase checkbox links its demo/test evidence and the next phase's entry decisions are resolved.

## Known Risks

- `Steamworks`: Steam integration is conceptually simple on paper but often messy in practice, especially around auth, lobby behavior, native SDK setup, and release testing.
- `Packaging`: `raylib-rs`, native `raylib` dependencies, and cross-platform release builds are likely to cause more friction than the core game logic.
- `Multiplayer sync`: Story state, combat state, reconnect flow, and deterministic-looking client behavior can become subtle quickly once multiple players act under latency.
- `Creator tooling`: the data-driven model is a strength, but it only pays off if validation and authoring tools arrive early enough.

## Tracking Note

This document keeps the roadmap at the planning level. Use `docs/progress-tracker.md` as the built-in execution checklist, and move day-to-day implementation detail into issues or a project board if the tracker becomes too granular.
