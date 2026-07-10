# Content Packs

## Core Idea

The game engine is proprietary, but the story format is meant to be open enough that creators can build new adventures without touching engine code.

The runtime should be able to load:

- story data
- character preset data
- theme/UI manifest data
- pack metadata used for compatibility checks

## Example Content Early

Example story metadata:

```json
{
  "content_schema_version": 1,
  "id": "cleanup_on_aisle_9",
  "title": "Cleanup on Aisle 9",
  "theme_id": "supermarket",
  "map_id": "aisle_9",
  "estimated_minutes": 25,
  "start_node": "intro_spill",
  "nodes": [],
  "encounters": [],
  "triggers": []
}
```

Example event/check node:

```json
{
  "id": "check_open_freezer",
  "kind": "check",
  "label": "Force open the frozen door",
  "stat": "strength",
  "effects": [],
  "outcomes": {
    "success": { "goto": "freezer_open", "effects": [] },
    "failure": { "goto": "alarm_triggers", "effects": [] },
    "critical_success": { "goto": "freezer_open_bonus", "effects": [] },
    "critical_failure": { "goto": "shelf_collapses", "effects": [] }
  }
}
```

The exact critical-branch fallback behavior is normative in `docs/mvp-contract.md`; this example shows all four branches explicitly.

Example character preset:

```json
{
  "content_schema_version": 1,
  "id": "banana_rogue",
  "name": "Banana Rogue",
  "group": "fruit",
  "stats": {
    "strength": 18,
    "agility": 62,
    "wit": 40,
    "charm": 27
  },
  "max_hp": 40,
  "resource_max": 30,
  "spells": ["peel_escape", "cheap_shot"],
  "starting_items": [],
  "sprite": "builtin:banana_rogue",
  "voice_barks": {}
}
```

Example theme manifest:

```json
{
  "content_schema_version": 1,
  "theme_id": "storage_room",
  "tileset": "builtin:storage_room_tileset",
  "props": "builtin:storage_room_props",
  "ambience_loop": "builtin:storage_room_hum",
  "exploration_music": "builtin:storage_room_explore",
  "combat_music": "builtin:storage_room_combat",
  "combat_backdrop": "builtin:storage_room_battle",
  "ui_variant": "cold_room"
}
```

Story metadata chooses a `theme_id`; the selected theme supplies its sole `ui_variant`, which resolves to visual skin tokens. The empty `nodes`/`encounters`/`triggers` arrays above keep the metadata example short and therefore are not a valid runnable story; complete required structures and bounds are normative in `docs/mvp-contract.md`.

Example pack manifest:

```json
{
  "content_schema_version": 1,
  "pack_id": "cleanup-pack",
  "version": "1.0.0",
  "checksum": "sha256:...",
  "game_rules_version": 1,
  "license_expression": "MIT",
  "component_licenses": {
    "stories/cleanup_on_aisle_9.json": {
      "license": "MIT",
      "attribution": "Example Creator"
    }
  },
  "includes": [
    "story:cleanup_on_aisle_9",
    "character:banana_rogue",
    "theme:supermarket"
  ]
}
```

## Aggregate Pack Model

One MVP session uses one aggregate archive and one compatibility tuple. The archive may contain any combination of the content categories below; they are categories inside the aggregate pack, not independently negotiated multiplayer dependencies.

### Story Packs

Define maps, events, choices, checks, and encounters. Stories should be fully runnable through the shared engine without custom code paths.

### Character Packs

Define premade food characters with stats, traits, spell lists, text flavor, and asset references.

### Theme Packs

Define environment-facing presentation: tileset, props, ambience loop, exploration music, combat music, combat backdrop, and UI variant references.

### Pack Manifests

Define aggregate pack identity, version, content schema version, game-rules version, checksum, SPDX license expression, per-component licensing/attribution where needed, and included content references so the server and players can verify compatibility.

## Creator Freedom Tiers

Creator-authored content rolls out in two tiers so the game can ship a consistent, low-moderation experience first and open up full customization once hub tooling exists.

