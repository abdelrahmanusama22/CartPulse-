# Payments

Primary owner: **Ahmed**.

Provider attempts, signed webhooks, reconciliation, COD confirmation and refunds.

## First task

Select provider and approved policies; use hosted/tokenized checkout, durable event deduplication and stable idempotency references; never trust browser success redirects.

## Implementation shape

Add handler.go (DTOs/routes), service.go (business rules), repository.go (narrow persistence boundary), types.go (domain inputs/statuses) and meaningful tests as needed. Do not introduce empty methods or successful dummy endpoints. Keep Gin context in the HTTP layer and context.Context in services. Place migrations under db/migrations and named SQL under db/queries/payments.

Review docs/architecture.md, api/README.md and docs/backend-plan.md before implementing cross-module behavior. This directory currently contains documentation/package scaffolding only.
