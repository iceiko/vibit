# Storage Objects Runtime Behavior Gate

Status: Accepted v0.1
Last updated: 2026-05-22
Scope: Gate-only boundary for future application-owned storage objects runtime behavior after the PostgreSQL adapter
Depends on: `docs/storage-objects-behavior-gate.md`, `docs/storage-objects-repository-boundary.md`, `docs/storage-objects-postgresql-adapter-gate.md`, `runtime/internal/modules/storage/repository.go`, `runtime/internal/platform/persistence/postgres/storage_object_repository.go`, `docs/runtime-protocol-adapter.md`, `docs/bound-identity-route-policy-gate.md`, `docs/runtime-session-validation-gate.md`, `docs/reference-game-server-alignment.md`
Canonical decision: `ADR-0116`

The paired Simplified Chinese translation is `docs/storage-objects-runtime-behavior-gate.zh-CN.md`. The English file is authoritative.

This document defines the storage objects runtime behavior gate. It is a gate artifact. It does not add runtime behavior implementation, runtime handlers, startup wiring, protocol routes, Protobuf source, generated output, repository interface changes, PostgreSQL adapter changes, migration changes, dependencies, authentication/session behavior changes, hosted deployments, release artifacts, public announcements, paid promotion, broad product module expansion, large object/blob storage, S3-compatible object storage, or direct Nakama/Pitaya API compatibility.

## 1. Core Rule

The storage objects runtime behavior gate record is:

```yaml
storage_objects_runtime_behavior_gate: defined
implementation_authorized_by_this_standard: false
completed_work_item: W-0208
decision: ADR-0116
check_rule: runtime.storage_objects_runtime_behavior_gate
source_adapter_implementation_decision: ADR-0115
source_adapter: runtime/internal/platform/persistence/postgres/storage_object_repository.go
source_repository_interface_decision: ADR-0113
repository_interface: runtime/internal/modules/storage.Repository
repository_interface_source: runtime/internal/modules/storage/repository.go
future_runtime_owner_candidate: runtime/internal/app
future_storage_application_package_candidate: runtime/internal/app/storage
future_runtime_service_source_candidate: runtime/internal/app/storage/service.go
future_runtime_service_test_candidate: runtime/internal/app/storage/service_test.go
request_identity_required: true
metadata_only_player_id_allowed_as_proof: false
first_owner_kind: player
owner_id_source: validated_request_identity_player_id
route_policy_requirement: request_token_required
service_application_owner: runtime/internal/app
repository_handoff: unit_of_work_storage_repository_factory
unit_of_work_handoff_required: true
runtime_behavior_gate_only: true
runtime_behavior_added: false
runtime_handlers_added: false
startup_wiring_added: false
protocol_route_added: false
protobuf_source_added: false
generated_output_added: false
repository_interface_changed: false
postgresql_adapter_changed: false
dependency_added: false
migration_added: false
authentication_session_behavior_changed: false
large_object_blob_storage_added: false
s3_compatible_object_storage_added: false
direct_nakama_pitaya_api_compatibility_added: false
future_runtime_behavior_implementation_work_item: W-0209
future_runtime_behavior_implementation_direction: storage_objects_runtime_behavior_implementation
```

## 2. Purpose

`W-0207` implemented the PostgreSQL adapter for `runtime/internal/modules/storage.Repository`. The next useful boundary is not a route or protocol change. The next useful boundary is the runtime behavior gate that defines how application code may later turn a validated route request into storage repository operations.

This gate records the future behavior shape before implementation:

- application ownership for the service;
- owner identity derivation from validated request identity;
- permission and route-policy posture;
- validation and conflict mapping expectations;
- unit-of-work and repository handoff;
- redaction rules;
- test expectations;
- stop conditions that keep protocol, generated output, authentication/session changes, and broader storage products out of this slice.

Nakama motivates a durable storage-object gameplay capability. Pitaya motivates keeping handlers, route policy, and persistence responsibilities separated. vibit adapts those references through explicit application-owned behavior and checks, not direct public API compatibility.

## 3. Ownership

Future runtime behavior is application-owned:

```yaml
future_runtime_owner_candidate: runtime/internal/app
future_storage_application_package_candidate: runtime/internal/app/storage
future_runtime_service_source_candidate: runtime/internal/app/storage/service.go
future_runtime_service_test_candidate: runtime/internal/app/storage/service_test.go
repository_interface_owner: runtime/internal/modules/storage
postgresql_adapter_owner: runtime/internal/platform/persistence/postgres
route_policy_owner: runtime/internal/app
protocol_adapter_owner: runtime/internal/platform/protocol/protobuf
websocket_transport_owner: runtime/internal/platform/transport/ws
```

Rules:

- Future service behavior may live under `runtime/internal/app/storage` or an equivalent application-owned package ratified by the implementation slice.
- The service may call `runtime/internal/modules/storage.Repository` only through application or unit-of-work dependencies.
- The service must not import PostgreSQL adapter packages, SQL row types, migration packages, WebSocket transport packages, generated Protobuf packages, S3 SDKs, or MinIO clients.
- The storage module remains the owner of storage-neutral value types, validation helpers, and repository error vocabulary.
- The PostgreSQL adapter remains persistence-only and must not derive request identity, route policy, or public protocol errors.
- Transport and protocol adapters must not own storage object permission or business behavior.

