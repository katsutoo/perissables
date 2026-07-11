# Les Périssables MVP Contract

Status: Locked for Phase 01 implementation; later phase changes follow the tracker and definition of done
Owner: Project team
Updated: 2026-07-11

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
10. Complete keyboard-only operation and at least two behavior-identical UI skin variants selected by theme data.

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
- An active session is any non-`Ended` session retained in memory or durable storage, including `Lobby`, `Running`, `Summary`, paused, and restored sessions. It counts from successful creation until the durable `Ended` transition and cleanup complete.
- A player is an occupied seat, whether connected or reserved during rejoin grace. A seat counts from committed admission until explicit leave in `Lobby`, session end, or grace expiry. A socket is not a player, and reconnect replacement never increments either capacity counter.
- Admission capacity is the remaining active-session and occupied-seat budget. Creation atomically reserves one session and one seat; join atomically reserves one seat. Rejoin uses an existing reservation and is allowed when new-admission capacity is zero. Concurrent operations are serialized by the admission owner and never oversubscribe a limit.
- Target rate: `60 FPS`
- Client fixed-timestep update rate: `60 Hz`
- The client processes at most `5` fixed updates in one rendered frame and caps accumulated unprocessed time at `250 ms`. Excess time is discarded once, increments `client_update_time_dropped_total`, and emits one rate-limited diagnostic; it must not cause an unbounded catch-up loop. Tests inject elapsed time at the exact cap and one update beyond it.
- Server simulation tick rate: `20 Hz`
- Client network input rate: at most one coalesced continuous-input message per server tick (`20 msg/s`). The `60 Hz` client update loop must not emit one network message per update. Discrete actions share the same bounded rate budget.
- Server state broadcast rate: `20 Hz` max. State broadcasts are latest-value coalesced. A server that falls behind may run at most one immediate catch-up tick, then records a missed-tick metric and resumes the normal schedule; it must not enter an unbounded catch-up loop.
- Hosted-server processing-latency target: from completion of protocol/authentication/rate/sequence validation for an expected-sequence input to handing its correlated `input_result` to the WebSocket writer after any required durable commit, p95 under `100 ms` and p99 under `200 ms`, with `64` active sessions of `4` players under the locked workload in `docs/benchmark-plan.md`. Internet RTT is reported separately and is not part of this gate.
- Tile size: `16x16`
- Primary resolutions for MVP QA: `1280x720`, `1920x1080`
- Playtest targets: median clean-account lobby-to-run start at or below `3 minutes`; median completed or failed run between `20` and `35 minutes`. Phase 17 uses at least `20` consented first-time participants and `30` completed-or-failed sessions, reports the median and bootstrap `95%` confidence interval, and does not reuse a participant's first-run onboarding measurement. Lobby timing starts when the authenticated lobby first becomes interactive and ends when the first authoritative `Running` state is rendered. Run timing starts at that transition and ends at the first authoritative `Summary`; abandoned sessions are reported separately and never silently excluded. These are balance targets, not telemetry collected by default.

## Locked Persistence Model (Release-ready v1)

This model is locked now, but implementation lands in Phase 16 and is part of the Release-ready milestone, not the Playable MVP milestone. Before Phase 16, restarts and deploys may destroy active runs and must be treated as known pre-release behavior.

- Run state is server-authoritative and persists server-side: the game server snapshots active sessions to durable storage (Railway volume) so restarts and deploys do not destroy runs.
- Session snapshots persist only `HMAC-SHA-256` of each rejoin token plus its key ID, generation, player/session binding, and absolute expiry. Compare digests in constant time. Raw bearer tokens exist only on the client and transiently while validating a presented token; they are never logged or persisted. The HMAC key is a managed server secret outside snapshots; retain the current and previous key until all snapshots/tokens using the previous key expire.
- Clients never submit saved run state; resume always happens through the rejoin flow into a server-restored session.
- Longer-term resume (the whole party returning after rejoin tokens expire) is post-MVP; if added, re-admittance is by validated Steam identity, never by extending token lifetime.
- `PersistedSessionSnapshot` is a server-only representation capped at `256 KiB`. It includes authoritative gameplay state, logical tick, state revision, stable RNG algorithm/version/state, pack identity, token digests, absolute `created_at_unix_ms`/`expires_at_unix_ms`/`last_observed_unix_ms`, and remaining monotonic gameplay-deadline durations. It is distinct from the client-visible `ClientResyncState`.
- `save_version` starts at `1`. A snapshot stores the current version after migration; readers support exactly versions `1` and the immediately previous positive version once version `2` exists. Future, zero, and older versions fail closed without partial restoration.
- After Phase 16, externally acknowledged authoritative transitions are group-committed to one bounded server WAL before their `state`, `event`, or `input_result` is emitted. The server batches all sessions for at most `20 ms`, appends length/checksum-framed records, `fsync`s once, then publishes acknowledgements. The WAL is capped at `64 MiB`; reaching the cap stops state-changing admission until compaction succeeds. Before Phase 16, restart rollback/loss remains explicit pre-release behavior and no durability claim applies.
- Dirty snapshots compact the WAL at most once per session per second. Run start, encounter resolution, run end, initial token issuance after Phase 16, and token rotation are durable checkpoints; a token is not returned until its digest/generation record is committed.
- Snapshot generations are immutable files named by monotonically increasing generation. Write a same-directory temporary file, `fsync` it, rename to its final generation, then `fsync` the directory. On startup, scan newest to oldest and select the newest checksum-valid generation whose WAL replay validates; retain the newest two valid generations and delete older ones only after the current generation is durable. At most one snapshot write and one pending dirty state exist per session.
- Each state-changing operation is a transaction: validate against committed state; clone or journal only the bounded fields it can change; stage gameplay state, RNG position, consumed `input_seq`, revision, and outputs; append and commit one WAL record; then atomically replace committed in-memory state and enqueue outputs. A failed append or `fsync` discards the entire staged transition, does not consume `input_seq`, does not increment `state_revision`, emits no semantic event/state/result, and returns correlated recoverable `persistence_unavailable` with the unchanged `next_input_seq`. Retrying that same sequence is safe. No observer may read staged state.
- Persistence retries are bounded to three attempts with delays of `100 ms`, `500 ms`, and `2 s`. While retrying, the owning session accepts heartbeats/resync but rejects new state-changing work with correlated `persistence_unavailable`. After the third failure it enters `PersistenceBlocked` for at most `30s`; one successful probe and commit restores service. Otherwise the owner commits no further mutation, sends a `session_terminated` notice if the writer is available, closes with `1011`, and retains the last valid durable generation for operator recovery.
- Every immutable snapshot generation records the inclusive WAL record offset it contains. Global compaction may replace the WAL only after every non-ended session has a checksum-valid, directory-synced snapshot covering a common offset; the low-watermark is the minimum covered offset across those sessions. Copy records after that offset into a new same-directory WAL, `fsync`, rename, and `fsync` the directory before deleting the old WAL. Ended-session tombstones remain until covered by the same low-watermark.
- Client save paths hold only non-authoritative local data (settings, keybinds, local preferences):
  - Linux: `~/.local/share/les-perissables/`
  - Windows: `%AppData%/LesPerissables/`

