# VAYRO backend implementation plan

Planning baseline: 10 November 2026. Stack confirmed by the user: Go with Gin. Launch model confirmed by the user: one store initially. Team: Abdulrhman, Waheed, Huda, Ahmed.

This is an implementation plan, not an existing backend audit. Ownership is proposed according to module boundaries; adjust it to the developers' experience and availability before estimating individual commitments.

## 1. Product scope and decisions

The Figma screen list includes dashboard overview, products/product editor, inventory, orders, customers, discounts, and store settings, with desktop/tablet/mobile layouts. The storefront overview shows home/catalog/product, bag, checkout, and confirmation layouts. Responsive layouts use the same backend endpoints.

The plan covers the backend needed for those journeys. Exact field requirements still need a screen-by-screen contract review; this review has not verified every design field or interaction.

### Proposed launch scope

- One store, one stock location, one selling currency; confirm the currency rather than infer it from mock data.
- Physical products with size/color variants, SKU, price, images, categories, and collections.
- Published storefront catalog, search/filter/sort, home content, navigation, cart, and checkout.
- Guest checkout plus optional customer accounts; guest orders remain visible only through a scoped opaque access token or verified ownership.
- Staff login, fixed permissions, product and stock management, orders, customers, discounts, basic reporting, and settings.
- One online payment provider if chosen, plus cash on delivery only if the business requires it.
- Configurable shipping zones/rates and manual shipment tracking initially.
- Full-order cancellation before shipment, full refunds where applicable, and separate stock-return actions.
- Transactional email for account recovery and order updates.

### Defer unless required for launch

Multiple stores/tenants, subscriptions, multiple warehouses, multiple currencies, loyalty, reviews, wishlists, advanced returns/exchanges, partial fulfillment/refunds, live carrier integrations, custom role builders, complex promotion stacking, and marketing automation.

### Decisions to resolve in week 1

| Decision | Proposed starting point | Why it matters | Coordinator |
|---|---|---|---|
| Currency and selling region | One explicitly chosen currency/region | Money precision, shipping, tax, payment availability | Ahmed |
| Payment provider and COD | One provider; COD only if needed | Credentials, webhook behavior, settlement and test scenarios | Ahmed |
| Taxes and displayed prices | Business-approved tax-inclusive/exclusive rules | Checkout totals and invoices | Ahmed |
| Shipping | One location; zone-based rates; manual tracking | Checkout eligibility and fulfillment | Abdulrhman |
| Guest/customer accounts | Both, with guest ownership token | Auth, checkout, history and customer matching | Huda + Abdulrhman |
| Staff permissions | Owner, manager, fulfillment | Sensitive mutations and customer visibility | Abdulrhman |
| Refund/return policy | Full refund initially; stock restocked explicitly | Accounting and inventory consistency | Ahmed + Waheed |
| Product/content fields | Map every design field to API/schema | Avoid UI/backend mismatch | Waheed + Huda |
| Domain/deployment arrangement | Storefront and dashboard proxy their API paths | Cookie scope, CSRF and CORS | Abdulrhman |

These choices are planning assumptions until agreed. Payment, tax and shipping integrations cannot be considered production-ready before the actual choices are validated.

## 2. Architecture and stack

Build a modular monolith: one Go repository, one API application, one worker executable, and one PostgreSQL database. Modules call each other through explicit Go service interfaces. No internal HTTP calls between modules.

```text
Storefront browser ─┐
                   ├─ HTTPS / reverse proxy ─ Gin API ─ PostgreSQL
Dashboard browser ─┘                              ├─ object storage / CDN
                                                └─ payment provider
PostgreSQL jobs + outbox ─ Go worker ─ email / payment reconciliation / exports
```

| Concern | Proposed choice | Implementation rule |
|---|---|---|
| HTTP API | Go + Gin | Pin supported versions when creating the repository |
| Database | PostgreSQL | Foreign keys, constraints, indexes, transactional writes |
| Query access | pgx/v5 + sqlc | Write SQL; generate typed query methods |
| Migrations | Versioned SQL with one migration runner | Never use runtime automatic schema mutation |
| Contract | OpenAPI 3.0.3 | Schemas/examples/errors reviewed before handlers |
| Authentication | Opaque server-side sessions | Separate staff and customer session audiences |
| Passwords | Argon2id through a maintained Go implementation | Store hash parameters; never store plaintext |
| Media | S3-compatible object storage + CDN | Private pending uploads, validated public product images |
| Background work | PostgreSQL jobs + transactional outbox | Durable retries, unique job keys, inspection/replay tooling |
| Logging | Go slog, JSON | Request IDs, redacted personal data and secrets |
| Local environment | Docker Compose | API, worker, PostgreSQL, local object storage, email test inbox |
| CI | Repository's chosen CI service | Formatting, generation, tests, vet, build, dependency scan |

