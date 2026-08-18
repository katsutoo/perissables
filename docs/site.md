# Public website

Authority: product boundaries come from `docs/mvp-contract.md`.

## Scope

One static marketing page for Les Périssables, plus a changelog and news feed.
Nothing else ships on the web.

Sections:

- game presentation — premise, characters, screenshots, trailer;
- changelog — released versions, newest first;
- news — short posts about updates and release milestones;
- Steam link — the store page is the only call to action; and
- footer — legal, privacy, and support contact.

## Shape

- Static HTML/CSS built from repository-owned sources; no application server and
  no database.
- Changelog and news are Markdown files in this repository, published by the
  same release process that ships the game build.
- Content stays consistent with the Steam store copy approved in Phase 12.
- The site never talks to the game server and is never required for gameplay.
