# Discounts

Primary owner: **Huda**.

Fixed/percentage discounts, eligibility, usage reservations and redemptions.

## First task

Agree rounding, one-code policy and transactional usage limits with checkout; quote validation is not redemption.

## Implementation shape

Add handler.go (DTOs/routes), service.go (business rules), repository.go (narrow persistence boundary), types.go (domain inputs/statuses) and meaningful tests as needed. Do not introduce empty methods or successful dummy endpoints. Keep Gin context in the HTTP layer and context.Context in services. Place migrations under db/migrations and named SQL under db/queries/discounts.

Review docs/architecture.md, api/README.md and docs/backend-plan.md before implementing cross-module behavior. This directory currently contains documentation/package scaffolding only.
