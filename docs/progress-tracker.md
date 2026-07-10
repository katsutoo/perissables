# Progress Tracker

Use this checklist as the execution board. Tick boxes as work is completed.

This tracker is intentionally detailed. Use `docs/mvp-contract.md` for locked MVP decisions, `docs/roadmap.md` for sequencing and phase grouping, and this file for step-by-step execution status.

## Pre-Implementation Lock Checklist

This checklist records status only. Exact values and behavior live in `docs/mvp-contract.md`; do not copy them here.

- [x] LOCK-A Product scope, gameplay bounds, platforms, and release milestones are normative.
- [x] LOCK-B Workspace, toolchain, dependency direction, CI, and repository boundaries are normative.
- [x] LOCK-C Protocol payloads, identity binding, versions, session lifecycle, replay, and resync are normative.
- [x] LOCK-D Runtime admission, rate, queue, timeout, heartbeat, overload, and shutdown bounds are normative.
- [x] LOCK-E Persistence projections, token digests, cadence, atomicity, migration, and failure policy are normative.
- [x] LOCK-F Content schema v1, aggregate pack model, portable archive limits, checksum canonicalization, and conformance corpus are normative.
- [x] LOCK-G Testing, QA, security, and benchmark evidence requirements are normative through their dedicated plans.
- [x] LOCK-H Community hub privacy, upload, licensing, moderation, and public-launch gates are normative.
- [x] LOCK-I Outside contributions remain closed pending legal review; reviewed end-user and third-party terms gate release, not private Phase 01 development.

## Phase 00 - Foundation And Scope Freeze

- [x] 00.1 Write `docs/mvp-contract.md` with goals, non-goals, and "not in MVP" list
- [x] 00.2 Freeze core constraints in `docs/mvp-contract.md`, including supported party size `2-4` and maximum `4`
- [x] 00.3 Freeze dice/check rules (`d100`, inclusive stat range `5..=70`, inclusive roll range `0..=100`, `000` crit success, `100` crit fail)
- [x] 00.4 Freeze networking scope for MVP (hosted server, no peer-to-peer)
- [x] 00.5 Define reusable phase-completion criteria and evidence rules in the roadmap and verification plans
- [x] 00.6 Create `docs/README.md` doc map and doc-boundary guidance
- [x] 00.7 Phase 00 complete

## Phase 01 - Repo Bootstrap

- [ ] 01.1 Initialize the five locked workspace members, root folders, committed `Cargo.lock`, and pinned `rust-toolchain.toml`
- [ ] 01.2 Create startup-only entrypoints for client/server plus minimal `game_core`, `shared`, and `integration_tests` crates with the locked dependency direction
- [ ] 01.3 Add logging bootstrap (`crates/shared/src/logging.rs`) using `tracing` and `tracing-subscriber`
- [ ] 01.4 Add local-dev-only `mise` tasks in `mise.toml`: `run-client`, `run-server`, `test`, `lint`, `security-scan`
- [ ] 01.5 Add the exact `--locked` baseline checks and pin the `cargo-audit` tool version
- [x] 01.6 Add legal files and close outside pull requests pending a reviewed policy with written contributor agreement
- [ ] 01.7 Create `les-perissables-stories` repo with MIT `LICENSE` and `README.md` (hosts the schema/validation crate `game_core` depends on from Phase 04 and the `storycheck` CLI from Phase 15)
- [ ] 01.8 Add starter CI that installs the pinned toolchain with `rustup`, runs locked baseline checks without `mise`, and compiles supported Linux/Windows targets
- [ ] 01.9 Phase 01 complete

## Phase 02 - Render Loop And Scene Skeleton

- [ ] 02.1 Create fixed-timestep game loop (`update` + `draw`) in client app layer
- [ ] 02.2 Add scene interface (`Enter`, `Update`, `Draw`, `Exit`)
- [ ] 02.3 Implement placeholder scenes: lobby, world, combat
- [ ] 02.4 Add scene router and debug scene switching keys
- [ ] 02.5 Add framerate/debug overlay that records frame/update distributions and missed updates without changing release behavior
- [ ] 02.6 Add audio manager with channel buses (`ambience`, `music`, `sfx`, `voice`) and volume controls
- [ ] 02.7 Add deterministic update-loop/scene-transition tests and release-build checks proving debug keys/characters/features are absent
- [ ] 02.8 Phase 02 complete

## Phase 03 - Tilemap World Prototype