### Authoritative RNG v1

- `rng_version: 1` is the IETF ChaCha20 stream construction with a per-session `256`-bit key, `96`-bit nonce, `32`-bit block counter, and little-endian `u32` words. The key is generated once from the operating-system CSPRNG; the nonce is independently generated from the same source. Neither value is supplied by a client or pack, and neither is exposed in client projections or logs.
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
- Single-instance enforcement: acquire an exclusive OS file lock on the mounted volume plus a persisted monotonically increasing fencing generation before restoration or admission. The lease records owner UUID and renewal time, renews every `5s`, expires after `20s`, and may be recovered only after expiry and successful lock acquisition. Every WAL/snapshot write includes and verifies the current fencing generation. `/readyz` remains unsuccessful until the lease is held and restoration is complete.
- Game-server config: clients receive an exact environment-specific production WebSocket origin allowlist from trusted build/environment config. Steam lobby metadata must not introduce an arbitrary endpoint.
- Session discovery: Steam lobbies/invites are discovery only; the authoritative server owns `session_id`, player IDs, game state, dice, combat, and story progression
- Steam lobby mapping: lobby metadata stores `session_id`, aggregate `pack_id`, pack `version`, pack checksum, `protocol_version`, `content_schema_version`, and `game_rules_version`. The endpoint comes only from the trusted allowlist.
- Production identity: production `join`/`rejoin` requires Steam auth/session-ticket validation before the server issues or accepts player/session credentials
- Operations baseline: expose `/healthz` for process liveness and `/readyz` for service readiness. Readiness requires the exclusive lease, completed restoration, writable persistence, and ability to serve existing sessions; exhausted session/player admission capacity does not make the service unready. New creation/join receives `server_full`, while reserved-seat rejoin remains available. Use structured logs and monitor deploy status, missed ticks, queue saturation, error rate, disconnect rate, reconnect failures, and admission capacity separately.
- Graceful shutdown: mark unready, stop new creation/join, send the `server_restarting` notice with `retry_after_ms`, checkpoint dirty sessions, await owned tasks, and exit within `20s`. Connections close with WebSocket `1012` only after their notice is handed to the writer. If the deadline is reached, exit non-zero and rely on the last durable checkpoint.
- Map format for MVP: TMX
- Release automation: GoReleaser for build/package generation only; it does not change licensing or grant public binary distribution rights
- Distribution policy: production desktop binaries ship through Steam depots, not public release pages
- Release targets: Linux + Windows only
- macOS policy: deferred (requires Apple Developer Program for signing/notarization workflow)
- Protocol versioning: `protocol_version` integer, start at `1`; every increment is breaking unless a future compatibility document says otherwise.
- Content schema versioning: `content_schema_version` integer, start at `1`; every increment is breaking. Reject unsupported values for story, character, theme, and pack manifests.
- Save schema versioning: `save_version` integer for server-side session snapshots, support the current version and the immediately previous positive version once one exists, with explicit migrators
- Client-settings versioning: `client_settings_version` integer starts at `1`; support the current version and the immediately previous positive version once one exists, with explicit pure migrators
- Game-rule versioning: `game_rules_version` integer, start at `1`; it changes when authoritative mechanics or built-in effect semantics become network-incompatible.

## Locked Protocol And Limits

### WebSocket Envelope v1

