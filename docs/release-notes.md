# Release Notes

Status: exploratory business and release notes, not the source of locked implementation scope.

## Current Direction

- Steam-first release.
- Native Linux and Windows builds first.
- GoReleaser may be used to build release artifacts, but that does not imply public binary distribution.
- Paid production binaries should ship through Steam only, not public GitHub/GitLab release pages.
- No browser version planned.
- macOS remains deferred unless signing and notarization work is adopted.

## Tentative Commercial Note

- Current pricing idea: around `2-3 EUR`.

## When to Update This File

Add a dated entry whenever a business, pricing, distribution, or release decision changes, even if it does not affect the MVP contract. This keeps release context separate from engineering scope without letting decisions disappear into chat history.

## Rule

If a release or business decision affects engineering scope in a hard way, copy the locked part into `docs/mvp-contract.md`. Otherwise it stays here.
