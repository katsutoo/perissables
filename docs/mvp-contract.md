# Les Périssables MVP Contract

Status: Locked for Phase 01 implementation; later phase changes follow the tracker and definition of done
Owner: Project team
Updated: 2026-07-24

This document is the single source of truth for locked MVP scope and implementation decisions.

Milestone terms are exact throughout the documentation: `Playable MVP` means Phases 00-13, `Release-ready` means Phases 00-18, and `Post-MVP` means work outside the Release-ready gate. Unqualified `MVP` refers to the product scope shared by those milestones, not to a delivery gate.

## Product Objective

Ship a funny, fast, multiplayer pixel-art RPG where players pick premade food characters, run story-driven adventures, and survive dice/combat events using one shared engine.

## MVP Scope (In)

1. Native desktop client (Rust + `raylib`/`raylib-rs`) and authoritative Rust server (`axum` + `tokio`).
2. Story engine driven by JSON (no story-specific hardcoded logic).
3. d100 dice system with locked rule set:
   - stat range `5..=70` (inclusive)
   - internal roll range `0..=100` (inclusive)
   - success when `roll <= stat`
   - critical success on `000` (internal `0`)
   - critical failure on `100`
   - this `000`/`100` asymmetry is intentional for flavor even though it creates `101` discrete outcomes (`000` through `100`)
   - each internal result is equiprobable; normal non-critical success probability is `stat / 101`, total success probability including critical success is `(stat + 1) / 101`, normal failure probability is `(99 - stat) / 101`, and each critical has probability `1 / 101`
4. Basic turn-based combat integrated with story encounters.
5. Preset character roster (no character build system, no leveling).
6. Run flow: lobby -> story -> end summary -> lobby.
7. Three environment themes:
   - supermarket
   - garden
   - storage_room
8. Core multiplayer for 2-4 players with server authority.
9. Event-driven audio system with channels for ambience, music, SFX, and tiny voice barks.
10. Complete keyboard-only operation and at least two behavior-identical built-in UI skin variants selected by theme data. Keyboard focus is always visible with an indicator at least `2` rendered pixels thick at both supported UI scales, and skinning may not change focus order, hit targets, controls, authority, or event behavior.

## Explicit Non-Goals (Not In MVP)

1. Skill trees, progression, talent builds, class customization.
2. Advanced inventory management or deep equipment systems.
3. Procedural map generation.
4. Voice chat.
5. Mod scripting with arbitrary code execution.
6. Story-specific custom UI logic (visual skinning only in MVP).
7. Mobile, browser, or console release targets.
8. Steam achievements; MVP Steamworks scope is ownership/authentication, lobbies/invites, and depot distribution only.

## Locked Runtime Constraints

- Max party size: `4`
- Max active sessions per server instance before admission is refused: `64`
- Max players per server instance before admission is refused: `256`
- For capacity accounting, an active session is any session from successful creation until its durable `Ended` transition plus active-state cleanup (seat release and tombstone) completes, including `Lobby`, `Running`, `Summary`, paused/restored sessions, and an `Ended` session still completing that cleanup. Retention-only snapshot/quarantine files remain under storage caps but do not consume an active-session slot. Lifecycle state and capacity release are therefore distinct until active-state cleanup finishes.
- A player is an occupied seat, whether connected or reserved during rejoin grace. A seat counts from committed admission until explicit leave in `Lobby`, session end, or grace expiry. A socket is not a player, and reconnect replacement never increments either capacity counter.
- Admission capacity is the remaining active-session and occupied-seat budget. Creation atomically reserves one session and one seat; join atomically reserves one seat. Rejoin uses an existing reservation and is allowed when new-admission capacity is zero. Concurrent operations are serialized by the admission owner and never oversubscribe a limit.
- Target rate: `60 FPS`
- Client fixed-timestep update rate: `60 Hz`
- The client processes at most `5` fixed updates in one rendered frame and caps accumulated unprocessed time at `250 ms`. Excess time is discarded once, increments `client_update_time_dropped_total`, and emits one rate-limited diagnostic; it must not cause an unbounded catch-up loop. Tests inject elapsed time at the exact cap and one update beyond it.
- Server simulation tick rate: `20 Hz`
- Client network input rate: at most one coalesced continuous-input message per server tick (`20 msg/s`). The `60 Hz` client update loop must not emit one network message per update. Discrete actions share the same bounded rate budget.
- Server state broadcast rate: `20 Hz` max. State broadcasts are latest-value coalesced. A server that falls behind may run at most one immediate catch-up tick, then records a missed-tick metric and resumes the normal schedule; it must not enter an unbounded catch-up loop.
- Hosted-server processing-latency target: from availability of the complete reassembled input frame at the server before JSON/protocol/authentication/rate/sequence validation through successful completion of the async WebSocket send for its correlated `input_result` after owner-queueing and any required durable commit, p95 under `100 ms` and p99 under `200 ms`, with `64` active sessions of `4` players under the locked workload in `docs/benchmark-plan.md`. The same-region open-loop scheduled-send-to-receive path is separately gated at p95 under `150 ms` and p99 under `300 ms`; public internet RTT outside that controlled path is reported separately.
- Tile size: `16x16`
- Primary resolutions for MVP QA: `1280x720`, `1920x1080`
- Playtest targets: median clean-account lobby-to-run start at or below `3 minutes`; median completed or failed run between `20` and `35 minutes`. Phase 17 uses at least `20` consented first-time participants and `30` completed-or-failed sessions, reports the median and bootstrap `95%` confidence interval, and does not reuse a participant's first-run onboarding measurement. Lobby timing starts when the authenticated lobby first becomes interactive and ends when the first authoritative `Running` state is rendered. Run timing starts at that transition and ends at the first authoritative `Summary`; abandoned sessions are reported separately and never silently excluded. These are balance targets, not telemetry collected by default.

## Locked Persistence Model (Release-ready v1)

This model is locked now, but implementation lands in Phase 16 and is part of the Release-ready milestone, not the Playable MVP milestone. Before Phase 16, restarts and deploys may destroy active runs and must be treated as known pre-release behavior.

- Run state is server-authoritative and persists server-side: the game server snapshots active sessions to durable storage (Railway volume) so restarts and deploys do not destroy runs.
- Session snapshots persist only `HMAC-SHA-256` digests of rejoin tokens. Each occupied seat stores its Steam-identity/session/player binding, non-secret `client_join_nonce`, current token key ID and generation, optional pending-handoff token key ID and generation, absolute token expiry, connection generation/state, `last_connection_observed_unix_ms`, and optional absolute `rejoin_grace_expires_at_unix_ms`. Compare digests in constant time. Raw bearer tokens exist only on the client and transiently on the server while generating, serializing, validating, or echoing a token; they are never logged or persisted. The HMAC key is a managed server secret outside snapshots; retain the current and previous key until all snapshots/tokens using the previous key expire.
- Clients never submit saved run state; resume always happens through the rejoin flow into a server-restored session.
- Longer-term resume (the whole party returning after rejoin tokens expire) is post-MVP; if added, re-admittance is by validated Steam identity, never by extending token lifetime.
- `PersistedSessionSnapshot` is a server-only representation capped at `256 KiB`. It includes authoritative gameplay state, logical tick, state revision, `next_event_id`, each player's committed `next_input_seq`, optional active-run RNG algorithm/version/state, pack identity, the complete per-seat token/connection/grace state above, absolute `created_at_unix_ms`/`expires_at_unix_ms`/`last_observed_unix_ms`, and remaining monotonic gameplay-deadline durations. Uncommitted unresolved-ledger entries are never snapshotted; after crash/restore clients retry from the committed frontier. Lobby snapshots contain no run RNG. It is distinct from the client-visible `ClientResyncState`.
- `save_version` starts at `1`. A normal snapshot reader supports its current version and the immediately previous positive version once version `2` exists; future, zero, and older versions fail closed without partial restoration. The sole exception is a retained rollback build prepared under the release rollback contract: it may carry an explicitly tested read-only migrator for exactly the next version while continuing to write its current version.
- After Phase 16, externally acknowledged authoritative transitions are group-committed to one bounded server WAL before their `state`, `event`, `input_result`, token handoff, or leave acknowledgement is emitted. The server batches all sessions for at most `20 ms`, appends length/checksum-framed records, `fsync`s once, then publishes acknowledgements. The WAL stops ordinary client mutations at a `56 MiB` soft limit. The final `8 MiB` is reserved for bounded lease/fencing, token-handoff completion, expiry/leave/end tombstones, snapshot-coverage, and compaction records; no prospective append may cross the `64 MiB` hard limit. Every per-seat reserve-class record is at most `4 KiB`, every per-session reserve-class record at most `64 KiB`, and one global lease/compaction control record at most `1 MiB`. Before allowing any reserve-region append, accounting preserves `1 MiB` for up to `256` seat obligations, `4 MiB` for up to `64` session obligations, `1 MiB` for global control, and `2 MiB` emergency margin; an obligation consumes its reservation rather than competing for new space. Reaching either applicable limit rejects new state-changing work until compaction succeeds. Before Phase 16, restart rollback/loss remains explicit pre-release behavior and no durability claim applies.
- Dirty snapshots compact the WAL at most once per session per second. Run start, encounter resolution, run end, initial token issuance after Phase 16, and token rotation are durable checkpoints; a token is not returned until its digest/generation record is committed.
- Snapshot generations are immutable files named by monotonically increasing generation. Write through one fixed same-directory temporary filename, `fsync` it, rename to its final generation, then `fsync` the directory. On startup, remove that stale temporary filename, inspect at most the newest `8` generation entries, and select the newest checksum-valid generation whose WAL replay validates; retain the newest two valid generations and delete older ones only after the current generation is durable. At most one snapshot write and one pending dirty state exist per session. If no generation and WAL combination validates and quarantine has reserved capacity, move the bounded inspected files to quarantine, persist an ended-session tombstone when possible, exclude the session from admission capacity, emit one redacted operator error, and make later join/rejoin return `session_not_found`. If quarantine cannot reserve its file/session/byte budget, leave the original bytes untouched, mark restoration blocked, keep `/readyz` unsuccessful, and require operator cleanup/recovery; never delete evidence to become ready. Restoration may succeed only after every damaged session is restored or durably quarantined.
- Each state-changing operation is a transaction: validate against committed state; clone or journal only the bounded fields it can change; stage gameplay state, RNG position, consumed `input_seq`, revision, and outputs; append and commit one WAL record; then atomically replace committed in-memory state and enqueue outputs. A failed append or `fsync` discards the entire staged transition, does not consume `input_seq`, does not increment `state_revision`, and emits no semantic event/state/result. It then clears every unresolved ledger in that failed session batch, resets each affected admission frontier to its lowest discarded sequence (or its committed frontier when empty), and returns correlated recoverable `persistence_unavailable` with that retryable `next_input_seq`. All higher discarded inputs are retried by the client from the returned frontier; none remains ambiguously admitted. No observer may read staged state.
- Persistence retries are bounded to three attempts with delays of `100 ms`, `500 ms`, and `2 s`. While retrying, the owning session accepts heartbeats/resync but rejects new state-changing work with correlated `persistence_unavailable`. After the third failure it enters `PersistenceBlocked` for at most `30s`; one successful probe and commit restores service. Otherwise the owner commits no further mutation, sends a `session_terminated` notice if the writer is available, closes with `1011`, and retains the last valid durable generation for operator recovery.
- Every immutable snapshot generation records the inclusive WAL record offset it contains. Global compaction may replace the WAL only after every non-ended session has a checksum-valid, directory-synced snapshot covering a common offset; the low-watermark is the minimum covered offset across those sessions. Copy records after that offset into a new same-directory WAL, `fsync`, rename, and `fsync` the directory before deleting the old WAL. Ended-session tombstones remain until covered by the same low-watermark.
- Each WAL group is one bounded length/checksum frame with header magic/version/fencing generation/payload length, payload, and a fixed footer magic plus checksum over header and payload. Length is validated against the reserve/ordinary record cap before allocation. A trailing frame missing its complete footer is unacknowledged crash residue and is preserved in diagnostics then truncated to the last complete frame. A complete frame with an invalid checksum or fence is corruption, never truncation.
- On complete-frame corruption at offset `C`, preserve the original global WAL in accounted quarantine. A session may restore only from a checksum-valid snapshot whose covered offset is at or beyond `C` and whose replay beginning after that snapshot validates to the WAL end; every session that would need the corrupt frame is quarantined under the session rule. Checkpoint all safely restored sessions under the new fence, start an empty WAL, and become ready only after the replacement WAL/snapshots/tombstones are directory-synced. If global/session quarantine cannot reserve capacity, leave original bytes untouched and remain unready.
- When no non-ended sessions exist, the compaction low-watermark is defined as the current WAL end. After a directory-synced terminal index covers all retained tombstones, compaction may replace the WAL with an empty file and then apply retention cleanup; the minimum of an empty set is never evaluated.
- All server persistence artifacts, including active/ended generations, WALs, quarantine, and temporary files, have a process-wide `1 GiB` accounted cap and `4,096`-file cap. Quarantine is additionally capped at `64` sessions and `256 MiB`. Session creation reserves `1 MiB` of that budget before admitting the first seat. Ended data remains subject to the retention policy but never bypasses these caps; creation fails with `server_full`, and mutation fails with `persistence_unavailable`, before crossing a file or byte limit. Cleanup never evicts an active or quarantined session implicitly.
- Client save paths hold only non-authoritative local data (settings, keybinds, local preferences):
  - Linux: `~/.local/share/les-perissables/`
  - Windows: `%AppData%/LesPerissables/`