- [x] 03.1 Map format locked for MVP: `TMX`
- [ ] 03.2 Load/render a sample supermarket map through a minimal `ThemeAssets` interface that Phase 10 extends without temporary hardcoded paths
- [ ] 03.3 Implement blocked tile collision and movement rejection
- [ ] 03.4 Implement 4-direction player movement and sprite facing/animation
- [ ] 03.5 Add camera follow and map bounds clamp
- [ ] 03.6 Add footsteps and exploration ambience playback in world scene
- [ ] 03.7 Add collision, map-edge, camera-clamp, missing-layer, and asset-failure tests at zero/one/limit boundaries
- [ ] 03.8 Phase 03 complete

## Phase 04 - Story Schema v1 (Core)

- [ ] 04.1 Implement the complete schema v1 already defined in `docs/mvp-contract.md` in `les-perissables-stories`, including every required field, tagged variant, cross-reference, custom-asset attestation gate, and locked bound
- [ ] 04.2 Implement the schema Rust structs + loader/validation in that MIT crate (strict validation, helpful path-based errors)
- [ ] 04.3 Have `game_core` depend on the MIT crate for the data model; its `story/` module holds runtime logic, not the schema definition
- [ ] 04.4 Add table-driven tests for exact valid/invalid values, unknown fields, reference failures, and boundary-adjacent limits
- [ ] 04.5 Add bounded fuzz targets and the shared canonicalization/checksum conformance corpus defined by `docs/test-strategy.md`
- [ ] 04.6 Publish a versioned identifier catalog for shipped map/theme/sprite/audio/spell IDs without redistributing proprietary asset bytes
- [ ] 04.7 Add an MIT canonical creator example in `les-perissables-stories` and a separate ARR built-in fixture in this repo
- [ ] 04.8 Phase 04 complete

## Phase 05 - Story Runtime State Machine

- [ ] 05.1 Implement runtime state struct (current node, flags, completed events)
- [ ] 05.2 Implement trigger resolver (tile position plus trigger ID)
- [ ] 05.3 Implement dialogue panel and branching choice handling
- [ ] 05.4 Implement side effects (set/unset flags, start encounter)
- [ ] 05.5 Add deterministic public-interface tests for branching, side-effect ordering, automatic-transition bounds, unreachable nodes, and dead ends
- [ ] 05.6 Phase 05 complete

## Phase 06 - Dice And Checks

- [ ] 06.1 Implement authoritative RNG v1 using the locked ChaCha20 state, rejection-sampling mapping, serialization, and consumption order for internal rolls `0..=100`
- [ ] 06.2 Implement resolver ordering: `0` critical success, `100` critical failure, then normal success/fail
- [ ] 06.3 Implement stat-based checks (`roll <= stat`), stat cap enforcement (`5..=70`)
- [ ] 06.4 Add UI formatting: display internal `0` as `000`
- [ ] 06.5 Exhaustively test resolver values `0..=100`, stat boundaries, probability formulas, formatting, and critical fallback without frequency-based RNG assertions
- [ ] 06.6 Add dice result SFX mapping (`success`, `failure`, `critical_success`, `critical_failure`)
- [ ] 06.7 Phase 06 complete

## Phase 07 - Combat Core v1

- [ ] 07.1 Implement combat state machine (start, player turn, enemy turn, resolution)
- [ ] 07.2 Implement action set (`attack`, `spell`, `item`, `pass`) with costs/effects
- [ ] 07.3 Implement hit/check resolution using d100 rules
- [ ] 07.4 Load encounter definitions from story data only (no hardcoded fights)
- [ ] 07.5 Use a hardcoded debug character only for Phases 05-07 so the vertical slice can exercise checks/combat before Phase 08 adds the real roster; remove the debug character path when Phase 08 lands
- [ ] 07.6 Return structured outcome to story runtime (win/loss/rewards)
- [ ] 07.7 Add combat music transition hooks (enter/exit combat)
- [ ] 07.8 Add spell-cast SFX and tiny enemy/character combat voice-bark hooks with cooldown
- [ ] 07.9 Add canonical combat transcript tests covering turn order, all actions, invalid-action non-mutation, costs, death, win, wipe, and event IDs
- [ ] 07.10 Phase 07 complete

## Phase 08 - Character Roster And Presets

- [ ] 08.1 Implement the already-defined character schema from the pinned MIT validation crate
- [ ] 08.2 Implement character loader and validation (`5..=70` for every stat)
- [ ] 08.3 Build character-select UI with lock-in flow
- [ ] 08.4 Enforce duplicate/invalid pick rules in headless `game_core`; Phase 13 reuses the same authoritative operation on the server
- [ ] 08.5 Add sample roster pack (fruit/vegetable/canned/frozen archetypes)
- [ ] 08.6 Add optional voice-bark references plus text fallback keys in character data
- [ ] 08.7 Add tests for stat/spell/roster limits, duplicate picks, simultaneous lock-in, invalid references, and no mutation on rejection
- [ ] 08.8 Phase 08 complete

