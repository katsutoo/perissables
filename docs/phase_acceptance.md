# Phase Acceptance Template

Use this template at the end of each phase before marking `Phase xx complete`.

## Phase Header

- Phase ID:
- Phase Name:
- Date:
- Owner:

## Goal Check

- Goal statement:
- Goal achieved: [ ] Yes [ ] No
- If No, blockers:

## Deliverables Check

- [ ] Deliverable 1 complete
- [ ] Deliverable 2 complete
- [ ] Deliverable 3 complete

Notes:

## Technical Quality Check

- [ ] No ignored errors in changed code
- [ ] No panic paths outside startup code
- [ ] Structured logging used where needed (`slog`)
- [ ] Scope remained within current phase

## Protocol, Security, And Content Conventions Check

- [ ] WebSocket envelope compatibility preserved (`type`, `schema_version`, `session_id`, `player_id`, `seq`, `payload`) when netcode changed
- [ ] Network limits enforced/configured (frame size, rate limits, heartbeat)
- [ ] TMX conventions respected (`ground`, `collision`, `events`; object naming + `event_id`)
- [ ] Asset conventions respected (sprite frame rules, naming patterns, audio formats)
- [ ] Audio/voice scope respected (text primary, tiny voice barks, no in-game voice chat)
- [ ] Licensing boundary respected (MIT data packs vs ARR runtime/assets)

## Testing Check

- [ ] Unit tests added/updated (`*_test.go` next to logic)
- [ ] Integration tests added/updated (if behavior spans systems)
- [ ] Graphics smoke test run (if render behavior changed)
- [ ] Manual playtest scenario executed

Commands run:

```bash
go test ./...
go vet ./...
staticcheck ./...
go tool govulncheck ./...
```

## Demo Evidence

- Demo artifact (gif/video/screenshot path):
- Logs or output sample:
- Short "what works now" summary:

## Risk And Follow-Up

- Known limitations:
- Deferred tasks (next phase candidate):
- Regression risk areas:

## Sign-Off

- Dev sign-off: [ ]
- Team sign-off: [ ]
- Phase accepted: [ ]
