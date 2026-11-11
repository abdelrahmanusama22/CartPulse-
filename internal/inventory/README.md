# Inventory

Primary owner: **Waheed**.

Single-location stock balances, immutable movements and order reservations.

## First task

Agree Reserve/Commit/Release transaction interfaces with Ahmed; specify available = on_hand - reserved and last-unit concurrency tests.

## Implementation shape

Add handler.go (DTOs/routes), service.go (business rules), repository.go (narrow persistence boundary), types.go (domain inputs/statuses) and meaningful tests as needed. Do not introduce empty methods or successful dummy endpoints. Keep Gin context in the HTTP layer and context.Context in services. Place migrations under db/migrations and named SQL under db/queries/inventory.

Review docs/architecture.md, api/README.md and docs/backend-plan.md before implementing cross-module behavior. This directory currently contains documentation/package scaffolding only.
