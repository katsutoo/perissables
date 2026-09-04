# Built-in content

Authority: game-facing content requirements come from
`docs/mvp-contract.md`. Exact schema fields, limits, canonical bytes, and
positive/rejection fixtures freeze with the versioned built-in schema.

## Purpose

Stories, preset characters, enemies, items, and presentation references are data
so the one release story can be edited without moving gameplay authority into
content. Engine-owned rules remain fixed.

One session uses one immutable aggregate content identity:

- pack ID;
- semantic version;
- canonical SHA-256 checksum;
- content schema version; and
- game-rules version.

The server loads the aggregate content from its release artifact.
Clients declare the same five-field identity before create/join/rejoin. The
contract requires exact equality, including the locale-independent aggregate
checksum; compatible-looking SemVer or schema values cannot substitute for it.
Phase 02 supplies admission fixtures and Phase 04 supplies canonical checksum
fixtures. Content identity does not establish client trust.

## Data and engine boundary

Content controls:

- story nodes, choices, checks, encounters, dialogue, and flavor;
- fixed HP and stat values within schema bounds;
- character and enemy spell/item composition from the engine catalog; and
- references to the built-in map, sprites, music, ambience, and SFX.

The engine controls authority, dice/combat behavior, effects, enemy random action
selection, movement, collision, timers, inventory, voting, UI behavior,
validation, compatibility, and resource enforcement. Content never executes
code.

Story choices declare a visible `default_choice_id` for abstention. The engine
owns ballot eligibility, tie resolution, and story-check actor selection.
Spells/items each reference one engine effect and one target type; assembling
character kits does not create multi-effect action scripts. Guard and modifier
replacement, item resolution, and full-inventory handling follow the contract.

## Explanatory shape

This snippet is intentionally incomplete and is not a conformance fixture:

```json
{
  "content_schema_version": 1,
  "id": "cleanup_on_aisle_9",
  "title": "Cleanup on Aisle 9",
  "theme_id": "supermarket",
  "map_id": "aisle_9",
  "estimated_minutes": 40,
  "start_node": "intro_spill",
  "nodes": [],
  "encounters": [],
  "triggers": []
}
```

Complete accepted examples and exact bytes belong to the repository conformance
fixtures, so this document cannot become a second schema.

## Validation ownership

`game_core` owns the built-in content DTOs and pure validation. The repository
owns canonical examples plus accepted and rejected fixtures. Startup validates
the immutable built-in aggregate before admitting a session.

Schema validation:

- rejects duplicate and unknown fields;
- enforces portable lowercase identifiers and paths;
- resolves every story, encounter, effect, map, sprite, and audio reference;
- requires each story vote's default to belong to its nonempty offered choice
  set; validates check stats, single-effect spell/item definitions, target
  types, recipient selectors, and positive bounded magnitudes;
- bounds text, collections, transition chains, graph size, map dimensions, and
  decoded resource dimensions;
- rejects unreachable required nodes and non-terminating automatic chains;
- disables TMX external entities, external resources, and parser network access;
  and
- verifies the canonical content checksum and release-content minimum.

Exact ceilings freeze with typical, large, limit, and rejected fixtures.
Released limits do not change silently.

## Map and asset direction

The MVP uses one finite orthogonal TMX supermarket map with a storage-room area,
`16x16` tiles, explicit collision/event objects, and no scripts or external
resolution.

Built-in presentation conventions:

- four-direction `16x16` character frames;
- Ogg Vorbis for music and ambience;
- signed 16-bit PCM WAV for short SFX; and
- mono/stereo `48 kHz` audio.

Presentation fallback never changes gameplay authority, collision, controls,
focus, timing, or events.

## Ownership boundary

The runtime, built-in assets, schema code, fixtures, and built-in content are All
Rights Reserved in `perissables`.
