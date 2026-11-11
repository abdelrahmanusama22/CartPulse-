# Contributing

1. Read README.md, docs/team.md, docs/architecture.md and your module README.
2. Create a focused branch, e.g. `feat/catalog-products` or `fix/cart-ownership`.
3. Agree the OpenAPI/schema/permission changes before implementing dependent endpoints.
4. Implement the service rules and persistence, then wire HTTP handlers. Never mount unprotected admin/account endpoints.
5. Add meaningful tests; update API and migration docs with the change.
6. Run `gofmt -w .`, `go vet ./...`, `go test ./...` and `go build ./...`.
7. Open a pull request and request another module owner's review.

Keep secrets and .env out of Git. Pin new dependencies in go.mod/go.sum. Avoid business logic in main.go, Gin handlers or generated queries. Do not modify merged migrations. Stock/order/payment changes require concurrency/idempotency evidence and review from an affected owner.

Do not add a fake health success, dummy auth principal or in-memory payment/stock implementation to make unfinished endpoints appear complete. The first implementation is a scaffold; follow the work packages in docs/backend-plan.md.