### Authoritative RNG v1

- `rng_version: 1` is the IETF ChaCha20 stream construction with a fresh per-run `256`-bit key, `96`-bit nonce, `32`-bit block counter, and little-endian `u32` words. For each accepted `start_run`, the session owner obtains a new key and independently generated nonce from the operating-system CSPRNG and supplies them to headless `game_core` as explicit `StartRunEntropy`; the key/nonce and `Lobby -> Running` transition are staged and committed atomically before any output. Neither value is supplied by a client or pack, and neither is exposed in client projections or logs. Returning to `Lobby` destroys the completed run's RNG state; a later run never resets or reuses that stream.
- Serialized state is exactly `rng_version`, key, nonce, block counter, and next word index `0..=15`. Restore resumes at that word without regeneration or skipping. Changing the algorithm, word interpretation, or consumption rules requires a new `rng_version` and `game_rules_version`.
- To map one `u32` word to `0..=100`, let `limit = floor(2^32 / 101) * 101`. Draw words until `word < limit`, then return `word % 101`. Rejected words advance the stream. Direct modulo without rejection is forbidden.
- Randomness is consumed only by a fully validated staged transition at the exact rule step that requests it. Malformed, unauthorized, stale-revision, unaffordable, out-of-turn, superseded, and otherwise rejected actions consume no RNG. If persistence fails, the staged RNG position is discarded with the transition. Within one committed transition, story checks consume first in authored effect order, then combat effects in stable actor-ID order.
- Tests use explicit serialized RNG states and known-answer ChaCha20 vectors. Production seeds never appear in fixtures, diagnostics, crash uploads, resync state, or benchmark artifacts.

## Locked Architecture Decisions

- HTTP/router stack: `axum` + `tower`
- Async runtime: `tokio`
- WebSocket transport: `axum` WebSockets
- Production transport: `wss://` (TLS terminated by reverse proxy)
- Production TLS uses platform-standard certificate-chain and hostname verification with no user bypass or plaintext fallback, minimum TLS `1.2`, and preference for TLS `1.3`. Railway origin ingress accepts only the configured proxy/platform path; direct origin bypass is not a supported credential-bearing route.
- Local development transport: `ws://localhost`
- WebSocket payload for MVP: JSON
- Production game-server hosting: Railway service running `crates/server`, with separate `staging` and `production` environments
- MVP server scaling policy: one active game-server instance per environment; do not enable multiple replicas until session state is externalized or sticky session allocation is implemented
- Single-instance enforcement: acquire an exclusive OS file lock on the mounted volume plus a persisted monotonically increasing fencing generation before restoration or admission. The lease records owner UUID, state (`held` or `released`), and renewal time; a held lease renews every `5s` and expires after `20s`. Crash recovery requires both expiry and successful file-lock acquisition. A clean shutdown, after final checkpoint/task completion, writes and directory-syncs `released` under the current fence before releasing the OS lock; a successor that acquires the lock may then increment the generation immediately without waiting for expiry. Every WAL/snapshot write includes and verifies the current fencing generation. `/readyz` remains unsuccessful until the lease is held and restoration is complete.
- Game-server config: clients receive an exact environment-specific production WebSocket origin allowlist from trusted build/environment config. Steam lobby metadata must not introduce an arbitrary endpoint.
- Session discovery: Steam lobbies/invites are discovery only; the authoritative server owns `session_id`, player IDs, game state, dice, combat, and story progression
- Steam lobby mapping: lobby metadata stores `session_id`, aggregate `pack_id`, pack `version`, pack checksum, `protocol_version`, `content_schema_version`, and `game_rules_version`. The endpoint comes only from the trusted allowlist.
- Production identity: production `join`/`rejoin` requires Steam auth/session-ticket validation before the server issues or accepts player/session credentials
- Operations baseline: expose `/healthz` for process liveness and `/readyz` for service readiness. Readiness requires the exclusive lease, completed restoration, writable persistence, and ability to serve existing sessions; exhausted session/player admission capacity does not make the service unready. New creation/join receives `server_full`, while reserved-seat rejoin remains available. Use structured logs and monitor deploy status, missed ticks, queue saturation, error rate, disconnect rate, reconnect failures, and admission capacity separately.
- Graceful shutdown: mark unready, stop new creation/join, send the `server_restarting` notice with `retry_after_ms`, checkpoint dirty sessions, await owned tasks, persist the clean lease release, and exit within `20s`. Connections close with WebSocket `1012` only after their notice is handed to the reserved writer-control slot. If the deadline is reached, do not write `released`; exit non-zero and rely on the last durable checkpoint plus crash lease expiry.
- Map format for MVP: TMX
- Release automation: GoReleaser for build/package generation only; it does not change licensing or grant public binary distribution rights
- Distribution policy: production desktop binaries ship through Steam depots, not public release pages
- Release targets: Linux + Windows only
- macOS policy: deferred (requires Apple Developer Program for signing/notarization workflow)
- Initial workspace members are exactly `client`, `server`, `game_core`, `shared`, and `integration_tests`, with package names `les-perissables-client`, `les-perissables-server`, `les-perissables-game-core`, `les-perissables-shared`, and `les-perissables-integration-tests`. `shared` owns protocol DTOs and IDs; `game_core` depends on `shared` and the pinned MIT pack crate; `client` and `server` depend on `shared` and `game_core`; `integration_tests` may depend on all workspace crates. Reverse dependencies and cycles are forbidden.
- Protocol versioning: `protocol_version` integer, start at `1`; every increment is breaking unless a future compatibility document says otherwise.
- Content schema versioning: `content_schema_version` integer, start at `1`; every increment is breaking. Reject unsupported values for story, character, theme, and pack manifests.
- Save schema versioning: `save_version` integer for server-side session snapshots, support the current version and the immediately previous positive version once one exists, with explicit migrators plus the narrow prepared-rollback reader exception above
- Client-settings versioning: `client_settings_version` integer starts at `1`; support the current version and the immediately previous positive version once one exists, with explicit pure migrators
- Game-rule versioning: `game_rules_version` integer, start at `1`; it changes when authoritative mechanics or built-in effect semantics become network-incompatible.

## Locked Protocol And Limits

### WebSocket Envelope v1

All messages must use:

The JSON snippets in this protocol section use human-readable placeholder IDs, tickets, tokens, checksums, and timestamps to show field shape. They are intentionally not parser/conformance fixtures; versioned fixtures use fully valid generated values and exact expected bytes.

```json
{
  "type": "input",
  "protocol_version": 1,
  "session_id": "ses_...",
  "player_id": "ply_...",
  "seq": 1,
  "payload": {}
}
```

Required fields are mandatory. Unknown fields and message types, missing fields, stale/duplicate sequence numbers, and unsupported protocol versions are rejected.

`session_id` and `player_id` are server-generated opaque identifiers with at least `128` bits of CSPRNG entropy, encoded as lowercase prefixes plus unpadded base64url, and capped at `64` ASCII bytes. Clients never select or derive them.

Bootstrap exception: a first `join` is sent before the server has issued player credentials. Its envelope uses empty `session_id` and `player_id`; its payload carries either create intent or the requested existing session. A pre-auth `error` response also uses empty IDs and correlates with the request envelope `seq`. A successful response issues the real identifiers, which are mandatory on every later message, including `rejoin`.

Allowed `type` values for v1:

- `join`
- `join_response`
- `input`
- `state`
- `event`
- `notice`
- `error`
- `ping`
- `pong`
- `rejoin`
- `rejoin_response`
- `resync_request`
- `resync_state`
- `resync_ack`
- `input_result`

Message directions for v1:

- Client -> server: `join`, `input`, `pong`, `rejoin`, `resync_request`, `resync_ack`
- Server -> client: `join_response`, `rejoin_response`, `state`, `event`, `notice`, `error`, `ping`, `resync_state`, `input_result`

Sequence rule:

- `seq` is an unsigned `64`-bit per-connection, per-direction transport counter and starts at `1`. It provides ordering only within one connection; it is not gameplay replay protection. A connection closes before incrementing past `u64::MAX`; wrapping is forbidden.

Join payload v1:

```json
{
  "mode": "join",
  "client_join_nonce": "AAAAAAAAAAAAAAAAAAAAAA",
  "requested_session_id": "ses_...",
  "steam_auth_ticket": "base64_ticket",
  "pack_id": "cleanup-pack",
  "pack_version": "1.0.0",
  "pack_checksum": "sha256:...",
  "content_schema_version": 1,
  "game_rules_version": 1
}
```

`mode` is `create` or `join`. `client_join_nonce` is generated once per user attempt and makes admission response recovery idempotent. For `create`, `requested_session_id` must be absent and successful admission creates a `Lobby`. For `join`, it is required and comes from Steam lobby metadata; a new seat may join only a non-ended `Lobby`. `Running` or `Summary` returns `session_not_joinable`, while `Ended` is indistinguishable from an unknown session and returns `session_not_found`. A `Lobby` whose max-party-size seats are all occupied returns `session_full`; that per-session limit is distinct from `server_full`, which reports exhausted instance-wide admission capacity. Both are checked after the session is resolved and before any seat is reserved. An identity that already owns a reserved seat with a different nonce receives `already_joined` and never receives a second seat. The same validated identity may repeat the exact mode/requested-session/nonce tuple after a lost `join_response`: the owner returns the same IDs, atomically replaces any old connection, commits a fresh initial-token digest/generation, and sends a new `join_response`. A mismatched tuple is never treated as recovery. The Steam ticket is required in production, sent only over WSS to an allowlisted origin, validated for the expected app and ownership before admission, and never logged or persisted. Pack and version fields are checked against the server's own validated aggregate pack; client values never establish authority.

Before rate-limit state allocation, authentication concurrency acquisition, or provider work, join fields obey these lexical bounds: `client_join_nonce` is unpadded RFC 4648 base64url ASCII decoding to exactly `16` bytes and is capped at `22` wire bytes; `steam_auth_ticket` decodes to `1..=4,096` bytes under the same alphabet; `requested_session_id` uses the server-ID grammar and cap; `pack_id` uses `ContentId`; `pack_version` is `1..=64` printable ASCII bytes and must parse as the locked SemVer subset; `pack_checksum` is exactly `sha256:` plus `64` lowercase hex digits; schema/rules versions are JSON integers `1..=u32::MAX`. Invalid base64, padding, non-ASCII, decoded oversize, or lexical mismatch rejects before Steam validation.

Join response payload v1:

```json
{
  "rejoin_token": "opaque_128_bit_minimum_random_token",
  "rejoin_grace_seconds": 600,
  "server_time_ms": 0,
  "state_revision": 1,
  "next_input_seq": 1
}
```

The `join_response` envelope carries the issued `session_id` and `player_id`. `server_time_ms` is milliseconds since Unix epoch sampled when the response is created and is presentation-only; authoritative deadlines use server monotonic time and are sent as non-negative `remaining_ms`, so wall-clock adjustment cannot change gameplay. The `rejoin_token` is opaque, server-generated with a CSPRNG at `128` bits of entropy minimum, never stored in Steam lobby metadata, logged, or persisted raw. It stays valid for the active session plus the `rejoin_grace_seconds` window after disconnect and is invalidated by recovery reissuance, explicit leave, session end, absolute expiry, or completion of a successful rejoin handoff. Initial issuance/reissuance is committed before response publication after Phase 16.

Rejoin payload v1:

```json
{
  "rejoin_token": "opaque_128_bit_minimum_random_token",
  "steam_auth_ticket": "base64_ticket"
}
```

For `rejoin`, `session_id` and `player_id` must be non-empty. The server binds the validated Steam identity, session, player, token generation, and a monotonically increasing connection generation. One player may have only one active connection; accepted rejoin atomically replaces and closes the old connection. A close/disconnect notification mutates seat/grace state only when its connection generation still matches the bound generation; stale takeover-close notifications are ignored. Rejoin is valid for a reserved seat in `Lobby`, `Running`, or `Summary`, including when new-admission capacity is zero.

`rejoin_token` is unpadded base64url ASCII representing `16..=32` decoded bytes and is capped at `64` wire bytes. Rejoin `steam_auth_ticket` uses the same rules as join. Token syntax/length, envelope IDs, and ticket decoding are checked before keyed-digest comparison or Steam validation.

Rejoin response payload v1:

```json
{
  "rejoin_token": "new_opaque_token",
  "rejoin_grace_seconds": 600,
  "server_time_ms": 0,
  "state_revision": 1,
  "next_input_seq": 1
}
```

Rejoin rotation is an acknowledged durable handoff rather than a one-message replacement:

- After validating a current token, the owner generates a new token, stages only its digest/key ID/new generation as `pending_handoff`, commits that record, binds the replacement connection, then sends `rejoin_response` followed by `resync_state`. The raw new token appears only in that response.
- The old current token remains valid only for restarting this handoff with the same validated Steam identity; it grants no gameplay input on the replacement connection. If `rejoin_response` was never received, the client retries with that old token, which invalidates the unknown prior pending digest and creates a new pending generation. If the response was received but the connection dropped before `resync_ack`, the client may present the known pending token to resume the same handoff; the server may echo that transiently presented raw token in `rejoin_response` without persisting it.
- A matching `resync_ack` commits promotion of the pending digest to current and invalidation of the old generation before any new `input` is accepted. Failure to commit leaves the sequence and both token generations unchanged and returns `persistence_unavailable`; the client retries `rejoin` with the current or pending token available for its observed handoff stage.
- Snapshots and WAL records preserve current/pending generations and bindings, so a crash cannot silently promote, lose, or extend a handoff. Explicit leave, session end, and absolute expiry invalidate both generations.

Other payloads v1:

| Type | Required payload fields |
| --- | --- |
| `input` | `input_seq: u64`, `based_on_revision: u64`, tagged `action` |
| `input_result` | `input_seq: u64`, `outcome`, `reason_code`, `state_revision: u64`, `next_input_seq: u64` |
| `state` | `state_revision: u64`, `logical_tick: u64`, `view: ClientView` |
| `event` | `state_revision: u64`, `events: EventItem[1..=64]` |
| `notice` | tagged payload: `server_restarting { retry_after_ms: u32 }` or `session_terminated { reason_code }` |
| `error` | `code`, `recoverable`, optional `request_seq`, `input_seq`, `expected_revision`, `next_input_seq`, `retry_after_ms`, and public `message` |
| `ping` | `nonce: u64`, `server_time_ms: u64` |
| `pong` | matching `nonce: u64` |
| `resync_request` | `last_state_revision: u64` |
| `resync_state` | `state_revision: u64`, `logical_tick: u64`, `next_input_seq: u64`, `view: ClientView` |
| `resync_ack` | matching `state_revision: u64` |

`outcome` is `applied`, `superseded`, or `rejected`. `reason_code` is `null` for applied input and otherwise one of `superseded`, `stale_revision`, `invalid_phase`, `invalid_action`, `not_owner`, `not_ready`, `duplicate_character`, `target_unavailable`, `item_unavailable`, `inventory_full`, `unaffordable`, `out_of_turn`, or `limit_reached`. A semantic rejection uses only `input_result`; `error` is reserved for envelope/auth/rate/sequence/service failures. When an error concerns an input, `input_seq` and `next_input_seq` are mandatory. `request_seq` is the inbound envelope sequence being answered.

`notice.kind` is `server_restarting` or `session_terminated`. `retry_after_ms` is mandatory only for restart and forbidden for termination. For termination, `reason_code` is mandatory and is `persistence_failed`, `clock_invalid`, `expired`, or `internal`; it is forbidden for restart. A notice is an out-of-band server control message, does not consume an input sequence or increment revision, and is ordered after all already-committed output bundles. Authenticated notices use the connection's session/player IDs. A pre-auth connection receives no notice: restart closes it with `1012`; unrecoverable shutdown closes it with `1011`. An authenticated restart sends `server_restarting` then closes `1012`; unrecoverable session termination sends `session_terminated` with a public reason then closes `1011`.

Malformed JSON, unknown fields/types, wrong direction, or authorization mismatch returns a stable error when safe and closes with WebSocket code `1008`. Oversized messages close with `1009`; service restart uses `1012`; temporary global overload uses `1013`; unrecoverable session failure uses `1011`. Public errors never include parser internals, tickets, tokens, digests, filesystem paths, or hidden state.

Input action v1 tags and fields:

| Tag | Fields | Valid phase | Queue class |
| --- | --- | --- | --- |
| `select_story` | `story_id: ContentId` | `Lobby`, owner only | discrete |
| `ready` | `character_id: ContentId` | `Lobby` | discrete |
| `unready` | none | `Lobby` | discrete |
| `start_run` | none | `Lobby`, owner only | discrete |
| `leave` | none | `Lobby` | discrete/terminal |
| `move` | `direction` (`up`, `down`, `left`, `right`, `none`) | `Running/world` | continuous/latest-value |
| `interact` | `target_id: EntityId` | `Running/world` | discrete |
| `loot` | `corpse_actor_id: EntityId`, `corpse_slot: u8` | `Running/world` | discrete |
| `story_vote` | `choice_id: ContentId` | `Running/story` | discrete |
| `combat_attack` | `target_actor_id: EntityId` | `Running/combat` | discrete |
| `combat_spell` | `spell_id: ContentId`, `target_actor_id: EntityId` | `Running/combat` | discrete |
| `combat_item` | `inventory_slot: u8`, optional `target_actor_id: EntityId` | `Running/combat` | discrete |
| `combat_pass` | none | `Running/combat` | discrete |
| `summary_ack` | none | `Summary` | discrete |

All tagged unions use a required `kind` field. `ContentId` matches `[a-z0-9][a-z0-9_-]{0,63}`. `EntityId` is a server-issued opaque identifier with the same entropy/length rules as player IDs. Human text is valid UTF-8, contains no control characters other than newline, and is capped at `2,000` Unicode scalar values. Arrays have the bounds below; unknown fields and duplicate JSON keys are rejected.

Client projection DTOs v1:

| DTO | Required fields and bounds |
| --- | --- |
| `ClientView` | `phase`, `self: PlayerView`, `party: PlayerView[1..=4]`, `available_actions: ActionKind[0..=16]`, optional recipient-only `self_inventory: InventorySlotView[0..=4]`, and exactly one phase detail |
| `PlayerView` | `player_id`, `display_name` (`1..=64` scalars), `connected: bool`, optional `character_id`; `ready: bool` is mandatory only in lobby, while `alive`, `hp: u16`, and `max_hp: u16` are mandatory only after a character is locked for a run |
| `LobbyView` | optional `selected_story: StoryPreview`, `stories: StoryPreview[1..=64]`, `owner_player_id`, `minimum_players: 2` |
| `StoryPreview` | `story_id`, `title` (`1..=128` scalars), `estimated_minutes: u8`, `theme_id` |
| `WorldView` | `map_id`, `self_position: Position`, `actors: VisibleActor[0..=128]`, `active_interaction_ids: EntityId[0..=32]`, `corpses: CorpseLootView[0..=4]` |
| `Position` | `x_subtiles: i32`, `y_subtiles: i32`, `facing`; one tile is exactly `256` subtiles |
| `VisibleActor` | `actor_id`, `position`, `sprite_id`, `alive: bool` |
| `InventorySlotView` | `slot: u8` in `0..=3`, `item_id: ContentId`; slots are unique and sorted ascending |
| `CorpseLootView` | `corpse_actor_id`, `position`, `exposed_items: InventorySlotView[0..=2]`; only server-visible, reachable loot is included |
| `StoryView` | `node_id`, `public_flags: ContentId[0..=256]`, `choices: ChoiceView[0..=8]`, `remaining_ms: u32` |
| `ChoiceView` | `choice_id`, `label` (`1..=256` scalars), `votes: u8` |
| `CombatView` | `actors: CombatActorView[1..=20]`, `turn_actor_id`, `round: u16`, `available_target_ids: EntityId[0..=20]`, `remaining_ms: u32` |
| `CombatActorView` | `actor_id`, `display_name` (`1..=64` scalars), `side`, `alive`, `hp`, `max_hp`, `resource: u8`, `resource_max: u8` |
| `SummaryView` | `result` (`completed` or `failed`), `deaths: EntityId[0..=4]`, `key_event_ids: EventId[0..=32]`, `acked_player_ids: PlayerId[0..=4]`, `discarded_rewards: RewardView[0..=16]`, `remaining_ms: u32` |
| `RewardView` | `item_id: ContentId`, `quantity: u8` in `1..=4`; entries are sorted by `item_id` |

`ClientView` contains exactly one detail field named for its `phase`: `lobby`, `world`, `story`, `combat`, or `summary`. `Lobby` and `Summary` session states map directly, while the `Running` state exposes one of the three running phases. `self_inventory` is mandatory from the accepted `start_run` through `Summary` for the recipient's own character and absent in `Lobby`; inventories of other players are server-only. Corpse exposure contains only the two eligible slots that the requesting player may currently loot and therefore may differ by recipient. `side` is `party` or `enemy`; `facing` is `up`, `down`, `left`, or `right`. Numeric fields use JSON integers and reject fractions. Optional means absent or a value, never explicit `null`, except `input_result.reason_code`.

`StoryView.node_id` selects immutable dialogue speaker/text or check label from the already checksum-matched local pack; these strings are not resent over the wire. The server still sends bounded choice labels and authoritative vote counts because availability may depend on current public flags. Failure to resolve the local node under the matched checksum is `pack_mismatch` and blocks input rather than displaying different content.

Protocol leaf types v1:

| Type | Wire definition |
| --- | --- |
| `ActionKind` | one of the exact input action tags in the action table |
| `PlayerId`, `EntityId` | opaque prefixed unpadded-base64url ASCII ID, `128` bits entropy minimum, max `64` bytes |
| `EventId` | JSON integer `1..=u64::MAX`, monotonic per session, no wrap |
| `ContentId` | lowercase ASCII matching `[a-z0-9][a-z0-9_-]{0,63}` |
| HP/resource fields | JSON integers in the gameplay bounds; current value never exceeds its declared maximum |
| `cue_id`, `sprite_id` | `ContentId`; must exist in the validated aggregate/built-in catalog |
| `channel` | `ambience`, `music`, `sfx`, or `voice` |
| dice `roll`, `stat`, `result` | integer `0..=100`, integer `5..=70`, and one of the four locked check outcomes |
| transition phase | `lobby`, `world`, `story`, `combat`, or `summary` |
| dropped ticks | JSON integer `1..=u32::MAX` |
| public reason code | one of the exact stable error, input-result, or notice reason codes, max `64` ASCII bytes |

`EventItem` has `event_id: EventId`, `kind`, and a kind-specific payload. Kinds are `audio_cue { cue_id, channel, optional actor_id }`, `dice_result { actor_id, roll, stat, result }`, `story_transition { from_node_id, to_node_id }`, `combat_transition { from_phase, to_phase }`, `actor_death { actor_id }`, `run_result { result, optional reason_code }`, and `time_dropped { ticks }`. `run_result.reason_code` is absent for normal completion/wipe and is `limit_reached` for a bounded-rules termination. One committed session batch emits at most one `event` envelope per connection containing at most `64` items and at most `48 KiB`; repeated low-value audio cues are coalesced by `(cue_id, channel, actor_id-or-absent)`. Content validation proves that every authored transition chain leaves an eight-item engine-event reserve and that every reachable committed batch fits the event count/byte caps. An unexpected required-event overflow rejects and rolls back every selected input in that staged batch with `limit_reached`; each consumed sequence receives a result and no semantic event/state is emitted.

Stable error codes v1 are `invalid_message`, `unauthorized`, `version_mismatch`, `pack_mismatch`, `session_not_found`, `session_not_joinable`, `already_joined`, `session_full`, `server_full`, `rate_limited`, `stale_input`, `input_gap`, `invalid_action`, `busy`, `persistence_unavailable`, and `internal`. That list is exhaustive for wire errors. `pack_unavailable` and `validation_resource_limit` are client-local codes reported by pack activation and pack validation respectively; they are never carried in a protocol envelope and are not public reason codes.

Revision, replay, and output rules:

