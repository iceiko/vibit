# ADR-0118: Storage Objects Protocol Route Gate

Status: Accepted
Date: 2026-05-22
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-22-define-storage-objects-protocol-route-gate/`

Related conversations:

- `conversations/2026-05-22-storage-objects-protocol-route-gate.md`

Related artifacts:

- `docs/storage-objects-protocol-route-gate.md`
- `docs/storage-objects-protocol-route-gate.zh-CN.md`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `modules/storage/module.yaml`
- `tools/vibit`
- `rules/check-rules.json`

## Context

`M-137/W-0209` implemented application-owned storage object runtime behavior under `runtime/internal/app/storage`. The service can get, list, put, and delete own-player storage objects while deriving owner identity from validated `app.RequestIdentity`, rejecting metadata-only identity before repository access, and using unit-of-work repository handoff.

The next bounded step is to define how this service may later become client-facing over the existing WebSocket/Protobuf protocol. That requires route names, message-shape posture, protected-route policy, protocol adapter ownership, generated-output expectations, redaction expectations, and stop conditions before `.proto` or route code exists.

Nakama is the reference for the product capability: durable player-owned storage objects with collection/key addressing, versioned conflict behavior, list/read/write/delete operations, and JSON game state are a common game backend surface. Pitaya is the reference for layering: acceptors, sessions, route handlers, serializers, and backend services remain separate. vibit adapts both references without copying either public API.

## Decision

Accept `docs/storage-objects-protocol-route-gate.md` as the gate for future storage object protocol routes.

The gate records:

- future route family `storage.GetOwnStorageObject`, `storage.ListOwnStorageObjects`, `storage.PutOwnStorageObject`, and `storage.DeleteOwnStorageObject`;
- query/command split for read and mutation routes;
- future Protobuf source candidate `proto/vibit/storage/v1/storage.proto`;
- future generated Go output candidate `runtime/internal/generated/proto/vibit/storage/v1/storage.pb.go`;
- candidate request/response messages and fields;
- unchanged envelope posture;
- `request_token_required` protected-route posture with authenticated wrapper requirement;
- validated authenticated request identity as the only owner source;
- refusal to treat metadata-only `player_id` or `session_id` as proof;
- no client-supplied owner ids in the first posture;
- future protocol bridge and application handler ownership;
- generated-output and no-hand-edit expectations;
- Nakama/Pitaya reference mapping;
- required future tests;
- stop conditions before route implementation, Protobuf source, generated output, startup wiring, repository/adapter/migration/dependency changes, authentication/session changes, hosted deployment, release artifacts, public announcements, paid promotion, object/blob storage, S3-compatible storage, or direct compatibility.

This ADR does not add protocol route implementation, Protobuf source, generated output, startup wiring, runtime handlers, repository interface changes, PostgreSQL adapter changes, dependencies, migration changes, authentication/session behavior changes, hosted deployments, release artifacts, public announcements, paid promotion, broad operations/admin behavior, public ACLs, admin search, group/guild/party/room/match storage scopes, batch writes, JSON patch, merge semantics, TTL, script hooks, large object/blob storage, S3-compatible object storage, or direct Nakama/Pitaya API compatibility.

## Alternatives Considered

- Add storage `.proto` messages and route handlers in the same slice.
- Copy Nakama storage object route names and permission model directly.
- Use a Pitaya-style opaque route string instead of vibit's `kind/module/name` envelope route identity.
- Allow client payloads to include owner ids.
- Treat envelope `player_id` or `session_id` metadata as sufficient owner proof.
- Put storage route behavior in WebSocket transport or PostgreSQL adapter packages.

## Rationale

The service behavior exists, but exposing it over the client protocol creates compatibility, identity, and generated-output risk. A gate-only slice lets the project preserve the useful Nakama-style storage object capability while keeping Pitaya-style route/session/handler separation and vibit's own contract-first protocol discipline.

The future route family should be player-owned and own-object only because the service is already built around validated player identity. Broader Nakama-like features such as public permissions, cross-user reads, system storage, script hooks, and batch operations are product-useful, but they need separate contracts and checks.

## Agent Reasoning Summary

The safest continuation from `W-0209` is a protocol route gate. It records the future route and message shape, connects the work to Nakama/Pitaya capability mapping, and prevents agents from jumping directly into `.proto`, generated code, startup wiring, or direct compatibility.

## Decision Weights

```yaml
decision_weights:
  prototype_ready_value: high
  reference_alignment: high
  identity_safety: high
  protocol_compatibility_control: high
  boundary_clarity: high
  implementation_risk: low
  generated_output_risk: deferred
confidence: high
```

## Consequences

- `docs/storage-objects-protocol-route-gate.md` and `.zh-CN.md` exist.
- `runtime.storage_objects_protocol_route_gate` becomes the repository check rule for this slice.
- `M-138/W-0210` is completed.
- The next bounded direction is storage objects protocol route implementation.
- Protobuf source/generated output, route implementation, startup wiring, repository/adapter/migration/dependency changes, authentication/session semantics, hosted deployment, release artifacts, large blob/S3 storage, and direct compatibility remain unchanged by this ADR.

## Reversal Conditions

Revisit this decision if:

- route policy changes the first protected-route posture for storage objects;
- request identity validation semantics change materially;
- storage object runtime behavior changes its owner derivation model;
- the project chooses direct Nakama or Pitaya API compatibility through a future ADR;
- storage objects expand to public ACLs, cross-owner access, group/room/match scopes, or large blob/S3-compatible storage before protocol exposure;
- generated-output standards change before the future implementation slice.

## Follow-Up

- Implement storage object protocol routes only after this gate is accepted.
- Keep Protobuf source/generated output, startup wiring, public ACLs, admin search, group/guild scopes, batch writes, JSON patch, merge semantics, TTL, script hooks, large blob storage, S3-compatible storage, and direct compatibility behind later bounded work items.
- Preserve metadata-only identity refusal in future route tests.
