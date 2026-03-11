# Docs Overview

This directory is split by purpose so product vision, locked MVP decisions, and implementation planning do not drift into one giant file.

## Start Here

- `docs/les-perissables-spec.md`: short project overview and doc map.
- `docs/product.md`: game vision, tone, player loop, and concrete content examples.
- `docs/mvp_contract.md`: locked MVP scope and implementation decisions. This is the single source of truth for anything marked locked.
- `docs/architecture.md`: client/server structure, stack rationale, package layout, and engineering practices.
- `docs/networking.md`: multiplayer model, server authority, sync approach, and security posture.
- `docs/content-packs.md`: data-driven story/theme/character pack model, examples, and repo boundary.
- `docs/roadmap.md`: milestone plan, phase-by-phase delivery path, definition of done, and known risks.
- `docs/progress_tracker.md`: detailed checkbox tracker for locks, phases, and completion status.
- `docs/phase_acceptance.md`: reusable template for accepting a completed phase.
- `docs/release_notes.md`: business and release notes that should stay separate from core technical design.

## Recommended Reading Order

1. Read `docs/les-perissables-spec.md` for the short overview.
2. Read `docs/product.md` to understand the game and example content.
3. Read `docs/mvp_contract.md` for the locked MVP boundaries.
4. Read `docs/architecture.md`, `docs/networking.md`, and `docs/content-packs.md` for implementation shape.
5. Read `docs/roadmap.md` for sequencing, `docs/progress_tracker.md` for execution status, and `docs/phase_acceptance.md` when closing a phase.

## Document Rules

- Put exact MVP locks only in `docs/mvp_contract.md`.
- Use the other docs for explanation, examples, rationale, and planning.
- Use `docs/progress_tracker.md` as the built-in execution checklist; move extra day-to-day task detail into issues or a project board if the tracker becomes too granular.