- `state_revision` starts at `1` for committed session creation and increments exactly once for each committed session batch that changes gameplay/lobby state or consumes one or more expected input sequences. A batch that only answers ping/resync does not increment it. It never wraps.
- Every `input` carries `input_seq` and `based_on_revision`. `input_seq` starts at `1`, is monotonic per player/session, survives reconnect, is persisted, and never wraps. Each player has one ordered unresolved-input ledger capped at `8` accepted entries total, including pending discrete input, the latest continuous input, and superseded entries awaiting a result. `resync_state.next_input_seq` is the next value not yet admitted to that ledger.
- `based_on_revision` must not be in the future and may trail the committed revision by at most `64`. Such an input is validated against current committed state, never historical state. A future or older revision is committed as `rejected/stale_revision`: it consumes that sequence, increments revision, and emits no gameplay event. The result tells the client to request resync. Duplicate/lower values return `stale_input` without consuming; gaps return `input_gap` followed by `resync_state`; the server never guesses missing intent.
- A syntactically valid, authenticated, rate-admitted expected input consumes its sequence only when its semantic result or supersession result commits. Parse, authentication, rate, gap, duplicate, ledger-full, mailbox-full, and persistence failures do not consume it. If the ledger or owner mailbox cannot reserve space before admission, return correlated recoverable `busy` with the unchanged `next_input_seq`; never create another result obligation.
- Ledger entries remain in ascending `input_seq`. Discrete entries are FIFO. At most one continuous entry remains selectable: a newer admitted continuous entry marks the older pending continuous entry `superseded`, but both occupy ledger slots until their ordered results commit. If an older continuous entry precedes the first discrete entry selected for a tick, resolve that continuous entry as `superseded` in the same batch before the discrete result. Thus every batch resolves one contiguous ledger prefix, never emits a result for a higher sequence while a lower sequence remains unresolved, and returns at most `8` results per player.
- On each simulation tick, the owner forms one session-wide batch containing at most one selected gameplay action per connected player plus the superseded prefix entries required above. It stages players in lexical `player_id` order against the preceding staged result and encodes bounded outputs before commit. It must reserve the exact required writer slots/bytes for every recipient before appending the WAL; failure disconnects only the slow recipient and retries staging without consuming any sequence. It then commits all gameplay/RNG changes, consumed sequences, one revision, and all outputs in one WAL record. Any persistence failure releases reservations and rolls back the whole batch. Semantic rejection of one input does not roll back another valid input; the event-bound failure rule above is the sole whole-batch semantic rollback.
- For one committed batch, the owner assigns event IDs, then hands each writer an indivisible ordered semantic bundle: optional single `event` envelope, optional single `state`, then that player's `0..=8` individual `input_result` envelopes in ascending `input_seq`. Each result envelope is capped at `512` bytes and all results carry the batch revision. The writer assigns transport `seq` in that order. No later revision may be handed to any writer before the earlier bundle. Resync state supersedes queued coalescible state but never a result or semantic event.

Session lifecycle v1:

- States are `Lobby`, `Running`, `Summary`, and `Ended`; only the session owner task mutates them. The creating player is owner until its committed `leave` or grace expiry in `Lobby`; ownership then transfers to the lowest lexical occupied `player_id`, even if that reserved seat is temporarily disconnected. Ownership is not consulted during `Running` or `Summary`. If the owner seat expires there, election is deferred until the transition back to `Lobby` and occurs before that lobby projection is emitted.
- One validated Steam identity may occupy one active seat and own one lobby. Per identity, session creation is limited to `2` attempts/minute (burst `2`) and join/rejoin to `10` attempts/minute (burst `10`).
- In `Lobby`, the owner selects one available story. Changing it clears every ready selection. Each connected player may ready one unique character or unready. Ready commands collected for the same simulation tick are resolved in lexical `player_id` order, so the lowest ID wins a duplicate pick; later picks reject with `duplicate_character`.
- `start_run` applies only when a story is selected, there are `2..=4` occupied seats, every occupied seat is connected and ready with a unique character, and no token handoff is pending. Only the owner may start. It commits `Lobby -> Running/story` at the selected story's `start_node`, captures the complete occupied roster, clears lobby-only acknowledgements, and initializes a fresh per-run RNG from the explicit `StartRunEntropy` input at word zero. A `return` node enters `Running/world`; later accepted TMX trigger interactions re-enter `Running/story` at their mapped node.
- `leave` is admitted only when that player's unresolved ledger is empty; otherwise it returns `busy` without sequence consumption. Admission marks the connection terminal and rejects later input. In a lobby batch, accepted leaves stage first in lexical `player_id` order before story/ready/start actions, so a concurrently leaving seat cannot be captured by `start_run`. Accepted `leave` invalidates both token generations, releases the seat, and transfers ownership before returning its terminal applied `input_result`; the writer then closes that connection with `1000`. Other players receive the resulting state. If the last seat leaves, the same transaction commits `Ended` and its durable tombstone instead of a lobby state. Socket close without accepted `leave` is always a recoverable disconnect and retains the seat through grace.
- Run completion or party wipe commits `Running -> Summary`. Connected players may `summary_ack`; after every connected player acknowledges or `60s` monotonic summary time elapses, the owner commits `Summary -> Lobby`, clears selected story, all ready/character selections, votes, run state, deadlines, event-dedup presentation state, and run RNG state, while preserving remaining seats, input sequences, and allowed local client settings. It elects the lowest lexical occupied player when the prior owner no longer has a seat. This supports repeated runs without reconnecting.
- A session has an absolute lifetime of `4 hours`, including reconnect grace. Persist `created_at_unix_ms`, `expires_at_unix_ms`, and `last_observed_unix_ms` from the authenticated server wall clock; use monotonic elapsed time between observations in one process. On restore, wall time at or after expiry ends the session and counts downtime. If wall time is earlier than persisted `last_observed_unix_ms` by more than `1s`, fail closed by ending the session with `clock_invalid` rather than extending credentials; smaller rollback clamps to the persisted value. Forward jumps may expire sessions early and are never undone. Accepted lobby actions, accepted gameplay input, rejoin, and resync acknowledgement count as activity. A connected lobby/run/summary with no such activity for `15 minutes` ends. A zero-connected lobby with one or more reserved seats expires after `10 minutes`; a zero-seat lobby is never retained and transitions directly to `Ended`.
- On observed socket loss, the owner resolves the disconnected player's admitted unresolved ledger prefix as `superseded`, durably sets the seat disconnected, and records `rejoin_grace_expires_at_unix_ms = min(now + 600s, token_expiry)` in one bounded transition before any later seat/lifecycle mutation. Those results need not be retained after the revision/next sequence are available through resync. While connected, `last_connection_observed_unix_ms` is refreshed in the bounded once-per-second session checkpoint. On process restore, every formerly connected seat becomes disconnected with grace expiry `min(last_connection_observed_unix_ms + 600s, token_expiry)`; downtime therefore consumes grace and a crash never grants a fresh window.
- A running or summary session with no connected players pauses simulation ticks, combat/story/summary deadlines, and idle time. It expires at the earlier of the last reserved seat's absolute `rejoin_grace_expires_at_unix_ms` or absolute session expiry. Rejoin resumes each gameplay deadline with its stored non-negative monotonic duration. Absolute session and per-seat grace expiry never pause and downtime counts against both.
- Token absolute expiry is session creation time plus `4 hours`; disconnect never extends it. Acknowledged rejoin handoff rotates the current generation but retains that cap. Expiry commits `Ended`, invalidates current/pending tokens, releases capacity, and removes persisted state through a durable tombstone after Phase 16.
- A disconnected seat remains reserved until grace expiry. In `Lobby`, expiry removes the seat, invalidates tokens, transfers ownership, and ends an empty session. In `Running`, expiry releases capacity and marks that player's actor abandoned/dead in one committed transition; it emits the normal death/run-result events, participates in wipe resolution, and leaves eligible corpse loot. In `Summary`, expiry releases capacity and removes that player from the required acknowledgement set. Abandoned player records are removed before the next `Lobby` projection; if no seat remains, transition to `Ended` instead.
- Every player-controlled combat turn has a `30-second` monotonic deadline and resolves to `pass` at expiry whether the player is connected or disconnected; the all-disconnected pause rule still applies. Story votes close after `60s`. Each eligible player has one current vote, and each accepted `story_vote` atomically creates or replaces it until closure. The choice with the most votes wins. Ties use the vote of the lowest lexical `player_id` among voters for tied choices, and no votes selects the first choice in schema order.
- `input_seq` establishes per-player order, not cross-player priority. Selection, supersession, and the discrete-over-continuous rule operate only through the bounded ordered ledger above.

Movement and logical tick v1:

- Authoritative positions are signed fixed-point `Position` values in subtiles; one `16x16` tile equals `256` subtiles on each axis. A living player moves `64` subtiles per `20 Hz` server tick while its held direction is not `none`, for a base speed of `5` tiles/second. Diagonal movement does not exist in v1.
- A `move` input replaces the held direction and remains active until another `move`, phase exit, disconnect, death, or blocked run transition sets it to `none`. The owner advances movement from an explicit `Tick { logical_tick }` command generated by the monotonic scheduler; ticks are state-machine inputs, not wall-clock reads inside `game_core`.
- For each tick, players are considered in lexical `player_id` order. Resolve X/Y as one cardinal displacement, then test map bounds, static blocked tiles, and finally the committed positions of actors already resolved for that tick. A blocked move leaves position unchanged but updates facing. Players not yet resolved retain their prior position for collision, so two players cannot swap or occupy one location in a tick.
- Actor collision uses one axis-aligned `192x192`-subtile box centered on the position. Map edges are blocked; positions and collision boxes must remain within validated map dimensions. Warps and interactions execute after all movement for the tick in lexical player order.
- A server delayed by more than one tick performs one immediate catch-up tick, records all additionally dropped logical ticks in `server_tick_time_dropped_total`, emits one bounded `time_dropped` diagnostic event, and resumes the next scheduled tick. Dropped ticks do not advance logical time or movement, but server-monotonic deadlines still elapse. The owner converts elapsed deadlines into explicit `DeadlineExpired { kind, deadline_id }` state-machine inputs; `game_core` never reads a clock. At an exact boundary, input received strictly before the deadline is staged first, then movement/interactions resolve, then equality belongs to expiry.

### Network Limits Defaults

- Max inbound WS frame size: `16 KiB`
- Max outbound WS frame size: `64 KiB`
- Max reassembled inbound message size: `16 KiB`; protocol JSON nesting depth is capped at `64` before DTO deserialization. WebSocket compression and binary gameplay messages are disabled for v1.
- Input rate limit: `20 msg/s` (burst `4`) per client
- Control rate limit: `5 msg/s` (burst `10`) per client
- Rate limits use continuously refilled token buckets initialized full. Refill uses injected monotonic nanoseconds and fixed-point token-nanoseconds with no floating point; an attempt is admitted when at least one whole token is available, consumes exactly one token, and `retry_after_ms` rounds the missing-token duration up to the next millisecond. Every received attempt consumes the applicable bucket before payload validation after only the bounded envelope classification needed to choose that bucket. Continuous client input is latest-value coalesced to at most `20 msg/s`; a client must not send at its `60 Hz` update rate.
- Heartbeat: the server initiates application `ping` every `10s`; the client answers with `pong`. Disconnect after `30s` without a valid pong.
- Control messages are `join`, `pong`, `rejoin`, `resync_request`, and `resync_ack`. Every state-changing lobby/gameplay `input` uses the input rate limit; server messages are bounded by message-size and broadcast-rate limits.
- Before parsing a Steam ticket or calling Steam, every `create`/`join`/`rejoin` uses the per-source-IP `20` attempts/minute (burst `5`) bucket. Existing-session `join`/`rejoin` additionally use a per-requested-session `30` attempts/minute (burst `10`) bucket; `create` has no session key and instead uses only the source bucket before authentication plus the authenticated per-identity creation limit. The global bucket table holds at most `4,096` entries, expires idle entries after `10 minutes`, evicts the least-recently-used idle entry only, and rejects with `rate_limited` when full with no idle entry. Unauthenticated failures never consume an authenticated gameplay budget.
- Steam validation has a global concurrency limit of `32`, a `3s` total timeout, and no automatic retry on one request. Exhaustion returns recoverable `busy` with `retry_after_ms` between `250` and `1,000`; timeout returns `unauthorized` without disclosing provider detail. A circuit opens for `15s` after `20` consecutive provider failures and permits one probe at a time while half-open.
- Source IP for throttling comes from the direct peer unless the connection came through an explicitly configured trusted proxy. Forwarding headers from any other peer are ignored; ambiguous forwarding chains fail closed.
- `ClientResyncState` is a per-player projection capped at one `64 KiB` outbound message. It excludes token digests, RNG state, other players' inventories/votes, unconsumed hidden triggers, and all other server-only data. The DTO allowlist plus analytical worst-case encoded-size proofs and boundary-complete generated states must establish that every reachable validated state fits; protocol fragmentation is out of scope for v1.
- Max simultaneous unauthenticated handshakes: `64`; max total sockets: `320`; handshake timeout: `10s`.
- Each session has one owner task and a bounded mailbox of `256` commands with reserved capacity for tick/deadline, disconnect, persistence-completion, expiry, and shutdown commands; gameplay ingress cannot consume those reserved slots and returns `busy` without consuming `input_seq` when its partition is full.
- Each connection has one byte-accounted outbound queue of `8` ordered bundle slots and at most `928 KiB` total encoded bytes. At most `7` slots hold semantic bundles; each contains at most one `48 KiB` event batch, one coalescible `64 KiB` state, and `0..=8` requester-only `512`-byte results, for at most `116 KiB`. One slot and `72 KiB` are reserved for control traffic (`join_response`, the ordered `rejoin_response` plus `resync_state` handoff, standalone resync, error, notice, ping). The `928 KiB` figure is the uniform `8`-slot ceiling (`8 * 116 KiB`); the realized worst case is `884 KiB` (`7 * 116 KiB` semantic plus `72 KiB` control). Pings coalesce, resync replaces only coalescible state, and no control or semantic result/event uses an unaccounted side queue. A request that cannot reserve its required output returns `busy` before sequence consumption; if even the reserved slot cannot drain, disconnect that slow client. Restart/termination notices use the reserved slot and are never queued behind replaceable state.
- No request path may spawn unbounded tasks or blocking work. Admission is refused when a global bounded resource is exhausted.
- The theoretical frame-cap ceiling is approximately `320 MiB/s` (`64 KiB * 20 Hz * 256 players`) before overhead. It is not an operating target; Phase 17 records actual p50/p95/p99 state sizes and network throughput and fails if practical instance bandwidth or queue bounds are exceeded.

