# Benchmark Plan

Status: Normative performance experiment
Owner: Project team
Updated: 2026-07-11

## Decisions And Gates

Benchmarks validate production-profile artifacts after correctness tests pass and are rerun on the final release artifacts when packaging can affect results. Server processing latency and client frame/update behavior are Release-ready gates; internet RTT, cold start, memory, network use, and artifact size are reported for capacity and regression decisions. A faster result with errors, dropped work, missed validation, or changed behavior fails.

## Environment Lock

- Build with the pinned Rust toolchain, committed lockfile, production release profile, exact target/features, and shipped assets.
- Before Phase 17 load testing, record and freeze the Railway service plan, region, CPU/memory limits, volume configuration, replica count, runtime variables, and deployment digest. Production must use the same or a demonstrably larger shape.
- Record the load-generator host/region, CPU/memory, OS/kernel, network path, tool versions, power/thermal state where applicable, and unrelated load.
- Establish the first client baseline on a named project-owned reference machine. Record immutable CPU, RAM, GPU, driver, OS, display, and power-mode details; future comparisons use the same machine or are labeled directional.

## Server Workload

- Artifact: deployed `les-perissables-server` release build with one active instance and fixture `benchmark-fixture-v1`. Its archive, canonical checksum, schema/rules versions, initial snapshots, action traces, expected counters, and expected payload-size distribution live under `benchmarks/fixtures/v1/` and are frozen before Phase 13 exits; changing any file creates `v2`, never silently updates v1.
- Fixture shape: one `512x512` CSV TMX map, `512` story nodes, `128` encounters, `64` characters, and `256` flags exercise maximum content-state bounds without making every client projection a maximum-size payload. The steady-state trace freezes encoded state payloads in the `2..=8 KiB` range with p50 at or below `4 KiB`, encoded event payloads in the `0.5..=4 KiB` range when present, and no resync payload during the normal measured interval. Its precomputed expected aggregate ingress plus egress at `100%` offered rate must be at or below `40 MiB/s`, leaving at least `10 MiB/s` beneath the gate for measured protocol/transport overhead; a fixture that cannot prove that budget is invalid before load testing. The fixture manifest records exact archive/file/expanded bytes, per-message expected encoded sizes, expected traffic by message class, and expected validation result. It uses synthetic text/media and no proprietary production account data. Boundary-size frames are measured separately by `max-payload-v1`.
- Load model: open-loop, `64` sessions with `4` authenticated clients each (`256` clients), `20` scheduled inputs per client per second (`5,120 msg/s` offered).
- Session distribution is fixed at `32` world, `16` story-vote, and `16` combat sessions. Trace seed is `0x4c505f42454e4348`. World clients send valid held-direction replacements; story clients alternate valid choices on the current open vote; in combat, only the actor current at batch start sends the fixture's first legal action while the other three send `combat_attack` with the reserved syntactically valid but nonexistent fixture ID `ent_invalid_benchmark_target`, guaranteeing `target_unavailable` even if the turn advances during sequential staging. Over a complete `10-minute` measured run the target is `3,072,000` scheduled inputs: `1,536,000` valid world, `768,000` valid story, `192,000` valid combat, and `576,000` expected target rejections. Fixture resets and actor rotation are part of the checked trace, not benchmark setup hidden from counters. Late/dropped starts, unexpected rejection, and deviation from these counts are failures, not discarded samples.
- Run shape: `2 minutes` warmup, `10 minutes` measured steady state, `5` independent process/deployment runs. Randomize candidate/baseline order for later comparisons.
- Processing latency starts after a complete expected-sequence input passes protocol, authentication, rate, and sequence validation; it ends when that input's `input_result` is handed to the connection's WebSocket writer after any required WAL commit. Correlate by `input_seq`; applied/rejected/superseded outcomes are reported separately.
- Also record scheduled client-send to matching client-receive latency from the same-region load generator. This includes network and queueing but has no fixed Release-ready gate.
- Authentication happens before warmup using `256` project-controlled Steam publisher test accounts explicitly authorized for staging. Tickets are validated normally and are never recorded. If the approved account pool or Steam test authorization is unavailable, the hosted benchmark is `BLOCKED`; do not add an auth bypass, reuse one identity across seats, or load-test Steam itself.

