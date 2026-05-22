# ADR-0117: Storage Objects Runtime Behavior Implementation

Status: Accepted
Date: 2026-05-22
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-22-implement-storage-objects-runtime-behavior/`

Related conversations:

- `conversations/2026-05-22-storage-objects-runtime-behavior-implementation.md`

Related artifacts:

- `runtime/internal/app/storage/service.go`
- `runtime/internal/app/storage/service_test.go`
- `modules/storage/module.yaml`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `.arch/modules.yaml`
- `tools/vibit`
- `rules/check-rules.json`

## Context

`M-136/W-0208` defined the storage objects runtime behavior gate after the storage-neutral repository interface and PostgreSQL adapter were in place. The gate authorized a later bounded application-owned runtime behavior implementation under `runtime/internal/app/storage`.

The storage module owns `runtime/internal/modules/storage.Repository` and storage-neutral value types. The PostgreSQL adapter owns SQL mapping and exposes `UnitOfWork.NewStorageObjectRepository`. Runtime behavior must therefore compose those pieces through the application unit-of-work boundary without importing persistence, protocol, transport, generated Protobuf, migration, or object-storage packages.

Nakama motivates durable player-owned storage objects as a useful game backend capability. Pitaya reinforces separating route/session context, handlers, and persistence. vibit adapts those references through validated request identity and application-owned behavior, not through direct public API compatibility.

## Decision

Implement storage objects runtime behavior under:

```text
runtime/internal/app/storage
```

The implementation adds:

- `Service` and `NewService`;
- `GetOwnStorageObject`;
- `ListOwnStorageObjects`;
- `PutOwnStorageObject`;
- `DeleteOwnStorageObject`;
- request and result vocabulary for application callers;
- stable public error codes for later handlers to map;
- validated-player identity owner derivation;
- metadata-only `player_id` and `session_id` refusal before repository access;
- collection, key, value JSON, value size, expected-version, and list-pagination validation before repository calls where possible;
- unit-of-work storage repository handoff through `NewStorageObjectRepository`;
- create-or-replace first `put` behavior when no expected version is supplied;
- expected-version update behavior when an expected version is supplied;
- soft-delete service behavior through repository delete;
- redacted conflict mapping for not found, already exists, version mismatch, invalid request, unavailable, and forbidden cases;
- focused fake-repository tests.

This ADR does not add runtime handlers or route registration, startup wiring, protocol routes, Protobuf sources, generated output, repository interface changes, PostgreSQL adapter changes, dependencies, migration changes, authentication/session behavior changes, hosted deployments, release artifacts, public announcements, paid promotion, broad operations/admin behavior, public ACLs, admin search, group/guild/party/room/match storage scopes, batch writes, JSON patch, merge semantics, TTL, script hooks, large object/blob storage, S3-compatible object storage, or direct Nakama/Pitaya API compatibility.

## Alternatives Considered

- Add WebSocket/Protobuf routes together with the service.
- Put the service under `runtime/internal/modules/storage`.
- Let client payloads provide the owner id.
- Treat metadata-only session fields as sufficient owner proof.
- Expose PostgreSQL adapter errors directly to callers.
- Require live PostgreSQL verification for this service slice.
- Copy Nakama storage object public API semantics directly.

## Rationale

The useful next step after the adapter is an application service that can be tested without transport or protocol wiring. This gives later protocol work a behavior surface to call while preserving the route, generated-output, and startup-composition decisions for later gates.

Deriving owner identity from `app.RequestIdentity` keeps the first storage-object posture aligned with the authentication and route-protection boundaries. Rejecting metadata-only identity before repository access prevents accidental cross-player storage access through unauthenticated envelope/session metadata.

Fake-repository tests are sufficient for this slice because persistence mapping was already covered in the PostgreSQL adapter. The runtime behavior risk is identity, validation, conflict mapping, and unit-of-work handoff, not live SQL execution.

## Agent Reasoning Summary

The safest continuation from `W-0208` was an application service only. It provides player-owned storage object behavior using already-ratified repository and adapter boundaries, preserves protocol and generated-output deferrals, and gives the next gate a concrete behavior surface to expose later.

## Decision Weights

```yaml
decision_weights:
  prototype_ready_value: high
  identity_safety: high
  boundary_clarity: high
  agent_readability: high
  implementation_risk: medium
  protocol_risk: deferred
  dependency_expansion: low
confidence: high
```

## Consequences

- `runtime/internal/app/storage/service.go` exists.
- `runtime/internal/app/storage/service_test.go` exists.
- `runtime.storage_objects_runtime_behavior_implementation` becomes the repository check rule for this slice.
- `M-137/W-0209` is completed.
- The next bounded direction is a storage objects protocol route gate.
- Protocol routes, Protobuf source/generated output, startup wiring, migrations, dependencies, authentication/session semantics, hosted deployment, release artifacts, large blob/S3 storage, and direct compatibility remain unchanged by this ADR.

## Reversal Conditions

Revisit this decision if:

- request identity validation semantics change materially;
- the storage repository interface changes materially;
- `PutOwnStorageObject` needs a different first create-vs-replace posture;
- route policy selects a different first protected-route posture for storage objects;
- storage object ownership expands beyond player-owned objects before protocol exposure;
- storage objects become large blob or S3-compatible object storage instead of small JSON game state;
- direct Nakama or Pitaya public API compatibility becomes an explicit project goal.

## Follow-Up

- Define storage object protocol routes and generated output in a later gate before exposing this service over WebSocket/Protobuf.
- Keep startup wiring, public ACLs, admin search, group/guild scopes, batch writes, JSON patch, merge semantics, TTL, script hooks, large blob storage, S3-compatible storage, and direct compatibility behind later gates.
- Preserve metadata-only identity refusal in future handlers.
