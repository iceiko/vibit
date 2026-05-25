# Plan

## Intake

- [x] Preserve the original request in `request.md`.
- [x] Clarify user-visible outcome, non-goals, unknowns, and acceptance criteria.
- [x] Map the request to the Nakama `friends_groups_and_parties` capability family.
- [x] Confirm Pitaya remains deferred unless a later ADR explicitly reactivates it.

## Files To Create

- [x] `docs/friends-relationship-repository-boundary.md`
- [x] `docs/friends-relationship-repository-boundary.zh-CN.md`
- [x] `decisions/ADR-0142-friends-relationship-repository-boundary.md`
- [x] `conversations/2026-05-25-friends-relationship-repository-boundary.md`

## Files To Edit

- [x] `.arch/work-items.yaml`
- [x] `.arch/runtime.yaml`
- [x] `.arch/conventions.yaml`
- [x] `.arch/contracts.yaml`
- [x] `.arch/reference.yaml`
- [x] `.arch/modules.yaml`
- [x] `README.md`
- [x] `README.zh-CN.md`
- [x] `docs/v0.1-alpha-goal.md`
- [x] `docs/v0.1-alpha-goal.zh-CN.md`
- [x] `docs/alpha-developer-flow.md`
- [x] `docs/alpha-developer-flow.zh-CN.md`
- [x] `docs/alpha-acceptance-checklist.md`
- [x] `docs/alpha-acceptance-checklist.zh-CN.md`
- [x] `docs/product-maturity-milestones.md`
- [x] `docs/product-maturity-milestones.zh-CN.md`
- [x] `docs/nakama-pitaya-product-parity-roadmap.md`
- [x] `docs/nakama-pitaya-product-parity-roadmap.zh-CN.md`
- [x] `AGENTS.md`
- [x] `AGENTS.zh-CN.md`
- [x] `runtime/AGENTS.md`
- [x] `runtime/AGENTS.zh-CN.md`
- [x] `modules/storage/module.yaml`
- [x] `modules/storage/AGENTS.md`
- [x] `modules/storage/AGENTS.zh-CN.md`
- [x] `tools/vibit`
- [x] `rules/check-rules.json`

## Contracts And Schemas

No public contracts, Protobuf source, generated output, or migrations are added. The change records future repository vocabulary only.

## Implementation Boundary

Allowed:

- boundary docs;
- ADR and conversation memory;
- change artifacts;
- manifests and guides;
- static repository checks.

Forbidden:

- repository interface implementation;
- PostgreSQL adapter implementation;
- runtime behavior;
- protocol routes;
- Protobuf source;
- generated output;
- migrations;
- dependencies;
- event/audit tables;
- SDK, hosted, distributed runtime, release artifact, or direct compatibility scope.

## Tests

- [x] Positive behavior tests: not applicable because no runtime behavior is added.
- [x] Negative behavior tests: static forbidden-scope checks.
- [x] Permission/authentication tests: not applicable because identity proof remains deferred.
- [x] Persistence/protocol/integration tests: not applicable because no adapter, protocol route, migration, or runtime behavior is added.
- [x] Repository checks: `runtime.friends_relationship_repository_boundary`.

## Verification Commands

- `node -c tools/vibit`
- `node tools/vibit inspect next --json`
- `node tools/vibit inspect rule runtime.friends_relationship_repository_boundary`
- `node tools/vibit check change define-friends-relationship-repository-boundary --json`
- `node tools/vibit check work --json`
- `node tools/vibit check runtime --json`
- `node tools/vibit check memory --json`
- `node tools/vibit check schemas --json`
- `node tools/vibit check all --json`
- `cd runtime && go test ./...`
- `git diff --check`

## Rollback Or Migration Notes

No migration is added or changed. Rolling back this slice means reverting documentation, ADR, conversation, manifests, static checks, and W-0234/W-0235 workflow metadata.
