# Pitaya-Aligned Dashboard And Admin Operations Boundary Gate

Status: Accepted v0.1
Last updated: 2026-06-01
Scope: Gate-only boundary for using Pitaya-aligned dashboard and admin operations vocabulary after the metrics and tracing source-first map
Depends on: `decisions/ADR-0184-select-next-pitaya-aligned-direction-after-metrics-tracing-map.md`, `decisions/ADR-0183-pitaya-aligned-metrics-tracing-source-first-map.md`, `docs/pitaya-aligned-metrics-tracing-boundary-gate.md`, `docs/reference-game-server-alignment.md`, `.arch/reference.yaml`
Canonical decision: `ADR-0185`

The paired Simplified Chinese translation is `docs/pitaya-aligned-dashboard-admin-operations-boundary-gate.zh-CN.md`. The English file is authoritative.

This document defines a dashboard and admin operations vocabulary gate only. It does not implement runtime endpoint behavior, metrics endpoints, tracing pipelines, observability pipelines, dashboards, admin console behavior, player/session/token inspectors, event/audit tables, transport behavior changes, protocol messages or routes, Protobuf source, generated output, repository interfaces, PostgreSQL adapters, migrations, dependencies, hosted deployment, SDK publication, release artifacts, distributed runtime behavior, or direct Nakama/Pitaya API compatibility.

## 1. Core Rule

The Pitaya-aligned dashboard and admin operations boundary gate record is:

```yaml
pitaya_aligned_dashboard_admin_operations_boundary_gate: defined
completed_work_item: W-0277
decision: ADR-0185
check_rule: runtime.pitaya_aligned_dashboard_admin_operations_boundary_gate
selection_decision: ADR-0184
metrics_tracing_source_first_map_decision: ADR-0183
metrics_tracing_boundary_gate_decision: ADR-0182
standard: docs/pitaya-aligned-dashboard-admin-operations-boundary-gate.md
translation: docs/pitaya-aligned-dashboard-admin-operations-boundary-gate.zh-CN.md
primary_product_reference: Nakama
pitaya_reference_status: dashboard_admin_operations_vocabulary_boundary_defined_for_future_architecture_planning
implementation_scope: gate_only_dashboard_admin_operations_vocabulary
future_implementation_work_item: W-0278
future_implementation_direction: implement_pitaya_aligned_dashboard_admin_operations_source_first_map
allowed_dashboard_admin_operations_vocabulary:
  - dashboard_admin_boundary
  - admin_operation_surface
  - operator_action_boundary
  - dashboard_read_model
  - source_first_operations_snapshot
  - inspector_redaction_posture
  - event_audit_deferral
  - local_operations_diagnostic_surface
  - future_admin_authorization_boundary
  - hosted_operations_deferral
current_source_first_dashboard_admin_operations_mapping:
  operations_snapshot:
    current: node_tools_vibit_inspect_operations
    future_vocabulary: source_first_operations_snapshot
    implementation_status: source_first_repository_inspection_only
  runtime_observability_map:
    current: node_tools_vibit_inspect_pitaya_observability
    future_vocabulary: dashboard_read_model
    implementation_status: no_dashboard_or_admin_console_behavior
  metrics_tracing_map:
    current: node_tools_vibit_inspect_pitaya_metrics_tracing
    future_vocabulary: dashboard_admin_boundary
    implementation_status: no_metrics_endpoint_or_tracing_pipeline_behavior
  health_readiness_version_config:
    current: existing_local_health_readiness_version_config_surfaces
    future_vocabulary: local_operations_diagnostic_surface
    implementation_status: no_new_runtime_endpoint_behavior
  route_inventory_snapshot:
    current: source_first_protocol_route_family_inventory
    future_vocabulary: dashboard_read_model
    implementation_status: no_protocol_route_or_wire_shape_change
  verification_posture:
    current: tools_vibit_check_work_runtime_memory_schemas_all
    future_vocabulary: admin_operation_surface
    implementation_status: no_admin_operation_execution_behavior
  inspector_redaction_policy:
    current: operations_and_metrics_tracing_redaction_requirements
    future_vocabulary: inspector_redaction_posture
    implementation_status: no_player_session_token_inspector_behavior
  audit_event_posture:
    current: no_event_audit_tables_or_admin_audit_runtime_behavior
    future_vocabulary: event_audit_deferral
    implementation_status: planning_vocabulary_only
  admin_authorization_posture:
    current: no_admin_user_role_permission_or_authentication_behavior
    future_vocabulary: future_admin_authorization_boundary
    implementation_status: planning_vocabulary_only
runtime_endpoint_behavior_added: false
metrics_endpoint_added: false
tracing_pipeline_added: false
observability_pipeline_added: false
dashboard_added: false
admin_console_added: false
player_session_token_inspector_added: false
event_audit_table_added: false
protocol_route_added: false
protobuf_source_added: false
generated_output_added: false
repository_interface_changed: false
postgresql_adapter_changed: false
migration_added: false
dependency_added: false
hosted_deployment_added: false
sdk_added: false
release_artifact_added: false
distributed_runtime_implementation_added: false
direct_nakama_pitaya_api_compatibility_added: false
```

## 2. Purpose

`ADR-0184` selected dashboard and admin operations as the next Pitaya-aligned direction after the metrics and tracing source-first map.

The risk is that agents may treat dashboard/admin vocabulary as permission to add UI, live operations endpoints, admin authentication, sensitive inspectors, audit tables, deployment surfaces, or third-party dependencies. This gate records vocabulary and source-first mapping only. It keeps the current local alpha runtime behavior unchanged and prepares a narrow source-first map follow-up.

## 3. Vocabulary