All messages must use:

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
- `resync_request`
- `resync_state`
- `resync_ack`
- `input_result`

Message directions for v1:

- Client -> server: `join`, `input`, `pong`, `rejoin`, `resync_request`, `resync_ack`
- Server -> client: `join_response`, `state`, `event`, `notice`, `error`, `ping`, `resync_state`, `input_result`

Sequence rule:

- `seq` is an unsigned `64`-bit per-connection, per-direction transport counter and starts at `1`. It provides ordering only within one connection; it is not gameplay replay protection. A connection closes before incrementing past `u64::MAX`; wrapping is forbidden.

Join payload v1:

```json
{
  "mode": "join",
  "requested_session_id": "ses_...",
  "steam_auth_ticket": "base64_ticket",
  "pack_id": "cleanup-pack",
  "pack_version": "1.0.0",
  "pack_checksum": "sha256:...",
  "content_schema_version": 1,
  "game_rules_version": 1
}
```

`mode` is `create` or `join`. For `create`, `requested_session_id` must be absent. For `join`, it is required and comes from Steam lobby metadata. The Steam ticket is required in production, sent only over WSS to an allowlisted origin, validated for the expected app and ownership before admission, and never logged or persisted. Pack and version fields are checked against the server's own validated aggregate pack; client values never establish authority.

Before rate-limit state allocation, authentication concurrency acquisition, or provider work, join fields obey these lexical bounds: `steam_auth_ticket` is unpadded RFC 4648 base64url ASCII decoding to `1..=4,096` bytes; `requested_session_id` uses the server-ID grammar and cap; `pack_id` uses `ContentId`; `pack_version` is `1..=64` printable ASCII bytes and must parse as the locked SemVer subset; `pack_checksum` is exactly `sha256:` plus `64` lowercase hex digits; schema/rules versions are JSON integers `1..=u32::MAX`. Invalid base64, padding, non-ASCII, decoded oversize, or lexical mismatch rejects before Steam validation.

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

The `join_response` envelope carries the issued `session_id` and `player_id`. `server_time_ms` is milliseconds since Unix epoch sampled when the response is created and is presentation-only; authoritative deadlines use server monotonic time and are sent as non-negative `remaining_ms`, so wall-clock adjustment cannot change gameplay. The `rejoin_token` is opaque, server-generated with a CSPRNG at `128` bits of entropy minimum, never stored in Steam lobby metadata, logged, or persisted raw. It stays valid for the active session plus the `rejoin_grace_seconds` window after disconnect, is rotated on successful `rejoin`, and is invalidated when the session ends.

Rejoin payload v1:

```json
{
  "rejoin_token": "opaque_128_bit_minimum_random_token",
  "steam_auth_ticket": "base64_ticket"
}
```

For `rejoin`, `session_id` and `player_id` must be non-empty. The server binds the validated Steam identity, session, player, token generation, and new connection. One player may have only one active connection; successful rejoin atomically replaces and closes the old connection. The server sends `resync_state` with a monotonic `state_revision`; the client must return `resync_ack` for that revision before new `input` is accepted.

`rejoin_token` is unpadded base64url ASCII representing `16..=32` decoded bytes and is capped at `64` wire bytes. Rejoin `steam_auth_ticket` uses the same rules as join. Token syntax/length, envelope IDs, and ticket decoding are checked before keyed-digest comparison or Steam validation.

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

`outcome` is `applied`, `superseded`, or `rejected`. `reason_code` is `null` for applied input and otherwise one of `superseded`, `stale_revision`, `invalid_phase`, `invalid_action`, `not_owner`, `not_ready`, `duplicate_character`, `target_unavailable`, `unaffordable`, `out_of_turn`, `queue_full`, or `limit_reached`. A semantic rejection uses only `input_result`; `error` is reserved for envelope/auth/rate/sequence/service failures. When an error concerns an input, `input_seq` and `next_input_seq` are mandatory. `request_seq` is the inbound envelope sequence being answered.

`notice.kind` is `server_restarting` or `session_terminated`. `retry_after_ms` is mandatory only for restart and forbidden for termination. For termination, `reason_code` is mandatory and is `persistence_failed`, `clock_invalid`, `expired`, or `internal`; it is forbidden for restart. A notice is an out-of-band server control message, does not consume an input sequence or increment revision, and is ordered after all already-committed output bundles. Authenticated notices use the connection's session/player IDs. A pre-auth connection receives no notice: restart closes it with `1012`; unrecoverable shutdown closes it with `1011`. An authenticated restart sends `server_restarting` then closes `1012`; unrecoverable session termination sends `session_terminated` with a public reason then closes `1011`.

Malformed JSON, unknown fields/types, wrong direction, or authorization mismatch returns a stable error when safe and closes with WebSocket code `1008`. Oversized messages close with `1009`; service restart uses `1012`; temporary global overload uses `1013`; unrecoverable session failure uses `1011`. Public errors never include parser internals, tickets, tokens, digests, filesystem paths, or hidden state.

Input action v1 tags and fields:

