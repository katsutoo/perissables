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
  "id": "cleanup_on_aisle_9",
  "title": "Cleanup on Aisle 9",
  "theme_id": "supermarket",
  "map": "maps/aisle_9.tmx",
  "estimated_minutes": 25,
  "start_node": "intro_spill"
}
```

Example event/check node:

```json
{
  "id": "check_open_freezer",
  "type": "check",
  "label": "Force open the frozen door",
  "check": {
    "stat": "strength",
    "dice": "d100"
  },
  "outcomes": {
    "success": { "goto": "freezer_open" },
    "failure": { "goto": "alarm_triggers" }
  }
}
```

Example character preset:

```json
{
  "id": "banana_rogue",
  "name": "Banana Rogue",
  "group": "fruit",
  "stats": {
    "strength": 18,
    "agility": 62,
    "wit": 40,
    "charm": 27
  },
  "spells": ["peel_escape", "cheap_shot"],
  "sprite": "char_fruit_banana_rogue.png"
}
```

Example theme manifest:

```json
{
  "theme_id": "storage_room",
  "tileset": "tileset_storage_room.png",
  "ambience": "music_storage_room_hum.ogg",
  "combat_backdrop": "storage_room_battle.png",
  "ui_variant": "cold_room"
}
```

Example pack manifest:

```json
{
  "pack_id": "cleanup-pack",
  "version": "1.0.0",
  "checksum": "sha256:...",
  "includes": ["story:cleanup_on_aisle_9", "theme:supermarket"]
}
```

## Pack Types

### Story Packs

Define maps, events, choices, checks, and encounters. Stories should be fully runnable through the shared engine without custom code paths.

### Character Packs

Define premade food characters with stats, traits, spell lists, text flavor, and asset references.

### Theme Packs

Define environment-facing presentation: tileset, ambience, combat backdrop, music, and UI skin references.

### Pack Manifests

Define pack identity, version, checksum, and included content references so multiplayer sessions can verify compatibility.

## Repo Boundary

Planned split:

- `les-perissables`: proprietary game runtime, built-in assets, client/server code
- `les-perissables-stories`: MIT-licensed schemas, tooling, examples, and creator docs

That boundary matters because the goal is to let creators author new data without granting rights to the game runtime itself.

## Compatibility Rule

Multiplayer sessions should require matching `pack_id`, `version`, and checksum across all players before a run starts. Exact validation rules stay locked in `docs/mvp_contract.md`.

## Authoring Principles

- Stories should stay declarative and readable.
- Schema validation errors should point to the exact failing path when possible.
- Tooling should catch cross-file problems before a pack is shared.
- Community content should feel first-class, but it should remain safely outside the proprietary runtime boundary.
