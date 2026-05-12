# Impact Analysis

## Affected Modules

No domain modules are affected.

## Module Ownership Impact

No ownership changes.

## Public Contract Impact

Adds a new CLI inspection command:

```bash
node tools/vibit inspect rule <rule-id>
```

Adds a new `rule_inspection` JSON output kind to `schema/inspect-output.schema.json`.

## Data And Migration Impact

No runtime data or migrations.

## Test Impact

Verification must confirm valid JSON output for an existing rule and clear failure for a missing rule.

## Documentation Impact

Update:

- README
- AGENTS
- Inspect output schema
- Change spec
- Conversation log

## Compatibility Risks

Low. This is an additive inspect command.