| Tag | Fields | Valid phase | Queue class |
| --- | --- | --- | --- |
| `select_story` | `story_id: ContentId` | `Lobby`, owner only | discrete |
| `ready` | `character_id: ContentId` | `Lobby` | discrete |
| `unready` | none | `Lobby` | discrete |
| `start_run` | none | `Lobby`, owner only | discrete |
| `move` | `direction` (`up`, `down`, `left`, `right`, `none`) | `Running/world` | continuous/latest-value |
| `interact` | `target_id: EntityId` | `Running/world` | discrete |
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
| `ClientView` | `phase`, `self: PlayerView`, `party: PlayerView[1..=4]`, `available_actions: ActionKind[0..=16]`, and exactly one phase detail |
| `PlayerView` | `player_id`, `display_name` (`1..=64` scalars), `connected: bool`, optional `character_id`; `ready: bool` is mandatory only in lobby, while `alive`, `hp: u16`, and `max_hp: u16` are mandatory only after a character is locked for a run |
| `LobbyView` | optional `selected_story: StoryPreview`, `stories: StoryPreview[1..=64]`, `owner_player_id`, `minimum_players: 2` |
| `StoryPreview` | `story_id`, `title` (`1..=128` scalars), `estimated_minutes: u8`, `theme_id` |
| `WorldView` | `map_id`, `self_position: Position`, `actors: VisibleActor[0..=128]`, `active_interaction_ids: EntityId[0..=32]` |
| `Position` | `x_subtiles: i32`, `y_subtiles: i32`, `facing`; one tile is exactly `256` subtiles |
| `VisibleActor` | `actor_id`, `position`, `sprite_id`, `alive: bool` |
| `StoryView` | `node_id`, `public_flags: ContentId[0..=256]`, `choices: ChoiceView[0..=8]`, `remaining_ms: u32` |
| `ChoiceView` | `choice_id`, `label` (`1..=256` scalars), `votes: u8` |
| `CombatView` | `actors: CombatActorView[1..=20]`, `turn_actor_id`, `round: u16`, `available_target_ids: EntityId[0..=20]`, `remaining_ms: u32` |
| `CombatActorView` | `actor_id`, `display_name`, `side`, `alive`, `hp`, `max_hp`, `resource: u8`, `resource_max: u8` |
| `SummaryView` | `result` (`completed` or `failed`), `deaths: EntityId[0..=4]`, `key_event_ids: EventId[0..=32]`, `acked_player_ids: PlayerId[0..=4]`, `remaining_ms: u32` |

`ClientView` contains exactly one detail field named for its `phase`: `lobby`, `world`, `story`, `combat`, or `summary`. `Lobby` and `Summary` session states map directly, while the `Running` state exposes one of the three running phases. `side` is `party` or `enemy`; `facing` is `up`, `down`, `left`, or `right`. Numeric fields use JSON integers and reject fractions. Optional means absent or a value, never explicit `null`, except `input_result.reason_code`.

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

`EventItem` has `event_id: EventId`, `kind`, and a kind-specific payload. Kinds are `audio_cue { cue_id, channel }`, `dice_result { actor_id, roll, stat, result }`, `story_transition { from_node_id, to_node_id }`, `combat_transition { from_phase, to_phase }`, `actor_death { actor_id }`, `run_result { result }`, and `time_dropped { ticks }`. One committed session batch emits at most one `event` envelope per connection containing at most `64` items and at most `48 KiB`; repeated low-value audio cues are coalesced by `(cue_id, actor_id)` and overflow of required semantic events rejects and rolls back the staged batch with `limit_reached` rather than disconnecting a valid client.

Stable error codes v1 are `invalid_message`, `unauthorized`, `version_mismatch`, `pack_mismatch`, `session_not_found`, `session_full`, `server_full`, `rate_limited`, `stale_input`, `input_gap`, `invalid_action`, `busy`, `persistence_unavailable`, and `internal`.

Revision, replay, and output rules:

- `state_revision` starts at `1` for committed session creation and increments exactly once for each committed session batch that changes gameplay/lobby state or consumes one or more expected input sequences. A batch that only answers ping/resync does not increment it. It never wraps.
- Every `input` carries `input_seq` and `based_on_revision`. `input_seq` starts at `1`, is monotonic per player/session, survives reconnect, is persisted, and never wraps. `resync_state` returns the next expected value.
- `based_on_revision` must not be in the future and may trail the committed revision by at most `64`. Such an input is validated against current committed state, never historical state. A future or older revision is committed as `rejected/stale_revision`: it consumes that sequence, increments revision, and emits no gameplay event. The result tells the client to request resync. Duplicate/lower values return `stale_input` without consuming; gaps return `input_gap` followed by `resync_state`; the server never guesses missing intent.
- A syntactically valid, authenticated, rate-admitted expected input consumes its sequence when semantically rejected or superseded. Parse, authentication, rate, gap, duplicate, and persistence failures do not consume it.
- On each simulation tick, the owner forms one session-wide batch containing at most one selected input per connected player plus ingress-generated superseded/queue-full results. It stages selected inputs in lexical `player_id` order against the preceding staged result, then commits all gameplay/RNG changes, consumed sequences, one revision, and all outputs in one WAL record. Any persistence or required-event-bound failure rolls back the whole batch; semantic rejection of one input does not roll back other valid inputs.
- For one committed batch, the owner assigns event IDs, then hands each writer an indivisible ordered output bundle: optional single `event` envelope, optional single `state`, then that player's `0..=8` individual `input_result` envelopes in ascending `input_seq`. Each result envelope is capped at `512` bytes and all results carry the batch revision. The writer assigns transport `seq` in that order. No later revision may be handed to any writer before the earlier bundle. Resync state supersedes queued coalescible state but never a result or semantic event.

