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
  storycheck/
    src/
      main.rs
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

The `storycheck` crate here is the CLI front-end. The reusable pack schema and validation rules it enforces live as an MIT library crate in the separate `les-perissables-stories` repo, so the game, `storycheck`, and the community hub all validate packs through the same code without sharing proprietary game-repo crates.

## Engineering Practices

- Standard library first, minimal crates beyond clear wins.
- Clear crate/module boundaries and explicit ownership.
- No ignored errors, no hidden panic paths outside startup.
- Use `thiserror` for domain/application errors and `anyhow` for startup/tooling context.
- Structured logging with `tracing` and `tracing-subscriber`.
- Table-driven and integration tests for parser, rules, and validation logic.
- Small end-to-end slices before broad feature expansion.

## Where Locks Live

This document explains the architecture, but exact locked MVP decisions such as protocol limits, release targets, save-path rules, and schema versioning live only in `docs/mvp-contract.md`.
