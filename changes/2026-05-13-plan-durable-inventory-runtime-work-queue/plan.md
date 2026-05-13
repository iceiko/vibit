# Plan

1. Inspect the current work queue and workflow standard.
2. Add remaining M-002 work items after W-0015.
3. Mark one conservative next item as `next_ready`.
4. Keep ask-first boundaries on steps that could change architecture.
5. Run work inspection and repository checks.
6. Record verification.

## Resulting Queue Shape

- `W-0016`: Plan durable inventory runtime work queue.
- `W-0017`: Add PostgreSQL configuration and transaction runner.
- `W-0018`: Add migration apply and status tooling.
- `W-0019`: Define disposable PostgreSQL verification environment.
- `W-0020`: Wire persistent inventory runtime composition.
- `W-0021`: Verify durable inventory runtime end to end.

This order keeps foundational runtime plumbing and migration execution ahead of live wiring, and keeps disposable integration environment policy explicit before making live PostgreSQL verification mandatory.
