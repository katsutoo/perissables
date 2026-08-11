# Content Packs

Authority: the game-facing content boundary comes from
`docs/mvp-contract.md`. Exact schema fields, limits, canonical bytes, and
positive/rejection fixtures are frozen in the versioned
`les-perissables-stories` validator release.

## Purpose

Stories, preset characters, and presentation references are data so new
adventures do not require engine changes. Engine-owned rules remain fixed.

One session uses one aggregate pack containing at least one story and any
permitted character/theme documents. Multiplayer negotiates one tuple:

- pack ID;
- semantic version;
- canonical SHA-256 checksum;
- content schema version; and
- game-rules version.

Hosted MVP servers load built-in packs from their immutable release artifact.
Clients cannot upload a pack, provide a URL, or make the server fetch one.

## Data And Engine Boundary

Pack-controlled:

- story nodes, choices, checks, encounters, dialogue, and flavor;
- preset stat/spell/item composition within engine bounds;
- references to permitted maps, themes, sprites, and audio.

Engine-controlled:

- authority, combat and dice behavior;
- spell/item effect implementations;
- movement, collision, timers, inventory, and voting rules;
- UI behavior, controls, focus order, and built-in skin tokens; and
- validation, compatibility, and resource enforcement.

No pack executes code.

## Illustrative Shape

This is explanatory, not a validator fixture:

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

The empty arrays make the snippet non-runnable. Complete accepted examples and
their exact bytes belong to the versioned conformance corpus, so documentation
cannot drift into a second schema.

## MVP Creator Scope: Tier 1

Tier 1 is reuse-only:

- new stories using shipped maps/themes;
- new character stat/spell/item compositions using shipped behavior and
  sprites; and
- text/flavor within the schema.

A Tier 1 archive contains only its root manifest and listed story/character
documents. It cannot contain theme documents, maps, tilesets, images, audio,
scripts, or unlisted members.

Tier 1 is available through safe local import and single-player creator testing.
Hosted community multiplayer is not part of MVP.

## Post-MVP: Tier 2 And Hub

Custom maps, themes, sprites, images, audio, public uploads, moderation,
publication signing/revocation, and hosted community multiplayer are Tier 2
post-MVP work.

Their trust, licensing, publication, and provisioning protocols are defined in
the future hub plan/repository only when implementation begins. This repository
will retain the final game-facing compatibility interface, not speculative hub
internals.

## Validation Ownership

The separate MIT `les-perissables-stories` repository owns:

- the `les-perissables-pack` Rust schema/validation crate;
- the `storycheck` CLI;
- canonical creator examples;
- accepted/rejected conformance vectors;
- checksum/path canonicalization; and
- authoring/troubleshooting documentation.

The proprietary game pins one released crate version and lockfile resolution.
It does not copy the schema into a second implementation.

## Validation Requirements

Community packs are hostile input. Before installation:

- enforce compressed, expanded, file-count, path, JSON/XML depth, text,
  collection, map, image, audio, and decoded-memory ceilings;
- accept only a portable regular-file ZIP subset;
- reject absolute/traversal/ambiguous paths, links, devices, nested/encrypted
  archives, duplicate JSON keys, unknown fields, and unresolved references;
- disable DTDs, external entities, XInclude, external URLs/files, and parser
  resolution;
- validate TMX semantics and story graph termination/reachability;
- compute the canonical checksum from validated logical content; and
- remove temporary output on every result.

Exact ceilings are frozen with schema v1 after typical, large, limit, and
rejected fixtures exist. Released limits do not expand silently.

Native/untrusted decode runs only in bounded killable workers under qualified
Linux and Windows sandbox profiles with no network, broad filesystem, inherited
secrets, or child-process escape. If a required release sandbox cannot be
installed, validation is unavailable rather than weakened.

## TMX And Asset Direction

MVP maps are finite orthogonal TMX maps with `16x16` tiles, explicit collision,
event objects, validated in-pack or built-in references, and no scripts or
external resolution. Precise allowed layers/objects and bounds belong to the
schema corpus.

Built-in presentation conventions:

- four-direction `16x16` character frames;
- Ogg Vorbis for ambience/music/voice;
- signed 16-bit PCM WAV for short SFX; and
- mono/stereo `48 kHz` audio.

Presentation fallback applies only after content passed validation. It never
changes gameplay authority, collision, controls, focus, timing, or events.

## Licensing Boundary

- `perissables` runtime, built-in assets, and built-in content are All Rights
  Reserved.
- `les-perissables-stories` schema, validator, tooling, and creator examples
  are MIT.
- Creator pack licenses apply to those packs only and grant no rights to the
  proprietary runtime or built-in assets.
