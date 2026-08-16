# Benchmark Plan

Status: Normative performance experiment policy
Owner: Sole developer
Updated: 2026-08-16

Performance is established by production-mode measurements, not by contract
detail or code inspection. Product targets come from `docs/mvp-contract.md`.

## Decisions Supported

1. **Phase 10 regional deployment and drain:** choose the Railway regional
   topology, session routing, drain deadline, and rollback procedure.
2. **Phase 12 capacity measurement:** determine instance sizing, admission
   limits, scaling needs, latency, reliability, and headroom.
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
- Separate warmup, measured work, setup, and drain/boundary work.
- Report independent runs in addition to operation samples.
- Randomize or interleave baseline/candidate order.
- Preserve timeouts, late starts, dropped work, and errors. Faster incorrect
  output fails.
- Predeclare practical thresholds after measuring the environment noise floor
  and before candidate comparison.
- Never average percentiles.

## Phase 10 regional deployment and drain spike

### Question

What is the smallest Railway-only topology within the entry-level subscription
budget that gives worldwide parties acceptable nearby latency to one owning
in-memory session while keeping drain and rollback simple for one developer?

### Workloads

- Create/join and latency probes worldwide against an initial
  Americas/Europe/Asia candidate topology.
- Same-area and cross-area parties through lobby, world, vote, combat, and
  summary traffic.
- Controlled latency-impairment playtests that freeze the highest acceptable
  cross-area movement/input latency before choosing the region set.
- Client disconnect/rejoin routed back to the owning live process.
- Mark one process unready, refuse new admission, finish short and full-length
  sessions, and expire the drain deadline.
- Forced stop, process crash, regional outage, replacement deployment, and
  rollback with the documented run-lost behavior.

### Metrics and decision record

Record p50/p95/p99/max scheduled-send latency by source/host region, route and
rejoin correctness, drain duration, sessions completed/lost, admission results,
CPU, peak RSS, network use, deployment duration, and operator steps.

The decision states the smallest Railway region set within the entry-level
budget, automatic worst-latency/median-latency placement, owning-process routing,
measured drain deadline, rollback procedure, headroom/uncertainty, rejected
alternatives, and the smallest workload that would invalidate it.

## Server capacity workload

The capacity workload is open-loop. Its single-instance maximum is a measurement
point, not a launch concurrency promise:

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

Steam authentication is qualified separately, outside capacity timing, with a
small pool of project-owned authorized test accounts. The capacity workload uses
synthetic identities issued only by an isolated non-release benchmark identity
adapter; that adapter cannot compile into release features/packages. The driver
uses a realistic identity/session distribution and the same post-auth session
and rate-limit paths as production without sending load to Steam.

Header spoofing, one identity reused across seats, and load against Steam itself
are forbidden. Missing authorized Steam staging makes the Steam authentication
check `BLOCKED`; it does not block an otherwise valid isolated capacity
measurement. Missing the production limiter path or release-adapter absence proof
makes the capacity gate `BLOCKED`.

### Measurement validity and targets

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
- Missing the `64`-session resource target informs instance sizing, scaling, or
  admission limits; it does not permit correctness loss or create an unsupported
  concurrency claim.

Every gated outcome class needs enough samples for its reported percentile;
rare classes are reported without a fabricated p99.

## Drain and boundary workloads

Run separately from steady state:

- clean drain, drain deadline, forced stop, process crash, and stable run loss;
- all-client reconnect spread within real limiter budgets and owning-process
  routing;
- one backpressured client per session while healthy clients continue;
- maximum valid state, event, result, and resync payloads; and
- regional endpoint loss and rollback through the safe mechanism defined by
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
| I/O | Network ingress/egress by region and route |
| Queues | occupancy distribution, refusal/overflow/coalescing |
| Client | presentation/update intervals, missed refreshes, CPU/RSS/GPU |
| Artifact | executable/package compressed and uncompressed size |
| Startup | Cold start, process-to-ready, drain, and replacement deployment |

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
