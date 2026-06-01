# Product Maturity Milestones

Status: Accepted v0.1
Last updated: 2026-05-31
Scope: Product maturity path from source alpha to production-useful vibit
Depends on: `docs/v0.1-alpha-goal.md`, `docs/nakama-pitaya-product-parity-roadmap.md`, `docs/first-alpha-user-discovery-loop.md`
Canonical decision: `ADR-0105`

The paired Simplified Chinese translation is `docs/product-maturity-milestones.zh-CN.md`. The English file is authoritative.

This document turns the maintainer's product-stage intent into durable milestones. The current source-first alpha proves that vibit has a real backend loop, but the product goal is larger: move from first alpha to prototype-ready foundation, then to a single-node production-candidate foundation, and finally to a Nakama-first open-source server framework where AI-native development and AI-native testing are the default user experience. This document is roadmap and feedback-triage guidance only. It does not authorize runtime behavior changes, protocol route changes, Protobuf source or generated output changes, migrations, dependencies, hosted deployments, additional release artifacts, public announcements beyond the GitHub release record, paid promotion, broad operations/admin behavior, authentication/session behavior changes, broad product module expansion, Pitaya-style distributed architecture, or direct Nakama/Pitaya API compatibility.

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
stage_4_nakama_first_ai_native_product: long_term_target
stage_4_nakama_pitaya_class_product: long_term_target
reference_posture_update: ADR-0127
primary_product_reference: Nakama
pitaya_reference_status: serializer_message_forwarding_vocabulary_boundary_defined_for_future_architecture_planning
ai_native_development_testing_goal: user_requirement_to_spec_tests_implementation_verification
feedback_intake_surface: .github/ISSUE_TEMPLATE/alpha-feedback.yml
feedback_intake_standard: docs/first-alpha-feedback-intake-surfaces.md
next_direction: select_next_pitaya_aligned_direction_after_serializer_message_forwarding_map
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
pitaya_distributed_architecture_added: false
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
- a complete user-facing AI-native feature request, test, implementation, and verification workflow.

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
- AI-native requirement intake that converts user requests into specs, acceptance criteria, test plans, tests, implementation boundaries, verification records, and durable project memory.
- Basic concurrency and failure-path verification for login, protected requests, logout, reconnect, and database failure behavior.

Exit criteria:

- A technically capable external developer can build a small prototype flow on top of vibit without changing internal architecture boundaries first.
- The main setup friction is documented and intentionally accepted or reduced.
- The next missing capability is a product decision, not an accidental unknown.
- At least one non-maintainer feedback item has been triaged into a bounded work item or explicitly deferred.
- A non-trivial future feature can follow the documented AI-native requirement and test workflow without inventing the process ad hoc.

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
- AI-generated or AI-maintained feature tests that are traceable to requirement and acceptance artifacts.

Exit criteria:

- A team can evaluate vibit for a real game backend project without treating the framework as a research artifact.
- Known production risks are documented as accepted, fixed, or blocking.
- Upgrade and operations expectations are explicit.
- The remaining gap to product-class parity is mostly breadth, not core reliability posture.

## 5. Stage 4: Nakama-First AI-Native Product

Status:

```text
stage_4_nakama_first_ai_native_product: long_term_target
stage_4_nakama_pitaya_class_product: long_term_target
target: ai_era_nakama_pitaya_class_server_framework
primary_product_reference: Nakama
pitaya_reference_status: deferred_future_architecture_reference
```

Purpose: become a serious open-source server framework in the same broad product class as Nakama, adapted around vibit's agent-native maintainability model. Pitaya-style distributed architecture remains future-only until a later ADR explicitly reactivates it.

The product-class target has two inseparable parts:

- Nakama-style backend capability coverage.
- AI-native development and testing, where user requirements produce specs, acceptance criteria, tests, implementation, verification, and durable records through agent assistance.

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
- distributed runtime, frontend/backend roles, RPC, and service discovery, only after a later architecture ADR selects that path;
- agent-native requirement, acceptance, test, implementation, and verification workflow.

Stage 4 does not imply direct Nakama/Pitaya API compatibility unless a later ADR explicitly adopts a compatibility surface.

## 6. Feedback Triage By Stage

Early user feedback should be mapped to one of these maturity buckets:

- `source_alpha_friction`: README, setup, checks, request-loop, runbook, or concept clarity blocks.
- `prototype_ready_gap`: missing shared online service, example flow, client ergonomics, or local development path blocks a prototype.
- `production_candidate_gap`: security, operations, packaging, migration, observability, performance, or failure behavior blocks real project evaluation.
- `product_class_gap`: social, competitive, matchmaking, match runtime, SDK, admin console, distributed runtime, extensibility breadth, or AI-native feature workflow gaps block Nakama-first product-class usefulness.
- `out_of_scope_for_now`: request is valid but must wait for explicit authorization or a later stage.

Feedback that asks for production claims, broad feature parity, hosted deployment, direct compatibility, binary/package/container publication, paid promotion, or public announcement should not be silently accepted. It should be routed through maintainer authorization.

## 7. Next Product Direction

The next product direction after selecting the post-session-routing Pitaya-aligned direction is:

```text
W-0264 Select next Pitaya-aligned direction after serializer and message forwarding map
```