Start without Redis or a separate message broker. If the API later runs across replicas, rate limiting must have shared state or be enforced at a shared gateway. Add caching only after measuring useful targets.

### Suggested layout

```text
cmd/api/
cmd/worker/
internal/platform/{config,database,http,logging,storage,jobs}/
internal/identity/
internal/catalog/
internal/inventory/
internal/customers/
internal/cart/
internal/discounts/
internal/checkout/
internal/orders/
internal/payments/
internal/store/
internal/content/
internal/reporting/
db/migrations/
db/queries/
internal/dbgen/
api/openapi.yaml
tests/integration/
```

Use handler -> service -> database/query boundary. Gin types remain in the HTTP layer. Services accept context.Context and domain inputs, own business rules, and use a shared unit-of-work boundary for cross-module database transactions. Avoid generic repositories that obscure stock/order SQL.

## 3. Ownership and collaboration

| Developer | Primary ownership | Secondary/review responsibility |
|---|---|---|
| Abdulrhman | Repository foundation, shared HTTP contract, identity/sessions, staff permissions, store/shipping settings, reporting, deployment/operations | Review Ahmed's payment/order correctness and Huda's access control |
| Waheed | Products, variants/options, categories/collections, media integration, inventory ledger and reservations | Review Huda's promotion/cart product rules; stock portion of checkout |
| Huda | Customer profiles/addresses, carts, discounts, storefront content/navigation, transactional notification templates/jobs | Review catalog API usability; assist reporting reconciliation |
| Ahmed | Checkout orchestration, orders, online payments/webhooks/reconciliation, fulfillment, cancellations and refunds | Review shared transaction/job foundations and inventory concurrency |

Each owner delivers migrations, queries, service logic, handlers, OpenAPI changes, tests, fixtures, and module documentation. Testing is part of every developer's tasks. Abdulrhman coordinates integration; this is not an assumption about seniority.

Shared contracts to approve by the end of week 1:

- Identity: Principal{ID, Audience, Permissions}; customer and staff privileges cannot be interchanged.
- Catalog: GetPurchasableVariants(ids) returns publish state, SKU, options, currency and current price.
- Inventory: Reserve(tx, orderID, lines, expiresAt), Commit(tx, reservationID), Release(tx, reservationID).
- Discounts: Quote(ctx, cart, customer), ReserveUsage(tx, orderID), CommitUsage(tx, orderID), ReleaseUsage(tx, orderID).
- Orders: checkout owns orchestration; one caller-owned transaction passes through stock, discount and order writes.
- Payments: CreatePayment, VerifyWebhook, GetPaymentStatus, Refund with stable provider idempotency references.
- Jobs/outbox: Enqueue(tx, eventType, aggregateID, payload, uniqueKey).
- Reporting: agreed definitions of paid revenue, refunds, orders, returning customers, average order value and category attribution.

Do not let each module invent incompatible transaction or error interfaces. All schema changes require review from affected owners. Merged migrations are immutable; resolve migration filename collisions before merge.

## 4. Data model