Session lifecycle v1:

- States are `Lobby`, `Running`, `Summary`, and `Ended`; only the session owner task mutates them. The creating player is owner until it leaves the lobby; ownership then transfers to the lowest lexical connected `player_id`. Ownership cannot change during `Running` or `Summary`.
- One validated Steam identity may occupy one active seat and own one lobby. Per identity, session creation is limited to `2` attempts/minute (burst `2`) and join/rejoin to `10` attempts/minute (burst `10`).
- In `Lobby`, the owner selects one available story. Changing it clears every ready selection. Each connected player may ready one unique character or unready. Ready commands collected for the same simulation tick are resolved in lexical `player_id` order, so the lowest ID wins a duplicate pick; later picks reject with `duplicate_character`.
- `start_run` applies only when a story is selected, `2..=4` seats are connected, and every connected player is ready with a unique character. Only the owner may start. It commits `Lobby -> Running/story` at the selected story's `start_node`, captures the selected story/roster, clears lobby-only acknowledgements, and initializes run state and RNG consumption at zero. A `return` node enters `Running/world`; later accepted TMX trigger interactions re-enter `Running/story` at their mapped node.
- Run completion or party wipe commits `Running -> Summary`. Connected players may `summary_ack`; after every connected player acknowledges or `60s` monotonic summary time elapses, the owner commits `Summary -> Lobby`, clears selected story, all ready/character selections, votes, run state, deadlines, event-dedup presentation state, and run RNG state, while preserving seats, ownership, input sequences, and allowed local client settings. This supports repeated runs without reconnecting.
- A session has an absolute lifetime of `4 hours`, including reconnect grace. Persist `created_at_unix_ms`, `expires_at_unix_ms`, and `last_observed_unix_ms` from the authenticated server wall clock; use monotonic elapsed time between observations in one process. On restore, wall time at or after expiry ends the session and counts downtime. If wall time is earlier than persisted `last_observed_unix_ms` by more than `1s`, fail closed by ending the session with `clock_invalid` rather than extending credentials; smaller rollback clamps to the persisted value. Forward jumps may expire sessions early and are never undone. Accepted lobby actions, accepted gameplay input, rejoin, and resync acknowledgement count as activity. A connected lobby/run/summary with no such activity for `15 minutes` ends. An empty lobby expires after `10 minutes`.
- A running or summary session with no connected players pauses simulation ticks, combat/story/summary deadlines, and idle time. It expires at the earlier of the `600-second` rejoin grace or absolute session expiry. Rejoin resumes each deadline with its stored non-negative remaining monotonic duration. Absolute expiry never pauses.
- Token absolute expiry is session creation time plus `4 hours`; disconnect never extends it. Tokens rotate on successful rejoin but retain that cap. Expiry commits `Ended`, invalidates tokens, releases capacity, and removes persisted state through a durable tombstone after Phase 16.
- A disconnected seat remains reserved during grace. A disconnected current combat actor is treated as `pass` after its paused/resumed `30-second` turn timeout. Story votes close after `60s`; the choice with the most votes wins. Ties use the vote of the lowest lexical `player_id` among voters for tied choices, and no votes selects the first choice in schema order.
- The session owner selects at most one gameplay input per connected player per simulation tick and stages them in lexical `player_id` order within the one session batch. A discrete action takes priority over continuous movement for that player's tick. Continuous input is latest-value coalesced; each superseded sequence receives `input_result`. Each player has a bounded queue of `8` discrete actions. Overflow consumes and rejects the newest expected sequence with `queue_full`. `input_seq` establishes per-player order, not cross-player priority.

Movement and logical tick v1:

- Authoritative positions are signed fixed-point `Position` values in subtiles; one `16x16` tile equals `256` subtiles on each axis. A living player moves `64` subtiles per `20 Hz` server tick while its held direction is not `none`, for a base speed of `5` tiles/second. Diagonal movement does not exist in v1.
- A `move` input replaces the held direction and remains active until another `move`, phase exit, disconnect, death, or blocked run transition sets it to `none`. The owner advances movement from an explicit `Tick { logical_tick }` command generated by the monotonic scheduler; ticks are state-machine inputs, not wall-clock reads inside `game_core`.
- For each tick, players are considered in lexical `player_id` order. Resolve X/Y as one cardinal displacement, then test map bounds, static blocked tiles, and finally the committed positions of actors already resolved for that tick. A blocked move leaves position unchanged but updates facing. Players not yet resolved retain their prior position for collision, so two players cannot swap or occupy one location in a tick.
- Actor collision uses one axis-aligned `192x192`-subtile box centered on the position. Map edges are blocked; positions and collision boxes must remain within validated map dimensions. Warps and interactions execute after all movement for the tick in lexical player order.
- A server delayed by more than one tick performs one immediate catch-up tick, records all additionally dropped logical ticks in `server_tick_time_dropped_total`, emits one bounded `time_dropped` diagnostic event, and resumes the next scheduled tick. Dropped ticks do not advance logical time, movement, or deadlines.

### Network Limits Defaults

