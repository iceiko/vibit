# Plan

## Files To Create

Pending implementation language decision.

Likely candidates:

- `tools/vibit`
- `tools/vibit.py`
- `package.json`
- `src/`

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
- YAML manifest loading if supported by selected language/tooling
- Deterministic reporting
- Module skeleton generation

## Tests

Depends on implementation language. At minimum:

- CLI help command
- Architecture check command
- Change check command
- Missing module check command

## Verification Commands

Pending implementation language decision.

## Rollback Or Migration Notes

Since this is the first CLI prototype, rollback is simply removing the CLI files and documentation references.
