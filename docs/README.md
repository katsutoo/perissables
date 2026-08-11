# Documentation

## Reading Order

1. `docs/product.md` — player experience, tone, and examples.
2. `docs/mvp-contract.md` — product, authority, compatibility, and release
   requirements.
3. `docs/architecture.md` and `docs/networking.md` — implementation shape.
4. `docs/content-packs.md` — data/engine and repository boundaries.
5. `docs/test-strategy.md`, `docs/qa-plan.md`,
   `docs/security-model.md`, and `docs/benchmark-plan.md` — evidence.
6. `docs/roadmap.md` and `docs/progress-tracker.md` — order and status.

## Authority

- `docs/mvp-contract.md` is authoritative for locked product and compatibility
  behavior.
- Supporting documents explain structure and verification. They do not create
  new product requirements.
- Evidence-gated implementation decisions remain explicitly provisional until
  their named phase records the result.
- Exact released protocol/content/save/settings schemas and byte fixtures are
  authoritative for their version and must agree with the contract.
- `docs/progress-tracker.md` contains phase gates only. Active work belongs in
  issues or a project board.

## Repository Boundaries

These documents cover the `perissables` game/runtime repository and its
game-facing integration contracts.

| Repository | Status | URL/pin |
| --- | --- | --- |
| `perissables` | Current proprietary game/runtime repository | [GitHub](https://github.com/nuggocto/perissables) |
| `les-perissables-stories` | Planned for Phase 01; not yet created/reviewed | Not assigned |
| `les-perissables-hub` | Post-MVP; not specified or reviewed | Not assigned |

The stories repository owns its MIT schema/validator/tooling documentation. A
future hub owns its application, OAuth, storage, moderation, security, QA, and
operations documentation. This repository retains only the final game-facing
pack compatibility and provisioning interface.

Update this table in the same change that creates or repins a related
repository. Do not claim an uncreated repository was reviewed.
