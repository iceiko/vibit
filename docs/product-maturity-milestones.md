# Product Maturity Milestones

Status: Accepted v0.1
Last updated: 2026-05-22
Scope: Product maturity path from source alpha to production-useful vibit
Depends on: `docs/v0.1-alpha-goal.md`, `docs/nakama-pitaya-product-parity-roadmap.md`, `docs/first-alpha-user-discovery-loop.md`
Canonical decision: `ADR-0105`

The paired Simplified Chinese translation is `docs/product-maturity-milestones.zh-CN.md`. The English file is authoritative.

This document turns the maintainer's product-stage intent into durable milestones. The current source-first alpha proves that vibit has a real backend loop, but the product goal is larger: move from first alpha to prototype-ready foundation, then to a single-node production-candidate foundation, and finally to a Nakama/Pitaya-class open-source server framework. This document is roadmap and feedback-triage guidance only. It does not authorize runtime behavior changes, protocol route changes, Protobuf source or generated output changes, migrations, dependencies, hosted deployments, additional release artifacts, public announcements beyond the GitHub release record, paid promotion, broad operations/admin behavior, authentication/session behavior changes, broad product module expansion, or direct Nakama/Pitaya API compatibility.

## 1. Core Rule

The product maturity record is:

```yaml
product_maturity_milestones: defined
completed_work_item: W-0197
decision: ADR-0105
check_rule: runtime.first_alpha_feedback_intake_surfaces
current_stage: source_first_alpha
current_release_identifier: v0.1.0-alpha.1
stage_1_source_first_alpha: reached
stage_2_prototype_ready_foundation: next_product_stage
stage_3_single_node_production_candidate_foundation: planned
stage_4_nakama_pitaya_class_product: long_term_target
feedback_intake_surface: .github/ISSUE_TEMPLATE/alpha-feedback.yml
feedback_intake_standard: docs/first-alpha-feedback-intake-surfaces.md
next_direction: prototype_ready_foundation_execution_plan
runtime_behavior_added: false
protocol_route_added: false
protobuf_source_added: false
generated_output_added: false
migration_added: false
dependency_added: false
hosted_deployment_added: false
additional_release_artifacts_authorized: false
public_announcements_beyond_github_release_authorized: false
paid_promotion_authorized: false
broad_operations_admin_behavior_added: false
authentication_session_behavior_changed: false
product_module_expansion_added: false
direct_nakama_pitaya_api_compatibility_added: false
```

## 2. Stage 1: Source-First Alpha

Status:

```text
stage_1_source_first_alpha: reached
release_identifier: v0.1.0-alpha.1
```

Purpose: prove that vibit is no longer only a design exercise.

The current alpha should be judged as a source-first developer alpha. It should attract technically capable developers who can run a local path, inspect the architecture, and give concrete feedback.

Stage 1 exists when a developer can:

- clone the repository;
- run repository checks;
- run Go tests;
- run the local alpha request loop;
- understand the WebSocket + Protobuf + PostgreSQL posture;
- understand local onboarding, device credential login, access-token validation, runtime sessions, connection binding, protected inventory, protected presence query, and logout;
- see the current limitations before mistaking the alpha for a production server distribution.

Stage 1 is already reached by `v0.1.0-alpha.1`.

Stage 1 is not:

- production deployment readiness;
- packaged distribution readiness;
- SDK readiness;
- hosted platform readiness;
- feature parity with Nakama, Pitaya, Colyseus, Pomelo, Agones, or custom production backends.

## 3. Stage 2: Prototype-Ready Foundation

Status:

```text
stage_2_prototype_ready_foundation: next_product_stage
suggested_release_band: v0.2_or_v0.3
```

Purpose: make vibit credible for a serious multiplayer or realtime backend prototype.

This stage means a developer can use vibit as the backend base for a small proof game or product prototype without first inventing the missing common online-service layer.

Required capability groups:

- Storage objects or a similarly general durable game-state surface beyond the module-local inventory proof slice.
- Clear presence/status semantics beyond the minimum current protected query.
- A first server push, broadcast, stream, or realtime messaging vocabulary.
- A minimal chat or realtime messaging slice, or a recorded decision that another shared online service is the more urgent prototype unlock.
- Better local startup ergonomics, including explicit setup, migration, and configuration flow.
- A realistic example client or example app path that demonstrates more than one request in isolation.
- Issue and feedback loops that convert external friction into bounded work items.
- Basic concurrency and failure-path verification for login, protected requests, logout, reconnect, and database failure behavior.

