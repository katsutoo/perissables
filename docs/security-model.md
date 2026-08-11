# Security Model

Status: Normative game/runtime threat scope
Owner: Project team
Updated: 2026-08-11

Product controls and compatibility behavior derive from
`docs/mvp-contract.md`. This model covers the game client/server and local
Tier 1 validation. A future community hub owns its own threat model.

## Assets And Objectives

Protect authoritative run integrity, player/session identity, Steam tickets,
rejoin tokens, hidden per-player state, local files, built-in content, service
availability, release artifacts, and recovery data.

Prevent a client or pack from:

- acting as another player or deciding authoritative outcomes;
- reading hidden state or credentials;
- replaying, reordering, or racing actions into invalid mutation;
- redirecting credentials to an untrusted endpoint;
- escaping pack storage/validation; or
- consuming unbounded CPU, memory, disk, sockets, tasks, queues, or retries.

## Threat Actors And Boundaries

- Unauthenticated internet clients may open connections and send malformed
  handshakes.
- Authenticated players control their clients, messages, order, timing, and
  Steam lobby metadata.
- Every community-pack byte remains hostile after prior validation.
- Native decoders, Steamworks, `raylib`, the selected database, dependencies,
  CI, Railway, Steam, and operator credentials are separate trust boundaries.
- Production and third-party services are not active-test targets without
  explicit authorization for the exact environment and actions.

## Threat Register

| Threat | Required control | Verification | Residual risk |
| --- | --- | --- | --- |
| Player/session takeover | Steam app/ownership validation; identity/session/player/token-generation binding; digest-stored rotating tokens; one connection | Wrong identity/app/token/generation, replay, dropped handoff, takeover, and expiry tests | A compromised Steam account/device remains authoritative until revocation/expiry |
| Client authority or hidden-state leak | Server-owned state machine; deny-by-default action checks; recipient views | 2/3/4-client transcripts and differential projection tests | New actions/projections require renewed review |
| Development identity in release | Compile-time feature separation and package inspection | Release builds and final packages prove adapter/symbol/config absence | Build misconfiguration remains possible until CI/package checks run |
| Malicious pack or parser escape | Strict archive/schema/path/resource bounds; disabled external resolution; killable sandboxed workers | Conformance, fuzz, and supported-OS escape/cleanup tests | Kernel/native decoder defects remain possible inside the sandbox |
| Resource exhaustion | Bounded admission, messages, queues, tasks, retries, storage, parsing, and fan-out | Boundary tests and authorized capacity experiments | One instance may refuse legitimate spikes |
| Persistence corruption or acknowledged-state loss | Phase 10 selected transactional design; explicit acknowledgement boundary; versioning; backup/rollback | Real-adapter integration faults plus clean/crash staging recovery | Whole-store/provider loss can require operator recovery |
| Native/dependency/build compromise | Pinned sources/checksums, minimal unsafe adapters, protected release jobs, advisories/licenses, reproducible artifacts | Wrapper tests/sanitizers where supported, dependency review, signature/provenance checks | Upstream compromise may evade known checks |

Owners are assigned in phase issues/evidence. Accepted residual risk records an
owner, rationale, review date, and release-owner approval.

## Required Controls

- Connect only to trusted environment endpoints. Production uses normal
  certificate-chain and hostname verification with no user bypass or plaintext
  fallback.
- Steam lobby metadata is discovery data, never endpoint or gameplay authority.
- Validate identity, authorization, phase, target, revision, sequence, rate, and
  resource bounds before state mutation.
- Store bearer credentials only as long as required. Never log or persist raw
  Steam tickets or rejoin tokens.
- Public errors expose stable codes, not parser internals, credentials, hidden
  state, or filesystem paths.
- Every queue, task set, retry loop, allocation, parser, archive, decoded
  resource, connection, and storage path has a justified bound before exposure.
- Trusted proxy configuration owns source attribution; forwarding headers from
  other peers are ignored.
- Pack validation accepts only the documented regular-file archive subset,
  rejects traversal/aliases/links/devices, disables XML external access, and
  performs native/untrusted decoding only after the required OS sandbox is
  installed.
- Validation workers inherit no secrets or broad handles, have no network or
  child-process capability, and are killed/reaped on timeout, cancellation, or
  limit breach.
- Release builds fail closed when a required supported-OS sandbox primitive is
  unavailable. Test escape fixtures cannot compile into packages.
- The selected persistence design keeps server-only state, versions formats,
  uses explicit transactions/parameterized operations, bounds recovery, and
  never reports uncommitted state as durable.
- Long-lived tasks have explicit owners, cancellation, and joined outcomes.
- Structured logs identify operations without message bodies, credentials, or
  personal data.

## Supply Chain And Release

- Pin Rust, `Cargo.lock`, CI tools/actions, Git revisions, native
  libraries/SDKs, build images, and the selected storage dependency.
- Review dependency licenses, sources, and advisories; scanner output is a lead,
  not proof of reachability.
- Run a local history-aware secret scan before public release. Never test a
  discovered credential; report its type/location and rotate through an
  authorized process.
- Signing keys are non-exportable and available only to protected release-tag
  workflows after artifact checksum approval.
- Final packages contain no local identity adapter, debug credential, test
  endpoint, developer character/action, secret, or unintended private symbol.

## Privacy And Retention

- No optional gameplay telemetry or voice recording by default.
- Use only gameplay-required platform/session identifiers.
- Crash upload, if added, is informed, opt-in, bounded, and redacted.
- Logs, backups, sessions, crash data, and raw QA/benchmark evidence have
  documented finite retention before production.
- Synthetic fixtures and redacted release summaries contain no credentials,
  account identifiers, or production records.

## Verification And Reporting

- Security tests use local or explicitly authorized isolated staging with
  synthetic/project-controlled data.
- Fuzzing follows `docs/test-strategy.md`; resource measurements follow
  `docs/benchmark-plan.md`.
- Stop once a risk is established. Do not access another party's data, use a
  discovered credential, load-test a third party, or suppress a credible issue
  because exploitation is unsafe.
- Findings record location, preconditions, impact, evidence, technical
  severity, remediation priority, confidence, fix, validation, and residual
  risk.
- A clean review or scan is not a claim that an unimplemented or untested system
  is secure.
