---
name: architecture
description: >
  Defines the Go app's DDD layered structure under internal/ (controllers,
  services, repositories, providers, models, presenters), dependency direction,
  and DTO rules. Use when the prompt-framework skill loads context, or when
  designing, implementing, or reviewing application structure and layer boundaries.
disable-model-invocation: true
---

# Architecture

Layered Domain-Driven Design under `internal/`. For how each layer is coded,
see [code-style](../code-style/SKILL.md).

## Package layout

```
internal/
  controllers/   # HTTP adapters (Echo)
  presenters/    # Request/response DTOs for controllers
  services/      # Domain / business logic
  repositories/  # Persistence interfaces
  providers/     # Interface implementations (go-pg, etc.)
  models/        # Domain entities and query params
```

One domain file per package (e.g. `session.go` in each layer that needs it).

## Dependency direction

```
Controllers → Services → Repositories (interfaces)
                              ↑
                         Providers (implementations)
```

- Controllers depend on **services** and **presenters** only
- Services depend on **repository interfaces** and **models** — never providers
- Providers implement repositories and talk to technology (Postgres, buckets, queues)
- Wire concrete providers at composition root (app startup), not inside services

## Layers

| Layer | Role |
|-------|------|
| Controllers | User-facing outbound adapters (routes, bind, HTTP status) |
| Services | Domain logic; timestamps, orchestration, rules |
| Repositories | Persistence contracts (swap implementations via DI) |
| Providers | Technology-specific repository implementations |

## DTOs

| Type | Used by | Purpose |
|------|---------|---------|
| Models | Providers, repositories, services | Internal entities and query params |
| Presenters | Controllers | Inbound bind types and outbound client shapes |

### Presenter rules

- Controllers bind HTTP input into presenters, then map to models via
  `ToParams` / `ToCreateXxx` / `ToPatchXxx`
- Use presenters for composite responses (e.g. paginated lists)
- Single-entity responses may return models when no dedicated presenter exists
- Prefer presenters when the client shape should diverge from storage/models

## Providers

- Adapters for a specific technology
- Implement a repository interface
- Map driver errors (e.g. `pg.ErrNoRows`) to domain errors from models
- Use **models** as DTOs

## Repositories

- Interfaces only — no technology imports
- Encapsulate providers so services stay storage-agnostic

## Services

- Semantically scoped domain logic (one service per domain area)
- May wrap one repository or compose other services
- Use **models** as DTOs; no Echo types, no presenters

## Controllers

- Register routes; apply auth middleware
- Bind → service call → JSON / NoContent
- Keep business logic minimal; push it to services when it grows
- Log with `slog`; report unexpected errors to Sentry