- Max inbound WS frame size: `16 KiB`
- Max outbound WS frame size: `64 KiB`
- Max reassembled inbound message size: `16 KiB`; WebSocket compression and binary gameplay messages are disabled for v1.
- Input rate limit: `20 msg/s` (burst `4`) per client
- Control rate limit: `5 msg/s` (burst `10`) per client
- Rate limits use token buckets. Every received attempt consumes capacity before payload validation. Continuous client input is latest-value coalesced to at most `20 msg/s`; a client must not send at its `60 Hz` update rate.
- Heartbeat: the server initiates application `ping` every `10s`; the client answers with `pong`. Disconnect after `30s` without a valid pong.
- Control messages are `join`, `pong`, `rejoin`, `resync_request`, and `resync_ack`. Every state-changing lobby/gameplay `input` uses the input rate limit; server messages are bounded by message-size and broadcast-rate limits.
- Before parsing a Steam ticket or calling Steam, `join`/`rejoin` use per-source-IP `20` attempts/minute (burst `5`) and per-requested-session `30` attempts/minute (burst `10`) buckets. The global bucket table holds at most `4,096` entries, expires idle entries after `10 minutes`, evicts the least-recently-used idle entry only, and rejects with `rate_limited` when full with no idle entry. Unauthenticated failures never consume an authenticated gameplay budget.
- Steam validation has a global concurrency limit of `32`, a `3s` total timeout, and no automatic retry on one request. Exhaustion returns recoverable `busy` with `retry_after_ms` between `250` and `1,000`; timeout returns `unauthorized` without disclosing provider detail. A circuit opens for `15s` after `20` consecutive provider failures and permits one probe at a time while half-open.
- Source IP for throttling comes from the direct peer unless the connection came through an explicitly configured trusted proxy. Forwarding headers from any other peer are ignored; ambiguous forwarding chains fail closed.
- `ClientResyncState` is a per-player projection capped at one `64 KiB` outbound message. It excludes token digests, RNG state, other server-only data, and information hidden from that player. Every reachable validated session state must produce a projection within this cap; protocol fragmentation is out of scope for v1.
- Max simultaneous unauthenticated handshakes: `64`; max total sockets: `320`; handshake timeout: `10s`.
- Each session has one owner task and a bounded mailbox of `256` commands. Each connection has an outbound queue of `8` ordered bundles, not individual event items; each bundle contains at most one `48 KiB` event batch, one coalescible `64 KiB` state, and `0..=8` requester-only `512`-byte results, for at most `116 KiB` encoded and `928 KiB` queued per connection. The input burst, one-result-per-sequence rule, one-event-batch-per-committed-batch rule, and one selected input/player/tick keep a conforming writer within this bound under the benchmark workload. If a bundle still cannot be queued because the peer is not draining, disconnect that slow client and require resync rather than growing memory.
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
- Max flags per story: `256`; max encounters per story: `128`; max automatic transitions handled for one intent: `64`.
- Max spells per character: `8`; max characters per aggregate pack: `64`.
- Max locally installed community packs: `32`; max aggregate installed-pack disk use: `4 GiB`. Accounting includes retained archives, extracted generations, metadata, and previous generations awaiting deletion. Import temporary files have a separate process-wide `256 MiB` cap; before extraction, require free space for current accounted use plus the declared final and temporary maxima. Import refuses before crossing a limit, removes temporary data on every outcome, and never evicts without user action.
- Max TMX file size: `4 MiB`.
- Max TMX map dimensions: `512x512` tiles; max `16` tilesets, `16` layers, `2` object layers, `1,024` objects, `4,096` properties, and `100,000` total XML elements per aggregate pack. Max XML depth is `32`, attributes per element `64`, attribute value `1 KiB`, text node `1 MiB`, and total decoded tile-layer data `2 MiB` per map.
- TMX tile layers use CSV encoding only in v1. Base64, gzip, zlib, zstd, infinite/chunked maps, templates, image layers, scripts, and parser-driven URL/file resolution are rejected. A TMX may reference an in-pack relative TSX or image path only after portable path normalization and membership validation; the loader opens that already-validated archive member directly and never asks the XML parser to resolve it.
- Image files are capped at `8 MiB`, `4096x4096` pixels, and `64 MiB` decoded bytes each. Aggregate decoded image memory during validation is capped at `128 MiB`.
- Audio files are capped at `16 MiB`, `10 minutes`, `2` channels, `48 kHz`, and `64 MiB` decoded bytes each.
- Client runtime retained decoded limits per active pack are `128 MiB` images plus `128 MiB` audio, with `256 MiB` total assets. During a theme/pack swap, process-wide decoded assets are capped at `512 MiB` and loader concurrency at `2`. Inspect headers and computed decoded sizes before allocation; refuse the load deterministically rather than implicitly evicting active assets. The headless server never decodes presentation media.
- TMX/XML parsing disables DTDs, external entities, XInclude, and all external resource resolution. Untrusted image/audio validation and decode run in a sandboxed worker subprocess with no network, a `2s` per-file wall timeout, process memory cap `256 MiB`, concurrency `2`, and declared decoded-output caps even when headers lie; timeout or limit kills the worker before its slot is reused. Validation and runtime loading enforce the same limits; community packs remain hostile after hub validation.

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
| `includes` | `1..=128` unique strings tagged `story:`, `character:`, or `theme:` followed by `ContentId`; every listed document must exist and unlisted content documents are rejected |
| `publication_attestation` | optional Tier 2 object defined under the content-tier rules; forbidden for unsigned Tier 1 authoring output |

