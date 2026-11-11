# VAYRO backend

Go + Gin backend scaffold for the VAYRO storefront and staff dashboard. Launch scope: one store. See [team ownership](docs/team.md) and the [complete implementation plan](docs/backend-plan.md).

## Implemented now

- Runnable API (`go run .` or `go run ./cmd/api`) and an explicitly idle worker entry point.
- Gin routing, JSON errors, server-generated request IDs, structured request logging, body-size limits and panic recovery.
- Environment configuration, HTTP timeouts, graceful shutdown, pgx pool wiring and database readiness.
- Health endpoint contract/tests, Go module checks, Docker development setup and GitHub CI.

Business endpoints, authentication, schemas/migrations, generated SQL, the queue, payment/storage/email adapters and integrations are **not implemented**. Their folders and ownership are ready for the team.

## Project layout

```text
main.go                       Convenience API entry point
cmd/api/                      API executable
cmd/worker/                   Worker executable (idle scaffold)
internal/app/                 Dependency wiring and lifecycle
internal/platform/            config, database, HTTP, logging, jobs, storage
internal/identity/            Abdulrhman: users, sessions, staff permissions
internal/catalog/             Waheed: products, variants, categories, collections
internal/inventory/           Waheed: balances, movements, reservations
internal/customers/           Huda: profiles, addresses, staff customer views
internal/cart/                Huda: guest/account carts and merge
internal/discounts/           Huda: promotion rules and redemption reservations
internal/content/             Huda: storefront content and navigation
internal/checkout/            Ahmed: authoritative pricing and checkout orchestration
internal/orders/              Ahmed: order lifecycle, fulfillment and cancellation
internal/payments/            Ahmed: provider integration, events and refunds
internal/store/               Abdulrhman: store/shipping settings
internal/reporting/           Abdulrhman: metrics and exports
internal/dbgen/               Reserved sqlc-generated query output
db/migrations/                Reviewed SQL migration pairs
db/queries/                   Named SQL queries, grouped by module
api/                          OpenAPI and API conventions
docs/                         Architecture, team tasks and detailed plan
tests/integration/            Future cross-module PostgreSQL tests
.github/                      CI, PR and issue templates
```

## Requirements

- Go 1.25 or later; pinned library versions are in go.mod/go.sum.
- Docker with Compose for the optional local PostgreSQL/container setup.
- Git. sqlc and golang-migrate CLIs are introduced/pinned when the first schema/query is added (see db/README.md).

The Go module name is `vayro/backend`; publishing to GitHub does not require changing it. If the team wants `github.com/OWNER/REPO`, update the module line and all `vayro/backend` imports together.

## Quick start: native Go, no database required for liveness

```powershell
go mod download
go run .
```

In a second terminal:

```powershell
Invoke-RestMethod http://localhost:8080/health/live
```

Liveness returns 200. Readiness intentionally returns 503 until DATABASE_URL is configured and PostgreSQL responds. This is not a complete commerce application. Unimplemented business routes return JSON 404.

## Local development with PostgreSQL

```powershell
Copy-Item .env.example .env
docker compose up -d postgres
$env:DATABASE_URL = 'postgres://vayro:vayro_local_only@localhost:5432/vayro?sslmode=disable'
go run ./cmd/api
```

The example assumes unchanged local credentials/port. If you customize .env, set DATABASE_URL to match. Native Go reads process environment variables; it does not load .env automatically. Docker Compose reads .env for interpolation.

Check `http://localhost:8080/health/ready`. It returns 200 after PostgreSQL responds. There are no domain migrations yet and the API does not auto-migrate. Health routes are `/health/*`, not `/api/v1/health/*`.

To run everything in containers instead:

```powershell
docker compose up --build -d
docker compose logs -f api worker
```

The worker logs `worker_scaffold_idle`; it consumes no jobs yet. Stop containers with `docker compose down`. The named PostgreSQL volume persists data. This Compose file is for local development, not production deployment.

## Checks

```powershell
gofmt -w .
go vet ./...
go test ./...
go build ./...
```

On a supported platform with a C toolchain: `go test -race ./...`. GitHub CI runs race tests on Linux. Current tests require no database; real PostgreSQL integration tests are added alongside actual schemas and services.

Optional make targets mirror the same commands. GNU Make is not required on Windows.

## Team workflow

1. Read your module README and docs/team.md.
2. Review the shared money, identity, transaction and API contracts before coding dependent modules.
3. Add real handler.go/service.go/repository.go/types.go files as needed; avoid empty placeholder implementations.
4. Update OpenAPI and migrations with each feature and submit a reviewed PR.

GitHub usernames and a repository URL have not been supplied, so CODEOWNERS and a remote are not guessed. No license is assumed; choose one before public distribution if appropriate.

## Upload to GitHub

Create an empty GitHub repository, then run in this project root:

```powershell
git init -b main
git add .
git commit -m "Initialize VAYRO Go backend scaffold"
git remote add origin https://github.com/YOUR-OWNER/YOUR-REPO.git
git push -u origin main
```

Replace the URL with the real repository. Review staged files before committing. .env is ignored and only .env.example belongs in Git. No repository is created or pushed by this scaffold.
