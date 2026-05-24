# Plan

## Intake

- [x] Preserve the original request in `request.md`.
- [x] Clarify user-visible outcome, non-goals, unknowns, and acceptance criteria.
- [x] Map the request to `agent_native_requirement_test_implementation_workflow`.
- [x] Confirm Pitaya remains deferred.

## Files To Create

- [x] `changes/_template/feature-request/request.md`
- [x] `changes/_template/feature-request/spec.yaml`
- [x] `changes/_template/feature-request/impact.md`
- [x] `changes/_template/feature-request/plan.md`
- [x] `changes/_template/feature-request/checklist.md`
- [x] `changes/_template/feature-request/verification.md`
- [x] `docs/agent-native-feature-request-scaffolding.md`
- [x] `docs/agent-native-feature-request-scaffolding.zh-CN.md`
- [x] `decisions/ADR-0137-agent-native-feature-request-scaffolding-implementation.md`
- [x] `conversations/2026-05-24-agent-native-feature-request-scaffolding-implementation.md`
- [x] `changes/2026-05-24-implement-agent-native-feature-request-scaffolding/`

## Files To Edit

- [x] `tools/vibit`
- [x] `rules/check-rules.json`
- [x] `.arch/work-items.yaml`
- [x] `.arch/runtime.yaml`
- [x] `.arch/reference.yaml`
- [x] `.arch/conventions.yaml`
- [x] `.arch/contracts.yaml`
- [x] `.arch/modules.yaml`
- [x] `modules/storage/module.yaml`
- [x] `AGENTS.md`
- [x] `AGENTS.zh-CN.md`
- [x] `runtime/AGENTS.md`
- [x] `runtime/AGENTS.zh-CN.md`
- [x] `modules/storage/AGENTS.md`
- [x] `modules/storage/AGENTS.zh-CN.md`
- [x] `README.md`
- [x] `README.zh-CN.md`
- [x] `docs/v0.1-alpha-goal.md`
- [x] `docs/v0.1-alpha-goal.zh-CN.md`
- [x] `docs/product-maturity-milestones.md`
- [x] `docs/product-maturity-milestones.zh-CN.md`
- [x] `docs/nakama-pitaya-product-parity-roadmap.md`
- [x] `docs/nakama-pitaya-product-parity-roadmap.zh-CN.md`

## Contracts And Schemas

- No runtime contracts or schemas are changed.
- The new command is a repository tooling command only.

## Implementation Boundary

- Allowed: docs, templates, tooling, checks, manifests, ADR, change memory, conversation memory, and agent guides.
- Forbidden: runtime, protocol, Protobuf, generated output, migration, dependency, persistence, startup wiring, SDK, hosted, distributed runtime, and direct compatibility scope.

## Tests

- [x] Tool syntax check.
- [x] Scaffold dry-run check.
- [x] Real scaffold invocation for this change directory.
- [x] Rule inspection.
- [x] Change/work/runtime/memory/schema/all checks.
- [x] Diff whitespace check.

## Verification Commands

- `node -c tools/vibit`
- `node tools/vibit scaffold feature scaffold-smoke --date 2026-05-24 --request "Smoke test feature request scaffold." --summary "Smoke test feature request scaffold." --dry-run`
- `node tools/vibit inspect rule runtime.agent_native_feature_request_scaffolding_implementation`
- `node tools/vibit check change implement-agent-native-feature-request-scaffolding --json`
- `node tools/vibit check work --json`
- `node tools/vibit check runtime --json`
- `node tools/vibit check memory --json`
- `node tools/vibit check schemas --json`
- `node tools/vibit check all --json`
- `git diff --check`

## Rollback Or Migration Notes

Rollback is limited to removing the template directory, command implementation, checks, and W-0229 metadata. No data migration or runtime rollback is needed.
