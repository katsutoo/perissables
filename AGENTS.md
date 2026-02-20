---
name: go
description: >
  Production-grade Go backend guidance focused on idiomatic code, clear
  architecture boundaries, reliability, and testability.
---

# Go Skill

## Core Philosophy

Be an honest, insightful programming partner. Challenge code when there are clear improvements in:

1. Efficiency
2. Readability and maintainability
3. Robustness and error handling
4. Scalability
5. Best practices and conventions
6. Security

If a solution is already idiomatic and solid, affirm it. Avoid churn and "refactor for refactor's sake." Explain the "why" behind suggestions.

## Rule Levels

Use this priority model to avoid over-constraining normal Go work.

- MUST: hard requirements for reliability, correctness, and security.
- SHOULD: strong defaults; deviate only with clear context.
- MAY: optional choices based on scale, constraints, or product needs.

## Code Quality Principles

- Prefer `range` loops over index loops unless index arithmetic is required.
- Keep functions focused and easy to scan; split only when it improves clarity.
- Make illegal states unrepresentable through types where practical.
- Handle errors explicitly; no silent failure paths.
- Minimize indirection and deep nesting.
- Favor composition over inheritance-style abstractions.
- Enforce quality gates with tooling, not style debates.

## Default Stack

These are default choices, not universal hard requirements.

| Category | Default | Level | Notes |
|----------|---------|-------|-------|
| Router | Chi | SHOULD | stdlib-friendly, small API surface |
| Database access | `pgx` + `sqlc` | SHOULD | SQL-first, type-safe generated code |
| Migrations | Goose (SQL migrations) | SHOULD | Prefer SQL migrations over Go migrations |
| Logging | `log/slog` | MUST | Structured logging only |
| Validation | `go-playground/validator` | SHOULD | Add custom validators as needed |
| Password hashing | Argon2id | MUST (if passwords) | `golang.org/x/crypto/argon2` |
| IDs | UUIDv7 | SHOULD | Time-sortable IDs; choose one ID strategy per service |
| Integration tests | `testcontainers-go` | SHOULD | Use when external deps are part of behavior |
| CI/CD | GitHub Actions | MAY | Use any equivalent CI platform |
| Caching | Redis or Dragonfly | MAY | Add only when measurable bottleneck exists |
| Message queues | RabbitMQ, RedPanda, NATS | MAY | Add only when async workflows justify it |
| API docs | Swag or ogen | MAY | Use where OpenAPI contract matters |
| Local dev reload | Air | MAY | Live reload in development via `air.toml` |

## Naming Conventions

| Type | Convention | Example |
|------|------------|---------|
| Files | `snake_case` | `user_repository.go` |
| Packages | lowercase, short, concrete | `handler`, `service` |
| Structs | PascalCase | `UserService` |
| Interfaces | PascalCase | `UserRepository`, `EmailSender` |
| Variables | camelCase | `userID`, `donationCount` |
| Constants | PascalCase (exported) | `MaxRetries` |
| Errors | `Err` prefix | `ErrUserNotFound` |
| Handlers | `<Entity>Handler` | `UserHandler` |
| Services | `<Entity>Service` | `UserService` |
| Repositories | `<Entity>Repository` | `UserRepository` |
| Constructors | `New<Type>` | `NewUserService` |
| Context keys | unexported type | `type ctxKey string` |
| Mutexes | `mu` prefix/suffix | `mu sync.Mutex`, `usersMu` |

Getter naming:

- Avoid `Get` for simple field accessors: `user.Name()`, `cfg.Timeout()`.
- Use `Get` when it implies lookup/I/O/side effects: `repo.GetByID(ctx, id)`.

ID and acronym casing:

- Use `ID`, not `Id` (`userID`, `GetByID`).
- Keep initialisms consistent (`httpClient`, `ServeHTTP`, `XMLName`).

## Interfaces and Abstractions

- Define interfaces at the consumer boundary, not in the implementation package by default.
- Start with concrete types; introduce interfaces when a second implementation or a test seam is truly needed.
- Keep interfaces small and behavior-focused.

## Project Layout (Template)

Use as a starting point. Adapt names to your domain.

