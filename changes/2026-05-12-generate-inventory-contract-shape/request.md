# Request

## Original Request

The maintainer asked:

> 继续

## Clarified Requirement

Start the first small generation step after contract inspection by generating a TypeScript contract shape for the existing inventory `GrantItem` command.

The scope is intentionally narrow:

- Generate one traceable TypeScript contract file from the existing YAML contract source.
- Declare the generated file in `modules/inventory/module.yaml`.
- Add verification that declared generated files exist and include trace metadata.
- Do not implement runtime business behavior yet.

Target command:

```bash
node tools/vibit generate contract --module inventory --type command --id GrantItem
```

## User-Visible Outcome

The repository demonstrates the first concrete step in the required chain:

```text
contract source -> generated shape -> handwritten extension point
```

Agents can regenerate the first contract shape from the source contract instead of hand-editing generated output.

## Non-Goals

- Do not implement the `GrantItem` command handler.
- Do not add HTTP routing.
- Do not choose a package manager.
- Do not add TypeScript compiler or runtime dependencies.
- Do not generate all inventory contracts in this change.
- Do not add full YAML parsing.

## Unknowns

- Whether future generated shapes should use interfaces, type aliases, or runtime validators as the stable public output.
- Whether a formal YAML parser should be introduced before generating every contract type.
- Whether generated files should eventually live under `generated/` or remain module-local.

## Acceptance Criteria

- [x] `node tools/vibit generate contract --module inventory --type command --id GrantItem` creates or refreshes the generated contract file.
- [x] The generated file is declared in `modules/inventory/module.yaml`.
- [x] The generated file includes source and generator trace metadata.
- [x] The generated file is not handwritten business logic.
- [x] `node tools/vibit check generated` verifies declared generated files.
- [x] `node tools/vibit check all --json` includes `check generated`.
- [x] English and Simplified Chinese docs mention the new command.
- [x] Verification is recorded.
