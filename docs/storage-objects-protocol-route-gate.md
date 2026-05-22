# Storage Objects Protocol Route Gate

Status: Accepted v0.1
Last updated: 2026-05-22
Scope: Gate-only boundary for future client-facing storage objects protocol routes after application runtime behavior
Depends on: `docs/storage-objects-behavior-gate.md`, `docs/storage-objects-runtime-behavior-gate.md`, `decisions/ADR-0117-storage-objects-runtime-behavior-implementation.md`, `docs/runtime-protocol-adapter.md`, `docs/game-protocol.md`, `docs/generated-output.md`, `docs/reference-game-server-alignment.md`, `docs/nakama-pitaya-product-parity-roadmap.md`
Canonical decision: `ADR-0118`

The paired Simplified Chinese translation is `docs/storage-objects-protocol-route-gate.zh-CN.md`. The English file is authoritative.

This document defines the storage objects protocol route gate. It is a gate artifact. It does not add protocol route implementation, Protobuf source, generated output, startup wiring, runtime handlers, repository interface changes, PostgreSQL adapter changes, migration changes, dependencies, authentication/session behavior changes, hosted deployments, release artifacts, public announcements, paid promotion, broad product module expansion, large object/blob storage, S3-compatible object storage, or direct Nakama/Pitaya API compatibility.

## 1. Core Rule

The storage objects protocol route gate record is:

```yaml
storage_objects_protocol_route_gate: defined
implementation_authorized_by_this_standard: false
completed_work_item: W-0210
decision: ADR-0118
check_rule: runtime.storage_objects_protocol_route_gate
source_runtime_behavior_implementation_decision: ADR-0117
source_runtime_behavior_implementation: runtime/internal/app/storage/service.go
source_runtime_behavior_tests: runtime/internal/app/storage/service_test.go
source_runtime_behavior_gate_decision: ADR-0116
source_repository_interface_decision: ADR-0113
repository_interface: runtime/internal/modules/storage.Repository
future_protocol_source_candidate: proto/vibit/storage/v1/storage.proto
future_generated_go_output_candidate: runtime/internal/generated/proto/vibit/storage/v1/storage.pb.go
future_protocol_bridge_candidate: runtime/internal/platform/protocol/protobuf/storage_bridge.go
future_protocol_bridge_test_candidate: runtime/internal/platform/protocol/protobuf/storage_bridge_test.go
future_application_handler_candidate: runtime/internal/app/bootstrap/storage.go
future_application_handler_test_candidate: runtime/internal/app/bootstrap/storage_test.go
route_policy_requirement: request_token_required
authenticated_wrapper_required: true
request_identity_source: validated_authenticated_request_identity
metadata_only_player_id_allowed_as_proof: false
metadata_only_session_id_allowed_as_proof: false
client_supplied_owner_id_allowed: false
first_owner_kind: player
first_payload_package: vibit.storage.v1
protobuf_envelope_change_status: unchanged
websocket_transport_credential_neutral: true
protocol_route_gate_only: true
protocol_route_added: false
protobuf_source_added: false
generated_output_added: false
runtime_handler_added: false
startup_wiring_added: false
repository_interface_changed: false
postgresql_adapter_changed: false
dependency_added: false
migration_added: false
authentication_session_behavior_changed: false
large_object_blob_storage_added: false
s3_compatible_object_storage_added: false
direct_nakama_pitaya_api_compatibility_added: false
future_protocol_route_implementation_work_item: W-0211
future_protocol_route_implementation_direction: storage_objects_protocol_route_implementation
```

## 2. Purpose

`W-0209` implemented application-owned storage object behavior under `runtime/internal/app/storage`. The next useful boundary is not route code or `.proto` generation. The next useful boundary is a protocol route gate that records how future WebSocket/Protobuf exposure should call that service without moving storage behavior into transport, generated files, or persistence adapters.