| Module | Main entities | Important constraints |
|---|---|---|
| Identity | users, credentials, sessions, reset_tokens, verification_tokens, staff_memberships, staff_invites | Unique normalized login identity; hashed session/token values; expiring single-use recovery/invite tokens |
| Store | store_settings, shipping_zones, shipping_rates, tax_rules | One store; validated timezone/currency; no raw provider secrets in editable settings |
| Catalog | products, product_variants, product_options, option_values, variant_option_values, categories, product_categories, collections, collection_products, product_media | Unique SKU and slug; unique option combination; valid publish status; one default variant even for simple products |
| Inventory | inventory_balances, inventory_movements, stock_reservations, reservation_items | One location initially; integer quantities; nonnegative availability; unique movement reference |
| Customers | customers, customer_addresses, customer_notes | Optional unique user link; guest checkout does not create account access automatically |
| Carts | carts, cart_items | Hashed opaque guest token; unique variant per cart; positive item quantity; expiry/version |
| Promotions | discounts, discount_targets, discount_reservations, discount_redemptions | Unique normalized code; valid dates; bounded usage; reservation/commit rules |
| Checkout/orders | checkout_quotes, orders, order_items, order_addresses, order_events | Unique order number; immutable item/address/pricing snapshots; lifecycle constraints |
| Payments | payment_attempts, payment_events, refunds | Unique provider event ID and payment reference; refund maximum enforced atomically |
| Fulfillment | shipments, shipment_items | Manual carrier/tracking initially; full-order fulfillment in MVP |
| Content | content_blocks, navigation_items | Validated structure, ordered published content, sanitized rich text |
| Operations | idempotency_records, outbox_events, jobs, export_jobs, audit_logs | Request hashes; durable state/retries; restricted personal-data access |

Use UUID identifiers and UTC timestamptz timestamps. Convert reporting boundaries using the store's configured timezone. Store money as integer minor units with currency, never floating-point values. Define rounding and discount/tax allocation once and reuse it for quotes, orders and refunds.

Order items snapshot title, SKU, options, unit amount, quantity, allocated discounts and taxes. Address edits and product price changes must not rewrite historical orders. Archive sold products rather than delete their history.

Add indexes for unique SKU/slug/code/provider IDs, order date/status/customer queries, cart ownership, stock variant lookup, job availability and foreign-key joins. Add indexes based on actual queries and inspect representative query plans before launch.

## 5. API conventions

All business routes below are prefixed /api/v1. Routes starting /admin require a staff session and an explicit permission. Storefront account routes require customer authentication; guest cart/order routes require a scoped ownership token.

- JSON uses snake_case. Separate request/response DTOs from generated database models.
- Lists support limit (default 20, maximum 100), cursor pagination, and allowlisted filters/sort fields. Use a stable tie-breaker such as ID.
- Use 200 for reads/updates, 201 for creation, 202 for accepted jobs, and 204 for successful no-content deletion/logout.
- Use 400 for malformed requests, 401 for missing/invalid identity, 403 for insufficient permission, 404 for inaccessible/missing resources, 409 for state/idempotency/version conflicts, 422 for valid JSON violating business rules, and 429 for throttling.
- Error example: {"error":{"code":"INSUFFICIENT_STOCK","message":"An item is no longer available.","fields":{"items.0.quantity":"Only 1 available"},"request_id":"..."}}.
- List example: {"data":[],"pagination":{"next_cursor":null,"has_more":false}}.
- Require Idempotency-Key on checkout creation, payment creation, refunds and inventory adjustments. Scope records by operation and authenticated/guest principal; store a request hash. Same key/different payload returns 409. Define expiry/replay policy and enforce durable domain uniqueness beyond that expiry.
- Use version/If-Match checks on competing staff edits where stale overwrites are harmful.
- OpenAPI documents every parameter, enum, auth requirement, request, response and expected failure. Publish a generated client/mock contract for frontend work.

### Endpoint inventory

Combined verbs in one row describe separate operations; detailed schemas are the next implementation artifact, not already completed work.

