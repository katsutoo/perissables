# Test Strategy

Status: Normative automated-test policy
Owner: Project team
Updated: 2026-08-11

Product and protocol behavior comes from `docs/mvp-contract.md`. This document
defines how automated tests earn confidence without becoming a second
implementation or a paperwork system.

## Principles

- Test observable behavior through public interfaces.
- One test protects one behavior; table cases may cover several inputs of that
  same behavior.
- Assert concrete outputs and negative space. A rejected action must not mutate
  state, consume a turn, spend resources, or emit gameplay events.
- Inject time, RNG state, IDs, transport schedules, storage failures, and
  identity-provider responses.
- No sleeps for synchronization, wall-clock dependence, shared mutable fixtures,
  order dependence, or retry-until-green.
- Use real deterministic collaborators when cheap. Substitute only the boundary
  that is unsafe, slow, nondeterministic, or external.

## Fail-First Evidence

Every new test is observed failing for the intended reason before acceptance.

- Regression fixes preserve the failing test name, command, and failure
  signature in the issue or change record.
- Critical protocol, persistence, security, and migration invariants preserve
  equivalent red/green evidence.
- Routine test-first feature work does not require a permanent test-only commit,
  tree identifier, or standalone evidence bundle.
- Pure refactors state why no new behavior test is needed.

## Layers

| Layer | Protects |
| --- | --- |
| Unit/property | Dice, combat, story rules, bounds, canonicalization, migrations |
| Schema/conformance | Shared accepted/rejected JSON, TMX, ZIP, paths, and bytes |
| State-machine transcript | Authoritative revisions, events, rejection, reset |
| Protocol contract | DTOs, directions, versions, sequences, projections, errors |
| Multi-client integration | 2/3/4-client convergence, reconnect, replay, capacity |
| Persistence integration | The selected store, commit boundary, crash, migration, corruption |
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
  gameplay phase, stable ordering, death/wipe, and three-run reset.
- Protocol tests cover fragmentation, malformed DTOs, wrong direction/version,
  stale/duplicate/gap sequences, recipient-specific projections, queue refusal,
  reconnect handoff, and slow writers.
- Identity tests cover wrong app/identity/token/generation, dropped handoff
  responses, takeover, expiry, and absence of the local adapter in release
  features.
- Content tests cover traversal, aliases, duplicate keys, decompression/resource
  limits, disabled XML external access, graph errors, checksum stability, and
  fail-closed sandbox behavior.
- Persistence tests are defined after Phase 10 and run against the real selected
  adapter. Internal begin/write/commit/checkpoint or equivalent fault injection,
  disk/full-space behavior, migration, crash, and corruption belong here rather
  than in manual QA.

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
