# Customers

Primary owner: **Huda**.

Customer profiles, owned addresses, staff customer views and notes.

## First task

Define customer/account links and ownership tests; never attach guest orders solely because an email string matches.

## Implementation shape

Add handler.go (DTOs/routes), service.go (business rules), repository.go (narrow persistence boundary), types.go (domain inputs/statuses) and meaningful tests as needed. Do not introduce empty methods or successful dummy endpoints. Keep Gin context in the HTTP layer and context.Context in services. Place migrations under db/migrations and named SQL under db/queries/customers.

Review docs/architecture.md, api/README.md and docs/backend-plan.md before implementing cross-module behavior. This directory currently contains documentation/package scaffolding only.
