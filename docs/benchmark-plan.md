# Benchmark Plan

Status: Normative performance experiment policy
Owner: Project team
Updated: 2026-08-11

Performance is established by production-mode measurements, not by contract
detail or code inspection. Product targets come from `docs/mvp-contract.md`.

## Decisions Supported

1. **Phase 10 storage decision:** choose the store, state representation,
   durability boundary, commit cadence, and operational budgets.
2. **Phase 12 release capacity:** determine whether the production-shaped server
   meets the session/player, latency, reliability, and headroom goals.
3. **Phase 12 client gate:** calibrate the reference machine and verify the 60
   FPS experience for frozen gameplay scenes.
4. **Phase 13 final validation:** rerun applicable gates against exact final
   server/package digests.

Phase 02 may collect directional timings for instrumentation sanity, but those
numbers are not release claims.

## Common Rules

- Correctness tests pass before timing.
- Use release/production profiles, exact features/targets, and shipped assets.
- Record Git SHA/dirty state, artifact digests, commands, tool versions,
  environment, workload version/checksum, duration, operations, errors, and raw
  result location.
- Separate warmup, measured work, setup, and recovery.
- Report independent runs in addition to operation samples.
- Randomize or interleave baseline/candidate order.
- Preserve timeouts, late starts, dropped work, and errors. Faster incorrect
  output fails.
- Predeclare practical thresholds after measuring the environment noise floor
  and before candidate comparison.
- Never average percentiles.

## Phase 10 Storage Spike

### Question

Which mature storage design gives the required recovery semantics with the
lowest operational and implementation cost?

At minimum compare bundled SQLite on the actual Railway volume with managed
PostgreSQL if SQLite misses a gate. Evaluate full snapshots, bounded
deltas/checkpoints, and candidate durability boundaries using implemented
session DTOs rather than invented byte blobs.

### Workloads

- Typical lobby, world, story, combat, and summary states.
- Maximum valid state generated through public rules/schema.
- One changing session, expected concurrent load, and the release capacity goal.
- Semantic actions, continuous movement, disconnect/rejoin handoff, and session
  end.
- Clean shutdown, immediate crash after acknowledged work, crash during a write,
  recovery, checkpoint/maintenance, and storage-unavailable behavior.

### Metrics And Decision Record

Record p50/p95/p99/max commit and acknowledgement latency, achieved operations/s,
bytes written, write amplification, fsyncs, CPU, peak RSS, volume growth,
recovery duration, correctness/errors, backup/restore behavior, and operator
steps.

The decision must state:

- production environment and independent run count;
- selected store, schema/state representation, and acknowledgement boundary;
- measured batching/cadence, queue, storage, and recovery budgets;
- headroom and uncertainty;
- rejected alternatives; and
- the smallest workload that would invalidate the decision.

A tiny local database or debug binary cannot establish this gate.

## Server Release Workload

The final capacity workload is open-loop:

- target `64` sessions with `4` authenticated clients each;
- up to `20` scheduled gameplay inputs/client/s;
- mixed world, story-vote, and combat sessions;
- frozen typical payload distribution plus separate maximum-payload tests;
- load points at `25%`, `50%`, `75%`, and `100%`;
- `2 minutes` warmup, `10 minutes` steady state, and `5` independent
  server/driver process runs per gated load point.

The Phase 12 fixture freezes exact seeds, traces, expected outcomes/reasons,
encoded-size distributions, and operation counts. Lower load points have
complete schedules; they are not truncated high-load traces.

Authentication setup occurs outside the measured interval using project-owned
authorized Steam test accounts. After the production limiter is implemented,
derive and provision the real source-IP topology needed to remain within every
source/session/identity bucket. Header spoofing, auth bypasses, one identity
across seats, and load against Steam itself are forbidden. Missing accounts,
authorized staging, or real limiter-keyed routes makes the hosted benchmark
`BLOCKED`.

