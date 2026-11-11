# Team ownership and first tasks

| Developer | Owned modules | First deliverable |
|---|---|---|
| Abdulrhman | platform, identity, store, reporting, deployment | Identity/session/permission contract; migrations/tool versions; shared transaction/job interfaces |
| Waheed | catalog, inventory, product media/storage integration | Product/variant/media ERD, catalog contract and stock reservation invariants |
| Huda | customers, cart, discounts, content, notification templates/jobs | Cart ownership/merge contract, customer boundaries, promotion examples |
| Ahmed | checkout, orders, payments, fulfillment | Order/payment/fulfillment state tables, provider interface and checkout contract |

First shared milestone: staff publishes a variant product and adds stock; guest finds it, creates an owned cart and receives an authoritative quote. Next: atomic order reservation and sandbox payment.

Each owner ships migrations, queries, services, handlers, OpenAPI, tests, fixtures and module documentation. Review across boundaries: Ahmed/Waheed for stock and checkout; Ahmed/Huda for discounts and pricing; Abdulrhman/owners for identity and shared foundations.

GitHub handles have not been provided. Add a real CODEOWNERS file after mapping people to their GitHub usernames; names in this document do not automatically assign GitHub reviews.

Before checkout coding, agree currency/region, provider/COD, tax treatment, shipping rates and cancellation/refund policies. See `backend-plan.md` for the full work packages and schedule.
