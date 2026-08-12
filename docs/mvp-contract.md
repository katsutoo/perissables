# Les Périssables MVP Contract

Status: Ready for Phase 01 implementation
Owner: Project team
Updated: 2026-08-12

This document is the source of truth for product, authority, compatibility, and
release requirements. It intentionally does not pre-design every queue, storage
record, retry, or deployment mechanism.

## Decision Classes

- **Locked contract:** changing it alters the product, a compatibility boundary,
  or a release promise. Update this document explicitly.
- **Initial safety ceiling:** a conservative bound that implementation may
  tighten at a phase boundary without expanding scope. Its versioned schema or
  protocol fixture becomes authoritative once released.
- **Evidence-gated decision:** no implementation is selected until the named
  spike or benchmark records the workload, environment, alternatives, and
  result.

If supporting documentation conflicts with a locked contract here, this document
wins. Verification plans define how requirements are tested without redefining
them.

## Product Objective

Ship a funny, fast, native multiplayer pixel-art RPG where two to four players
pick premade food characters, complete short data-driven adventures, and survive
dice/combat events through one authoritative server.

## Milestones

- **Authoritative vertical slice:** Phases 00-02.
- **Playable MVP:** Phases 00-09.
- **Release-ready:** Phases 00-13.
- **Post-MVP:** work outside the Release-ready gate, tracked separately.

## MVP Scope

### In

1. Native Linux and Windows client written in Rust with `raylib`.
2. Authoritative Rust server using `axum`, `tower`, and `tokio`.
3. Server-owned lobby, movement, story, dice, combat, inventory, and run state.
4. JSON-driven stories and character/theme data; TMX maps.
5. Premade characters, no builds or leveling.
6. Compact turn-based combat with attack, spell, item, and pass.
7. Run flow: lobby -> story/world/combat -> summary -> lobby.
8. Three built-in themes: `supermarket`, `garden`, and `storage_room`.
9. Event-driven ambience, music, SFX, and optional short voice barks.
10. Keyboard-only operation and at least two behavior-identical built-in UI
    variants.
11. Steam ownership/authentication, lobbies/invites, and depot distribution.
12. Safe local validation and single-player creator testing for reuse-only
    Tier 1 packs.

### Out

- Character builds, leveling, skill trees, or deep equipment.
- Procedural maps, voice chat, scripting, or arbitrary pack code.
- Mobile, browser, console, or macOS releases.
- Steam achievements.
- Hosted community-pack multiplayer.
- Public UGC hub, accounts, comments, likes, moderation, publication signing,
  custom maps/media, and Tier 2 packs.

Post-MVP features do not reserve implementation detail in this contract.

## Locked Gameplay Rules

### Dice

- Character stats are integers in `5..=70`.
- Internal rolls are equiprobable integers in `0..=100`.
- `0` displays as `000` and is critical success.
- `100` is critical failure.
- Otherwise, `roll <= stat` succeeds.
- Normal success probability is `stat / 101`; total success including the
  critical is `(stat + 1) / 101`.

Production randomness comes from a versioned CSPRNG state supplied explicitly to
headless game rules. Clients and packs never seed it. Rejected actions consume no
randomness. Tests use fixed known states, not frequency assertions.

### Story

- Stories are declarative state machines; no story-specific runtime code.
- Node kinds are dialogue, check, encounter, transition, return, and end.
- Checks have required success/failure branches; missing critical branches fall
  back to the corresponding normal branch while retaining critical feedback.
- Effects are limited to flag and item operations supported by the engine.
- Choices are resolved by authoritative votes. Each eligible player has one
  replaceable vote; ties use the lowest lexical player ID among tied voters.
- Automatic transition chains, events, collections, and text are bounded by the
  versioned content schema.

### Combat And Inventory

- Turn order is descending agility, then lexical actor ID.
- Every turn allows exactly one attack, spell, item, or pass.
- Attack checks strength. Normal success deals
  `max(1, floor(strength / 5))`; critical success doubles it; failures deal
  zero.
