# API contracts

`openapi.yaml` describes implemented health endpoints only. The business route inventory lives in `docs/backend-plan.md`; those routes are planned, not implemented.

Business prefix: `/api/v1`.

- `/auth/*`: customer identity and invitation acceptance.
- `/account/*`: authenticated customer profile, addresses, and history.
- `/storefront/*`: public catalog/content and ownership-checked cart/checkout/order actions.
- `/admin/auth/*`: staff session creation/revocation.
- `/admin/*`: staff-only operations with explicit permission checks.
- `/webhooks/payments/:provider`: raw-body signature verification and durable event deduplication.

Do not mount admin/account routes until their auth middleware exists. Staff/customer audiences must be checked separately. Cookie-authenticated mutations need CSRF protections. Unimplemented endpoints currently return the normal JSON 404 response.

JSON names use snake_case. Document pagination, status codes, validation errors and authorization per operation. Error shape: `{"error":{"code":"...","message":"...","request_id":"..."}}`. Add field errors to this contract when validation is implemented. Use `platform/http.WriteError` for consistent transport errors.

Use `ShouldBindJSON` with explicit errors; service validation still enforces business rules. Server-calculated money uses integer minor units plus currency. Checkout, payment creation, refunds and inventory adjustments need durable idempotency handling before release.

When wiring modules, create explicit Gin route groups in `internal/platform/http/router.go` and pass constructed services into handlers. Keep Gin context out of business service interfaces. Introduce the agreed principal/transaction types with the first real identity/checkout implementations.
