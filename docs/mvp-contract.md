# Les Périssables MVP Contract

Status: Locked for pre-implementation planning
Owner: Project team
Updated: 2026-07-02

This document is the single source of truth for locked MVP scope and implementation decisions.

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
4. Basic turn-based combat integrated with story encounters.
5. Preset character roster (no character build system, no leveling).
6. Run flow: lobby -> story -> end summary -> lobby.
7. Three environment themes:
   - supermarket
   - garden
   - storage_room
8. Core multiplayer for 2-4 players with server authority.
9. Event-driven audio system with channels for ambience, music, SFX, and tiny voice barks.

## Explicit Non-Goals (Not In MVP)

1. Skill trees, progression, talent builds, class customization.
2. Advanced inventory management or deep equipment systems.
3. Procedural map generation.
4. Voice chat.
5. Mod scripting with arbitrary code execution.
6. Story-specific custom UI logic (visual skinning only in MVP).
7. Mobile, browser, or console release targets.

## Locked Runtime Constraints

- Max party size: `4`
- Max active sessions per server instance before admission is refused: `64`
- Max players per server instance before admission is refused: `256`
- Target rate: `60 FPS`
- Client fixed-timestep update rate: `60 Hz`
- Server simulation tick rate: `20 Hz`
- Server state broadcast rate: `20 Hz` max; broadcasts may be coalesced, but authoritative ticks must not be skipped silently
- Hosted-server latency target for Phase 17 benchmarking: p95 input->authoritative-state-broadcast under `100 ms`, p99 under `200 ms`, with `64` active sessions of `4` players on the production-equivalent Railway instance
- Tile size: `16x16`
- Primary resolutions for MVP QA: `1280x720`, `1920x1080`

## Locked Persistence Model (Release-ready v1)

This model is locked now, but implementation lands in Phase 16 and is part of the Release-ready milestone, not the Playable MVP milestone. Before Phase 16, restarts and deploys may destroy active runs and must be treated as known pre-release behavior.

- Run state is server-authoritative and persists server-side: the game server snapshots active sessions to durable storage (Railway volume) so restarts and deploys do not destroy runs.
- Rejoin tokens and their expiry windows are part of the session snapshot, so a server restart does not invalidate reconnects that would otherwise still be allowed.
- Clients never submit saved run state; resume always happens through the rejoin flow into a server-restored session.
- Longer-term resume (the whole party returning after rejoin tokens expire) is post-MVP; if added, re-admittance is by validated Steam identity, never by extending token lifetime.
- Client save paths hold only non-authoritative local data (settings, keybinds, local preferences):
  - Linux: `~/.local/share/les-perissables/`
  - Windows: `%AppData%/LesPerissables/`

## Locked Architecture Decisions

- HTTP/router stack: `axum` + `tower`
- Async runtime: `tokio`
- WebSocket transport: `axum` WebSockets
- Production transport: `wss://` (TLS terminated by reverse proxy)
- Local development transport: `ws://localhost`
- WebSocket payload for MVP: JSON
- Production game-server hosting: Railway service running `crates/server`, with separate `staging` and `production` environments
- MVP server scaling policy: one active game-server instance per environment; do not enable multiple replicas until session state is externalized or sticky session allocation is implemented
- Game-server config: clients receive the production WebSocket base URL from build/environment config, never from hardcoded gameplay logic
- Session discovery: Steam lobbies/invites are discovery only; the authoritative server owns `session_id`, player IDs, game state, dice, combat, and story progression
- Steam lobby mapping: lobby metadata stores the server WebSocket URL, `session_id`, `pack_id`, `version`, checksum, and protocol/schema versions so invitees connect to the same authoritative session
- Production identity: production `join`/`rejoin` requires Steam auth/session-ticket validation before the server issues or accepts player/session credentials
- Operations baseline: expose `/healthz` and `/readyz`, use structured logs, and monitor deploy status, error rate, disconnect rate, and reconnect failures through the hosting platform
- Map format for MVP: TMX
- Release automation: GoReleaser for build/package generation only; it does not change licensing or grant public binary distribution rights
- Distribution policy: production desktop binaries ship through Steam depots, not public release pages
- Release targets: Linux + Windows only
- macOS policy: deferred (requires Apple Developer Program for signing/notarization workflow)
- Content schema versioning: `schema_version` integer, start at `1`, reject unsupported major versions for story, character, theme, and pack manifests
- Save schema versioning: `save_version` integer for server-side session snapshots and any versioned client-local files, support current + previous version with explicit migrators

