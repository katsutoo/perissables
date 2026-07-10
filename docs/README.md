# Docs Overview

This directory is split by purpose so product vision, locked MVP decisions, and implementation planning do not drift into one giant file.

## Start Here

- `docs/les-perissables-spec.md` — short project overview and doc map.
- `docs/product.md` — game vision, tone, player loop, and concrete content examples.
- `docs/mvp-contract.md` — locked MVP scope and implementation decisions. Single source of truth for anything marked locked.
- `docs/architecture.md` — client/server structure, stack rationale, crate/module layout, and engineering practices.
- `docs/networking.md` — multiplayer model, server authority, sync approach, and security posture.
- `docs/content-packs.md` — data-driven story/theme/character pack model, examples, and repo boundary.
- `docs/test-strategy.md` — automated-test layers, determinism rules, boundaries, fuzzing, and evidence.
- `docs/qa-plan.md` — release-artifact QA matrix, scenarios, verdicts, and evidence.
- `docs/security-model.md` — threat boundaries, required controls, and safe verification policy.
- `docs/benchmark-plan.md` — production-mode workloads, metrics, gates, and raw-result policy.
- `docs/roadmap.md` — milestone plan, phase-by-phase delivery path, definition of done, and known risks.
- `docs/progress-tracker.md` — checkbox tracker for locks, phases, and completion status.

## Recommended Reading Order

1. `docs/les-perissables-spec.md` — short overview.
2. `docs/product.md` — game vision and example content.
3. `docs/mvp-contract.md` — locked MVP boundaries (single source of truth).
4. `docs/architecture.md`, `docs/networking.md`, `docs/content-packs.md` — implementation shape.
5. `docs/test-strategy.md`, `docs/qa-plan.md`, `docs/security-model.md`, `docs/benchmark-plan.md` — verification.
6. `docs/roadmap.md` and `docs/progress-tracker.md` — sequencing and status.

## Document Rules

- Put exact MVP locks only in `docs/mvp-contract.md`.
- Dedicated verification plans are normative for how locked behavior is tested, secured, exercised, and measured; they must not redefine product or protocol values.
- Use the other docs for explanation, examples, rationale, and planning. If explanatory text sounds normative, replace it with a contract reference.
- Use `docs/progress-tracker.md` as the built-in execution checklist; move extra day-to-day task detail into issues or a project board if the tracker becomes too granular.