Nakama motivates the product surface: client-facing storage objects are a common game backend capability with collection/key identity, owner scope, permissions, version conflict behavior, list/read/write/delete operations, and durable JSON game state. vibit should cover that capability class.

Pitaya motivates the architecture posture: acceptors, sessions, route handlers, serializers, and backend services should remain separated. vibit adapts that by keeping WebSocket transport credential-neutral, keeping Protobuf payload bridging explicit, and invoking application-owned route handlers that call application-owned storage services.

This gate records the future protocol shape before implementation:

- candidate route names;
- candidate request and response message shapes;
- route protection and identity handoff posture;
- protocol adapter, application handler, and startup ownership;
- generated-output expectations;
- error mapping and redaction expectations;
- Nakama/Pitaya reference mapping;
- stop conditions that keep implementation and generated artifacts out of this slice.

## 3. Future Route Surface

The first route family should expose own-player storage object operations only:

```yaml
candidate_routes:
  - kind: query
    module: storage
    name: GetOwnStorageObject
    route_id: storage.GetOwnStorageObject
    service_method: GetOwnStorageObject
  - kind: query
    module: storage
    name: ListOwnStorageObjects
    route_id: storage.ListOwnStorageObjects
    service_method: ListOwnStorageObjects
  - kind: command
    module: storage
    name: PutOwnStorageObject
    route_id: storage.PutOwnStorageObject
    service_method: PutOwnStorageObject
  - kind: command
    module: storage
    name: DeleteOwnStorageObject
    route_id: storage.DeleteOwnStorageObject
    service_method: DeleteOwnStorageObject
```

Rules:

- The route names must stay vibit-native and must not copy Nakama route paths or Pitaya route naming.
- `GetOwnStorageObject` and `ListOwnStorageObjects` are queries.
- `PutOwnStorageObject` and `DeleteOwnStorageObject` are commands.
- The first route family is only for the validated player owner. It must not expose arbitrary owner ids.
- Public ACLs, admin search, group/guild/party/room/match scopes, batch writes, JSON patch, merge semantics, TTL, script hooks, and large object/blob storage remain deferred.
- Future route implementation must register routes explicitly. No catch-all storage route or reflective handler is allowed.

## 4. Protocol Shape

The first storage object protocol source candidate is:

```text
proto/vibit/storage/v1/storage.proto
```

The first generated output candidate is:

```text
runtime/internal/generated/proto/vibit/storage/v1/storage.pb.go
```

The first Protobuf package candidate is:

```text
vibit.storage.v1
```

Candidate messages:

```yaml
messages:
  StorageObject:
    fields:
      collection: string
      key: string
      value_json: string
      version: int64
      created_at: string
      updated_at: string
  GetOwnStorageObjectRequest:
    fields:
      collection: string
      key: string
  GetOwnStorageObjectResponse:
    fields:
      object: StorageObject
  ListOwnStorageObjectsRequest:
    fields:
      collection: string
      limit: int32
      after_key: string
  ListOwnStorageObjectsResponse:
    fields:
      objects: repeated StorageObject
      next_key: string
  PutOwnStorageObjectRequest:
    fields:
      collection: string
      key: string
      value_json: string
      expected_version: int64
  PutOwnStorageObjectResponse:
    fields:
      object: StorageObject
      version: int64
  DeleteOwnStorageObjectRequest:
    fields:
      collection: string
      key: string
      expected_version: int64
  DeleteOwnStorageObjectResponse:
    fields:
      deleted: bool
      version: int64
```

Rules:

- The existing `proto/vibit/protocol/v1/envelope.proto` must remain unchanged unless a later protocol ADR explicitly changes envelope semantics.
- `value_json` is not log-safe. It must not appear in default error messages, logs, route policy diagnostics, or test names.
- Time values should use RFC3339Nano UTC text when exposed.
- `expected_version` should use `0` or field absence semantics only if the implementation ADR makes the no-precondition posture explicit. The service already has optional expected-version vocabulary; future Protobuf mapping must preserve that distinction without inventing merge semantics.
- The protocol shape must not include `owner_id`, `player_id`, `session_id`, raw access tokens, credential material, lookup digests, verifier digests, SQL details, blob bytes, S3 bucket names, or transport metadata.

