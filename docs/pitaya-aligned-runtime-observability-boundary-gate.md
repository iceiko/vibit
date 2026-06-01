# Pitaya-Aligned Runtime Observability Boundary Gate

Status: Accepted v0.1
Last updated: 2026-06-01
Scope: Gate-only boundary for using Pitaya-aligned runtime observability vocabulary after the session binding, kick/disconnect, and session data follow-up direction selection
Depends on: `decisions/ADR-0178-select-next-pitaya-aligned-direction-after-session-binding-kick-disconnect-session-data-map.md`, `decisions/ADR-0177-pitaya-aligned-session-binding-kick-disconnect-session-data-source-first-map.md`, `docs/minimum-operations-inspection-surface-gate.md`, `docs/runtime-runbook.md`, `docs/reference-game-server-alignment.md`, `.arch/reference.yaml`
Canonical decision: `ADR-0179`

The paired Simplified Chinese translation is `docs/pitaya-aligned-runtime-observability-boundary-gate.zh-CN.md`. The English file is authoritative.

This document defines a runtime observability vocabulary gate only. It does not implement runtime endpoint behavior, metrics endpoints, tracing pipelines, observability pipelines, dashboards, admin console behavior, player/session/token inspectors, event/audit tables, transport behavior changes, protocol messages or routes, Protobuf source, generated output, repository interfaces, PostgreSQL adapters, migrations, dependencies, hosted deployment, SDK publication, release artifacts, distributed runtime behavior, or direct Nakama/Pitaya API compatibility.

## 1. Core Rule

The Pitaya-aligned runtime observability boundary gate record is:

```yaml
pitaya_aligned_runtime_observability_boundary_gate: defined
completed_work_item: W-0271
decision: ADR-0179
check_rule: runtime.pitaya_aligned_runtime_observability_boundary_gate
selection_decision: ADR-0178
session_lifecycle_source_first_map_decision: ADR-0177
minimum_operations_inspection_gate_decision: ADR-0152
minimum_operations_inspection_source_first_decision: ADR-0153
standard: docs/pitaya-aligned-runtime-observability-boundary-gate.md
translation: docs/pitaya-aligned-runtime-observability-boundary-gate.zh-CN.md
primary_product_reference: Nakama
pitaya_reference_status: runtime_observability_vocabulary_boundary_defined_for_future_architecture_planning
implementation_scope: gate_only_runtime_observability_vocabulary
future_implementation_work_item: W-0272
future_implementation_direction: implement_pitaya_aligned_runtime_observability_source_first_map
allowed_runtime_observability_vocabulary:
  - runtime_observability_boundary
  - operations_snapshot
  - health_readiness_signal
  - version_release_posture
  - configuration_posture
  - route_inventory_snapshot
  - verification_posture
  - redaction_posture
  - deferred_operations_surface
  - node_local_runtime_surface
current_source_first_runtime_observability_mapping:
  minimum_operations_inspection:
    current: node_tools_vibit_inspect_operations
    future_vocabulary: operations_snapshot
    implementation_status: source_first_repository_inspection_only
  health_and_readiness:
    current: existing_local_healthz_and_readyz_endpoint_summary
    future_vocabulary: health_readiness_signal
    implementation_status: no_new_runtime_endpoint_behavior
  version_surface:
    current: existing_local_version_endpoint_summary
    future_vocabulary: version_release_posture
    implementation_status: no_release_or_hosted_surface_change
  configuration_surface:
    current: existing_redacted_configz_endpoint_summary
    future_vocabulary: configuration_posture
    implementation_status: no_secret_or_dsn_exposure
  route_inventory:
    current: source_first_protocol_route_family_inventory
    future_vocabulary: route_inventory_snapshot
    implementation_status: no_protocol_route_or_wire_shape_change
  repository_verification:
    current: tools_vibit_check_work_runtime_memory_schemas_all
    future_vocabulary: verification_posture
    implementation_status: no_new_test_framework_or_dependency
  redaction_policy:
    current: existing_operations_redaction_requirements
    future_vocabulary: redaction_posture
    implementation_status: no_identifier_or_payload_log_safety_expansion
  deferred_operations_surfaces:
    current: admin_metrics_tracing_dashboard_inspectors_deferred
    future_vocabulary: deferred_operations_surface
    implementation_status: planning_vocabulary_only
runtime_endpoint_behavior_added: false
metrics_endpoint_added: false
tracing_pipeline_added: false
observability_pipeline_added: false
dashboard_added: false
admin_console_added: false
player_session_token_inspector_added: false
event_audit_table_added: false
session_binding_behavior_added: false
kick_disconnect_behavior_added: false
session_data_behavior_added: false
session_data_persistence_added: false
acceptor_behavior_added: false
tcp_acceptor_added: false
websocket_behavior_changed: false
connection_lifecycle_behavior_changed: false
route_handler_implementation_added: false
handler_routing_behavior_added: false
handler_pipeline_behavior_added: false
pipeline_middleware_behavior_added: false
backend_route_targeting_added: false
cluster_safe_session_routing_added: false
session_location_registry_added: false
connection_owner_node_registry_added: false
routing_epoch_behavior_added: false
remote_connection_handoff_added: false
distributed_session_routing_added: false
distributed_group_implementation_added: false
room_broadcast_fanout_added: false
broadcast_delivery_guarantee_added: false
stream_subscription_added: false
service_discovery_implementation_added: false
service_registry_added: false
service_selector_added: false
node_registry_added: false
server_identity_added: false
server_to_server_rpc_implementation_added: false
remote_call_behavior_added: false
frontend_server_implementation_added: false
backend_server_implementation_added: false
frontend_backend_server_roles_added: false
distributed_runtime_implementation_added: false
runtime_behavior_added: false
protocol_route_added: false
protobuf_source_added: false
generated_output_added: false
repository_interface_changed: false
postgresql_adapter_changed: false
migration_added: false
dependency_added: false
authentication_session_behavior_changed: false
hosted_deployment_added: false
sdk_added: false
release_artifact_added: false
direct_nakama_pitaya_api_compatibility_added: false
```