### Content Pack Loading Limits

- Pack container: one `.lp-pack.zip` archive using stored or DEFLATE entries only. Nested archives and encrypted entries are rejected.
- Max compressed archive size: `64 MiB`; max total expanded bytes: `128 MiB`; max compression ratio per entry and whole archive: `20:1`. Ratio is `expanded_bytes / max(compressed_bytes, 1)`; a zero-byte compressed entry with non-zero output is rejected.
- Max files per pack: `512`.
- Accept regular files only. Reject symlinks, hardlinks, devices, absolute paths, empty or dot segments, backslashes, control characters, Windows reserved names, trailing dots/spaces, and case-insensitive path collisions.
- Paths use portable ASCII POSIX segments matching `[a-z0-9][a-z0-9._-]{0,63}` with a maximum total length of `240` bytes.
- Max general file size: `16 MiB`.
- Max single JSON file size: `1 MiB`.
- Max JSON nesting depth: `64`; max string length: `64 KiB`; max `100,000` total JSON tokens, `25,000` aggregate object members, and `25,000` aggregate array elements per pack. Duplicate object keys are rejected before deserialization.
- Max story nodes per story: `512`.
- Max choices per node: `8`.
- Stories per aggregate pack: `1..=64`; validation rejects a session pack with zero stories or a manifest whose `includes` list exceeds this story count even though the mixed-category list has its separate cap.
- Max flags per story: `256`; max encounters per story: `128`; max automatic transitions handled for one intent: `56`, reserving at least `8` event items for dice/combat/death/run/diagnostic events within the `64`-item batch cap.
- Max spells per character: `8`; max characters per aggregate pack: `64`.
- Max locally installed community packs: `32`; max aggregate installed-pack disk use: `4 GiB`. Accounting includes retained archives, extracted generations, metadata, and previous generations awaiting deletion. Import temporary files have a separate process-wide `256 MiB` cap; before extraction, require free space for current accounted use plus the declared final and temporary maxima. Import refuses before crossing a limit, removes temporary data on every outcome, and never evicts without user action.
- Max TMX file size: `4 MiB`.
- Max TMX map dimensions: `512x512` tiles. Per map: max `16` tilesets, `16` layers, `2` object layers, `1,024` objects, and `4,096` properties. Per aggregate pack: max `100,000` total XML elements. Max XML depth is `32`, attributes per element `64`, and attribute value `1 KiB`. Per map, each CSV tile-layer text node is capped at `2 MiB` and total decoded tile-layer data at `3 MiB`, where decoded size counts exactly `4` bytes per GID. Both caps are sized for the maximum layer set at maximum dimensions: `512 * 512 * 4` bytes is `1 MiB` per layer, so `ground`, `collision`, and the optional `decor` layer reach `3 MiB` together, and one layer's CSV text reaches `2 MiB` at the widest `7`-digit GID plus separator.
- TMX tile layers use CSV encoding only in v1. Base64, gzip, zlib, zstd, infinite/chunked maps, templates, image layers, inline tilesets, scripts, and parser-driven URL/file resolution are rejected. A TMX may reference only an in-pack relative TSX after portable path normalization and membership validation; that TSX may reference validated in-pack images. The loader opens those already-validated archive members directly and never asks the XML parser to resolve them.
- Image files are capped at `8 MiB`, `4096x4096` pixels, and `64 MiB` decoded bytes each. Aggregate decoded image memory during validation is capped at `128 MiB`.
- Audio files are capped at `16 MiB`, `10 minutes`, `2` channels, `48 kHz`, and `64 MiB` decoded bytes each.
- Client runtime retained decoded limits per active pack are `128 MiB` images plus `128 MiB` audio, with `256 MiB` total assets. During a theme/pack swap, process-wide decoded assets are capped at `512 MiB` and loader concurrency at `2`. Inspect headers and computed decoded sizes before allocation; refuse the load deterministically rather than implicitly evicting active assets. The headless server never decodes presentation media.
- TMX/XML parsing disables DTDs, external entities, XInclude, and all external resource resolution. Untrusted image/audio validation and decode run in a sandboxed worker subprocess with no network, a `2s` per-file wall timeout, process memory cap `256 MiB`, concurrency `2`, and declared decoded-output caps even when headers lie; timeout or limit kills the worker before its slot is reused. The validation supervisor reserves worker RSS, decoded output, IPC, and temporary bytes before spawn; all aggregate/media workers for one process share a `1 GiB` reservation ceiling and a `256 MiB` temporary-disk ceiling, and unreserved jobs remain in a bounded queue of `8` packs/`512 MiB` compressed bytes before rejection. Workers inherit no secrets or unrelated descriptors, receive only pre-opened input/output handles, cannot create child processes, and are confined to a fresh temporary directory. Validation and runtime loading enforce the same limits; community packs remain hostile after hub validation.

## Locked Pack Compatibility And Checksums

MVP sessions use exactly one aggregate pack archive. That archive may include story, character, and theme content and may reference the versioned built-in identifier catalog. Multiplayer sessions require the server and every player to match all of:

- `pack_id`
- `version`
- canonical SHA-256 checksum
- `content_schema_version`
- `game_rules_version`

`pack_manifest.json` is required at the archive root. Every content/pack identifier matches `[a-z0-9][a-z0-9_-]{0,63}`, is lowercase ASCII, and is unique in its document namespace; IDs are compared byte-for-byte without normalization. File path segments use the separate dot-permitting path grammar above. Pack `version` is valid SemVer 2.0.0 without build metadata for MVP compatibility comparisons. The checksum identifies canonical logical content: ZIP metadata, entry order/compression, and insignificant manifest representation may differ while the validated path/byte set is equal. It does not prove raw archive-byte equality, authorship, safety, or trust.

Checksum v1 rules:

- The checksum covers the whole pack, not individual files.
- The pack is treated as the validated set of portable relative POSIX paths plus bytes defined by the loading limits above.
- Sort paths lexicographically by UTF-8 bytes before hashing.
- Hash input starts with `les-perissables-pack-v1\n`.
- For each sorted file, append `path\0length\0bytes\n`, where `length` is the decimal byte length of the hashed bytes.
- For `pack_manifest.json`, reject duplicate keys, omit the top-level `checksum` and `publication_attestation` fields, then serialize with RFC 8785 JSON Canonicalization Scheme before hashing. The attestation is detached metadata over the resulting checksum and cannot affect logical-content identity.
- All other files are hashed as raw bytes.
- The manifest stores the checksum as lowercase hex prefixed with `sha256:`.

### Publication Attestation v1

All signatures use Ed25519. Public keys are exactly `32` bytes and signatures exactly `64` bytes, encoded as unpadded RFC 4648 base64url. Signed bytes are the listed ASCII domain prefix followed immediately by the RFC 8785 canonical JSON encoding of the object with its signature field omitted. `JcsUint` means a JSON integer in `1..=9,007,199,254,740,991` (`2^53-1`), preserving exact RFC 8785/IEEE-754 identity; larger values, fractions, exponent forms that do not canonicalize to the same integer, and negative/zero values are rejected. Unknown/duplicate fields, non-canonical value types, invalid encodings, and signatures over any other representation are rejected.

`PublicationAttestationV1` is the optional `pack_manifest.json.publication_attestation` object:

| Field | Exact rule |
| --- | --- |
| `attestation_version` | integer, exactly `1` |
| `pack_id` | exactly the containing manifest `pack_id` |
| `pack_version` | exactly the containing manifest `version` |
| `pack_checksum` | exactly the recomputed canonical pack checksum |
| `publication_generation` | `JcsUint`; immutable for one published object generation |
| `policy_tier` | exactly `tier2` |
| `issued_at_unix_ms` | `JcsUint` |
| `expires_at_unix_ms` | `JcsUint`, greater than issuance and no later than checked issuance plus `30 days` |
| `signing_key_id` | `ContentId` selecting one active publication key |
| `signature` | unpadded base64url Ed25519 signature |

Its domain prefix is `les-perissables-publication-attestation-v1\n`. The signed identity includes no archive URL or mutable listing metadata. One attestation authorizes only the exact logical checksum/generation/tier tuple and never authorizes a different archive with the same listing name.

`PublicationKeySetV1` is fetched from the one hub origin pinned in trusted build configuration and cached outside pack-controlled storage:

| Field | Exact rule |
| --- | --- |
| `key_set_version` | integer, exactly `1` |
| `generation` | `JcsUint`; any positive generation initializes an empty cache, otherwise it must be strictly newer than the replaced cache |
| `issued_at_unix_ms`, `expires_at_unix_ms` | `JcsUint`; expiry is greater than issuance and no later than checked issuance plus `90 days` |
| `keys` | `1..=16` entries with unique `key_id` |
| `root_signature` | signature by the offline root key pinned in the game release |

Each key entry has exactly `key_id: ContentId`, `usage` (`publication` or `revocation`), `public_key`, `not_before_unix_ms: JcsUint`, and `not_after_unix_ms: JcsUint`; `not_before < not_after`, the interval is contained by the key-set lifetime, and checked arithmetic must not overflow. The key-set domain prefix is `les-perissables-publication-key-set-v1\n`, and signing omits `root_signature`. Root rotation requires a game release that pins the new root while retaining the old root for the prior key-set overlap; hub data alone cannot replace the root of trust.

`PublicationRevocationListV1` has exactly `revocation_version: 1`, `generation: JcsUint`, `issued_at_unix_ms: JcsUint`, `expires_at_unix_ms: JcsUint`, `signing_key_id`, `revoked_key_ids: ContentId[0..=16]`, `revoked_publications: RevokedPublication[0..=1024]`, and `signature`. Any positive generation initializes an empty cache; replacement must be strictly newer. Expiry is greater than issuance and no later than checked issuance plus `24 hours`. Key IDs are unique and sorted lexicographically by ASCII bytes. `RevokedPublication` has exactly `pack_id`, `pack_version`, and `publication_generation: JcsUint`; entries are unique and sorted by `pack_id` ASCII bytes, then `pack_version` ASCII bytes, then numeric generation. The signing key must be active with `usage: revocation`. Its domain prefix is `les-perissables-publication-revocations-v1\n`, and signing omits `signature`.

Verification order is fixed: recompute pack checksum; establish trusted time; require `attestation.issued_at <= trusted_time < attestation.expires_at`; validate exact attestation identity and checked lifetime; validate the root-signed key set with `issued_at <= trusted_time < expires_at`; require the selected key's `not_before <= trusted_time < not_after` and correct usage; validate the revocation-list signature with `issued_at <= trusted_time < expires_at`; reject revoked key IDs and publication tuples; then verify the attestation signature. Trusted time is `max(system_wall_time, last_trusted_hub_time)`; authenticated hub responses advance and persist `last_trusted_hub_time` in pack metadata, never move it backward, and a local clock more than `5 minutes` behind that value fails closed. Every duration addition is checked before comparison. If the key set or revocation list cannot be refreshed before expiry, Tier 2 activation is unavailable offline; Tier 1 remains usable. No validation worker performs this network fetch.

Attestation validity is checked at import and immediately before each run activation. A run that has already activated an immutable pack may finish using its captured validation generation if metadata merely expires while offline; no new run may start. A newly fetched revocation or release-embedded emergency denylist prevents further asset loads and terminates a creator-test run at the next scene boundary with `pack_unavailable`. Conformance vectors cover exact signed bytes, Unicode canonicalization, wrong domain/version/identity/key usage, root and online-key rotation, future/expired times, stale/downgraded generations, offline expiry, and both revocation forms.

The game loader, `storycheck`, and the community hub must all use the same pinned release of the MIT validation crate for path normalization, limits, canonicalization, and checksum behavior. That crate owns a shared conformance corpus with positive digests and rejection vectors for empty packs, reordered files, Unicode JSON values, duplicate keys, traversal, Windows aliases, symlinks, boundary sizes, checksum stability with/without detached attestation, valid signatures, wrong checksums/key IDs, expiry, and revocation.

Hosted Release-ready multiplayer serves built-in aggregate packs only. The server loads them from its immutable release artifact and verifies their canonical checksums at startup. Clients cannot upload, provide a URL, or cause the game server to fetch a pack. Community-pack multiplayer remains disabled until the hub provides an operator-controlled, authenticated provisioning feed with immutable checksum-addressed objects, publication attestations, bounded cache/storage/retention, and rollback; adding that feed is a post-MVP contract change. Local Tier 1 community packs remain available for non-hosted validation and single-player creator testing only.

