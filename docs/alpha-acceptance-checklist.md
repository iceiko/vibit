# Alpha Acceptance Checklist

Status: Draft v0.1
Last updated: 2026-05-21
Scope: Local v0.1 alpha acceptance criteria for vibit

The paired Simplified Chinese translation is `docs/alpha-acceptance-checklist.zh-CN.md`. The English file is authoritative.

This checklist is not a release declaration. It is the maintainer-facing and contributor-facing checklist for deciding whether the local `v0.1 alpha` developer flow is ready to inspect locally and consider for later publishing.

## 1. Purpose

The first alpha should be developer-usable before it is release-published. This checklist defines what must be true for the local flow to be accepted as alpha-ready:

- a contributor can understand what vibit is,
- prepare local prerequisites and configuration,
- run or verify the current runtime,
- understand the authenticated gameplay flow,
- run the repository checks,
- confirm secrets remain redacted,
- and find the next contribution entry point.

This checklist does not authorize release publishing, release packaging, runtime behavior changes, protocol changes, generated output changes, migrations, dependencies, broad operations/admin behavior, product module expansion, or direct Nakama/Pitaya API compatibility.

## 2. Checklist States

Use these states when reviewing the alpha flow:

- `Ready`: the repository already contains the required artifact or behavior and it is verified by tests, checks, or explicit documentation.
- `Manual`: the step requires a developer to perform local setup, such as installing tools or applying migrations.
- `Deferred`: the item is intentionally outside the first alpha or awaits a later work item.
- `Blocked`: the alpha cannot be accepted until the item is resolved.

## 3. Repository Intake

- [x] `README.md` and `README.zh-CN.md` explain that vibit is an agent-native server framework and that the current state is pre-alpha.
- [x] `docs/v0.1-alpha-goal.md` and `docs/v0.1-alpha-goal.zh-CN.md` define the short-term `v0.1 alpha` target.
- [x] `AGENTS.md`, `AGENTS.zh-CN.md`, `runtime/AGENTS.md`, and `runtime/AGENTS.zh-CN.md` point agents to the current continuation queue.
- [x] `.arch/work-items.yaml` records the current continuation state.
- [x] `node tools/vibit inspect next` identifies the current `W-0268 Define Pitaya-aligned session binding, kick/disconnect, and session data boundary gate` continuation step.
- [x] `node tools/vibit inspect operations --json` reports the source-first local operations posture, route families, redaction flags, and Pitaya deferred architecture mapping.
- [x] `node tools/vibit inspect pitaya-vocabulary --json` reports the source-first Pitaya vocabulary map and deferrals.
- [x] `node tools/vibit inspect pitaya-roles --json` reports the source-first Pitaya frontend/backend role map and deferrals.
- [x] `node tools/vibit inspect pitaya-rpc --json` reports the source-first Pitaya server-to-server RPC map and deferrals.
- [x] `node tools/vibit inspect pitaya-discovery --json` reports the source-first Pitaya service discovery map and deferrals.
- [x] `node tools/vibit inspect pitaya-groups --json` reports the source-first Pitaya distributed group and broadcast map and deferrals.
- [x] `node tools/vibit inspect pitaya-sessions --json` reports the source-first Pitaya cluster-safe session routing map and deferrals.
- [x] `node tools/vibit inspect pitaya-routes --json` reports the source-first Pitaya route handler pipeline map and deferrals.
- [x] `node tools/vibit inspect pitaya-serializer-forwarding --json` reports the source-first Pitaya serializer and message forwarding map and deferrals.
- [x] `node tools/vibit inspect pitaya-acceptor-connection --json` reports the source-first Pitaya acceptor and connection lifecycle map and deferrals.
- [x] The current `next_direction: define_pitaya_aligned_session_binding_kick_disconnect_session_data_boundary_gate` follows the completed friends relationship protocol route local proof, W-0243 post-proof selection, W-0244 minimum operations inspection surface gate, W-0245 source-first operations inspection implementation, W-0246 Pitaya vocabulary gate, W-0247 source-first vocabulary map, W-0248 frontend/backend role boundary gate, W-0249 source-first role map, W-0250 server-to-server RPC boundary gate, W-0251 server-to-server RPC source-first map, W-0252 service discovery boundary gate, W-0253 service discovery source-first map, W-0254 distributed group/broadcast boundary gate, W-0255 distributed group/broadcast source-first map, W-0256 cluster-safe session routing boundary gate, W-0257 cluster-safe session routing source-first map, W-0258 next Pitaya-aligned direction selection, W-0259 route handler pipeline boundary gate, W-0260 route handler pipeline source-first map, W-0261 next Pitaya-aligned direction selection after the route handler pipeline map, W-0262 serializer and message forwarding boundary gate, W-0263 serializer and message forwarding source-first map, W-0264 next Pitaya-aligned direction selection after the serializer and message forwarding map, W-0265 acceptor and connection lifecycle boundary gate, W-0266 acceptor and connection lifecycle source-first map, and W-0267 next Pitaya-aligned direction selection after the acceptor and connection lifecycle map.
- [x] `docs/prototype-ready-local-development-path-package.md` packages the repeatable source-first local path.

