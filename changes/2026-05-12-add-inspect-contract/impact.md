# Impact Analysis

## Affected Modules

The inventory module is used as the first verified contract inspection target.

## Module Ownership Impact

No ownership changes.

## Public Contract Impact

Adds a CLI inspection command:

```bash
node tools/vibit inspect contract --module <module> --type <type> --id <id>
```

Supported type values:

- `command` or `commands`
- `query` or `queries`
- `event` or `events`
- `error` or `errors`
- `permission` or `permissions`

Adds a new `contract_inspection` JSON output kind to `schema/inspect-output.schema.json`.

## Data And Migration Impact

No runtime data or migrations.

## Test Impact

Verification must confirm valid JSON output for at least one command contract and one catalog contract, plus clear failure for a missing contract.

## Documentation Impact

Update:

- README
- AGENTS
- Inspect output schema
- Change spec
- Conversation log

## Compatibility Risks

Low. This is an additive inspect command.

The implementation remains a lightweight source reader and intentionally does not claim full YAML parsing or payload schema validation.
