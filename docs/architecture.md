# Architecture

## Architecture Goals

- Keep the runtime small, explicit, and easy to reason about.
- Put game authority on the server, not the client.
- Make stories, characters, and themes data-driven from the start.
- Favor a small dependency set and boring infrastructure where possible.

## Stack Overview

| Layer | Technology | Purpose |
| --- | --- | --- |
| Game client | Go + `raylib-go` | Windowing, render loop, input, scene transitions, UI drawing, audio hooks |
| Game server | Go + `chi` | Routing, middleware, health endpoints, session orchestration |
| Networking | `github.com/coder/websocket` over `net/http` | Client input transport and authoritative state/event updates |
| Story content | JSON + validator tooling | Story definitions, branching events, encounters, metadata |
| Theme system | Asset manifests + per-theme packs | Tilesets, ambience, combat backdrops, UI skin references |
| Quality baseline | `go test`, `go vet`, `staticcheck`, `govulncheck` | Reliability and security hygiene |
| Release automation | GoReleaser | Tagged Linux/Windows builds, packaging, checksums |
| Distribution | Steam + Steamworks SDK | Lobbies, invites, achievements, native distribution |

## System Shape

```text
[Story JSON + Theme/Character Data]
               |
               v
        [Authoritative Go Server]
   (session state, checks, combat, sync)
               |
        WebSocket event/state flow
               |
               v
         [Go + raylib-go Client]
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
cmd/
  client/
  server/

internal/
  app/
  scene/
  render/
  input/
  netcode/
  game/
  story/
  theme/
  combat/
  character/
  item/
  dice/
  logx/

assets/
  themes/
  ui/
  sprites/
  audio/

stories/
  builtin/
  community/

test/
  integration/
  graphics/
  testutil/
```

## Engineering Practices

- Standard library first, minimal dependencies.
- Clear package boundaries and constructor injection.
- No ignored errors, no hidden panic paths outside startup.
- Structured logging with `slog`.
- Table-driven tests for parser, rules, and validation logic.
- Small end-to-end slices before broad feature expansion.

## Where Locks Live

This document explains the architecture, but exact locked MVP decisions such as protocol limits, release targets, save-path rules, and schema versioning live only in `docs/mvp_contract.md`.
