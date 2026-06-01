# Progress Tracker

Use this checklist as the execution board. Tick boxes as work is completed.

This tracker is intentionally detailed. Use `docs/mvp_contract.md` for locked MVP decisions, `docs/roadmap.md` for sequencing and phase grouping, and this file for step-by-step execution status.

## Pre-Implementation Lock Checklist

- [x] LOCK-1 Map format selected: `TMX`
- [x] LOCK-2 WebSocket transport selected: `axum` WebSockets on `tokio`
- [x] LOCK-3 HTTP/router stack selected: `axum` + `tower`
- [x] LOCK-4 Production transport fixed: `wss://` with reverse-proxy TLS termination
- [x] LOCK-5 WebSocket payload format fixed for MVP: `JSON`
- [x] LOCK-6 Story/theme schema versioning fixed: `schema_version` integer, start at `1`, reject unsupported major versions
- [x] LOCK-7 Save-file versioning/migration fixed: `save_version` integer, support current + previous version with explicit migrators
- [x] LOCK-8 Release targets fixed: Linux + Windows only; macOS deferred
- [x] LOCK-9 WS envelope fixed: `type`, `schema_version`, `session_id`, `player_id`, `seq`, `payload`
- [x] LOCK-10 Network limits fixed: frame caps, per-client rate limits, heartbeat interval/timeout
- [x] LOCK-11 TMX conventions fixed: layer names plus object naming/property rules
- [x] LOCK-12 Asset conventions fixed: sprite sheet/frame order/naming plus audio formats
- [x] LOCK-13 Repo split/legal baseline fixed: ARR main repo plus MIT stories repo with legal files on day 1
- [x] LOCK-14 CI baseline fixed: `cargo fmt --all --check`, `cargo clippy --all-targets --all-features -- -D warnings`, `cargo test --all-features`, `cargo audit`
- [x] LOCK-15 Audio scope fixed: `ambience`/`music`/`sfx`/`voice` channels plus gameplay audio events
- [x] LOCK-16 Voice chat policy fixed: no in-game voice chat; external apps only
- [x] LOCK-17 Community content/licensing boundary fixed: MIT data packs plus ARR runtime/assets
- [x] LOCK-18 Save path conventions fixed: Linux `~/.local/share/les-perissables/`, Windows `%AppData%/LesPerissables/`
- [x] LOCK-19 Privacy baseline fixed: minimum data, no default telemetry, opt-in crash upload if added later
- [x] LOCK-20 Community website plan fixed: separate repo (`les-perissables-hub`), post-MVP, two stages (landing page first, then community hub soon after game launch); data packs only, no runtime binaries
- [x] LOCK-21 Creator content tiers fixed: Tier 1 reuse-only at first release, Tier 2 original assets later (gated on hub moderation + asset validation); presentation/narrative are data, rules/spell-behaviors/UI-behavior stay in the engine
- [x] LOCK-22 Hub stack/hosting/data fixed: Rust `axum` + `maud` + `htmx`, deployed on Railway, Railway Postgres via `sqlx` with migrations; PlanetScale as switch-later option (not Neon); Toasty deferred until post-1.0
- [x] LOCK-23 Hub identity fixed: accounts via Discord + GitHub OAuth only (no homegrown email/password); store opaque provider ID + display name; uploading account owns/attributes its packs
- [x] LOCK-24 Hub is a UGC social platform: share packs, like, comment, sort-by-likes; moderation (report/flag, admin delete/ban, anti-spam) and privacy (policy + account/content deletion) ship at Stage 2 launch

## Phase 00 - Foundation And Scope Freeze

- [x] 00.1 Write `docs/mvp_contract.md` with goals, non-goals, and "not in MVP" list
- [x] 00.2 Freeze core constraints: party size (`4`), tile size (`16x16`), target FPS (`60`), target resolutions (`1280x720`, `1920x1080`)
- [x] 00.3 Freeze dice/check rules (`d100`, stat range `5-70`, `000` crit success, `100` crit fail)
- [x] 00.4 Freeze networking scope for MVP (hosted server, no peer-to-peer)
- [x] 00.5 Define reusable phase-completion criteria for future phases
- [x] 00.6 Create `docs/README.md` doc map and doc-boundary guidance
- [x] 00.7 Phase 00 complete

## Phase 01 - Repo Bootstrap