## Locked Content Conventions

### Schema v1

- Story, character, theme, and pack documents use `content_schema_version: 1`. Every field below is present in the Phase 04 schema; later phases implement behavior rather than mutate schema v1.
- Story metadata requires `theme_id`. The referenced theme supplies the single `ui_variant` field used to select visual skin tokens; `default_skin` is not a separate field in v1.
- Story checks resolve to `success`, `failure`, `critical_success`, or `critical_failure`. `success` and `failure` branches are required. Missing critical branches fall back to their corresponding normal branch while retaining critical UI/audio feedback.
- Unknown fields are rejected. Identifiers are portable ASCII, at most `64` bytes, and unique within their namespace.
- Character presets contain at most `8` spells. Duplicate character picks are forbidden in MVP sessions.
- The schema supports custom asset references from v1. The release runtime loads custom media only when the pack carries a valid hub publication attestation described under the content-tier rules; unsigned/local Tier 1 packs may reference only the shipped identifier catalog.

All schema-v1 documents are UTF-8 JSON objects, use the exact field names below, reject duplicate/unknown fields, and obey the aggregate parser bounds. `ContentId` uses the identifier grammar above; `AssetRef` is either `builtin:<ContentId>` or `pack:<portable-relative-path>`. `pack:` is Tier 2 and follows its attestation gate. `LocalizedText` is a non-empty UTF-8 string of at most `2,000` Unicode scalar values in v1; localization tables are post-MVP.

Pack manifest v1 (`pack_manifest.json`):

| Field | Type and rule |
| --- | --- |
| `content_schema_version` | integer, exactly `1` |
| `pack_id` | `ContentId` |
| `version` | SemVer 2.0.0 string without build metadata, max `64` bytes |
| `checksum` | `sha256:` plus `64` lowercase hex digits |
| `game_rules_version` | integer, exactly `1` |
| `license_expression` | SPDX expression string, `1..=256` bytes |
| `component_licenses` | object with `0..=512` portable paths; each value has required `license` (`1..=128` bytes) and optional `attribution` (`1..=512` scalars) |
| `includes` | `1..=128` unique strings tagged `story:`, `character:`, or `theme:` followed by `ContentId`; `1..=64` entries must be tagged `story:`; every listed document must exist and unlisted content documents are rejected |
| `publication_attestation` | optional `PublicationAttestationV1` defined under the attestation contract; forbidden for Tier 1 authoring output |

Story document v1 (`stories/<story_id>.json`):

| Field | Type and rule |
| --- | --- |
| `content_schema_version` | integer, exactly `1` |
| `id` | `ContentId`, equal to filename stem |
| `title` | `LocalizedText`, max `128` scalars |
| `theme_id`, `map_id`, `start_node` | `ContentId`; references must resolve; `map_id` resolves exactly one built-in catalog map or `maps/<map_id>.tmx` archive member |
| `estimated_minutes` | integer `1..=120` |
| `nodes` | `1..=512` unique `StoryNode` objects |
| `encounters` | `0..=128` unique `Encounter` objects |
| `triggers` | `1..=128` unique `StoryTrigger` objects |

A `StoryTrigger` has `event_id: ContentId`, `node_id: ContentId`, and `once: bool`. Every TMX `events`-layer object `event_id` must map to exactly one trigger and every trigger must map to one object and reachable node. An accepted world `interact` for that server-visible and reachable event enters `story` at `node_id`; a once-only trigger is marked consumed in the same staged transaction.

Every active story sequence carries one authoritative `story_actor_player_id: PlayerId`. Trigger entry uses the interacting living player. `start_run` and any other system entry use the lowest lexical living `player_id`. Votes and automatic transitions retain the current player rather than changing it to the last voter. Checks read that player's selected character stat, `dice_result.actor_id` uses the associated character `EntityId`, and `Effect` target `actor` refers to that character/inventory. If grace expiry or another committed transition makes the character non-living before the next check/effect, rebind once to the lowest lexical living player; if none exists, enter failed summary. The player/character mapping and context are persisted and included in deterministic transcripts.

`StoryNode` has required `id: ContentId`, `kind`, and `effects: Effect[0..=16]`. Its remaining required fields depend on `kind`:

| Node kind | Required fields |
| --- | --- |
| `dialogue` | `speaker_id: ContentId`, `text: LocalizedText`, `choices: Choice[1..=8]`; each choice has unique `id`, `label: LocalizedText` capped at `256` scalars, `goto`, and optional `requires_flag` |
| `check` | `label: LocalizedText`, `stat` (`strength`, `agility`, `wit`, `charm`), `outcomes`; outcomes require `success` and `failure` and optionally contain `critical_success`/`critical_failure`, each a `Branch` |
| `encounter` | `encounter_id: ContentId`, `on_win: ContentId`; a party wipe bypasses story continuation and enters failed summary |
| `transition` | `goto: ContentId`; executes automatically after common effects and counts toward the automatic-transition limit |
| `return` | no additional fields; commits the current story sequence back to `world` at the unchanged authoritative positions |
| `end` | `result` (`completed` or `failed`) |

A `Branch` has `goto: ContentId` and `effects: Effect[0..=16]`. An `Effect` is one of `set_flag { flag_id }`, `unset_flag { flag_id }`, `give_item { item_id, recipient }`, or `remove_item { item_id, owner }`; `recipient`/`owner` is `actor`, `lowest_hp_living`, or `lowest_player_id_living` under the story-actor rule above. Flags referenced anywhere count toward the story's `256`-flag limit. Every `goto`, trigger, encounter, item, speaker, map, theme, and built-in behavior reference must resolve during validation. The graph roots are `start_node` and every trigger `node_id`. From every root, each path must reach `return`, `end`, an unresolved player choice, check, or encounter within `56` automatic transitions; every node must be reachable from at least one root, unbounded automatic cycles are rejected, and the worst-case required event count must fit the reserved event budget.

An `Encounter` has `id`, `enemies: Enemy[1..=16]`, and `rewards: Reward[0..=16]`. `Enemy` has unique `id`, `name: LocalizedText` capped at `64` scalars, four stats in `5..=70`, `max_hp: 1..=999`, `resource_max: 0..=100`, and `spell_ids: ContentId[0..=8]` from the built-in catalog. `Reward` has `item_id` and `quantity: 1..=4`; total reward items are capped at `16` before inventory placement.

Character document v1 (`characters/<character_id>.json`):

| Field | Type and rule |
| --- | --- |
| `content_schema_version` | integer, exactly `1` |
| `id` | `ContentId`, equal to filename stem |
| `name` | `LocalizedText`, max `64` scalars |
| `group` | `ContentId` |
| `stats` | exactly `strength`, `agility`, `wit`, and `charm`, each integer `5..=70` |
| `max_hp` | integer `1..=999` |
| `resource_max` | integer `0..=100` |
| `spells` | `0..=8` unique built-in spell IDs |
| `starting_items` | `0..=4` unique built-in item IDs |
| `sprite` | `AssetRef` |
| `voice_barks` | object with `0..=16` event-key to `AssetRef` entries; text fallback remains mandatory in engine data |

Theme document v1 (`themes/<theme_id>.json`):

| Field | Type and rule |
| --- | --- |
| `content_schema_version` | integer, exactly `1` |
| `theme_id` | `ContentId`, equal to filename stem |
| `tileset`, `props`, `ambience_loop`, `exploration_music`, `combat_music`, `combat_backdrop` | `AssetRef` |
| `ui_variant` | `ContentId`; references a built-in behavior-neutral skin-token set |

Schema validation occurs in this order: archive/path/resource bounds, JSON syntax/duplicate keys, structural field/type/bounds, identifier uniqueness, cross-reference resolution, TMX semantics, graph/automatic-transition and worst-case event/projection-size validation, license/attestation policy, then canonical checksum. ZIP/path/checksum and JSON/TMX/XML parsing/semantic validation run only in the sandboxed no-network workers described above, with at most `4` aggregate workers and `2` media workers subject to the shared reservation ceiling. Each stage has a `5s` wall deadline and one pack has a `15s` aggregate deadline; an individual aggregate worker has a `512 MiB` RSS cap. Deadline, cancellation, worker crash, or limit kills and reaps the subprocess, removes temporary output, returns `validation_resource_limit`, and installs nothing; parser work never runs in a non-killable in-process thread. A failure returns one bounded path-based diagnostic and performs no media decode or network access before cheaper structural checks pass. Tier 1 rejects every non-manifest/content file and every in-pack map/media member; Tier 2 classifies permission from the complete validated archive member set, so an unreferenced custom file cannot bypass attestation.

### TMX

- A custom map has the only valid path `maps/<map_id>.tmx`; `<map_id>` is `ContentId` and equals the filename stem. It must not collide with a built-in map ID. A story map reference resolves exactly one source. Any `maps/` member makes the pack Tier 2 and requires a valid publication attestation.
- Maps are finite orthogonal TMX maps with width and height each in `1..=512`, `tilewidth="16"`, `tileheight="16"`, and no infinite/chunked data. Layer names are unique, and the complete layer set is exactly tile layers `ground` and `collision`, object layer `events`, and optional tile layer `decor`; every extra group/image/tile/object layer is rejected. Each tile layer matches map dimensions and its CSV contains exactly `width * height` unsigned GIDs. Flip/rotation flag bits and unknown layer/object types are rejected.
- GID `0` in `collision` is walkable and every non-zero GID is blocked. `ground` must contain a valid non-zero tile for every walkable cell. Inline tilesets are rejected. Each of at most `16` TMX tileset entries has a strictly increasing non-zero `firstgid` and references one validated `tilesets/<tileset_id>.tsx` member. A TSX has `tilewidth="16"`, `tileheight="16"`, `tilecount: 1..=65,535`, no external source, and validated in-pack image references only. Checked intervals `[firstgid, firstgid + tilecount)` may not overlap or overflow, and every non-zero CSV GID resolves exactly one interval; image bytes are reached only through that TSX, never by interpreting a GID as a path. External URLs and parser resolution remain forbidden. Map boundaries are blocked regardless of tile data.
- `events` accepts only axis-aligned rectangle objects with integer pixel coordinates/sizes aligned to the `16x16` grid, no rotation, and names `evt_<id>`, `spawn_<id>`, or `warp_<id>`, where `<id>` is `ContentId`. Names and semantic IDs are unique. Unknown properties and object shapes are rejected.
- Exactly one `spawn_party` rectangle is required. It is exactly `32x32` pixels over four walkable cells; row-major cell centers are assigned to players in lexical `player_id` order when `start_run` stages initial positions. Additional `spawn_<id>` objects are exactly one walkable tile and are warp destinations.
- An `evt_<id>` object is exactly one tile and has exactly `event_id: ContentId`, equal to `<id>`. At runtime the server issues its opaque `EntityId`. It appears in a player's `active_interaction_ids` only when that living player's center is in the event tile or one cardinally adjacent walkable tile; the submitted ID must be in that recipient-specific list at commit time. Diagonal or remote interaction is rejected with `target_unavailable`.
- A `warp_<id>` object is exactly one walkable tile and has exactly `target_spawn_id: ContentId`, which resolves an additional spawn in the same map. Entering the warp after movement places the actor at that destination if its collision box is free; otherwise the move remains at the pre-warp position. A warp never chains again in the same tick.
- Semantic validation proves spawn/collision bounds, event/trigger bijection, warp targets, story map resolution, and that every player collision box remains inside the map at every spawn/warp destination. It evaluates every reachable player-center tile under the exact cardinal interaction rule and rejects a map where more than `32` event IDs could enter one recipient's `active_interaction_ids`. Positive and rejection conformance vectors fix all of these interpretations.

### Assets

- Character frames: `16x16`
- Direction row order: down, left, right, up
- Frames per direction: `4` (`0` idle, `1-3` walk)
- Character naming pattern: `char_<group>_<name>.png`
- Theme tileset naming pattern: `tileset_<theme_id>.png`
- Theme prop atlas naming pattern: `props_<theme_id>.png`
- Combat backdrop naming pattern: `backdrop_<theme_id>_<name>.png`
- Ambience loops: Ogg container with Vorbis audio, `.ogg` (`ambience_<theme_id>_<name>.ogg`)
- Music: Ogg container with Vorbis audio, `.ogg` (`music_<theme_id>_<track>.ogg`)
- Voice bark clips: Ogg container with Vorbis audio, `.ogg` (`voice_<actor>_<line_id>.ogg`)
- SFX: uncompressed little-endian PCM WAV, `.wav`, with `16`-bit signed samples (`sfx_<category>_<name>.wav`)
- Accepted audio is mono or stereo at `48 kHz`; validators reject other codecs, sample formats, channel counts, or rates rather than relying on platform decoder behavior

### Gameplay Bounds v1

