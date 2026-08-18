# Production inventory

Status: Planning inventory
Owner: Sole developer
Updated: 2026-08-16

Product requirements and release counts belong to `docs/mvp-contract.md`. This
file tracks how approved scope becomes finished screens, content, assets, and
store deliverables. A `TBD` is a decision to make, not permission to invent
scope during implementation.

## Product flow and screens

```text
startup/connection
  -> lobby and character selection
  -> world exploration
  -> dialogue, choice, and check feedback
  -> combat
  -> run summary
  -> cleared lobby
```

Required supporting surfaces include settings/keybindings, reconnect/resync
progress, compatibility/update-required notices, and public/run-lost errors.
Exact layouts and overlays are approved before their implementation phase; rough
wireframes belong with the corresponding phase evidence rather than a second UI
specification here.

## Inventory

| Area | Known constraint | Release quantity | Owner/source | Status |
| --- | --- | --- | --- | --- |
| Stories and dialogue | One route; `35-45` minutes; `25-40` beats | One story | Owner and credited friends | Planned |
| Maps and tilesets | One orthogonal TMX supermarket; `16x16` tiles | One map with storage-room area | Owner, credited friends, purchased sources if used | Planned |
| Playable characters | Premade food characters; four-direction `16x16` frames | Four | Owner and credited friends | Planned |
| Enemies and bosses | Random legal kit action/target; mandatory boss | Five normal types and one boss | Owner and credited friends | Planned |
| Spells, items, effects | Engine-owned effect catalog | Eight spells, eight items, five effects | Project | Planned |
| UI and icons | One scalable keyboard-operable UI | One complete set | Owner, credited friends, purchased sources if used | Planned |
| Animation and VFX | Presentation only; no authority changes | TBD after wireframes | Owner and credited friends | Planned |
| Music and ambience | Ogg Vorbis; local playback from semantic cues | Four tracks and two ambience loops | Owner, credited friends, purchased sources if used | Planned |
| SFX | Signed 16-bit PCM WAV; mono/stereo `48 kHz` | At least twenty | Owner, credited friends, purchased sources if used | Planned |
| Voice | Not in MVP | None | None | Cut |
| Fonts and text | `en`, `fr`, `zh-Hans`, `zh-Hant`, `ja`, `ko`, `de`, `es`, `th` | One packaged fallback set and nine reviewed locales | Owner, fluent credited friends, licensed fonts | Planned |
| Steam store media | Capsule art, screenshots, trailer, copy | One approved set localized for nine locales | Owner and credited friends | Planned |

Remaining `TBD` values are locked before production work in that area is
estimated or accepted.

## Asset pipeline and naming

- Preserve editable source files and record the tool/version needed to export
  them; runtime assets contain only approved exports.
- Every friend contribution needs written permission or assignment covering game
  and Steam distribution, modification, credit, and continued use. Friendship or
  a verbal agreement is not provenance.
- Runtime IDs and paths follow the portable lowercase rules from the content
  contract. One manifest maps each shipped asset to its source, creator,
  license, approval, and export settings.
- Prefer reproducible command-line exports when the source tool supports them.
  Otherwise record the manual export steps and review the resulting checksum.
- Keep temporary, source-tool, and unlicensed files out of release packages.
- The approved font set must render all nine launch locales, French diacritics
  in the title, keyboard prompts, and fallback glyphs without network access.

## Acceptance

An asset or content item is complete only when:

- its owner/source and license or employment ownership are recorded;
- its runtime format, dimensions, duration, and decoded resource use pass the
  applicable validator limits;
- it works in the shipped entrypoint at required resolutions and UI scales;
- missing or failed presentation uses the documented safe fallback;
- its source and exported checksum are recorded; and
- no placeholder or debug asset remains in release features/packages.

## Early playtest record

Character names, HP/stat sheets, spell assignments, and the five normal-enemy
plus boss kits remain Phase 05 content work. Enemy kits may be spell-only. Their
final values must be recorded before Phase 08 asset/content production begins.

Each scheduled playtest records the build/content revision, participant relation
to the project, party size, run duration, abandoned runs, rules questions,
observed idle time, confusing choices, defects, and whether participants wanted
to replay. Phase gates use the smallest useful session and do not invent a
numeric fun score before baseline evidence exists.

## Non-code release inventory

Phase 12 prepares and assigns owners for Steam app/store access, capsule art,
screenshots, trailer, store copy, age/content disclosures, launch languages,
pricing proposal, privacy/support contacts, incident handling, server cost,
shutdown policy, and the public marketing site (`docs/site.md`). Phase 13 approves and publishes the final set against the exact
release artifacts.