## Phase 09 - Death, Loot, And Run Continuation

- [ ] 09.1 Add character death state and remove dead actors from turn queue
- [ ] 09.2 Add corpse interaction with the locked maximum of `2` exposed eligible items
- [ ] 09.3 Enforce `4` inventory slots and exact transfer rejection/non-mutation rules
- [ ] 09.4 Implement run continuation logic (alive party continues)
- [ ] 09.5 Implement wipe logic (all dead -> run failed summary)
- [ ] 09.6 Add tests for zero/full inventory, corpse exposure `0/1/2`, transfer limits, dead-turn removal, continuation, and wipe
- [ ] 09.7 Phase 09 complete

## Phase 10 - Theme Packs (Supermarket, Garden, Storage Room)

- [ ] 10.1 Implement loading for the `theme_id` already required by content schema v1
- [ ] 10.2 Implement the existing theme manifest with one `ui_variant` field; do not add `default_skin`
- [ ] 10.3 Implement theme asset loader/unloader with missing-asset fallbacks
- [ ] 10.4 Create `assets/themes/supermarket`, `assets/themes/garden`, and `assets/themes/storage_room`
- [ ] 10.5 Verify story switch applies full environment swap for all three themes without logic changes
- [ ] 10.6 Add per-theme audio manifest and fallback rules, including shared default SFX
- [ ] 10.7 Add tests for all themes, missing/invalid assets, decoded-resource caps, fallback behavior, and unload/reload state
- [ ] 10.8 Phase 10 complete

## Phase 11 - UI Variants Per Story

- [ ] 11.1 Add `ui_variant` support from the theme manifest and derive story UI skin from the selected `theme_id`
- [ ] 11.2 Define skin tokens (frame sprites, button sprites, font refs, color tokens)
- [ ] 11.3 Refactor UI draw code to read tokens instead of hardcoded values
- [ ] 11.4 Implement at least two complete skins and fallback behavior
- [ ] 11.5 Verify identical controls/behavior across all skins
- [ ] 11.6 Add renderer/token tests proving skins cannot change hit targets, focus order, controls, authority, or event behavior
- [ ] 11.7 Phase 11 complete

## Phase 12 - Lobby And Story Rotation

- [ ] 12.1 Build lobby scene that lists available built-in stories
- [ ] 12.2 Show story preview metadata (title, estimated duration, theme)
- [ ] 12.3 Add end-of-run summary screen (result, deaths, key events)
- [ ] 12.4 Add reset hooks that clear transient run state safely
- [ ] 12.5 Implement the locked `select_story`, `ready`, `unready`, `start_run`, and `summary_ack` input actions, owner/minimum/all-ready checks, and summary timeout
- [ ] 12.6 Confirm loop: finish run -> summary acknowledgement/timeout -> cleared lobby -> pick next story -> launch
- [ ] 12.7 Add three-run transcript tests proving reset removes transient state while preserving seats, ownership, sequences, and allowed local settings
- [ ] 12.8 Add deterministic same-tick duplicate-character arbitration tests in lexical player-ID order
- [ ] 12.9 Phase 12 complete

## Phase 13 - Multiplayer Authoritative Server

- [ ] 13.1 Generate the already-locked leaf DTO schemas from `les-perissables-shared`, check in positive/rejection conformance vectors, and implement every message, tagged action/view/event, output-order, revision, and error-correlation rule from `docs/mvp-contract.md`
- [ ] 13.2 Implement server session lifecycle and lobby-to-run transition
- [ ] 13.3 Move all authority server-side (movement, story state, combat, dice)
- [ ] 13.4 Implement client intent messages only (never trust client outcomes)
- [ ] 13.5 Add strict payload validation plus unknown-message rejection
- [ ] 13.6 Enforce locked message, handshake, connection, mailbox, writer-queue, rate, timeout, and slow-client limits
- [ ] 13.7 Implement CSPRNG rejoin tokens with keyed digest storage, identity/session/player/generation binding, rotation, invalidation, and grace accounting
- [ ] 13.8 Run deterministic scripted `2`, `3`, and `4` client convergence scenarios plus all locked admission boundaries
- [ ] 13.9 Emit semantic audio-cue events only (clients play local audio; no streamed audio payloads)
- [ ] 13.10 Phase 13 complete

