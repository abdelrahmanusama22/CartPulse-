# Store

Primary owner: **Abdulrhman**.

Single-store settings, public configuration, shipping zones and shipping rates.

## First task

Agree currency, timezone and shipping rules; expose only public settings to the storefront and never provider secrets.

## Implementation shape

Add handler.go (DTOs/routes), service.go (business rules), repository.go (narrow persistence boundary), types.go (domain inputs/statuses) and meaningful tests as needed. Do not introduce empty methods or successful dummy endpoints. Keep Gin context in the HTTP layer and context.Context in services. Place migrations under db/migrations and named SQL under db/queries/store.

Review docs/architecture.md, api/README.md and docs/backend-plan.md before implementing cross-module behavior. This directory currently contains documentation/package scaffolding only.
