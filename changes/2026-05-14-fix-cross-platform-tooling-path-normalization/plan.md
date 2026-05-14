# Plan

## Files To Create

- `changes/2026-05-14-fix-cross-platform-tooling-path-normalization/`

## Files To Edit

- `tools/vibit`
- `docs/agent-tooling.md`
- `docs/agent-tooling.zh-CN.md`
- `AGENTS.md`
- `AGENTS.zh-CN.md`
- `.arch/conventions.yaml`

## Generated Artifacts

None.

## Handwritten Logic

Add path normalization helpers to `tools/vibit`:

- `toRepoPath(value)`
- `pathStartsWith(value, prefix)`
- `normalizeArtifact(value)`

Use these helpers for repository-relative path output and path boundary checks.

## Tests

No dedicated test harness exists for Windows in this environment. Verification should include syntax checks, agent-tooling checks, runtime checks, and full repository checks.

## Verification Commands

```bash
node -c tools/vibit
node tools/vibit check agent-tooling --json
node tools/vibit check runtime --json
node tools/vibit check change fix-cross-platform-tooling-path-normalization --json
node tools/vibit check all --json
git diff --check
```

## Rollback Notes

Rollback means restoring raw `path.relative()` behavior and raw path `startsWith` comparisons. That would reintroduce Windows false positives and platform-specific JSON paths.