## Phase 14 - Sync Robustness And Reconnect

- [ ] 14.1 Implement transport `seq`, persistent per-player `input_seq`, state revisions, and event IDs with locked `u64` no-wrap behavior
- [ ] 14.2 Implement identity-bound reconnect and single-connection takeover
- [ ] 14.3 Implement `ClientResyncState` plus matching `resync_ack` before accepting inputs
- [ ] 14.4 Add ping/pong heartbeat and idle-timeout disconnect rules
- [ ] 14.5 Add server-side sanity checks for illegal movement/actions
- [ ] 14.6 Add duplicate/replay input protection using sequence validation
- [ ] 14.7 Add seeded, simulated-transport chaos tests with injected clocks, recorded schedules, no sleeps, and repeated parallel runs
- [ ] 14.8 Phase 14 complete

## Phase 15 - Content Tooling For Story Creators

- [ ] 15.1 Build the `storycheck` CLI in `les-perissables-stories` on top of the existing MIT schema/validation crate (from Phase 04) for story/character/theme schema and reference validation
- [ ] 15.2 Validate cross-file references (story -> character -> theme -> assets)
- [ ] 15.3 Validate the aggregate tuple, license/attribution, portable archive, RFC 8785 canonicalization, and canonical SHA-256 checksum
- [ ] 15.4 Add `--dry-run` graph walk for branching reachability and dead ends
- [ ] 15.5 Write creator docs with a minimal first-story tutorial plus licensing boundaries
- [ ] 15.6 Add example packs and a common-error troubleshooting section
- [ ] 15.7 Implement safe local import: validate archive, copy to a temporary per-user directory, fsync, atomically install by aggregate identity/checksum, and never execute/load directly from an untrusted archive
- [ ] 15.8 Add Tier 1 community-pack discovery, single-player creator-test selection, update/removal, conflict, and exact final/temporary disk-quota behavior; hosted community multiplayer and unattested Tier 2 media remain disabled
- [ ] 15.9 Phase 15 complete

## Phase 16 - Save System And Session Persistence

- [ ] 16.1 Implement server-side snapshot `save_version: 1` and the locked current-plus-previous migration/rejection strategy
- [ ] 16.2 Serialize `PersistedSessionSnapshot`, including logical tick, revision, RNG state/version, pack identity, token digests/generations/expiries, and never raw bearer tokens
- [ ] 16.3 Implement fenced lease writes, staged-transition commit/publish/discard semantics, bounded group-commit WAL, common-low-watermark compaction, immutable snapshot generations, checkpoints, retries, fsync ordering, and newest-valid restore selection
- [ ] 16.4 Restore snapshotted sessions on server startup so restarts/deploys do not destroy runs
- [ ] 16.5 Implement resume validation using selected story ID plus aggregate pack ID/version/checksum, content schema version, game-rules version, and save version; clients never submit run state
- [ ] 16.6 Add corruption handling (backup snapshot slot plus safe session-teardown message to affected players)
- [ ] 16.7 Add current/previous/future/corrupt/oversized fixtures and fault injection at every persistence checkpoint
- [ ] 16.8 Keep client-local saves for non-authoritative data only at locked per-OS roots
- [ ] 16.9 Run deterministic release-mode local `restore-v1` and `compaction-v1` qualification at maximum fixture state, including every rename/fsync interruption; Phase 17 reruns hosted capacity measurements
- [ ] 16.10 Phase 16 complete

## Phase 17 - QA, Balance, And Performance

- [ ] 17.1 Create the consented playtest protocol using the contract's exact timing endpoints, minimum participant/session counts, abandonment/reuse rules, bootstrap confidence interval, difficulty/completion bands, readability tasks, and fixed survey questions/scales
- [ ] 17.2 Collect consented/manual playtest balancing data with the locked sample policy, expected probability comparisons, anonymization, retention, and no default telemetry
- [ ] 17.3 Provision production-equivalent staging first, record its Railway plan/region/resources, and run `docs/benchmark-plan.md` against production-profile core artifacts
- [ ] 17.4 Profile client/server hotspots and reduce avoidable allocations only after benchmark results identify bottlenecks
- [ ] 17.5 Run the full regression suite introduced alongside each behavior; do not defer missing parser/branching/combat/dice coverage to this phase
- [ ] 17.6 Execute the Phase 17 pre-release scenarios in `docs/qa-plan.md`; close blocker/critical bugs and preserve evidence for all pass/fail results without claiming final Release-ready `PASS`
- [ ] 17.7 Phase 17 complete