| Owner | Method | Path | Purpose |
|---|---|---|---|
| Abdulrhman | GET | /auth/csrf | Obtain CSRF token bound to the appropriate session context |
| Abdulrhman | POST | /auth/register | Customer registration; never grants staff privileges |
| Abdulrhman | POST | /auth/login | Customer login |
| Abdulrhman | POST | /auth/logout | Revoke customer session |
| Abdulrhman | POST | /auth/forgot-password | Non-enumerating recovery request |
| Abdulrhman | POST | /auth/reset-password | Consume single-use reset token; revoke affected sessions |
| Abdulrhman | POST | /auth/verify-email | Consume verification token |
| Abdulrhman | POST | /auth/resend-verification | Throttled verification request |
| Abdulrhman | POST | /admin/auth/login | Staff login with staff-membership check |
| Abdulrhman | POST | /admin/auth/logout | Revoke staff session |
| Abdulrhman | GET | /admin/me | Current staff profile and permissions |
| Abdulrhman | GET | /admin/staff | List staff; owner permission |
| Abdulrhman | POST | /admin/staff/invitations | Invite staff with a bounded role |
| Abdulrhman | POST | /auth/staff-invitations/accept | Consume invitation; authenticated identity must match invitation |
| Abdulrhman | PATCH | /admin/staff/:id | Change role/status; preserve at least one active owner |
| Waheed | GET, POST | /admin/products | Filtered list; create draft product |
| Waheed | GET, PATCH | /admin/products/:id | Product detail; validated edit including publish/archive state |
| Waheed | POST | /admin/products/:id/variants | Add variant and option combination |
| Waheed | PATCH | /admin/products/:id/variants/:variant_id | Edit/archive variant; check parent ownership |
| Waheed | POST | /admin/media/uploads | Issue restricted upload target/object ID |
| Waheed | POST | /admin/media/:id/complete | Verify ownership, size/type and image decode before publication |
| Waheed | POST | /admin/products/:id/media | Attach verified media and ordering |
| Waheed | PATCH, DELETE | /admin/products/:id/media/:media_id | Reorder/set alt text; detach media |
| Waheed | GET, POST | /admin/categories | List/create categories |
| Waheed | PATCH, DELETE | /admin/categories/:id | Update/remove category with defined reassignment behavior |
| Waheed | GET, POST | /admin/collections | List/create manually curated collections |
| Waheed | GET, PATCH | /admin/collections/:id | Read/edit collection and product membership |
| Waheed | GET | /admin/inventory | Stock levels, available/reserved and low-stock filters |
| Waheed | GET | /admin/inventory/:variant_id/movements | Immutable adjustment/reservation/shipment history |
| Waheed | POST | /admin/inventory/adjustments | Atomic adjustment with reason and idempotency |
| Waheed | GET | /storefront/products | Published catalog with search/filter/sort |
| Waheed | GET | /storefront/products/:slug | Product, purchasable variants and public media |
| Waheed | GET | /storefront/categories | Published categories |
| Waheed | GET | /storefront/collections | Published collections |
| Waheed | GET | /storefront/collections/:slug/products | Published products in collection |
| Huda | GET, PATCH | /account | Customer profile; identity changes use verification rules |
| Huda | GET, POST | /account/addresses | List/create own addresses |
| Huda | PATCH, DELETE | /account/addresses/:id | Update/delete own address |
| Huda | GET | /admin/customers | Search/filter customer records |
| Huda | GET | /admin/customers/:id | Profile, aggregate spend and order references |
| Huda | POST | /admin/customers/:id/notes | Audited staff-only note |
| Huda | POST | /storefront/carts | Create guest/account cart and ownership context |
| Huda | GET | /storefront/carts/:id | Read owned cart; return current pricing/availability |
| Huda | POST | /storefront/carts/:id/items | Add variant quantity |
| Huda | PATCH, DELETE | /storefront/carts/:id/items/:item_id | Set positive quantity or remove owned item |
| Huda | PUT, DELETE | /storefront/carts/:id/discount | Set/clear one code; preview eligibility only |
| Huda | POST | /storefront/carts/:id/merge | Merge owned guest cart into signed-in cart after login |
| Huda | GET, POST | /admin/discounts | List/create percentage or fixed discounts |
| Huda | GET, PATCH | /admin/discounts/:id | Detail, rules, schedule, enable/disable |
| Huda | GET, PATCH | /admin/content/home | Manage approved homepage sections |
| Huda | GET, PATCH | /admin/content/navigation | Manage navigation links/order |
| Huda | GET | /storefront/home | Published homepage content and curated products |
| Huda | GET | /storefront/navigation | Published navigation |
| Ahmed | POST | /storefront/checkout/quotes | Revalidate cart and compute shipping/tax/discount; issue short-lived quote |
| Ahmed | POST | /storefront/checkouts | Consume valid quote; reserve stock/discount; create pending order |
| Ahmed | POST | /storefront/orders/:id/payments | Create/reuse provider attempt for an owned order |
| Ahmed | GET | /storefront/orders/:id | Token-scoped guest/customer order and payment state |
| Ahmed | GET | /account/orders | Authenticated customer's order history |
| Ahmed | GET | /account/orders/:id | Authenticated customer's order detail |
| Ahmed | POST | /webhooks/payments/:provider | Verify raw-body signature; durably deduplicate/process provider event |
| Ahmed | GET | /admin/orders | Search/filter/paginate orders |
| Ahmed | GET | /admin/orders/:id | Order snapshots, timeline, payment and shipment |
| Ahmed | POST | /admin/orders/:id/notes | Append staff note/event |
| Ahmed | POST | /admin/orders/:id/cancel | Validated cancellation; inventory/payment compensation |
| Ahmed | POST | /admin/orders/:id/shipments | Record full-order shipment; commit stock transition |
| Ahmed | PATCH | /admin/orders/:id/shipments/:shipment_id | Update tracking/delivery; validate parent and state |
| Ahmed | POST | /admin/orders/:id/payment-confirmations | Restricted audited COD collection confirmation; reject online-payment use |
| Ahmed | POST | /admin/orders/:id/refunds | Request idempotent full refund with reason |
| Abdulrhman | GET, PATCH | /admin/settings/store | Name, support details, timezone and approved settings |
| Abdulrhman | GET, POST | /admin/shipping/zones | Manage geographic eligibility |
| Abdulrhman | PATCH, DELETE | /admin/shipping/zones/:id | Edit/remove zone |
| Abdulrhman | GET, POST | /admin/shipping/rates | Manage flat/free-over-threshold rates |
| Abdulrhman | PATCH, DELETE | /admin/shipping/rates/:id | Edit/remove rate |
| Abdulrhman | GET | /storefront/config | Public name, currency, supported methods and public settings only |
| Abdulrhman | GET | /admin/analytics/overview | KPI totals and comparison period |
| Abdulrhman | GET | /admin/analytics/revenue | Daily/weekly revenue series |
| Abdulrhman | GET | /admin/analytics/categories | Category sales distribution |
| Abdulrhman | GET | /admin/analytics/bestsellers | Top products by explicit metric |
| Abdulrhman | POST | /admin/exports | Queue permitted order/product/report export |
| Abdulrhman | GET | /admin/exports/:id | Job state and short-lived private download link |
| Abdulrhman | GET | /admin/audit-logs | Restricted searchable staff action trail |

