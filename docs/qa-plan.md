# QA Plan

Status: Normative release-behavior policy
Owner: Sole developer
Updated: 2026-09-12

QA exercises the product through a supported shipped entrypoint. Internal state
machines, parser matrices, and process-lifecycle fault injection belong to
automated tests. QA observes user-visible consequences and session effects.

## Verdicts

- `PASS`: every applicable required journey passed on the identified artifact
  with no release-blocking defect.
- `PASS WITH KNOWN ISSUES`: required journeys passed; only documented
  non-blocking defects remain with an owner and accepted residual risk.
- `FAIL`: a release criterion failed or a blocker/critical defect is
  reproducible.
- `BLOCKED`: a prerequisite prevented meaningful execution.
- `INCONCLUSIVE`: evidence is ambiguous, contradictory, or intermittent.

Every report also recommends `ship`, `hold`, or `no recommendation`.
Release-ready requires final `PASS`; the release owner records the decision.

## Evidence

Record:

- Git SHA, dirty state, artifact digest, build profile/features, and command;
- target environment/URL and non-secret configuration;
- OS build, GPU/driver, resolution, UI scale, and input mode;
- authorized Steam test-account roles and synthetic gameplay data;
- exact steps, expected result, actual result, visible and durable effects;
- logs or captures with secrets/personal data redacted; and
- cleanup and untested residual risk.

Preserve the first intermittent failure. Diagnostic reruns use the same artifact
and a stated protocol; a later pass does not erase it.

## Supported Matrix

| Target | Required environment |
| --- | --- |
| Linux | Exact Phase 01-pinned Steam Linux Runtime/container and supported x86_64 host baseline; `1280x720` and `1920x1080` |
| Windows | Frozen Windows 11 x86_64 build; `1280x720` and `1920x1080` |

Record the Linux runtime/container digest, host distribution, windowed/fullscreen
mode, UI scale, GPU, and driver. Keyboard-only
operation through the single shipped UI is required. Other Linux distributions,
Steam Deck, controllers, and additional display modes are exploratory until
added to the support contract.

## Required Journeys

1. **Clean install and startup.** Launch from a clean user-data directory,
   exercise settings creation/recovery through the UI, relaunch, and confirm no
   authoritative state or secret is stored locally.
2. **Three-run loop.** Complete lobby -> selection -> ready/start -> story ->
   check -> combat -> summary -> lobby three times without stale state. Observe
   the displayed timeout choice when everyone abstains, check-actor attribution,
   effect replacement, item consumption, and full-inventory reward feedback.
3. **Presentation and accessibility.** Exercise the supermarket presentation
   and single scalable UI at required resolutions/scales in all nine launch
   locales, with keyboard-only navigation, visible focus, correct wrapping/font
   fallback, independent audio channels, and presentation fallbacks.
4. **Hosted convergence.** Complete deterministic 2-, 3-, and 4-player staging
   runs through private/friends Steam invites, including owner handoff and
   successful owner-authorized game admission; verify automatic worst-latency
   Railway region selection for same-area and cross-area parties, observe stable
   fifth-seat and server-capacity refusal without affecting admitted players,
   and confirm no public lobby browser or matchmaking is exposed.
5. **Reconnect.** Disconnect or terminate the client during world, story, and
   combat; use a fresh Steam ticket to reclaim the reserved seat, acknowledge
   resync, exercise takeover, and confirm no duplicate visible action/audio
   event. Confirm discarded ballots stay discarded, a vacant leader is elected
   on eligible rejoin, and voluntary departure permits the last living player
   to finish an existing run but not start a new solo run. Deadline interleavings
   belong to automated transcripts.
6. **Drain and run loss.** Exercise a clean drain, drain deadline, forced stop,
   and isolated process crash at declared user-visible points. Keep clients
   connected: idle lobbies close with maintenance guidance, new runs cannot
   start, and an active run finishes through summary before its session closes.
   Reserved-seat rejoin works for non-ended running/summary sessions. An idle or
   completed session must not show run loss; interrupted runs use the stable
   run-lost outcome with no false recovery claim. Deadline interleavings and
   internal cleanup assertions belong to process-lifecycle tests.
7. **Trust boundaries.** Observe invalid/expired Steam proof and compatibility
   failures through supported connection flows. Wrong protocol shows
   `update_required`; a valid unequal content identity shows `content_mismatch`
   with update/restart guidance. A failed rejoin preserves the reservation and
   existing connection. Controlled mismatching peer builds have their own
   recorded digests. Oversized/malformed traffic and per-field mismatch matrices
   belong to protocol integration tests, as do direct uninvited joins and grant
   forgery/replay cases. Confirm no leaked internals or external content fetch.
8. **Final package and rollback.** On exact Phase 13 candidates, inspect package
   contents, prove development identity/debug paths and secrets are absent,
   launch both targets from clean caches, verify signatures/checksums/depot
   layout, and complete the isolated rollback drill.

Do not mutate the exact candidate to manufacture an internal failure and then
describe the result as candidate QA. Modified negative packages or faulting
environments are separate diagnostic artifacts with their own digest.

## Operations Oracles

- `/healthz` reports process liveness only.
- `/readyz` remains false until built-in content and required initialization
  complete, and becomes false before graceful drain.
- Capacity exhaustion rejects new admission but does not make healthy existing
  sessions or reserved-seat rejoin unready.
- Shutdown stops admission and run starts, ends idle lobbies and completed
  summaries, and permits reserved-seat rejoin only to non-ended running/summary
  sessions. It awaits bounded owned work without requiring clients to leave and
  distinguishes maintenance closure from an interrupted run.
- Public errors are stable and redacted; internal logs retain useful structured
  context without credentials or personal data.

## Severity And Stages

- `Blocker`: authority corruption, credential exposure, universal startup
  failure, or no safe workaround.
- `Critical`: a major supported journey is unavailable, cross-player authority
  fails, or repeatable drain/run-loss behavior violates the contract.
- `Major`/`Minor`: degraded behavior with a safe workaround or limited
  presentation impact.

Phase 09 may record a Playable-MVP report. Phase 12 records a production-profile
pre-release report. Phase 13 reruns every applicable journey on exact final
digests and is the only final Release-ready QA verdict.

Synthetic capacity tests use a separate non-distributable server build and are
not final-package QA. Release evidence pairs both builds using the artifact
matrix in `docs/benchmark-plan.md`; the production digest always uses real Steam
authentication with the authorized account pool.
