# Les Périssables MVP Contract

Status: Ready for Phase 01 implementation
Owner: Sole developer
Updated: 2026-09-12

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
4. JSON-driven built-in stories and character/theme data; TMX maps.
5. Premade characters, no builds or leveling.
6. Compact turn-based combat with attack, spell, item, and pass.
7. Run flow: lobby -> story/world/combat -> summary -> lobby.
8. One supermarket theme, with the storage room as an area of its map.
9. Event-driven ambience, music, and SFX without voice playback.
10. One readable, scalable, keyboard-operable UI.
11. Steam ownership/authentication, private/friends lobbies, invites, and depot
    distribution.

### Out

- Character builds, leveling, skill trees, deep equipment, account progression,
  or long-term saves.
- Starting a solo run. New runs require two to four players; an already-started
  run may continue with one remaining player after departures or deaths.
  One-player starts exist only in development and automated tests.
- Procedural maps, voice chat, voice barks, scripting, or arbitrary content
  code.
- Additional themes and UI skins/variants.
- Mobile, browser, console, or macOS releases.
- Public lobby browsing, matchmaking, mid-run kicking, or Steam achievements.
- Durable recovery of an active run after a server process crash.

Excluded features do not reserve implementation detail in this contract.

## Release Content Minimum

| Content | Release minimum |
| --- | --- |
| Built-in story | One handcrafted `35-45` minute route |
| World map | One supermarket TMX map with a storage-room area |
| Playable characters | Four |
| Normal enemy types | Five |
| Bosses | One mandatory boss; no surviving route can bypass it |
| Character spells | Eight total |
| Items | Eight total |
| Combat encounters | Two normal encounters and one boss encounter per run |
| Checks | Four to six presented per run |
| Major choices | Three presented per run |
| Dialogue/choice beats | `25-40` presented per run |
| Music | Four tracks: lobby, exploration, combat, and boss |
| Ambience | Two loops |
| SFX | At least twenty distinct effects |
| Voice | None |

Choices may alter local events, checks, rewards, dialogue, and encounter details,
but they do not create substantially different routes or bypass the boss.

## Locked Gameplay Rules

### Dice

- The six character stats are Strength, Perception, Chance, Dexterity, Charisma,
  and Education.
- Character stats are integers in `5..=70`.
- Internal rolls are equiprobable integers in `0..=100`.
- `0` displays as `000` and is critical success.
- `100` is critical failure.
- Otherwise, `roll <= stat` succeeds.
- Normal success probability is `stat / 101`; total success including the
  critical is `(stat + 1) / 101`.

Production randomness comes from a versioned CSPRNG state supplied explicitly to
headless game rules. Clients and content never seed it. Rejected actions consume no
randomness. Tests use fixed known states, not frequency assertions.

### World And Interaction

- Exploration uses continuous cardinal movement without diagonal movement.
- The server owns movement speed, collision, and final position.
- An interaction targets the valid object directly in front of the character
  within one tile.

### Story And Voting

- Stories are declarative state machines; no story-specific runtime code.
- Node kinds are dialogue, check, encounter, transition, return, and end.
- Checks have required success/failure branches; missing critical branches fall
  back to the corresponding normal branch while retaining critical feedback.
- At each story check, the actor is the connected living leader, otherwise the
  connected living player with the lowest lexical player ID. This applies after
  both direct interactions and group choices; the interacting player does not
  automatically roll. Content names the stat, not a different actor policy.
- Actor selection, reading the actor's current stat/modifiers, rolling, and
  committing the outcome are one authoritative operation. Presentation cannot
  postpone or reroll it. If nobody qualifies, the node waits without consuming
  randomness or modifiers and retries when a living player reconnects. Run
  termination cancels the pending node. A disconnect processed before resolution
  excludes that actor; one processed afterward cannot change the result.
- Effects are limited to flag and item operations supported by the engine.
- At run start, the server uses the run CSPRNG to choose one leader uniformly
  from the occupied party.
- Connected living players have one replaceable vote. Disconnect, death, leave,
  or seat expiry discards that player's ballots immediately. Rejoin does not
  restore a discarded ballot; the player may vote again while the vote is open.
