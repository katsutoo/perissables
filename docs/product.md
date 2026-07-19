# Product Vision

## Premise

Les Périssables is a stupid, funny, multiplayer tabletop-style RPG set inside a grocery store. Two to four players pick premade supermarket characters like a Banana Rogue, Canned Beans Paladin, or Leek Bard, then try to survive short story-driven runs full of bad decisions, dice rolls, and lethal combat.

The joke works because the structure is serious even when the characters are not. The game borrows the clarity of tabletop co-op RPGs, but strips out the slow parts: no leveling, no character builds, no giant inventory screens, and no endless campaign commitment.

## Concrete Example First

The project is data-driven, so the easiest way to understand it is through example content.

Example story pack:

- Story title: `Cleanup on Aisle 9`
- Theme: `supermarket`
- Map: one compact TMX grocery floor with blocked shelves, freezers, and checkout lanes
- Hook: the night crew vanished after opening a pallet of cursed discount spices
- Core choices: inspect the spill, break into storage, calm a panicking NPC, or flee deeper into the store
- Failure escalation: alarms, possessed shopping carts, freezer-burn damage/flavor, and an encounter with a Frozen Pizza Golem

Example party:

- Banana Rogue: high agility, slippery utility, cowardly flavor text
- Canned Beans Paladin: tanky, loud, righteous, probably dented
- Leek Bard: support spells, morale boosts, and bad vegetable puns
- Fish Sticks Mage: risky burst damage and freezer-themed magic

One run should feel like this:

1. Players meet in lobby and lock in premade characters.
2. The group enters a short story with a distinct map and theme.
3. They move through the store, trigger events, and make branching choices.
4. Checks resolve through the shared dice system.
5. Encounters switch into a compact turn-based combat screen.
6. The party either completes the run or dies trying, then returns to lobby.

## Player Experience Goals

- Fast onboarding: pick a character and start immediately.
- Strong replayability: swap stories, themes, and party combinations without changing engine code.
- High readability: simple top-down exploration and compact combat UI.
- Social comedy: funny characters, bad luck, and party chaos generate the stories players retell.
- Lethal but short runs: failure is part of the joke, not a multi-hour punishment.

## Visual Direction

- 2D top-down pixel art during exploration, similar in readability to classic Pokemon or early Zelda.
- Separate combat presentation with party and enemies facing off in a dedicated battle view.
- Minimal interface surface: lobby, world scene, dialog/choice box, combat actions, HP bars, dice feedback, and tiny inventory slots.
- Theme swaps should change the atmosphere of a run without changing the underlying rules.

## Story And Theme Model

Stories are data, not code. A story file defines:

- map reference
- triggerable events
- branching choices
- checks
- encounters
- theme selection

That allows the same runtime to support very different tones:

- `supermarket`: fluorescent aisles, carts, checkout chaos
- `garden`: softer paths, overgrown produce, compost enemies
- `storage_room`: cramped backrooms, pallets, cold lighting, heavier survival tone

The core rule is simple: the engine stays shared, while flavor comes from data packs.

## Combat And Dice

Combat is turn-based, dice-driven, and intentionally compact. Characters have preset stats and spells. Encounters come from story data, not hardcoded scene scripts.

The game uses a d100-style system with intentionally unusual critical extremes for flavor. Exact rule definitions, resolution order, and locked limits live in `docs/mvp-contract.md`.

## Platform And Release Shape

- Native desktop game, not a browser game.
- Multiplayer is online through a hosted authoritative server.
- Steam lobbies and invites help players find/join sessions, but gameplay authority stays on the server.
- Steam-first distribution is the expected path.

Business-facing release notes and pricing decisions stay out of these docs; they are tracked separately when needed.
