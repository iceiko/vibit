# Nakama-First Product Capability Roadmap

Status: Draft v0.2
Last updated: 2026-05-31
Scope: Product roadmap standard for vibit's Nakama-first game backend target

The paired Simplified Chinese translation is `docs/nakama-pitaya-product-parity-roadmap.zh-CN.md`. The English file is authoritative.

## 1. Purpose

This standard originally upgraded the reference posture from "learn from Nakama and Pitaya" to a product target. `ADR-0127` refines that posture: Nakama is now the primary product capability reference, while Pitaya is deferred to a future architecture reference for distributed Go game server concerns.

```text
nakama_pitaya_product_parity_roadmap: ratified
completed_work_item: W-0168
decision: ADR-0078
check_rule: runtime.reference_product_parity_roadmap
reference_posture_update: ADR-0127
primary_product_reference: Nakama
pitaya_reference_status: distributed_group_broadcast_vocabulary_boundary_defined_for_future_architecture_planning
parity_goal: nakama_first_same_class_common_capability_coverage
api_compatibility_goal: false
direct_nakama_pitaya_api_compatibility_added: false
implementation_authorized_by_this_standard: roadmap_only
```

Historical `ADR-0078` check markers are retained for traceability. They are superseded as current planning guidance by `ADR-0127`, but older repository checks still use them to confirm the original roadmap lineage:

```text
parity_goal: nakama_pitaya_same_class_common_capability_coverage
recommended_next_direction: define_protocol_logout_route_gate
second_direction: define_transport_close_handoff_gate
third_direction: define_reconnect_connection_epoch_gate
fourth_direction: define_protocol_session_carrier_gate
first_module_expansion_after_lifecycle: define_presence_lifecycle_gate
```

Historical `W-0220` workflow markers are retained for repository check traceability. They record the direction that opened the now-completed workflow pilot, not the current next-ready item:

```text
recommended_next_direction: pilot_nakama_aligned_feature_request_workflow
```

Historical `W-0221` workflow pilot markers are retained for repository check traceability. They record the direction that opened the now-completed presence/status proof-hardening slice, not the current next-ready item:

```text
recommended_next_direction: presence_status_local_proof_hardening
```

vibit should become a Nakama-class open-source game backend framework with AI-native development and AI-native testing as the product purpose. This means vibit must cover the common product capability families that game teams expect from Nakama-style systems while preserving vibit's core differentiator: users state backend requirements in ordinary product language, and AI agents turn those requirements into bounded specs, acceptance criteria, tests, implementation, verification records, ADRs, manifests, and repository checks.

Product parity means comparable capability coverage and operational usefulness. It does not mean direct API compatibility, copied public routes, copied data models, copied clustering internals, or a commitment to follow every implementation detail from Nakama. Pitaya-style frontend/backend, RPC, service discovery, cluster groups, and distributed topology remain future concerns until a later ADR reintroduces Pitaya as an active architecture reference.

## 2. Product Target

Nakama is the primary reference for broad game backend product coverage:

- Accounts, authentication, users, and sessions.
- Storage objects and durable game state.
- Friends, groups, parties, chat, status, presence, and notifications.
- Leaderboards, tournaments, rewards, currencies, and other metagame systems.
- Match listing, matchmaking, realtime multiplayer, and authoritative match runtime.
- Server runtime customization, hooks, RPC-like extension points, streams, console, metrics, and operations.
- Client SDK and sample application ergonomics.

Pitaya is not a current product planning driver. It is a deferred future architecture reference for later distributed Go game server topology:

- Acceptors and connection lifecycle.
- Sessions, binding, kick/disconnect, and session data.
- Handler routing, pipelines, serializers, and message forwarding.
- Groups, broadcast, multicast, and push.
- Frontend/backend server roles.
- Server-to-server RPC, service discovery, cluster mode, monitoring, and tracing.

Agents should not use Pitaya to pull cluster/RPC/frontend-backend work into the current prototype-ready foundation. Pitaya vocabulary may remain useful when documenting why transport, session, protocol, application, and backend service concerns stay separated, but Nakama decides the product capability priorities for the current roadmap.

