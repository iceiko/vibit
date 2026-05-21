# Plan

## Files To Create

- `docs/release-execution-authorization-gate.md`
- `docs/release-execution-authorization-gate.zh-CN.md`
- `decisions/ADR-0100-release-execution-authorization-gate.md`
- `conversations/2026-05-21-release-execution-authorization-gate.md`
- `changes/2026-05-21-define-release-execution-authorization-gate/`

## Files To Edit

- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `README.md`
- `README.zh-CN.md`
- `docs/v0.1-alpha-goal.md`
- `docs/v0.1-alpha-goal.zh-CN.md`
- `docs/alpha-developer-flow.md`
- `docs/alpha-developer-flow.zh-CN.md`
- `docs/alpha-acceptance-checklist.md`
- `docs/alpha-acceptance-checklist.zh-CN.md`
- `docs/runtime-runbook.md`
- `docs/runtime-runbook.zh-CN.md`
- `docs/release-execution-preparation-gate.md`
- `docs/release-execution-preparation-gate.zh-CN.md`
- `AGENTS.md`
- `AGENTS.zh-CN.md`
- `runtime/AGENTS.md`
- `runtime/AGENTS.zh-CN.md`
- `tools/vibit`
- `rules/check-rules.json`

## Generated Artifacts

None.

## Handwritten Logic

No runtime handwritten logic is included. Tooling changes are limited to repository checks for the release execution authorization gate.

## Tests

Use existing Go tests and `tools/vibit` checks.

## Verification Commands

```bash
node -c tools/vibit
node tools/vibit inspect next
node tools/vibit check change define-release-execution-authorization-gate --json
node tools/vibit check work --json
node tools/vibit check runtime --json
node tools/vibit check memory --json
node tools/vibit check schemas --json
node tools/vibit check all --json
cd runtime && go test ./cmd/vibit-server
cd runtime && go test ./...
examples/local-alpha-request-loop.sh
git diff --check
```

## Rollback Or Migration Notes

Rollback is a documentation/tooling revert. No data migration rollback is needed.
