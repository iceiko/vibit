# Conversation: Inventory Contracts

Date: 2026-05-12  
Participants: Maintainer, Agent  
Related changes:

- `changes/2026-05-12-add-inventory-contracts/`

Related artifacts:

- `.arch/contracts.yaml`
- `contracts/inventory/commands/GrantItem.yaml`
- `contracts/inventory/queries/GetInventory.yaml`
- `contracts/inventory/events/ItemGranted.yaml`
- `contracts/inventory/errors/inventory_errors.yaml`
- `contracts/inventory/permissions/inventory_permissions.yaml`
- `decisions/ADR-0007-yaml-contract-source-format.md`

## Context

The inventory module had been prepared as the first runtime proof slice. The next step was to declare source contracts before writing runtime implementation or generators.

## Maintainer Narrative

The maintainer asked:

```text
继续进行下一步
```

## Agent Response Summary

The agent selected YAML as the first contract source format and recorded the choice in ADR-0007.

The agent added a contract registry and first inventory contract files for the prepared proof slice:

- `GrantItem`
- `GetInventory`
- `ItemGranted`
- Inventory errors
- Inventory permissions

No runtime implementation, generator, package dependency, or persistence layer was added.

## Decisions

- Use YAML as the first contract source format.
- Keep contract files as source artifacts, not generated output.
- Embed JSON-Schema-like payload shapes inside YAML contracts for future generation and validation.

## Artifacts

- Added `.arch/contracts.yaml`.
- Added `contracts/inventory/`.
- Added `decisions/ADR-0007-yaml-contract-source-format.md`.
- Added `changes/2026-05-12-add-inventory-contracts/`.

## Open Questions

- Exact generated TypeScript output paths.
- Exact runtime validation library.
- Exact package manager and test runner.
- Initial persistence adapter.

## Follow-Up

- Add a contract check once the source format stabilizes.
- Generate TypeScript types, validators, dispatch shapes, and fixtures from the inventory contracts when runtime scaffolding starts.

## Redaction Notes

No secrets, tokens, credentials, private account details, or unrelated private data are recorded in this log.
