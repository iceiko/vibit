# Plan

1. Add parser support for generated inspection filters.
2. Filter expected generated contract shapes by module and type.
3. Preserve unfiltered output behavior.
4. Update agent tooling docs.
5. Run focused inspection and agent tooling checks.

## Files To Edit

- `tools/vibit`
- `docs/agent-tooling.md`
- `docs/agent-tooling.zh-CN.md`
- `.arch/work-items.yaml`

## Verification Commands

- `node -c tools/vibit`
- `node tools/vibit inspect generated --module inventory --json`
- `node tools/vibit inspect generated --type command --json`
- `node tools/vibit check agent-tooling --json`
- `node tools/vibit check work --json`
