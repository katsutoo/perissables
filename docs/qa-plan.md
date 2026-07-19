# QA Plan

Status: Normative release-behavior policy
Owner: Project team
Updated: 2026-07-19

Authority: release-behavior methodology and verdicts are normative here; expected product/protocol values derive from the named locked sections of `docs/mvp-contract.md`.

## Verdicts

- `PASS`: all applicable acceptance scenarios passed on the identified artifact and no known open defect affects a required scenario.
- `PASS WITH KNOWN ISSUES`: all required scenarios passed, but one or more documented `Major`/`Minor` defects remain with an owner, target date, workaround, and explicit residual-risk acceptance.
- `FAIL`: a release criterion failed or a blocker/critical defect is reproducible.
- `BLOCKED`: a prerequisite prevented meaningful execution.
- `INCONCLUSIVE`: evidence is ambiguous, contradictory, or intermittently failing.

Every report has `stage: pre_release` or `stage: final` plus exactly one verdict from the exhaustive list above. Phase 17 requires `stage: pre_release` `PASS` against the production-profile server image and both OS client-candidate digests for scenarios `1..=9` plus `11`; it is not a Release-ready verdict. After Phase 18 creates Steam/package/signing/rollback artifacts, every scenario including `10` is rerun with `stage: final`, and the final benchmark rerun is linked from the report. Release-ready requires final `PASS`; `PASS WITH KNOWN ISSUES` does not satisfy that gate. Retrying until green, undeclared artifact/configuration changes, or hiding an intermittent result is forbidden. Scenario-defined restart/deploy/rollback transitions are allowed only when all before/after digests/configurations and the intended final deployed state are declared before the run.

Every report also gives a separate release recommendation: `ship`, `hold`, or `no recommendation`. A final `PASS` is required before QA may recommend `ship`; `FAIL` requires `hold`; `BLOCKED` or `INCONCLUSIVE` requires `no recommendation` unless a documented release criterion independently requires `hold`. The named release owner makes the decision and records any departure from the QA recommendation.

## Evidence Record

Every run records the Git SHA and dirty state, every artifact checksum/digest used, build command/profile/features, target environment and URL, before/after non-secret configuration, intended final deployed state, OS/runtime/GPU/driver, display resolution/scaling, account role, synthetic test data, exact steps, expected result, actual result, logs/network/durable effects, cleanup, and artifact location under `artifacts/qa/<git-sha>/<run-id>/`. Secrets, tickets, tokens, provider subjects, and personal data are redacted.

## Release-ready Client Matrix

| Target | Required environment |
| --- | --- |
| Linux | Ubuntu 24.04 LTS x86_64; freeze and record exact ISO/image digest plus package snapshot before the run; `1280x720` and `1920x1080` |
| Windows | Windows 11 24H2 x86_64; freeze and record full OS build/UBR and update IDs before the run; `1280x720` and `1920x1080` |

Record GPU/driver and windowed/fullscreen mode. Keyboard-only operation and all shipped behavior-identical built-in UI variants (at least two) are product requirements from `docs/mvp-contract.md`; any additional supported controller gets its own matrix row before release. Other distributions and Steam Deck are exploratory until explicitly added to the support contract.

## Required Scenarios

