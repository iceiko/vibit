# Impact Analysis

## Affected Modules

`inventory` is affected because this change generates the first contract shape for `GrantItem`.

## Module Ownership Impact

No ownership changes.

## Public Contract Impact

Adds CLI tooling commands:

```bash
node tools/vibit generate contract --module <module> --type <type> --id <id>
node tools/vibit check generated
node tools/vibit check generated --json
```

No runtime public command, query, event, error, or permission semantics change. The existing `GrantItem` contract remains the source of truth.

## Generated File Impact

Adds:

```text
modules/inventory/generated/contracts/GrantItem.generated.ts
```

The file is generated output and must not be hand-edited by ordinary agents. If it is wrong, update the source contract, generator, or template.

## Data And Migration Impact

No runtime data or migrations.

## Test Impact

Adds generated-file verification that checks:

- Declared generated files exist.
- Generated files include a generated marker.
- Generated files include source trace metadata.
- Generated files include generator trace metadata.

## Documentation Impact

Update:

- README
- AGENTS
- Module manifest generated files declaration
- Change spec
- Conversation log

## Compatibility Risks

Low. This is the first generated output and is additive.

The generated TypeScript file is not compiled yet because the repository has not introduced a TypeScript package, package manager, or test runner.