## 2. Purpose

`ADR-0178` selected runtime observability as the next Pitaya-aligned direction after the session lifecycle source-first map. The risk is that agents may treat observability as permission to add metrics endpoints, tracing pipelines, dashboard behavior, admin operations, or live state inspectors.

This gate records vocabulary and mapping only. It connects the existing source-first operations inspection surface to future Pitaya-aligned runtime observability planning while keeping the current local alpha runtime behavior unchanged.

## 3. Vocabulary

Allowed runtime observability vocabulary:

- `runtime_observability_boundary`: future planning vocabulary for the boundary that owns runtime inspection and operational visibility.
- `operations_snapshot`: future planning vocabulary for a source-first snapshot of runtime posture, route families, verification, and deferred operations state.
- `health_readiness_signal`: future planning vocabulary for liveness and readiness status. Current behavior remains `/healthz` and `/readyz`.
- `version_release_posture`: future planning vocabulary for source and runtime version posture. Current behavior remains `/version`.
- `configuration_posture`: future planning vocabulary for redacted configuration state. Current behavior remains `/configz`.
- `route_inventory_snapshot`: future planning vocabulary for committed route family inventory. It is not a protocol or dispatch change.
- `verification_posture`: future planning vocabulary for repository checks and known warnings.
- `redaction_posture`: future planning vocabulary for what observability surfaces may expose.
- `deferred_operations_surface`: future planning vocabulary for admin, metrics, tracing, dashboard, and inspector surfaces that are still deferred.
- `node_local_runtime_surface`: future planning vocabulary for local single-process operational facts. It is not distributed node telemetry.

Forbidden vocabulary use:

- Do not introduce concrete public API, package, route, method, wire, handler, metrics, trace, dashboard, or admin compatibility names from Pitaya or Nakama.
- Do not use runtime observability vocabulary as permission to add runtime endpoints, metrics endpoints, tracing pipelines, observability pipelines, dashboards, admin console behavior, live player/session/token inspectors, event/audit tables, dependencies, protocol messages, generated output, persistence, hosted surfaces, SDKs, or distributed runtime behavior.
- Do not classify raw tokens, credentials, lookup digests, verifier digests, verifier keys, DSNs, headers, cookies, query strings, subprotocol values, remote addresses, database payloads, local secret file contents, route payloads, session data payloads, or concrete transport metadata as log-safe in this gate.

## 4. Current Mapping

