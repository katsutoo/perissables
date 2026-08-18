# Progress Tracker

This file tracks phase gates, not day-to-day tasks. Active implementation work,
owners, and small subtasks belong in issues or a project board. A phase checkbox
links its relevant code, CI, test, QA, security, benchmark, or decision evidence.

`docs/mvp-contract.md` owns product and compatibility requirements.
`docs/roadmap.md` owns execution order.

## Phase 00 - Scope Freeze

- [x] Product scope, non-goals, authority model, and milestones are clear.
- [x] Documentation authority and repository boundaries are clear.
- [x] Evidence-gated decisions are distinguished from locked contracts.
- [x] Phase 00 complete.

## Phase 01 - Repository Bootstrap

- [ ] Create the five-member Rust workspace with resolver 3, add the root
  `rust-toolchain.toml`, and commit `Cargo.lock`.
- [ ] Add minimal client/server entrypoints, shared logging, local `mise`
  pins for tools outside the Rust toolchain, local tasks, and the locked CI
  checks.
- [ ] Pin reproducible Linux/Windows native dependency acquisition, select the
  exact Steam Linux Runtime/container, and compile both minimal targets.
- [ ] Obtain Steamworks partner/app access, which is not currently available,
  then compile disposable authentication plus lobby/invite API spikes on both
  targets.
- [ ] Check Railway availability and entry-level subscription limits for an
  initial Americas/Europe/Asia candidate topology.
- [ ] Keep the versioned built-in content schema and fixtures in this
  repository.
- [ ] Phase 01 complete.

## Phase 02 - Early Authoritative Vertical Slice

- [ ] Two local test identities create/join one lobby through the real protocol.
- [ ] One map interaction, dice check, legal combat action, and rejected action
  resolve on the server.
- [ ] Both clients converge on the same summary revisions/events.
- [ ] Release feature checks prove the local identity adapter is absent.
- [ ] Run a headless/internal rules playtest and record unclear check, combat,
  and turn feedback for Phase 05.
- [ ] Phase 02 complete.

## Phase 03 - World Runtime

- [ ] Implement fixed client update/render loops and server-owned cardinal
  movement.
- [ ] Load one TMX map with collision, bounds, camera, and interaction behavior.
- [ ] Extend the authoritative transcript with movement and interaction
  boundary/error coverage.
- [ ] Phase 03 complete.

## Phase 04 - Story Schema And Runtime

- [ ] Implement built-in schema v1 and its strict repository-owned
  validator/conformance corpus.
- [ ] Implement declarative story nodes, choices, checks, effects, encounters,
  and bounded transitions.
- [ ] Freeze parser/resource ceilings from typical, large, limit, and rejected
  fixtures.
- [ ] Phase 04 complete.

## Phase 05 - Combat, Characters, Death, And Loot

- [ ] Complete roster, dice, combat, spell/item, death, loot, victory, and wipe
  behavior through authoritative state machines.
- [ ] Add deterministic public transcripts for legal, rejected, boundary, and
  rollback behavior.
- [ ] Remove every temporary debug character/action from release features.
- [ ] Run a blind combat playtest outside the implementation team and record
  turn clarity, idle time, encounter length, rules questions, and desire to
  replay.
- [ ] Phase 05 complete.

## Phase 06 - Complete Run Loop And Presentation

- [ ] Complete story selection, ready/start, summary/reset, and three repeated
  runs without stale state.
- [ ] Ship the supermarket presentation and storage-room area, one scalable
  keyboard-operable UI, settings, and fallback behavior.
- [ ] Complete event-driven ambience, music, and SFX without voice.
- [ ] Run the first complete story as a blind playtest and resolve blockers in
  comprehension, party downtime, run length, or willingness to replay.
- [ ] Phase 06 complete.

## Phase 07 - Multiplayer And Identity Hardening

- [ ] Add authorized staging Steam validation plus end-to-end lobby creation,
  join, and invites; prove expected app/ownership binding.