Operational routes outside /api/v1: GET /health/live, GET /health/ready. Restrict metrics/diagnostics to internal access. Tax configuration endpoints are added only after the chosen tax model is defined; do not expose an undefined tax-rule editor.

## 6. Business rules that must precede endpoint coding

### Identity and authorization

- Separate customer and staff session cookies, audiences and permissions. No public endpoint can create a staff membership.
- Secure, HttpOnly, host-scoped cookies with an explicit SameSite policy. Prefer same-origin API proxying for each frontend. Enforce CSRF tokens and trusted Origin checks for cookie-authenticated mutations; use narrow CORS allowlists if cross-origin calls are necessary.
- Session rotation after login/privilege change, idle/absolute expiry, logout revocation, throttled login/recovery, and password-reset invalidation.
- Fixed permissions such as catalog.write, inventory.adjust, orders.fulfill, orders.refund, customers.read, analytics.read and staff.manage. Approve a role-to-permission matrix, including which roles may see customer addresses.
- Ownership checks apply to each object and nested resource, not just route groups. An order ID, cart ID or matching email is not proof of ownership.
- Attach historical guest orders to an account only after verified identity/ownership. Never grant order access merely because someone registered the same email address.

### Inventory and order lifecycle

- Available = on_hand - reserved. Cart activity alone does not reserve stock.
- Checkout reserves stock for a bounded payment window. Use deterministic row-lock ordering or atomic conditional updates, within the order transaction, so two buyers cannot buy the last unit.
- Payment success makes an online order eligible for fulfillment; stock remains reserved until shipment. Shipment reduces both on_hand and reserved once. COD may become fulfillable before collection under the chosen policy.
- Confirmed online orders and accepted COD orders no longer use the temporary payment-reservation expiry. Their reservations remain held until fulfillment or an explicit cancellation; the expiry worker must exclude them.
- Expiry or valid cancellation releases unshipped reservations once. Inventory adjustments cannot reduce on_hand below reserved unless an explicit shortage workflow exists.
- Worker expiry and payment webhooks lock and validate the same order/reservation state. A late success after cancellation/expiry goes to reconciliation/refund or an explicit re-reservation attempt; never blindly revive an order or make stock negative.
- Order lifecycle: awaiting_payment -> confirmed -> processing -> shipped -> delivered; eligible unshipped orders may transition to cancelled. COD confirmation rules are explicit.
- Payment lifecycle is separate: pending, succeeded, failed, expired, refund_pending, refunded. Fulfillment lifecycle is separate: unfulfilled, shipped, delivered. A refund does not automatically mean goods have returned.

