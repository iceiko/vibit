# Impact Analysis

## Affected Modules

- `inventory`: used as the first module for inspect command verification.

## Module Ownership Impact

No ownership changes.

## Public Contract Impact

Adds CLI commands:

```text
vibit inspect boundary --from <module> --to <module>
vibit inspect module <module>
```

Adds permission concept:

```text
generated_file_override
```

This is a standards-level permission concept, not yet a runtime permission system.

## Data And Migration Impact

No runtime data or migrations.

## Test Impact

Verification should include:

- `node tools/vibit check all`
- `node tools/vibit inspect module inventory`
- `node tools/vibit inspect boundary --from inventory --to player`
- Secret scan

## Documentation Impact

Adds:

- ADR-Agent decision record standard
- Decision directory
- Decision template
- First decision record
- Conversation log for external feedback

Updates:

- Constitution
- README
- AGENTS
- Architecture conventions

## Compatibility Risks

Low. This is additive. The main risk is over-recording decision rationale or storing private chain-of-thought. The standard should require concise public rationale, not hidden reasoning dumps.
