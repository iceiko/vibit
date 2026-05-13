# Impact Analysis

## Affected Modules

- `runtime`: gains a migration source verification command and runtime guide updates.
- `inventory`: its first migration source becomes machine-checked through table and ownership references.

## Module Ownership Impact

No data ownership changes. Inventory remains the owner of `inventory_accounts`, `inventory_items`, and `inventory_item_grants`.

## Public Contract Impact

No command, query, event, error, permission, Protobuf, or WebSocket protocol changes.

The `tools/vibit` CLI gains a new public command:

```bash
node tools/vibit check migrations
node tools/vibit check migrations --json
```

## Data And Migration Impact

No new migration files are added. The existing first migration is now validated for:

- root path
- deterministic SQL naming
- `-- +goose Up`
- `-- +goose Down`
- owning module trace
- required inventory table references
- absence of unapproved Go migrations

## Test Impact

The check itself is exercised through CLI verification commands. Go runtime behavior is unchanged.

## Documentation Impact

Update AGENTS guidance, runtime guidance, schema validation docs, PostgreSQL persistence docs, architecture manifests, and rule catalog metadata.

## Compatibility Risks

Low. This adds verification coverage and can fail future invalid migration sources earlier, but it does not change runtime behavior.