### Checkout and pricing

1. Load the owned cart and validate current product publish state, variants, quantity limits and currency.
2. Compute authoritative amounts using current prices, eligible discount, address-based shipping and the approved tax rules.
3. Create a short-lived versioned quote with amount breakdown; checkout revalidates it and returns 409/422 if repricing requires customer review.
4. Within one database transaction, claim idempotency, reserve stock and bounded discount usage, snapshot order/address/items, and record pending payment work/outbox.
5. Commit the database transaction before calling a remote payment provider. Create/retry the provider session using a stable reference and provider-supported idempotency behavior.
6. On provider timeout, treat the result as unknown and reconcile before creating another charge. Use a worker to recover interrupted calls and expired reservations.
7. Confirm online payment only from a verified provider event or authenticated server-side provider query. The browser success redirect is informational.
8. Consume promotions once and enqueue confirmation email transactionally. Repeated events/jobs must have no additional effect.

No database transaction stays open across external network calls. Never accept client-supplied totals or payment status as authoritative.

Use provider-hosted checkout or provider-tokenized payment fields; the VAYRO API must not receive or persist raw card details. Lock the selling currency after the first financial transaction unless a deliberate migration is implemented. Changing configured tax or shipping rules affects new quotes, not historical order snapshots.

### Promotions, refunds, reporting

- MVP has one code per cart, fixed or percentage, dates, minimum subtotal, optional product/collection scope and global usage limit. Specify rounding, eligible lines and fixed-discount allocation.
- Quote validation is not redemption. Reserve bounded usage during order creation; commit/release under the same explicit order rules. Per-customer limits require a verified customer identity; do not promise reliable per-person limits for anonymous guests.
- Full refunds have an atomic refundable-balance guard and pending state; provider retries cannot issue duplicate refunds. COD refunds require an explicitly chosen manual process and an audited reference.
- Restocking is a separate inventory adjustment after receipt/inspection; cancellation before shipment can release reservations automatically.
- Define the dashboard's displayed revenue precisely. Suggested metric: successful collected merchandise amounts after discounts, excluding shipping/tax, less successful merchandise refunds in their effective period. COD counts when collection is recorded. Label the UI consistently if a different metric is chosen.
- Define order count/AOV cohort, returning-customer identity, comparison-period boundaries and category attribution from historical line snapshots. Category allocation must not double-count products in multiple categories; use one reporting category per order line or an agreed allocation.
- For an empty comparison period return an explicit unavailable percentage rather than divide by zero. Keep a ledger-backed reconciliation fixture for all reports.

## 7. Work packages and acceptance criteria

| ID | Owner | Work package | Depends on | Acceptance evidence |
|---|---|---|---|---|
| A1 | Abdulrhman | Bootstrap, config, HTTP/errors, OpenAPI, CI, migration/job skeleton | Week-1 decisions | Fresh clone starts documented local stack; CI builds API/worker and applies migrations |
| A2 | Abdulrhman | Customer/staff identity, CSRF, sessions, roles, invites | A1 | Login/logout/reset work; invalid roles and cross-audience access rejected |
| A3 | Abdulrhman | Store/shipping settings and audit plumbing | A1 | Public config leaks no secrets; checkout can query deterministic shipping rates |
| A4 | Abdulrhman | Reporting, private exports, operations | O2/O3 event definitions | Known fixtures reconcile; restore/retry/export runbooks demonstrated |
| W1 | Waheed | Product/category/collection schema and API | A1 | Valid unique SKU/options; draft products hidden; variants archived safely |
| W2 | Waheed | Verified media upload and storefront catalog | W1, storage adapter | Invalid/foreign uploads rejected; filters and publication rules tested |
| W3 | Waheed | Inventory ledger/reserve/commit/release | W1, shared transaction boundary | Last-unit concurrency passes; replay cannot double-adjust; expiry safe |
| W4 | Waheed | Search/stock query tuning and catalog integration | W2/W3 | Realistic catalog queries measured; dashboard/storefront consume contracts |
| H1 | Huda | Profiles, addresses, staff customer views | A2 | Cross-customer access rejected; guest/account matching doesn't grant access |
| H2 | Huda | Carts, ownership, merge, server pricing inputs | W1, A2 contracts | Add/update/remove/merge deterministic; cart totals never trusted at checkout |
| H3 | Huda | Discount engine, reservations and API | H2, shared transaction boundary | Date/threshold/usage/concurrency cases pass; one code rule consistent |
| H4 | Huda | Home/navigation content and notification jobs | W2, A1 jobs, O2 events | Sanitized published content; duplicate events do not resend logical notification |
| O1 | Ahmed | Order schema/state machine and provider interface | A1, week-1 policies | State transitions specified and tested before external integration |
| O2 | Ahmed | Quotes and transactional checkout | A3, W3, H2/H3 | Order + stock + promo atomic; changed quotes fail clearly; retries reuse order |
| O3 | Ahmed | Payment attempts, webhook inbox and reconciliation | O1/O2, provider sandbox | Valid/invalid/replayed/out-of-order/late events tested; unknown timeouts recover |
| O4 | Ahmed | Fulfillment, COD, cancellation and full refunds | O3, W3 | Invalid transitions blocked; refund/stock effects idempotent; end-to-end lifecycle passes |