## 4. Local Prerequisites

- [ ] Manual: Go is installed for runtime tests and local server execution.
- [ ] Manual: Node.js is installed for `tools/vibit` checks.
- [ ] Manual: PostgreSQL is available when testing the persistent runtime path.
- [ ] Manual: Buf and Protobuf tooling are available when regenerating Protobuf output.
- [x] `docs/runtime-runbook.md` documents the memory and PostgreSQL startup paths.
- [x] The repository checks do not require live PostgreSQL by default.

## 5. Configuration And Secret Handling

- [x] `VIBIT_RUNTIME_STORE=memory` remains the default bootstrap path.
- [x] `VIBIT_RUNTIME_STORE=postgres` is documented as the current alpha runtime composition.
- [x] PostgreSQL runtime startup requires `VIBIT_POSTGRES_DSN` and verifier key environment variables.
- [x] `examples/local.prototype.env.example` provides a placeholder-only local configuration checklist.
- [x] `.gitignore` ignores `.vibit.local.env`, `.env.local`, and `.env.*.local`.
- [x] Verifier key material, raw device credentials, raw access tokens, DSNs, digests, headers, cookies, query strings, subprotocol values, remote addresses, and concrete transport metadata are documented as not log-safe.
- [x] The local alpha request-loop script does not print raw credentials, raw tokens, verifier keys, DSNs, digests, or transport metadata.
- [x] `/configz` reports only redacted posture and includes `secrets_redacted: true`.

## 6. Database And Migration Posture

- [x] PostgreSQL migration sources exist under `runtime/migrations/postgres`.
- [x] Normal runtime startup does not apply migrations automatically.
- [x] Migration apply/status tooling exists through the repository runtime tooling.
- [ ] Manual: a local PostgreSQL database is prepared before using `VIBIT_RUNTIME_STORE=postgres`.
- [ ] Manual: required SQL migrations are applied or verified before using a fresh PostgreSQL database.
- [ ] Manual: optional live PostgreSQL verification uses `VIBIT_POSTGRES_TEST_DSN` and a disposable database.

## 7. Runtime Surface

- [x] The gameplay WebSocket endpoint is `/v1/ws`.
- [x] `/v1/ws` expects binary `vibit.protocol.v1.Envelope` Protobuf bytes, not JSON.
- [x] `/healthz` reports process health.
- [x] `/readyz` reports readiness posture, runtime store, and WebSocket path.
- [x] `/version` reports the pre-alpha runtime version.
- [x] `/configz` reports only redacted runtime posture.
- [x] These HTTP status endpoints are local troubleshooting surfaces, not a production operations API, admin console, metrics backend, or gameplay protocol route.

## 8. Authenticated Gameplay Flow

- [x] Local onboarding exists as application service behavior through `OnboardLocalPlayerWithDeviceCredential`.
- [x] Local onboarding is not exposed as a public WebSocket, Protobuf, HTTP, CLI, or startup auto-creation surface.
- [x] Device credential login is exposed through the `runtime.authentication.AuthenticateWithDeviceCredential` protocol route.
- [x] Login returns an opaque access token and runtime session metadata.
- [x] First-message connection binding is exposed through the `runtime.authentication.BindConnection` protocol route.
- [x] Protected inventory grant/read requests use `AuthenticatedRequest`.
- [x] Protected presence query is available through `runtime.presence.GetPlayerPresence`.
- [x] Protected own-player storage object put/get/list/delete is proven through `storage.GetOwnStorageObject`, `storage.ListOwnStorageObjects`, `storage.PutOwnStorageObject`, and `storage.DeleteOwnStorageObject`.
- [x] Logout is exposed through the `runtime.authentication.LogoutAccessToken` protocol route.
- [x] A post-logout protected request using the same token is rejected.
- [x] The focused authenticated gameplay E2E test proves onboarding -> login -> connection binding -> protected inventory -> presence query -> logout -> post-logout rejection.

