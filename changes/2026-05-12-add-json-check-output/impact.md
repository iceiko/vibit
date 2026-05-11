# Impact Analysis

## Affected Modules

- `inventory`: used by `check module inventory --json`.

## Module Ownership Impact

No ownership changes.

## Public Contract Impact

Adds JSON output variants for check commands.

Default text behavior should remain unchanged.

## Data And Migration Impact

No runtime data or migrations.

## Test Impact

Verification must confirm both text and JSON output paths.

## Documentation Impact

Update:

- README
- AGENTS
- Schema validation standard
- Check output schema
- Conversation log

## Compatibility Risks

Medium. This introduces a machine-readable output shape that future agents may rely on. Mark schema version `0.1` and avoid over-promising final stability.
