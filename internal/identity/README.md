# Identity

Primary owner: **Abdulrhman**.

Customer/staff identities, sessions, recovery, invitations and fixed permissions.

## First task

Define principal/session audiences and role permissions; implement identity schema and auth middleware before mounting protected routes.

## Implementation shape

Add handler.go (DTOs/routes), service.go (business rules), repository.go (narrow persistence boundary), types.go (domain inputs/statuses) and meaningful tests as needed. Do not introduce empty methods or successful dummy endpoints. Keep Gin context in the HTTP layer and context.Context in services. Place migrations under db/migrations and named SQL under db/queries/identity.

Review docs/architecture.md, api/README.md and docs/backend-plan.md before implementing cross-module behavior. This directory currently contains documentation/package scaffolding only.