### Server Gates

- Processing latency p95 below `100 ms`, p99 below `200 ms`.
- Same-region scheduled-send-to-receive p95 below `150 ms`, p99 below
  `300 ms`.
- Offered and achieved rates, late starts, timeouts, and every expected/actual
  outcome count agree with the frozen workload.
- No unexpected error, disconnect, unreported missed tick, or healthy-client
  queue overflow.
- Sustained CPU, peak RSS, volume, and network use retain at least `30%`
  measured allocation headroom.
- p99 bounded-queue occupancy remains below the frozen Phase 12 headroom gate.
- Storage commit/recovery gates use the Phase 10 decision rather than obsolete
  pre-spike assumptions.

Every gated outcome class needs enough samples for its reported percentile;
rare classes are reported without a fabricated p99.

## Recovery And Boundary Workloads

Run separately from steady state:

- clean restart and crash recovery at the selected acknowledgement boundary;
- all-client reconnect spread within real limiter budgets;
- one backpressured client per session while healthy clients continue;
- selected-store maintenance/checkpoint behavior;
- maximum valid state, event, result, and resync payloads; and
- storage full/unavailable behavior through the safe test mechanism defined by
  Phase 10.

Each workload declares setup, stop conditions, expected state/checksum,
timeout, cleanup, and whether it is a correctness gate or diagnostic.

## Client Calibration And Workload

Use one named project-owned reference machine. Record CPU, RAM, GPU, driver, OS,
display, power/thermal mode, exact build, and unrelated load.

Before setting the gate:

1. measure actual refresh and presentation callback behavior with the smallest
   production render path;
2. measure timer/driver interval noise and run-to-run variation;
3. freeze a tolerance and missed-refresh budget that exceed that noise but still
   protect the 60 FPS experience; and
4. publish the calibration artifact before candidate measurement.

The gate is then p99 presentation interval at or below one measured refresh
interval plus the frozen calibrated tolerance, within the frozen missed-refresh
budget, with no unexplained update backlog, no silently dropped logical update,
and no frame above `100 ms`.

Frozen `world`, `story`, and `combat` traces run at both required
resolutions after warmup for at least `10 minutes`, with `5` independent
runs per OS/resolution/scene. Report presentation and update distributions
separately. Exercise the fixed-update catch-up/drop policy in its own
correctness workload rather than mixing a forced stall into normal frame data.

## Required Metrics

| Area | Metrics |
| --- | --- |
| Latency | p50, p95, p99, max, samples, histogram range/precision |
| Load | offered/achieved ops/s, late starts, ticks/broadcasts |
| Reliability | outcomes, errors, timeouts, disconnects, dropped/coalesced work |
| Resources | CPU user/system, peak RSS, allocations when useful |
| I/O | database bytes/fsyncs, volume growth, network ingress/egress |
| Queues | occupancy distribution, refusal/overflow/coalescing |
| Client | presentation/update intervals, missed refreshes, CPU/RSS/GPU |
| Artifact | executable/package compressed and uncompressed size |
| Startup | cold start, process-to-ready, restored sessions |

## Harness And Artifacts

The repository-owned driver uses shared protocol DTOs, an open-loop monotonic
scheduler, bounded tasks/connections, explicit expected results, and lossless
raw samples or compatible histograms. Calibrate it independently and prove that
it sustains at least `125%` of target offered rate without becoming the
bottleneck.

Store versioned synthetic fixtures in source control. Keep redacted run
manifests, raw samples/histograms, logs, commands, and summaries under
`artifacts/benchmarks/<git-sha>/<run-id>/`; raw operational artifacts remain
out of Git and follow finite retention.

Verdicts are `PASS`, `FAIL`, `BLOCKED`, or `INCONCLUSIVE`. Uncertainty
crossing a practical threshold is `INCONCLUSIVE`. Profiling explains a
measured result but never substitutes for an unprofiled benchmark.
