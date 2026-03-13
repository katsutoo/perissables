# Les Perissables

Les Perissables is a funny, multiplayer, tabletop-style RPG set inside a grocery store. Players pick premade food characters, explore short story-driven scenarios, make terrible decisions, roll dice, and try not to die in aisle-themed combat.

The project works because the joke is loud but the scope is disciplined:

- premade characters instead of full builds
- short replayable runs instead of giant campaigns
- one shared engine with data-driven stories and themes
- authoritative multiplayer instead of peer-to-peer chaos

## Read The Docs By Intent

- For the product vision and concrete content examples, read `docs/product.md`.
- For exact MVP locks and non-goals, read `docs/mvp_contract.md`.
- For technical stack and codebase structure, read `docs/architecture.md`.
- For multiplayer model and authority boundaries, read `docs/networking.md`.
- For story/theme/character pack design, read `docs/content-packs.md`.
- For delivery planning and known risks, read `docs/roadmap.md`.
- For delivery status and phase completion tracking, read `docs/progress_tracker.md`.
- For release and pricing notes, read `docs/release_notes.md`.

## Document Boundary

`docs/mvp_contract.md` is the single source of truth for locked MVP decisions. The other docs are for explanation, examples, rationale, and planning.

## Why This Shape

The original spec mixed pitch, technical architecture, legal policy, roadmap, and task tracking into one long file. Splitting it keeps each document easier to maintain and lowers the chance that locked scope will drift from explanatory notes.