- Spells and items use immutable engine-owned effect definitions referenced by
  content IDs. Packs compose known behavior; they do not define code.
- Characters have four inventory slots.
- Dead characters leave only server-approved loot; a corpse exposes at most two
  eligible items to a reachable living player.
- Combat ends on enemy defeat, party wipe, or a bounded engine limit. Invalid
  actions do not mutate gameplay state or consume a turn.

## Runtime Targets And Safety Budgets

Locked product targets:

- Party size: `2..=4`.
- Client presentation: target `60 FPS`.
- Client fixed update: `60 Hz`, with bounded catch-up and dropped-time
  diagnostics.
- Server simulation: `20 Hz`.
- Client gameplay input: at most `20 messages/s`; continuous input is
  coalesced.
- Server state publication: at most `20 Hz`, latest-value coalesced.
- Primary QA resolutions: `1280x720` and `1920x1080`.
- Keyboard focus is always visible and at least two rendered pixels thick at
  supported UI scales.

Initial release capacity goal:

- `64` active sessions and `256` occupied seats on one server instance.
- Rejoin uses an existing seat reservation and remains possible when new
  admission is full.
- Capacity is a Phase 12 measured release gate, not a claim about unbuilt code.
  If the production shape cannot meet it with required headroom, the team
  changes the deployment/capacity plan before release rather than weakening
  correctness.

Initial latency goal under the frozen Phase 12 workload:

- Server processing p95 below `100 ms`, p99 below `200 ms`.
- Same-region scheduled-send-to-correlated-receive p95 below `150 ms`, p99
  below `300 ms`.

Client frame gates are calibrated on the named reference machine before
candidate measurement. They require the 60 FPS target, a predeclared missed
refresh budget, no unexplained update backlog, and no frame above `100 ms`.
No sub-millisecond tolerance is locked before the timer/driver noise floor is
measured.

Playtest targets:

- Median clean-account lobby-to-run start at or below `3 minutes`.
- Median completed or failed run between `20` and `35 minutes`.
- Phase 12 records consent, sample counts, abandoned sessions, medians, and
  uncertainty; no default telemetry is added.

## Architecture Contract

### Workspace

The virtual workspace uses resolver `3`, Rust `1.97.1`, Edition 2024, and
MSRV `1.97.1`. Phase 01 adds a root `rust-toolchain.toml` as the only Rust
toolchain pin, including `rustfmt` and Clippy. `mise` must not declare or
install Rust; `mise.toml` is reserved for tools outside the Rust toolchain and
local task aliases. Initial members are exactly:

- `les-perissables-client`
- `les-perissables-server`
- `les-perissables-game-core`
- `les-perissables-shared`
- `les-perissables-integration-tests`

`shared` owns protocol DTOs and IDs. `game_core` depends on `shared` and
the pinned MIT pack crate. Client and server depend on both. Integration tests
may depend on every member. Cycles and reverse dependencies are forbidden.

### Authority From The First Slice

Phase 02 uses the real client/server boundary for lobby creation/join, one map
interaction, one dice check, one combat action, and summary convergence.

- Client: input collection, rendering, local UI state, and local audio playback.
- Server: identity binding, session lifecycle, validation, movement, story,
  dice, combat, inventory, revisions, and events.
- `game_core`: deterministic headless rules with explicit time and randomness.

A local test identity adapter is permitted only in tests and non-release
development builds. Release features/packages must prove that it is absent.

### Async And Native Boundaries

- Startup owns one supervised task tree.
- Session, connection, heartbeat, identity-provider, and persistence work has
  explicit ownership, cancellation, timeouts, and bounded concurrency.
- Dropping a task handle may not detach correctness-critical work.
- Blocking native work never runs directly on Tokio workers.
- `unsafe` is forbidden in domain crates. Native adapters expose safe owned
  types and document every unsafe block with a local safety contract.
- Recoverable malformed input, I/O, and dependency failures return explicit
  errors; they do not panic.

## Protocol Contract

### Transport And Envelope

