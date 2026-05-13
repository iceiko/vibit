# ADR-0020: PostgreSQL Persistence Boundary

Status: Accepted  
Date: 2026-05-13  
Decision Makers: Agent  
Related changes:

- `changes/2026-05-13-define-postgresql-repository-and-migration-boundary/`

Related conversations:

- `conversations/2026-05-12-persistence-direction.md`

Related artifacts:

- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/work-items.yaml`
- `docs/postgresql-persistence-boundary.md`
- `modules/inventory/module.yaml`
- `decisions/ADR-0011-postgresql-and-object-storage-persistence.md`
- `decisions/ADR-0013-first-go-runtime-dependencies.md`
- `decisions/ADR-0014-go-runtime-layout-and-boundaries.md`

## Context

The first WebSocket Protobuf request loop now works without persistence. The next milestone is durable inventory runtime behavior backed by PostgreSQL.

PostgreSQL, `pgx/v5`, and `goose/v3` are already accepted by earlier decisions. What remains undefined is the implementation boundary: where SQL lives, where transaction ownership lives, how inventory state is locked, how migrations are structured, and how agents should verify persistence work.

This needs a durable decision before code is added because persistence mistakes are difficult to infer from local code alone. A repository that "works" locally can still violate transaction ownership, hide permission checks, race capacity invariants, or make later event delivery unreliable.

## Decision

Adopt `docs/postgresql-persistence-boundary.md` as the standard for the first PostgreSQL persistence implementation.

The first durable module state is inventory. PostgreSQL persistence must be added in this order:

1. Define the persistence boundary and work sequence.
2. Add command-safe inventory repository mutation locking.
3. Add a transaction boundary skeleton.
4. Add SQL-first inventory migrations.
5. Add migration verification.
6. Add a PostgreSQL inventory repository adapter.
7. Wire persistent runtime configuration only after the adapter and verification path exist.

The inventory persistent grant flow must use a transaction-bound repository and an explicit inventory account row lock for `player_id` before reading current inventory and applying a grant.

`pgx/v5` remains owned by `runtime/internal/platform/persistence/postgres/`. `goose/v3` remains owned by `runtime/internal/platform/migrations/`. SQL migration source files live under `runtime/migrations/postgres/`.

MinIO and S3-compatible object storage remain out of scope for durable inventory state.

## Alternatives Considered

- Let the PostgreSQL repository open its own transactions internally.
- Use database constraints as the only capacity protection.
- Use a PostgreSQL advisory lock instead of an inventory account row lock.
- Add the first SQL migration and repository adapter without a persistence boundary standard.
- Make MinIO part of the durable inventory milestone.
- Defer all transaction boundary work until after a repository adapter exists.

## Rationale

Application-owned unit-of-work boundaries are easier for agents to inspect than repositories that silently start transactions. They also leave room for durable event recording and future outbox behavior without rewriting every command handler.

An explicit inventory account row lock is more readable than a PostgreSQL advisory lock and maps directly to the inventory aggregate being protected. It gives agents a concrete sequence to follow when preserving capacity invariants under concurrent grants.

SQL-first migrations keep schema changes inspectable and diffable. Goose markers are simple enough for agents to validate before deeper migration tooling exists.

Keeping MinIO out of this milestone preserves the accepted split: PostgreSQL is for authoritative transactional state, while S3-compatible object storage is for large artifacts when a concrete use case exists.

## Agent Reasoning Summary

The durable inventory path should move from boundary to lock contract to transaction skeleton to migration to adapter. This sequence reduces the chance that a future agent writes a locally passing PostgreSQL repository that violates capacity, event, or ownership rules.

## Decision Weights

```yaml
decision_weights:
  agent_context: high
  human_ergonomics: high
  implementation_cost: medium
  reversibility: medium
  long_term_maintainability: high
confidence: high
```

## Consequences

- W-0010 defines the boundary only; it does not add PostgreSQL implementation code.
- The next implementation work should add command-safe inventory mutation locking before a PostgreSQL adapter.
- SQL migrations must use deterministic sequence names and `goose` Up/Down markers.
- Persistent command flows must not rely on repositories opening hidden transactions.
- PostgreSQL integration verification may be gated until a disposable database test environment exists, but skipped verification must be recorded.

## Reversal Conditions

Revisit this decision if the explicit inventory account row lock creates unacceptable contention, if a stronger event log or outbox standard must be implemented before inventory persistence, or if early PostgreSQL adapter work shows the transaction boundary needs a different package owner.

## Follow-Up

- Add command-safe inventory repository mutation locking.
- Add transaction boundary interfaces under `runtime/internal/platform/tx/`.
- Add the first SQL inventory migration under `runtime/migrations/postgres/`.
- Add migration verification to `tools/vibit`.
- Add a PostgreSQL inventory repository adapter behind the inventory repository interface.
