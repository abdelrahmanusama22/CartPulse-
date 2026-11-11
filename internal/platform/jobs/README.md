# Durable jobs and outbox - foundation owner: Abdulrhman

The worker entry point starts and shuts down, but there is no queue consumer yet. Design the job/outbox schema and transaction-aware enqueue interface before implementing handlers.

Required behavior: durable enqueue in the caller transaction, unique keys, claiming/leases, bounded retry/backoff, dead-letter state, cancellation and operator replay. Do not execute job side effects inside request goroutines as a substitute for durability.

Consumers: Huda owns notification templates/delivery jobs; Ahmed owns payment reconciliation; Waheed + Ahmed define stock expiry behavior; Abdulrhman owns export jobs. Critical side effects must be idempotent even if a job is delivered multiple times.
