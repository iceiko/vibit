# Impact Analysis

## Affected Modules

No domain modules are affected.

## Module Ownership Impact

No ownership changes.

## Public Contract Impact

Adds a new CLI inspection command:

```bash
node tools/vibit inspect memory
```

Adds a new `memory_inspection` JSON output kind to `schema/inspect-output.schema.json`.

## Event Impact

No events are added, changed, or removed.

## Permission Impact

No permissions are added, changed, or removed.

## Data And Migration Impact

No runtime data or migrations.

## Test Impact

Verification must confirm valid JSON output and basic counts for changes, conversations, and decisions.

## Documentation Impact

Update:

- README
- AGENTS
- Inspect output schema
- Change spec
- Conversation log

## Compatibility Risks

Low. This is an additive inspect command.
