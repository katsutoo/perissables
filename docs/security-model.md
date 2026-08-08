# Security Model

Status: Normative review and verification scope
Owner: Project team
Updated: 2026-08-09

Authority: threat scope and verification policy are normative here; concrete product/protocol controls and limits derive from the named locked sections of `docs/mvp-contract.md`.

## Objectives And Assets

Protect authoritative run integrity, player/session identity, Steam proofs, rejoin tokens, unpublished packs, OAuth sessions, uploader content, and service availability. Prevent one player or pack from reading hidden state, acting as another player, changing outcomes, escaping pack storage, exhausting unbounded resources, or sending credentials to an untrusted endpoint.

## Threat Actors And Boundaries

- Unauthenticated internet clients can open connections and submit malformed handshakes.
- Authenticated players can send arbitrary, reordered, duplicated, stale, or state-invalid messages and manipulate Steam lobby metadata.
- Community packs and every contained byte are hostile even after hub validation.
- A compromised client cannot be trusted with authoritative outcomes or hidden state.
- Railway volumes/backups, CI logs/artifacts, R2 objects, OAuth callbacks, dependencies, and administrator sessions are distinct trust boundaries.
- Production, Steam, Railway, OAuth providers, and third-party services are not active-test targets without explicit written authorization for the exact environment and actions.

## Threat Register

This register is the minimum tracked set. Each implementation phase links its applicable tests or review evidence and records any changed residual risk; an unresolved High-priority threat blocks the phase that exposes it.

| Threat | Primary control | Owner | Verification criterion | Priority | Residual risk |
| --- | --- | --- | --- | --- | --- |
| Session or player takeover through forged/replayed credentials | Steam proof validation, identity/session/player/token-generation binding, keyed token digests, rotation, absolute expiry, one active connection | Server identity owner | Negative join/rejoin/takeover tests prove wrong identity, token, generation, app, and expiry cannot gain authority and leak no secret | High | Compromised Steam account or player device remains authoritative until provider/session revocation |
| Cross-player, hidden-state, or client-authority violation | Server-owned state machine, deny-by-default action ownership, per-player projection | Game-core and protocol owners | Transcript and 2/3/4-client tests prove unauthorized actions do not mutate state and projections contain no hidden fields | High | Bugs in a newly added projection/action require renewed review and tests |
| Malicious pack traversal, parser exploit, decompression/resource exhaustion, or external fetch | Portable archive subset, strict bounds, fail-closed Linux Landlock/seccomp and Windows AppContainer/Job Object workers, direct validated-member access | Pack-validation owner | Conformance/fuzz and supported-OS sandbox-negative suites cover every parser, limit, socket/file/process escape, teardown, and unavailable-control path; no panic, escape, fetch, descendant, residue, or acceptance beyond bounds | High | Kernel/platform sandbox defects and native decoder vulnerabilities remain possible inside the bounded worker |
| Persistence corruption, rollback, split ownership, or acknowledged-state loss | Exclusive process lock, pinned SQLite WAL with `synchronous=FULL`, staged transaction commit-before-publish, row checksums, bounded recovery | Persistence owner | Actual-volume qualification plus fault injection at SQLite begin/write/commit/checkpoint, process crash, disk-full, lock, row-corruption, and database-corruption paths restores exactly or fails closed without treating uncommitted state as acknowledged | High | Whole-volume or database-page loss can terminate sessions and require operator backup recovery |
| CPU, memory, socket, queue, disk, retry, or provider exhaustion | Contract-wide admission and resource bounds, backpressure, rate limits, circuit breakers, slow-consumer disconnect | Server operations owner | Boundary tests plus authorized benchmark workloads demonstrate refusal before oversubscription and stable healthy-client service | High | A single instance has finite capacity and may refuse new sessions under legitimate spikes |
| Native FFI, dependency, or build-supply-chain compromise | Pinned toolchain/SDK/source checksums, minimal adapters, unsafe contracts, protected release jobs, advisory/license/source policy | Release owner | Clean reproducible builds, dependency review, wrapper tests/sanitizers where supported, signature/provenance verification | High | Upstream compromise may evade known-advisory and reproducibility controls |
| Hub account/object authorization bypass or confused deputy | Provider-scoped identity, CSRF/session controls, server-side ownership checks, private immutable validation state | Hub identity/storage owner | Two-account and moderator-role tests cover read/edit/delete/publish/link paths and object state transitions | High at Phase 19 | OAuth/provider compromise and moderator abuse require revocation/audit response |
| UGC abuse, illegal content, privacy failure, or deletion failure | Private-by-default launch gate, moderation/reporting, quotas, rights terms, minimized retention, reference-counted deletion | Hub policy/moderation owner | Pre-public-launch policy review and QA prove reporting, takedown, ban, account/content deletion, backup-expiry disclosure, and audit access | High at Phase 20 | Provider backups expire asynchronously and cannot promise immediate physical erasure |