Story document v1 (`stories/<story_id>.json`):

| Field | Type and rule |
| --- | --- |
| `content_schema_version` | integer, exactly `1` |
| `id` | `ContentId`, equal to filename stem |
| `title` | `LocalizedText`, max `128` scalars |
| `theme_id`, `map_id`, `start_node` | `ContentId`; references must resolve |
| `estimated_minutes` | integer `1..=120` |
| `nodes` | `1..=512` unique `StoryNode` objects |
| `encounters` | `0..=128` unique `Encounter` objects |
| `triggers` | `1..=128` unique `StoryTrigger` objects |

A `StoryTrigger` has `event_id: ContentId`, `node_id: ContentId`, and `once: bool`. Every TMX `events`-layer object `event_id` must map to exactly one trigger and every trigger must map to one object and reachable node. An accepted world `interact` for that server-visible event enters `story` at `node_id`; a once-only trigger is marked consumed in the same staged transaction.

`StoryNode` has required `id: ContentId`, `kind`, and `effects: Effect[0..=16]`. Its remaining required fields depend on `kind`:

| Node kind | Required fields |
| --- | --- |
| `dialogue` | `speaker_id: ContentId`, `text: LocalizedText`, `choices: Choice[1..=8]`; each choice has unique `id`, `label`, `goto`, and optional `requires_flag` |
| `check` | `label: LocalizedText`, `stat` (`strength`, `agility`, `wit`, `charm`), `outcomes`; outcomes require `success` and `failure` and optionally contain `critical_success`/`critical_failure`, each a `Branch` |
| `encounter` | `encounter_id: ContentId`, `on_win: ContentId`; a party wipe bypasses story continuation and enters failed summary |
| `transition` | `goto: ContentId`; executes automatically after common effects and counts toward the automatic-transition limit |
| `return` | no additional fields; commits the current story sequence back to `world` at the unchanged authoritative positions |
| `end` | `result` (`completed` or `failed`) |

A `Branch` has `goto: ContentId` and `effects: Effect[0..=16]`. An `Effect` is one of `set_flag { flag_id }`, `unset_flag { flag_id }`, `give_item { item_id, recipient }`, or `remove_item { item_id, owner }`; `recipient`/`owner` is `actor`, `lowest_hp_living`, or `lowest_player_id_living`. Flags referenced anywhere count toward the story's `256`-flag limit. Every `goto`, trigger, encounter, item, speaker, map, theme, and built-in behavior reference must resolve during validation. The graph roots are `start_node` and every trigger `node_id`. From every root, each path must reach `return`, `end`, an unresolved player choice, check, or encounter within `64` automatic transitions; every node must be reachable from at least one root, and unbounded automatic cycles are rejected.

An `Encounter` has `id`, `enemies: Enemy[1..=16]`, and `rewards: Reward[0..=16]`. `Enemy` has unique `id`, `name`, four stats in `5..=70`, `max_hp: 1..=999`, `resource_max: 0..=100`, and `spell_ids: ContentId[0..=8]` from the built-in catalog. `Reward` has `item_id` and `quantity: 1..=4`; total reward items are capped at `16` before inventory placement.

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

Schema validation occurs in this order: archive/path/resource bounds, JSON syntax/duplicate keys, structural field/type/bounds, identifier uniqueness, cross-reference resolution, graph/automatic-transition validation, license/attestation policy, then canonical checksum. ZIP/path/checksum and JSON/TMX/XML parsing/semantic validation run only in sandboxed no-network subprocess workers, with at most `4` aggregate workers and `2` media workers. Each stage has a `5s` wall deadline and one pack has a `15s` aggregate deadline, `512 MiB` subprocess RSS cap, and `256 MiB` temporary-disk cap. Deadline, cancellation, worker crash, or limit kills and reaps the subprocess, removes temporary output, returns `validation_resource_limit`, and installs nothing; parser work never runs in a non-killable in-process thread. A failure returns one bounded path-based diagnostic and performs no media decode or network access before cheaper structural checks pass.

### TMX