```text
project/
├── cmd/
│   └── server/
│       └── main.go                  # entrypoint
│
├── internal/
│   ├── handler/
│   │   ├── auth_handler.go
│   │   ├── user_handler.go
│   │   └── health_handler.go
│   │
│   ├── service/
│   │   ├── auth_service.go
│   │   ├── user_service.go
│   │   └── transaction.go           # TxManager
│   │
│   ├── repository/
│   │   └── user_repo.go
│   │
│   ├── middleware/
│   │   ├── auth.go
│   │   ├── logging.go
│   │   ├── cors.go
│   │   ├── security.go
│   │   └── rate_limit.go
│   │
│   ├── model/
│   │   └── user.go
│   │
│   ├── database/
│   │   └── sqlc/                    # generated -- NEVER edit
│   │       ├── db.go
│   │       ├── models.go
│   │       └── querier.go
│   │
│   ├── apierror/
│   │   └── errors.go
│   │
│   ├── validator/
│   │   └── validator.go
│   │
│   └── client/
│       └── email_client.go
│
├── sql/
│   ├── sqlc.yaml
│   ├── schema/
│   │   ├── 20260219120000_create_users.sql
│   │   └── 20260219121000_create_sessions.sql
│   └── queries/
│       ├── users.sql
│       └── sessions.sql
│
├── test/
│   ├── integration/
│   │   └── user_test.go
│   └── testutil/
│       └── fixtures.go
│
├── .github/workflows/ci.yml
├── Dockerfile
├── docker-compose.yml
├── air.toml
├── justfile
├── go.mod
├── go.sum
└── README.md
```

## Request Flow

```text
HTTP Request
    -> Handler   (transport parsing, response mapping)
    -> Service   (business rules, orchestration, transactions)
    -> Repository(DB queries + mapping)
    -> Database
```

### Layer Responsibilities

| Layer | Does | Does NOT | Test Focus |
|-------|------|----------|------------|
| Handler | Parse/validate transport fields, call service, write response | Domain business logic, raw SQL | Serialization, status codes, request validation |
| Service | Business rules, transactions, orchestration | HTTP concerns, raw SQL | Rules, rollback behavior, error paths |
| Repository | SQL queries, persistence mapping | Domain policy, HTTP concerns | Query correctness, mapping, DB errors |

## Anti-Patterns to Avoid

No god handlers/services:

```text
bad:  handler/sports.go (1600+ lines, many unrelated entities)
good: handler/team_handler.go, handler/player_handler.go
```

Constructor injection only:

```go
// bad
h.Rating = NewRatingHandler(...)
h.Rating.SetCache(cache)

// good
ratingService := service.NewRatingService(repo, cache, notifier)
h.Rating = NewRatingHandler(ratingService)
```

Never ignore errors:

- No `data, _ := fn()` unless there is explicit, documented justification.
- No `fmt.Println` for operational logs; use `slog`.
- No `panic` outside startup/bootstrapping paths in `main`.

## Error Handling

| Situation | Use |
|-----------|-----|
| Define domain error | `var ErrX = errors.New("x")` |
| Return domain error | `return ErrUserNotFound` |
| Wrap with context | `fmt.Errorf("doing X: %w", err)` |
| Check specific error | `errors.Is(err, ErrX)` |
| Check error type | `errors.As(err, &target)` |

Rules:

- Always use `%w` when wrapping errors.
- Keep transport-level error mapping in handlers, not in services/repositories.
- Add stable machine-readable error codes for API consumers.

## Logging with slog

- Use structured fields (`slog.String`, `slog.Int`, `slog.Any`, etc.).
- Include request-scoped keys where available (`request_id`, `user_id`, `trace_id`).
- Log errors with explicit `err` field.

```go
slog.Info("user_created",
    slog.String("user_id", user.ID.String()),
    slog.Duration("duration", time.Since(start)),
)
```

## Validation

- Validate request DTOs at handler boundaries.
- Enforce domain invariants in service layer.
- For JSON decoding, prefer `DisallowUnknownFields` where strict contracts are desired.
- Keep custom validators in one package and reuse via tags.

Example tag: `validate:"required,email,uuidv7,min=3,max=50"`.

## Password Hashing (Argon2id)

- Use Argon2id for password storage.
- Store algorithm parameters alongside the hash for future tuning.
- Tune memory/time cost to match hardware and latency budget.

