# Architecture

## Architecture Goals

- Keep the runtime small, explicit, and easy to reason about.
- Put game authority on the server, not the client.
- Make stories, characters, and themes data-driven from the start.
- Favor a small dependency set and boring infrastructure where possible.

## Stack Overview

| Layer | Technology | Purpose |
| --- | --- | --- |
| Game client | Rust + `raylib` (`raylib-rs`) | Windowing, render loop, input, scene transitions, UI drawing, audio hooks |
| Game server | Rust + `axum` + `tokio` + `tower` | Routing, middleware, health endpoints, session orchestration |
| Game server hosting | Railway service, single active instance per environment for MVP | Staging/production deployment, WSS endpoint, health/readiness checks, structured operational logs |
| Networking | `axum` WebSockets on `tokio` | Client input transport and authoritative state/event updates |
| Serialization | `serde` + `serde_json` | Protocol payloads, story parsing, save/load data, tooling I/O |
| Story content | JSON + validator tooling | Story definitions, branching events, encounters, metadata |
| Theme system | Asset manifests + per-theme packs | Tilesets, props, ambience, music, combat backdrops, UI skin references |
| Error handling | `thiserror` + `anyhow` | Domain errors plus startup/tooling context with explicit boundaries |
| Quality baseline | Locked `cargo --locked` checks plus test, QA, security, and benchmark plans | Reproducible reliability evidence; exact commands and toolchain live in `docs/mvp-contract.md` |
| Release automation | GoReleaser | Tagged Linux/Windows builds, packaging, checksums |
| Distribution | Steam + Steamworks SDK | Authentication/ownership, lobbies, invites, and native depot distribution; achievements are post-MVP |

## System Shape

```text
[Story JSON + Theme/Character Data]
               |
               v
      [Authoritative Rust Server]
   (session state, checks, combat, sync)
               |
        WebSocket event/state flow
               |
               v
    [Rust + raylib-rs Client]
    (render, input, local audio playback)
```

## Client Responsibilities

- Run the render loop and scene transitions.
- Collect player input and turn it into intents.
- Render world, dialog, combat, and UI state.
- Play local audio in response to server-approved gameplay events.

## Server Responsibilities

- Own party state, story state, combat turns, and dice outcomes.
- Validate inputs and reject illegal transitions.
- Broadcast authoritative state snapshots and semantic events.
- Handle reconnect, resync, and multiplayer session lifecycle.
- After Phase 16, snapshot active sessions to durable storage so Release-ready restarts and deploys do not destroy runs. Before Phase 16, run loss on restart is explicit pre-release behavior.

## Headless game_core Rule

Phases 03-09 build gameplay before Phase 13 moves authority server-side. That migration stays a relocation instead of a rewrite only if `game_core` is headless from day one:

- No `raylib` types, rendering, input, or audio dependencies anywhere in `game_core`.
- No filesystem, network, wall-clock, environment-variable, or process-global access. Storage I/O stays in `server`; `game_core::save` contains versioned DTOs and pure migrators only.
- All randomness and logical time are explicit state-machine inputs. The stable RNG algorithm/version/state and logical tick are serializable so restore and deterministic tests continue exactly.
- Rules use logical durations rather than assuming the caller's update frequency. The client may call the core from a `60 Hz` presentation loop before Phase 13, but authoritative rule advancement remains compatible with the server's `20 Hz` tick.
- State advances only through explicit player intents or explicit scheduler `Tick` inputs and returns state revisions plus stable events/results; the client renders from those and never reaches into mutable game logic.

If a phase 03-09 feature is tempting to implement in the client "just for now", it goes in `game_core` behind an intent instead.

## Intended Project Structure