## 5. Route Protection And Identity Handoff

The first route-policy posture is:

```yaml
route_policy_requirement: request_token_required
authenticated_wrapper_required: true
request_identity_source: validated_authenticated_request_identity
metadata_only_player_id_allowed_as_proof: false
metadata_only_session_id_allowed_as_proof: false
client_supplied_owner_id_allowed: false
```

Rules:

- Future storage routes must be protected gameplay routes.
- Future handlers must receive a validated `app.RequestIdentity` from the existing protected-route flow.
- Metadata-only `player_id` or `session_id` from envelope/session metadata must never become storage owner proof.
- Client payloads must not choose owner ids.
- The service remains responsible for rejecting invalid identity before repository access.
- This gate does not change authentication, token validation, session persistence, first-message binding, WebSocket handshake behavior, bound-identity policy, or route-protection semantics.

## 6. Future Route Flow

Future implementation must preserve this sequence:

```yaml
storage_route_flow:
  - websocket_transport_receives_binary_frame_without_reading_credentials
  - protobuf_adapter_decodes_existing_envelope
  - route_policy_requires_authenticated_request_wrapper
  - route_policy_validates_access_token_through_existing authentication/session behavior
  - protobuf_adapter_decodes storage request payload
  - protocol_bridge_maps_payload_to runtime/internal/app/storage request
  - application_handler_calls storage.Service own-object method
  - storage_service_derives_owner_from validated app.RequestIdentity
  - storage_service_uses unit-of-work NewStorageObjectRepository handoff
  - protocol_bridge_maps storage service result to storage response payload
  - protobuf_adapter_returns success or existing error envelope
```

Rules:

- WebSocket transport remains credential-neutral and payload-neutral.
- Protobuf adapter may decode and encode storage payloads but must not own permission decisions, repository calls, or storage behavior.
- Application handler registration belongs in `runtime/internal/app/bootstrap` or an equivalent application-composition package.
- The handler may call only the application storage service. It must not call repositories directly, open ad hoc SQL, import PostgreSQL adapters, or parse transport metadata as proof.
- Normal query/command transaction wrapping may be used only through existing application dispatch/transaction boundaries. This gate does not add startup composition.

## 7. Error Mapping

Future protocol behavior must map service public errors through existing application error envelopes:

```yaml
service_public_errors:
  STORAGE_OBJECT_INVALID_REQUEST: application_error_same_code
  STORAGE_OBJECT_NOT_FOUND: application_error_same_code
  STORAGE_OBJECT_ALREADY_EXISTS: application_error_same_code
  STORAGE_OBJECT_VERSION_MISMATCH: application_error_same_code
  STORAGE_OBJECT_UNAVAILABLE: application_error_same_code
  STORAGE_OBJECT_FORBIDDEN: application_error_same_code
```

Rules:

- Not-found, owner mismatch, and deleted-object cases must not leak cross-player existence.
- Version mismatch may be public as a conflict class.
- Errors must not include raw JSON values, owner ids beyond already validated caller identity, raw token material, credential material, lookup digests, verifier digests, HMAC input/output, SQL strings, database errors, DSNs, headers, cookies, query strings, WebSocket subprotocol values, remote addresses, connection ids, or session ids.
- Protocol adapter errors for malformed payloads must stay distinct from service validation errors without exposing payload contents.

## 8. Ownership

Future implementation must preserve these owners:

```yaml
storage_service_owner: runtime/internal/app/storage
application_handler_owner: runtime/internal/app/bootstrap
route_policy_owner: runtime/internal/app
protocol_adapter_owner: runtime/internal/platform/protocol/protobuf
storage_repository_interface_owner: runtime/internal/modules/storage
postgresql_adapter_owner: runtime/internal/platform/persistence/postgres
websocket_transport_owner: runtime/internal/platform/transport/ws
startup_owner: runtime/cmd/vibit-server
generated_output_owner: runtime/internal/generated/proto/vibit/storage/v1
protobuf_source_owner: proto/vibit/storage/v1
```

