# Request

## Original Request

Continue ten work items unless maintainer confirmation is required.

## Clarified Requirement

Advance `W-0015` by adding the first PostgreSQL adapter that implements the inventory repository interface behind the platform persistence boundary.

## User-Visible Outcome

Maintainers and agents can inspect a concrete PostgreSQL repository adapter under:

```text
runtime/internal/platform/persistence/postgres/
```

The adapter implements:

- `inventory.Repository.GetInventory`
- `inventory.Repository.LockInventoryForMutation`
- locked `inventory.MutationLock.GetInventory`
- locked `inventory.MutationLock.GrantItem`

`GrantItem` records the quantity change and the durable grant row using the same executor supplied by the application-owned unit of work.

## Non-Goals

- Do not wire PostgreSQL into process startup yet.
- Do not add PostgreSQL connection pooling or configuration yet.
- Do not add migration apply or rollback tooling yet.
- Do not make live PostgreSQL integration tests mandatory.
- Do not add an outbox or event publication standard.
- Do not replace the inventory account row-lock concurrency model.
- Do not introduce object storage or MinIO.

## Unknowns

- The concrete runtime composition that binds a pgx transaction to inventory handlers remains deferred.
- Live PostgreSQL integration testing remains deferred until the project defines a disposable PostgreSQL test environment standard.
- Durable event publication outside the transaction remains deferred until an outbox or event delivery standard exists.

## Acceptance Criteria

- [ ] Add a PostgreSQL inventory repository adapter under the allowed platform persistence package.
- [ ] Keep `pgx/v5` imports inside `runtime/internal/platform/persistence/postgres/`.
- [ ] Bind repository behavior to an injected executor instead of opening hidden write transactions.
- [ ] Implement inventory account row locking through `SELECT ... FOR UPDATE`.
- [ ] Upsert inventory item quantities through the locked mutation path.
- [ ] Record `inventory_item_grants` rows in the same executor path as the quantity mutation.
- [ ] Add focused tests for SQL shape, row mapping, lock behavior, grant recording, and transaction-control absence.
- [ ] Record skipped live PostgreSQL integration verification explicitly.
