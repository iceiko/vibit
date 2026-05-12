# Request

## Original Request

The maintainer asked:

> 继续

## Clarified Requirement

Implement the smallest handwritten runtime slice for the existing `GrantItem` inventory command.

The implementation must consume the generated contract shape and stay inside declared inventory extension points:

- Command handler
- Inventory repository
- Inventory capacity policy
- Inventory permission policy
- Focused runtime tests

The implementation should remain dependency-free and use Node.js built-in tooling for the first executable test path.

## User-Visible Outcome

The repository demonstrates the next step in the proof chain:

```text
contract source -> generated shape -> handwritten logic -> tests -> verification
```

Agents can run one command to verify the handwritten `GrantItem` behavior.

## Non-Goals

- Do not implement HTTP routes.
- Do not implement `GetInventory`.
- Do not implement persistent storage.
- Do not choose a package manager.
- Do not add external dependencies.
- Do not hand-edit generated files.
- Do not broaden inventory into item catalog, player account, currency, reward, quest, or match behavior.

## Unknowns

- Whether Node's built-in TypeScript type stripping should remain the long-term test path.
- Whether the first runtime package should later add `tsc` and a package manager.
- Whether repository interfaces should become generated or remain handwritten extension points.

## Acceptance Criteria

- [x] `GrantItem` handler is implemented under the inventory command extension point.
- [x] The handler imports and uses the generated `GrantItem` contract shape.
- [x] Inventory repository and policies are implemented inside the inventory module.
- [x] Runtime tests cover successful grants, invalid quantity, capacity exceeded, and permission denied.
- [x] `node tools/vibit check runtime` runs the runtime tests.
- [x] `node tools/vibit check all --json` includes `check runtime`.
- [x] Generated files remain unchanged except through the generator.
- [x] English and Simplified Chinese docs mention the runtime check.
- [x] Verification is recorded.
