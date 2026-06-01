# Alpha Developer Flow

Status: Draft v0.1
Last updated: 2026-05-21
Scope: Packaged local developer journey for vibit's v0.1 alpha path

The paired Simplified Chinese translation is `docs/alpha-developer-flow.zh-CN.md`. The English file is authoritative.

This document packages the existing local alpha entry points into one developer journey. It is not a release declaration and does not authorize release publishing, release packaging, hosted deployment, runtime behavior changes, protocol changes, generated output changes, migrations, dependencies, broad operations/admin behavior, product module expansion, or direct Nakama/Pitaya API compatibility.

## 1. Purpose

The local alpha path now has the pieces a technically capable contributor needs to inspect vibit:

- project positioning in `README.md`,
- runtime startup and verification notes in `docs/runtime-runbook.md`,
- a redacted request-loop script at `examples/local-alpha-request-loop.sh`,
- local status endpoints at `/healthz`, `/readyz`, `/version`, and `/configz`,
- acceptance criteria in `docs/alpha-acceptance-checklist.md`,
- and continuation state in `.arch/work-items.yaml`.

This document connects those pieces so a contributor can follow the same sequence without hunting through project memory.

## 2. Current Package State

```text
local_alpha_developer_flow_packaged: true
release_declared: false
release_publishing_authorized_by_this_document: false
release_packaging_authorized_by_this_document: false
release_publishing_decision_gate_defined: true
release_execution_preparation_gate_defined: true
release_execution_authorization_gate_defined: true
release_execution_maintainer_decision_recorded: true
release_identifier_artifact_plan_defined: true
release_execution_final_authorization_recorded: true
proposed_release_identifier: v0.1.0-alpha.1
authorized_release_identifier: v0.1.0-alpha.1
prototype_ready_local_development_path_package_implemented: true
prototype_ready_local_development_path_package: docs/prototype-ready-local-development-path-package.md
storage_objects_protocol_route_gate_defined: true
storage_objects_protocol_route_gate: docs/storage-objects-protocol-route-gate.md
storage_objects_protocol_route_gate_decision: ADR-0118
storage_objects_protocol_route_implementation_completed: true
storage_objects_protocol_route_implementation_decision: ADR-0119
storage_objects_protocol_route_local_proof_completed: true
storage_objects_protocol_route_local_proof_decision: ADR-0120
first_server_push_realtime_messaging_gate_defined: true
first_server_push_realtime_messaging_gate_decision: ADR-0122
first_server_push_realtime_messaging_runtime_slice_completed: true
first_server_push_realtime_messaging_runtime_slice_decision: ADR-0123
next_alpha_direction_after_realtime_runtime_slice_selected: true
next_alpha_direction_after_realtime_runtime_slice_decision: ADR-0124
realtime_protocol_websocket_outbound_delivery_gate_defined: true
realtime_protocol_websocket_outbound_delivery_gate_decision: ADR-0125
realtime_protocol_websocket_outbound_delivery_implementation_completed: true
realtime_protocol_websocket_outbound_delivery_implementation_decision: ADR-0126
next_alpha_direction_after_realtime_outbound_delivery_slice_decision: ADR-0127
agent_native_feature_request_test_workflow_defined: true
agent_native_feature_request_test_workflow_decision: ADR-0128
agent_native_feature_request_test_workflow_standard: docs/agent-native-feature-request-test-workflow.md
primary_product_reference: Nakama
pitaya_reference_status: deferred_future_architecture_reference
ai_native_development_testing_goal: user_requirement_to_spec_tests_implementation_verification
minimum_operations_inspection_surface_gate_decision: ADR-0152
minimum_operations_inspection_surface_gate_check_rule: runtime.minimum_operations_inspection_surface_gate
minimum_operations_inspection_source_first_surface_decision: ADR-0153
minimum_operations_inspection_source_first_surface_check_rule: runtime.minimum_operations_inspection_source_first_surface_implementation
next_work_item: W-0256 Define Pitaya-aligned cluster-safe session routing boundary gate
next_direction: define_pitaya_aligned_cluster_safe_session_routing_boundary_gate
next_work_status: next_ready
```

The repository has a source-first alpha authorization for `v0.1.0-alpha.1`. The packaged flow is ready for local review, and the authorized release surface remains Git tag, GitHub release record, release notes, and GitHub source archive only. Binaries, packages, containers, checksums, signing/provenance artifacts, hosted deployment, install scripts, registry publication, and public announcements beyond the GitHub release record remain deferred.

The prototype-ready local development path package is recorded in `docs/prototype-ready-local-development-path-package.md`. Its example and local placeholder configuration entry points are `examples/README.md` and `examples/local.prototype.env.example`.

