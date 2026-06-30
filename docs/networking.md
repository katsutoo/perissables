# Networking

## Multiplayer Model

The game uses an authoritative server. Clients never decide movement success, combat outcomes, dice results, or story progression. They send intent; the server validates it, advances shared state, and publishes the result.

That model fits the project well because it keeps co-op sync readable, limits cheating, and lets story logic stay in one place.

## High-Level Flow

1. Player joins or rejoins a session.
2. Client sends ready/input messages.
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

For implementation, keep client->server and server->client message directions explicit. `join_response` is the server message that issues `session_id`, `player_id`, and the opaque reconnect credential; `rejoin` must present that credential and receive `resync_state` before new gameplay input is accepted.

## Production Session Discovery

Steam lobbies and invites are discovery UX, not gameplay authority. A lobby points players at the hosted authoritative server session by carrying metadata such as the WebSocket URL, `session_id`, protocol/schema versions, and pack identity/checksum. The server still validates Steam auth/session tickets, issues player credentials, and owns all state transitions.

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
- Use sequence numbers to drop stale or replayed input.
- Use heartbeat and timeout rules to detect dead connections.
- Keep gameplay-critical logic off the client.

## Reconnect And Sync

Reconnect matters because short co-op runs feel bad if one temporary disconnect destroys the session.

The plan is:

- issue session/rejoin credentials during join flow
- keep reconnect tokens opaque, server-generated, and out of Steam lobby metadata
- require a reconnect handshake before accepting new input
- send a full authoritative resync before returning a player to active play
- reject out-of-order or replayed input after reconnect

## Audio Rule

The server does not stream audio. It emits gameplay events, and each client plays matching local sounds. That keeps bandwidth lower and avoids making voice/audio transport part of the gameplay protocol.
