# Cart

Primary owner: **Huda**.

Owned guest/account carts, item updates, cart merge and price previews.

## First task

Define opaque guest ownership tokens, versioned cart changes and deterministic merge behavior; cart activity does not reserve stock.

## Implementation shape

Add handler.go (DTOs/routes), service.go (business rules), repository.go (narrow persistence boundary), types.go (domain inputs/statuses) and meaningful tests as needed. Do not introduce empty methods or successful dummy endpoints. Keep Gin context in the HTTP layer and context.Context in services. Place migrations under db/migrations and named SQL under db/queries/cart.

Review docs/architecture.md, api/README.md and docs/backend-plan.md before implementing cross-module behavior. This directory currently contains documentation/package scaffolding only.