Trace references for the current packaged flow include `docs/first-alpha-user-discovery-loop.md`, `docs/first-alpha-feedback-intake-surfaces.md`, `docs/product-maturity-milestones.md`, `docs/prototype-ready-foundation-execution-plan.md`, `docs/prototype-ready-local-development-path-gate.md`, `docs/storage-objects-behavior-gate.md`, and `docs/storage-objects-persistence-schema-gate.md`.

## 3. Recommended Journey

1. Read `README.md` to understand vibit, its pre-alpha state, and the agent-native server framework goal.
2. Read `docs/v0.1-alpha-goal.md` for the short-term `v0.1 alpha` target and long-term Nakama/Pitaya-class direction.
3. Read `docs/alpha-acceptance-checklist.md` to see which alpha items are ready, manual, deferred, or blocked.
4. Install local prerequisites: Go, Node.js, PostgreSQL when testing the persistent runtime path, and Buf/Protobuf tooling only when regenerating Protobuf output.
5. Run static repository checks:

```bash
node tools/vibit inspect next
node tools/vibit check work --json
node tools/vibit check all --json
```

6. Run Go tests:

```bash
cd runtime
go test ./...
```

7. Run the redacted local alpha request loop from the repository root:

```bash
examples/local-alpha-request-loop.sh
```

8. Use `docs/runtime-runbook.md` when starting the server manually or evaluating the PostgreSQL runtime path.
9. Read `docs/prototype-ready-local-development-path-package.md` when evaluating the repeatable source-first local path.
10. Use `.arch/work-items.yaml` and `node tools/vibit inspect next` to find the next bounded contribution.

## 4. Runtime Entry Points

The current process exposes:

```text
/v1/ws
/healthz
/readyz
/version
/configz
```

`/v1/ws` is the gameplay WebSocket endpoint. It expects binary `vibit.protocol.v1.Envelope` Protobuf bytes. JSON is not accepted on this endpoint.

`/healthz`, `/readyz`, `/version`, and `/configz` are local alpha troubleshooting endpoints. `/configz` reports redacted runtime posture and includes `secrets_redacted: true`. These endpoints are not a production operations API, admin console, metrics backend, gameplay protocol route, release artifact, or hosted deployment surface.

## 5. Local Proof Flow

The packaged local proof flow is:

```text
local onboarding
-> device credential login
-> connection binding
-> protected inventory grant/read
-> protected presence query
-> protected own-player storage object put/get/list/delete
-> logout
-> post-logout protected request rejection
```

The executable entry point is:

```bash
examples/local-alpha-request-loop.sh
```

The script wraps the focused Go E2E proof:

```bash
cd runtime && go test ./internal/platform/protocol/protobuf -run 'TestAuthenticatedGameplayE2EUsesExistingOnboardingLoginBindingInventoryPresenceAndLogout|TestStorageObjectsProtocolRouteLocalAlphaFlow' -v
```

The proof uses existing runtime protocol handlers. It does not require live PostgreSQL, committed verifier keys, raw credentials, raw access tokens, DSNs, digests, or a hand-built WebSocket client.

## 6. PostgreSQL Path

The PostgreSQL runtime path is the current alpha runtime composition, but it has manual setup requirements:

- prepare a local PostgreSQL database,
- apply or verify SQL migrations explicitly,
- set `VIBIT_RUNTIME_STORE=postgres`,
- set `VIBIT_POSTGRES_DSN`,
- set all authentication verifier key environment variables,
- avoid committing local verifier keys or DSNs.

Normal server startup does not apply migrations automatically. Optional live PostgreSQL verification remains opt-in through `VIBIT_POSTGRES_TEST_DSN` and a disposable database.

## 7. Redaction Contract

Do not record or commit:

- raw device credential text or bytes,
- raw access tokens,
- credential or token lookup digests,
- credential or token verifier digests,
- HMAC input or output bytes,
- verifier key values,
- concrete verifier key set ids,
- PostgreSQL DSNs with credentials,
- headers, cookies, query strings, WebSocket subprotocol values, or remote addresses that may carry secrets.

The request-loop script and `/configz` surface are part of this redaction posture.

## 8. Contribution Entry Point

The next contribution path is always machine-readable:

```bash
node tools/vibit inspect next
```

