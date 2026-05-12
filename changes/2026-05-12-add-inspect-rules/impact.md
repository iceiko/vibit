# Impact Analysis

## Affected Modules

No domain modules are affected.

## Module Ownership Impact

No ownership changes.

## Public Contract Impact

Adds two CLI inspection forms:

```bash
node tools/vibit inspect rules
node tools/vibit inspect rules --category <category>
```

Adds a new `rules_inspection` JSON output kind to `schema/inspect-output.schema.json`.

## Data And Migration Impact

No runtime data or migrations.

## Test Impact

Verification must confirm valid JSON output and category filtering.

## Documentation Impact

Update:

- README
- AGENTS
- Inspect output schema
- Change spec
- Conversation log

## Compatibility Risks

Low. This is an additive inspect command.
