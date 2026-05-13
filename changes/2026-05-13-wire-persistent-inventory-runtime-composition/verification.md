# Verification

Verified:

- `cd runtime && go test ./...`
- `cd runtime && go vet ./...`
- `node tools/vibit check architecture --json`
- `node tools/vibit check schemas --json`
- `node tools/vibit check module inventory --json`
- `node tools/vibit check runtime --json`
- `node tools/vibit check change wire-persistent-inventory-runtime-composition --json`
- `node tools/vibit check all --json`

Not verified:

- Live PostgreSQL end-to-end request-loop verification. This remains the next work item, `W-0021`, and requires the disposable PostgreSQL verification environment or an explicit unavailability record.
- Migration apply/status execution against a live PostgreSQL database. This remains part of `W-0021`.

Not applicable:

- New migrations; this change uses the existing inventory migration source.
- Contract changes; inventory command, query, event, permission, error, and Protobuf contracts did not change.
- Automatic startup migrations; normal server startup intentionally does not apply migrations.