Owners are roles until named individuals are assigned in the phase evidence. Accepted residual risk requires an owner, rationale, target review date, and release-owner approval; it is not implied by passing tests.

## Required Controls

- Connect only to the trusted environment allowlist; lobby metadata cannot choose an endpoint. Production uses standard certificate-chain/hostname verification, TLS `1.2+` with no user bypass or plaintext fallback, and no direct origin path that bypasses the trusted proxy. Validate Steam proof for the expected app/ownership before issuing or accepting credentials.
- Bind connection authority to validated identity, session, player, and token generation. Envelope IDs are consistency checks, not authorization credentials. Rejoin rotation uses the durable pending-handoff/acknowledgement state; neither a dropped response nor an old generation may create two authoritative connections.
- Store rejoin and web-session secrets as keyed digests with rotation and absolute expiry. Never log or persist raw tickets/tokens, and redact them from errors, traces, screenshots, and benchmark artifacts.
- Enforce every contract limit before allocation or expensive parsing where possible. Bound sockets, handshakes, tasks, queues, retries, decoded data, parser depth, fan-out, storage, and session lifetime.
- Treat checksums as canonical logical-content equality, not raw archive-byte equality or authenticity. The server computes/loads its own validated pack and never trusts a client-supplied digest as authority. Tier 2 authenticity uses only the exact root-signed key-set, fresh revocation-list, and domain-separated attestation contract; stale metadata fails closed before activation.
- Accept only the locked regular-file ZIP subset; disable XML external access and all parser/network side effects; use isolated temporary storage and remove it on every outcome. Release validation starts only after every locked Linux or Windows sandbox control is installed, and unsupported/unavailable controls fail closed without a bypass setting.
- Deny authorization by default for pack edit/delete, account deletion, moderation, and administrative actions. Audit admin actions without storing sensitive payloads.
- OAuth uses provider-scoped subjects, state, PKCE where supported, exact callbacks, secure `HttpOnly`/`SameSite` cookies, CSRF protection, session rotation/revocation, and reauthentication for account linking or admin-sensitive actions.
- Uploaded objects remain private and immutable through validation. Publication follows the contract's idempotent R2/Postgres state machine with quota reservation, expiry cleanup, reconciliation, reference-counted unlisting/deletion, and no claim of a cross-service transaction. Overwriting a validated generation is forbidden. Signing keys are separated by publication/revocation usage; the offline root never signs packs or runs in the hub application.
- Persistence uses one pinned SQLite build and one supervised connection owner, preserves storage/free-space headroom for terminal work, caps all database/WAL/temporary artifacts and startup validation, persists per-seat grace/token-handoff state in the same atomic session row, and marks an invalid row quarantined without deleting its bytes. Database-level corruption keeps readiness false for operator restore; corrupt state never re-enters admission automatically.
- Render all UGC as escaped plain text under restrictive CSP and security headers. Downloads use fixed safe content types, `nosniff`, attachment disposition, and a cookie-less public origin; private object URLs are short-lived and single-object.
- Enforce the contract retention table for logs, backups, sessions, moderation records, crash uploads, and verification artifacts. Deletion requests remove live data promptly and document finite provider-backup expiry rather than promising impossible immediate erasure.

## Supply Chain And Release

- Pin Rust, Cargo.lock, CI tools, Git revisions, native libraries, and build images. Review dependency licenses and advisories; scanner output is a lead until applicability and reachability are established.
- Run local secret scanning before public release. Never authenticate with a discovered credential; stop, record only its type/location or approved fingerprint, notify the owner, and rotate through an authorized process.
- Release artifacts contain no debug credentials, private symbols unless intentionally shipped, test endpoints, developer-only character paths, or public-release uploads prohibited by policy.

## Verification And Reporting

- Derive negative tests from each boundary and run them only against local or isolated staging with synthetic accounts/data.
- Fuzzing follows `docs/test-strategy.md`; load/resource measurements follow `docs/benchmark-plan.md` and never intentionally deny service.
- Findings record location, preconditions, impact, safe evidence, technical severity, remediation priority, confidence, fix, and validation plan.
- Preserve the first failure, stop when risk is established, and do not access another party's data or use discovered credentials.
- A clean review or scan is not proof of security. Release reporting states reviewed, not applicable, blocked, and untested areas plus residual risk.
