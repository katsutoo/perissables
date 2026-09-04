# Test Strategy

Status: Normative automated-test policy
Owner: Sole developer
Updated: 2026-09-04

Product and protocol behavior comes from `docs/mvp-contract.md`. This document
defines how automated tests earn confidence without becoming a second
implementation or a paperwork system.

## Principles

- Test observable behavior through public interfaces.
- One test protects one behavior; table cases may cover several inputs of that
  same behavior.
- Assert concrete outputs and negative space. A rejected action must not mutate
  state, consume a turn, spend resources, or emit gameplay events.
- Inject time, RNG state, IDs, transport schedules, process-stop signals, and
  identity-provider responses.
- No sleeps for synchronization, wall-clock dependence, shared mutable fixtures,
  order dependence, or retry-until-green.
- Use real deterministic collaborators when cheap. Substitute only the boundary
  that is unsafe, slow, nondeterministic, or external.

## Fail-First Evidence

Every new test is observed failing for the intended reason before acceptance.

- Regression fixes preserve the failing test name, command, and failure
  signature in the issue or change record.
- Critical protocol, session-drain, security, and content invariants preserve
  equivalent red/green evidence.
- Routine test-first feature work does not require a permanent test-only commit,
  tree identifier, or standalone evidence bundle.
- Pure refactors state why no new behavior test is needed.

## Layers

| Layer | Protects |
| --- | --- |
| Unit/property | Dice, combat, story rules, bounds, and canonicalization |
| Schema/conformance | Accepted/rejected built-in JSON, TMX, paths, and bytes |
| State-machine transcript | Authoritative revisions, events, rejection, reset |
| Protocol contract | DTOs, directions, versions, sequences, projections, errors |
| Multi-client integration | 2/3/4-client convergence, reconnect, replay, capacity |
| Process lifecycle | Admission stop, drain, deadline, forced stop, and run loss |
| Release smoke | Shipped startup and critical user journeys on supported OSes |

Release smoke is automated end-to-end coverage only when it launches the shipped
entrypoint. Internal server/component tests remain integration tests.

## Selecting Cases

Use boundary analysis, not mechanical case multiplication.

- Cover empty/zero, one, a representative value, the meaningful limit, just
  beyond the limit, malformed input, and the credible operational error where
  those cases can change behavior.
- Do not generate limit-minus-one/limit/limit-plus-one cases for every field when
  a shared validator or property states the same invariant once.
- Use exhaustive cases where the state space is deliberately small, such as all
  d100 results.
- Use property tests for broad invariants such as serialization round trips,
  projection secrecy, queue accounting, canonical checksums, and no mutation on
  rejection.
- Keep fixtures minimal. Production-sized fixtures belong only where size or
  concurrency is the behavior under test.

## Required High-Risk Coverage

- Authoritative transcripts cover one legal and one illegal action at every
  gameplay phase, stable ordering, leader/tied votes, corpse-loot votes,
  spell-only enemy kits, each spell result class and critical redirection
  polarity, death/wipe, and three-run reset.
- Protocol tests cover fragmentation, malformed DTOs, wrong direction/version,
  stale/duplicate/gap sequences, recipient-specific projections, queue refusal,
  reconnect handoff, and slow writers.
- Identity tests cover wrong app/identity, stale Steam proof, takeover, seat
  reclaim, expiry, and absence of local/benchmark adapters in release features.
- Content tests cover duplicate/unknown fields, resource limits, disabled TMX
  external access, graph errors, reference errors, checksum stability, and the
  exact release-content minimum.
- Process-lifecycle tests cover admission refusal after unready, reserved-seat
  rejoin while draining, natural completion, deadline expiry, forced stop, and
  stable run-loss behavior.

The following focused transcripts close the documentation review findings. Add
them with their owning phase; these are acceptance cases, not tests already run.

| Phase | Behavior and required oracle |
| --- | --- |
| 02, extended in 04 | Exact five-field content match admits; change each field independently to get `content_mismatch` without seat allocation or gameplay state. Missing/malformed fields and wrong protocol retain their distinct errors. Phase 04 adds canonical checksum and locale-selection cases. |
| 04 | Zero story ballots select the declared default; one winner, leader tie, and lexical tie each select the specified choice without RNG consumption. Disconnect discards an existing ballot; rejoin cannot resurrect it. A fresh ballot before closure counts; at the deadline it is rejected. |
| 04 | Check actors follow eligible leader then lexical fallback even when a different player initiated the interaction. Use unequal stats and a fixed roll that distinguishes them. An empty candidate set leaves the node waiting without RNG/modifier consumption; reconnect resolves once; ending cancels. Disconnect before versus after atomic resolution cannot produce a second roll. |
| 05 | Repeated normal/critical guards replace remaining hit counts; zero damage preserves guard. Repeated bonuses/penalties replace their own slots, combine once, clamp, and both expire on every actual check including criticals. Death/encounter end clears effects; disconnect does not. |
| 05 | Items resolve without RNG, spend exactly one item/turn on legal use, preserve next-check modifiers, and spend nothing on rejection. Full-HP use still spends; full-inventory story grants skip without replacing items or blocking progression. Multiple grants use authored order. |
| 05 | No loot ballots or no valid recipient means no assignment. Revalidate recipients at closure; simultaneous awards competing for one free slot resolve in the specified corpse/slot order without overflow or duplicate grants. |
| 06, transport qualification in 07 | Leader expires while all others are disconnected: vacant role, no random draw, then one election on eligible rejoin. Empty lobby ownership similarly recovers. Leave/expiry in each state releases the right seat and removes active characters/items without loot. One-player continuation, no-living-player wipe, no-seat end, and reset to a two-player start requirement are distinct outcomes. |
| 07 | Due expiry wins over rejoin or a due vote. Dead spectators cannot advance gameplay while every living player is disconnected; gameplay resumes with remaining time, while seat/session/drain expiry keeps running. A content-mismatched takeover neither replaces the connection nor extends grace. |

Phase 13 build checks also establish the production/capacity-build relationship
from `docs/benchmark-plan.md`, including adapter absence in the production
artifact. Behavioral transcripts do not prove equal performance between builds.

## Fuzzing

Fuzz parsers, canonicalization, protocol decoding, and state-machine boundaries
in bounded isolated processes.

- Keep targets narrow and free of external I/O.
- Bound input size, time, memory, recursion, operations, and concurrency.
- Preserve/minimize failures and promote useful cases to deterministic
  regressions.
- A panic, hang, escape, external fetch, unbounded allocation, or cleanup leak
  is a failure.

## CI And Evidence

- Run the pinned commands from `docs/mvp-contract.md` on a clean checkout.
- Run tests in parallel; repeat only concurrency/chaos suites under recorded
  deterministic seeds.
- Preserve the first intermittent failure. The suite remains failing or
  `INCONCLUSIVE` until the cause is understood.
- Coverage percentage is informational. Missing behavior is not excused by a
  high number.
- Phase evidence links CI, important test names, fixture/schema versions, and
  any manual or environment-specific result.
