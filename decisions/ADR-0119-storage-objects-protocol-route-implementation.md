# ADR-0119: Storage Objects Protocol Route Implementation

Status: Accepted
Date: 2026-05-22
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-22-implement-storage-objects-protocol-route/`

Related conversations:

- `conversations/2026-05-22-storage-objects-protocol-route-implementation.md`

Related artifacts:

- `proto/vibit/storage/v1/storage.proto`
- `runtime/internal/generated/proto/vibit/storage/v1/storage.pb.go`
- `runtime/internal/app/storage/routes.go`
- `runtime/internal/app/bootstrap/storage.go`
- `runtime/internal/app/bootstrap/storage_test.go`
- `runtime/internal/platform/protocol/protobuf/storage_bridge.go`
- `runtime/internal/platform/protocol/protobuf/storage_bridge_test.go`
- `runtime/cmd/vibit-server/main.go`
- `runtime/cmd/vibit-server/main_test.go`
- `modules/storage/module.yaml`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `.arch/modules.yaml`
- `.arch/protocol.yaml`
- `tools/vibit`
- `rules/check-rules.json`

## Context

`M-138/W-0210` accepted `ADR-0118` and defined the storage objects protocol route gate. The gate authorized a later bounded implementation of own-player storage object routes over the existing WebSocket/Protobuf request flow.

The application service already exists under `runtime/internal/app/storage`. It owns storage object behavior, derives player ownership from validated request identity, refuses metadata-only identity before repository access, and uses the storage repository through the application unit-of-work boundary.

The useful next step is to expose that already-ratified behavior through explicit route names and Protobuf request/response payloads while keeping transport, protocol adaptation, application handlers, runtime service behavior, and persistence separated.

Nakama provides the capability pressure: durable player-owned storage objects with collection/key/value/version semantics and read/list/write/delete operations are a core game backend feature. Pitaya provides the layering pressure: route handling, session/identity context, serializer/protocol adapter, and backend service behavior should remain separate. vibit adapts both references without direct public API compatibility.

## Decision

Implement the own-player storage object protocol route family authorized by `ADR-0118`.

The implementation adds:

- Protobuf source `proto/vibit/storage/v1/storage.proto`;
- generated Go Protobuf output `runtime/internal/generated/proto/vibit/storage/v1/storage.pb.go` through Buf;
- route keys under `runtime/internal/app/storage/routes.go`;
- routes `storage.GetOwnStorageObject`, `storage.ListOwnStorageObjects`, `storage.PutOwnStorageObject`, and `storage.DeleteOwnStorageObject`;
- protocol bridge mapping under `runtime/internal/platform/protocol/protobuf/storage_bridge.go`;
- payload registry integration for storage request and response messages;
- bootstrap handlers under `runtime/internal/app/bootstrap/storage.go`;
- PostgreSQL startup composition that constructs the storage service and registers route handlers;
- transaction bypass for storage write routes because the storage service owns its unit of work;
- focused tests for route mapping, optional expected versions, response mapping, request identity handoff, redacted handler errors, protected-route wrapper enforcement, and startup object-id generation;
- repository check coverage through `runtime.storage_objects_protocol_route_implementation`.

The Protobuf payloads intentionally omit client-supplied owner ids, player ids, session ids, access-token fields, digest material, blob fields, and object-storage bucket fields. Owner identity continues to come from the authenticated request identity passed into the application service.

This ADR does not add public ACLs, cross-owner access, admin search, group/guild/party/room/match storage scopes, batch writes, JSON patch, merge semantics, TTL, script hooks, large object/blob storage, S3-compatible object storage, repository interface changes, PostgreSQL adapter changes, migration changes, dependency additions, authentication/session behavior changes, hosted deployments, release artifacts, public announcements, paid promotion, or direct Nakama/Pitaya API compatibility.

## Alternatives Considered

- Keep storage object behavior application-only until a larger client example exists.
- Add protocol routes but defer Protobuf source and generated output.
- Let client payloads include owner identity.
- Put storage route behavior inside the WebSocket transport package.
- Add Nakama-like ACL and cross-owner behavior in the first route slice.
- Combine this route implementation with a local end-to-end proof in the same slice.

## Rationale

The route gate already selected the safe surface: own-player get, list, put, and delete. Implementing that route family now makes storage objects visible through the same WebSocket/Protobuf architecture used by existing gameplay routes without broadening product scope.

Keeping owner identity out of the payload is the important safety property. The authenticated wrapper and validated request identity remain the only owner source. That preserves the earlier authentication/session boundaries and prevents metadata-only envelope fields from becoming proof.

Separating the route implementation from the next local proof keeps this slice focused on protocol shape, bridge mapping, handler registration, startup composition, and focused unit tests. The next work item can then exercise the route family through the local alpha request flow without changing protocol shape.

## Agent Reasoning Summary

The smallest product-useful continuation after the route gate is to wire the existing storage object service into the established route/protocol/bootstrap layers. This advances Nakama-class storage capability while preserving Pitaya-style layering and vibit's own generated-output and identity-safety rules.

## Decision Weights

```yaml
decision_weights:
  prototype_ready_value: high
  reference_alignment: high
  identity_safety: high
  protocol_compatibility_control: high
  boundary_clarity: high
  implementation_risk: medium
  generated_output_risk: medium
confidence: high
```

## Consequences

- `proto/vibit/storage/v1/storage.proto` exists.
- `runtime/internal/generated/proto/vibit/storage/v1/storage.pb.go` exists and is generated through Buf.
- `runtime/internal/app/storage/routes.go` exposes the storage route keys.
- `runtime/internal/app/bootstrap/storage.go` registers application handlers.
- `runtime/internal/platform/protocol/protobuf/storage_bridge.go` maps route payloads.
- PostgreSQL startup composition registers the storage route family.
- `runtime.storage_objects_protocol_route_implementation` becomes the repository check rule for this slice.
- `M-139/W-0211` is completed.
- The next bounded direction is `W-0212 Prove storage objects protocol route in local alpha request flow`.
- Broader storage features, direct compatibility, production deployment, release artifact expansion, and social/realtime/matchmaking modules remain deferred.

## Reversal Conditions

Revisit this decision if:

- route policy for protected gameplay requests changes materially;
- request identity validation semantics change materially;
- storage objects expand beyond own-player scope before the first local proof;
- generated-output standards change the Protobuf source or generated Go output path;
- direct Nakama or Pitaya public API compatibility becomes an explicit future goal through a separate ADR;
- the local proof reveals that the route surface cannot support the intended alpha request flow without changing protocol shape.

## Follow-Up

- Add a narrow local proof for own-player storage object put/get/list/delete in `W-0212`.
- Preserve the completed Protobuf shape unless the proof finds a specific compatibility issue.
- Keep public ACLs, cross-owner reads/writes, group or match storage scopes, batch operations, JSON patch/merge, TTL, script hooks, blob/S3 storage, repository/adapter/migration changes, authentication/session changes, and direct compatibility behind later bounded work items.
