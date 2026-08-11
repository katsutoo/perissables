# Architecture

Authority: product and compatibility requirements live in
`docs/mvp-contract.md`. This document explains the intended structure and the
reasoning behind it. Storage mechanisms and measured capacity budgets remain
provisional until their named decision gates pass.

## Architecture Goals

- Exercise the real client/server boundary from the first playable slice.
- Keep all gameplay rules headless, deterministic, and server-authoritative.
- Keep stories, characters, and themes data-driven.
- Prefer small, explicit modules and bounded resources.
- Lock implementation details only after evidence makes the trade-off clear.

## Stack

| Layer | Technology | Responsibility |
| --- | --- | --- |
| Client | Rust + `raylib` (`raylib-rs`) | Window, rendering, input, UI, audio |
| Server | Rust + `axum` + `tower` + `tokio` | Sessions, authority, health, networking |
| Shared rules | Headless Rust crate | Story, dice, combat, character, and run state machines |
| Protocol | `serde` + `serde_json` over WebSockets | Versioned intents, views, events, and errors |
| Content | JSON + TMX + validator crate | Built-in and creator-authored declarative data |
| Persistence | Evidence-gated in Phase 10 | Durable sessions after a measured SQLite/Postgres decision |
| Hosting | Railway | Staging and production server environments |
| Distribution | Steam | Ownership, lobbies/invites, and Linux/Windows depots |

The current storage candidate is bundled SQLite on a Railway volume. It is not
an architecture lock. Phase 10 measures it against the actual state shape and
deployment environment before the project selects SQLite, managed PostgreSQL,
or a different mature transactional store.

## System Shape

```text
[Validated story/theme/character data]
                  |
                  v
        [Headless game_core]
                  |
                  v
       [Authoritative server]
          |             |
   versioned views   persistence adapter
          |
       WebSocket
          |
          v
            [Client]
     render, input, audio
```

The client sends intent and renders recipient-specific authoritative views. It
never decides movement success, story outcomes, dice, combat, inventory, or run
completion.

## Early Authoritative Vertical Slice

Phase 02 is the first executable product proof. It must use the real client,
server, protocol, and headless rules path for:

1. creating one lobby and joining it with two local test identities;
2. starting one tiny run;
3. moving to and activating one map interaction;
4. resolving one server-generated dice check;
5. resolving one legal and one rejected combat action; and
6. returning both clients to the same summary state.

The local identity adapter exists only for local tests and development builds.
Production features and packages must not compile it. Steam identity replaces
that adapter without changing gameplay authority.

Later phases deepen this slice instead of relocating client-owned gameplay to the
server. A feature is not part of the game until it travels through the same
intent -> authority -> view/event path.

## Crate Boundaries

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
      identity/
  game_core/
    src/
      game/
      story/
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

assets/
stories/
  builtin/
```

Workspace members are `client`, `server`, `game_core`, `shared`, and
`integration_tests`. Package names use the `les-perissables-` prefix.

Dependency direction is one-way:

- `shared` owns protocol DTOs, IDs, and common logging setup.
- `game_core` depends on `shared` and the pinned MIT pack crate.
- `client` and `server` depend on `shared` and `game_core`.
- `integration_tests` may depend on every workspace crate.

Reverse dependencies and cycles are forbidden.

## Headless game_core

`game_core` has no `raylib`, network, filesystem, environment, wall-clock, or
process-global access.

- State changes only through typed intents, ticks, deadlines, and explicit
  start-run entropy.
- Randomness and time are injected.
- Operations return new revisions plus stable events/results.
- Persistence DTOs and pure migrators may live here; storage I/O may not.
- Tests can replay the same transcript without a renderer or socket.

## Runtime Ownership

- One session owner serializes mutation for one session.
- Long-lived tasks belong to one supervised startup/task tree.
- Every task, timeout, queue, and blocking operation is bounded.
- Critical tasks are cancelled and joined explicitly; dropped handles may not
  silently detach authoritative work.
- Native `raylib` and Steamworks calls stay behind small safe adapters.
- `unsafe` is forbidden in domain crates. Adapter unsafe blocks require a
  local `SAFETY` explanation and tests for lifetime/thread-affinity contracts.

## Persistence Decision Gate

Phase 10 builds a production-shaped spike after representative state DTOs and
network behavior exist. The decision compares at least:

- bundled SQLite on the actual Railway volume;
- managed PostgreSQL when SQLite cannot meet the gates;
- full snapshots versus bounded deltas/checkpoints;
- durability at every acknowledged action versus semantic-boundary durability;
- commit cadence, write amplification, recovery time, and operational burden.

The spike covers typical and maximum valid states, clean restart, crash recovery,
checkpoint behavior, concurrent sessions, and storage failure. It records
latency distributions, throughput, bytes written, fsyncs, CPU, peak RSS, and
recovery correctness.

Only the selected design is then added to `docs/mvp-contract.md`, including its
schema, acknowledgement boundary, queue/cadence budgets, and rollback rules.
Until that decision, no document may claim that a particular queue size,
snapshot cap, pragma set, or commit interval is final.

## Engineering Practices

- Rust `1.97.1`, Edition 2024, virtual-workspace resolver 3.
- Commit `Cargo.lock`; CI and release builds use `--locked`.
- Prefer the standard library and concrete types until an abstraction has more
  than one useful consumer or implementation.
- Use typed errors at reusable boundaries and contextual errors in executable
  orchestration. Recoverable input, I/O, and dependency failures do not panic.
- Use structured `tracing` fields without secrets or personal data.
- Keep transport, rules, presentation, and persistence separable at their real
  ownership boundaries; do not create layers solely to match a diagram.
- Add tests with the behavior that introduces them. Profile only measured hot
  paths.
- Follow `docs/test-strategy.md`, `docs/qa-plan.md`,
  `docs/security-model.md`, and `docs/benchmark-plan.md` for evidence.
