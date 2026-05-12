# Request

## Original Request

The maintainer asked:

> 继续

## Clarified Requirement

Implement the smallest query-side runtime slice for the existing `GetInventory` inventory query.

The implementation must follow the same proof chain as `GrantItem`:

```text
contract source -> generated shape -> handwritten logic -> tests -> verification
```

Because the current generator only supports command contracts, this change also extends contract generation to query contracts before adding handwritten query logic.

## User-Visible Outcome

The inventory module demonstrates both write-side and read-side behavior:

- `GrantItem` can record inventory state.
- `GetInventory` can read inventory state without mutating it.
- Agents can verify both runtime paths with `node tools/vibit check runtime`.

## Non-Goals

- Do not add HTTP routes.
- Do not add persistent storage.
- Do not choose a package manager.
- Do not add external dependencies.
- Do not implement item catalog, player account, currency, reward, quest, or match behavior.
- Do not hand-edit generated files.
- Do not broaden query generation beyond the current YAML contract shape needed by `GetInventory`.

## Unknowns

- Whether query and command generated shapes should later share a richer reusable generator model.
- Whether read permission checks should remain in the current permission policy or become generated from the permission catalog later.
- Whether a future repository interface should be generated from module manifests.

## Acceptance Criteria

- [x] `GetInventory` contract shape is generated from `contracts/inventory/queries/GetInventory.yaml`.
- [x] The generated query shape is declared in `modules/inventory/module.yaml`.
- [x] `GetInventory` handler is implemented under the inventory query extension point.
- [x] The handler imports and uses the generated `GetInventory` contract shape.
- [x] Permission policy supports `inventory_read`.
- [x] Runtime tests cover successful inventory reads, empty reads, non-mutation, sorted item output, identity stability, and permission denial.
- [x] `node tools/vibit check runtime` runs both command and query tests.
- [x] `node tools/vibit check all --json` includes the runtime check and passes.
- [x] English and Simplified Chinese module docs mention the generated `GetInventory` shape and runtime path.
- [x] Verification is recorded.