vibit must adapt these capabilities into its own model:

- Contract-first public behavior.
- Module-owned invariants.
- Generated repeatable structure.
- Application-owned lifecycle policy.
- Transport/protocol/domain separation.
- Repository and persistence boundaries.
- Agent-readable guides and checkable architecture.
- AI-native requirement intake, test planning, implementation, and verification.

## 2.1 AI-Native Product Purpose

vibit's product purpose is not only to run a game backend. Its product purpose is to make backend development itself AI-native.

When a user describes a new backend requirement, the intended product workflow is:

```text
user requirement
-> AI-written bounded requirement spec
-> AI-written acceptance criteria
-> AI-written test plan
-> AI-written or updated tests
-> AI implementation inside declared boundaries
-> AI-run verification
-> AI-updated docs, manifests, ADRs, and change records
```

The architecture exists to make that workflow reliable. Contracts, module ownership, generated files, repository checks, redaction rules, and verification commands are not internal bureaucracy; they are the mechanism that lets an AI agent safely convert user requirements into tested backend behavior.

Every meaningful future feature should explain:

- the user requirement it satisfies;
- the Nakama-style product capability family it maps to;
- the acceptance criteria a user can judge;
- the positive and negative tests that prove behavior;
- the implementation boundary;
- the verification commands that were run;
- the artifacts updated so future AI agents can continue safely.

## 3. Parity Capability Families

The following capability families are first-class roadmap scope:

```text
parity_capability_families:
  - identity_authentication_sessions
  - connection_lifecycle_reconnect_logout
  - storage_objects_and_durable_game_state
  - presence_status_and_notifications
  - chat_streams_and_realtime_messaging
  - friends_groups_and_parties
  - leaderboards_tournaments_and_competitive_systems
  - economy_inventory_rewards_currencies_and_progression
  - matchmaking_match_listing_and_room_lifecycle
  - realtime_multiplayer_and_authoritative_match_runtime
  - server_runtime_hooks_rpc_and_custom_logic
  - admin_console_metrics_observability_and_operations
  - client_sdks_examples_and_developer_experience
  - distributed_runtime_frontend_backend_rpc_and_service_discovery
  - agent_native_requirement_test_implementation_workflow
```

Each family must eventually have:

- A module or runtime subsystem owner.
- A semantic contract surface where public behavior exists.
- A protocol surface where client/server messages exist.
- A storage boundary where durable state exists.
- Tests for invariants and error behavior.
- Repository checks or architecture checks for the most important boundaries.
- English and Simplified Chinese documentation when public-facing.
- Requirement, acceptance, and test-plan records for non-trivial user-facing changes.

## 4. Current Status

Current completed foundation:

- Agent-native project governance, change specs, ADRs, conversation memory, and checks.
- Go runtime layout.
- WebSocket transport and Protobuf envelope.
- PostgreSQL persistence and migration tooling.
- Inventory proof slice.
- Player account repository and PostgreSQL adapter.
- Device credential login, opaque access-token validation, logout, runtime session persistence, session validation, and route protection.
- First-message connection binding.
- Single-process active connection registry.
- Single-process WebSocket close policy that invalidates registry records without concrete socket close handoff.

Current gap:

- The source-first alpha is visible, and the prototype-ready foundation execution plan has advanced through storage objects and the first realtime outbound delivery foundation.
- The next gap is to make the AI-native requirement-to-test-to-implementation loop explicit before broadening product modules. Without this, agents can keep adding slices, but users do not yet receive the promised "state a requirement and have AI handle spec, tests, implementation, verification, and records" workflow.

## 5. Phase Plan

### Phase 2R: Runtime Lifecycle Closure

Status:

```text
phase_2r_runtime_lifecycle_closure: active
current_near_term_priority: protocol_logout_and_connection_lifecycle
```

Goal: finish the login, route protection, runtime session, connection binding, logout, close intent, concrete close, reconnect, and session-carrier loop before expanding product modules.