The release execution maintainer decision, release identifier plan, final authorization, first alpha discovery and feedback surfaces, prototype-ready execution plan, storage objects path, realtime outbound delivery path, Nakama-first direction, agent-native feature request/test workflow, Nakama-aligned presence/status workflow pilot, presence/status local proof hardening, authenticated gameplay failure-path verification, next Nakama prototype-ready capability selection, local alpha example client path gate, local alpha example client path implementation, feature request scaffolding selection, feature request scaffolding gate, feature request scaffolding implementation, scaffolded Nakama feature request intake pilot, friends relationship lifecycle gate, friends relationship persistence schema gate, friends relationship migration source, friends relationship repository boundary, friends relationship repository interface implementation, friends relationship PostgreSQL adapter gate, friends relationship PostgreSQL adapter implementation, friends relationship runtime behavior gate, friends relationship runtime behavior implementation, friends relationship protocol route gate, friends relationship protocol route implementation, friends relationship protocol route local proof, minimum operations inspection surface gate, minimum operations inspection source-first surface implementation, Pitaya-aligned distributed runtime vocabulary reactivation gate, Pitaya-aligned distributed runtime vocabulary source-first map, Pitaya-aligned frontend/backend role boundary gate, Pitaya-aligned frontend/backend role source-first map, Pitaya-aligned server-to-server RPC boundary gate, Pitaya-aligned server-to-server RPC source-first map, Pitaya-aligned service discovery boundary gate, Pitaya-aligned service discovery source-first map, Pitaya-aligned distributed group and broadcast boundary gate, and Pitaya-aligned distributed group and broadcast source-first map remain recorded in their existing artifacts. The agent-native feature request and test workflow is recorded in `docs/agent-native-feature-request-test-workflow.md` and `ADR-0128`; the friends protocol route local proof decision is recorded in `ADR-0150` with check rule `runtime.friends_relationship_protocol_route_local_proof`; the minimum operations inspection surface gate is recorded in `ADR-0152` with check rule `runtime.minimum_operations_inspection_surface_gate`; the source-first operations inspection implementation is recorded in `ADR-0153` with check rule `runtime.minimum_operations_inspection_source_first_surface_implementation`; the Pitaya vocabulary gate is recorded in `ADR-0154` with check rule `runtime.pitaya_aligned_distributed_runtime_vocabulary_reactivation_gate`; the source-first Pitaya vocabulary map is recorded in `ADR-0155` with check rule `runtime.pitaya_aligned_distributed_runtime_vocabulary_source_first_map`; the frontend/backend role boundary gate is recorded in `ADR-0156` with check rule `runtime.pitaya_aligned_frontend_backend_role_boundary_gate`; the frontend/backend role source-first map is recorded in `ADR-0157` with check rule `runtime.pitaya_aligned_frontend_backend_role_source_first_map` and `node tools/vibit inspect pitaya-roles --json`; the server-to-server RPC boundary gate is recorded in `ADR-0158` with check rule `runtime.pitaya_aligned_server_to_server_rpc_boundary_gate`; the server-to-server RPC source-first map is recorded in `ADR-0159` with check rule `runtime.pitaya_aligned_server_to_server_rpc_source_first_map` and `node tools/vibit inspect pitaya-rpc --json`; the service discovery boundary gate is recorded in `ADR-0160` with check rule `runtime.pitaya_aligned_service_discovery_boundary_gate`; the service discovery source-first map is recorded in `ADR-0161` with check rule `runtime.pitaya_aligned_service_discovery_source_first_map` and `node tools/vibit inspect pitaya-discovery --json`; the distributed group and broadcast boundary gate is recorded in `ADR-0162` with check rule `runtime.pitaya_aligned_distributed_group_broadcast_boundary_gate`; and the distributed group and broadcast source-first map is recorded in `ADR-0163` with check rule `runtime.pitaya_aligned_distributed_group_broadcast_source_first_map` and `node tools/vibit inspect pitaya-groups --json`. The next work is `W-0256 Define Pitaya-aligned cluster-safe session routing boundary gate`; next_direction: define_pitaya_aligned_cluster_safe_session_routing_boundary_gate.

## 9. Deferred Work

The following remain deferred until later explicit work items:

- publishing `v0.1 alpha`,
- selecting release identifiers for execution,
- creating release tags, binaries, archives, containers, packages, checksums, provenance files, or hosted deployments,
- adding a public local onboarding protocol route,
- adding production signup, external identity providers, password login, account recovery, account merge, or multi-device linking,
- adding broad operations/admin behavior, metrics backend integration, or production observability,
- adding chat, friends, groups, parties, matchmaking, match runtime, SDKs, distributed runtime, or direct Nakama/Pitaya API compatibility.

## 10. Verification

Use this command set when checking the packaged flow:

```bash
node -c tools/vibit
node tools/vibit inspect next
node tools/vibit check work --json
node tools/vibit check runtime --json
node tools/vibit check memory --json
node tools/vibit check schemas --json
node tools/vibit check all --json
cd runtime && go test ./cmd/vibit-server
cd runtime && go test ./...
examples/local-alpha-request-loop.sh
git diff --check
```