The agent-native feature request and test workflow is recorded in `docs/agent-native-feature-request-test-workflow.md` and `ADR-0128`. `ADR-0164` defined the Pitaya-aligned cluster-safe session routing boundary gate with check rule `runtime.pitaya_aligned_cluster_safe_session_routing_boundary_gate`; `ADR-0165` implemented `node tools/vibit inspect pitaya-sessions --json` with check rule `runtime.pitaya_aligned_cluster_safe_session_routing_source_first_map`; `ADR-0166` selected `define_pitaya_aligned_route_handler_pipeline_boundary_gate` with check rule `runtime.next_pitaya_aligned_direction_after_cluster_safe_session_routing_map`; `ADR-0167` defined the Pitaya-aligned route handler pipeline boundary gate with check rule `runtime.pitaya_aligned_route_handler_pipeline_boundary_gate`; `ADR-0168` implemented `node tools/vibit inspect pitaya-routes --json` with check rule `runtime.pitaya_aligned_route_handler_pipeline_source_first_map`; `ADR-0169` selected `define_pitaya_aligned_serializer_message_forwarding_boundary_gate` with check rule `runtime.next_pitaya_aligned_direction_after_route_handler_pipeline_map`; `ADR-0170` defined the Pitaya-aligned serializer and message forwarding boundary gate with check rule `runtime.pitaya_aligned_serializer_message_forwarding_boundary_gate`; and `ADR-0171` implemented `node tools/vibit inspect pitaya-serializer-forwarding --json` with check rule `runtime.pitaya_aligned_serializer_message_forwarding_source_first_map`. The next work is `W-0264 Select next Pitaya-aligned direction after serializer and message forwarding map`; it should only select the next bounded direction while keeping serializer behavior, message forwarding behavior, route handler implementation, handler routing behavior, handler pipeline behavior, pipeline middleware behavior, backend route targeting, protocol shape changes, service behavior changes, repository interface changes, PostgreSQL adapter changes, migrations, dependencies, stream subscriptions, chat rooms, groups, parties, distributed group behavior, group membership registries, broadcast fanout, delivery guarantees, distributed runtime implementation, frontend/backend server role implementation, server-to-server RPC implementation, remote calls, service discovery implementation, service registries, service selectors, node identity, cluster-safe session routing behavior, session location registries, connection owner node registries, routing epoch behavior, session route targets, remote connection handoff, distributed session routing, matchmaking, match runtime, SDK publication, hosted deployments, release artifacts, public announcements, event/audit tables, admin endpoints, metrics endpoints, observability pipelines, dashboards, and direct compatibility deferred unless a later explicit work item authorizes them.

Stage 2 trace references include `docs/prototype-ready-local-development-path-package.md`, `docs/storage-objects-behavior-gate.md`, `docs/storage-objects-persistence-schema-gate.md`, `runtime/migrations/postgres/000006_create_storage_objects.sql`, `docs/storage-objects-repository-boundary.md`, `runtime/migrations/postgres/000007_create_friend_relationships.sql`, and `docs/friends-relationship-repository-boundary.md`.

Friends relationship trace references include `decisions/ADR-0142-friends-relationship-repository-boundary.md`, `decisions/ADR-0144-friends-relationship-postgresql-adapter-gate.md`, `decisions/ADR-0145-friends-relationship-postgresql-adapter-implementation.md`, `decisions/ADR-0146-friends-relationship-runtime-behavior-gate.md`, `decisions/ADR-0147-friends-relationship-runtime-behavior-implementation.md`, `decisions/ADR-0148-friends-relationship-protocol-route-gate.md`, `decisions/ADR-0149-friends-relationship-protocol-route-implementation.md`, and `decisions/ADR-0150-friends-relationship-protocol-route-local-proof.md`.

The prior directions `W-0198 Define prototype-ready foundation execution plan`, `W-0199 Define prototype-ready local development path gate`, `W-0200 Implement prototype-ready local development path package`, `W-0201 Define storage objects behavior gate`, `W-0202 Define storage objects persistence schema gate`, `W-0203 Add storage objects migration source`, `W-0204 Define storage objects repository boundary`, `W-0205 Implement storage-neutral storage objects repository interface`, `W-0206 Define storage objects PostgreSQL adapter gate`, `W-0207 Implement storage objects PostgreSQL adapter`, `W-0208 Define storage objects runtime behavior gate`, `W-0209 Implement storage objects runtime behavior`, `W-0210 Define storage objects protocol route gate`, `W-0211 Implement storage objects protocol route`, `W-0212 Prove storage objects protocol route in local alpha request flow`, `W-0213 Confirm next alpha direction after storage objects local proof`, `W-0214 Define first server push and realtime messaging gate`, and `W-0215 Implement first server push and realtime messaging runtime slice` are completed and remain the trace from feedback intake into the Stage 2 execution plan and its first product capability path.

The recorded candidate focus areas are now:

- harden presence/status local proof through close and offline cases;
- strengthen concurrency and failure-path verification for the existing authenticated loop;
- implement the local alpha example client or example app path;
- define the minimum operations inspection surface needed before serious prototype use.

Future work must preserve ask-first boundaries for runtime behavior, protocol, generated output, dependencies, repository interfaces, storage adapters, operations/admin breadth, release artifacts, hosted deployments, public announcements, authentication/session behavior changes, broad product modules, large object/blob storage, S3-compatible object storage, Pitaya-style distributed architecture, and direct Nakama/Pitaya compatibility.

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
- add Pitaya-style distributed architecture before a later ADR reactivates it;
- add direct Nakama/Pitaya API compatibility;
- declare the current alpha production-ready.
