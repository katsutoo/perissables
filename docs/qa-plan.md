# QA Plan

Status: Normative release-behavior policy
Owner: Project team
Updated: 2026-07-11

## Verdicts

- `PASS`: all applicable acceptance scenarios passed on the identified artifact and no known open defect affects a required scenario.
- `PASS WITH KNOWN ISSUES`: all required scenarios passed, but one or more documented `Major`/`Minor` defects remain with an owner, target date, workaround, and explicit residual-risk acceptance.
- `FAIL`: a release criterion failed or a blocker/critical defect is reproducible.
- `BLOCKED`: a prerequisite prevented meaningful execution.
- `INCONCLUSIVE`: evidence is ambiguous, contradictory, or intermittently failing.

Every report has `stage: pre_release` or `stage: final` plus exactly one verdict from the exhaustive list above. Phase 17 requires `stage: pre_release` `PASS` against the production-profile core artifact and scenarios `1..=9` plus `11`; it is not a Release-ready verdict. After Phase 18 creates Steam/package/signing/rollback artifacts, every scenario including `10` is rerun with `stage: final`. Release-ready requires final `PASS`; `PASS WITH KNOWN ISSUES` does not satisfy that gate. Retrying until green, changing the tested artifact/configuration, or hiding an intermittent result is forbidden.

Every report also gives a separate release recommendation: `ship`, `hold`, or `no recommendation`. A final `PASS` is required before QA may recommend `ship`; `FAIL` requires `hold`; `BLOCKED` or `INCONCLUSIVE` requires `no recommendation` unless a documented release criterion independently requires `hold`. The named release owner makes the decision and records any departure from the QA recommendation.

## Evidence Record

Every run records the Git SHA and dirty state, artifact checksum, build command/profile/features, target environment and URL, non-secret configuration, OS/runtime/GPU/driver, display resolution/scaling, account role, synthetic test data, exact steps, expected result, actual result, logs/network/durable effects, cleanup, and artifact location. Secrets, tickets, tokens, provider subjects, and personal data are redacted.

## Release-ready Client Matrix

| Target | Required environment |
| --- | --- |
| Linux | Ubuntu 24.04 LTS x86_64; freeze and record exact ISO/image digest plus package snapshot before the run; `1280x720` and `1920x1080` |
| Windows | Windows 11 24H2 x86_64; freeze and record full OS build/UBR and update IDs before the run; `1280x720` and `1920x1080` |

Record GPU/driver and windowed/fullscreen mode. Keyboard-only operation and two behavior-identical UI variants are product requirements from `docs/mvp-contract.md`; any additional supported controller gets its own matrix row before release. Other distributions and Steam Deck are exploratory until explicitly added to the support contract.

## Required Scenarios

1. Install/start with a clean user-data directory; missing optional local settings recover to documented defaults without creating authoritative state.
2. Complete the canonical lobby -> story selection -> unique character lock -> all-ready owner start -> story -> checks -> combat -> summary acknowledgement/timeout -> reset lobby flow three consecutive times without stale state.
3. Exercise all three themes and both UI variants with identical action/focus order and hit targets. At each required resolution and `100%`/`200%` OS scaling, all required text fits its declared box without clipping/overlap, keyboard focus is always visible with at least a `2px` indicator, every action is reachable without a pointer, and each audio channel independently reaches `0%` and `100%` without changing another channel. Missing assets produce the documented fallback rather than a crash.
4. Run deterministic 2-, 3-, and 4-player hosted sessions and verify revision/event convergence plus fifth-player and capacity rejection.
5. Disconnect during world, story choice, and combat turns; rejoin, acknowledge resync, and verify no duplicated action/audio event or lost authoritative state.
6. Restart/deploy the hosted server during a run, restore the exact durable state, rejoin within grace, and finish the run.
7. Exercise invalid/expired Steam proof, token, version, checksum, sequence, rate, and frame/message boundaries and verify stable errors without leaked internals.
8. Inject persistence unavailability, corrupt current snapshot, unavailable backup, and graceful-shutdown deadline behavior; verify the documented fail-closed result.
9. In the local single-player creator-test path, install a valid Tier 1 aggregate community pack and reject traversal, alias, oversized, malformed, unsupported-version, checksum-mismatch, and unattested Tier 2 packs without external access or residue. Verify hosted multiplayer remains built-in-only.
10. Build/package both targets, inspect contents for debug-only characters/keys/features and secrets, launch the shipped binary, and verify Steam depot layout plus rollback artifact.
11. Restore the current save and client-settings versions plus the immediately previous positive version of each when one exists through the tested artifact; verify exact migrated state, idempotent re-open, malformed/oversized/future/too-old rejection, documented default recovery, and rollback behavior without destructive downgrade. Version `1` records previous-version coverage as not applicable rather than fabricating a version `0` fixture.

## Operations Oracles

- `/healthz` succeeds while the process event loop is alive and does not claim dependency readiness.
- `/readyz` fails until the exclusive lease, restoration, and writable persistence prerequisites are satisfied; it remains successful when only new-admission capacity is exhausted and fails before graceful drain begins.
- Expected operational failures use stable external codes and structured internal context without credentials or personal data.
- Shutdown stops admission, preserves the last durable checkpoint, awaits owned tasks, and exits within the locked deadline.

## Severity And Cleanup

- `Blocker`: data/authority corruption, credential exposure, unrecoverable run loss, universal startup failure, or no safe workaround.
- `Critical`: major supported flow unavailable, cross-player authority failure, repeatable crash, or severe persistence/reconnect failure.
- `Major` and `Minor`: degraded behavior with a safe workaround or limited presentation impact.
- Use synthetic accounts/data, avoid production and real third-party effects, restore local/staging state, and report anything that could not be cleaned up.
- Reports list every scenario as passed, failed, blocked, inconclusive, or not applicable; they identify the tested build/environment/account roles, release recommendation, decision owner, cleanup result, and untested residual risk.