## Recovery And Bound Workloads

- `restore-v1`: stop one isolated staging instance with `64` maximum-size durable sessions, deploy the identical digest, and measure process start through restored `/readyz`. Gate: all sessions checksum/replay exactly, no admission before completion, restoration at or below `20s`, and peak RSS within the normal headroom gate.
- `reconnect-v1`: after steady state, disconnect all `256` clients and reconnect their reserved seats uniformly over `30s`. Gate: no new capacity consumption, all resync acknowledgements complete within `10s` of each connection, zero duplicate semantic events, and zero unexpected errors.
- `slow-writer-v1`: one synthetic client per session stops reading for `5s` while others continue. Gate: only those `64` clients are disconnected through the documented bounded queue path; owner mailboxes remain below `50%` p99 and healthy clients remain connected.
- `compaction-v1`: begin at `48 MiB` WAL with all sessions dirty, run one compaction while processing `50%` offered rate, and inject a restart at each rename/fsync checkpoint in separate runs. The trace is precomputed to append at most `8 MiB` before durable replacement, leaving `8 MiB` below the hard cap. Gate: exact restore, no acknowledged loss, WAL at or below `16 MiB` afterward, no cap-blocked input, and processing latency remains within the normal gate.
- `max-payload-v1`: with gameplay input paused, test one largest-valid frame type at a time: resync, state, then a `64`-item event batch. For each type, cover all clients in four deterministic waves of `64` and start the next wave only after every prior writer queue drains; record bytes plus drain time for each wave. Gate: every frame remains within its individual cap, only the intended bundle is queued, zero queue overflow/disconnect occurs, queues return to zero within `5s` per wave, and the driver records actual peak and average throughput. The steady-state `50 MiB/s` budget does not apply to an instantaneous diagnostic burst; this workload cannot establish sustained capacity.

## Client Workload

- Artifact: shipped Linux/Windows release package, not a debug build or dev launcher.
- Run each required resolution for `10 minutes` in world traversal, dense dialog/UI, and representative four-player combat scenes after `2 minutes` warmup.
- Perform `5` independent runs per OS/resolution/scene on the frozen reference hardware.
- The 60 FPS gate requires frame-time p99 at or below `16.67 ms`, no unexplained update backlog, no silently missed fixed updates, and no frame above `100 ms` outside an explicitly identified OS/device interruption. The locked five-update/`250 ms` accumulator cap is exercised separately; dropped logical time must increment the exact diagnostic once. Report p50/p95/p99/max frame and update time separately.

## Required Metrics

| Area | Metrics |
| --- | --- |
| Server latency | p50, p95, p99, max, sample count, histogram range/precision |
| Load | offered/achieved msg/s, state broadcasts/s, tick interval/jitter, missed ticks |
| Reliability | errors by expected/unexpected kind, timeouts, disconnects, queue saturation, dropped/coalesced states |
| Resources | CPU user/system, peak RSS, allocations if available, disk writes/fsyncs, network bytes |
| Payload | p50/p95/p99/max inbound, state, event, resync, and persisted-snapshot sizes |
| Client | frame/update distributions, missed updates, CPU, peak RSS, GPU utilization when available |
| Artifact | compressed/uncompressed package and executable size, stripped/symbol status |
| Startup | cold process-to-ready, restoration time and restored session count |

Before the first gated run, convert the frozen Railway limits into numeric report values. Required headroom gates are sustained CPU at or below `70%` of allocation, peak RSS and volume use at or below `70%`, aggregate instance ingress plus egress at or below the fixed `50 MiB/s` MVP budget, p99 session-mailbox and writer-queue occupancy at or below `50%`, zero queue overflow for healthy clients, and p99 WAL group-commit latency below `50 ms`. WAL latency starts when a staged record enters the group and ends after `fsync`, so it includes the locked maximum `20 ms` batching delay. Client peak RSS is capped at `1 GiB` on the reference machine.

