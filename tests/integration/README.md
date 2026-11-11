# PostgreSQL integration tests

Place cross-module database tests here when actual schemas/services exist. Unit/handler tests stay beside their packages.

Use an isolated test database, apply migrations, and clean fixtures safely. Never point destructive fixtures at a shared development or production database. Add a PostgreSQL CI service and explicit test DSN when the first integration test is introduced.

Priority scenarios: last-unit concurrency, rollback of order/stock/promotion writes, repeated checkout, duplicate/late payment events, expiry races, cancellation, refund guards and cross-customer access. The scaffold's current Go tests do not require PostgreSQL.
