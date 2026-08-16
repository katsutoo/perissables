# Documentation

## Reading Order

1. `docs/product.md` — player experience, tone, and examples.
2. `docs/mvp-contract.md` — product, authority, compatibility, and release
   requirements.
3. `docs/production.md` — screens, asset inventory, provenance, and non-code
   release deliverables.
4. `docs/architecture.md` and `docs/networking.md` — implementation shape.
5. `docs/content-packs.md` — data/engine and repository boundaries.
6. `docs/test-strategy.md`, `docs/qa-plan.md`,
   `docs/security-model.md`, and `docs/benchmark-plan.md` — evidence.
7. `docs/roadmap.md` and `docs/progress-tracker.md` — order and status.

## Authority

- `docs/mvp-contract.md` is authoritative for locked product and compatibility
  behavior.
- Supporting documents explain structure, production inventory, and
  verification. They do not create new product requirements.
- Evidence-gated implementation decisions remain explicitly provisional until
  their named phase records the result.
- Exact released protocol/content/settings schemas and byte fixtures are
  authoritative for their version and must agree with the contract.
- `docs/progress-tracker.md` contains phase gates only. Active work belongs in
  issues or a project board.

## Repository Boundaries

These documents cover the `perissables` game/runtime repository and its
game-facing integration contracts.

| Repository | Status | URL/pin |
| --- | --- | --- |
| `perissables` | Current proprietary game/runtime, schema, built-in content, and assets | [GitHub](https://github.com/nuggocto/perissables) |
| `les-perissables-hub` | Post-MVP; not specified or reviewed | Not assigned |

Creator tooling and a separate schema repository are post-MVP and do not exist
yet. A future hub owns its application, identity, storage, moderation, security,
QA, and operations documentation. No game-facing creator compatibility interface
is reserved before that work begins.

Update this table in the same change that creates or repins a related
repository. Do not claim an uncreated repository was reviewed.
