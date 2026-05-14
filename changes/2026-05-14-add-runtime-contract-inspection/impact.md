# Impact

## Affected Areas

- `tools/vibit`
- `.arch/work-items.yaml`

## Tooling Impact

`inspect contract` now handles runtime-owned contract families registered in `.arch/contracts.yaml`, starting with `runtime_contracts.session_validation`.

## Contract Impact

No contract source changes are introduced by this tooling step.

## Runtime Impact

No Go runtime behavior changes.

## Compatibility Risks

Low. The output keeps the existing `contract_inspection` shape and adds runtime-specific fields under existing objects.