## 4. Request Identity And Owner Derivation

The first runtime behavior posture is player-owned:

```yaml
first_owner_kind: player
owner_id_source: validated_request_identity_player_id
request_identity_required: true
metadata_only_player_id_allowed_as_proof: false
metadata_only_session_id_allowed_as_proof: false
client_supplied_owner_id_allowed: false
owner_id_overrides_allowed: false
```

Rules:

- A future storage operation must derive `storage.StorageObjectOwner{Kind: player, ID: identity.PlayerID}` from a validated `app.RequestIdentity`.
- `RequestIdentity.Status` must be `validated`.
- `RequestIdentity.ActorKind` must be `player`.
- `RequestIdentity.PlayerIDValidated` must be true.
- `RequestIdentity.PlayerID` must be non-empty and must match the actor identity when both are present.
- Metadata-only `player_id` from envelope/session metadata must never satisfy this gate.
- A persisted `session_id` alone must never become proof.
- Client payloads must not choose another owner id in the first posture.

This gate does not change `RequestIdentity`, access-token validation, bound connection identity, durable runtime session validation, or WebSocket handshake behavior. It only records the identity requirements that future storage behavior must require before repository access.

## 5. Future Runtime Behavior Shape

The future first implementation may expose an application service with these candidate operations:

```yaml
candidate_operations:
  - get_own_storage_object
  - list_own_storage_objects
  - put_own_storage_object
  - delete_own_storage_object
```

Recommended first posture:

- `get` reads one active object by collection and key for the validated player.
- `list` reads active objects in one collection for the validated player with bounded pagination.
- `put` creates or replaces one object for the validated player and returns the current version.
- `delete` soft-deletes one object for the validated player and returns stable success or conflict behavior.

Rules:

- Runtime behavior must use server-derived owner identity.
- Runtime behavior must validate collection, key, value shape, value size, expected version, and list limit before repository calls.
- Runtime behavior must not expose cross-owner reads, cross-owner writes, public ACLs, admin bypass, group/guild scopes, room scopes, match scopes, batch writes, JSON patch, merge semantics, TTL, or server script hooks in the first implementation.
- Runtime behavior must not add public protocol routes or generated output unless a later protocol gate authorizes them.

## 6. Candidate Application Service Shape

The first implementation slice should define a small application-owned service boundary. Candidate inputs and outputs:

```yaml
candidate_request_fields:
  - request_identity
  - collection
  - key
  - value_json
  - expected_version
  - list_limit
  - after_object_key
candidate_result_fields:
  - object
  - objects
  - next_object_key
  - version
  - public_error_code
```

Rules:

- The service should accept already-normalized application identity, not raw tokens, cookies, headers, WebSocket subprotocol values, or envelope proof carriers.
- The service should call storage module normalizers before repository handoff.
- The service should keep raw JSON values out of default errors and logs.
- The service should expose stable public error codes or classes for runtime handlers to map later.
- The service should not add route registration, Protobuf conversion, startup composition, or command/query dispatch wiring in the gate slice.

## 7. Validation Rules

Future runtime behavior must enforce validation before persistence:

```yaml
validation_required:
  - validated_player_identity
  - collection_non_empty_length_bounded
  - key_non_empty_length_bounded
  - value_top_level_json_object
  - value_size_bounded
  - expected_version_positive_when_present
  - list_limit_bounded
  - pagination_cursor_bounded
```

Rules:

- Collection and key validation should reuse storage module normalization rules unless a future contract explicitly tightens protocol-visible syntax.
- Value JSON is not log-safe and must remain copied or immutable across boundaries.
- Missing expected version behavior must be explicit in implementation tests.
- Invalid input must fail before repository mutation when possible.
- Repository unavailable errors must remain redacted.

## 8. Permission And Route Policy Posture

The first route-policy posture is conservative:

```yaml
route_policy_requirement: request_token_required
public_storage_routes_allowed: false
bound_connection_required_by_this_gate: false
session_validated_required_by_this_gate: false
bound_session_required_by_this_gate: false
```

Candidate permission families for later public contracts:

- read own storage object;
- list own storage objects;
- write own storage object;
- delete own storage object.

Rules:

- Storage object routes must be protected routes.
- The first posture should use the existing `request_token_required` route protection family unless a later route-policy ADR changes named routes.
- Public routes must not read or mutate storage objects.
- Bound connection identity and durable session validation may remain available for future route families, but this gate does not require them or change ordinary protected route behavior.
- Metadata-only identity must fail closed.

## 9. Conflict And Error Mapping

Future runtime behavior must map storage repository errors into stable application classes:

```yaml
candidate_public_error_classes:
  - STORAGE_OBJECT_INVALID_REQUEST
  - STORAGE_OBJECT_NOT_FOUND
  - STORAGE_OBJECT_ALREADY_EXISTS
  - STORAGE_OBJECT_VERSION_MISMATCH
  - STORAGE_OBJECT_UNAVAILABLE
  - STORAGE_OBJECT_FORBIDDEN
```