## Locked Protocol And Limits

### WebSocket Envelope v1

All messages must use:

```json
{
  "type": "input",
  "schema_version": 1,
  "session_id": "ses_...",
  "player_id": "ply_...",
  "seq": 1,
  "payload": {}
}
```

Required fields are mandatory. Unknown message types, missing fields, stale/duplicate sequence numbers, and unsupported schema versions are rejected.

Bootstrap exception: a first `join` is sent before the server has issued IDs. For `join` only, `session_id` and `player_id` are empty strings; the server's join response issues the real values, which are mandatory on every later message (including `rejoin`).

Allowed `type` values for v1:

- `join`
- `join_response`
- `ready`
- `input`
- `state`
- `event`
- `error`
- `ping`
- `pong`
- `rejoin`
- `resync_request`
- `resync_state`

Message directions for v1:

- Client -> server: `join`, `ready`, `input`, `ping`, `pong`, `rejoin`, `resync_request`
- Server -> client: `join_response`, `state`, `event`, `error`, `ping`, `pong`, `resync_state`

Sequence rule:

- `seq` is a per-connection, per-direction monotonic counter and starts at `1` for both client->server and server->client traffic.

Join response payload v1:

```json
{
  "rejoin_token": "opaque_128_bit_minimum_random_token",
  "rejoin_grace_seconds": 600,
  "server_time_ms": 0
}
```

The `join_response` envelope carries the issued `session_id` and `player_id`. The `rejoin_token` is opaque, server-generated with a CSPRNG at `128` bits of entropy minimum, never stored in Steam lobby metadata, and never logged. It stays valid for the active session plus the `rejoin_grace_seconds` window after disconnect, is rotated on successful `rejoin`, and is invalidated when the session ends.

Rejoin payload v1:

```json
{
  "rejoin_token": "opaque_128_bit_minimum_random_token"
}
```

For `rejoin`, `session_id` and `player_id` must be non-empty. If accepted, the server sends `resync_state` before accepting new `input` messages from that connection.

### Network Limits Defaults

- Max inbound WS frame size: `16 KiB`
- Max outbound WS frame size: `64 KiB`
- Input rate limit: `20 msg/s` (burst `40`) per client
- Control rate limit: `5 msg/s` (burst `10`) per client
- Heartbeat: every `10s`, disconnect after `30s` timeout
- Control messages are `join`, `ready`, `ping`, `pong`, `rejoin`, and `resync_request`. Gameplay `input` uses the input rate limit; server messages are bounded by frame size and broadcast-rate limits.
- `join` and `rejoin` additionally require per-IP and per-session throttling before token/session validation so token brute-force attempts are rate-limited even when they fail authentication.
- `resync_state` must fit in one `64 KiB` outbound frame. If a full session snapshot would exceed that cap, the server rejects the session shape during validation; protocol fragmentation is out of scope for v1.
- Max serialized session snapshot size: `256 KiB`.

### Content Pack Loading Limits

- Max pack archive size accepted by tooling/runtime: `64 MiB`.
- Max files per pack: `512`.
- Max relative path length: `240` bytes UTF-8.
- Max single JSON file size: `1 MiB`.
- Max story nodes per story: `512`.
- Max choices per node: `8`.
- Max TMX file size: `4 MiB`.
- TMX/XML parsing must disable external entities and external resource resolution, enforce parser depth/size limits, and treat all community packs as hostile input.

## Locked Pack Compatibility And Checksums

Multiplayer sessions require every player to match all of:

- `pack_id`
- `version`
- canonical SHA-256 checksum

Checksum v1 rules:

- The checksum covers the whole pack, not individual files.
- The pack is treated as a set of relative POSIX paths plus bytes.
- Reject absolute paths, `..`, symlinks, duplicate paths, and paths outside the pack root.
- Sort paths lexicographically by UTF-8 bytes before hashing.
- Hash input starts with `les-perissables-pack-v1\n`.
- For each sorted file, append `path\0length\0bytes\n`, where `length` is the decimal byte length of the hashed bytes.
- For `pack_manifest.json`, hash canonical JSON with the top-level `checksum` field omitted; object keys are sorted, UTF-8 is used, and insignificant whitespace is removed.
- All other files are hashed as raw bytes.
- The manifest stores the checksum as lowercase hex prefixed with `sha256:`.

The game loader, `storycheck`, and the community hub must all use the same MIT validation crate implementation for this checksum so compatibility checks cannot drift.

## Locked Content Conventions

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
- Ambience loops: `.ogg` (`ambience_<theme_id>_<name>.ogg`)
- Music: `.ogg` (`music_<theme_id>_<track>.ogg`)
- SFX: `.wav` (`sfx_<category>_<name>.wav`)
- Voice bark clips: `.ogg` (`voice_<actor>_<line_id>.ogg`)
- Target sample rate: `48 kHz`

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

- Main repo (`les-perissables`): add `COPYRIGHT`, ARR `LICENSE`, `README.md` on day 1
- Stories repo (`les-perissables-stories`): add MIT `LICENSE`, `README.md` on day 1
- The stories repo hosts both the MIT pack schema/validation crate and the `storycheck` CLI, so creators can validate packs without any proprietary game-repo code
- `mise` is for local developer convenience only: `mise.toml` may pin tools and expose tasks, but CI must not depend on `mise`
- CI installs Rust through `rustup`/standard Rust tooling and runs the locked `cargo` checks directly
- Minimum CI checks on first commit:
  - `cargo fmt --all --check`
  - `cargo clippy --all-targets --all-features -- -D warnings`
  - `cargo test --all-features`
  - `cargo audit`

## Privacy And Data Minimization

- No optional telemetry by default in the game client/server.
- No in-game voice chat and no voice recording storage.
- Use only gameplay-required platform/session identifiers.
- If crash upload is added, it must be opt-in and clearly documented.
- Community hub (separate context): the hub collects account data and user-generated content - OAuth provider IDs, display names, shared packs, likes, and comments - which is distinct from the game's minimal-PII stance. The hub must publish a privacy policy, store only what these features require, and support account and content deletion.

## Community Website (Post-MVP)

- Build website as separate repo: `les-perissables-hub`. Keeping it separate preserves the boundary between the ARR game runtime and a public, content-facing site.
- Roll it out in two stages inside that repo:
  - Stage 1 (may ship before the game launches): a simple landing page that points the domain at the project, links the Steam page/wishlist, and links community channels (e.g. Discord). Content-only: no accounts.
  - Stage 2 (implemented soon after the game ships, but not publicly launched until moderation/privacy gates are complete): a community hub where players sign in, share content packs, and discover others' packs.