```yaml
current_source_first_runtime_observability_mapping:
  minimum_operations_inspection:
    current: node tools/vibit inspect operations --json
    future_vocabulary: operations_snapshot
    status: source_first_repository_inspection_only
  health_and_readiness:
    current: existing /healthz and /readyz local troubleshooting endpoints
    future_vocabulary: health_readiness_signal
    status: no_new_runtime_endpoint_behavior
  version_surface:
    current: existing /version local troubleshooting endpoint
    future_vocabulary: version_release_posture
    status: no_release_or_hosted_surface_change
  configuration_surface:
    current: existing redacted /configz local troubleshooting endpoint
    future_vocabulary: configuration_posture
    status: no_secret_or_dsn_exposure
  route_inventory:
    current: source-first route family inventory from committed source
    future_vocabulary: route_inventory_snapshot
    status: no_protocol_route_or_wire_shape_change
  repository_verification:
    current: tools/vibit repository checks and known warning posture
    future_vocabulary: verification_posture
    status: no_new_test_framework_or_dependency
  redaction_policy:
    current: minimum operations inspection redaction requirements
    future_vocabulary: redaction_posture
    status: no_identifier_or_payload_log_safety_expansion
  deferred_operations_surfaces:
    current: admin, metrics, tracing, dashboard, inspectors, hosted, SDK, and distributed operations deferred
    future_vocabulary: deferred_operations_surface
    status: planning_vocabulary_only
```

## 5. Ownership

Runtime observability vocabulary ownership:

```yaml
architecture_vocabulary_owner:
  - docs/pitaya-aligned-runtime-observability-boundary-gate.md
  - .arch/reference.yaml
  - .arch/runtime.yaml
source_first_map_candidate_owner:
  - tools/vibit
runtime_endpoint_owner: unchanged
metrics_owner: deferred
tracing_owner: deferred
dashboard_owner: deferred
admin_console_owner: deferred
protocol_owner: unchanged
persistence_owner: unchanged
dependency_owner: unchanged
```

Rules:

- Documentation and manifests may define runtime observability vocabulary and current source-first mapping.
- `tools/vibit` may later emit a source-first runtime observability map if a follow-up implementation work item authorizes it.
- Existing `/healthz`, `/readyz`, `/version`, and `/configz` endpoints remain the only runtime troubleshooting endpoints accepted by this gate.
- Runtime behavior, transport behavior, protocol payloads, repository interfaces, migrations, generated output, startup wiring, metrics, tracing, dashboards, admin console behavior, event/audit tables, dependencies, hosted surfaces, SDKs, and distributed runtime behavior remain unchanged by this gate.

## 6. Nakama And Pitaya Mapping

Nakama remains the primary product reference for admin console, metrics, observability, and operations capability pressure. Pitaya remains an architecture vocabulary reference for runtime, node, session, route, RPC, service discovery, and group concerns.

This gate maps those references into vibit-owned vocabulary only. It does not create direct compatibility, public API parity, metrics parity, tracing parity, dashboard parity, or runtime behavior.

## 7. Stop Conditions

Stop and require a later bounded work item before adding:

- Runtime endpoint behavior, metrics endpoints, tracing pipelines, observability pipelines, dashboards, admin console behavior, player/session/token inspectors, event/audit tables, logging policy changes, dependency changes, hosted operations surfaces, SDK publication, release artifacts, or direct Nakama/Pitaya API compatibility.
- Protocol messages or routes, Protobuf sources, generated output, repository interfaces, PostgreSQL adapters, migrations, database inspection, arbitrary state dumps, or persistence changes.
- Session binding behavior, kick/disconnect behavior, session data behavior or persistence, acceptor behavior, WebSocket behavior changes, connection lifecycle behavior changes, route handler implementation, backend route targeting, cluster-safe session routing, service discovery, RPC, frontend/backend roles, or distributed runtime behavior.

## 8. Verification

The repository check rule for this boundary is:

```text
runtime.pitaya_aligned_runtime_observability_boundary_gate
```

Verification commands:

```sh
node -c tools/vibit
node tools/vibit inspect next --json
node tools/vibit inspect rule runtime.pitaya_aligned_runtime_observability_boundary_gate
node tools/vibit check change define-pitaya-aligned-runtime-observability-boundary-gate --json
node tools/vibit check work --json
node tools/vibit check runtime --json
node tools/vibit check memory --json
node tools/vibit check schemas --json
node tools/vibit check all --json
git diff --check
```

The gate itself does not require Go tests because it adds no Go runtime behavior.
