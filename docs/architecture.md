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
| Theme system | Asset manifests + per-theme packs | Tilesets, ambience, combat backdrops, UI skin references |
| Error handling | `thiserror` + `anyhow` | Domain errors plus startup/tooling context with explicit boundaries |
| Quality baseline | `cargo fmt --all --check`, `cargo clippy --all-targets --all-features -- -D warnings`, `cargo test --all-features`, `cargo audit` | Reliability and security hygiene (`cargo audit` requires `cargo install cargo-audit` or `cargo binstall cargo-audit`) |
| Release automation | GoReleaser | Tagged Linux/Windows builds, packaging, checksums |
| Distribution | Steam + Steamworks SDK | Lobbies, invites, achievements, native distribution |

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
- Snapshot active sessions to durable storage so restarts and deploys do not destroy runs.

## Headless game_core Rule

Phases 03-09 build gameplay before Phase 13 moves authority server-side. That migration stays a relocation instead of a rewrite only if `game_core` is headless from day one:

- No `raylib` types, rendering, input, or audio dependencies anywhere in `game_core`.
- All randomness injected (seeded RNG passed in), never created internally.
- State advances only through explicit intents and returns events/results; the client renders from those, it never reaches into game logic.

If a phase 03-09 feature is tempting to implement in the client "just for now", it goes in `game_core` behind an intent instead.

## Intended Project Structure

```text
Cargo.toml
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

assets/
  themes/
  ui/
  sprites/
  audio/

stories/
  builtin/
  community/

tests/
  integration/
  graphics/
  testutil/
```

`mise.toml` is local-development tooling only. CI installs Rust with `rustup`/standard Rust tooling and runs the locked `cargo` checks directly rather than invoking `mise` tasks.

The `storycheck` CLI and the reusable pack schema/validation rules it enforces both live in the separate MIT `les-perissables-stories` repo; the CLI is a thin front-end over that library crate. Keeping the CLI out of the proprietary repo means creators can install and run the validator without any access to game code, while the game loader and the community hub still validate packs through the same shared crate.

## Engineering Practices

- Standard library first, minimal crates beyond clear wins.
- Clear crate/module boundaries and explicit ownership.
- No ignored errors, no hidden panic paths outside startup.
- Use `thiserror` for domain/application errors and `anyhow` for startup/tooling context.
- Structured logging with `tracing` and `tracing-subscriber`.
- Table-driven and integration tests for parser, rules, and validation logic.
- Small end-to-end slices before broad feature expansion.

## Where Locks Live

This document explains the architecture, but exact locked MVP decisions such as protocol limits, release targets, persistence rules, and schema versioning live only in `docs/mvp-contract.md`.