- Every accepted input executes against a bounded staged copy. Validate phase, actor, target, turn, resource, inventory, and references before consuming RNG. Apply the direct action, then authored branch effects in array order, then automatic `goto` transitions and their effects in order. Emit events in the same order after state is committed.
- If an effect, reference, checked arithmetic operation, event bound, or the `56`-automatic-transition limit fails, discard the complete staged gameplay/RNG mutation. Commit only consumed replay metadata and return `rejected/limit_reached` or `rejected/invalid_action`; partial flags, damage, rewards, movement, or events are forbidden.
- Story votes do not consume RNG. On vote closure, resolve the plurality rule, apply the selected branch transaction once, and clear all votes before the next node becomes visible.
- Each living character has `4` inventory slots. The recipient sees its own occupied slots through `self_inventory`. After combat, dead party actors remain as world corpses at the encounter location; a corpse exposes at most `2` carried items whose immutable catalog entry has `lootable: true`, selected by ascending source slot then item ID. Fewer are exposed when fewer exist.
- `loot` transfers exactly the requested currently exposed `corpse_slot` to the requesting living player's lowest free inventory slot in one transaction. The corpse must be in that recipient's `WorldView.corpses`, and the request must satisfy the same cardinal reachability rule as interaction. A stale/missing slot rejects with `item_unavailable`; no free slot rejects with `inventory_full`. Rejection changes neither inventory nor corpse and consumes no RNG.
- HP/max HP: `0..=999`; action resource pools: `0..=100`; individual costs: `0..=100`; individual damage/healing effects: `0..=999`; inventory slot index: `0..=3`.
- Encounters contain at most `16` enemies and `20` total actors. Combat is capped at `256` rounds and `5,120` actor turns; reaching either cap ends the encounter as a failed run with a stable limit event.
- Arithmetic uses checked operations in a wider intermediate type, then clamps HP/resources to their validated maxima. Overflow, underflow, or an out-of-range authored value is a validation/programmer error, never wrapping behavior.
- Combat turn order is descending agility. Ties are resolved by stable lexical actor ID. Each actor performs exactly one of `attack`, `spell`, `item`, or `pass` per turn.
- Invalid, unaffordable, dead-actor, or out-of-turn actions return `input_result` with `outcome: rejected` and the applicable stable `reason_code`; they never return an `error` envelope for the semantic rejection. They do not mutate gameplay state or consume the turn. Their expected `input_seq` still advances replay metadata and is durably recorded after Phase 16.
- An actor at zero HP is dead and removed from future turns. Combat ends when all enemies are dead or all player characters are dead; a party wipe produces the failed-run summary.
- `attack` checks the attacker's `strength` using the locked d100 resolver. Normal success deals `max(1, floor(strength / 5))` damage; critical success deals twice that value; failure and critical failure deal zero. Critical failure has no additional self-damage in v1. Damage is applied after the dice event and before death events.
- A spell is immutable built-in engine data with exact `id`, `check_stat`, `resource_cost`, `target_side`, `effect_kind`, `normal_amount`, and `critical_amount`. Phase 07 will check in the v1 catalog at `stories/builtin/game_rules_v1.json`, canonicalize it with RFC 8785, and record its SHA-256 in the acceptance transcript and release manifest. `effect_kind` is `damage`, `heal`, or `restore_resource`; no status-effect scripting exists in v1. A cast deducts cost, performs one d100 check, applies normal/critical amount on success, and applies zero effect on failure. Cost remains spent on a failed roll but not on an invalid action.
- An item is immutable built-in engine data with exact `id`, `target_side`, `effect_kind`, `amount`, and `lootable: bool`. A valid use removes exactly one item before applying its non-random effect. Invalid target or full-resource no-op rejects without consumption. Creator packs may compose only catalog IDs and cannot define behavior.
- `give_item` places one item in the recipient's lowest free slot; no free slot fails the whole authored transition with `limit_reached`. `remove_item` removes the matching item in the owner's lowest occupied slot; no match fails the whole transition with `invalid_action`. Both failures follow the staged full-rollback rule, consume no RNG beyond earlier successfully staged rule steps that are themselves rolled back, and emit no partial inventory event/state.
- On an enemy turn, choose the first spell in lexical spell-ID order that is affordable and has a legal target; otherwise attack. For damage choose the living party actor with the lowest HP percentage by cross multiplication, then lowest absolute HP, then lexical actor ID. For healing/resource restoration choose the legal ally with the greatest missing amount, then lexical ID. Enemy actions use the same validation, RNG, cost, effect, and event ordering as player actions.
- After victory, sort rewards by `item_id`, then offer each unit to living players in lexical `player_id` order with a free slot. Unplaced rewards are discarded and reported in the summary; they never overfill inventory. Corpse exposure chooses catalog items with `lootable: true` by ascending inventory slot and then item ID. Loot transfer is atomic and never consumes RNG.
- A round increments after every actor present at round start has either acted or died. New actors cannot spawn in v1. Combat is capped at `256` rounds and `5,120` actor turns; reaching either cap stages a failed-run summary with one `limit_reached` run event.

### Audio And Voice Scope

- Gameplay audio includes footsteps, object-use sounds, dice result sounds, spell SFX, ambience, and combat transitions.
- Voice is tiny and optional (short barks/short lines only); text remains primary.
- No in-game voice chat (players use external voice apps like Discord).
- Multiplayer audio is event-driven: server emits gameplay events, clients play local audio.

Audio channels for v1:

- `ambience`
- `music`
- `sfx`
- `voice`

## Locked Repo/Legal/CI Baseline

- Main repo (`perissables`; project/package prefix `les-perissables`): add `COPYRIGHT`, ARR `LICENSE`, `README.md` on day 1
- Stories repo (`les-perissables-stories`): add MIT `LICENSE`, `README.md` when Phase 01 creates it
- The stories repo will host both the MIT pack schema/validation crate and the `storycheck` CLI, so creators can validate packs without any proprietary game-repo code
- Shared validation package name: `les-perissables-pack`. Release it from `les-perissables-stories` to crates.io with a matching repository tag; game and hub consumers pin the same exact crate release in workspace dependencies and `Cargo.lock`. Production builds do not follow a moving Git branch.
- Rust toolchain: `1.95.0`, edition `2024`, and MSRV `1.95.0`, pinned in `rust-toolchain.toml`. Toolchain changes are explicit contract updates.
- Commit the workspace `Cargo.lock`; CI and release builds use `--locked`. Git dependencies require an immutable revision, and published crates use exact compatible version requirements plus lockfile resolution.
- `mise` is for local developer convenience only: `mise.toml` may pin tools and expose tasks, but CI must not depend on `mise`
- CI installs Rust through `rustup`/standard Rust tooling, pins `cargo-audit` to `0.22.1`, and runs the locked `cargo` checks directly. Initial features must be mutually compatible; if that changes, replace `--all-features` with an explicit supported-feature matrix.
- Minimum CI checks on first commit:
  - `cargo fmt --all --check`
  - `cargo clippy --locked --all-targets --all-features -- -D warnings`
  - `cargo test --locked --all-features`
  - `cargo test --locked --doc --all-features`
  - `RUSTDOCFLAGS="-D warnings" cargo doc --locked --no-deps --all-features`
  - `cargo audit --file Cargo.lock`
  - `cargo deny check`
- `cargo-deny` and its configuration are pinned in the same tool manifest as `cargo-audit`. Advisory scanning establishes a lead, not reachability; any exception records the advisory/license/source, owner, rationale, expiry, and validation evidence.
- Before Phase 01 completes, build documentation pins the `raylib-rs` crate, native `raylib` source/archive and checksum, system-versus-bundled policy, static/dynamic linkage, target prerequisites, and license obligations. It also pins the Steamworks Rust/native SDK acquisition and license policy even though authenticated integration lands later. No native dependency may be inferred from a developer machine.
- Linux and Windows release targets must compile in CI before Phase 02 completes. The exact native dependency documents and checksums are completion evidence for that gate.
- Reproducibility acceptance builds each target twice from separate clean checkouts with the same pinned source, lockfile, toolchain, native SDKs, target, release profile, environment, and `SOURCE_DATE_EPOCH`. Version/timestamp metadata is normalized. After stripping per policy, executable, packaged assets, manifest, and archive checksums must match byte-for-byte; any documented platform-signature envelope is compared after removing only that envelope. Commands, container/image digest, checksums, and any allowed nondeterminism are retained.
- Windows release executables and installers are Authenticode-signed. The private key resides in a non-exportable managed signing service or hardware-backed CI identity, is unavailable to pull-request jobs and developers, and signs only protected release-tag workflows after artifact checksum approval. CI verifies chain, subject, timestamp, and file digest before depot handoff. Linux packages publish SHA-256 checksums and provenance; signing/provenance secrets follow the same protected-job boundary.
- Outside pull requests are not accepted until a qualified legal review supplies an explicit inbound contribution policy requiring an appropriate written contributor agreement. Steam distribution likewise requires reviewed end-user terms and a third-party notice/license inventory before Release-ready exit.

## Release Rollback v1

- Every final candidate produces `artifacts/releases/<version>/rollback-manifest.json` before promotion. It records the candidate and retained-prior Git SHA, server image digest, Linux/Windows package and Steam depot manifest IDs, protocol/content/rules/save/settings versions, built-in pack checksums, non-secret configuration snapshot digest, database/volume migration state, exact rollback commands, operator, and evidence paths. The prior server image and both prior depots remain immutable and available until the next release clears its rollback window.
- The server rollback procedure marks the candidate unready, stops admission, checkpoints, performs the clean lease release, deploys the exact retained digest/config, restores every session, verifies `/healthz` and `/readyz`, and runs create/join/rejoin plus one continued-run transcript before traffic returns. Clean process-start-to-ready must meet the `20s` clean-restore gate; the complete staging rollback drill, including smoke checks, must finish within `5 minutes`.
- Client rollback switches the isolated Steam test branch to the recorded prior depot manifests for both OS targets, installs from a clean cache, verifies signatures/checksums, launches, and completes the compatibility transcript. Production rollback uses the same recorded manifests but is never exercised destructively as a release test.
- A candidate is rollback-compatible only when the retained server can read every authoritative save version the candidate may write and the coordinated client/server pair uses a compatible protocol/rules/pack tuple. A save-version change therefore uses expand-and-contract: before any release writes version `N`, the retained rollback build must already contain and pass a read-only `N` migrator while still writing `N-1`. If that preparatory release was not deployed, the version-changing candidate cannot pass Release-ready. An older client may preserve and ignore a newer settings file and use canonical defaults as already specified, because settings are non-authoritative. No rollback path may rewrite a newer settings/save file destructively or restore an older snapshot that loses an acknowledged transition.
- A breaking protocol/rules change rolls client depots and server back as one coordinated operation; mixed incompatible versions remain rejected with `version_mismatch`. Failure of any digest, compatibility, state, deadline, or smoke oracle aborts the rollback drill and blocks release rather than improvising another artifact.

## Privacy And Data Minimization

- No optional telemetry by default in the game client/server.
- No in-game voice chat and no voice recording storage.
- Use only gameplay-required platform/session identifiers.
- If crash upload is added, it must be opt-in and clearly documented.
- Community hub (separate context): the hub collects account data and user-generated content - OAuth provider IDs, display names, shared packs, likes, and comments - which is distinct from the game's minimal-PII stance. The hub must publish a privacy policy, store only what these features require, and support account and content deletion.

Retention defaults are maxima unless law, an active security incident hold, or an explicit user request requires a documented exception:

| Data class | Storage/access | Retention and deletion |
| --- | --- | --- |
| Game structured logs | Railway logging; on-call/maintainer only; no tickets/tokens/message bodies | `14 days`, then deletion |
| Durable active-session snapshots/WAL | Railway volume; game service and authorized operator | Delete after committed session end plus `24 hours`; encrypted provider backups expire within `30 days` and are not used to resurrect deleted sessions |
| Local crash logs | Player device; player controls access | Rotate at `5` files or `20 MiB`; delete through settings/uninstall guidance |
| Opt-in crash uploads | Restricted crash store; maintainers only; redacted before send | `30 days`; consent screen states fields and deletion route |
| QA/benchmark artifacts | Restricted project storage; maintainers/testers | Raw logs/samples `90 days`, summaries without identifiers retained for release history |
| Hub account/content | Postgres/R2; service and authorized moderation roles | Until account/content deletion or policy-defined takedown; public bytes removed through reference-counted deletion, backups expire within `30 days` |
| Moderation audit records | Restricted Postgres tables; authorized moderators only | `1 year`, with payload minimization and provider subject pseudonymization after account deletion unless legally required |

Versioned synthetic conformance and benchmark fixtures, workload definitions, expected counters, minimized non-sensitive fuzz regressions, and redacted summaries are source artifacts rather than raw operational data and are retained for the supported release history. They contain no account identifiers, credentials, production logs, or personal data. The `90-day` limit applies to raw run samples, captures, and operational logs.

## Client-Local Settings v1