Required near-term gates:

1. `define_protocol_logout_route_gate`
2. `define_transport_close_handoff_gate`
3. `define_reconnect_connection_epoch_gate`
4. `define_protocol_session_carrier_gate`
5. `define_presence_lifecycle_gate`
6. `strengthen_operations_observability_and_admin_tooling`

This phase absorbs Nakama's explicit account/session/socket lifecycle pressure before higher-level social or multiplayer modules depend on unstable lifecycle behavior. Transport, protocol, application, and backend service separation remains a vibit architecture rule; it no longer requires Pitaya to be an active near-term planning driver.

### Phase 3: Shared Online Services

Goal: add the common always-on game backend services that many features depend on.

Candidate work:

- Storage object contracts and permissions beyond module-local inventory state.
- Presence and status lifecycle.
- Notifications.
- Chat, streams, and realtime messaging.
- Server push and broadcast vocabulary.
- Admin inspection for players, sessions, tokens, and active connections.

### Phase 2P: Prototype-Ready Foundation

Status:

```text
phase_2p_prototype_ready_foundation: next_product_stage
standard: docs/product-maturity-milestones.md
execution_plan: docs/prototype-ready-foundation-execution-plan.md
local_development_path_gate: docs/prototype-ready-local-development-path-gate.md
local_development_path_package: docs/prototype-ready-local-development-path-package.md
storage_objects_behavior_gate: docs/storage-objects-behavior-gate.md
storage_objects_persistence_schema_gate: docs/storage-objects-persistence-schema-gate.md
storage_objects_repository_boundary: docs/storage-objects-repository-boundary.md
storage_objects_runtime_behavior_implementation: runtime/internal/app/storage/service.go
storage_objects_protocol_route_gate: docs/storage-objects-protocol-route-gate.md
storage_objects_protocol_route_gate_decision: ADR-0118
storage_objects_protocol_route_implementation: proto/vibit/storage/v1/storage.proto
storage_objects_protocol_route_implementation_decision: ADR-0119
storage_objects_protocol_route_local_proof: runtime/internal/platform/protocol/protobuf/authenticated_gameplay_e2e_test.go
storage_objects_protocol_route_local_proof_decision: ADR-0120
first_server_push_realtime_messaging_gate: docs/first-server-push-realtime-messaging-gate.md
first_server_push_realtime_messaging_gate_decision: ADR-0122
first_server_push_realtime_messaging_runtime_slice: runtime/internal/app/realtime/service.go
first_server_push_realtime_messaging_runtime_slice_decision: ADR-0123
next_alpha_direction_after_realtime_runtime_slice_decision: ADR-0124
realtime_protocol_websocket_outbound_delivery_implementation: proto/vibit/realtime/v1/realtime.proto
realtime_protocol_websocket_outbound_delivery_implementation_decision: ADR-0126
agent_native_feature_request_test_workflow: docs/agent-native-feature-request-test-workflow.md
agent_native_feature_request_test_workflow_decision: ADR-0128
next_nakama_prototype_ready_capability_selection_decision: ADR-0132
agent_native_feature_request_scaffolding_gate: docs/agent-native-feature-request-scaffolding-gate.md
agent_native_feature_request_scaffolding_gate_decision: ADR-0136
scaffolded_nakama_feature_request_intake_pilot_decision: ADR-0138
friends_relationship_lifecycle_gate_decision: ADR-0139
friends_relationship_migration_source_decision: ADR-0141
friends_relationship_repository_boundary_decision: ADR-0142
friends_relationship_protocol_route_gate_decision: ADR-0148
friends_relationship_protocol_route_implementation_decision: ADR-0149
friends_relationship_protocol_route_local_proof_decision: ADR-0150
minimum_operations_inspection_surface_gate_decision: ADR-0152
minimum_operations_inspection_surface_gate_check_rule: runtime.minimum_operations_inspection_surface_gate
minimum_operations_inspection_source_first_surface_implementation_decision: ADR-0153
minimum_operations_inspection_source_first_surface_implementation_check_rule: runtime.minimum_operations_inspection_source_first_surface_implementation
pitaya_aligned_server_to_server_rpc_boundary_gate_decision: ADR-0158
pitaya_aligned_server_to_server_rpc_boundary_gate_check_rule: runtime.pitaya_aligned_server_to_server_rpc_boundary_gate
pitaya_aligned_server_to_server_rpc_source_first_map_decision: ADR-0159
pitaya_aligned_server_to_server_rpc_source_first_map_check_rule: runtime.pitaya_aligned_server_to_server_rpc_source_first_map
pitaya_aligned_service_discovery_boundary_gate_decision: ADR-0160
pitaya_aligned_service_discovery_boundary_gate_check_rule: runtime.pitaya_aligned_service_discovery_boundary_gate
pitaya_aligned_service_discovery_source_first_map_decision: ADR-0161
pitaya_aligned_service_discovery_source_first_map_check_rule: runtime.pitaya_aligned_service_discovery_source_first_map
historical_after_adr_0127:
  pitaya_reference_status: deferred_future_architecture_reference
next_work_item: W-0257 Implement Pitaya-aligned cluster-safe session routing source-first map
```

