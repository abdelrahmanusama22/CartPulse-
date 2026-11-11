# Checkout

Primary owner: **Ahmed**.

Server-authoritative quotes and orchestration of order, stock and promotion reservation.

## First task

Agree money/rounding and caller-owned transaction interfaces; atomically reserve stock/discounts and create an idempotent order before remote payment calls.

## Implementation shape

Add handler.go (DTOs/routes), service.go (business rules), repository.go (narrow persistence boundary), types.go (domain inputs/statuses) and meaningful tests as needed. Do not introduce empty methods or successful dummy endpoints. Keep Gin context in the HTTP layer and context.Context in services. Place migrations under db/migrations and named SQL under db/queries/checkout.

Review docs/architecture.md, api/README.md and docs/backend-plan.md before implementing cross-module behavior. This directory currently contains documentation/package scaffolding only.