- The client stores `client_settings.json` under the locked per-user data root as one strict UTF-8 JSON object capped at `64 KiB`. It never contains authoritative run state, Steam tickets, rejoin tokens, provider subjects, or server endpoints.
- The canonical version-1 defaults are exactly:

```json
{
  "client_settings_version": 1,
  "volumes": { "ambience": 80, "music": 80, "sfx": 100, "voice": 100 },
  "key_bindings": {
    "move_up": ["W", "ArrowUp"],
    "move_down": ["S", "ArrowDown"],
    "move_left": ["A", "ArrowLeft"],
    "move_right": ["D", "ArrowRight"],
    "interact": ["E"],
    "ui_confirm": ["Enter", "Space"],
    "ui_cancel": ["Escape"],
    "ui_next": ["Tab"],
    "ui_previous": ["Shift+Tab"],
    "combat_attack": ["1"],
    "combat_spell": ["2"],
    "combat_item": ["3"],
    "combat_pass": ["4"]
  },
  "display": { "mode": "windowed", "width": 1280, "height": 720, "ui_scale_percent": 100 },
  "accessibility": { "reduced_motion": false }
}
```

- Every volume is an integer `0..=100`. Display mode is `windowed` or `fullscreen`; width/height is exactly `1280x720` or `1920x1080`; UI scale is `100` or `200`. Supported chord values are ASCII `A` through `Z`, `0` through `9`, `ArrowUp`, `ArrowDown`, `ArrowLeft`, `ArrowRight`, `Enter`, `Space`, `Escape`, `Tab`, and `Shift+Tab`. Each listed action is mandatory with `1..=2` unique supported key chords; unknown actions/chords, unknown fields, fractions, and explicit `null` are rejected.
- Active conflict groups are exact: `world` contains movement, `interact`, and all four UI actions; `combat` contains all four combat and all four UI actions; `menu` contains the four UI actions and is used by lobby/story/summary. One chord cannot bind two actions in the same group. A complete keyboard-only path requires all movement/interact actions, `ui_next`, `ui_previous`, `ui_confirm`, `ui_cancel`, and all four combat actions to retain at least one non-conflicting chord, allowing launch, story/character selection, ready/start, movement/interaction/loot focus, story voting, every combat action, and summary acknowledgement without a pointer.
- Missing settings use the exact object above. Malformed, oversized, unsupported-future, or too-old settings are preserved for diagnosis, ignored with one redacted diagnostic, and replaced in memory by those defaults while the first interactive screen clearly reports recovery; startup must not panic.
- Writes use a same-directory temporary file, file sync, atomic rename, and directory sync where the platform supports it. A failed write leaves the previous valid file intact.
- Migration is pure and idempotent. Version `1` has no previous positive version; once version `2` exists, readers support exactly `2` and `1`. Downgrade never destructively rewrites a newer file.

## Presentation Fallback v1

- Missing or invalid required JSON, TMX, catalog IDs, or archive members are validation failures; the pack/run is not activated. An invalid/expired/revoked Tier 2 attestation returns local `pack_unavailable` and never falls back to loading untrusted bytes.
- Fallback applies only when an asset that already passed validation cannot be decoded/created by the client at presentation time. Sprite, prop, tileset, and backdrop presentation use built-in `builtin:missing_texture`; world collision and authority still use the validated server map. A failed ambience/music/voice clip becomes silence. A failed SFX uses `builtin:default_sfx`. A missing/corrupt built-in skin-token resource uses the complete built-in `classic` variant; an unknown `ui_variant` remains a validation error.
- `builtin:missing_texture`, `builtin:default_sfx`, and the complete `classic` skin are boot-critical embedded package resources verified by checksum before the first interactive screen. If any is absent/corrupt, startup fails with one stable package-integrity error; fallback never recurses.
- Each fallback emits one redacted diagnostic per `(pack_checksum, asset_path, failure_kind)` per run. Fallback never changes hit targets, focus order, controls, collision, timing, authority, or event handling, and it never retries or allocates without the normal asset limits.

## Community Website (Post-MVP)

- Build website as separate repo: `les-perissables-hub`. Keeping it separate preserves the boundary between the ARR game runtime and a public, content-facing site.
- Roll it out in two stages inside that repo:
  - Stage 1 (may ship before the game launches): a simple landing page that points the domain at the project, links the Steam page/wishlist, and links community channels (e.g. Discord). Content-only: no accounts.
  - Stage 2 (implemented soon after the game ships, but not publicly launched until moderation/privacy gates are complete): a community hub where players sign in, share content packs, and discover others' packs.
- Hub purpose: host/discover community data packs (story/character/theme packs), not runtime binaries. Players still need to own the game on Steam to run any pack.
- The hub is a user-generated-content (UGC) social platform: signed-in users can share packs, like them, comment on them, and sort/browse by likes.
- Hub names, comments, tags, attribution, and pack metadata are plain text rendered with contextual escaping; raw HTML, scriptable markup, inline event handlers, and user-controlled template fragments are forbidden. Responses use a restrictive nonce/hash-based CSP, `X-Content-Type-Options: nosniff`, `Referrer-Policy: strict-origin-when-cross-origin`, frame denial, and HSTS after HTTPS is verified.
- Pack downloads use an allowlisted non-executable content type plus `Content-Disposition: attachment` with a server-generated filename. Public downloads are served from a cookie-less R2/CDN origin separate from OAuth/session cookie scope; authenticated private downloads use short-lived single-object URLs and never reflect a user-supplied content type.
- Identity (locked): authentication via Discord and GitHub OAuth only - no homegrown email/password system. Store unique `(provider, provider_subject)` plus display name. Account linking is explicit and reauthenticated; never merge accounts by display name.
- Moderation is a launch requirement, not a later add-on: report/flag flow, admin delete/ban actions, and anti-spam/upload limits ship before the Stage 2 UGC hub is publicly opened.
- Privacy is a launch requirement: publish a privacy policy and support account/content deletion. See "Privacy And Data Minimization" for how the hub's data handling differs from the game.
- Baseline hub bounds before any external account: one active upload/validation per account, four validations globally per instance, presigned upload expiry `15 minutes`, `20` stored packs and `2 GiB` per account, `10` upload attempts per day, comment length `2,000` Unicode scalar values, and pagination capped at `100` items. Phase 20 may tighten these values from measured abuse/operational evidence but may not remove bounds.
- Locked tech stack: Rust `axum` + `maud` (server-rendered HTML) + `htmx` (interactivity), as one app (a single crate, not a workspace) that starts as the Stage 1 landing page and grows into the Stage 2 hub. For upload validation it depends on the MIT pack schema/validation library crate published from `les-perissables-stories`, so hub-side checks match the game exactly without pulling in proprietary game-repo code.
- Locked hosting/data: deploy on Railway; database is Railway Postgres accessed via `sqlx` with migrations. If the database is ever outgrown, switch to PlanetScale (Postgres); Neon is explicitly not used. Toasty ORM was evaluated and deferred until it is post-1.0/stable.
- Locked storage/CDN: uploaded pack/asset files live in Cloudflare R2, accessed from the Railway app through its S3-compatible API; Postgres stores metadata, listing control, checksum, immutable object generation, and the logical-content-addressed R2 key. Use the bounded upload state machine `Reserved -> UploadedPrivate -> Validating -> ValidatedPrivate -> Publishing -> Published`, with `Rejected`/`Expired` terminal cleanup states. Reserve account quota before issuing a `15-minute` presigned upload, verify stored size/checksum, and expire/delete abandoned private objects within `1 hour`; run a daily bounded orphan reconciliation.
- After the Phase 20 public gate, copy validated bytes before committing the `Published` database state; retries are idempotent by upload ID/checksum. Before that gate, non-synthetic objects remain private. Unlisting removes discovery immediately, purges CDN entries, and decrements the object's reference count; delete public bytes only when no listing references that logical checksum and retention/takedown rules permit it. R2 and Postgres are not one transaction, so reconciliation repairs partial copy/database outcomes. Cloudflare provides DNS/CDN in front of Railway. Re-evaluate platform placement from measured requirements rather than assuming Workers would require D1 or a full application rewrite.
- Locked license: the hub is proprietary code and ships All Rights Reserved, like the game (`COPYRIGHT` + ARR `LICENSE` added when the repo is created). This is independent of the MIT schema/tooling crate and the explicitly licensed community packs it serves.

## Community Content And Licensing Boundary

The project has a dual model:

- `perissables` (main game repository; `les-perissables` package prefix): proprietary, All Rights Reserved
- planned `les-perissables-stories` schema, validator, tooling, and bundled creator examples: MIT once that repository is created

Community creators may:

- Create/share story JSON packs
- Create/share character preset JSON packs
- Create/share theme manifests and original presentation assets

Community creators may not:

- Redistribute proprietary game code
- Redistribute proprietary built-in game assets outside the shipped game

Important runtime rule:

- Community packs are data consumed by the proprietary runtime; they do not grant rights to the engine.
- Users still need the game install/ownership to run packs.
- Multiplayer sessions require matching aggregate `pack_id`, pack `version`, canonical SHA-256 checksum, `content_schema_version`, and `game_rules_version` across server and players.
- Pack manifests include those values plus an SPDX `license_expression` and any required per-component license/attribution map.

Creator content tiers (rollout):

- Tier 1 (reuse-only) at first creator release: packs may add new stories, new character stat/spell combinations, and flavor, but must reference already-shipped maps, themes, sprites, audio, and spells. No new asset files. This keeps packs instantly consistent and minimizes moderation/validation load.
- Tier 2 (original assets) later: packs may also ship original tilesets, sprites, maps, audio, combat backdrops, and theme manifests, conforming to the locked asset/TMX conventions above. Theme manifests still select one built-in behavior-neutral `ui_variant`; custom skin-token documents are not schema v1. Tier 2 is enabled only once the hub has submission rules, asset/format validation, and the moderation/abuse controls from the community website plan.
- The schema allows custom asset references and maps from day one, but release builds keep all Tier 2 members disabled unless `PublicationAttestationV1`, its root-signed key set, and the current revocation list verify under the exact attestation contract above. Local developer builds may enable synthetic fixtures only through an explicit non-release feature proven absent from packaged binaries.
- Engine boundary (unchanged by either tier): presentation (art and audio) and narrative (story branching, checks, encounters, character stat/spell composition) are data; UI skin behavior/tokens, combat rules, the d100 system, and spell behaviors remain in the proprietary engine. Creators reskin and re-author the world; they do not change how the game plays.
- Sharing/licensing: the uploading account controls its listing but does not thereby prove copyright ownership. Before upload, terms require the uploader to represent that it has the necessary rights and grant the hub permission to host and distribute the pack. Tier 1 JSON/data uses MIT or CC0; Tier 2 original assets use CC0 or CC BY 4.0 with required attribution. The manifest carries an SPDX `license_expression` for the aggregate plus a required per-path/component license-and-attribution map whenever more than one license or attribution applies. Takedown and repeat-infringer procedures are required before public UGC launch.

## Milestone Exit Criteria

### Vertical Slice Exit

- One canonical map/story scenario completes from start to summary using the recorded acceptance transcript.
- Exhaustive resolver tests cover internal rolls `0..=100`; deterministic fixtures cover branching and critical fallback behavior.
- The canonical combat transcript verifies turn order, all four action types, invalid-action non-mutation, death, win, and wipe behavior.
- no full multiplayer requirement yet
- may use one hardcoded debug character until the Phase 08 data-driven roster exists

### Playable MVP Exit

- Scripted deterministic `2`, `3`, and `4` client scenarios converge on identical state revisions and event IDs.
- Admission boundary tests cover session `64/65`, player `256/257`, a fifth player, simultaneous joins, leave/re-admit, and stable rejection codes.
- The lobby -> run -> summary -> lobby scenario completes three consecutive times without stale state on every supported client platform.
- All three themes load through story metadata with identical controls and authoritative behavior.
- durable run persistence is not required until the Release-ready milestone, even though the persistence model is locked above

### Release-ready Exit

- Deterministic disconnect, reconnect, late-message, duplicate-input, takeover, and resync-ack tests pass without sleeps or unexplained intermittent results.
- The current save version and the immediately previous positive version when one exists restore exact authoritative state after injected failures at every write/rename checkpoint; corrupt/future/oversized snapshots fail safely.
- Hosted staging smoke tests cover create, join, rejoin, deploy/restart restore, graceful shutdown, and rollback on both target OS clients.
- All locked client/server budgets in `docs/benchmark-plan.md` pass on release artifacts with no unexpected errors or disconnects.
- `docs/qa-plan.md` reports `PASS`, and blocker/critical defects are zero.
- Reviewed end-user terms, contribution policy decision, third-party notices, and community-content terms are present where applicable.

## Change Control

- Any feature outside this contract is added to post-MVP backlog.
- Scope changes only happen between phases, never inside an active phase.
- If scope grows, timeline updates must be acknowledged before coding continues.
- Business-facing release notes and pricing decisions stay out of these docs; track them separately when needed.
