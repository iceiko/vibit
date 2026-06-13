# Pitaya-Aligned Service Export Boundary Gate 中文版

Status: Accepted v0.1
Last updated: 2026-06-13
Scope: Gate-only boundary for using Pitaya-aligned service export vocabulary after the currency wallet protocol route implementation
Depends on: `decisions/ADR-0210-currency-wallet-protocol-route-implementation.md`, `docs/reference-game-server-alignment.md`, `.arch/reference.yaml`
Canonical decision: `ADR-0211`

说明：本文件是 `docs/pitaya-aligned-service-export-boundary-gate.md` 的简体中文译本。英文版本是权威版本。

This document defines service export vocabulary only. It does not implement service export behavior, component service registration, backend service export behavior, runtime behavior, protocol messages or routes, Protobuf source, generated output, repository interfaces, PostgreSQL adapters, migrations, dependencies, hosted deployment, SDK publication, release artifacts, distributed runtime behavior, or direct Nakama/Pitaya API compatibility.

## 1. Core Rule

```yaml
pitaya_aligned_service_export_boundary_gate: defined
completed_work_item: W-0303
decision: ADR-0211
check_rule: runtime.pitaya_aligned_service_export_boundary_gate
source_currency_wallet_protocol_route_implementation_decision: ADR-0210
standard: docs/pitaya-aligned-service-export-boundary-gate.md
translation: docs/pitaya-aligned-service-export-boundary-gate.zh-CN.md
primary_product_reference: Nakama
pitaya_reference_status: service_export_vocabulary_boundary_defined_for_future_architecture_planning
implementation_scope: gate_only_service_export_vocabulary
future_implementation_work_item: W-0304
future_implementation_direction: implement_pitaya_aligned_service_export_source_first_map
allowed_service_export_vocabulary:
  - service_export_boundary
  - explicit_service_catalog_source
  - handler_export_deferral
  - component_export_handoff
  - backend_export_deferral
  - distributed_service_export_deferral
runtime_behavior_added: false
service_export_behavior_added: false
remote_call_dispatch_behavior_added: false
server_to_server_rpc_behavior_added: false
frontend_message_forwarding_behavior_added: false
backend_service_route_ownership_behavior_added: false
backend_service_targeting_behavior_added: false
cluster_event_bus_behavior_added: false
cluster_pubsub_behavior_added: false
node_event_fanout_behavior_added: false
outbox_event_delivery_behavior_added: false
service_discovery_behavior_added: false
frontend_backend_role_runtime_behavior_added: false
distributed_runtime_implementation_added: false
runtime_endpoint_behavior_added: false
protocol_route_added: false
protobuf_source_added: false
generated_output_added: false
repository_interface_changed: false
postgresql_adapter_changed: false
migration_added: false
dependency_added: false
persistence_added: false
startup_wiring_added: false
authentication_session_behavior_changed: false
hosted_deployment_added: false
sdk_added: false
release_artifact_added: false
direct_nakama_pitaya_api_compatibility_added: false
```

## 2. Purpose

The purpose is to keep service export visible as a future Pitaya-class architecture concern without authorizing implementation. Agents may use this vocabulary for planning, inspection, ADRs, change specs, and future source-first maps only.

## 3. Vocabulary

- `service_export_boundary`: planning vocabulary for service export without implementing service export behavior.
- `explicit_service_catalog_source`: planning vocabulary for service export without implementing service export behavior.
- `handler_export_deferral`: planning vocabulary for service export without implementing service export behavior.
- `component_export_handoff`: planning vocabulary for service export without implementing service export behavior.
- `backend_export_deferral`: planning vocabulary for service export without implementing service export behavior.
- `distributed_service_export_deferral`: planning vocabulary for service export without implementing service export behavior.

Forbidden use:

- Do not introduce concrete public API, package, route, method, wire, handler, dashboard, metrics, trace, admin, console, or inspector compatibility names from Pitaya or Nakama.
- Do not use this vocabulary as permission to add service export behavior, component service registration, backend service export behavior, runtime behavior, protocol messages, generated output, persistence, dependencies, hosted surfaces, SDKs, release artifacts, or distributed runtime behavior.
- Do not classify raw tokens, credentials, lookup digests, verifier digests, verifier keys, DSNs, headers, cookies, query strings, subprotocol values, remote addresses, database payloads, local secret file contents, route payloads, event payloads, or concrete transport metadata as log-safe in this gate.

## 4. Current Mapping Candidate

```yaml
current_source_first_service_export_mapping:
  application_services:
    current: runtime/internal/app service packages
    future_vocabulary: explicit_service_catalog_source
    status: source_first_repository_inspection_only
  route_handlers:
    current: runtime/internal/app/bootstrap route handler sources
    future_vocabulary: handler_export_deferral
    status: no_service_export_behavior
  module_services:
    current: runtime/internal/modules storage-neutral domain interfaces
    future_vocabulary: component_export_handoff
    status: no_component_service_registration
  backend_exports:
    current: deferred
    future_vocabulary: backend_export_deferral
    status: no_backend_service_export_behavior
  distributed_exports:
    current: deferred
    future_vocabulary: distributed_service_export_deferral
    status: no_distributed_runtime_behavior
```

## 5. Ownership

```yaml
architecture_vocabulary_owner:
  - docs/pitaya-aligned-service-export-boundary-gate.md
  - .arch/reference.yaml
  - .arch/runtime.yaml
source_first_map_candidate_owner:
  - tools/vibit
runtime_behavior_owner: deferred
protocol_owner: unchanged
persistence_owner: unchanged
distributed_runtime_owner: deferred
```

Rules:

- Documentation and manifests may define service export vocabulary and current source-first mapping.
- `tools/vibit` may emit a source-first service export map because W-0304 is separately recorded.
- Existing runtime, protocol, route, repository, persistence, and transport behavior remain unchanged by this gate.

## 6. Reference Mapping

Nakama remains the primary product reference for broad game backend product capability pressure. Pitaya remains an architecture vocabulary reference for distributed runtime, components, handlers, services, sessions, routes, remote calls, service discovery, groups, lifecycle hooks, and operational concerns.

This gate maps those references into vibit-owned vocabulary only. It does not create direct compatibility, public API parity, runtime behavior, or distributed behavior.

## 7. Stop Conditions

Stop and require a later bounded work item before adding:

- service export behavior;
- component service registration;
- backend service export behavior;
- runtime behavior;
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
- direct Nakama/Pitaya API compatibility;

## 8. Verification

Repository verification:

```text
node tools/vibit inspect rule runtime.pitaya_aligned_service_export_boundary_gate
node tools/vibit check change define-pitaya-aligned-service-export-boundary-gate --json
node tools/vibit check runtime --json
node tools/vibit check work --json
```
