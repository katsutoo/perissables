# Architecture

Authority: exact workspace, toolchain, runtime, persistence, protocol, and release values derive from the corresponding locked sections of `docs/mvp-contract.md`. This document explains structure and rationale only.

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
| Serialization | `serde` + `serde_json` | Protocol payloads, story parsing, versioned persisted state, tooling I/O |
| Game persistence | SQLite through pinned bundled `rusqlite`, owned by one supervised blocking thread | Transactional active-session durability and crash recovery on the single-instance Railway volume |
| Story content | JSON + validator tooling | Story definitions, branching events, encounters, metadata |
| Theme system | Asset manifests + per-theme packs | Tilesets, props, ambience, music, combat backdrops, built-in behavior-neutral UI variant references |
| Error handling | `thiserror` + `anyhow` | Domain errors plus startup/tooling context with explicit boundaries |
| Quality baseline | Locked `cargo --locked` format, lint, test, doctest, docs, advisory, license/source-policy checks plus test, QA, security, and benchmark plans | Reproducible reliability evidence; exact commands and toolchain live in `docs/mvp-contract.md` |
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
- After Phase 16, transactionally persist active sessions in SQLite so Release-ready restarts and deploys do not destroy runs. Before Phase 16, run loss on restart is explicit pre-release behavior.

## Headless game_core Rule

Phases 03-09 build gameplay before Phase 13 moves authority server-side. That migration stays a relocation instead of a rewrite only if `game_core` is headless from day one:

- No `raylib` types, rendering, input, or audio dependencies anywhere in `game_core`.
- No filesystem, network, wall-clock, environment-variable, or process-global access. Storage I/O stays in `server`; `game_core::save` contains versioned DTOs and pure migrators only.
- All randomness and logical time are explicit state-machine inputs. The session owner supplies fresh per-run CSPRNG entropy through `StartRunEntropy`; stable RNG algorithm/version/state and logical tick are serializable so restore and deterministic tests continue exactly.
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
      persistence/
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

The exact initial workspace members, package names, and dependency direction are locked in `docs/mvp-contract.md` under "Locked Architecture Decisions." The structure above visualizes that contract; it does not redefine it.

`stories/builtin/` will contain ARR built-in game content. MIT schemas, conformance fixtures, creator examples, and the canonical authoring tutorial will live only in `les-perissables-stories` after that repository is created in Phase 01. Downloaded community packs will be installed under the locked per-user data root at runtime and never committed under a main-repo `stories/community/` source directory.

Phase 01 will add `rust-toolchain.toml` pinned to Rust `1.97.1` with edition/MSRV `2024`/`1.97.1`, set virtual-workspace `resolver = "3"`, commit `Cargo.lock`, and add local-development-only `mise.toml`. CI will install Rust with `rustup`/standard Rust tooling and run the locked `cargo --locked` checks directly rather than invoking `mise` tasks. Phase 01 also pins the bundled SQLite source through `rusqlite` and runs the contract's storage, sandbox, Steam-account/source-IP, signing, and native-build feasibility gates before later phases depend on them. These values derive from "Locked Repo/Legal/CI Baseline" and "Locked Persistence Model" in `docs/mvp-contract.md`.

Phase 01 will place workspace-level integration tests in the dedicated `crates/integration_tests` member so Cargo runs them in CI. Per-crate tests will remain in each crate's own `tests/` directory. Graphics behavior that depends on `raylib` or a real display is verified through deterministic renderer/unit seams and the release-artifact matrix in `docs/qa-plan.md` unless a headless harness is explicitly added; do not add a root-level `tests/graphics` directory that CI silently ignores.

After the separate MIT `les-perissables-stories` repository is created, the `storycheck` CLI and reusable pack schema/validation rules will live there; the CLI will be a thin front-end over that library crate. Keeping the CLI out of the proprietary repo lets creators install and run the validator without game-code access, while the game loader and future community hub validate packs through the same pinned shared crate.

## Engineering Practices

- Standard library first, minimal crates beyond clear wins.
- Clear crate/module boundaries and explicit ownership.
- Operational failures return explicit errors with context, including startup/configuration failures. Panics are reserved for documented programmer-error invariants; no malformed input or unavailable dependency may trigger one.
- Use `thiserror` for domain/application errors and `anyhow` for startup/tooling context.
- Structured logging with `tracing` and `tracing-subscriber`.
- Startup owns one supervised task tree. Long-lived connection, session, persistence, heartbeat, and shutdown tasks are tracked in bounded `JoinSet`s or equivalent owners; cancellation is explicit, every join result is observed, and dropping a handle must not detach correctness-critical work. A session-task panic marks that session unavailable, records a redacted invariant failure, invokes the bounded safe-termination path, and cannot silently leave authority running elsewhere.
- Blocking native/library work uses `spawn_blocking` only behind the global bounded concurrency limits in the contract. SQLite is the deliberate exception to per-operation blocking jobs: one supervised dedicated thread owns its connection, bounded request queue, group transactions, checkpoint, and close lifecycle. Shutdown stops admission, signals cancellation, waits for async owners and that persistence thread, and accounts for non-cancellable blocking work inside the `20s` deadline.
- `unsafe` is forbidden in game/domain crates. Native `raylib` and Steamworks calls stay behind small adapter modules; each unsafe block has a local `SAFETY` contract, public unsafe APIs are forbidden, callback lifetimes/thread affinity are tested, and Miri or platform sanitizers cover owned unsafe wrappers where supported.
- Phase 01 build documentation fixes native `raylib` and Steamworks crate/SDK sources, checksums, linkage, target prerequisites, and licenses before client work depends on a developer machine. Adapters expose safe owned Rust types and keep native handles, callbacks, and thread affinity out of `game_core` and `shared`.
- Follow `docs/test-strategy.md`, `docs/qa-plan.md`, `docs/security-model.md`, and `docs/benchmark-plan.md` for verification evidence and phase gates.
- Small end-to-end slices before broad feature expansion.

## Where Locks Live

This document explains the architecture, but exact locked MVP decisions such as protocol limits, release targets, persistence rules, and schema versioning live only in `docs/mvp-contract.md`.
