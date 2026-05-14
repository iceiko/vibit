# Plan

1. Run focused inspect commands.
2. Run focused checks.
3. Add missing change specs required by work queue validation.
4. Mark `W-0040` completed after checks pass.
5. Add exactly one next-ready follow-up work item for continued tooling hardening.
6. Run repository checks and Go tests.

## Files To Edit

- `.arch/work-items.yaml`
- Change spec files under this directory.

## Verification Commands

- `node tools/vibit inspect next --json`
- `node tools/vibit inspect contracts --json`
- `node tools/vibit inspect generated --json`
- `node tools/vibit inspect reference --json`
- `node tools/vibit check agent-tooling --json`
- `node tools/vibit check generated --json`
- `node tools/vibit check work --json`
- `node tools/vibit check all --json`
- `cd runtime && go test ./...`
- `cd runtime && go vet ./...`
- `git diff --check`
