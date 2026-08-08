# Networking

Authority: exact message, lifecycle, sequence, rate, queue, timeout, and identity behavior derives from "Locked Protocol And Limits" in `docs/mvp-contract.md`. This document is explanatory.

## Multiplayer Model

The game uses an authoritative server. Clients never decide movement success, combat outcomes, dice results, or story progression. They send intent; the server validates it, advances shared state, and publishes the result.

That model fits the project well because it keeps co-op sync readable, limits cheating, and lets story logic stay in one place.

## High-Level Flow

1. Player creates/joins a lobby, rejoins a reserved seat, or explicitly leaves a lobby.
2. Client sends sequenced `input` actions for lobby readiness and gameplay.
3. Server validates the request against current state.
4. Server updates the canonical session state.
5. Server emits snapshots and events to all clients.
6. Clients render the new world/combat/dialog state and play local feedback.

## Why WebSockets

- Good fit for a small real-time co-op game with frequent but lightweight state updates.
- Simple enough for MVP implementation and debugging.
- Works cleanly with an `axum`/`tokio` server and reverse-proxy TLS termination.

## Message Shape

All gameplay traffic follows one shared envelope shape, with message-specific payloads inside it. Exact field requirements, allowed message types, sequence rules, and size/rate limits are locked in `docs/mvp-contract.md`.

For design purposes, the important rule is:

- envelope stays stable
- payload varies by message type
- unknown or stale messages are rejected

For implementation, keep directions explicit. `join` carries create/join intent, the target lobby when joining, Steam proof, and the aggregate compatibility tuple. `join_response` issues server-generated player/session IDs plus the initial rejoin token. `rejoin` presents fresh Steam proof plus a current or pending-handoff token; the server durably returns the candidate token in `rejoin_response`, sends `resync_state`, and promotes rotation only after `resync_ack`. Authenticated `leave` is a terminal lobby input that releases the seat after commit. Exact payloads remain in "Locked Protocol And Limits" in `docs/mvp-contract.md`.

## Production Session Discovery

Steam lobbies and invites are discovery UX, not endpoint or gameplay authority. Lobby metadata carries the target `session_id`, independent protocol/content/rules versions, and aggregate pack identity/checksum. New seats may join only the authoritative `Lobby`; reserved seats use rejoin in any non-ended phase. The WebSocket origin comes only from the trusted environment allowlist, so manipulated lobby metadata cannot redirect Steam proof. The server validates tickets, computes its own compatibility values, issues identifiers/tokens, and owns all state transitions.

## Authority Boundaries

The server owns:

- session and lobby lifecycle
- map positions and collision validity
- event triggers and story progression
- dice rolls and check resolution
- combat turns, action validity, and outcomes
- reconnect and state resync decisions

The client owns:

- rendering
- local input collection
- local audio playback after receiving approved gameplay events
- temporary presentation state that does not affect gameplay authority

## Security Posture For MVP

- Validate every inbound message shape before acting on it.
- Reject unknown message types and impossible state transitions.
- Enforce frame size limits and per-client rate limits.
- Treat control messages separately from gameplay input, enforce pre-auth handshake/global bounds, and throttle failed `join`/`rejoin` attempts before expensive validation.
- Use per-connection transport sequences for ordering and persisted per-player input sequences for replay protection across reconnects.
- Use heartbeat and timeout rules to detect dead connections.
- Bound the unresolved-input ledger, reserved session-mailbox partitions, and one unified byte-accounted writer queue; coalesce only replaceable state/control and disconnect slow consumers instead of growing memory.
- Keep gameplay-critical logic off the client.

## Reconnect And Sync

Reconnect matters because short co-op runs feel bad if one temporary disconnect destroys the session.

The plan is:

- issue a CSPRNG-backed rejoin token alongside server-generated session/player IDs during join
- keep reconnect tokens opaque, server-generated, and out of Steam lobby metadata
- commit and return a pending rotation token, then require resync acknowledgement before promoting it or accepting new input
- send a bounded player-specific authoritative resync and require acknowledgement before returning a player to active play
- reject out-of-order or replayed input after reconnect
- after Phase 16, commit versioned session state through the dedicated SQLite persistence owner so a Release-ready restart or deploy looks like an ordinary disconnect/reconnect; before Phase 16, run loss remains explicit pre-release behavior

## Audio Rule

The server does not stream audio. It emits gameplay events, and each client plays matching local sounds. That keeps bandwidth lower and avoids making voice/audio transport part of the gameplay protocol.
