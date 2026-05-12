# Impact Analysis

## Affected Modules

No domain modules are affected.

## Module Ownership Impact

No ownership changes.

## Public Contract Impact

Adds a new CLI inspection command:

```bash
node tools/vibit inspect change <change-id>
```

Adds a new `change_inspection` JSON output kind to `schema/inspect-output.schema.json`.

## Event Impact

No events are added, changed, or removed.

## Permission Impact

No permissions are added, changed, or removed.

## Data And Migration Impact

No runtime data or migrations.

## Test Impact

Verification must confirm valid JSON output for an existing change and a structured `exists: false` JSON result for a missing change.

## Documentation Impact

Update:

- README
- AGENTS
- Inspect output schema
- Change spec
- Conversation log

## Compatibility Risks

Low. This is an additive inspect command.
