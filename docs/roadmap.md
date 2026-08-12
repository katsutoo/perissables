# Roadmap

## Planning Rule

This roadmap is an execution order, not a calendar promise. For one developer,
the previous `4-8 month` release range remains an aspirational hypothesis until
the Phase 02 authoritative slice, native target builds, Steam access, and Phase
10 storage decision have produced evidence.

Re-estimate after each of those gates. A failed feasibility gate changes the
plan before dependent implementation begins; it does not get hidden inside a
later phase.

## Milestones

| Milestone | Phases | Outcome |
| --- | --- | --- |
| Authoritative vertical slice | 00-02 | Two clients complete one tiny server-owned run |
| Playable MVP | 00-09 | Complete repeatable 2-4 player game loop with built-in content |
| Release-ready | 00-13 | Durable sessions, measured performance, Steam packages, and production operations |
| Post-MVP | Separate backlog/repositories | Community hub, public UGC, Tier 2 custom assets, and other extensions |

## Phase Plan

### Phase 00 - Scope Freeze

Define the player experience, product boundaries, authority model, and completion
evidence. Phase 00 is documentation-only.

Exit: the current contract and documentation map agree, with no implementation
detail labeled final before its evidence gate.

### Phase 01 - Repository Bootstrap

Create the five-member Rust workspace, add the root pinned toolchain, commit the
lockfile, and add minimal client/server entrypoints, logging, CI, legal baseline,
and local development tasks. Keep Rust owned by `rust-toolchain.toml`; use
`mise` only for tools outside the Rust toolchain and local task aliases. Create
the separate MIT `les-perissables-stories` repository and pin its relationship.

Run only bootstrap-relevant feasibility checks:

- Linux and Windows native dependencies can be acquired reproducibly.
- Both targets compile minimal client/server shells.
- CI can run the locked Rust checks.

Steam account pools, distributed load-generator IPs, production signing, pack
sandbox qualification, and storage selection are later gates owned by the phases
that need them.

### Phase 02 - Early Authoritative Vertical Slice

Build the first runnable product through the real server boundary:

- two local test identities create/join one lobby;
- the server starts one tiny run;
- one map interaction triggers one check;
- one legal and one rejected combat action resolve;
- both clients converge on the same summary.

Use real protocol DTOs, revisions, events, errors, and bounded queues from the
start. Restart durability, Steam authentication, polished rendering, and broad
content are not required yet.

### Phase 03 - World Runtime

Implement fixed-timestep client presentation, server-owned cardinal movement,
TMX loading/collision, camera behavior, interactions, and the minimal audio
event path. Extend the Phase 02 transcript rather than creating a separate
single-player rules path.

### Phase 04 - Story Schema And Runtime

Publish schema v1 from `les-perissables-stories`, implement strict validation
and a declarative story state machine, and add branching, choices, checks,
effects, encounters, and return/end transitions. Freeze exact parser/resource
ceilings with the validator corpus, not before it exists.

### Phase 05 - Combat, Characters, Death, And Loot

Complete the preset roster, dice/combat rules, items, death, corpse loot,
victory/wipe behavior, and deterministic authoritative transcripts. Remove every
temporary debug character and action before exit.

### Phase 06 - Complete Run Loop And Presentation

Finish story selection, ready/start, summary/reset, three built-in themes,
behavior-neutral UI variants, keyboard-only operation, settings, fallback
assets, music/SFX/voice channels, and three-run reset coverage.

### Phase 07 - Multiplayer And Identity Hardening

Add Steam ticket validation for isolated staging, identity-bound rejoin tokens,
single-connection takeover, acknowledged resync, replay protection, heartbeat,
rate limits, slow-client handling, and deterministic transport-chaos tests.

Local test identity remains available only to tests and non-release development
builds and is proven absent from production packages.

### Phase 08 - Creator Tooling And Tier 1 Packs

Build `storycheck`, canonical checksum/path validation, the creator tutorial,
safe local import, and single-player creator testing for reuse-only Tier 1 packs.
Hosted multiplayer remains built-in-only. Custom maps/media and public UGC are
post-MVP.

### Phase 09 - Playable MVP Stabilization

Run deterministic 2-, 3-, and 4-client scenarios, repeated run loops, reconnect
coverage, supported-resolution presentation checks, and risk-based regression
tests. Resolve blocker/critical defects before declaring Playable MVP.

### Phase 10 - Persistence Spike And Decision

Measure the actual state and write workload on production-shaped Railway
infrastructure. Compare bundled SQLite with managed PostgreSQL when appropriate,
and compare snapshot/checkpoint strategies and durability boundaries.

Publish one decision record containing:

- selected store and operational topology;
- durable acknowledgement semantics;
- schema/version/migration policy;
- measured commit cadence and queue/storage budgets;
- clean/crash recovery behavior;
- backup and rollback approach; and
- rejected alternatives with measured reasons.

Only this decision gates durable-session implementation. It does not block
Phases 02-09.

### Phase 11 - Durable Sessions

Implement the Phase 10 decision behind the server persistence adapter. Persist
only server-owned state, never raw bearer tokens or client-submitted saves.
Exercise migration, storage failure, clean shutdown, crash recovery, corruption,
and rollback compatibility with real integration tests.

### Phase 12 - Release QA, Balance, And Performance

Provision production-equivalent staging. Run consented playtests, the release
QA matrix, calibrated client measurements, server capacity/load experiments,
recovery workloads, and profiling of measured bottlenecks. Freeze release
budgets only after the reference environment and noise floor are recorded.

### Phase 13 - Steam Packaging And Production

Create reproducible signed Linux/Windows packages, complete Steam
lobbies/invites and depot layout, deploy the production server, drill rollback,
and rerun release QA plus performance gates against the exact final digests.
Complete end-user terms, third-party notices, and asset provenance.

## Post-MVP Boundary

The community hub, accounts, likes/comments, public uploads, moderation,
publication signing, custom maps/media, and Tier 2 content belong to a separate
post-MVP plan and repository. This game repository retains only the game-facing
pack compatibility interface once that work begins.

## Definition Of Done

For every implementation phase:

- The acceptance behavior and failure oracle are clear before implementation.
- New behavior has focused public-interface tests for meaningful valid,
  boundary, and error paths.
- CI passes on a clean checkout for supported features and targets.
- Security-sensitive changes update the threat model and negative coverage.
- Performance-sensitive choices use the smallest production-mode measurement
  that can support the decision.
- Manual QA uses the shipped entrypoint only when a user journey exists.
- No blocker/critical defect or unexplained intermittent failure remains.
- The tracker links evidence; active implementation detail lives in issues.