Goal: move from a source-first alpha that developers can inspect to a foundation they can use for a serious small prototype.

Candidate work:

- Reduce local setup, migration, and configuration friction. Completed by `W-0200`.
- Define and then add a clearer example client or example app path. Completed by `W-0225` and `W-0226`.
- Define first general storage-object behavior beyond the inventory proof slice. Completed by `W-0201`.
- Define first storage objects persistence schema posture. Completed by `W-0202`.
- Add the first storage objects migration source. Completed by `W-0203`.
- Define the storage objects repository boundary. Completed by `W-0204`.
- Implement the storage-neutral storage objects repository interface. Completed by `W-0205`.
- Define and implement the storage objects PostgreSQL adapter. Completed by `W-0206` and `W-0207`.
- Define and implement storage objects runtime behavior. Completed by `W-0208` and `W-0209`.
- Define and implement the storage objects protocol route family. Completed by `W-0210` and `W-0211`.
- Prove storage object routes through the local alpha request flow. Completed by `W-0212`.
- Confirm the next alpha direction after storage object local proof. Completed by `W-0213`.
- Define the first server push and realtime messaging gate. Completed by `W-0214`.
- Implement the first server push and realtime messaging runtime slice. Completed by `W-0215`.
- Confirm the next alpha direction after the realtime runtime slice. Completed by `W-0216`.
- Implement the realtime protocol and WebSocket outbound delivery slice. Completed by `W-0218`.
- Confirm the next alpha direction after realtime outbound delivery. Completed by `W-0219`.
- Define the agent-native feature request and test workflow. Completed by `W-0220`.
- Pilot the Nakama-aligned feature request workflow. Completed by `W-0221`.
- Harden presence/status local proof through close and offline cases. Next.
- Strengthen concurrency and failure-path verification around the authenticated gameplay loop.
- Define the minimal operations inspection surface needed before serious prototype use.

This phase does not declare production readiness. It selects the next smallest product-useful slices while preserving source-alpha honesty and the existing ask-first boundaries.

### Phase 4: Social And Competitive Modules

Goal: cover Nakama-style metagame and social surfaces.

Candidate modules:

- Friends.
- Groups.
- Parties.
- Leaderboards.
- Tournaments.
- Wallet/currency.
- Rewards/claims.
- Quests/progression.

### Phase 5: Matchmaking And Match Runtime

Goal: cover the multiplayer capability families expected from Nakama. Routing, room, and broadcast implementation details remain vibit-owned until a future architecture ADR reintroduces Pitaya-style distributed vocabulary.

Candidate work:

- Match listing.
- Matchmaking tickets and criteria.
- Room lifecycle.
- Broadcast groups and target scopes.
- Authoritative match runtime contracts.
- Relayed realtime multiplayer contracts.
- Reconnect/replay decisions.