Run the workload at `25%`, `50%`, `75%`, and `100%` of offered rate with the same seeded state machine. The gated `100%` point must remain below saturation by the headroom rules; do not increase load until failure or affect production/third-party availability.

## Harness And Raw Result Contract

- The repository-owned benchmark driver is a non-shipping binary target in `les-perissables-integration-tests`, built with the pinned toolchain and the release profile. It uses the shared protocol DTOs, the frozen trace, an open-loop monotonic scheduler, and a bounded connection/task pool; it does not duplicate protocol serialization or game rules.
- Every scheduled operation records `scheduled_send_ns`, `actual_send_ns`, `receive_ns`, session/player/input sequence, expected outcome, actual outcome, encoded request/response bytes, timeout status, and run ID. Scheduled-send latency is the end-to-end open-loop measurement; actual-send latency is retained only to diagnose driver delay. Server processing latency comes from the server's correlated monotonic instrumentation and is never reconstructed from client clocks.
- Histograms use one pinned HDR-compatible implementation with `1 us` lowest discernible value, `60 s` highest trackable value, and three significant digits. Values beyond the range and timed-out operations are counted explicitly and fail the run; they are never clamped or omitted. Compatible raw histograms may be merged, but run-level percentiles are never averaged.
- The driver writes newline-delimited JSON operation records, one JSON run manifest, and lossless histogram files. The manifest includes schema version, Git SHA/dirty state, artifact digest, exact commands, tool versions, fixture version/checksum, seed, run order, monotonic clock source, start/end wall time, offered/achieved rate, warmup/measurement durations, client/session counts, timeout policy, histogram configuration, expected/actual counters, error totals, and host/environment capture.
- Warmup records are tagged and excluded by the fixed time boundary rather than deleted. Setup, authentication, and fixture reset happen outside the measured interval. The driver verifies expected state/counters before and after timing and exits non-zero on correctness, count, late-start, timeout, or serialization mismatch.
- Each independent run uses a fresh server process or deployment and a fresh driver process. Later baseline/candidate comparisons generate a seeded randomized or interleaved order before execution and retain that order; rerunning only an unfavorable side is forbidden.
- Before accepting the harness, run a loopback calibration with a no-op protocol fixture to establish scheduler delay, timestamp overhead, maximum sustainable generator rate, CPU use, socket/file-descriptor headroom, and the measurement noise floor. A gated run is invalid if the driver exceeds `70%` CPU allocation, misses more than `0.1%` of scheduled starts, or cannot sustain at least `125%` of the target offered rate in calibration.

## Pass And Interpretation

- Server processing latency passes only when every independent run satisfies p95 `<100 ms` and p99 `<200 ms`, all resource/headroom gates above, zero unexpected errors/disconnects, and no unreported missed tick.
- Client behavior passes only when every required matrix run meets the frame/update gate without correctness failure.
- Predeclare a practical regression budget before baseline/candidate comparison and measure the machine's noise floor. Report effect size, independent run count, harness uncertainty or run-to-run range; uncertainty crossing the budget is `INCONCLUSIVE`.
- Never average percentiles or infer them from a mean. Merge compatible raw histograms or report run-level distributions.

## Artifacts

Store exact build and benchmark commands, redacted configuration, environment capture, workload seed, raw NDJSON samples/histograms, run manifest, logs, and a Markdown summary under `artifacts/benchmarks/<git-sha>/<run-id>/`. Keep raw run artifacts out of source control and retain them according to the contract; versioned synthetic fixtures, workload definitions, expected counters, and redacted summaries remain reproducible source/release-history artifacts. Profiling output explains a measured result and never substitutes for an unprofiled benchmark.
