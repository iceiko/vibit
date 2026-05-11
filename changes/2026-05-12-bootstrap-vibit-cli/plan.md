# Plan

## Implementation Language Decision

Use Node.js with only standard-library APIs for the first CLI prototype.

Reasoning:

- Node.js is available in the current Termux environment.
- Python is not currently installed.
- The first CLI only needs filesystem checks, deterministic reporting, and template generation.
- Avoiding npm dependencies keeps the first executable standard small and easy for agents to inspect.
- A future packaging layer can expose the same CLI through npm without changing the command contract.

## Files To Create

- `tools/vibit`

## Files To Edit

- `README.md`
- `README.zh-CN.md`
- `AGENTS.md`
- `AGENTS.zh-CN.md`
- `.arch/conventions.yaml`
- Possibly `.arch/generation.yaml`

## Generated Artifacts

Potentially generated module skeletons under `modules/<module>/`.

## Handwritten Logic

- CLI argument parsing
- Repository root discovery
- Required file checks
- Lightweight manifest text checks without external YAML dependencies
- Deterministic reporting
- Module skeleton generation

## Tests

- CLI help command
- Architecture check command
- Change check command
- Missing module check command

## Verification Commands

- `node tools/vibit --help`
- `node tools/vibit check architecture`
- `node tools/vibit check change bootstrap-vibit-cli`
- `node tools/vibit check module inventory`

## Rollback Or Migration Notes

Since this is the first CLI prototype, rollback is simply removing the CLI files and documentation references.
