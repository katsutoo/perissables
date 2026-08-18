# Security model

Status: Normative game/runtime threat scope
Owner: Sole developer
Updated: 2026-08-16

Product controls and compatibility behavior derive from
`docs/mvp-contract.md`. This model covers the shipped game client/server and
built-in content.

## Assets and objectives

Protect authoritative run integrity, player/session identity, Steam tickets,
hidden per-player state, local files, built-in content, service availability,
release artifacts, and operator credentials.

Prevent a client from acting as another player, deciding authoritative outcomes,
reading hidden state, replaying actions into invalid mutation, redirecting Steam
proof to an untrusted endpoint, or consuming unbounded resources.

## Threat actors and boundaries

- Unauthenticated internet clients may open connections and send malformed
  handshakes.
- Authenticated players control their clients, messages, order, timing, and
  Steam lobby metadata.
- Steamworks, `raylib`, native decoders, dependencies, CI, Railway, Steam, and
  operator credentials are separate trust boundaries.
- Production and third-party services are not active-test targets without
  explicit authorization for the exact environment and actions.

## Threat register

| Threat | Required control | Verification | Residual risk |
| --- | --- | --- | --- |
| Player/session takeover | Fresh Steam app/ownership validation; identity/session/player binding; one connection | Wrong identity/app, replay, takeover, seat reclaim, and expiry tests | A compromised Steam account/device remains authoritative |
| Client authority or hidden-state leak | Server-owned state machine; deny-by-default action checks; recipient views | 2/3/4-client transcripts and differential projection tests | New actions/projections require renewed review |
| Development identity in release | Compile-time feature separation and package inspection | Final packages prove adapter/symbol/config absence | Build misconfiguration remains possible until package checks run |
| Resource exhaustion | Bounded admission, messages, queues, tasks, retries, parsing, and fan-out | Boundary tests and authorized capacity experiments | One region may refuse legitimate spikes |
| Active-run loss | In-memory-only sessions; stable run-lost error; graceful drain | Clean drain, deadline, forced-stop, and crash tests | A process or regional outage can end active runs |
| Native/dependency/build compromise | Pinned sources/checksums, minimal unsafe adapters, protected release jobs, advisories/licenses, reproducible artifacts | Wrapper tests/sanitizers where supported, dependency review, signature/provenance checks | Upstream compromise may evade known checks |

Accepted residual risk records an owner, rationale, review date, and release-owner
approval.

## Required controls

- Connect only to trusted environment endpoints. Production uses normal
  certificate-chain and hostname verification with no bypass or plaintext
  fallback.
- Steam lobby metadata is discovery data, never endpoint or gameplay authority.
- Validate identity, authorization, phase, target, revision, sequence, rate, and
  resource bounds before state mutation.
- Keep Steam tickets only as long as validation requires. Never log or persist
  them raw.
- A fresh validated Steam identity can reclaim only its own reserved in-memory
  seat.
- Public errors expose stable codes, not parser internals, credentials, hidden
  state, or filesystem paths.
- Every queue, task set, retry loop, allocation, parser, decoded built-in
  resource, connection, and drain path has a justified bound before exposure.
- Trusted proxy configuration owns source attribution; forwarding headers from
  other peers are ignored.
- Built-in JSON/TMX parsing disables external entities, external resources, and
  network access and enforces schema/resource bounds before session admission.
- Long-lived tasks have explicit owners, cancellation, and joined outcomes.
- Structured logs identify operations without message bodies, credentials, or
  personal data.
- A draining process refuses new admission before shutdown. A forced stop never
  claims a clean drain or recovered run.

## Supply chain and release

- Pin Rust, `Cargo.lock`, CI tools/actions, Git revisions, native libraries/SDKs,
  build images, and the Steam Linux Runtime/container.
- Review dependency licenses, sources, and advisories; scanner output is a lead,
  not proof of reachability.
- Run a local history-aware secret scan before public release. Never test a
  discovered credential; rotate it through an authorized process.
- Signing keys are non-exportable and available only to protected release-tag
  workflows after artifact checksum approval.
- Final packages contain no local/benchmark identity adapter, debug credential,
  test endpoint, developer character/action, secret, or unintended private
  symbol.

## Privacy and retention

- No optional gameplay telemetry, voice recording, or voice playback ships by
  default.
- Use only gameplay-required platform/session identifiers.
- Crash upload, if added, is informed, opt-in, bounded, and redacted.
- Logs, sessions, crash data, and raw QA/benchmark evidence have documented
  finite retention before production.
- Synthetic fixtures and redacted release summaries contain no credentials,
  account identifiers, or production records.

## Verification and reporting

- Security tests use local or explicitly authorized isolated staging with
  synthetic/project-controlled data.
- Fuzzing follows `docs/test-strategy.md`; resource measurements follow
  `docs/benchmark-plan.md`.
- Stop once a risk is established. Do not access another party's data, use a
  discovered credential, load-test a third party, or suppress a credible issue
  because exploitation is unsafe.
- Findings record location, preconditions, impact, evidence, severity,
  remediation priority, confidence, fix, validation, and residual risk.
- A clean review or scan is not a claim that an unimplemented or untested system
  is secure.
