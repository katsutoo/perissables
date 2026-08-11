# QA Plan

Status: Normative release-behavior policy
Owner: Project team
Updated: 2026-08-11

QA exercises the product through a supported shipped entrypoint. Internal state
machines, parser matrices, and storage fault injection belong to automated
tests. QA observes the user-visible consequences and durable effects.

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
- synthetic account roles/data;
- exact steps, expected result, actual result, visible and durable effects;
- logs or captures with secrets/personal data redacted; and
- cleanup and untested residual risk.

Preserve the first intermittent failure. Diagnostic reruns use the same artifact
and a stated protocol; a later pass does not erase it.

## Supported Matrix

| Target | Required environment |
| --- | --- |
| Linux | Frozen Ubuntu 24.04 LTS x86_64 image; `1280x720` and `1920x1080` |
| Windows | Frozen Windows 11 x86_64 build; `1280x720` and `1920x1080` |

Record windowed/fullscreen mode, UI scale, GPU, and driver. Keyboard-only
operation and every shipped UI variant are required. Other Linux distributions,
Steam Deck, controllers, and additional display modes are exploratory until
added to the support contract.

## Required Journeys

1. **Clean install and startup.** Launch from a clean user-data directory,
   exercise settings creation/recovery through the UI, relaunch, and confirm no
   authoritative state or secret is stored locally.
2. **Three-run loop.** Complete lobby -> selection -> ready/start -> story ->
   check -> combat -> summary -> lobby three times without stale state.
3. **Presentation and accessibility.** Exercise all themes and UI variants at
   required resolutions/scales with keyboard-only navigation, visible focus,
   readable text, independent audio channels, and presentation fallbacks.
4. **Hosted convergence.** Complete deterministic 2-, 3-, and 4-player staging
   runs; observe stable fifth-seat/server-capacity refusal without affecting
   admitted players.
5. **Reconnect.** Disconnect during world, story, and combat; exercise dropped
   handoff response, acknowledged resync, takeover, token expiry, and no
   duplicate visible action/audio event.
6. **Supported recovery.** After Phase 11, restart cleanly and crash the isolated
   staging process at declared user-visible points. Rejoin and verify behavior
   matches the selected storage acknowledgement boundary. Internal storage
   begin/write/commit/checkpoint matrices remain integration tests.
7. **Trust boundaries.** Exercise invalid/expired Steam proof, wrong version or
   pack identity, oversized traffic, and valid/invalid Tier 1 import through
   normal user entrypoints. Confirm stable public errors, no leaked internals,
   no external fetch, and no residue.
8. **Final package and rollback.** On exact Phase 13 candidates, inspect package
   contents, prove development identity/debug paths and secrets are absent,
   launch both targets from clean caches, verify signatures/checksums/depot
   layout, and complete the isolated rollback drill.

Do not mutate the exact candidate to manufacture an internal failure and then
describe the result as candidate QA. Modified negative packages or faulting
environments are separate diagnostic artifacts with their own digest.

## Operations Oracles

- `/healthz` reports process liveness only.
- `/readyz` remains false until required initialization and selected storage
  recovery are complete, and becomes false before graceful drain.
- Capacity exhaustion rejects new admission but does not make healthy existing
  sessions or reserved-seat rejoin unready.
- Shutdown stops admission, preserves the selected durable boundary, awaits
  owned work within its measured deadline, and does not claim a clean close
  after forced termination.
- Public errors are stable and redacted; internal logs retain useful structured
  context without credentials or personal data.

## Severity And Stages

- `Blocker`: authority/data corruption, credential exposure, universal startup
  failure, or no safe workaround.
- `Critical`: a major supported journey is unavailable, cross-player authority
  fails, or a repeatable crash/recovery failure occurs.
- `Major`/`Minor`: degraded behavior with a safe workaround or limited
  presentation impact.

Phase 09 may record a Playable-MVP report. Phase 12 records a production-profile
pre-release report. Phase 13 reruns every applicable journey on exact final
digests and is the only final Release-ready QA verdict.