Parallelism uses agreed contracts, fixtures and provider fakes until real dependencies are merged. Fakes unblock coding, but acceptance requires integration against actual PostgreSQL and the real shared services.

## 8. Delivery sequence

Planning estimate: eight full-time working weeks for an MVP candidate, with up to two additional weeks for integration/provider/operational contingency. This is not a commitment; re-estimate after the first week's screen mapping, provider selection and team-capacity review. Dates are intentionally relative until a start date and working calendar are agreed.

| Week | Abdulrhman | Waheed | Huda | Ahmed | Shared exit gate |
|---|---|---|---|---|---|
| 1 | A1; identity/permission design | Catalog schema/field mapping | Customer/cart/promo/content contracts | Order/payment state models; provider spike | ERD, OpenAPI baseline, policies and module transaction contracts approved; CI/local stack runs |
| 2 | A2 sessions/staff/customer identity | W1 catalog CRUD/variants | H1 profiles; H2 cart skeleton with contract fixtures | O1 order schema; provider adapter/fake | Staff can create product; storefront can read it; authorization works |
| 3 | A3 shipping/store settings; audit | W2 media/public catalog; W3 stock | H2 real cart; H3 promotion quote rules | O2 quote and checkout against module contracts | Browser-to-API catalog/cart/quote journey works in staging |
| 4 | Integration support; job/outbox hardening | W3 concurrency, expiry interfaces | H3 discount reservation; H4 content | O2 real atomic checkout; payment session creation | Guest checkout creates exactly one reserved pending order under retries |
| 5 | Reporting queries and permission review | Stock/payment race tests; W4 | Notification jobs; customer dashboard integration | O3 signed webhooks/reconciliation; COD if selected | Sandbox purchase confirms once; expiry/failure releases correctly |
| 6 | A4 KPIs/private exports | Inventory dashboard; fulfill/cancel integration | Content/discount/customer UI integration | O4 shipment/cancellation/full refund | Complete purchase-to-delivery/refund journey; reporting fixtures reconcile |
| 7 | Deployment/backups/observability | Query/load tuning; media abuse tests | Ownership/security/notification retry tests | Provider outage/replay/refund recovery tests | All launch UI contracts integrated; critical failure cases pass |
| 8 | Restore rehearsal and release coordination | Catalog/inventory acceptance fixes | Storefront/customer acceptance fixes | Checkout/payment/order acceptance fixes | Production candidate and documented release/rollback checklist |

Show one integrated demo each week. Share OpenAPI/mock responses from week 1 so frontend implementation can proceed without waiting for all endpoints. The critical path is inventory + cart + shipping/discount quote -> transactional checkout -> verified payment -> fulfillment/reporting.

## 9. Quality and release gates

Every merge request must include contract/schema changes, meaningful tests, and a reviewer from another module. High-risk order/payment/stock changes require review by an affected owner.

### Automated checks

- gofmt, go vet, go test ./..., go test -race ./..., build API/worker, dependency vulnerability checks.
- sqlc generation stays in sync; OpenAPI linting and response-contract checks.
- Apply migrations to a clean database and upgrade an existing fixture database. Favor forward fixes and compatible expand/contract migrations over destructive rollback.
- Handler tests through net/http/httptest; service rule tests; PostgreSQL integration tests for constraints/locking/transactions.
- Provider sandbox integration tests and recorded fake scenarios for webhook retries, timeout/unknown outcome and reconciliation.

