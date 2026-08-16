# Roadmap

## Planning Rule

This roadmap is an execution order, not a calendar promise. For one developer,
the previous `4-8 month` release range remains an aspirational hypothesis until
the Phase 02 authoritative slice, native target builds, Steam access, and Phase
10 regional deployment decision have produced evidence.

Re-estimate after each of those gates. A failed feasibility gate changes the
plan before dependent implementation begins; it does not get hidden inside a
later phase.

## Milestones

| Milestone | Phases | Outcome |
| --- | --- | --- |
| Authoritative vertical slice | 00-02 | Two clients complete one tiny server-owned run |
| Playable MVP | 00-09 | Complete repeatable 2-4 player game loop with built-in content |
| Release-ready | 00-13 | Regional service, measured performance, Steam packages, and production operations |
| Post-MVP | Separate backlog/repositories | Creator tooling, community content, additional themes/UI, durable runs, and other extensions |

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
`mise` only for tools outside the Rust toolchain and local task aliases. Keep the
built-in content schema and fixtures in this repository.

Run only bootstrap-relevant feasibility checks:

- Linux and Windows native dependencies can be acquired reproducibly.
- The exact Steam Linux Runtime/container is selected and pinned.
- Both targets compile minimal client/server shells.
- Steamworks partner/app access is obtained, and minimal authentication plus
  lobby/invite APIs compile on both targets in a disposable spike.
- Railway availability and entry-level subscription limits are checked for an
  initial Americas/Europe/Asia candidate topology; nearby worldwide latency is
  measured later rather than requiring a physical deployment in every area.
- CI can run the locked Rust checks.

Large Steam account pools, production integration, production signing, and
regional deployment are later gates owned by the phases that need them.

### Phase 02 - Early Authoritative Vertical Slice

Build the first runnable product through the real server boundary:

- two local test identities create/join one lobby;
- the server starts one tiny run;
- one map interaction triggers one check;
- one legal and one rejected combat action resolve;
- both clients converge on the same summary.

Use real protocol DTOs, revisions, events, errors, and bounded queues from the
start. Restart durability, Steam authentication, polished rendering, and broad
content are not required yet. Before broad rendering, run the slice as a
headless/internal rules playtest and record unclear check, combat, and turn
feedback for Phase 05.

### Phase 03 - World Runtime

Implement fixed-timestep client presentation, server-owned cardinal movement,
TMX loading/collision, camera behavior, interactions, and the minimal audio
event path. Extend the Phase 02 transcript rather than creating a separate
single-player rules path.

### Phase 04 - Story Schema And Runtime

Implement the internal built-in schema v1, strict validation, and a declarative
story state machine. Add branching, choices, checks, effects, encounters, and
return/end transitions. Freeze exact parser/resource ceilings with repository
conformance fixtures, not before they exist.

### Phase 05 - Combat, Characters, Death, And Loot

Complete the preset roster, dice/combat rules, items, death, corpse loot,
victory/wipe behavior, and deterministic authoritative transcripts. Remove every
temporary debug character and action before exit. Run a small blind combat
playtest with people outside implementation and record turn clarity, idle time,
encounter length, rules questions, and desire to replay before locking the
combat loop.

### Phase 06 - Complete Run Loop And Presentation

Finish story selection, ready/start, summary/reset, the supermarket presentation
and storage-room area, one scalable keyboard-operable UI, settings, fallback
assets, music/SFX channels, and three-run reset coverage. Run the first complete
story as a blind playtest and resolve any blocker in comprehension, party
downtime, run length, or willingness to replay before exit.

### Phase 07 - Multiplayer And Identity Hardening

Add Steam ticket validation and end-to-end Steam lobby creation, join, and
invite handling in isolated staging. Add fresh-ticket reserved-seat reclaim,
single-connection takeover, acknowledged resync, replay protection, heartbeat,
rate limits, slow-client handling, and deterministic transport-chaos tests.

Local and benchmark identity adapters remain available only to tests and
non-release builds and are proven absent from production packages.

### Phase 08 - Built-in content production

Complete the one-story release minimum: one supermarket map with storage-room
area, four playable characters, five normal enemy types, one mandatory boss,
eight spells, eight items, the required encounters/checks/choices/dialogue,
four music tracks, two ambience loops, and at least twenty SFX. Integrate the
assets produced by the owner and credited friends, plus any purchased sources,
with written rights/provenance. Produce fluent collaborator-reviewed text for
all nine launch locales.

### Phase 09 - Playable MVP Stabilization

Run deterministic 2-, 3-, and 4-client scenarios, repeated run loops, reconnect
coverage, supported-resolution presentation checks, and risk-based regression
tests. Resolve blocker/critical defects before declaring Playable MVP.

### Phase 10 - Regional deployment and drain decision

Measure candidate Railway regions, automatic worst-latency player-to-session
placement, owning-process rejoin routing, graceful drain, forced stop, regional
loss, and rollback using the implemented game workload. Freeze acceptable
cross-area latency from controlled impairment playtests before choosing the
smallest region set.

Publish one decision record containing:

- the smallest Railway-only region set within the entry-level budget and its
  automatic worst-latency placement policy;
- cross-region latency and uncertainty;
- owning-process session/rejoin routing;
- measured drain deadline and admission behavior;
- forced-stop and run-lost behavior;
- deployment and rollback procedure; and
- rejected alternatives with measured reasons.

### Phase 11 - Regional service and graceful drain

Implement the Phase 10 decision. A draining process becomes unready, stops new
admission, preserves live sessions and reserved-seat rejoin until the measured
deadline, then exits. Exercise regional routing, process replacement, forced
stop, crash/run-loss handling, observability, and rollback with production-shaped
integration tests.

### Phase 12 - Release QA, Balance, And Performance

Provision production-equivalent staging. Run consented playtests, the release
QA matrix, calibrated client measurements, server capacity/load experiments,
regional drain/boundary workloads, and profiling of measured bottlenecks. Freeze release
budgets only after the reference environment and noise floor are recorded.

Prepare and assign owners for Steam store copy and media, age/content
disclosures, launch languages, pricing, privacy/support contacts, incident
handling, server operating cost, and shutdown policy.

### Phase 13 - Steam Packaging And Production

Create reproducible signed Linux/Windows packages, verify the Phase 07 Steam
lobby/invite flow against final depots, deploy the production server, drill
rollback, and rerun release QA plus performance gates against the exact final
digests. Approve and publish the store materials, disclosures, end-user terms,
privacy/support contacts, third-party notices, asset provenance, and operating
plan.

## Post-MVP Boundary

Creator tooling, `storycheck`, local/archive import, a separate MIT schema
repository, community hosting, accounts, likes/comments, public uploads,
moderation, publication signing, custom maps/media, additional themes/UI, voice,
durable active runs, and normal solo play belong to post-MVP plans. No
compatibility interface is reserved until one of those projects begins.

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