Rules:

- Not-found, owner mismatch, and deleted-object cases must avoid cross-player existence leaks.
- Version mismatch may be public as a conflict class, but stored values, raw JSON, SQL details, driver errors, DSNs, credentials, token material, verifier digests, and route proof carriers must not leak.
- Repository `storage_unavailable` errors must map to a retryable or unavailable class without exposing platform internals.
- Permission failure must happen before repository access when the request identity is not validated.
- Runtime behavior must not add authentication/token/session failure detail beyond existing application route-protection classes.

## 10. Unit-Of-Work And Repository Handoff

Future runtime behavior should use the existing application transaction boundary:

```yaml
repository_handoff: unit_of_work_storage_repository_factory
unit_of_work_handoff_required: true
service_starts_transactions: false
service_commits_transactions: false
service_rolls_back_transactions: false
repository_factory_candidate: NewStorageObjectRepository
```

Rules:

- Mutating operations should run inside the existing command unit-of-work boundary when called from command dispatch.
- The service should obtain a `storage.Repository` from an application dependency or from a unit-of-work capability ratified by the implementation.
- The service must not import PostgreSQL packages just to create repositories.
- Query behavior may use an explicit query repository provider if the implementation chooses a read path, but the owner derivation and validation rules remain the same.
- This gate does not change `TransactionalDispatcher` or startup wiring.

## 11. Relationship To Protocol

This gate does not add protocol behavior:

```yaml
storage_protocol_routes_added: false
protobuf_storage_messages_added: false
generated_storage_output_added: false
existing_envelope_changed: false
websocket_transport_changed: false
```

Rules:

- Future protocol work must define exact route names, module/name routing, request/response messages, generated output, and error mapping in a separate bounded work item.
- Candidate route names from `docs/storage-objects-behavior-gate.md` remain planning vocabulary only.
- WebSocket transport must remain credential-neutral.
- Protobuf adapters must not derive permissions or owner identity.

## 12. Relationship To Authentication And Session

This gate uses existing validated request identity vocabulary but does not change authentication or session behavior:

```yaml
access_token_validation_changed: false
request_identity_shape_changed: false
session_validation_changed: false
websocket_handshake_authentication_changed: false
first_message_connection_binding_changed: false
metadata_only_player_id_allowed_as_proof: false
```

Rules:

- Access-token validation remains the current proof path for protected routes.
- Durable session validation remains separately owned.
- First-message connection binding does not authorize ordinary storage operations by this gate.
- WebSocket handshake authentication remains deferred.
- The implementation slice must not relax metadata-only identity protections.

## 13. Test Expectations

The later implementation slice should add focused tests for:

```yaml
future_tests:
  - metadata_only_identity_rejected_before_repository
  - validated_player_identity_derives_owner
  - client_owner_id_ignored_or_rejected
  - get_maps_repository_object_to_runtime_result
  - list_is_owner_collection_scoped_and_bounded
  - put_validates_json_object_and_version
  - delete_checks_expected_version_when_supplied
  - owner_mismatch_and_not_found_do_not_leak_existence
  - repository_errors_are_redacted
  - unit_of_work_repository_handoff_is_used_for_mutations
  - no_protocol_route_or_generated_output_required
```

Rules:

- Tests should use fake repositories or fake unit-of-work providers where possible.
- Default repository checks must not require live PostgreSQL.
- Tests must not print raw values, DSNs, credentials, access tokens, verifier material, lookup digests, verifier digests, cookies, headers, query strings, or WebSocket subprotocol values.

## 14. Stop Conditions

Stop and require a later bounded work item before adding:

- runtime behavior implementation;
- runtime handlers or route registration;
- protocol routes;
- Protobuf source files;
- generated output;
- repository interface changes;
- PostgreSQL adapter changes;
- migration changes;
- new dependencies;
- startup wiring;
- authentication/session behavior changes;
- WebSocket handshake authentication;
- public ACLs, admin storage search, group/guild/party/room/match storage scopes, batch writes, JSON patch, merge semantics, TTL, or script hooks;
- large object/blob storage;
- S3-compatible object storage;
- hosted deployments;
- release artifacts;
- public announcements beyond the GitHub release record;
- paid promotion;
- direct Nakama/Pitaya API compatibility.

## 15. Verification

The gate is verified by repository checks:

```yaml
check_rule: runtime.storage_objects_runtime_behavior_gate
required_commands:
  - node -c tools/vibit
  - node tools/vibit inspect next --json
  - node tools/vibit inspect rule runtime.storage_objects_runtime_behavior_gate
  - node tools/vibit check change define-storage-objects-runtime-behavior-gate --json
  - node tools/vibit check module storage --json
  - node tools/vibit check work --json
  - node tools/vibit check runtime --json
  - node tools/vibit check memory --json
  - node tools/vibit check schemas --json
  - node tools/vibit check all --json
  - cd runtime && go test ./...
  - git diff --check
```

No live PostgreSQL verification is required for this gate because it adds no runtime behavior or SQL execution.