### Phase 6: Runtime Extensibility And Developer Experience

Goal: make vibit usable as a framework, not only as a server codebase.

Candidate work:

- Server runtime hooks.
- Server-side custom logic/RPC surface.
- Module scaffolding and generator hardening.
- Client SDKs and examples.
- Local development runbook.
- Admin console and CLI workflows.

### Phase 7: Distributed Runtime

Goal: introduce distributed topology only after single-process semantics are stable and after a later ADR decides whether Pitaya is the right active architecture reference again.

Candidate work:

- Frontend/backend server role split.
- Server-to-server RPC.
- Service discovery.
- Distributed groups and push.
- Cluster-safe session routing.
- Multi-node presence and matchmaking behavior.

This phase must not start by weakening the single-process contracts. It should lift proven single-process semantics into distributed adapters.

## 6. Development Method

Future work should use this ordering for each capability family:

1. Requirement intake: restate the user's requirement, user-visible behavior, non-goals, and risks.
2. Nakama reference review: record the Nakama product capability being covered or explicitly note when no Nakama mapping applies.
3. Acceptance criteria: define user-visible success and failure behavior before implementation.
4. Test plan: define positive, negative, permission, failure-path, persistence, protocol, and integration tests that apply to the slice.
5. Vibit ownership: decide the module, runtime, platform, generated, and operations owners.
6. Semantic contract: define commands, queries, events, errors, permissions, and invariants.
7. Protocol contract: define wire messages only after semantic behavior is stable.
8. Persistence boundary: define tables, repositories, indexes, and redaction before adapters.
9. Application behavior: implement the smallest vertical slice.
10. Operations surface: add inspection, metrics, admin actions, or runbook guidance when the behavior needs operators.
11. Verification: add focused Go tests and repository checks, then run them.
12. Memory: update change specs, ADRs, manifests, and conversation logs.

The preferred implementation shape remains:

```text
user requirement -> spec -> acceptance criteria -> test plan -> tests -> contract -> generated shape -> logic -> checks -> docs
```

## 7. Near-Term Recommendation

After authenticated failure-path proof, next capability selection, the example client path gate, the example client path implementation, the follow-up scaffolding selection, the feature request scaffolding gate, the feature request scaffolding implementation, the scaffolded Nakama intake pilot, the friends relationship lifecycle gate, the friends relationship persistence schema gate, the friends relationship migration source, the friends relationship repository boundary, the storage-neutral friends relationship repository interface implementation, the friends relationship PostgreSQL adapter gate, the friends relationship PostgreSQL adapter implementation, the friends relationship runtime behavior gate, the friends relationship runtime behavior implementation, the friends relationship protocol route gate, the friends relationship protocol route implementation, the friends relationship protocol route local proof, the post-proof next capability selection, the minimum operations inspection surface gate, the source-first operations inspection implementation, the Pitaya distributed runtime vocabulary map, the frontend/backend role map, the server-to-server RPC boundary gate, the server-to-server RPC source-first map, the service discovery boundary gate, and the service discovery source-first map, the next concrete work should not jump directly to chat, groups, matchmaking, match runtime, SDK publication, hosted demos, distributed runtime implementation, broad social modules, event/audit tables, generated client libraries, new protocol shape, admin endpoints, metrics endpoints, observability pipelines, dashboards, frontend/backend server role implementation, server-to-server RPC implementation, remote calls, service discovery implementation, service registries, service selectors, node identity, distributed groups, room broadcast fanout, or cluster-safe session routing. The next concrete work should define only a gate-only Pitaya-aligned distributed group and broadcast boundary:

```text
next_work_item: W-0257 Implement Pitaya-aligned cluster-safe session routing source-first map
recommended_next_direction: implement_pitaya_aligned_cluster_safe_session_routing_source_first_map
primary_product_reference: Nakama
pitaya_reference_status: distributed_group_broadcast_vocabulary_boundary_defined_for_future_architecture_planning
ai_native_development_testing_goal: user_requirement_to_spec_tests_implementation_verification
```