## 9. Verification Commands

Run from the repository root unless noted:

```bash
node tools/vibit inspect next
node tools/vibit inspect operations --json
node tools/vibit check work --json
node tools/vibit check all --json
git diff --check
cd runtime && go test ./...
examples/local-alpha-request-loop.sh
```

Optional focused checks:

```bash
cd runtime && go test ./cmd/vibit-server
cd runtime && go test ./internal/platform/protocol/protobuf -run TestAuthenticatedGameplayE2EUsesExistingOnboardingLoginBindingInventoryPresenceAndLogout -v
cd runtime && go test ./internal/platform/protocol/protobuf -run TestStorageObjectsProtocolRouteLocalAlphaFlow -v
```

Optional live PostgreSQL verification remains opt-in and requires a disposable database:

```bash
cd runtime && VIBIT_POSTGRES_TEST_DSN='postgres://user:pass@127.0.0.1:5432/vibit_test?sslmode=disable' VIBIT_POSTGRES_TEST_ALLOW_DESTRUCTIVE=1 go test ./internal/platform/protocol/protobuf -run TestPostgresPersistentInventoryRequestLoop -v
```

## 10. Contribution Entry Point

The alpha developer flow is now packaged in:

```text
docs/alpha-developer-flow.md
```

The release, discovery, prototype-ready, storage objects, realtime outbound delivery, Nakama-first direction, agent-native feature request/test workflow, Nakama-aligned presence/status pilot, presence/status local proof hardening, authenticated gameplay failure-path verification, next Nakama prototype-ready capability selection, local alpha example client path gate, local alpha example client path implementation, feature request scaffolding, scaffolded Nakama feature request intake pilot, friends relationship lifecycle gate, friends relationship persistence schema gate, friends relationship migration source, friends relationship repository boundary, friends relationship repository interface implementation, friends relationship PostgreSQL adapter gate, friends relationship runtime behavior gate, friends relationship protocol route gate, friends relationship protocol route implementation, friends relationship protocol route local proof, next Nakama prototype-ready capability after friends route proof, minimum operations inspection surface gate, minimum operations inspection source-first surface implementation, Pitaya-aligned distributed runtime vocabulary reactivation gate, Pitaya-aligned distributed runtime vocabulary source-first map, Pitaya-aligned frontend/backend role boundary gate, Pitaya-aligned frontend/backend role source-first map, Pitaya-aligned server-to-server RPC boundary gate, Pitaya-aligned server-to-server RPC source-first map, Pitaya-aligned service discovery boundary gate, Pitaya-aligned service discovery source-first map, Pitaya-aligned distributed group and broadcast boundary gate, Pitaya-aligned distributed group and broadcast source-first map, Pitaya-aligned cluster-safe session routing boundary gate, Pitaya-aligned cluster-safe session routing source-first map, and next Pitaya-aligned direction selection artifacts remain recorded in their existing documents. The friends route local proof is recorded in `ADR-0150` with check rule `runtime.friends_relationship_protocol_route_local_proof`; the post-proof selection is recorded in `ADR-0151` with check rule `runtime.next_nakama_prototype_ready_capability_after_friends_route_proof`; the operations gate is recorded in `ADR-0152` with check rule `runtime.minimum_operations_inspection_surface_gate`; the operations inspection implementation is recorded in `ADR-0153` with check rule `runtime.minimum_operations_inspection_source_first_surface_implementation`; the Pitaya vocabulary gate is recorded in `ADR-0154` with check rule `runtime.pitaya_aligned_distributed_runtime_vocabulary_reactivation_gate`; the source-first Pitaya vocabulary map is recorded in `ADR-0155` with check rule `runtime.pitaya_aligned_distributed_runtime_vocabulary_source_first_map`; the frontend/backend role boundary gate is recorded in `ADR-0156` with check rule `runtime.pitaya_aligned_frontend_backend_role_boundary_gate`; the frontend/backend role source-first map is recorded in `ADR-0157` with check rule `runtime.pitaya_aligned_frontend_backend_role_source_first_map` and `node tools/vibit inspect pitaya-roles --json`; the server-to-server RPC boundary gate is recorded in `ADR-0158` with check rule `runtime.pitaya_aligned_server_to_server_rpc_boundary_gate`; the server-to-server RPC source-first map is recorded in `ADR-0159` with check rule `runtime.pitaya_aligned_server_to_server_rpc_source_first_map` and `node tools/vibit inspect pitaya-rpc --json`; the service discovery boundary gate is recorded in `ADR-0160` with check rule `runtime.pitaya_aligned_service_discovery_boundary_gate`; the service discovery source-first map is recorded in `ADR-0161` with check rule `runtime.pitaya_aligned_service_discovery_source_first_map` and `node tools/vibit inspect pitaya-discovery --json`; the distributed group and broadcast boundary gate is recorded in `ADR-0162` with check rule `runtime.pitaya_aligned_distributed_group_broadcast_boundary_gate`; the distributed group and broadcast source-first map is recorded in `ADR-0163` with check rule `runtime.pitaya_aligned_distributed_group_broadcast_source_first_map` and `node tools/vibit inspect pitaya-groups --json`; the cluster-safe session routing boundary gate is recorded in `ADR-0164` with check rule `runtime.pitaya_aligned_cluster_safe_session_routing_boundary_gate`; the cluster-safe session routing source-first map is recorded in `ADR-0165` with check rule `runtime.pitaya_aligned_cluster_safe_session_routing_source_first_map` and `node tools/vibit inspect pitaya-sessions --json`; and the W-0258 direction selection is recorded in `ADR-0166` with check rule `runtime.next_pitaya_aligned_direction_after_cluster_safe_session_routing_map`. The next work is `W-0259 Define Pitaya-aligned route handler pipeline boundary gate`; next_direction: define_pitaya_aligned_route_handler_pipeline_boundary_gate.

