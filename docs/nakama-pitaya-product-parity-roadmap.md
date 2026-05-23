# Nakama And Pitaya Product Parity Roadmap

Status: Draft v0.1
Last updated: 2026-05-18
Scope: Product roadmap standard for vibit's Nakama/Pitaya-class game backend target

The paired Simplified Chinese translation is `docs/nakama-pitaya-product-parity-roadmap.zh-CN.md`. The English file is authoritative.

## 1. Purpose

This standard upgrades the reference posture from "learn from Nakama and Pitaya" to a product target:

```text
nakama_pitaya_product_parity_roadmap: ratified
completed_work_item: W-0168
decision: ADR-0078
check_rule: runtime.reference_product_parity_roadmap
parity_goal: nakama_pitaya_same_class_common_capability_coverage
api_compatibility_goal: false
direct_nakama_pitaya_api_compatibility_added: false
implementation_authorized_by_this_standard: roadmap_only
```

vibit should become a Nakama/Pitaya-class open-source game backend framework. This means vibit must cover the common product capability families that game teams expect from those systems, while preserving vibit's core differentiator: agent-native maintainability through explicit contracts, manifests, generation, tests, ADRs, and repository checks.

Product parity means comparable capability coverage and operational usefulness. It does not mean direct API compatibility, copied public routes, copied data models, copied clustering internals, or a commitment to follow every implementation detail from either project.

## 2. Product Target

Nakama remains the primary reference for broad game backend product coverage:

- Accounts, authentication, users, and sessions.
- Storage objects and durable game state.
- Friends, groups, parties, chat, status, presence, and notifications.
- Leaderboards, tournaments, rewards, currencies, and other metagame systems.
- Match listing, matchmaking, realtime multiplayer, and authoritative match runtime.
- Server runtime customization, hooks, RPC-like extension points, streams, console, metrics, and operations.
- Client SDK and sample application ergonomics.

Pitaya remains the primary reference for Go game server architecture vocabulary:

- Acceptors and connection lifecycle.
- Sessions, binding, kick/disconnect, and session data.
- Handler routing, pipelines, serializers, and message forwarding.
- Groups, broadcast, multicast, and push.
- Frontend/backend server roles.
- Server-to-server RPC, service discovery, cluster mode, monitoring, and tracing.

vibit must adapt these capabilities into its own model:

- Contract-first public behavior.
- Module-owned invariants.
- Generated repeatable structure.
- Application-owned lifecycle policy.
- Transport/protocol/domain separation.
- Repository and persistence boundaries.
- Agent-readable guides and checkable architecture.

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
```

Each family must eventually have:

- A module or runtime subsystem owner.
- A semantic contract surface where public behavior exists.
- A protocol surface where client/server messages exist.
- A storage boundary where durable state exists.
- Tests for invariants and error behavior.
- Repository checks or architecture checks for the most important boundaries.
- English and Simplified Chinese documentation when public-facing.

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

- The source-first alpha is visible, but the next product stage still needs an explicit prototype-ready foundation execution plan.
- The prototype-ready foundation execution plan is recorded, and the next work must define the local development path gate before implementation slices broaden shared online services.

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

This phase absorbs Nakama's explicit account/session/socket lifecycle pressure and Pitaya's acceptor/session/handler separation before higher-level social or multiplayer modules depend on unstable lifecycle behavior.

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
next_work_item: W-0213 Confirm next alpha direction after storage objects local proof
```

Goal: move from a source-first alpha that developers can inspect to a foundation they can use for a serious small prototype.

Candidate work:

- Reduce local setup, migration, and configuration friction. Completed by `W-0200`.
- Add a clearer example client or example app path.
- Define first general storage-object behavior beyond the inventory proof slice. Completed by `W-0201`.
- Define first storage objects persistence schema posture. Completed by `W-0202`.
- Add the first storage objects migration source. Completed by `W-0203`.
- Define the storage objects repository boundary. Completed by `W-0204`.
- Implement the storage-neutral storage objects repository interface. Completed by `W-0205`.
- Define and implement the storage objects PostgreSQL adapter. Completed by `W-0206` and `W-0207`.
- Define and implement storage objects runtime behavior. Completed by `W-0208` and `W-0209`.
- Define and implement the storage objects protocol route family. Completed by `W-0210` and `W-0211`.
- Prove storage object routes through the local alpha request flow. Completed by `W-0212`.
- Confirm the next alpha direction after storage object local proof. Next.
- Define first realtime messaging, stream, broadcast, or server-push behavior.
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

Goal: cover the multiplayer capability families expected from Nakama and the routing/group semantics informed by Pitaya.

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

Goal: introduce Pitaya-style distributed topology only after single-process semantics are stable.

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

1. Reference review: record the Nakama/Pitaya capability being covered.
2. Vibit ownership: decide the module, runtime, platform, generated, and operations owners.
3. Semantic contract: define commands, queries, events, errors, permissions, and invariants.
4. Protocol contract: define wire messages only after semantic behavior is stable.
5. Persistence boundary: define tables, repositories, indexes, and redaction before adapters.
6. Application behavior: implement the smallest vertical slice.
7. Operations surface: add inspection, metrics, admin actions, or runbook guidance when the behavior needs operators.
8. Verification: add focused Go tests and repository checks.
9. Memory: update change specs, ADRs, manifests, and conversation logs.

The preferred implementation shape remains:

```text
requirement -> reference mapping -> spec -> contract -> generated shape -> logic -> tests -> checks -> docs
```

## 7. Near-Term Recommendation

After this roadmap is ratified, the next concrete work should not jump directly to chat, groups, matchmaking, or match runtime. The next concrete work should finish the runtime lifecycle foundation:

```text
recommended_next_direction: define_protocol_logout_route_gate
second_direction: define_transport_close_handoff_gate
third_direction: define_reconnect_connection_epoch_gate
fourth_direction: define_protocol_session_carrier_gate
first_module_expansion_after_lifecycle: define_presence_lifecycle_gate
```

Rationale:

- Logout exists as service behavior but is not exposed as a client protocol route.
- Close policy can invalidate registry records but cannot yet close concrete WebSocket sockets.
- Reconnect and duplicate connection behavior should not be designed after presence or match runtime depends on it.
- Protocol session carriers must be explicit before clients can safely reason about runtime sessions.
- Presence and chat depend on connection/session lifecycle semantics.

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

## 9. Agent Rules

Agents must:

- Treat Nakama/Pitaya-class common capability coverage as a product requirement, not only as background reading.
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
