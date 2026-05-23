# ADR-0120: Storage Objects Protocol Route Local Proof

Status: Accepted
Date: 2026-05-23
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-22-prove-storage-objects-protocol-route-local-flow/`

Related conversations:

- `conversations/2026-05-23-storage-objects-protocol-route-local-proof.md`

Related artifacts:

- `runtime/internal/platform/protocol/protobuf/authenticated_gameplay_e2e_test.go`
- `examples/local-alpha-request-loop.sh`
- `examples/README.md`
- `examples/README.zh-CN.md`
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

`M-139/W-0211` implemented the own-player storage object protocol route family. It added `vibit.storage.v1` Protobuf payloads, route keys, protocol bridge mapping, application bootstrap handlers, startup registration, and focused route tests for `storage.GetOwnStorageObject`, `storage.ListOwnStorageObjects`, `storage.PutOwnStorageObject`, and `storage.DeleteOwnStorageObject`.

The next useful step is not new storage behavior. The next useful step is to prove that the route family works inside the same local alpha request flow that already proves onboarding, device credential login, first-message connection binding, protected inventory, protected presence, logout, and revoked-token rejection.

Nakama provides the capability reference: durable player-owned storage objects with collection/key/value/version semantics and read/list/write/delete operations are a common game backend primitive. Pitaya provides the layering reference: transport, session metadata, route handling, serializer/protocol adapter, application handler, service behavior, and repository handoff should remain separate. vibit adapts both references without direct public API compatibility.

## Decision

Add a proof-only local alpha E2E slice for the completed storage object route family.

The proof:

- extends the existing authenticated gameplay E2E fixture with storage route registration;
- adds a test-only in-memory implementation of the existing `runtime/internal/modules/storage.Repository` interface;
- adds `TestStorageObjectsProtocolRouteLocalAlphaFlow`;
- exercises local onboarding, device credential login, first-message connection binding, authenticated `PutOwnStorageObject`, authenticated `GetOwnStorageObject`, authenticated `ListOwnStorageObjects`, authenticated `DeleteOwnStorageObject`, and a post-delete not-found response;
- checks error redaction against the access token, one-time device credential, and proof value text;
- updates `examples/local-alpha-request-loop.sh` to run the existing gameplay proof and the storage route proof;
- updates examples documentation with the exact command, request-flow shape, redaction expectations, and Nakama/Pitaya alignment notes;
- registers `runtime.storage_objects_protocol_route_local_proof` as the check rule for this slice.

This ADR does not add protocol messages, protocol routes, generated output, storage service behavior, repository interface changes, PostgreSQL adapter changes, migration changes, dependency additions, authentication/session behavior changes, route-protection semantic changes, hosted deployments, release artifacts, public announcements, paid promotion, public ACLs, cross-owner access, admin search, group/guild/party/room/match storage scopes, batch writes, JSON patch, merge semantics, TTL, script hooks, large object/blob storage, S3-compatible object storage, production memory storage behavior, or direct Nakama/Pitaya API compatibility.

No direct Nakama or Pitaya API compatibility is added by this proof.

## Alternatives Considered

- Prove only one storage route instead of the full put/get/list/delete family.
- Add a standalone example client before proving route behavior through Go E2E tests.
- Use the PostgreSQL adapter for this proof by requiring a live database.
- Add production in-memory storage startup behavior.
- Combine the local proof with new storage features such as ACLs, batch writes, or group scopes.

## Rationale

The existing local alpha request-loop script is already the most reliable source-first proof path for new contributors. Extending it with storage object put/get/list/delete demonstrates the newly implemented Nakama-class storage capability without creating a new runtime surface.

Keeping the repository implementation test-only avoids creating an accidental second production storage backend. The real production storage path remains the existing PostgreSQL adapter and startup composition.

The proof is valuable because it validates the Pitaya-style separation already chosen by vibit: transport carries bytes, Protobuf adapts payloads, route protection validates request proof, application handlers inject validated identity, the storage service owns behavior, and the repository boundary owns persistence handoff.

## Agent Reasoning Summary

The smallest product-useful continuation after route implementation is an end-to-end local alpha proof. It makes the feature demonstrable to external contributors while preserving vibit's bounded workflow and avoiding premature expansion into ACLs, batch APIs, blob storage, or direct compatibility.

## Decision Weights

```yaml
decision_weights:
  prototype_ready_value: high
  reference_alignment: high
  identity_safety: high
  protocol_compatibility_control: high
  boundary_clarity: high
  implementation_risk: low
  generated_output_risk: none
confidence: high
```

## Consequences

- `TestStorageObjectsProtocolRouteLocalAlphaFlow` proves storage object put/get/list/delete over the existing `FrameHandler` path.
- `examples/local-alpha-request-loop.sh` now runs both the authenticated gameplay proof and the storage object route proof.
- `examples/README.md` and `examples/README.zh-CN.md` document the storage proof.
- `runtime.storage_objects_protocol_route_local_proof` becomes the repository check rule for this slice.
- `M-140/W-0212` is completed.
- The next bounded direction is a confirmation step after storage object local proof.
- Broader storage features, direct compatibility, production memory storage behavior, deployment, release artifact expansion, and larger product modules remain deferred.

## Reversal Conditions

Revisit this decision if:

- the storage route family changes protocol shape;
- route protection stops using request-level authenticated wrappers for storage routes;
- storage service owner derivation changes away from validated request identity;
- the local proof needs a live PostgreSQL database to catch meaningful behavior;
- direct Nakama or Pitaya public API compatibility becomes an explicit future goal through a separate ADR;
- production memory storage behavior becomes an explicit product requirement.

## Follow-Up

- Confirm the next alpha direction after storage object local proof.
- Keep public ACLs, cross-owner reads/writes, group or match storage scopes, batch operations, JSON patch/merge, TTL, script hooks, blob/S3 storage, repository/adapter/migration changes, authentication/session changes, production memory storage behavior, and direct compatibility behind later bounded work items.