Trace references for alpha acceptance include `docs/first-alpha-user-discovery-loop.md`, `docs/first-alpha-feedback-intake-surfaces.md`, `docs/product-maturity-milestones.md`, `docs/prototype-ready-foundation-execution-plan.md`, `docs/prototype-ready-local-development-path-gate.md`, `docs/storage-objects-behavior-gate.md`, and `docs/storage-objects-persistence-schema-gate.md`.

## 11. Release Deferrals

The following remain deferred until a later explicit work item:

- Creating additional release artifacts beyond the GitHub source archive.
- Creating release binaries, packages, containers, checksums, provenance files, signing artifacts, install scripts, registry publications, or hosted deployments.
- Adding a public local onboarding protocol route.
- Adding production signup, external identity providers, password login, account recovery, account merge, or multi-device linking.
- Adding broad operations/admin behavior, metrics backend integration, or production observability.
- Adding chat, friends, groups, parties, matchmaking, match runtime, SDKs, distributed runtime, or direct Nakama/Pitaya API compatibility.

## 12. Current Acceptance Result

The local alpha flow is checkable and packaged for local developer review, but the repository is still pre-alpha until a later release-publishing work item explicitly declares otherwise.

Current result:

```text
local_alpha_flow_checkable: true
local_alpha_developer_flow_packaged: true
release_publishing_decision_gate_defined: true
release_execution_preparation_gate_defined: true
release_execution_authorization_gate_defined: true
release_execution_maintainer_decision_recorded: true
release_identifier_artifact_plan_defined: true
release_execution_final_authorization_recorded: true
proposed_release_identifier: v0.1.0-alpha.1
authorized_release_identifier: v0.1.0-alpha.1
release_declared: true
release_publishing_authorized_by_this_checklist: true
prototype_ready_local_development_path_package_implemented: true
storage_objects_protocol_route_gate_defined: true
storage_objects_protocol_route_implementation_completed: true
storage_objects_protocol_route_local_proof_completed: true
first_server_push_realtime_messaging_gate_defined: true
first_server_push_realtime_messaging_runtime_slice_completed: true
realtime_protocol_websocket_outbound_delivery_gate_defined: true
realtime_protocol_websocket_outbound_delivery_implementation_completed: true
agent_native_feature_request_test_workflow_defined: true
operations_inspection_source_first_surface_implemented: true
pitaya_deferred_architecture_map_recorded: true
next_direction: implement_pitaya_aligned_metrics_tracing_source_first_map
next_work_status: next_ready
```

Current next work item: `W-0275 Implement Pitaya-aligned metrics and tracing source-first map`.

Current check rule: `runtime.pitaya_aligned_metrics_tracing_boundary_gate`.

Current inspection commands: `node tools/vibit inspect pitaya-session-lifecycle --json` and `node tools/vibit inspect pitaya-observability --json`.
