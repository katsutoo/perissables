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
- Create/join/rejoin carry all five aggregate content identity fields. Match
  them exactly before seat allocation, connection takeover, or gameplay state
  disclosure. `content_mismatch` preserves the existing reservation/connection;
  it cannot grant authority or trigger a content download. Identity fields and
  rejection fixtures start in Phase 02; canonical content fixtures follow in
  Phase 04.
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

Production create/join/rejoin validates a fresh Steam proof for the expected app
and ownership. The server issues all session/player IDs.

- Joining an existing session also requires the contract's server-held join
  grant for the exact identity/session, authorized by the current connected
  lobby owner. Steam metadata and an applicant's membership claim cannot
  substitute for that grant.
  Allocation consumes it atomically; rejected admission leaves it unconsumed.
  Phase 02 proves this with local identities; Phase 07 connects owner approvals
  to the Steam invite/membership flow before game admission.
- Raw Steam tickets are never logged or persisted.
- A returning Steam identity may reclaim only its own reserved in-memory seat,
  without needing a new join grant.
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
Disconnect preserves a seat through bounded grace. Explicit leave in any
non-ended state or grace expiry releases it. The contract's departure table
defines character/inventory removal, leader vacancy, last-player continuation,
and wipe/end behavior. Rejoin at or after the grace deadline cannot reclaim the
seat. A mismatched rejoin cannot evict an already-connected player.

The session owner applies due seat expiries before vote closure and new input.
Disconnect discards ballots; acknowledged rejoin permits fresh ballots only
while the vote is open. Gameplay pauses when no living player is connected;
absolute expiry and drain deadlines continue. Lifecycle changes and atomic story
check resolution share the same serialized authority path.

Session state exists only in the owning server process. Drain refuses new
admission, grants, and run starts. Idle lobbies end immediately; active runs
finish through summary and then end instead of returning to that process's lobby.
Reserved-seat rejoin remains available for non-ended running/summary sessions.
Summary and drain deadlines cannot be extended by reconnects, and cleanup never
waits indefinitely for peers to disconnect.

Idle/completed sessions receive the contract's maintenance notice. A drain
deadline, process crash, or forced stop that interrupts a run uses the stable
run-lost outcome. Regional routing must keep permitted rejoin traffic on the
owning process while it lives; maintenance does not migrate session state.

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
