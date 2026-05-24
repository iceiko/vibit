# Plan

## Intake

- [x] Preserve the original request in `request.md`.
- [x] Clarify user-visible outcome, non-goals, unknowns, and acceptance criteria.
- [x] Map the request to the Nakama `friends_groups_and_parties` capability family.
- [x] Confirm Pitaya remains deferred unless a later ADR explicitly reactivates it.

## Files To Create

- `docs/friends-relationship-lifecycle-gate.md`
- `docs/friends-relationship-lifecycle-gate.zh-CN.md`
- `decisions/ADR-0139-friends-relationship-lifecycle-gate.md`
- `conversations/2026-05-24-friends-relationship-lifecycle-gate.md`

## Files To Edit

- `changes/2026-05-24-define-friends-relationship-lifecycle-gate/*`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `.arch/modules.yaml`
- `modules/storage/module.yaml`
- `README.md`
- `README.zh-CN.md`
- `AGENTS.md`
- `AGENTS.zh-CN.md`
- `runtime/AGENTS.md`
- `runtime/AGENTS.zh-CN.md`
- `modules/storage/AGENTS.md`
- `modules/storage/AGENTS.zh-CN.md`
- `docs/v0.1-alpha-goal.md`
- `docs/v0.1-alpha-goal.zh-CN.md`
- `docs/alpha-developer-flow.md`
- `docs/alpha-developer-flow.zh-CN.md`
- `docs/alpha-acceptance-checklist.md`
- `docs/alpha-acceptance-checklist.zh-CN.md`
- `docs/product-maturity-milestones.md`
- `docs/product-maturity-milestones.zh-CN.md`
- `docs/nakama-pitaya-product-parity-roadmap.md`
- `docs/nakama-pitaya-product-parity-roadmap.zh-CN.md`
- `tools/vibit`
- `rules/check-rules.json`

## Contracts And Schemas

- Record future command, query, event, error, permission, actor-relative state, invariant, and test vocabulary in the gate.
- Do not create contract source files in this slice.
- Do not create protocol payloads, generated output, migrations, repositories, adapters, or runtime handlers in this slice.

## Implementation Boundary

- Allowed: docs, ADR, conversation log, change artifacts, architecture manifests, repository guides, roadmap docs, `tools/vibit`, and rule catalog updates required to register the gate and move continuation to W-0232.
- Forbidden: runtime behavior, protocol routes, Protobuf source, generated output, migrations, repository interfaces, PostgreSQL adapters, dependencies, startup wiring, authentication/session changes, delivery guarantees, stream subscriptions, chat rooms, groups, parties, broadcast fanout, matchmaking, match runtime, operations/admin behavior, SDK publication, generated client libraries, hosted deployments, release artifacts, distributed runtime, and direct compatibility.

## Tests

- [x] Positive future behavior tests recorded.
- [x] Negative future behavior tests recorded.
- [x] Permission/authentication future tests recorded.
- [x] Persistence tests deferred to W-0232 with rationale.
- [x] Protocol/integration tests deferred until route/runtime work is authorized.
- [x] Repository checks planned.

## Verification Commands

- `node -c tools/vibit`
- `node tools/vibit inspect next --json`
- `node tools/vibit inspect rule runtime.friends_relationship_lifecycle_gate`
- `node tools/vibit check change define-friends-relationship-lifecycle-gate --json`
- `node tools/vibit check work --json`
- `node tools/vibit check runtime --json`
- `node tools/vibit check memory --json`
- `node tools/vibit check schemas --json`
- `node tools/vibit check all --json`
- `git diff --check`

## Rollback Or Migration Notes

This gate adds no migration source and no runtime behavior. Rollback is documentation/tooling state only: remove the gate docs, ADR, conversation, change artifacts, manifest progression, and rule registration if the decision is reversed before dependent work begins.
