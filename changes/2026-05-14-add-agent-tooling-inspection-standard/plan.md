# Plan

1. Add `docs/agent-tooling.md`.
2. Add `docs/agent-tooling.zh-CN.md`.
3. Add the new inspect commands.
4. Add `check agent-tooling`.
5. Include the check in `check all`.
6. Update AGENTS guidance, conventions, runtime/reference manifests, generated-output docs, and rule catalog.
7. Run focused checks.

## Files To Edit

- `tools/vibit`
- `docs/agent-tooling.md`
- `docs/agent-tooling.zh-CN.md`
- `docs/generated-output.md`
- `docs/generated-output.zh-CN.md`
- `AGENTS.md`
- `AGENTS.zh-CN.md`
- `.arch/conventions.yaml`
- `.arch/runtime.yaml`
- `.arch/reference.yaml`
- `rules/check-rules.json`
- `.arch/work-items.yaml`

## Generated Artifacts

None.

## Verification Commands

- `node -c tools/vibit`
- `node tools/vibit inspect next --json`
- `node tools/vibit inspect contracts --json`
- `node tools/vibit inspect generated --json`
- `node tools/vibit inspect reference --json`
- `node tools/vibit check agent-tooling --json`
