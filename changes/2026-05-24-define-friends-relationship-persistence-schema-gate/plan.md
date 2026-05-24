# Plan

## Intake

- [x] Preserve the original request in `request.md`.
- [x] Clarify user-visible outcome, non-goals, unknowns, and acceptance criteria.
- [x] Map the request to Nakama `friends_groups_and_parties`.
- [x] Confirm Pitaya remains deferred.

## Files To Create

- [x] `docs/friends-relationship-persistence-schema-gate.md`
- [x] `docs/friends-relationship-persistence-schema-gate.zh-CN.md`
- [x] `decisions/ADR-0140-friends-relationship-persistence-schema-gate.md`
- [x] `conversations/2026-05-24-friends-relationship-persistence-schema-gate.md`
- [x] `changes/2026-05-24-define-friends-relationship-persistence-schema-gate/`

## Files To Edit

- [x] `.arch/work-items.yaml`
- [x] `.arch/runtime.yaml`
- [x] `.arch/conventions.yaml`
- [x] `.arch/contracts.yaml`
- [x] `.arch/reference.yaml`
- [x] `.arch/modules.yaml`
- [x] `modules/storage/module.yaml`
- [x] `README.md`
- [x] `README.zh-CN.md`
- [x] `AGENTS.md`
- [x] `AGENTS.zh-CN.md`
- [x] `runtime/AGENTS.md`
- [x] `runtime/AGENTS.zh-CN.md`
- [x] `modules/storage/AGENTS.md`
- [x] `modules/storage/AGENTS.zh-CN.md`
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
- [x] `tools/vibit`
- [x] `rules/check-rules.json`

## Contracts And Schemas

No public contract source, Protobuf source, generated output, repository interface, SQL migration source, or runtime schema is added.

The gate records a future `friend_relationships` table candidate and the migration source candidate `runtime/migrations/postgres/000007_create_friend_relationships.sql`.

## Implementation Boundary

Allowed: documentation, ADR, conversation memory, change artifacts, architecture manifests, README/AGENTS continuation pointers, rule catalog, and repository check tooling.

Forbidden: SQL migration source creation, `friend_relationships` table creation, runtime behavior, protocol routes, Protobuf source, generated output, repository interfaces, PostgreSQL adapters, dependencies, startup wiring, authentication/session changes, chat, groups, parties, matchmaking, match runtime, operations/admin behavior, SDK publication, hosted deployments, release artifacts, distributed runtime, and direct compatibility.

## Tests

- [x] Positive behavior tests are not applicable because no behavior changes.
- [x] Negative behavior tests are not applicable because no behavior changes.
- [x] Permission/authentication tests remain future runtime behavior requirements.
- [x] Persistence checks are planned for `W-0233`.
- [x] Repository checks cover this gate.

## Verification Commands

- `node -c tools/vibit`
- `node tools/vibit inspect next --json`
- `node tools/vibit inspect rule runtime.friends_relationship_persistence_schema_gate`
- `node tools/vibit check change define-friends-relationship-persistence-schema-gate --json`
- `node tools/vibit check work --json`
- `node tools/vibit check runtime --json`
- `node tools/vibit check memory --json`
- `node tools/vibit check schemas --json`
- `node tools/vibit check all --json`
- `git diff --check`

## Rollback Or Migration Notes

No migration source or runtime behavior is added, so rollback is documentation/check-rule reversal only.