Starting points (review periodically against current OWASP guidance):

- Memory-heavy profile: `m=47104, t=1, p=1`
- CPU-heavy profile: `m=19456, t=2, p=1`

## Context Handling

- `context.Context` is the first parameter, named `ctx`.
- Never store context in structs.
- Avoid `context.Background()` in request paths; use it only at process roots.
- Propagate context to DB, HTTP, cache, and queue calls.
- Check cancellation in long-running loops and worker operations.

## Concurrency and Goroutines

- Every started goroutine must have a clear owner and shutdown path.
- Use `errgroup.WithContext` for fan-out work that should cancel together.
- Avoid unbounded goroutine creation; use worker pools or backpressure.
- Close channels from the sender side only.
- Run race detection regularly: `go test -race ./...`.

## HTTP Standards and Safety

- Use `net/http` status constants (`http.StatusOK`, not numeric literals).
- Set server timeouts explicitly (`ReadHeaderTimeout`, `ReadTimeout`, `WriteTimeout`, `IdleTimeout`).
- Limit request body size for JSON endpoints (`http.MaxBytesReader`).
- Use explicit client timeouts and sane transport settings for outbound calls.
- Keep response envelopes consistent per API surface.

Recommended JSON envelope patterns:

- Success: `{ "data": T }`
- Error: `{ "error": string, "code": string, "details": map }`
- Paginated: `{ "data": []T, "total": int, "page": int, "per_page": int }`

## Testing Standards

### Unit Tests

- Prefer table-driven tests when one behavior has many input/output cases.
- Cover happy path, error path, and edge cases.
- Use `t.Run` names that describe behavior, not only inputs.
- Use `t.Parallel()` for independent tests.
- Use `t.Helper()` for helper functions.
- Use `t.Cleanup()` for shared teardown and `defer` for local, immediate cleanup.

### Integration Tests

- Use real dependencies (often via `testcontainers-go`) when behavior depends on them.
- Isolate state with per-test schema/database or cleanup in `t.Cleanup()`.
- Keep integration test setup deterministic and CI-friendly.

### Other Test Types

- Add fuzz tests for parser/decoder and boundary-heavy logic (`go test -fuzz`).
- Add benchmarks for hot paths (`go test -bench=.`).
- Use golden tests where output shape matters and is intentionally stable.

### Contract Testing

For external APIs, prefer record/replay (`go-vcr`) or contract tests (Pact) when API drift can break CI.

## Goose Migrations

- Location: `sql/schema/`.
- Naming (default): `YYYYMMDDHHMMSS_description.sql`.
- Include both `-- +goose Up` and `-- +goose Down`.
- Keep migrations reversible whenever possible; document any intentional one-way migration.
- Pick one naming convention per repository and stay consistent.

## sqlc Configuration

- Config path: `sql/sqlc.yaml`.
- Use `pgx/v5` as `sql_package`.
- Common options: `emit_json_tags`, `emit_interface`.
- Output path: `internal/database/sqlc` (generated code, never hand-edit).

## Code Quality Checklist

Before committing:

- [ ] Business logic lives in service layer, not handlers.
- [ ] No ignored errors without explicit justification.
- [ ] Structured logging with `slog` (no `fmt.Println`).
- [ ] Input validation present for external-facing endpoints.
- [ ] Context propagated through all I/O boundaries.
- [ ] Goroutines have cancellation/shutdown path.
- [ ] Transaction boundaries are explicit and tested.
- [ ] Table-driven tests added where they improve coverage clarity.
- [ ] Integration tests added/updated when external deps are involved.
- [ ] Migrations include valid Up/Down and match repo naming convention.
- [ ] `sqlc` regenerated when schema/queries change.
- [ ] Never manually edit `internal/database/sqlc/`.
- [ ] `go test ./...` passes.
- [ ] `go test -race ./...` passes for non-trivial concurrent code.
- [ ] `go vet ./...` passes.
- [ ] `staticcheck ./...` passes.
- [ ] `govulncheck ./...` reviewed for known CVEs.
- [ ] Graceful shutdown handles SIGINT/SIGTERM.
- [ ] Production DB URLs use secure TLS settings (for Postgres, `sslmode=require` or stronger).
- [ ] No `panic` in request path.