- Production uses `wss://`; local development uses `ws://localhost`.
- Steam lobby metadata is discovery only and cannot choose an arbitrary server
  endpoint.
- JSON is the v1 gameplay encoding.
- Every message has `type`, `protocol_version`, `session_id`,
  `player_id`, per-direction transport `seq`, and a typed `payload`.
- Unknown fields/types/directions and unsupported versions are rejected.
- IDs are server-generated opaque values with at least 128 bits of CSPRNG
  entropy.

Message families are join/rejoin, input/result, state/event, resync,
ping/pong, notice, and error. Exact DTOs and byte fixtures are frozen alongside
the Phase 02 protocol implementation.

### Ordering And Replay

- Each player/session has a monotonic `input_seq` that survives reconnect.
- Each session has a monotonic `state_revision` and `event_id`.
- Authenticated expected input receives exactly one applied, superseded, or
  rejected result.
- Gaps, duplicates, stale revisions, queue refusal, and service errors have
  stable outcomes and never cause ambiguous mutation.
- Recipient-specific views exclude secrets, RNG state, hidden triggers, other
  players' private inventory, and other server-only fields.
- Every ingress queue, session mailbox, unresolved-input ledger, writer queue,
  collection, task fan-out, and retry loop is bounded. Exact internal capacities
  are implementation budgets justified by boundary tests and measurements, not
  protocol promises.

Initial wire safety ceilings:

- Reassembled inbound message: `16 KiB`.
- Outbound message: `64 KiB`.
- JSON nesting depth: `64`.
- WebSocket compression and binary gameplay messages: disabled in v1.
- Application heartbeat every `10s`; disconnect after `30s` without a valid
  response.

### Identity And Rejoin

- Production join/rejoin requires a Steam ticket validated for the expected app
  and ownership before authority is granted.
- Tickets and bearer tokens are never logged or persisted raw.
- Rejoin tokens are opaque, random, identity/session/player-bound, digest-stored,
  rotated on acknowledged handoff, and absolutely expiring.
- One player has at most one authoritative connection.
- A disconnected seat remains reserved for an initial `10 minute` grace
  window, bounded by an initial `4 hour` session/token lifetime.
- Rejoin receives a recipient-specific resync and acknowledges it before new
  gameplay input is accepted.
- Production connections use the trusted environment endpoint allowlist with
  normal certificate and hostname verification and no plaintext fallback.

### Session Lifecycle

States are `Lobby`, `Running`, `Summary`, and `Ended`.

- New seats join only a lobby; reserved seats may rejoin non-ended states.
- The lobby owner selects a story. All occupied seats must be connected, ready,
  and use unique characters before start.
- Run completion or wipe enters summary.
- Acknowledgement or bounded timeout returns remaining seats to a cleared lobby.
- Explicit leave releases a lobby seat; socket loss preserves it through grace.
- Empty/expired sessions end and release capacity.
- Combat turns and story votes have monotonic deadlines; all-disconnected runs
  pause gameplay deadlines but not absolute credential/session expiry.

## Evidence-Gated Persistence

Release-ready requires server-owned active runs to survive supported clean
deploys and process crashes according to a documented acknowledgement boundary.
The mechanism is not selected yet.

Phase 10 must measure a production-shaped prototype using representative and
maximum valid state on the actual Railway environment. It compares bundled
SQLite with a mature managed transactional alternative when SQLite misses a
gate, and evaluates:

- full snapshots, deltas, and semantic checkpoints;
- which acknowledged operations require immediate durability;
- commit cadence and batching;
- write amplification, fsync behavior, latency, CPU, memory, and volume use;
- clean restart, crash recovery, corruption, backup, and rollback; and
- operational simplicity for one developer.

Non-negotiable persistence properties:

- Clients never submit authoritative save state.
- Schema and save formats are explicitly versioned and migrated.
- Raw Steam tickets and rejoin tokens are never persisted.
- The store uses parameterized operations and explicit transaction boundaries.
- Tokio workers never perform blocking storage calls.
- Queues, retries, startup validation, artifacts, and recovery work are bounded.
- An uncommitted transition is never reported as durable.
- Corrupt or future state fails safely without being silently deleted.

