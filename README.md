# go-take-home-test

Go implementation of the Healthtech-1 form ingest / transform pipeline.

The system accepts registration forms from an unreliable third party, persists them, geocodes addresses, transforms them into the FORM-BOT schema, and notifies the team by email — with idempotency, retries, and a recovery path when a step fails.

## Requirements (brief)

| Requirement | Approach |
|---|---|
| Ingest via `/ingest` into a real database | SQLite + Bun; `ingested_forms` |
| Conform to the agreed provider schema | Validation on `models.IngestedForm` |
| Attach lat/long from postcode | Mock Ideal Postcodes lookup in transform worker |
| Transform to FORM-BOT schema | `ToTransformedForm` + camelCase presenter for the email body |
| Recover after a failed step | `POST /retry` re-queues transform or send-to-bot |
| At-least-once delivery from provider | Content fingerprint (`UNIQUE`) dedupes identical payloads |
| Never give FORM-BOT the same form twice | `sent_to_bot` flag; skip if already sent |
| Notify `happyforms@bots.com` after success | Mock email after transform; queue retries until send succeeds |

Schemas live in:

- Ingest: [`internal/models/ingested.go`](internal/models/ingested.go)
- Transformed / FORM-BOT: [`internal/models/transformed.go`](internal/models/transformed.go), outbound shape in [`internal/presenters/transform.go`](internal/presenters/transform.go)

Example payloads: [`forms/examples/`](forms/examples/).

---

## Quick start

### Prerequisites

- Go **1.25+**
- Dependencies are vendored under `vendor/` (no network required for build)

### Run the server

From the repo root:

```bash
cd go-take-home-test
go run ./cmd/server
```

Server listens on **`:8080`** by default (`PORT` overrides).

On first start the app:

1. Creates `./data/db.sqlite` (directory included)
2. Applies SQL migrations from `internal/migrations/`
3. Writes `./migration_version.txt`

```text
Server is running on http://localhost:8080
⇨ http server started on [::]:8080
```

**Port and queue must match.** The in-process queue POSTs workers at `http://localhost:8080/workers/...` (see `BaseURL` in `internal/app/app.go`). If you change `PORT`, keep that base URL in sync or workers will fail.

### Ingest an example form

```bash
curl -s -X POST http://localhost:8080/ingest \
  -H 'Content-Type: application/json' \
  -d @forms/examples/person_one.json | jq
```

### Retry a failed form

```bash
curl -s -X POST http://localhost:8080/retry \
  -H 'Content-Type: application/json' \
  -d '{"session_id":"YOUR_SESSION_ID"}' | jq
```

You can also pass `"ingested_id": 1`.

---

## Configuration

| Variable | Default | Meaning |
|---|---|---|
| `PORT` | `8080` | HTTP listen port |
| `HEALTH_MOCK_RELIABLE` | unset | Mock postcode + email behaviour (see below) |

### Mock reliability (`HEALTH_MOCK_RELIABLE`)

Used by the mock Ideal Postcodes and SendGrid providers:

| Value | Behaviour |
|---|---|
| `1` | Always succeed (no artificial delay) — best for local demos |
| `0` | Always fail — useful to exercise failure / `/retry` |
| unset | ~95% success with ~1s latency |

App defaults (SQLite path, migrations, queue retries) are in [`internal/app/app.go`](internal/app/app.go).

---

## HTTP API

### Public

| Method | Path | Purpose |
|---|---|---|
| `POST` | `/ingest` | Validate, persist (or reuse by fingerprint), queue transform |
| `POST` | `/retry` | Re-queue transform or send-to-bot for an existing form |

### Internal workers (invoked by the mock queue)

| Method | Path | Purpose |
|---|---|---|
| `POST` | `/workers/transform` | Map fields, geocode, write `transformed_forms`, queue send-to-bot |
| `POST` | `/workers/send-to-bot` | Email FORM-BOT payload to `happyforms@bots.com`, set `sent_to_bot` |
| `POST` | `/workers/ingest` | Same handler as `/ingest` (mounted for symmetry) |
| `POST` | `/workers/retry` | Same handler as `/retry` |

The queue implementation POSTs JSON to `{BaseURL}/workers/{queueName}` with exponential backoff (`QueueMaxAttempts`, `QueueInitialBackoff`).

---

## Pipeline

```text
POST /ingest
    │
    ├─ bind + fingerprint (SHA-256 of full ingest JSON)
    ├─ validate against agreed schema
    ├─ GetOrCreate by fingerprint  ──duplicate──► ack / re-queue if incomplete
    │
    └─ queue "transform"
            │
            ├─ load ingested row
            ├─ map → transformed schema (name split, gender map, DOB parse)
            ├─ postcode → lat/long
            ├─ insert transformed_forms (UNIQUE ingested_form_id)
            │
            └─ queue "send-to-bot"
                    │
                    ├─ if sent_to_bot already → no-op
                    ├─ email happyforms@bots.com (JSON body)
                    └─ mark sent_to_bot + transform_log
```

`/retry` inspects state:

- No transformed row → re-queue **transform**
- Transformed but not sent → re-queue **send-to-bot**
- Already `sent_to_bot` → acknowledge without sending again

---

## Architecture

Layered layout under `internal/` (controllers → services → repository interfaces ← providers). Wired in [`internal/app/app.go`](internal/app/app.go).

