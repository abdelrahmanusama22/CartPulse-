# Database workflow

PostgreSQL is the source of truth. Migrations and queries start empty so the team can review the ERD before committing domain schemas.

1. Use `golang-migrate` compatible pairs: `000001_create_users.up.sql` and `000001_create_users.down.sql` in `db/migrations/`.
2. Agree migration numbers with other developers. Never modify a merged migration; add a new one.
3. Run migrations explicitly against the intended database. The API does not auto-migrate.
4. Add named queries under `db/queries/<module>/`. Example annotation: `-- name: GetProduct :one`.
5. Install an agreed, pinned sqlc CLI version and run `sqlc generate` after the first schema/query exists.
6. Commit generated files in `internal/dbgen/`; do not edit them manually.
7. Add PostgreSQL integration tests for constraints, stock concurrency and checkout transactions.

Example once migration files exist:

```powershell
migrate -path db/migrations -database $env:DATABASE_URL up
sqlc generate
```

No migration runner or sqlc CLI binary is bundled. Neither command is expected to generate/apply anything before the first migration and query are written. Select/pin their versions in the foundation task. No generated query code is required by the initial API.

Use pgx transactions passed through explicit service boundaries for order + reservation + discount writes. Do not open separate transactions in each participating module. Never hold a database transaction across a payment-provider HTTP call.