The Phase 10 decision updates this section with the chosen store, schema,
durability/acknowledgement semantics, measured budgets, failure policy, and
rollback contract before Phase 11 implementation begins.

Before Phase 11, restart/deploy run loss is an explicit pre-release limitation.

## Content Contract

### Repository Boundary

- `perissables`: proprietary runtime, server, built-in content, and assets.
- `les-perissables-stories`: MIT schema/validation crate, `storycheck` CLI,
  conformance corpus, creator examples, and authoring documentation.

The game, validator, and any future service use the same pinned release of the
MIT pack crate. Production builds do not follow a moving Git branch.

### Aggregate Packs

One session uses one aggregate pack identity:

- `pack_id`
- SemVer `version`
- canonical SHA-256 checksum
- `content_schema_version`
- `game_rules_version`

Release-ready hosted multiplayer serves built-in packs from the immutable server
artifact. Clients cannot upload a pack or make the game server fetch one.

Tier 1 local packs may contain stories and character compositions that reference
the shipped identifier catalog. They cannot contain theme documents, maps,
images, audio, scripts, or other custom members.

Tier 2 custom maps/media and their publication/trust system are post-MVP and are
not specified here.

### Schema And Safety

Schema v1 defines strict story, character, theme, and manifest DTOs with:

- required known fields and duplicate-key rejection;
- portable lowercase identifiers and paths;
- explicit cross-reference and graph validation;
- bounded text, collections, nesting, files, and decoded resources;
- canonical checksum fixtures; and
- exact positive/rejection conformance vectors.

Exact parser/resource ceilings are frozen with the schema implementation after
the corpus demonstrates typical, large, limit, and rejected inputs. Limits may
tighten between pre-release schema revisions; they never expand silently in a
released version.

Community archives and every contained byte are hostile input. Validation:

- accepts a minimal regular-file ZIP subset;
- rejects traversal, links, devices, aliases, nested/encrypted archives, and
  case/path ambiguity;
- disables XML external entities/resources and parser network access;
- inspects media bounds before allocation;
- uses bounded, killable workers for native/untrusted decoding; and
- fails closed if a release-required OS sandbox cannot be installed.

The Linux and Windows sandbox profiles are qualified before Phase 08 exposes
creator import. They are not a Phase 01 bootstrap gate.

## Presentation And Local Settings

- Theme/UI skinning changes presentation only, never controls, hit targets,
  focus order, authority, collision, or events.
- Validated assets that fail at presentation time use bounded built-in
  texture/SFX/silence/classic-UI fallbacks.
- Missing required content or untrusted custom media rejects activation; it does
  not fall back into unsafe loading.
- Boot-critical fallback resources are package-integrity checked.

Client settings contain only display, volumes, keybindings, accessibility, and
other non-authoritative preferences.

- Linux data root: `$XDG_DATA_HOME/les-perissables/` when
  `XDG_DATA_HOME` is set to an absolute path, otherwise
  `~/.local/share/les-perissables/`.
- Windows data root: `%AppData%/LesPerissables/`.
- Settings writes use a same-directory temporary file, sync, and atomic replace
  where supported.
- Malformed/future settings are preserved for diagnosis and replaced only in
  memory by safe defaults; startup does not panic.

Exact default bindings and schema are frozen with the Phase 06 settings fixture,
where conflicts and keyboard-only reachability can be tested against the real UI.

## Security And Privacy

- The server validates syntax, identity, authorization, phase, target,
  sequence, rate, and resource bounds before mutation.
- Public errors contain stable codes, not parser internals, credentials,
  filesystem paths, or hidden state.
- Source IP attribution trusts forwarding headers only from configured proxies.
- No request path creates unbounded tasks, allocations, queues, retries, or
  blocking work.
- Logs and verification artifacts redact tickets, tokens, provider subjects,
  message bodies, and personal data.
