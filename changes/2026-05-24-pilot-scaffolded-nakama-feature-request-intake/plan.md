# Plan

## Intake

- [x] Preserve the original request in `request.md`.
- [x] Clarify user-visible outcome, non-goals, unknowns, and acceptance criteria.
- [x] Map the request to the Nakama `friends_groups_and_parties` capability family.
- [x] Confirm Pitaya remains deferred.

## Files To Create

- `decisions/ADR-0138-scaffolded-nakama-feature-request-intake-pilot.md`
- `conversations/2026-05-24-scaffolded-nakama-feature-request-intake-pilot.md`

## Files To Edit

- `changes/2026-05-24-pilot-scaffolded-nakama-feature-request-intake/*`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `.arch/modules.yaml`
- `modules/storage/module.yaml`
- `AGENTS.md`
- `AGENTS.zh-CN.md`
- `runtime/AGENTS.md`
- `runtime/AGENTS.zh-CN.md`
- `modules/storage/AGENTS.md`
- `modules/storage/AGENTS.zh-CN.md`
- `README.md`
- `README.zh-CN.md`
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
- `rules/check-rules.json`
- `tools/vibit`

## Contracts And Schemas

This pilot records future contract expectations only. W-0231 must define friendship commands, queries, events, errors, permissions, invariants, route/protocol posture, generated output posture, and migration posture before implementation.

## Implementation Boundary

- Allowed: source-first change artifacts, ADR, conversation memory, architecture manifests, repository docs, check-rule catalog, and `tools/vibit` repository checks.
- Forbidden: runtime behavior, protocol routes, Protobuf source, generated output, migrations, dependencies, SDK publication, hosted deployments, release artifacts, distributed runtime, Pitaya distributed architecture, direct Nakama/Pitaya API compatibility, or friendship implementation code.

## Tests

- [x] Positive behavior tests planned for future W-0231.
- [x] Negative behavior tests planned for future W-0231.
- [x] Permission/authentication tests planned for future W-0231.
- [x] Persistence/protocol/integration tests deferred with rationale.
- [x] Repository checks planned for this pilot.

## Verification Commands

- `node -c tools/vibit`
- `node tools/vibit inspect next --json`
- `node tools/vibit inspect rule runtime.scaffolded_nakama_feature_request_intake_pilot`
- `node tools/vibit check change pilot-scaffolded-nakama-feature-request-intake --json`
- `node tools/vibit check work --json`
- `node tools/vibit check runtime --json`
- `node tools/vibit check memory --json`
- `node tools/vibit check schemas --json`
- `node tools/vibit check all --json`
- `git diff --check`

## Rollback Or Migration Notes

No runtime, data, migration, dependency, or generated-output rollback is required because this pilot only changes source-first planning, memory, and repository checks. Reversal would update the work queue and ADRs before any friendship behavior implementation starts.