Rules:

- Storage runtime behavior remains in `runtime/internal/app/storage`.
- Protocol bridge code may map payload fields only.
- Persistence code remains storage adapter behavior only.
- Generated output must be produced from `.proto` sources and must not be hand-edited.
- Startup wiring, route registration, and generated output remain behind the later implementation work item.

## 9. Required Future Tests

The future implementation slice must add focused tests for:

```yaml
required_tests:
  proto_source_and_generated_output_include_storage_messages
  storage_routes_are_registered_only_when_storage_service_is_composed
  storage_routes_are_protected_and_require_authenticated_wrapper
  storage_routes_do_not_accept_metadata_only_player_id_or_session_id
  storage_route_payloads_do_not_include_owner_id
  get_request_maps_to_GetOwnStorageObject
  list_request_maps_to_ListOwnStorageObjects
  put_request_maps_to_PutOwnStorageObject
  delete_request_maps_to_DeleteOwnStorageObject
  storage_success_maps_service_results_to_response_payloads
  storage_public_errors_map_to_error_envelopes
  storage_errors_do_not_leak_value_json_or_repository_details
  websocket_transport_remains_credential_neutral
  existing_protobuf_envelope_remains_unchanged
```

Live PostgreSQL verification remains optional and must not be required by default repository checks.

## 10. Nakama And Pitaya Reference Mapping

Nakama reference mapping:

- Adopt the product expectation that durable player-owned storage objects are a first-class game backend capability.
- Adapt collection/key/value/version/list/read/write/delete semantics into vibit's service and protocol model.
- Adapt optimistic version conflict semantics without copying Nakama route paths, permission integers, data model names, server runtime APIs, JavaScript/Lua hook model, or direct compatibility.
- Defer public ACLs, cross-user reads, system-owned storage, admin search, batch writes, match/party/group scoped storage, TTL, script hooks, and storage object API compatibility.

Pitaya reference mapping:

- Adopt separation between acceptors, sessions, routes, serializers, handlers, and backend behavior.
- Adapt handler routing by keeping `kind/module/name` route identity explicit and by making application handlers call application services.
- Keep WebSocket acceptor behavior credential-neutral and storage-neutral.
- Defer Pitaya route naming compatibility, frontend/backend cluster routing, remote calls, groups integration, distributed push, and RPC/session propagation.

## 11. Stop Conditions

Stop and open a later bounded work item before any of the following:

- protocol route implementation;
- Protobuf source creation;
- generated output creation or editing;
- protocol bridge implementation;
- application route registration;
- startup wiring;
- repository interface changes;
- PostgreSQL adapter changes;
- migration changes;
- new dependencies;
- authentication/session behavior changes;
- route-protection semantic changes;
- public ACLs or cross-owner access;
- admin search;
- group/guild/party/room/match storage scopes;
- batch writes;
- JSON patch or merge semantics;
- TTL or script hooks;
- large object/blob storage;
- S3-compatible object storage;
- hosted deployments;
- release artifacts;
- public announcements;
- paid promotion;
- direct Nakama/Pitaya API compatibility.

## 12. Verification

The repository check rule for this gate is:

```text
runtime.storage_objects_protocol_route_gate
```

Expected verification for this gate:

```bash
node -c tools/vibit
node tools/vibit inspect next --json
node tools/vibit inspect rule runtime.storage_objects_protocol_route_gate
node tools/vibit check change define-storage-objects-protocol-route-gate --json
node tools/vibit check module storage --json
node tools/vibit check work --json
node tools/vibit check runtime --json
node tools/vibit check memory --json
node tools/vibit check schemas --json
node tools/vibit check all --json
git diff --check
```

Go tests are not required for this gate because it must not add Go runtime behavior.
