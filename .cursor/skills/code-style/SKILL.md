---
name: code-style
description: >
  Defines Go package layout and coding patterns for models, presenters,
  repositories, providers, services, and controllers (Echo, go-pg, slog, Sentry).
  Use when the prompt-framework skill loads context for writing or editing code,
  or when creating a new domain feature end-to-end.
disable-model-invocation: true
---

# Code Style

Example-driven conventions for this Go API. New domains (users, sessions, …)
should mirror this shape. Architecture and dependency rules live in
[architecture](../architecture/SKILL.md); full session walkthrough in
[examples.md](examples.md).

## Package layout

```
internal/
  controllers/
  presenters/
  models/
  providers/
  repositories/
  services/
```

One file per domain per package (e.g. `session.go`).

## Stack

| Concern | Choice |
|---------|--------|
| HTTP | Echo (`echo.Group`, middleware on routes) |
| DB | go-pg (`*pg.DB`, `Model`, `WherePK`, `SelectAndCount`) |
| Logging | `slog` with structured fields |
| Error tracking | `sentry.CaptureError` from controllers |
| Domain errors | `var ErrXxx = fmt.Errorf(...)` on models |
| Not-found | Map `pg.ErrNoRows` → domain error in providers |

## Models

- Struct with `json` + `pg` tags; PK via `pg:",pk"`
- Sentinel errors beside the type
- Filters/pagination as `XxxParams` (`Page`, `PerPage`, filter fields)
- IDs are `int64` unless the domain already uses another type

## Presenters

- Request/query types with `query` / `json` tags for Echo `Bind`
- Mappers: `ToParams()`, `ToCreateXxx()`, `ToPatchXxx()`
- Pagination defaults in `ToParams()`: page `1`, per_page `20`
- Patch returns `(*models.T, []string)` — entity + columns to update
- List responses: `Total`, `Page`, `PerPage`, `Data`

## Repositories

- Interface in `repositories`: `Create`, `FindAll`, `FindByID`, `Patch`, `Update`, `Delete`
- All methods take `context.Context`
- `FindAll` → `(count int, items []T, err error)`
- `Patch` → `columns ...string`

## Providers

- Unexported `xxxRepository` with `*pg.DB`
- `NewXxxRepository(db *pg.DB) repositories.XxxRepository`
- Implement the repository interface
- `slog.Error` on failure before returning

## Services

- Exported `XxxService` interface; unexported `xxxService`
- `NewXxxService(repo repositories.XxxRepository) XxxService`
- Set `CreatedAt` / `UpdatedAt` on create/update; on patch append `"updated_at"`
- Thin over the repo unless real domain logic is required
- Method names: `GetByID` / `GetAll` at service; `FindByID` / `FindAll` at repo

## Controllers

- Exported `XxxController` with `Routes(g *echo.Group)`
- Unexported `xxxController`; `NewXxxController(svc services.XxxService)`
- Routes + auth middleware in `Routes`
- Swagger godoc on each handler
- Flow: bind presenter → service → `JSON` / `NoContent`
- Bind/parse errors → `400`; service errors → `500` + Sentry
- Empty lists → `[]models.T{}` (never null JSON)

## Generating code

1. Respect [architecture](../architecture/SKILL.md) boundaries
2. Copy the session pattern in [examples.md](examples.md); rename for the domain
3. Keep IDs, tags, and method names consistent across all six layers
4. Read examples only when implementing or substantially extending a domain