### Tier 1 - Reuse Only (first creator release)

Creators author with the assets the game already ships:

- New stories: branching, choices, checks, encounters, dialogue, flavor.
- New characters: new stat lines and spell loadouts drawn from existing spells, using existing sprites.
- Existing themes only: pick from the shipped themes (supermarket, garden, storage_room).

No new asset files are added, so every pack looks and sounds on-brand. Tier 1 still requires schema, semantic/reference, graph, checksum, path, and resource-limit validation because references to shipped maps, spells, themes, and nodes can be invalid or hostile.

### Tier 2 - Original Assets (later)

Creators may additionally ship their own presentation so a pack fully matches its own setting (for example a haunted mansion or a space station rather than a grocery store):

- Original tilesets, maps, character sprites, music, ambience, SFX, and combat backdrops.
- Original UI/theme manifests (visual skin variants).

All original assets must conform to the locked conventions in `docs/mvp-contract.md` (16x16 tiles, sprite frame order/naming, TMX layer/object rules, audio formats and sample rate). Tier 2 is enabled only after the community hub has submission rules, asset/format validation, and moderation/abuse controls, because arbitrary uploaded art and audio raise moderation, licensing/IP, distribution, and untrusted-file-handling concerns.

The schema supports custom asset references from the start. Release builds load them only with the signed, unexpired, non-revoked hub publication attestation defined in `docs/mvp-contract.md`; ordinary side-loading cannot bypass the Tier 2 gate.

## Presentation vs Engine Boundary

A simple rule governs what creators can and cannot change:

- Data (creator-controllable): presentation (tilesets, sprites, audio, backdrops, UI skins) and narrative (story branching, checks, encounters, character stat/spell composition, flavor text).
- Engine (fixed, proprietary): combat rules, the d100 dice system and its locked limits, spell behaviors/effects, and UI behavior (visual skinning only, per the MVP non-goals).

Creators re-author and reskin the world to fit their own idea, but they play by the same rules and reuse the engine's spell/effect library. New mechanics or new spell effects require engine support and are out of scope for data packs.

## Sharing And Attribution

Packs are shared through the community hub (`les-perissables-hub`), where creators sign in with Discord or GitHub. The uploading account controls its listing and supplies rights/attribution information; uploading alone is not proof of copyright ownership. Other players can like and comment after the public UGC gate. See `docs/mvp-contract.md` for the normative licensing, identity, moderation, and privacy requirements.

## Repo Boundary

Planned split:

- `perissables`: proprietary game runtime, built-in assets, client/server code (`les-perissables` package prefix)
- `les-perissables-stories`: MIT-licensed schemas, tooling, examples, and creator docs

That boundary matters because the goal is to let creators author new data without granting rights to the game runtime itself.

The reusable pack schema and validation rules live as an MIT library crate in `les-perissables-stories`, alongside the `storycheck` CLI that fronts it, so creators can validate packs without touching the proprietary repo. Everything that validates packs depends on the same pinned crate release and conformance corpus: the game loader, `storycheck`, and the proprietary `les-perissables-hub`. The hub is a single-crate app, not a multi-crate workspace.

## Compatibility Rule

Multiplayer sessions require one matching aggregate compatibility tuple across the server and all players. Release-ready hosted servers serve built-in packs from their immutable artifact; client-driven community-pack fetch/upload is post-MVP. Exact fields, hashing, provisioning boundary, and validation rules stay locked in `docs/mvp-contract.md`.

## Untrusted Input Rule

Community packs are hostile input until validated. Pack loading and `storycheck` must enforce the locked size/path/count limits, reject path traversal and symlinks, parse TMX/XML with external entities and external resources disabled, and keep parser fuzz targets for JSON/TMX boundary cases.

## Authoring Principles

- Stories should stay declarative and readable.
- Schema validation errors should point to the exact failing path when possible.
- Tooling should catch cross-file problems before a pack is shared.
- Community content should feel first-class, but it should remain safely outside the proprietary runtime boundary.