### Mandatory commerce/security scenarios

1. Two simultaneous checkouts for one available unit: only one reserves it.
2. Checkout retry with same key returns the original order; changed payload returns conflict.
3. Duplicate payment event produces one state change and one logical confirmation notification.
4. Payment success races with reservation expiry/cancellation: stock/order/payment reconcile safely.
5. Global discount last-use race: at most one valid reservation/redemption wins.
6. Client changes price, currency, totals, publish state or item IDs: server rejects/recomputes.
7. Customer A cannot read/update customer B's cart, address or order; staff role restrictions hold.
8. Refunding or shipping twice cannot repeat financial/inventory effects; refund after failed provider call remains recoverable.
9. Product/address edits do not change old orders; archiving prevents new purchases without breaking history.
10. Provider/email/storage outage has bounded timeouts, durable retry and visible operational state.
11. Metrics match a seeded dataset containing discounts, tax, shipping, COD, cancellation, refunds and comparison dates.
12. Malicious uploads and guessed upload IDs are rejected; rich content sanitized; secrets and personal data excluded from logs.

### Production readiness

- Separate development, staging and production data, storage and provider credentials.
- HTTPS, trusted-proxy configuration, restrictive CORS/CSRF/session policies, request-size limits and server/database timeouts.
- Database backups plus a successful restore rehearsal; business-agreed recovery objectives.
- Monitor error rate, latency, database pool pressure, job/outbox lag, payment reconciliation backlog, stuck orders and stock anomalies. Agree alert thresholds using staging load measurements.
- Run load tests using an agreed launch traffic/catalog/order forecast; set latency targets in week 1 rather than invent capacity guarantees.
- API/worker graceful shutdown, readiness checks, compatible migrations, release smoke test and application rollback procedure.
- Restricted job/webhook replay procedure with audit records; no new admin UI is required if a documented operator CLI is sufficient.
- Owner account provisioned securely; sandbox credentials replaced; actual provider webhook verification and refund path validated.

An endpoint is complete when its contract, migrations, permissions, service rules, tests, audit effects and staging integration are complete. A successful local request alone is not completion.

## 10. First assignments

Start with these concrete deliverables before splitting into isolated endpoint work:

1. Abdulrhman: bootstrap API/worker/local environment, common error/auth context, migration discipline, CI and OpenAPI skeleton.
2. Waheed: map product/editor/inventory fields to product/variant/media/stock schemas and write catalog endpoint contracts.
3. Huda: map storefront/customer/cart/discount fields, cart ownership/merge behavior and promotion examples; draft their contracts.
4. Ahmed: draft order/payment/fulfillment state tables, checkout amount examples and provider evaluation/sandbox checklist.
5. All four: review ERD, money/rounding, permission matrix, shared transaction signatures and one end-to-end API example.

First integrated milestone: staff creates and publishes a variant product; a guest sees it, puts it in an owned cart, and receives a server-calculated checkout quote. Follow that with reservation and sandbox payment, rather than implementing every dashboard screen before checkout integration.

## Technical references

- [Gin middleware and route groups](https://gin-gonic.com/en/docs/middleware/using-middleware/): group authentication/authorization and shared middleware.
- [Gin binding and validation](https://gin-gonic.com/en/docs/binding/binding-and-validation/): use ShouldBind-style methods for controlled JSON errors, plus service-level validation.
- [Go relational database guide](https://go.dev/doc/database/): transaction and connection-pool fundamentals.
- [Go cancellation](https://go.dev/doc/database/cancel-operations): propagate request context and bounded timeouts to queries/providers.
- [PostgreSQL SELECT locking](https://www.postgresql.org/docs/current/sql-select.html): row locks and concurrent stock/order transitions.
- [sqlc documentation](https://docs.sqlc.dev/): typed Go queries generated from SQL.
- [OWASP session management](https://cheatsheetseries.owasp.org/cheatsheets/Session_Management_Cheat_Sheet): session handling and cookie controls.
- [OWASP CSRF prevention](https://cheatsheetseries.owasp.org/cheatsheets/Cross-Site_Request_Forgery_Prevention_Cheat_Sheet.html): CSRF tokens and SameSite defense in depth.