- Each story vote declares one `default_choice_id`, shown as the timeout choice.
  Validation requires a nonempty choice set containing that default. Offered
  choices and the default remain fixed while the vote is open.
- A vote remains open for `45 seconds` of unpaused gameplay time. At closure,
  count only ballots from currently connected living players for valid choices.
  Votes consume no randomness. Resolve according to this table:

| Ballots at closure | Outcome |
| --- | --- |
| None | Story selects `default_choice_id`; loot remains unassigned |
| One choice has the highest count | That choice wins |
| Top choices tie and the connected living leader voted for one of them | The leader's choice wins |
| Top choices tie without an eligible leader ballot among them | The tied choice supported by the lowest lexical player ID wins |

- A ballot processed at or after the deadline is rejected. The session owner
  applies due seat expiries before resolving due votes or accepting new input.
  Other accepted events follow the session's serialized order. Ending a run
  cancels its open votes without applying their choices or granting loot.
- A disconnected leader keeps the role through the seat-rejoin grace window but
  cannot break ties while absent. A connected dead leader may transfer
  leadership to one connected living player; this is their only permitted
  gameplay-related control while spectating. If the leader explicitly leaves or
  the seat expires, the role becomes vacant. If connected living candidates
  exist, the server immediately selects a replacement uniformly using the run
  CSPRNG; otherwise it leaves the role vacant without a random draw. A vacant
  role is filled by the same rule when an eligible player next reconnects.
  Leadership changes are authoritative events and cannot postpone vote closure.
- Automatic transition chains, events, collections, and text are bounded by the
  versioned content schema.

### Combat And Inventory

- Characters and enemies have fixed positive HP values supplied by validated
  built-in content; HP is not derived from the six stats.
- Turn order is descending Dexterity, then lexical actor ID.
- Every turn allows exactly one attack, spell, item, or pass.
- Combat has no grid, range, or positional movement. An action chooses from its
  currently valid targets.
- Attack checks Strength. Normal success deals
  `max(1, floor(Strength / 5))`; critical success doubles it; failures deal
  zero.
- Every spell declares its check stat, engine-owned effect ID, target type, and
  per-encounter charges. Player, normal-enemy, and boss spells all roll checks;
  healing and support spells can fail. Casting always consumes the turn and one
  charge.
- On normal success, the spell applies its content-defined magnitude. Critical
  success doubles direct damage, healing, and check-modifier magnitude; a
  critical guard protects against the next two hits instead of one.
- On normal failure, the spell has no effect. On critical failure, it still
  consumes the charge and redirects by effect polarity using the run CSPRNG:
  direct damage or a penalty targets a uniformly random living member of the
  caster's side, including the caster; healing, guard, or a bonus targets a
  uniformly random living opponent. If no redirected target exists, the spell
  has no effect.
- Direct damage and healing use a fixed content-defined amount, with healing
  capped at maximum HP. Guard changes the next incoming positive damage to
  `max(1, floor(damage / 2))`, then consumes one protected hit. A next-check bonus
  or penalty adds or subtracts its content-defined amount once; the modified
  check target is clamped to `0..=100`.
- Each actor holds at most one guard, one next-check bonus, and one next-check
  penalty. Reapplying the same kind replaces its remaining value, even if the
  new value is weaker; kinds never accumulate. Normal guard sets one protected
  hit and critical guard sets two. Zero damage does not consume guard.
- The next attack, spell, or story check uses `stat + bonus - penalty`, clamps
  once, and consumes both modifiers even on a critical roll. Compute with a
  signed intermediate wide enough for the validated magnitudes. Critical `0`
  and `100` take precedence over the modified target. Rejected actions, passes,
  and item uses consume no check modifiers. Guard and modifiers clear on death,
  removal, encounter end, or run reset; disconnect alone does not clear them.
- Content composes these five known effects and cannot define code.
  Each spell or item applies exactly one effect to one target; composition means
  assembling a kit of actions, not ordering multiple effects in one action.
- Enemies do not use an AI subsystem. On their turn, the server uses the run
  CSPRNG to choose uniformly from legal actions in their kit, then uniformly
  from targets valid for that action. Enemy kits may contain only spells and do
  not require a basic attack. With no legal action, they pass.