Rationale:

- The project now has an explicit user-facing AI-native development workflow and has piloted it through presence/status and authenticated failure-path proof.
- The existing source-first alpha capabilities now have a readable local example path and source-first feature request scaffold.
- `ADR-0138` used that scaffold on a bounded Nakama request intake, selected `friends_groups_and_parties`, and opened `W-0231 Define friends relationship lifecycle gate` before any friendship implementation work.
- `ADR-0139` defined the friend relationship lifecycle semantic gate and opened `W-0232 Define friends relationship persistence schema gate`.
- `ADR-0140` defined the friends relationship persistence schema gate and opened `W-0233 Add friends relationship migration source`.
- `ADR-0141` added only the friends relationship migration source and opened `W-0234 Define friends relationship repository boundary`.
- `ADR-0142` defined the friends relationship repository boundary and opened `W-0235 Implement storage-neutral friends relationship repository interface`.
- `ADR-0143` implemented the storage-neutral friends relationship repository interface, registered `runtime.friends_relationship_repository_interface_implementation`, and opened `W-0236 Define friends relationship PostgreSQL adapter gate`.
- `ADR-0144` defined the friends relationship PostgreSQL adapter gate, registered `runtime.friends_relationship_postgresql_adapter_gate`, and opened `W-0237 Implement friends relationship PostgreSQL adapter`.
- `ADR-0145` implemented the friends relationship PostgreSQL adapter, registered `runtime.friends_relationship_postgresql_adapter_implementation`, and opened `W-0238 Define friends relationship runtime behavior gate`.
- `ADR-0146` defined the friends relationship runtime behavior gate, registered `runtime.friends_relationship_runtime_behavior_gate`, and opened `W-0239 Implement friends relationship runtime behavior`.
- `ADR-0147` implemented the friends relationship runtime behavior service, registered `runtime.friends_relationship_runtime_behavior_implementation`, and opened `W-0240 Define friends relationship protocol route gate`.
- `ADR-0148` defined the friends relationship protocol route gate, registered `runtime.friends_relationship_protocol_route_gate`, and opened `W-0241 Implement friends relationship protocol route`.
- `ADR-0149` implemented the friends relationship protocol route, registered `runtime.friends_relationship_protocol_route_implementation`, and opened `W-0242 Prove friends relationship protocol route in local alpha request flow`.
- `ADR-0150` proved the friends relationship protocol route locally and registered `runtime.friends_relationship_protocol_route_local_proof`.
- `ADR-0152` defined the minimum operations inspection surface gate, registered `runtime.minimum_operations_inspection_surface_gate`, and opened `W-0245 Implement minimum operations inspection source-first surface`.
- `ADR-0153` implemented the source-first operations inspection surface, registered `runtime.minimum_operations_inspection_source_first_surface_implementation`, recorded the Pitaya deferred architecture map, and opened `W-0246 Define Pitaya-aligned distributed runtime vocabulary reactivation gate`.
- `ADR-0154` defined the Pitaya-aligned distributed runtime vocabulary reactivation gate, registered `runtime.pitaya_aligned_distributed_runtime_vocabulary_reactivation_gate`, and opened `W-0247 Implement Pitaya-aligned distributed runtime vocabulary source-first map`.
- `ADR-0155` implemented `node tools/vibit inspect pitaya-vocabulary --json`, registered `runtime.pitaya_aligned_distributed_runtime_vocabulary_source_first_map`, and opened `W-0248 Define Pitaya-aligned frontend/backend role boundary gate`.
- `ADR-0156` defined the Pitaya-aligned frontend/backend role boundary gate, registered `runtime.pitaya_aligned_frontend_backend_role_boundary_gate`, and opened `W-0249 Implement Pitaya-aligned frontend/backend role source-first map`.
- `ADR-0157` implemented `node tools/vibit inspect pitaya-roles --json`, registered `runtime.pitaya_aligned_frontend_backend_role_source_first_map`, and opened `W-0250 Define Pitaya-aligned server-to-server RPC boundary gate`.
- `ADR-0158` defined the Pitaya-aligned server-to-server RPC boundary gate, registered `runtime.pitaya_aligned_server_to_server_rpc_boundary_gate`, and opened `W-0251 Implement Pitaya-aligned server-to-server RPC source-first map`.
- `ADR-0159` implemented `node tools/vibit inspect pitaya-rpc --json`, registered `runtime.pitaya_aligned_server_to_server_rpc_source_first_map`, and opened `W-0252 Define Pitaya-aligned service discovery boundary gate`.
- `ADR-0160` defined the Pitaya-aligned service discovery boundary gate, registered `runtime.pitaya_aligned_service_discovery_boundary_gate`, and opened `W-0253 Implement Pitaya-aligned service discovery source-first map`.
- `ADR-0161` implemented `node tools/vibit inspect pitaya-discovery --json`, registered `runtime.pitaya_aligned_service_discovery_source_first_map`, and opened `W-0254 Define Pitaya-aligned distributed group and broadcast boundary gate`.
- `ADR-0162` defined the Pitaya-aligned distributed group and broadcast boundary gate, registered `runtime.pitaya_aligned_distributed_group_broadcast_boundary_gate`, and opened `W-0257 Implement Pitaya-aligned cluster-safe session routing source-first map`.
- Friend relationship lifecycle is a core Nakama-class social graph primitive, and its protected protocol route family is now proven through the local alpha request flow.
- The operations inspection surface now makes local alpha status, route families, verification posture, and deferred operations surfaces inspectable without adding runtime endpoints.
- The service discovery source-first map now makes future Pitaya-aligned service discovery, registry, instance, and selector vocabulary inspectable before groups, parties, chat targeting, invites, matchmaking filters, or match runtime social context depend on broader scope.
- The next source-first map should make distributed group and broadcast vocabulary inspectable without adding group implementation, group membership registries, room broadcast fanout, delivery guarantees, distributed runtime behavior, service discovery behavior, RPC behavior, protocol shape, persistence, dependencies, or direct compatibility.
- Nakama-first product planning still prevents near-term scope from being split between product breadth and Pitaya-style distributed implementation.
- Pitaya-style cluster/RPC/frontend-backend concerns should stay deferred to vocabulary records and source-first maps until a later explicit implementation work item authorizes them.

