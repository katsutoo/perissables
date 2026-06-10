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

## Creator Freedom Tiers

Creator-authored content rolls out in two tiers so the game can ship a consistent, low-moderation experience first and open up full customization once hub tooling exists.

### Tier 1 - Reuse Only (first creator release)

Creators author with the assets the game already ships:

- New stories: branching, choices, checks, encounters, dialogue, flavor.
- New characters: new stat lines and spell loadouts drawn from existing spells, using existing sprites.
- Existing themes only: pick from the shipped themes (supermarket, garden, storage_room).

No new asset files are added, so every pack looks and sounds on-brand, plays cleanly in multiplayer, and needs only schema validation. This is the recommended default and covers the large majority of "make your own adventure" cases without any art skill.

### Tier 2 - Original Assets (later)

Creators may additionally ship their own presentation so a pack fully matches its own setting (for example a haunted mansion or a space station rather than a grocery store):

- Original tilesets, maps, character sprites, music, ambience, SFX, and combat backdrops.
- Original UI/theme manifests (visual skin variants).

All original assets must conform to the locked conventions in `docs/mvp-contract.md` (16x16 tiles, sprite frame order/naming, TMX layer/object rules, audio formats and sample rate). Tier 2 is enabled only after the community hub has submission rules, asset/format validation, and moderation/abuse controls, because arbitrary uploaded art and audio raise moderation, licensing/IP, distribution, and untrusted-file-handling concerns.

The schema supports custom asset references from the start (see the theme manifest example above), so enabling Tier 2 is a hub/policy rollout, not an engine change.

## Presentation vs Engine Boundary

A simple rule governs what creators can and cannot change:

- Data (creator-controllable): presentation (tilesets, sprites, audio, backdrops, UI skins) and narrative (story branching, checks, encounters, character stat/spell composition, flavor text).
- Engine (fixed, proprietary): combat rules, the d100 dice system and its locked limits, spell behaviors/effects, and UI behavior (visual skinning only, per the MVP non-goals).

Creators re-author and reskin the world to fit their own idea, but they play by the same rules and reuse the engine's spell/effect library. New mechanics or new spell effects require engine support and are out of scope for data packs.

## Sharing And Attribution

Packs are shared through the community hub (`les-perissables-hub`), where creators sign in with Discord or GitHub. The uploading account owns its packs (it can update or remove them) and is the attribution shown to other players, who can like and comment on packs. See `docs/mvp-contract.md` for the hub's locked stack, hosting, identity, moderation, and privacy decisions.

## Repo Boundary

Planned split:

- `les-perissables`: proprietary game runtime, built-in assets, client/server code
- `les-perissables-stories`: MIT-licensed schemas, tooling, examples, and creator docs

That boundary matters because the goal is to let creators author new data without granting rights to the game runtime itself.

The reusable pack schema and validation rules live as an MIT library crate in `les-perissables-stories`. Everything that validates packs depends on that one crate: the game's loader, the `storycheck` CLI, and the community hub (`les-perissables-hub`). That keeps validation identical everywhere and lets the separate, non-proprietary hub reuse it without depending on any proprietary game-repo code. The hub itself is a single-crate app, not a multi-crate workspace.

## Compatibility Rule

Multiplayer sessions should require matching `pack_id`, `version`, and checksum across all players before a run starts. Exact validation rules stay locked in `docs/mvp-contract.md`.

## Authoring Principles

- Stories should stay declarative and readable.
- Schema validation errors should point to the exact failing path when possible.
- Tooling should catch cross-file problems before a pack is shared.
- Community content should feel first-class, but it should remain safely outside the proprietary runtime boundary.