- Characters have four inventory slots. Items come from story rewards and
  combat loot and are consumed on use.
- Every item declares one of the five effect IDs, its positive magnitude where
  applicable, and one target type: living ally including self, living opponent,
  or self. Damage/penalty items target opponents; healing/guard/bonus items
  target allies or self. Items are usable only on the holder's combat turn,
  apply the normal effect without a roll or critical outcome, and consume one
  item and the turn atomically. Enemy-held items follow the same rules.
- A legal item use at full HP or on an already-buffed target still consumes the
  item and turn; replacement rules apply. Invalid phase, item, or target rejects
  the action without spending anything. Items do not stack within a slot and
  enter the lowest-index free slot.
- Story item grants resolve in authored order for their designated recipients.
  A recipient without a living reserved character or a free slot receives no
  item; that grant is skipped with recipient-appropriate feedback. There is no
  automatic discard of existing items, replacement, rerouting, or pending reward
  queue. The story continues and does not roll back other completed effects.
- Dead players spectate. Apart from the leader-transfer control above, they
  cannot act or vote and cannot be targeted by actions requiring a living target.
- After combat, a corpse exposes at most two server-approved eligible items to
  the living party; no world-distance or combat-position test applies. For each
  item, connected living players cast one replaceable vote for an eligible
  living recipient with a free slot. The leader tie rule and lexical fallback
  apply. An unassigned item disappears when the `45 second` loot vote closes.
- Loot uses the ballot eligibility and timeout table above, with no default
  recipient. Recipient eligibility is checked again at closure. Concurrent loot
  votes resolve in lexical corpse ID, then corpse loot-slot order; each corpse's
  published loot slots are fixed for that loot phase, even after an assignment.
  Each assignment commits before the next vote checks free slots. Invalid
  recipient ballots are discarded, then the remaining ballots are tallied.
- A player combat turn lasts `30 seconds`, then becomes an authoritative pass.
- Combat ends on enemy defeat, party wipe, or a bounded engine limit. Invalid
  actions do not mutate gameplay state, consume randomness, or consume a turn.

## Runtime Targets And Safety Budgets

Locked product targets:

- Party size at run start: `2..=4`; departures may reduce an active run to one.
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

Initial capacity measurement point:

- Phase 12 measures `64` active sessions and `256` occupied seats on one server
  instance; this is a benchmark point, not a launch concurrency promise.
- Rejoin uses an existing seat reservation and remains possible when new
  admission is full.
- Results determine instance sizing, admission limits, and scaling. Release does
  not claim unsupported concurrency from an unbuilt or unmeasured deployment.

Initial latency goal under the frozen Phase 12 workload:

- Server processing p95 below `100 ms`, p99 below `200 ms`.
- Same-region scheduled-send-to-correlated-receive p95 below `150 ms`, p99
  below `300 ms`.

Players may connect worldwide through Railway. Physical deployment in every
geographic area is not required; nearby measured latency is the requirement. A
party may span areas, but one server process in one region owns its in-memory
session. Region selection automatically minimizes the party's worst measured
latency, then median latency, then lexical region ID. Phase 10 tests an initial
Americas/Europe/Asia topology and freezes the smallest Railway-only deployment
that provides acceptable play within the entry-level subscription budget.
Additional Railway regions follow observed player demand rather than launching
speculatively.

Client frame gates are calibrated on the named reference machine before
candidate measurement. They require the 60 FPS target, a predeclared missed
refresh budget, no unexplained update backlog, and no frame above `100 ms`.
No sub-millisecond tolerance is locked before the timer/driver noise floor is
measured.

Playtest targets:

- Median clean-account lobby-to-run start at or below `3 minutes`.
- Median completed or failed run between `35` and `45 minutes`.
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

`shared` owns protocol DTOs and IDs. `game_core` depends on `shared` and owns
pure built-in content validation. The server depends on `shared` and `game_core`;
the client depends on `shared` only and renders authoritative views without
linking gameplay rules. Integration tests may depend on every member. Cycles and
reverse dependencies are forbidden.

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
- Session, connection, heartbeat, identity-provider, and deployment-drain work
  has explicit ownership, cancellation, timeouts, and bounded concurrency.
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
- Before a seat is admitted, authentication, create, join, and rejoin messages
  use a pre-session envelope containing `type`, `protocol_version`,
  per-direction transport `seq`, and a typed `payload`. Claimed session/player
  identifiers, when needed for rejoin, remain untrusted payload fields.