1. Install/start with a clean user-data directory; missing optional local settings produce the exact canonical `Client-Local Settings v1` object in memory without creating authoritative state. Exercise every field/conflict-group/keyboard-path boundary and verify malformed/future settings are preserved, reported once, and never destructively rewritten. Inject failure at temporary write, file sync, rename, and directory sync; reopen and prove the previous valid file remains byte-identical.
2. Complete the canonical lobby -> story selection -> unique character lock -> all-ready owner start -> story -> checks -> combat -> summary acknowledgement/timeout -> reset lobby flow three consecutive times without stale state.
3. Exercise all three themes and every shipped built-in UI variant with identical action/focus order and hit targets. At each required resolution and the contract's `100`/`200` UI scales, all required text fits its declared box without clipping/overlap, keyboard focus is always visible with at least a `2px` indicator, every action is reachable without a pointer, and each audio channel independently reaches `0` and `100` without changing another channel. Validated presentation-time image/audio/skin failures use the exact `Presentation Fallback v1` assets/silence/diagnostic behavior; missing required content and unattested Tier 2 bytes reject activation instead of falling back. Corrupt each boot-critical fallback resource in an isolated package and verify deterministic package-integrity startup failure without recursion.
4. Run deterministic 2-, 3-, and 4-player hosted sessions and verify revision/event convergence plus fifth-player and capacity rejection.
5. Drop the first initial `join_response` and prove the same identity/mode/session/`client_join_nonce` recovers the same IDs with a fresh token and no duplicate seat/session. Then disconnect during world, story choice, and combat turns; complete and durably acknowledge token rotation plus resync, repeat the disconnect using the new token, and verify old/pending generations, generation-guarded takeover closes, duplicated action/audio events, and authoritative state follow the handoff contract. Drop `rejoin_response` before receipt and prove the old token creates a replacement handoff; separately receive the response then drop before `resync_ack` and prove the known pending token resumes it. Neither path may grant two active connections.
6. Run both `restore-clean-v1` and `restore-crash-v1` during a run. Measure separately: clean shutdown plus synced `released` marker at or below `20s`; clean successor process-start-to-ready at or below `20s`; crash lease wait/lock acquisition at or below `20s`; post-acquisition restoration at or below `20s`; and total crash process-start-to-ready at or below `40s`. Verify exact durable state, no split owner/admission, rejoin within persisted grace, and successful run completion.
7. Exercise invalid/expired Steam proof, token, version, checksum, sequence, rate, and frame/message boundaries and verify stable errors without leaked internals.
8. Inject persistence unavailability, soft/hard WAL limits, trailing incomplete WAL frame, complete interior WAL corruption, corrupt current snapshot, no valid backup, artifact/file/global-and-session-quarantine caps, and graceful-shutdown deadline behavior. Verify exact unchanged sequences/revisions, reserved terminal writes, bounded scans, safe trailing-frame truncation, selective/global WAL recovery, damaged-session quarantine plus `session_not_found`, readiness only after durable recovery/quarantine, operator evidence, and no clean lease release on deadline failure.
9. In the local single-player creator-test path, install a valid Tier 1 aggregate community pack and reject traversal, alias, oversized, malformed, unsupported-version, checksum-mismatch, and unattested Tier 2 packs without external access or residue. Verify hosted multiplayer remains built-in-only.
10. Build/package both targets, inspect contents for debug-only characters/keys/features and secrets, launch the shipped binary, verify Steam depot layout, validate every digest/version/command in `rollback-manifest.json`, and execute the complete isolated-staging server plus Steam test-branch rollback within its locked deadlines. Verify exact continued authoritative state, coordinated protocol compatibility, signatures/checksums, and no improvised artifact or destructive downgrade.
11. Restore the current save and client-settings versions plus the immediately previous positive version of each when one exists through the tested artifact; verify exact migrated state, idempotent re-open, malformed/oversized/future/too-old rejection, the exact canonical settings defaults, the prepared next-version rollback reader when applicable, and rollback behavior without destructive downgrade. Version `1` records previous-version coverage as not applicable rather than fabricating a version `0` fixture.

## Operations Oracles

- `/healthz` succeeds while the process event loop is alive and does not claim dependency readiness.
- `/readyz` fails until the exclusive lease, restoration/quarantine pass, and writable persistence prerequisites are satisfied; it remains successful when only new-admission capacity is exhausted and fails before graceful drain begins. Clean successors consume the synced release marker immediately; crash successors cannot become ready before lease expiry and lock acquisition.
- Expected operational failures use stable external codes and structured internal context without credentials or personal data.
- Shutdown stops admission, preserves the last durable checkpoint, awaits owned tasks, and exits within the locked deadline.

## Severity And Cleanup

- `Blocker`: data/authority corruption, credential exposure, unrecoverable run loss, universal startup failure, or no safe workaround.
- `Critical`: major supported flow unavailable, cross-player authority failure, repeatable crash, or severe persistence/reconnect failure.
- `Major` and `Minor`: degraded behavior with a safe workaround or limited presentation impact.
- Use synthetic accounts/data, avoid production and real third-party effects, restore local/staging state, and report anything that could not be cleaned up.
- Reports list every scenario as passed, failed, blocked, inconclusive, or not applicable; they identify the tested build/environment/account roles, release recommendation, decision owner, cleanup result, and untested residual risk.