- Hub purpose: host/discover community data packs (story/character/theme packs), not runtime binaries. Players still need to own the game on Steam to run any pack.
- The hub is a user-generated-content (UGC) social platform: signed-in users can share packs, like them, comment on them, and sort/browse by likes.
- Identity (locked): authentication via Discord and GitHub OAuth only - no homegrown email/password system. Store an opaque provider ID plus display name; the uploading account owns its packs (edit/delete) and is the attribution shown to others.
- Moderation is a launch requirement, not a later add-on: report/flag flow, admin delete/ban actions, and anti-spam/upload limits ship before the Stage 2 UGC hub is publicly opened.
- Privacy is a launch requirement: publish a privacy policy and support account/content deletion. See "Privacy And Data Minimization" for how the hub's data handling differs from the game.
- Locked tech stack: Rust `axum` + `maud` (server-rendered HTML) + `htmx` (interactivity), as one app (a single crate, not a workspace) that starts as the Stage 1 landing page and grows into the Stage 2 hub. For upload validation it depends on the MIT pack schema/validation library crate published from `les-perissables-stories`, so hub-side checks match the game exactly without pulling in proprietary game-repo code.
- Locked hosting/data: deploy on Railway; database is Railway Postgres accessed via `sqlx` with migrations. If the database is ever outgrown, switch to PlanetScale (Postgres); Neon is explicitly not used. Toasty ORM was evaluated and deferred until it is post-1.0/stable.
- Locked storage/CDN: uploaded pack/asset files live in Cloudflare R2 (S3-compatible, zero-egress), accessed from the Railway app via `aws-sdk-s3`/`object_store`; Postgres stores metadata, ownership, and the R2 object key. Uploads use presigned URLs and a validate-before-publish flow (private bucket -> schema/asset validation -> public). Cloudflare also provides DNS + CDN in front of the Railway app. The app and Postgres stay on Railway; the app is not moved to Cloudflare Workers (which would force a WASM rewrite and SQLite/D1).
- Locked license: the hub is your proprietary code and ships All Rights Reserved, like the game (`COPYRIGHT` + ARR `LICENSE` added when the repo is created). This is independent of the MIT data packs it serves and the MIT schema/validation crate it depends on.

## Community Content And Licensing Boundary

The project has a dual model:

- `les-perissables` (main game): proprietary, All Rights Reserved
- `les-perissables-stories` (data/specs): MIT

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
- Multiplayer sessions require matching `pack_id`, `version`, and canonical SHA-256 checksum across players.
- Pack manifest support must include `pack_id`, `version`, `schema_version`, and checksum validation in tooling/runtime checks using the locked checksum rules above.

Creator content tiers (rollout):

- Tier 1 (reuse-only) at first creator release: packs may add new stories, new character stat/spell combinations, and flavor, but must reference already-shipped maps, themes, sprites, audio, and spells. No new asset files. This keeps packs instantly consistent and minimizes moderation/validation load.
- Tier 2 (original assets) later: packs may also ship original tilesets, sprites, maps, audio, combat backdrops, and UI/theme manifests, conforming to the locked asset/TMX conventions above. Enabled only once the hub has submission rules, asset/format validation, and the moderation/abuse controls from the community website plan.
- The story/theme/character schema must allow custom asset references from day one so Tier 2 needs no re-architecture; the rollout gates uploads/acceptance, not engine capability.
- Engine boundary (unchanged by either tier): presentation (art, audio, UI skin) and narrative (story branching, checks, encounters, character stat/spell composition) are data; combat rules, the d100 system, spell behaviors, and UI behavior remain in the proprietary engine. Creators reskin and re-author the world; they do not change how the game plays.
- Sharing/attribution: packs are shared through an authenticated hub account (Discord/GitHub); the uploading account owns and can update/remove its packs and is the displayed attribution. Other signed-in users can like and comment on packs.

## Milestone Exit Criteria

### Vertical Slice Exit

- One playable map + one playable story
- d100 checks working
- basic combat working
- no full multiplayer requirement yet
- may use one hardcoded debug character until the Phase 08 data-driven roster exists

### MVP Exit

- Authoritative multiplayer stable (2-4 players)
- lobby loop is complete
- three themes selectable through story metadata
- core run loop stable for repeated sessions
- durable run persistence is not required until the Release-ready milestone, even though the persistence model is locked above

### Release-ready Exit

- reconnect/resync robustness complete
- server-side session snapshots survive restarts/deploys
- production hosted-server smoke tests cover join, rejoin, and restore

## Change Control

- Any feature outside this contract is added to post-MVP backlog.
- Scope changes only happen between phases, never inside an active phase.
- If scope grows, timeline updates must be acknowledged before coding continues.
- Pricing and other business-facing release notes belong in `docs/release-notes.md` unless they become locked MVP constraints.
