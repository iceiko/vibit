# Request

## Original Request

The maintainer asked to continue after the first inventory contract source files were added.

## Clarified Requirement

Turn the manual contract registry and source-file consistency check into a first-class CLI command:

```bash
node tools/vibit check contracts
node tools/vibit check contracts --json
```

Include the new check in `node tools/vibit check all`.

## User-Visible Outcome

Agents can verify contract registry consistency without manually reading `.arch/contracts.yaml` and checking every referenced file.

## Non-Goals

- Do not add an external YAML parsing dependency.
- Do not implement full JSON Schema validation yet.
- Do not generate TypeScript types or validators yet.
- Do not change runtime contracts except where needed to satisfy the new check.

## Unknowns

- Exact future contract schema validator.
- Exact generated TypeScript output paths.
- Whether a dedicated YAML parser should be adopted later.

## Acceptance Criteria

- `node tools/vibit check contracts` validates `.arch/contracts.yaml` and referenced contract files.
- `node tools/vibit check contracts --json` returns machine-readable check results.
- `node tools/vibit check all --json` includes the contracts check.
- Rule metadata exists for contract check failures.
- README and AGENTS command lists mention the new command in English and Chinese.
- Verification is run and recorded.