- No optional telemetry or voice recording is enabled by default.
- Crash upload, if added, is informed, opt-in, bounded, and redacted.
- Security testing targets only local or explicitly authorized isolated staging
  with synthetic accounts/data.

Retention values are deployment policy, recorded before production and reviewed
with end-user terms. Synthetic conformance fixtures and redacted summaries may
remain with source/release history; raw operational evidence has a finite
documented retention.

## Repository, Legal, And CI Baseline

- Main repository is All Rights Reserved with `LICENSE`, `COPYRIGHT`, and
  `README.md`.
- Outside pull requests remain closed until legal review provides an inbound
  contribution policy.
- `Cargo.lock` is committed. Git dependencies use immutable revisions.
- Phase 01 adds `rust-toolchain.toml` to own Rust, Cargo, `rustfmt`, and Clippy.
  `mise` is local convenience for other development tools and task aliases
  only; CI invokes the pinned tools directly rather than invoking `mise` tasks.
- Native `raylib`, Steamworks, and the selected storage dependency are pinned
  with acquisition, checksum, linkage, target, update, and license evidence
  before their release use.

Minimum CI after bootstrap:

```text
cargo fmt --all -- --check
CARGO_BUILD_WARNINGS=deny cargo clippy --locked --all-targets --all-features
cargo test --locked --all-features
cargo test --locked --doc --all-features
RUSTDOCFLAGS="-D warnings" cargo doc --locked --no-deps --all-features
cargo audit --file Cargo.lock
cargo deny check
```

If features become mutually exclusive, replace `--all-features` with the
documented supported matrix.

Release builds are reproducible from pinned inputs, Windows artifacts are signed
through a protected non-exportable workflow, and Linux artifacts publish
checksums/provenance. Exact package/rollback requirements are frozen before
Phase 13 release candidates exist.

## Verification Contract

- Tests protect public behavior and meaningful boundaries; they do not mirror
  every implementation detail.
- Regression fixes and critical invariants are observed failing for the intended
  reason before acceptance. Routine new tests do not require a permanent
  test-only commit or evidence bundle.
- Property and conformance tests cover broad parser/protocol/state invariants.
- Integration tests own internal storage fault injection, migration, and crash
  correctness.
- QA exercises critical user journeys through the shipped entrypoint and
  observes externally visible recovery; it does not duplicate internal fault
  matrices.
- Benchmarks use production-mode artifacts, frozen workloads, calibrated
  environments, raw results, independent runs, correctness/error counts, and
  predeclared gates that exceed measured noise.
- No test or benchmark retries until green, hides intermittent failures, or
  trades correctness for speed.

## Milestone Exit Criteria

### Authoritative Vertical Slice

- Two clients use the real server path from lobby through one interaction,
  check, combat action, and shared summary.
- A legal action mutates authoritative state; an illegal action is rejected
  without mutation.
- Both clients converge on the same revisions/events.
- The client contains no gameplay-authority shortcut.

### Playable MVP

- Deterministic 2-, 3-, and 4-client full runs converge.
- Three consecutive lobby -> run -> summary -> lobby loops retain no stale run
  state.
- Rejoin, replay rejection, heartbeat, capacity refusal, and slow-client
  behavior are bounded and tested.
- All themes and UI variants preserve controls and authority.
- Tier 1 local validation is safe on supported OS baselines.
- Durable restart recovery is not required until Release-ready.

### Release-ready

- The Phase 10 storage decision is implemented and clean/crash recovery matches
  its documented acknowledgement boundary.
- Release QA passes on exact Linux/Windows package digests.
- Calibrated client and server capacity/performance gates pass with no
  unexpected errors, dropped correctness work, or hidden saturation.
- Steam ownership, lobbies/invites, depots, production deployment, signing,
  rollback, legal terms, third-party notices, and asset provenance are complete.

## Change Control

- Scope changes occur between phases and update affected milestones.
- Evidence-gated details are locked only after their named decision record.
- A released protocol/content/save/settings version never changes silently.
- Features outside this contract stay in the post-MVP backlog.