- Required layers: `ground`, `collision`, `events`
- Optional layer: `decor`
- Object naming: `evt_<id>`, `spawn_<id>`, `warp_<id>`
- Required event object property: `event_id`

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
- If an effect, reference, checked arithmetic operation, event bound, or the `64`-automatic-transition limit fails, discard the complete staged gameplay/RNG mutation. Commit only consumed replay metadata and return `rejected/limit_reached` or `rejected/invalid_action`; partial flags, damage, rewards, movement, or events are forbidden.
- Story votes do not consume RNG. On vote closure, resolve the plurality rule, apply the selected branch transaction once, and clear all votes before the next node becomes visible.
- Each living character has `4` inventory slots. A corpse exposes at most `2` eligible carried items; fewer are exposed when fewer exist. Item transfer cannot exceed the receiver's free slots.
- HP/max HP: `0..=999`; action resource pools: `0..=100`; individual costs: `0..=100`; individual damage/healing effects: `0..=999`; inventory slot index: `0..=3`.
- Encounters contain at most `16` enemies and `20` total actors. Combat is capped at `256` rounds and `5,120` actor turns; reaching either cap ends the encounter as a failed run with a stable limit event.
- Arithmetic uses checked operations in a wider intermediate type, then clamps HP/resources to their validated maxima. Overflow, underflow, or an out-of-range authored value is a validation/programmer error, never wrapping behavior.
- Combat turn order is descending agility. Ties are resolved by stable lexical actor ID. Each actor performs exactly one of `attack`, `spell`, `item`, or `pass` per turn.
- Invalid, unaffordable, dead-actor, or out-of-turn actions return `input_result` with `outcome: rejected` and the applicable stable `reason_code`; they never return an `error` envelope for the semantic rejection. They do not mutate gameplay state or consume the turn. Their expected `input_seq` still advances replay metadata and is durably recorded after Phase 16.
- An actor at zero HP is dead and removed from future turns. Combat ends when all enemies are dead or all player characters are dead; a party wipe produces the failed-run summary.
- `attack` checks the attacker's `strength` using the locked d100 resolver. Normal success deals `max(1, floor(strength / 5))` damage; critical success deals twice that value; failure and critical failure deal zero. Critical failure has no additional self-damage in v1. Damage is applied after the dice event and before death events.
- A spell is immutable built-in engine data with exact `id`, `check_stat`, `resource_cost`, `target_side`, `effect_kind`, `normal_amount`, and `critical_amount`. The v1 catalog is checked in as `game_rules_v1.json`, canonicalized with RFC 8785, and its SHA-256 is recorded in the Phase 07 acceptance transcript and release manifest. `effect_kind` is `damage`, `heal`, or `restore_resource`; no status-effect scripting exists in v1. A cast deducts cost, performs one d100 check, applies normal/critical amount on success, and applies zero effect on failure. Cost remains spent on a failed roll but not on an invalid action.
- An item is immutable built-in engine data with `target_side`, `effect_kind`, and `amount`. A valid use removes exactly one item before applying its non-random effect. Invalid target or full-resource no-op rejects without consumption. Creator packs may compose only catalog IDs and cannot define behavior.
- On an enemy turn, choose the first spell in lexical spell-ID order that is affordable and has a legal target; otherwise attack. For damage choose the living party actor with the lowest HP percentage by cross multiplication, then lowest absolute HP, then lexical actor ID. For healing/resource restoration choose the legal ally with the greatest missing amount, then lexical ID. Enemy actions use the same validation, RNG, cost, effect, and event ordering as player actions.
- After victory, sort rewards by `item_id`, then offer each unit to living players in lexical `player_id` order with a free slot. Unplaced rewards are discarded and reported in the summary; they never overfill inventory. Corpse exposure chooses eligible non-bound items by ascending inventory slot and then item ID. Loot transfer is atomic and never consumes RNG.
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
- Stories repo (`les-perissables-stories`): add MIT `LICENSE`, `README.md` on day 1
- The stories repo hosts both the MIT pack schema/validation crate and the `storycheck` CLI, so creators can validate packs without any proprietary game-repo code
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

- The client stores one UTF-8 JSON object capped at `64 KiB` with required `client_settings_version: 1`, audio-channel volumes, key bindings, display mode, and behavior-neutral accessibility/presentation preferences. It never contains authoritative run state, Steam tickets, rejoin tokens, provider subjects, or server endpoints.
- Missing settings use documented built-in defaults. Malformed, oversized, unsupported-future, or too-old settings are preserved for diagnosis, ignored with one redacted diagnostic, and replaced in memory by defaults while the first interactive screen clearly reports recovery; startup must not panic.
- Writes use a same-directory temporary file, file sync, atomic rename, and directory sync where the platform supports it. A failed write leaves the previous valid file intact.
- Migration is pure and idempotent. Version `1` has no previous positive version; once version `2` exists, readers support exactly `2` and `1`. Downgrade never destructively rewrites a newer file.

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
- `les-perissables-stories` schema, validator, tooling, and bundled creator examples: MIT

Community creators may:

- Create/share story JSON packs
- Create/share character preset JSON packs
- Create/share UI/theme manifests and original assets

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
- Tier 2 (original assets) later: packs may also ship original tilesets, sprites, maps, audio, combat backdrops, and UI/theme manifests, conforming to the locked asset/TMX conventions above. Enabled only once the hub has submission rules, asset/format validation, and the moderation/abuse controls from the community website plan.
- The schema allows custom asset references from day one, but release builds keep `pack:` media loading disabled unless `publication_attestation` verifies. The attestation contains `pack_id`, version, canonical checksum, publication generation, issued/expiry times, policy tier, and signing key ID, signed as canonical JSON with Ed25519 by an online publication key whose public key is shipped in a signed key set. Attestations expire within `30 days`; the client refreshes them through the hub. Key rotation overlaps old/new public keys; a signed revocation list can revoke key IDs or pack generations. Failure, expiry, or revocation rejects custom media and never falls back to loading it. Local developer builds may enable synthetic fixtures only through an explicit non-release feature proven absent from packaged binaries.
- Engine boundary (unchanged by either tier): presentation (art, audio, UI skin) and narrative (story branching, checks, encounters, character stat/spell composition) are data; combat rules, the d100 system, spell behaviors, and UI behavior remain in the proprietary engine. Creators reskin and re-author the world; they do not change how the game plays.
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