Exit criteria:

- A technically capable external developer can build a small prototype flow on top of vibit without changing internal architecture boundaries first.
- The main setup friction is documented and intentionally accepted or reduced.
- The next missing capability is a product decision, not an accidental unknown.
- At least one non-maintainer feedback item has been triaged into a bounded work item or explicitly deferred.

Stage 2 is still not production-ready. It may remain single-node, local-first, and source-first, but it should feel useful rather than merely inspectable.

## 4. Stage 3: Single-Node Production-Candidate Foundation

Status:

```text
stage_3_single_node_production_candidate_foundation: planned
suggested_release_band: v0.4_or_v0.5
```

Purpose: make vibit a serious single-node foundation that an external team could evaluate for real project development.

This stage does not mean broad distributed scalability or full game-backend product parity. It means the single-node server path has enough security, operations, packaging, and common backend capability to be a production-candidate foundation.

Required capability groups:

- Stable lifecycle semantics for login, token validation, runtime sessions, connection binding, logout, close handoff, reconnect, and session carriers.
- Hardened PostgreSQL path: migrations, indexes, transaction boundaries, connection handling, failure behavior, and upgrade notes.
- Production configuration posture for secrets, redaction, validation, fail-closed behavior, and environment separation.
- Observable operations baseline: structured logs, metrics or metrics boundary, health/readiness posture, and admin inspection path for players, sessions, tokens, and active connections.
- Release distribution beyond source archive, after explicit authorization: likely container image, checksums, versioned release notes, and upgrade documentation.
- Stable client ergonomics: SDK, client helper, or documented protocol client example.
- Security review for authentication, sessions, token redaction, configuration leakage, and route permissions.
- Concurrency, soak, and failure-mode verification appropriate to single-node usage.

Exit criteria:

- A team can evaluate vibit for a real game backend project without treating the framework as a research artifact.
- Known production risks are documented as accepted, fixed, or blocking.
- Upgrade and operations expectations are explicit.
- The remaining gap to product-class parity is mostly breadth, not core reliability posture.

## 5. Stage 4: Nakama/Pitaya-Class Product

Status:

```text
stage_4_nakama_pitaya_class_product: long_term_target
target: ai_era_nakama_pitaya_class_server_framework
```

Purpose: become a serious open-source server framework in the same broad product class as Nakama and Pitaya, adapted around vibit's agent-native maintainability model.

This stage requires common capability coverage across the roadmap families already recorded in `docs/nakama-pitaya-product-parity-roadmap.md`:

- identity, authentication, and sessions;
- connection lifecycle, reconnect, and logout;
- storage objects and durable game state;
- presence, status, and notifications;
- chat, streams, and realtime messaging;
- friends, groups, and parties;
- leaderboards, tournaments, and competitive systems;
- economy, inventory, rewards, currencies, and progression;
- matchmaking, match listing, and room lifecycle;
- realtime multiplayer and authoritative match runtime;
- server runtime hooks, RPC, and custom logic;
- admin console, metrics, observability, and operations;
- client SDKs, examples, and developer experience;
- distributed runtime, frontend/backend roles, RPC, and service discovery.

Stage 4 does not imply direct Nakama/Pitaya API compatibility unless a later ADR explicitly adopts a compatibility surface.

## 6. Feedback Triage By Stage

Early user feedback should be mapped to one of these maturity buckets:

- `source_alpha_friction`: README, setup, checks, request-loop, runbook, or concept clarity blocks.
- `prototype_ready_gap`: missing shared online service, example flow, client ergonomics, or local development path blocks a prototype.
- `production_candidate_gap`: security, operations, packaging, migration, observability, performance, or failure behavior blocks real project evaluation.
- `product_class_gap`: social, competitive, matchmaking, match runtime, SDK, admin console, distributed runtime, or extensibility breadth blocks Nakama/Pitaya-class usefulness.
- `out_of_scope_for_now`: request is valid but must wait for explicit authorization or a later stage.

Feedback that asks for production claims, broad feature parity, hosted deployment, direct compatibility, binary/package/container publication, paid promotion, or public announcement should not be silently accepted. It should be routed through maintainer authorization.

