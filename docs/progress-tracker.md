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

- [ ] Create the five-member Rust workspace with resolver 3, pinned
  `rust-toolchain.toml`, and committed `Cargo.lock`.
- [ ] Add minimal client/server entrypoints, shared logging, local `mise`
  tasks, and the locked CI checks.
- [ ] Pin reproducible Linux/Windows native dependency acquisition and compile
  both minimal targets.
- [ ] Create and record the MIT `les-perissables-stories` repository
  relationship.
- [ ] Phase 01 complete.

## Phase 02 - Early Authoritative Vertical Slice

- [ ] Two local test identities create/join one lobby through the real protocol.
- [ ] One map interaction, dice check, legal combat action, and rejected action
  resolve on the server.
- [ ] Both clients converge on the same summary revisions/events.
- [ ] Release feature checks prove the local identity adapter is absent.
- [ ] Phase 02 complete.

## Phase 03 - World Runtime

- [ ] Implement fixed client update/render loops and server-owned cardinal
  movement.
- [ ] Load one TMX map with collision, bounds, camera, and interaction behavior.
- [ ] Extend the authoritative transcript with movement and interaction
  boundary/error coverage.
- [ ] Phase 03 complete.

## Phase 04 - Story Schema And Runtime

- [ ] Publish schema v1 and its strict validator/conformance corpus from
  `les-perissables-stories`.
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
- [ ] Phase 05 complete.

## Phase 06 - Complete Run Loop And Presentation

- [ ] Complete story selection, ready/start, summary/reset, and three repeated
  runs without stale state.
- [ ] Ship three themes, at least two behavior-neutral UI variants, keyboard-only
  operation, settings, and fallback behavior.
- [ ] Complete event-driven ambience, music, SFX, and optional voice barks.
- [ ] Phase 06 complete.

## Phase 07 - Multiplayer And Identity Hardening

- [ ] Add authorized staging Steam validation and prove expected app/ownership
  binding.
- [ ] Implement rejoin rotation, takeover, acknowledged resync, heartbeat,
  replay rejection, and stable public errors.
- [ ] Qualify bounded rate, mailbox/writer, slow-client, and transport-chaos
  behavior.
- [ ] Record the authorized Steam test-account availability needed for final
  capacity testing.
- [ ] Phase 07 complete.

## Phase 08 - Creator Tooling And Tier 1

- [ ] Build `storycheck`, canonical checksum/path validation, creator docs, and
  reusable conformance fixtures.
- [ ] Implement safe local import and single-player creator testing for Tier 1
  stories/characters only.
- [ ] Qualify fail-closed Linux and Windows validation sandboxes before exposing
  creator import.
- [ ] Prove hosted multiplayer remains built-in-only and custom media/maps are
  rejected.
- [ ] Phase 08 complete.

## Phase 09 - Playable MVP Stabilization

- [ ] Deterministic 2-, 3-, and 4-client full runs converge.
- [ ] Repeated run, reconnect, capacity, slow-client, settings, theme, and UI
  regression scenarios pass.
- [ ] No blocker/critical defect or unexplained intermittent result remains.
- [ ] Phase 09 complete: Playable MVP.

## Phase 10 - Persistence Spike And Decision

- [ ] Freeze representative and maximum valid session-state fixtures from the
  implemented game.
- [ ] Measure storage alternatives and durability strategies on the actual
  production-shaped Railway environment.
- [ ] Record latency, throughput, bytes/fsyncs, CPU/RSS, recovery, corruption,
  backup, rollback, and operational trade-offs.
- [ ] Update the MVP contract with the selected store, schema,
  acknowledgement boundary, and measured budgets.
- [ ] Phase 10 complete.

## Phase 11 - Durable Sessions

- [ ] Implement the selected persistence adapter and versioned server-owned
  state.
- [ ] Pass integration fault tests for migration, storage failure, restart,
  crash, corruption, and rollback compatibility.
- [ ] Pass release-mode clean/crash recovery on the actual staging environment.
- [ ] Phase 11 complete.

## Phase 12 - Release QA, Balance, And Performance

- [ ] Freeze the reference client machine, staging shape, benchmark workloads,
  timer/driver noise floor, and practical release gates.
- [ ] Complete consented playtests and report onboarding/run timing with
  abandonment and uncertainty.
- [ ] Run the pre-release QA matrix and calibrated client/server benchmarks on
  production-profile candidate digests.
- [ ] Provision the authorized account and real source-IP topology required for
  the maximum hosted workload.
- [ ] Phase 12 complete.

## Phase 13 - Steam Packaging And Production

- [ ] Build reproducible Linux/Windows packages and sign Windows artifacts
  through the protected signing workflow.
- [ ] Complete Steam lobbies/invites, depot layout, production deployment,
  backup, observability, and rollback drill.
- [ ] Rerun final QA and performance gates against exact final digests.
- [ ] Complete end-user terms, third-party notices, contribution-policy
  decision, and asset provenance.
- [ ] Phase 13 complete: Release-ready.

## Post-MVP

Track the community hub, public UGC, moderation, publication signing, custom
maps/media, Tier 2 packs, achievements, and additional platforms in separate
backlogs and repositories. They are not Phase 01-13 completion gates.