Allowed dashboard and admin operations vocabulary:

- `dashboard_admin_boundary`: future planning vocabulary for the boundary that owns dashboard/admin operations semantics.
- `admin_operation_surface`: future planning vocabulary for an operator-facing command, query, or action surface.
- `operator_action_boundary`: future planning vocabulary for actions that would later require authorization and audit posture.
- `dashboard_read_model`: future planning vocabulary for source-derived read models that may later inform dashboards.
- `source_first_operations_snapshot`: future planning vocabulary for repository-derived operations facts.
- `inspector_redaction_posture`: future planning vocabulary for what sensitive inspection surfaces may expose.
- `event_audit_deferral`: future planning vocabulary for audit and event table work that is still deferred.
- `local_operations_diagnostic_surface`: future planning vocabulary for local diagnostic facts already represented in source-first surfaces.
- `future_admin_authorization_boundary`: future planning vocabulary for admin authorization that is not selected here.
- `hosted_operations_deferral`: future planning vocabulary for hosted operations work that remains deferred.

Forbidden vocabulary use:

- Do not introduce concrete public API, package, route, method, wire, handler, dashboard, metrics, trace, admin, console, or inspector compatibility names from Pitaya or Nakama.
- Do not use dashboard/admin vocabulary as permission to add runtime endpoints, metrics endpoints, tracing pipelines, observability pipelines, dashboards, admin console behavior, live player/session/token inspectors, event/audit tables, dependencies, protocol messages, generated output, persistence, hosted surfaces, SDKs, release artifacts, or distributed runtime behavior.
- Do not classify raw tokens, credentials, lookup digests, verifier digests, verifier keys, DSNs, headers, cookies, query strings, subprotocol values, remote addresses, database payloads, local secret file contents, route payloads, session data payloads, dashboard payloads, admin console payloads, or concrete transport metadata as log-safe in this gate.

## 4. Current Mapping

```yaml
current_source_first_dashboard_admin_operations_mapping:
  operations_snapshot:
    current: node tools/vibit inspect operations --json
    future_vocabulary: source_first_operations_snapshot
    status: source_first_repository_inspection_only
  runtime_observability_map:
    current: node tools/vibit inspect pitaya-observability --json
    future_vocabulary: dashboard_read_model
    status: no_dashboard_or_admin_console_behavior
  metrics_tracing_map:
    current: node tools/vibit inspect pitaya-metrics-tracing --json
    future_vocabulary: dashboard_admin_boundary
    status: no_metrics_endpoint_or_tracing_pipeline_behavior
  verification_posture:
    current: tools/vibit repository checks and known warning posture
    future_vocabulary: admin_operation_surface
    status: no_admin_operation_execution_behavior
  redaction_policy:
    current: operations and metrics/tracing redaction requirements
    future_vocabulary: inspector_redaction_posture
    status: no_sensitive_inspector_behavior
  audit_event_posture:
    current: no event/audit tables or admin audit runtime behavior
    future_vocabulary: event_audit_deferral
    status: planning_vocabulary_only
```

## 5. Ownership

Dashboard and admin operations vocabulary ownership:

```yaml
architecture_vocabulary_owner:
  - docs/pitaya-aligned-dashboard-admin-operations-boundary-gate.md
  - .arch/reference.yaml
  - .arch/runtime.yaml
source_first_map_candidate_owner:
  - tools/vibit
runtime_endpoint_owner: unchanged
metrics_endpoint_owner: deferred
tracing_pipeline_owner: deferred
dashboard_owner: deferred
admin_console_owner: deferred
inspector_owner: deferred
event_audit_owner: deferred
protocol_owner: unchanged
persistence_owner: unchanged
dependency_owner: unchanged
```

Rules:

- Documentation and manifests may define dashboard/admin operations vocabulary and current source-first mapping.
- `tools/vibit` may later emit a source-first dashboard/admin operations map if a follow-up implementation work item authorizes it.
- Existing `/healthz`, `/readyz`, `/version`, and `/configz` endpoints remain unchanged by this gate.
- Runtime behavior, transport behavior, protocol payloads, repository interfaces, migrations, generated output, startup wiring, metrics endpoints, tracing pipelines, dashboards, admin console behavior, sensitive inspectors, event/audit tables, dependencies, hosted surfaces, SDKs, release artifacts, and distributed runtime behavior remain unchanged by this gate.

## 6. Nakama And Pitaya Mapping

Nakama remains the primary product reference for admin console, metrics, observability, and operations capability pressure. Pitaya remains an architecture vocabulary reference for runtime, node, session, route, RPC, service discovery, group, and operational concerns.

This gate maps those references into vibit-owned vocabulary only. It does not create direct compatibility, public API parity, dashboard parity, admin console parity, metrics parity, tracing parity, or runtime behavior.

## 7. Stop Conditions

Stop and require a later bounded work item before adding:

- runtime endpoint behavior;
- metrics endpoints;
- tracing pipelines;
- observability pipelines;
- dashboards;
- admin console behavior;
- player/session/token inspectors;
- event/audit tables;
- admin users, roles, permissions, or authentication behavior;
- protocol messages or routes;
- Protobuf source;
- generated output;
- repository interfaces;
- PostgreSQL adapters;
- migrations;
- dependencies;
- hosted deployment;
- SDK publication;
- release artifacts;
- distributed runtime behavior;
- direct Nakama/Pitaya API compatibility.

## 8. Verification

The repository check rule is `runtime.pitaya_aligned_dashboard_admin_operations_boundary_gate`.

The check verifies the standard, translation, ADR, change artifacts, manifest references, next-ready state, vocabulary markers, and explicit implementation deferrals.