## Phase 18 - Steam Packaging, Production Server, And Release Readiness

- [ ] 18.1 Add GoReleaser config for Rust Linux/Windows release automation
- [ ] 18.2 Add reproducible Linux/Windows build scripts and version stamping; perform two clean builds per target and compare normalized artifacts byte-for-byte with retained commands/image digests/checksums
- [ ] 18.3 Package runtime assets and verify path handling in release builds
- [ ] 18.4 Integrate scoped Steamworks features (ownership/authentication, lobbies, and invites); achievements remain post-MVP
- [ ] 18.5 Bind production multiplayer identity to Steam auth/session tickets
- [ ] 18.6 Promote the benchmarked staging shape to production with trusted WSS allowlist, health/readiness semantics, lease enforcement, structured logs, and graceful drain
- [ ] 18.7 Store session metadata without an endpoint: `session_id`, protocol/content/rules versions, aggregate pack ID/version/checksum
- [ ] 18.8 Add bounded local redacted crash-log rotation and hotfix playbook; any upload is separately informed, opt-in, redacted, and subject to the contract retention table
- [ ] 18.9 Sign Windows artifacts through the protected non-exportable signing workflow and verify signature/timestamp/digest; publish Linux checksums/provenance
- [ ] 18.10 Rerun every `docs/qa-plan.md` scenario on the final Steam/package candidate, including create/join/rejoin/deploy-restore/rollback on both target OSes, and require final `PASS`
- [ ] 18.11 Complete reviewed end-user terms, third-party notices, asset provenance, Steam depot handoff, and contribution-policy decision
- [ ] 18.12 Phase 18 complete

## Phase 19 - Community Web Hub Foundation (Separate Repo)

- [ ] 19.1 Create separate repository `les-perissables-hub` with legal files on day 1 (`COPYRIGHT` + ARR `LICENSE`); scaffold Rust `axum` + `maud` + `htmx` app
- [ ] 19.2 Add pinned locked CI checks plus an exact SQLx migration check against a pinned Postgres service image
- [ ] 19.3 Deploy on Railway; wire the domain via Cloudflare DNS/CDN with TLS in front of the Railway app
- [ ] 19.4 Configure secrets/env (Discord + GitHub OAuth credentials, `DATABASE_URL`, session signing key, R2 credentials)
- [ ] 19.5 Ship Stage 1 landing page (domain, Steam/wishlist link, community links; content-only, no accounts) - may go live before game launch
- [ ] 19.6 Add Railway Postgres + `sqlx` with migrations and backups; define account/pack/like/comment schema
- [ ] 19.7 Add provider-scoped OAuth identities, state/PKCE, exact callback allowlists, secure sessions, CSRF protection, revocation, and explicit reauthenticated linking
- [ ] 19.8 Add immutable R2 upload/validation/publication with signed bounds, private random keys, server-side checksum, content-addressed public keys, and transactional metadata
- [ ] 19.9 Add gated/private Tier 1 sharing with full schema/semantic/reference/graph/checksum/path/resource validation using the pinned MIT crate
- [ ] 19.10 Add pack listing metadata for aggregate identity, versions, checksum, license, attribution, tags, and immutable object generation
- [ ] 19.11 Add safe pack downloads through the cookie-less origin with fixed content type, `nosniff`, attachment disposition, generated filename, and short-lived single-object URLs for private content
- [ ] 19.12 Add likes, comments, and sort/browse-by-likes
- [ ] 19.13 Publish privacy terms, account/content deletion, uploader rights grant, basic quotas/rate limits, and keep all non-synthetic objects private
- [ ] 19.14 Add creator-doc links plus validator-integration guidance
- [ ] 19.15 Phase 19 complete; Stage 2 UGC remains private/not publicly launched until Phase 20 is complete

## Phase 20 - Community Moderation And Trust

- [ ] 20.1 Publish moderation policy and submission rules (including original-asset licensing/IP rules)
- [ ] 20.2 Add report/flag flow plus admin moderation actions (delete content, ban accounts)
- [ ] 20.3 Expand baseline quotas into production anti-spam/abuse controls and moderation audit logs
- [ ] 20.4 Verify privacy/deletion behavior and publish the final policy before opening external registration
- [ ] 20.5 Add asset/format validation for uploads (dimensions, frame counts, formats, sizes, checksums)
- [ ] 20.6 Implement signed publication attestations, key rotation/revocation, and client verification; enable Tier 2 only once moderation + asset validation + attestation delivery are in place
- [ ] 20.7 Phase 20 complete; public Stage 2 UGC launch gate is satisfied