- After the server binds an identity to a seat, every session message also has
  server-issued `session_id` and `player_id` envelope fields.
- Unknown fields/types/directions are rejected.
- Production supports exactly one gameplay protocol version. Older/newer clients
  receive a stable `update_required` error before admission and cannot mutate
  state. Deployments drain admitted sessions before removing their server
  version.
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

- Production create/join/rejoin requires a fresh Steam ticket validated for the
  expected app and ownership before authority is granted.
- Tickets are never logged or persisted raw, and the client stores no rejoin
  bearer token locally.
- A validated returning Steam identity may reclaim only its own reserved seat.
- One player has at most one authoritative connection.
- A disconnected seat remains reserved for an initial `10 minute` grace window,
  bounded by an initial `4 hour` in-memory session lifetime.
- Rejoin receives a recipient-specific resync and acknowledges it before new
  gameplay input is accepted.
- Production connections use the trusted environment endpoint allowlist with
  normal certificate and hostname verification and no plaintext fallback.

### New-seat authorization

Steam authentication proves identity and app ownership. New-seat admission also
requires permission to enter the requested game session:

- Successful session creation admits its authenticated creator as lobby owner.
  Every subsequent new seat requires an unexpired, unconsumed server-held join
  grant bound to that session and the joining identity.
- Only the current connected lobby owner may issue or revoke a grant through
  their authenticated session connection, while the session is in `Lobby` and
  the process is accepting new admission. Other players cannot approve themselves
  or another identity. A session ID, Steam lobby ID, ticket, or client claim of
  membership alone grants no admission rights.
- The owner's client requests grants for members admitted by Steam to the
  associated private/friends lobby. Phase 07 connects this to the normal Steam
  invite/membership flow, including owner handoff. The server trusts the current
  owner's explicit authorization, not an applicant's reported Steam membership
  or lobby metadata. No separate approval screen is required.
- Keep at most one pending grant per identity/session. Grant count, issuance
  rate, and absolute lifetime have tested ceilings frozen in Phase 02. Grants
  expire without extending session lifetime. Owner change, run start, drain,
  and session end revoke all pending grants. Leave or seat removal clears any
  grant for that identity; later new-seat admission requires a new owner grant.
- After the existing protocol, identity, and content checks, the session owner
  checks the grant and other admission conditions and consumes the grant
  atomically with successful seat allocation. Failure does not consume a grant,
  change seats/reservations, or disclose gameplay state. An otherwise valid join
  without a valid grant returns `join_not_authorized`; existing compatibility,
  capacity, lifecycle, and drain errors retain their meanings.
- Reserved-seat rejoin uses the authenticated reservation instead of a join
  grant. Grant expiry or owner change cannot revoke a valid reservation.

Phase 02 implements this policy with local identities and freezes its DTOs and
acceptance fixtures. Phase 07 binds it to validated Steam identities and proves
that the invite flow authorizes the intended identity before game admission.
Synthetic capacity clients use the same owner-grant path.

### Session Lifecycle

States are `Lobby`, `Running`, `Summary`, and `Ended`.

- Steam lobbies are private or friends-only and joinable by invite. MVP has no
  public lobby browser or matchmaking.
- New seats join only a lobby; reserved seats may rejoin non-ended states.
- The lobby owner selects a story and may remove a seat only while the session
  is in `Lobby`. There is no mid-run kick.
- If the lobby owner disconnects or leaves, ownership transfers to the
  longest-connected remaining player, then lexical player ID. With no connected
  player, ownership is vacant; apply the same rule on the next join/rejoin.
  Lobby ownership is separate from run leadership and requires no living actor.
- Starting requires two to four occupied seats, all connected, ready, and using
  unique characters, and a process that is not draining. Membership or
  character-selection changes clear readiness.
- Run completion or wipe enters summary.
- Acknowledgement or bounded timeout returns remaining seats to a cleared lobby,
  except during drain, when the session enters `Ended` as specified below.
