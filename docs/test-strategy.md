# Test Strategy

Status: Normative verification policy
Owner: Project team
Updated: 2026-07-10

## Purpose

Tests protect observable behavior and locked invariants. They do not replace QA, security review, or benchmarks. Every behavior is tested in the phase that introduces it; Phase 17 runs the resulting suite rather than adding overdue regression coverage.

## Test Rules

- Test public interfaces and exact state/event/error contracts, not private implementation details.
- Every new regression/behavior test must be observed failing for the intended reason before the implementation or fix is accepted. Preserve the pre-fix commit/tree identifier, exact command, test name, and expected failure signature in the phase evidence; then record the passing run against the implementation. For new code developed test-first, the failing test-only commit supplies this evidence. Pure refactors that add no behavior test state why this rule is not applicable.
- Use visible arrange-act-assert structure. One test protects one behavior; table cases may cover multiple values of that behavior.
- Assert concrete outputs and negative space: rejected input must not mutate gameplay state, consume a turn, or emit a gameplay event. An authenticated expected `input_seq` still advances replay metadata and, after Phase 16, writes that metadata/outcome to the WAL as required by the protocol contract.
- Operational-error tests assert the stable error kind and durable side effects, not merely that an error occurred.
- No sleeps, wall-clock dependence, random seeds from the environment, shared mutable fixtures, order dependence, or retry-until-green behavior.
- Inject clocks, logical ticks, RNG state, IDs, transport schedules, storage faults, and external-service responses.
- An unexplained intermittent failure is a failing suite. Preserve its first evidence and report the phase `INCONCLUSIVE` until resolved.

## Layers

| Layer | Protects | Required examples |
| --- | --- | --- |
| Unit/property | Pure rules and bounded arithmetic | All dice values, limits, turn ordering, migrations, path grammar |
| Schema/conformance | Shared byte-level behavior | Valid/invalid documents, RFC 8785 vectors, checksum corpus, portable paths |
| State-machine transcript | Authoritative transitions | Story branches, combat actions, death/wipe, reset, stable revisions/events |
| Protocol contract | Wire compatibility and rejection | Every payload/direction, versions, sequence boundaries, close/error behavior |
| Multi-client integration | Convergence and authority | Scripted 2/3/4 clients, simultaneous joins, stale/replayed input, reconnect |
| Persistence integration | Crash-safe continuation | Current/previous migration, corrupt/future/oversized state, every write fault point |
| Release-artifact smoke | Shipped behavior | Startup, full run, hosted join/rejoin/restore on Linux and Windows |

## Mandatory Boundaries

- Test zero, one, typical, limit minus one, limit, and limit plus one for every count, size, queue, rate, timeout, index, and version boundary.
- Exhaustively test the d100 resolver for `0..=100`; test RNG v1 with official ChaCha20 known-answer vectors, rejection-threshold words, serialized resume at every word index, transaction rollback, and fixed expected outputs, never probabilistic frequency assertions.
- Inject separate fake monotonic and Unix wall clocks. Use monotonic time for token buckets, heartbeat, turn/vote/summary deadlines, idle timeout, snapshot cadence, and retry delays. Test absolute session/token expiry across restart downtime, `1s` wall-clock rollback tolerance, rollback beyond `1s` fail-closed behavior, and forward jumps. Test just before, at, and just after every boundary.
- Run protocol tests with deterministic fragmentation, duplicate, stale, gap, reconnect, and slow-writer schedules. Persist the seed and schedule for every failure.
- Verify every valid persisted state produces a `ClientResyncState` within its separate cap and never contains token digests, RNG state, or hidden player data.

## Pack Fuzzing

- Fuzz JSON, ZIP, path, checksum, image-header, audio-header, and TMX/XML boundaries in an isolated local process with no network access.
- PR smoke target: `60s` per parser family, `512 MiB` RSS ceiling, `1s` per-input timeout.
- Scheduled target: at least `30 minutes` per parser family using retained corpora.
- Required oracle: no panic, hang, external fetch, traversal, unbounded allocation, cleanup leak, inconsistent checksum, or acceptance beyond a locked limit.
- Preserve minimized reproducers and their exact toolchain/seed. Never skip a discovered input to restore green CI.

## CI And Evidence

- Run the pinned commands from `docs/mvp-contract.md` on a clean checkout with committed `Cargo.lock`.
- Run tests in parallel and repeat concurrency/chaos suites under at least three recorded seeds before phase completion.
- If features become mutually exclusive, replace `--all-features` with an explicit supported matrix in CI and this document.
- A phase-completion checkbox links the CI run, relevant test names, fixture/conformance version, red/green evidence required above, and any manual evidence. Coverage percentage is informational and never substitutes for missing behavior.