## 7. Next Product Direction

The next product direction after the local development path package is:

```text
W-0213 Confirm next alpha direction after storage objects local proof
```

The prototype-ready execution plan is now recorded in `docs/prototype-ready-foundation-execution-plan.md`, the local development path gate is recorded in `docs/prototype-ready-local-development-path-gate.md`, the local development path package is recorded in `docs/prototype-ready-local-development-path-package.md`, the storage objects behavior gate is recorded in `docs/storage-objects-behavior-gate.md`, the storage objects persistence schema gate is recorded in `docs/storage-objects-persistence-schema-gate.md`, the storage objects migration source is recorded in `runtime/migrations/postgres/000006_create_storage_objects.sql`, the storage objects repository boundary is recorded in `docs/storage-objects-repository-boundary.md`, the storage objects repository interface is recorded in `runtime/internal/modules/storage/repository.go`, the storage objects PostgreSQL adapter gate is recorded in `docs/storage-objects-postgresql-adapter-gate.md`, the storage objects PostgreSQL adapter is recorded in `runtime/internal/platform/persistence/postgres/storage_object_repository.go`, the storage objects runtime behavior gate is recorded in `docs/storage-objects-runtime-behavior-gate.md`, the storage objects runtime behavior implementation is recorded in `runtime/internal/app/storage/service.go`, the storage objects protocol route gate is recorded in `docs/storage-objects-protocol-route-gate.md` and `ADR-0118`, the storage objects protocol route implementation is recorded in `proto/vibit/storage/v1/storage.proto`, `runtime/internal/app/bootstrap/storage.go`, `runtime/internal/platform/protocol/protobuf/storage_bridge.go`, and `ADR-0119`, and the storage objects protocol route local proof is recorded in `runtime/internal/platform/protocol/protobuf/authenticated_gameplay_e2e_test.go`, `examples/local-alpha-request-loop.sh`, and `ADR-0120`. The next work should confirm one bounded alpha direction after storage object local proof, keeping Nakama-style capability planning and Pitaya-style layering alignment without implementing new runtime behavior or direct compatibility.

The prior directions `W-0198 Define prototype-ready foundation execution plan`, `W-0199 Define prototype-ready local development path gate`, `W-0200 Implement prototype-ready local development path package`, `W-0201 Define storage objects behavior gate`, `W-0202 Define storage objects persistence schema gate`, `W-0203 Add storage objects migration source`, `W-0204 Define storage objects repository boundary`, `W-0205 Implement storage-neutral storage objects repository interface`, `W-0206 Define storage objects PostgreSQL adapter gate`, `W-0207 Implement storage objects PostgreSQL adapter`, `W-0208 Define storage objects runtime behavior gate`, `W-0209 Implement storage objects runtime behavior`, `W-0210 Define storage objects protocol route gate`, `W-0211 Implement storage objects protocol route`, and `W-0212 Prove storage objects protocol route in local alpha request flow` are completed and remain the trace from feedback intake into the Stage 2 execution plan and its first product capability path.

The recorded candidate focus areas remain:

- reduce local setup friction;
- add a clearer example client or example app path;
- prove the storage objects protocol route in the local alpha request flow;
- define first realtime messaging or server-push slice;
- strengthen concurrency and failure-path verification for the existing authenticated loop;
- define the minimum operations inspection surface needed before serious prototype use.

Future work must preserve ask-first boundaries for runtime behavior, protocol, generated output, dependencies, repository interfaces, storage adapters, operations/admin breadth, release artifacts, hosted deployments, public announcements, authentication/session behavior changes, broad product modules, large object/blob storage, S3-compatible object storage, and direct Nakama/Pitaya compatibility.

## 8. Non-Authorization

This document records product maturity goals. It does not itself:

- implement runtime behavior;
- add protocol routes;
- add Protobuf source or generated output;
- add migrations;
- add dependencies;
- publish binaries, packages, containers, checksums, signing/provenance artifacts, install scripts, registry artifacts, or hosted deployments;
- authorize broad public announcements or paid promotion;
- add operations/admin behavior;
- change authentication/session behavior;
- add broad product modules;
- add direct Nakama/Pitaya API compatibility;
- declare the current alpha production-ready.