- [ ] 01.1 Initialize Rust workspace and root folders (`Cargo.toml`, `mise.toml`, `crates/`, `assets/`, `stories/`, `docs/`)
- [ ] 01.2 Create `crates/client/src/main.rs` and `crates/server/src/main.rs` with startup wiring only
- [ ] 01.3 Add logging bootstrap (`crates/shared/src/logging.rs`) using `tracing` and `tracing-subscriber`
- [ ] 01.4 Add `mise` tasks in `mise.toml`: `run-client`, `run-server`, `test`, `lint`, `security-scan`
- [ ] 01.5 Add baseline checks (`cargo fmt --all --check`, `cargo clippy --all-targets --all-features -- -D warnings`, `cargo test --all-features`, `cargo audit`)
- [ ] 01.6 Add legal files (`COPYRIGHT`, ARR `LICENSE`) in main repo scaffold
- [ ] 01.7 Add starter CI workflow at `.github/workflows/ci.yml` with locked baseline checks
- [ ] 01.8 Phase 01 complete

## Phase 02 - Render Loop And Scene Skeleton

- [ ] 02.1 Create fixed-timestep game loop (`update` + `draw`) in client app layer
- [ ] 02.2 Add scene interface (`Enter`, `Update`, `Draw`, `Exit`)
- [ ] 02.3 Implement placeholder scenes: lobby, world, combat
- [ ] 02.4 Add scene router and debug scene switching keys
- [ ] 02.5 Add framerate/debug overlay to confirm stable loop behavior
- [ ] 02.6 Add audio manager with channel buses (`ambience`, `music`, `sfx`, `voice`) and volume controls
- [ ] 02.7 Phase 02 complete

## Phase 03 - Tilemap World Prototype

- [x] 03.1 Map format locked for MVP: `TMX`
- [ ] 03.2 Load and render one sample supermarket map with layers
- [ ] 03.3 Implement blocked tile collision and movement rejection
- [ ] 03.4 Implement 4-direction player movement and sprite facing/animation
- [ ] 03.5 Add camera follow and map bounds clamp
- [ ] 03.6 Add footsteps and exploration ambience playback in world scene
- [ ] 03.7 Phase 03 complete

## Phase 04 - Story Schema v1 (Core)

- [ ] 04.1 Define schema files for story metadata, events, choices, checks, encounters
- [ ] 04.2 Implement Rust structs in `crates/game_core/src/story` matching schema fields
- [ ] 04.3 Implement loader with strict validation and helpful path-based errors
- [ ] 04.4 Add table-driven tests for valid/invalid story files
- [ ] 04.5 Add one canonical example story in `stories/builtin/`
- [ ] 04.6 Phase 04 complete

## Phase 05 - Story Runtime State Machine

- [ ] 05.1 Implement runtime state struct (current node, flags, completed events)
- [ ] 05.2 Implement trigger resolver (tile position plus trigger ID)
- [ ] 05.3 Implement dialogue panel and branching choice handling
- [ ] 05.4 Implement side effects (set/unset flags, start encounter)
- [ ] 05.5 Add tests for branching paths and unreachable-node detection
- [ ] 05.6 Phase 05 complete

## Phase 06 - Dice And Checks

- [ ] 06.1 Implement d100 roll generator with internal roll range `0..100`
- [ ] 06.2 Implement resolver ordering: `0` critical success, `100` critical failure, then normal success/fail
- [ ] 06.3 Implement stat-based checks (`roll <= stat`), stat cap enforcement (`5..70`)
- [ ] 06.4 Add UI formatting: display internal `0` as `000`
- [ ] 06.5 Add table-driven tests covering boundaries (`0`, `1`, stat, `stat+1`, `100`)
- [ ] 06.6 Add dice result SFX mapping (`success`, `failure`, `critical_success`, `critical_failure`)
- [ ] 06.7 Phase 06 complete

## Phase 07 - Combat Core v1

- [ ] 07.1 Implement combat state machine (start, player turn, enemy turn, resolution)
- [ ] 07.2 Implement action set (`attack`, `spell`, `item`, `pass`) with costs/effects
- [ ] 07.3 Implement hit/check resolution using d100 rules
- [ ] 07.4 Load encounter definitions from story data only (no hardcoded fights)
- [ ] 07.5 Return structured outcome to story runtime (win/loss/rewards)
- [ ] 07.6 Add combat music transition hooks (enter/exit combat)
- [ ] 07.7 Add spell-cast SFX and tiny enemy/character combat voice-bark hooks with cooldown
- [ ] 07.8 Phase 07 complete