- Empty/expired sessions end and release capacity.
- Combat turns and story/loot votes have monotonic deadlines. When no living
  player is connected, gameplay and its deadlines pause, including enemy turns
  and automatic story progression. Resume with the remaining durations when a
  living player reconnects; dead spectators cannot keep gameplay running alone.
  Seat, credential, session, and deployment-drain expiry remain absolute.

Departure transitions are authoritative and apply before selecting replacement
leaders or resolving subsequent gameplay:

| Event/state | Seat and gameplay outcome |
| --- | --- |
| Socket loss in any non-ended state | Reserve the seat for grace; discard ballots. Keep its character, HP, inventory, and effects. While gameplay is unpaused, the disconnected character remains targetable and its combat turns time out to pass. |
| Explicit leave or grace expiry in `Lobby` | Release the seat and selection; clear readiness. There is no run inventory. |
| Explicit leave or grace expiry in `Running` | Release the seat and remove its character from targets and turn order; discard inventory/effects without creating corpse loot. Retain only bounded summary history. End its current turn if needed. |
| Explicit leave or grace expiry in `Summary` | Release the seat and its acknowledgement obligation; preserve the already-produced summary for remaining players. |
| No occupied seats remain | Enter `Ended` and release session capacity. |
| Seats remain but no living reserved characters remain in `Running` | Enter the wipe summary; cancel pending story/check/loot work. |
| One living reserved character remains in `Running` | Continue the existing run, pausing if that player is disconnected. A later run still needs at least two players. |

Explicit leave is allowed in every non-ended state and ends the right to reclaim
that reservation. An expired or departed player cannot join the ongoing run as
a new seat. At a grace deadline, expiry wins over rejoin. Session expiry enters
`Ended` even if gameplay is paused.

## In-Memory Sessions And Draining

MVP session and run state exists only in server memory. The game has no database,
long-term save, or account progression.

- Clients never submit authoritative save state.
- A server process crash ends its active lobbies and runs. Reconnecting clients
  receive a stable run-lost error and return to lobby rather than restoring
  partial state.
- A supported deployment enters drain once, becomes unready, and refuses new
  session creation, new-seat admission, join grants, and `StartRun` with
  `server_draining`. Drain and run-start commits have an authoritative order:
  starts committed before drain may finish; starts processed at or after it are
  rejected even if their requests were already queued. No new run may commit
  after the process reports itself unready for drain.
- Drain applies these terminal transitions without waiting for players to leave:

| Session state | Drain behavior |
| --- | --- |
| `Lobby` | Immediately enter `Ended` and release seats, reservations, and grants. |
| `Running` | Let only the current run continue. Completion/wipe enters `Summary`; it can never return to `Lobby` on this process. |
| `Summary` | Preserve the result until acknowledgement or the bounded summary timeout, then enter `Ended` and release reservations. |
| `Ended` | Release remaining owned resources; no admission or rejoin. |

- Reserved-seat rejoin remains available only for non-ended `Running` and
  `Summary` sessions on the live draining process, subject to the existing
  identity, content, and expiry checks. It cannot reset summary or drain
  deadlines. During drain, the summary timeout is absolute, continues without
  connected players, and is capped by the drain deadline.
- Ending an idle or completed session sends a bounded `server_draining` notice
  that returns clients to the connection/lobby entry flow. They can create or
  join a new session on an admitting process; no seat or run is migrated.
  Maintenance closure of an idle/completed session is not reported as run loss.
- Once all sessions are `Ended`, bounded connection cleanup and joining owned
  tasks let the process exit without waiting for peer disconnects. At the drain
  deadline, end any remaining sessions; unfinished runs use the stable run-lost
  outcome. Record natural completion, deadline expiry, and forced termination
  separately. A deadline or forced stop that interrupts a run is not a clean drain.
- Phase 10 measures and freezes the drain deadline, regional deployment shape,
  and rollback procedure on Railway. It does not select a database.

## Content Contract

All content, schema code, validation, fixtures, and assets live in the
proprietary `perissables` repository.

One session uses the built-in aggregate content identity:

- `pack_id`
- SemVer `version`
- canonical SHA-256 checksum
- `content_schema_version`
- `game_rules_version`

The server loads built-in content from its immutable release artifact.

Create, join, and rejoin carry the client's complete aggregate content identity
in the pre-session payload. After syntax/protocol and identity validation, but
before allocating a seat, taking over a connection, or sending gameplay state,
the server requires exact equality of all five fields with its loaded aggregate.
SemVer ranges and matching schema versions alone do not establish compatibility.
Existing sessions retain the aggregate they started with throughout drain.

A well-formed unequal identity returns `content_mismatch`; malformed or missing
fields use the protocol's malformed-input error, and unsupported protocol
versions still return `update_required`. Failure neither mutates the session nor
releases, extends, or takes over a reservation. The client offers update/restart
guidance and does not fetch content from lobby metadata or retry automatically.
The identity is a compatibility claim, not proof of client integrity or authority.

Phase 02 freezes the identity envelope and mismatch fixtures with its slice
content. Phase 04 freezes canonical aggregate/checksum fixtures. The checksum
is locale-independent and covers the shipped aggregate, including all packaged
locales; choosing another language does not change identity.

Schema v1 defines strict story, character, theme, map-reference, and manifest
DTOs with required known fields, duplicate-key rejection, portable lowercase
identifiers/paths, explicit cross-reference and graph validation, bounded text
and collections, canonical checksum fixtures, and exact positive/rejection
vectors. TMX parsing disables external entities, external resources, and parser
network access.

Exact parser and resource ceilings are frozen with typical, large, limit, and
rejected built-in fixtures. Validation exists to reject a bad repository-owned
build before release.

## Presentation And Local Settings

English is the source and fallback language. Launch locales are English (`en`),
French (`fr`), Simplified Chinese (`zh-Hans`), Traditional Chinese (`zh-Hant`),
Japanese (`ja`), Korean (`ko`), German (`de`), international Spanish (`es`), and
Thai (`th`).

- Every player-facing string is externalized and supports Unicode, wrapping,
  locale-appropriate line breaking, and packaged font fallback without network
  access.
- Release text and store copy receive review from fluent credited collaborators
  in every non-English launch locale.
- The single UI and supermarket presentation never control authority,
  collision, or events.
- Built-in assets that fail at presentation time use bounded texture, SFX,
  silence, and classic-UI fallbacks.
- Missing required built-in content rejects activation.
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
- Native `raylib` and Steamworks dependencies are pinned with acquisition,
  checksum, linkage, target, update, and license evidence before release use.

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
- Integration tests own session drain, forced-stop, and run-loss correctness.
- QA exercises critical user journeys through the shipped entrypoint and
  observes externally visible drain and run-loss behavior.
- Benchmarks use production-mode artifacts, frozen workloads, calibrated
  environments, raw results, independent runs, correctness/error counts, and
  predeclared gates that exceed measured noise.
- Capacity gates use a separately identified production-profile build with the
  isolated synthetic identity adapter. Exact production artifacts exclude that
  adapter and undergo real Steam authentication, release QA, client performance,
  and bounded server performance/lifecycle checks with authorized accounts.
  Phase 13 requires both sets of evidence, paired by final source revision and
  the build-equivalence rules in `docs/benchmark-plan.md`; synthetic capacity
  results are never attributed to a production digest.
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
- The supermarket presentation and single scalable UI preserve controls and
  authority.
- Built-in content meets the release minimum and passes schema validation.

### Release-ready

- The Phase 10 regional deployment, graceful drain, forced-stop, and rollback
  behavior matches the measured operating contract.
- Release QA passes on exact Linux/Windows package digests.
- Calibrated client and server capacity/performance gates pass with no
  unexpected errors, dropped correctness work, or hidden saturation.
- Steam ownership, lobbies/invites, depots, production deployment, signing,
  rollback, store assets/disclosures, legal terms, privacy/support contacts,
  third-party notices, asset provenance, and the operating-cost/shutdown plan are
  complete.

## Change Control

- Scope changes occur between phases and update affected milestones.
- Evidence-gated details are locked only after their named decision record.
- A released protocol/content/settings version never changes silently.
- Features outside this contract stay in the post-MVP backlog.