```text
Cargo.toml
Cargo.lock
rust-toolchain.toml
mise.toml

crates/
  client/
    src/
      main.rs
      app/
      scene/
      render/
      input/
      audio/
  server/
    src/
      main.rs
      startup.rs
      router.rs
      session/
      netcode/
  game_core/
    src/
      game/
      story/
      theme/
      combat/
      character/
      item/
      dice/
      save/
  shared/
    src/
      protocol/
      ids/
      logging.rs
  integration_tests/
    tests/
      protocol.rs
      story_runtime.rs

assets/
  themes/
  ui/
  sprites/
  audio/

stories/
  builtin/
```

The initial workspace members are `client`, `server`, `game_core`, `shared`, and `integration_tests`; their package names are `les-perissables-client`, `les-perissables-server`, `les-perissables-game-core`, `les-perissables-shared`, and `les-perissables-integration-tests`. `shared` owns protocol DTOs and IDs, `game_core` depends on `shared` plus the pinned MIT schema crate, `client` and `server` depend on `shared` and `game_core`, and `integration_tests` may depend on all workspace crates. Reverse dependencies and cycles are forbidden.

`stories/builtin/` contains ARR built-in game content. MIT schemas, conformance fixtures, creator examples, and the canonical authoring tutorial live only in `les-perissables-stories`. Downloaded community packs are installed under the locked per-user data root at runtime and are never committed under a main-repo `stories/community/` source directory.

`rust-toolchain.toml` pins Rust `1.95.0` with edition/MSRV `2024`/`1.95.0`; `Cargo.lock` is committed. `mise.toml` is local-development tooling only. CI installs Rust with `rustup`/standard Rust tooling and runs the locked `cargo --locked` checks directly rather than invoking `mise` tasks.

Workspace-level integration tests live in a dedicated `crates/integration_tests` member so Cargo runs them in CI. Per-crate tests remain in each crate's own `tests/` directory. Graphics behavior that depends on `raylib` or a real display is verified through deterministic renderer/unit seams and the release-artifact matrix in `docs/qa-plan.md` unless a headless harness is explicitly added; do not add a root-level `tests/graphics` directory that CI silently ignores.

The `storycheck` CLI and the reusable pack schema/validation rules it enforces both live in the separate MIT `les-perissables-stories` repo; the CLI is a thin front-end over that library crate. Keeping the CLI out of the proprietary repo means creators can install and run the validator without any access to game code, while the game loader and the community hub still validate packs through the same shared crate.

## Engineering Practices

- Standard library first, minimal crates beyond clear wins.
- Clear crate/module boundaries and explicit ownership.
- Operational failures return explicit errors with context, including startup/configuration failures. Panics are reserved for documented programmer-error invariants; no malformed input or unavailable dependency may trigger one.
- Use `thiserror` for domain/application errors and `anyhow` for startup/tooling context.
- Structured logging with `tracing` and `tracing-subscriber`.
- Startup owns one supervised task tree. Long-lived connection, session, persistence, heartbeat, and shutdown tasks are tracked in bounded `JoinSet`s or equivalent owners; cancellation is explicit, every join result is observed, and dropping a handle must not detach correctness-critical work. A session-task panic marks that session unavailable, records a redacted invariant failure, invokes the bounded safe-termination path, and cannot silently leave authority running elsewhere.
- Blocking native/library work uses `spawn_blocking` only behind the global bounded concurrency limits in the contract. Shutdown stops admission, signals cancellation, waits for async owners, and accounts for non-cancellable blocking work inside the `20s` deadline.
- `unsafe` is forbidden in game/domain crates. Native `raylib` and Steamworks calls stay behind small adapter modules; each unsafe block has a local `SAFETY` contract, public unsafe APIs are forbidden, callback lifetimes/thread affinity are tested, and Miri or platform sanitizers cover owned unsafe wrappers where supported.
- Follow `docs/test-strategy.md`, `docs/qa-plan.md`, `docs/security-model.md`, and `docs/benchmark-plan.md` for verification evidence and phase gates.
- Small end-to-end slices before broad feature expansion.

## Where Locks Live

This document explains the architecture, but exact locked MVP decisions such as protocol limits, release targets, persistence rules, and schema versioning live only in `docs/mvp-contract.md`.
