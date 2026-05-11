# Impact Analysis

## Affected Modules

- `inventory`: included in aggregate module checks because it is registered in `.arch/modules.yaml`.

## Module Ownership Impact

No module ownership changes.

## Public Contract Impact

Adds one CLI command:

```text
vibit check all
```

Existing commands remain:

- `vibit check architecture`
- `vibit check change <change-id>`
- `vibit check module <module>`
- `vibit generate module <module>`

## Data And Migration Impact

No runtime data or migrations.

## Test Impact

Verification should include:

- CLI help output
- Aggregate check
- Individual architecture check
- Individual change check
- Individual module check

## Documentation Impact

Update:

- README
- AGENTS
- CLI change docs
- Conversation log

## Compatibility Risks

Low. The command is additive. Main risk is brittle discovery logic if manifests evolve.
