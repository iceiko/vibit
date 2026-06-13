# Pitaya-Aligned Remote Call Dispatch Boundary Gate 中文版

Status: Accepted v0.1
Last updated: 2026-06-13
Scope: Gate-only boundary for using Pitaya-aligned remote call dispatch vocabulary after the currency wallet protocol route implementation
Depends on: `decisions/ADR-0210-currency-wallet-protocol-route-implementation.md`, `docs/reference-game-server-alignment.md`, `.arch/reference.yaml`
Canonical decision: `ADR-0213`

说明：本文件是 `docs/pitaya-aligned-remote-call-dispatch-boundary-gate.md` 的简体中文译本。英文版本是权威版本。

This document defines remote call dispatch vocabulary only. It does not implement remote call dispatch behavior, server-to-server RPC behavior, remote timeout or retry behavior, runtime behavior, protocol messages or routes, Protobuf source, generated output, repository interfaces, PostgreSQL adapters, migrations, dependencies, hosted deployment, SDK publication, release artifacts, distributed runtime behavior, or direct Nakama/Pitaya API compatibility.

## 1. Core Rule

```yaml
pitaya_aligned_remote_call_dispatch_boundary_gate: defined
completed_work_item: W-0305
decision: ADR-0213
check_rule: runtime.pitaya_aligned_remote_call_dispatch_boundary_gate
source_currency_wallet_protocol_route_implementation_decision: ADR-0210
standard: docs/pitaya-aligned-remote-call-dispatch-boundary-gate.md
translation: docs/pitaya-aligned-remote-call-dispatch-boundary-gate.zh-CN.md
primary_product_reference: Nakama
pitaya_reference_status: remote_call_dispatch_vocabulary_boundary_defined_for_future_architecture_planning
implementation_scope: gate_only_remote_call_dispatch_vocabulary
future_implementation_work_item: W-0306
future_implementation_direction: implement_pitaya_aligned_remote_call_dispatch_source_first_map
allowed_remote_call_dispatch_vocabulary:
  - remote_call_dispatch_boundary
  - local_application_dispatch_source
  - server_to_server_rpc_deferral
  - timeout_retry_deferral
  - remote_error_mapping_deferral
  - distributed_targeting_deferral
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

The purpose is to keep remote call dispatch visible as a future Pitaya-class architecture concern without authorizing implementation. Agents may use this vocabulary for planning, inspection, ADRs, change specs, and future source-first maps only.

## 3. Vocabulary

- `remote_call_dispatch_boundary`: planning vocabulary for remote call dispatch without implementing remote call dispatch behavior.
- `local_application_dispatch_source`: planning vocabulary for remote call dispatch without implementing remote call dispatch behavior.
- `server_to_server_rpc_deferral`: planning vocabulary for remote call dispatch without implementing remote call dispatch behavior.
- `timeout_retry_deferral`: planning vocabulary for remote call dispatch without implementing remote call dispatch behavior.
- `remote_error_mapping_deferral`: planning vocabulary for remote call dispatch without implementing remote call dispatch behavior.
- `distributed_targeting_deferral`: planning vocabulary for remote call dispatch without implementing remote call dispatch behavior.

Forbidden use:

- Do not introduce concrete public API, package, route, method, wire, handler, dashboard, metrics, trace, admin, console, or inspector compatibility names from Pitaya or Nakama.
- Do not use this vocabulary as permission to add remote call dispatch behavior, server-to-server RPC behavior, remote timeout or retry behavior, runtime behavior, protocol messages, generated output, persistence, dependencies, hosted surfaces, SDKs, release artifacts, or distributed runtime behavior.
- Do not classify raw tokens, credentials, lookup digests, verifier digests, verifier keys, DSNs, headers, cookies, query strings, subprotocol values, remote addresses, database payloads, local secret file contents, route payloads, event payloads, or concrete transport metadata as log-safe in this gate.

## 4. Current Mapping Candidate

```yaml
current_source_first_remote_call_dispatch_mapping:
  local_application_dispatch:
    current: runtime/internal/app dispatcher and transactional dispatcher sources
    future_vocabulary: local_application_dispatch_source
    status: single_process_dispatch_only
  protocol_route_dispatch:
    current: runtime/internal/platform/protocol/protobuf bridge sources
    future_vocabulary: remote_error_mapping_deferral
    status: no_remote_call_dispatch_behavior
  server_to_server_rpc:
    current: deferred
    future_vocabulary: server_to_server_rpc_deferral
    status: no_server_to_server_rpc_behavior
  remote_timeouts_and_retries:
    current: deferred
    future_vocabulary: timeout_retry_deferral
    status: no_remote_timeout_retry_behavior
  distributed_targets:
    current: deferred
    future_vocabulary: distributed_targeting_deferral
    status: no_distributed_runtime_behavior
```

## 5. Ownership

```yaml
architecture_vocabulary_owner:
  - docs/pitaya-aligned-remote-call-dispatch-boundary-gate.md
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

- Documentation and manifests may define remote call dispatch vocabulary and current source-first mapping.
- `tools/vibit` may emit a source-first remote call dispatch map because W-0306 is separately recorded.
- Existing runtime, protocol, route, repository, persistence, and transport behavior remain unchanged by this gate.

## 6. Reference Mapping

Nakama remains the primary product reference for broad game backend product capability pressure. Pitaya remains an architecture vocabulary reference for distributed runtime, components, handlers, services, sessions, routes, remote calls, service discovery, groups, lifecycle hooks, and operational concerns.

This gate maps those references into vibit-owned vocabulary only. It does not create direct compatibility, public API parity, runtime behavior, or distributed behavior.

## 7. Stop Conditions

Stop and require a later bounded work item before adding:

- remote call dispatch behavior;
- server-to-server RPC behavior;
- remote timeout or retry behavior;
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
node tools/vibit inspect rule runtime.pitaya_aligned_remote_call_dispatch_boundary_gate
node tools/vibit check change define-pitaya-aligned-remote-call-dispatch-boundary-gate --json
node tools/vibit check runtime --json
node tools/vibit check work --json
```