## Phase 08 - Character Roster And Presets

- [ ] 08.1 Define character preset schema (stats, spells, traits, sprite IDs)
- [ ] 08.2 Implement character loader and validation (`stat <= 70`)
- [ ] 08.3 Build character-select UI with lock-in flow
- [ ] 08.4 Enforce duplicate/invalid pick rules server-side
- [ ] 08.5 Add sample roster pack (fruit/vegetable/canned/frozen archetypes)
- [ ] 08.6 Add optional voice-bark references plus text fallback keys in character data
- [ ] 08.7 Phase 08 complete

## Phase 09 - Death, Loot, And Run Continuation

- [ ] 09.1 Add character death state and remove dead actors from turn queue
- [ ] 09.2 Add corpse interaction with loot limit (`1-2` items)
- [ ] 09.3 Add tiny inventory slot constraints and item transfer rules
- [ ] 09.4 Implement run continuation logic (alive party continues)
- [ ] 09.5 Implement wipe logic (all dead -> run failed summary)
- [ ] 09.6 Phase 09 complete

## Phase 10 - Theme Packs (Supermarket, Garden, Storage Room)

- [ ] 10.1 Add `theme_id` to story metadata schema and validator
- [ ] 10.2 Define theme manifest format (tiles, props, ambience loop, exploration music, combat music, combat backdrop, default skin)
- [ ] 10.3 Implement theme asset loader/unloader with missing-asset fallbacks
- [ ] 10.4 Create `assets/themes/supermarket`, `assets/themes/garden`, and `assets/themes/storage_room`
- [ ] 10.5 Verify story switch applies full environment swap for all three themes without logic changes
- [ ] 10.6 Add per-theme audio manifest and fallback rules, including shared default SFX
- [ ] 10.7 Phase 10 complete

## Phase 11 - UI Variants Per Story

- [ ] 11.1 Add `ui_variant` support in schema or derive it from theme manifest
- [ ] 11.2 Define skin tokens (frame sprites, button sprites, font refs, color tokens)
- [ ] 11.3 Refactor UI draw code to read tokens instead of hardcoded values
- [ ] 11.4 Implement at least two complete skins and fallback behavior
- [ ] 11.5 Verify identical controls/behavior across all skins
- [ ] 11.6 Phase 11 complete

## Phase 12 - Lobby And Story Rotation

- [ ] 12.1 Build lobby scene that lists available built-in stories
- [ ] 12.2 Show story preview metadata (title, estimated duration, theme)
- [ ] 12.3 Add end-of-run summary screen (result, deaths, key events)
- [ ] 12.4 Add reset hooks that clear transient run state safely
- [ ] 12.5 Confirm loop: finish run -> lobby -> pick next story -> launch
- [ ] 12.6 Phase 12 complete

## Phase 13 - Multiplayer Authoritative Server

- [ ] 13.1 Define network protocol messages (`join`, `ready`, `input`, `state`, `event`, `error`)
- [ ] 13.2 Implement server session lifecycle and lobby-to-run transition
- [ ] 13.3 Move all authority server-side (movement, story state, combat, dice)
- [ ] 13.4 Implement client intent messages only (never trust client outcomes)
- [ ] 13.5 Add strict payload validation plus unknown-message rejection
- [ ] 13.6 Enforce max frame size and per-client/session rate limits
- [ ] 13.7 Implement session/rejoin token issue/validate/expiry flow
- [ ] 13.8 Validate `2-4` player synchronization in one local test session
- [ ] 13.9 Emit semantic audio-cue events only (clients play local audio; no streamed audio payloads)
- [ ] 13.10 Phase 13 complete

## Phase 14 - Sync Robustness And Reconnect

- [ ] 14.1 Add sequence numbers/tick IDs to state updates
- [ ] 14.2 Add reconnect handshake (`session_id`, `player_id`, auth/rejoin token)
- [ ] 14.3 Implement full state resync on reconnect before inputs are accepted
- [ ] 14.4 Add ping/pong heartbeat and idle-timeout disconnect rules
- [ ] 14.5 Add server-side sanity checks for illegal movement/actions
- [ ] 14.6 Add duplicate/replay input protection using sequence validation
- [ ] 14.7 Add chaos tests (disconnect/reconnect/late packets) for stability
- [ ] 14.8 Phase 14 complete

