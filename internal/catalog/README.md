# Catalog

Primary owner: **Waheed**.

Products, variants/options, categories, collections, verified product media and public catalog reads.

## First task

Map product editor fields to an ERD and OpenAPI contract; enforce unique SKU/options and draft/published/archive behavior.

## Implementation shape

Add handler.go (DTOs/routes), service.go (business rules), repository.go (narrow persistence boundary), types.go (domain inputs/statuses) and meaningful tests as needed. Do not introduce empty methods or successful dummy endpoints. Keep Gin context in the HTTP layer and context.Context in services. Place migrations under db/migrations and named SQL under db/queries/catalog.

Review docs/architecture.md, api/README.md and docs/backend-plan.md before implementing cross-module behavior. This directory currently contains documentation/package scaffolding only.
