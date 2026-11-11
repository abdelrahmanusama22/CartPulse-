# Architecture and module rules

One store initially. One modular Go backend and PostgreSQL database, with an API and a separate worker executable. Business modules communicate through Go interfaces, not internal HTTP calls.

For each module add these files as real work begins:

```text
handler.go         HTTP request/response DTOs, binding and route registration
service.go         Business rules, orchestration and transaction boundaries
repository.go      Narrow persistence interface or adapter around generated SQL
types.go           Domain types, inputs, statuses and errors
*_test.go          Meaningful unit/handler tests
```

Empty implementations and fake business endpoints are intentionally absent. `doc.go` and module READMEs preserve the initial folders in Git. Split a module into subpackages only when its actual size warrants it.

- `internal/app`: dependency construction and executable lifecycle.
- `internal/platform`: configuration, HTTP conventions, logging, database, storage/jobs adapters.
- `internal/<business module>`: explicit business ownership.
- `db`: versioned schema migrations and named SQL queries.
- `internal/dbgen`: generated SQL access, not HTTP response types.
- `api`: OpenAPI and transport conventions.
- `tests/integration`: real cross-module/database checks.

Pass context.Context through services and queries. Handle ownership and permission checks for every object. Use separate staff/customer session audiences, not a trusted client-supplied role. Configure CSRF and exact CORS origins when authentication is implemented.

Stock reservation, order creation and bounded discount usage must share a caller-owned transaction. Payments, email and storage calls occur outside that transaction; use a durable outbox and idempotent reconciliation. Currency and rounding rules are shared. Orders snapshot historical prices, SKU/options, taxes, discounts and addresses.

The current readiness endpoint checks connectivity only; add migration/critical-dependency checks as implemented. The worker is explicitly idle until durable jobs are implemented. No auth, payment or commerce behavior is claimed by this scaffold.