## Phase 15 - Content Tooling For Story Creators

- [ ] 15.1 Create `crates/storycheck` CLI for story/character/theme schema and reference validation
- [ ] 15.2 Validate cross-file references (story -> character -> theme -> assets)
- [ ] 15.3 Add pack manifest checks (`pack_id`, `version`, `checksum`)
- [ ] 15.4 Add `--dry-run` graph walk for branching reachability and dead ends
- [ ] 15.5 Write creator docs with a minimal first-story tutorial plus licensing boundaries
- [ ] 15.6 Add example packs and a common-error troubleshooting section
- [ ] 15.7 Phase 15 complete

## Phase 16 - Save System And Session Persistence

- [ ] 16.1 Define save schema with version field and migration strategy
- [ ] 16.2 Serialize run state (party, flags, position, current node, encounter)
- [ ] 16.3 Implement atomic save writes (temp file + fsync + rename)
- [ ] 16.4 Use locked per-OS save roots (Linux `~/.local/share/les-perissables/`, Windows `%AppData%/LesPerissables/`)
- [ ] 16.5 Implement load/resume validation (story ID/version compatibility)
- [ ] 16.6 Add corruption handling (backup slot plus user-facing recovery message)
- [ ] 16.7 Phase 16 complete

## Phase 17 - QA, Balance, And Performance

- [ ] 17.1 Create playtest checklist for pacing, difficulty, and comedy tone
- [ ] 17.2 Collect balancing data (success rates per stat/check type)
- [ ] 17.3 Profile client/server hotspots and reduce avoidable allocations
- [ ] 17.4 Add regression tests for parser, branching, combat math, and dice edges
- [ ] 17.5 Close blocker/critical bugs and verify no regressions
- [ ] 17.6 Phase 17 complete

## Phase 18 - Steam Packaging And Release Readiness

- [ ] 18.1 Add GoReleaser config for Rust Linux/Windows release automation
- [ ] 18.2 Add reproducible Linux/Windows build scripts and version stamping
- [ ] 18.3 Package runtime assets and verify path handling in release builds
- [ ] 18.4 Integrate scoped Steamworks features (lobbies/invites/achievements)
- [ ] 18.5 Bind production multiplayer identity to Steam auth/session tickets
- [ ] 18.6 Add crash log/reporting path and hotfix playbook
- [ ] 18.7 Run release-candidate smoke tests on both target OSes
- [ ] 18.8 Phase 18 complete

## Phase 19 - Community Web Hub Foundation (Separate Repo)

- [ ] 19.1 Create separate repository `les-perissables-hub`; scaffold Rust `axum` + `maud` + `htmx` app deployed on Railway
- [ ] 19.2 Ship Stage 1 landing page (domain, Steam/wishlist link, community links; content-only, no accounts) - may go live before game launch
- [ ] 19.3 Add Railway Postgres + `sqlx` with migrations; define account/pack/like/comment schema
- [ ] 19.4 Add authentication via Discord + GitHub OAuth (store provider ID + display name)
- [ ] 19.5 Add pack sharing/upload flow owned by the uploading account (Tier 1 reuse-only; no runtime binaries), with `storycheck`/`shared` validation
- [ ] 19.6 Add pack listing pages and metadata model (`pack_id`, `version`, `checksum`, `tags`)
- [ ] 19.7 Add safe pack download flow
- [ ] 19.8 Add likes, comments, and sort/browse-by-likes
- [ ] 19.9 Add creator-doc links plus validator-integration guidance
- [ ] 19.10 Phase 19 complete

## Phase 20 - Community Moderation And Trust

- [ ] 20.1 Publish moderation policy and submission rules (including original-asset licensing/IP rules)
- [ ] 20.2 Add report/flag flow plus admin moderation actions (delete content, ban accounts)
- [ ] 20.3 Add anti-spam/abuse controls and upload limits
- [ ] 20.4 Publish privacy policy and support account/content deletion (minimal data handling)
- [ ] 20.5 Add asset/format validation for uploads (dimensions, frame counts, formats, sizes, checksums)
- [ ] 20.6 Enable Tier 2 (original-asset) packs once moderation + asset validation are in place
- [ ] 20.7 Phase 20 complete