```text
cmd/server/          process entrypoint
internal/
  app/               composition root + Echo setup + E2E tests
  controllers/       HTTP adapters (bind, status codes, orchestration of services)
  presenters/        request/response DTOs (ingest bind, bot email shape, worker payloads)
  services/          domain logic (create, get-or-create, patch, send email, queue)
  repositories/      persistence / outbound interfaces only
  providers/         Bun/SQLite, mock queue, mock email, mock postcode
  models/            entities, validation, transform mapping, domain errors
  migrations/        SQL schema up/down
forms/examples/      sample provider payloads
```

| Layer | Responsibility |
|---|---|
| **Controllers** | Bind presenters, call services, return JSON. Keep HTTP concerns here. |
| **Presenters** | Inbound bind types (`IngestedForm`), worker messages, FORM-BOT camelCase email body. |
| **Services** | Business rules: timestamps, `GetOrCreate` by fingerprint, patches, thin wrappers over repos. |
| **Repositories** | Interfaces (`FindByFingerprint`, `Create`, `SendEmail`, …) — no Bun/HTTP imports. |
| **Providers** | Concrete adapters: Bun CRUD, file-based migrations, HTTP queue, flaky mocks. |
| **Models** | Stored entities, query params, `Validate`, `ToTransformedForm`, shared errors. |

Dependency direction:

```text
Controllers → Services → Repositories (interfaces)
                              ↑
                         Providers (implementations)
```

---

## Database schema

SQLite file: `./data/db.sqlite`  
Migration source: [`internal/migrations/20260720041429_init.up.sql`](internal/migrations/20260720041429_init.up.sql)  
Version marker: `./migration_version.txt`

| Table | Role |
|---|---|
| `ingested_forms` | Raw provider payload (+ `fingerprint UNIQUE`, `status`) |
| `transformed_forms` | FORM-BOT-ready row (`ingested_form_id UNIQUE`, `sent_to_bot`, lat/long) |
| `transform_logs` | Outcome of send-to-bot / email attempts |

Design notes:

- **Fingerprint** = SHA-256 of the marshalled ingest presenter; identical redeliveries share one row.
- **`ingested_form_id UNIQUE`** prevents duplicate transformed rows for the same ingest.
- **`sent_to_bot`** is the gate for “never send the same form twice.”
- Address on ingest is stored as JSON text; transformed form flattens address lines + coordinates.

---

## Design decisions (short)

1. **Dedupe by content fingerprint**, not `session_id` alone — provider may reuse or change session IDs; identical full payloads collapse to one form.
2. **Synchronous mock queue with retries** — enough to model failure without running Redis/SQS; workers stay ordinary HTTP handlers.
3. **Idempotent transform / send-to-bot** — safe to re-queue after partial failure.
4. **Public `/retry`** — recover after a code fix without re-ingesting from the provider.
5. **Strict validation on ingest** — reject schema violations with `400`; extra JSON fields are ignored (tolerant of additive drift).

---

## Testing

### Run tests

```bash
cd go-take-home-test
go test ./...
```

Verbose E2E package:

```bash
go test ./internal/app/ -v
```

Tests use `t.TempDir()` for SQLite and `migration_version.txt`, so they **do not** write `./data/db.sqlite`. Each case spins up `httptest.Server` with `BaseURL` pointed at itself so the mock queue can call `/workers/*`.

### What the tests cover

| Test | What it asserts |
|---|---|
| `TestIngestEndToEndReliable` | Happy path for `person_one/two/three.json` with reliable mocks (`HEALTH_MOCK_RELIABLE=1`) |
| `TestIngestEndToEndUnreliableMocksFail` | With mocks forced to fail (`=0`), ingest returns `500` after queue retries exhaust |
| `TestIngestEndToEndUnreliableMocksEventuallySucceed` | Default flaky mocks (~95%) eventually succeed given enough attempts |
| `TestIngestDeduplicatesByFingerprint` | Identical payload twice → same `id` and `fingerprint` |
| `TestRetryAfterFailedTransform` | Failed ingest under bad mocks, then `/retry` succeeds after mocks are made reliable |

Example fixtures: [`forms/examples/`](forms/examples/).

---

## Project layout (top level)

```text
go-take-home-test/
├── cmd/server/           main
├── forms/examples/       sample JSON forms
├── internal/             application code (see Architecture)
├── vendor/               vendored modules
├── data/                 created at runtime (db.sqlite) — gitignored if present
├── migration_version.txt created at runtime
├── go.mod
├── go.sum
└── README.md
```

---

## Troubleshooting

| Symptom | Likely cause |
|---|---|
| `listen tcp :8080: bind: address already in use` | Another process (or previous `go run`) still holds the port |
| Queue / transform `500`, postcode errors | Unset or `HEALTH_MOCK_RELIABLE=0`; set `=1` for local demos |
| Workers not called after ingest | `PORT` ≠ host in app `BaseURL` (default both `8080`) |
| No `data/db.sqlite` after tests | Expected — tests use temp dirs; run the server to create the file |

---

## Stack

- [Echo](https://echo.labstack.com/) — HTTP
- [Bun](https://bun.uptrace.dev/) + SQLite — persistence
- Standard library `log/slog` — logging