- [ ] Implement fresh-ticket reserved-seat reclaim, takeover, acknowledged
  resync, heartbeat, replay rejection, and stable public errors.
- [ ] Qualify bounded rate, mailbox/writer, slow-client, and transport-chaos
  behavior.
- [ ] Record the small authorized Steam test-account pool needed for identity
  qualification, separate from synthetic isolated capacity identities.
- [ ] Phase 07 complete.

## Phase 08 - Built-in content production

- [ ] Complete one `35-45` minute story, one supermarket map with storage-room
  area, four characters, five normal enemy types, and one mandatory boss.
- [ ] Complete eight character spells, eight items, two normal encounters plus
  the boss, four to six checks, three major choices, and `25-40` presented
  dialogue/choice beats.
- [ ] Complete four music tracks, two ambience loops, and at least twenty SFX.
- [ ] Record written rights/provenance for assets produced by the owner and
  credited friends plus any purchased sources; complete fluent
  collaborator-reviewed text for all nine launch locales.
- [ ] Phase 08 complete.

## Phase 09 - Playable MVP Stabilization

- [ ] Deterministic 2-, 3-, and 4-client full runs converge.
- [ ] Repeated run, reconnect, capacity, slow-client, settings, supermarket
  presentation, and single-UI regression scenarios pass.
- [ ] No blocker/critical defect or unexplained intermittent result remains.
- [ ] Phase 09 complete: Playable MVP.

## Phase 10 - Regional deployment and drain decision

- [ ] Freeze representative regional latency, full-run, reconnect, drain,
  forced-stop, and rollback workloads from the implemented game.
- [ ] Measure candidate Railway regions, automatic worst-latency
  player-to-session placement, and owning-process routing on production-shaped
  infrastructure; freeze acceptable cross-area latency with controlled
  impairment playtests.
- [ ] Record latency by source/host region, drain duration, sessions
  completed/lost, deployment time, CPU/RSS/network, and operator steps.
- [ ] Freeze the smallest Railway-only region set within the entry-level budget,
  automatic worst-latency placement/rejoin routing, drain deadline, run-lost
  behavior, and rollback procedure.
- [ ] Phase 10 complete.

## Phase 11 - Regional service and graceful drain

- [ ] Implement the selected regional placement and owning-process routing.
- [ ] Stop admission before drain while preserving live sessions and
  reserved-seat rejoin through the measured deadline.
- [ ] Pass production-shaped regional loss, process replacement, forced-stop,
  crash/run-loss, observability, and rollback tests.
- [ ] Phase 11 complete.

## Phase 12 - Release QA, Balance, And Performance

- [ ] Freeze the reference client machine, staging shape, benchmark workloads,
  timer/driver noise floor, and practical release gates.
- [ ] Complete consented playtests and report onboarding/run timing with
  abandonment and uncertainty.
- [ ] Run the pre-release QA matrix and calibrated client/server benchmarks on
  production-profile candidate digests.
- [ ] Qualify Steam authentication with the authorized account pool and run
  capacity tests with the isolated non-release synthetic identity adapter.
- [ ] Prepare and assign owners for Steam store media/copy, age/content
  disclosures, launch languages, pricing, privacy/support contacts, incident
  handling, server operating cost, and shutdown policy.
- [ ] Phase 12 complete.

## Phase 13 - Steam Packaging And Production

- [ ] Build reproducible Linux/Windows packages and sign Windows artifacts
  through the protected signing workflow.
- [ ] Verify the Phase 07 Steam lobby/invite flow against final depots; complete
  production deployment, observability, and the rollback drill.
- [ ] Rerun final QA and performance gates against exact final digests.
- [ ] Approve and publish store materials, age/content disclosures, pricing,
  launch languages, end-user terms, privacy/support contacts, third-party
  notices, contribution-policy decision, asset provenance, and the operating
  plan.
- [ ] Phase 13 complete: Release-ready.

## Post-MVP

Track additional themes/UI, voice, durable active runs, normal solo play,
achievements, and additional platforms in separate post-MVP backlogs. They are
not Phase 01-13 completion gates.
