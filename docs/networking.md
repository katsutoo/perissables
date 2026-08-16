# Networking

Authority: protocol, identity, lifecycle, and safety requirements live in
`docs/mvp-contract.md`. Exact v1 DTOs and byte fixtures are frozen with the
Phase 02 implementation.

## Model

The server is authoritative from the first playable slice. Clients send intent;
the server validates identity and current state, advances headless rules, and
publishes recipient-specific views, events, and results.

The server owns:

- lobby/session lifecycle and ownership;
- positions, collision, interactions, and story progression;
- dice, combat, inventory, death, and run completion;
- input ordering, revisions, event IDs, deadlines, and reconnect; and
- every decision that can affect another player.

The client owns rendering, input collection, local UI/presentation state, and
audio playback after an authoritative event.

## Early Vertical Slice

Phase 02 proves the complete path before broad gameplay:

```text
client intent
    -> WebSocket DTO
    -> session owner
    -> headless game_core
    -> revision + event/result
    -> recipient view
    -> both clients converge
```

The slice covers lobby create/join, one interaction, one dice check, one legal
combat action, one rejected action, and summary. Later behavior extends this
path; it is not migrated from client authority.

## Transport

WebSockets fit the small co-op update model and the `axum`/`tokio` server.
Production uses WSS through trusted environment endpoints. Steam lobby metadata
contains discovery/compatibility data, never arbitrary endpoints or authority.

JSON v1 prioritizes debuggability. Before admission, authentication,
create/join, and rejoin use the pre-session envelope defined by the contract.
After identity/seat binding, session messages add server-issued session/player
IDs. Every envelope carries a type, protocol version, per-direction transport
sequence, and typed payload. Gameplay inputs also carry persisted per-player
order and a based-on revision.

Protocol rules:

- DTO direction and unknown-field behavior are explicit.
- Production supports one gameplay protocol version. Unsupported clients receive
  `update_required` before admission; deployments drain admitted sessions before
  removing that server version.
- IDs are opaque and server-generated.
- Every admitted expected input receives exactly one result.
- Recipient projections expose no secret or other-player private state.
- State is coalescible; semantic events/results are not silently discarded.
- Malformed, unauthorized, stale, duplicate, gap, rate, queue, and service
  failures have stable bounded behavior.
- Every message, ledger, mailbox, writer queue, timeout, retry, and fan-out is
  bounded.

Internal queue capacities are implementation budgets justified by tests and
measurements. They are not copied into protocol documentation unless a client
must know them.

## Identity

Production join/rejoin validates a fresh Steam proof for the expected app and
ownership. The server issues all session/player IDs.

- Raw Steam tickets are never logged or persisted.
- A returning Steam identity may reclaim only its own reserved in-memory seat.
- One player has at most one authoritative connection.
- The client acknowledges recipient-specific resync before new gameplay input.
- Stale takeover-close notifications cannot disconnect the replacement.
- The client stores no rejoin bearer token locally.

A local identity adapter supports deterministic development and tests before
Steam staging is available. It cannot compile into release features/packages.

## Lifecycle And Recovery

Steam lobbies are private or friends-only and invite-based; there is no public
browser or matchmaking. New seats join only a lobby; reserved seats may rejoin a
non-ended session.
Disconnect preserves a seat through bounded grace. Explicit lobby leave or
expiry releases it.

Session state exists only in the owning server process. A graceful deploy marks
that process unready, refuses new admission, and preserves existing connections
and reserved-seat rejoin while sessions drain. A process crash or forced stop
ends remaining runs; clients receive the stable run-lost outcome rather than a
partial restore. Regional routing must keep rejoin traffic on the owning process
while it lives.

## Slow And Failed Peers

State updates are latest-value coalesced. Required semantic events/results use
bounded ordered bundles. A client that cannot drain its bounded writer budget is
disconnected without allowing unbounded memory growth or stalling healthy
sessions.

Heartbeats detect dead connections. Rate limits and pre-auth concurrency bounds
protect expensive identity validation. Trusted proxy configuration, not
user-supplied forwarding headers, determines source attribution.

## Audio

The server never streams game audio. It emits semantic cue events, and each
client plays its matching local asset. This keeps audio out of transport
authority and bandwidth budgets.