## 8. Non-Goals

This roadmap does not:

- Implement any new runtime behavior by itself.
- Add protocol logout routes.
- Add concrete WebSocket close handoff.
- Add reconnect, resume, duplicate replacement, or session carrier behavior.
- Add presence, chat, friends, groups, parties, leaderboards, tournaments, matchmaking, or match runtime code.
- Add admin console, SDKs, cluster, RPC, service discovery, or distributed runtime behavior.
- Add dependencies.
- Commit to Nakama or Pitaya API compatibility.
- Treat Pitaya as a current product driver before a later ADR reactivates it.
- Let AI implement a non-trivial user request without a spec, acceptance criteria, test plan, tests or explicit verification rationale.

## 9. Agent Rules

Agents must:

- Treat Nakama-class common capability coverage as a product requirement, not only as background reading.
- Treat AI-native requirement intake, testing, implementation, and verification as the product purpose, not only as an internal workflow.
- Map new major work to one roadmap family.
- Prefer lifecycle and shared services before high-level social or multiplayer modules.
- Preserve vibit's agent-native constraints even when matching reference capabilities.
- Record adopted, adapted, deferred, or rejected reference patterns in ADRs or change specs.
- Keep direct API compatibility behind an explicit future compatibility ADR.

Agents must not:

- Add product features directly in WebSocket transport handlers.
- Add direct external API compatibility accidentally.
- Skip semantic contracts because a reference system already has a similar feature.
- Start distributed runtime work before single-process semantics are tested.
- Treat product parity as permission to weaken generated-file, redaction, permission, session, or module-boundary rules.
- Skip acceptance criteria or test planning for non-trivial user-facing requirements.
