# Conversation: Storage Objects Protocol Route Implementation

Date: 2026-05-22
Participants: Maintainer, Agent
Related changes:

- `changes/2026-05-22-implement-storage-objects-protocol-route/`

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
- `decisions/ADR-0119-storage-objects-protocol-route-implementation.md`
- `changes/2026-05-22-implement-storage-objects-protocol-route/`
- `rules/check-rules.json`
- `tools/vibit`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `.arch/modules.yaml`
- `.arch/protocol.yaml`
- `modules/storage/module.yaml`
- `modules/storage/AGENTS.md`
- `modules/storage/AGENTS.zh-CN.md`

## Context

`M-138/W-0210` completed the storage objects protocol route gate. It authorized a bounded route implementation for the own-player storage object route family while preserving repository interface changes, PostgreSQL adapter changes, migrations, dependencies, authentication/session changes, broad storage features, hosted deployment, release artifact expansion, paid promotion, public announcements, and direct Nakama/Pitaya API compatibility as deferred concerns.

The next-ready work item was `W-0211 Implement storage objects protocol route`.

## Maintainer Narrative

The maintainer asked:

```text
继续推进。注意要贴合nakama pitaya 进行提交和推送。
```

English summary: continue the work, keep the Nakama/Pitaya alignment explicit, then commit and push the result.

## Agent Response Summary

The agent advanced one bounded work item and implemented the own-player storage object protocol route family.

The work added:

- `vibit.storage.v1` Protobuf request and response messages;
- generated Go Protobuf output through Buf;
- route keys for `storage.GetOwnStorageObject`, `storage.ListOwnStorageObjects`, `storage.PutOwnStorageObject`, and `storage.DeleteOwnStorageObject`;
- protocol bridge mapping between Protobuf payloads and application storage payloads;
- bootstrap route handlers that inject validated request identity and map redacted storage errors;
- PostgreSQL startup registration for the storage route family;
- transaction bypass for storage write routes because the storage service owns its unit of work;
- focused route, bridge, protected-route, and startup tests;
- ADR, change spec, manifest, check-rule, and continuation updates.

## Decisions

- Complete `M-139/W-0211`.
- Accept `ADR-0119`.
- Add `runtime.storage_objects_protocol_route_implementation`.
- Keep storage object routes own-player only.
- Keep owner identity out of the Protobuf payload and derive it from validated request identity.
- Keep public ACLs, cross-owner access, batch writes, JSON patch/merge, TTL, script hooks, blob/S3 storage, repository/adapter/migration/dependency changes, authentication/session changes, and direct compatibility deferred.
- Select `W-0212 Prove storage objects protocol route in local alpha request flow` as the next bounded direction.

## Nakama And Pitaya Reference Basis

Nakama guided the capability pressure: player-owned storage objects are a common durable game-state primitive, and the route family covers read/list/write/delete behavior with collection/key/value/version semantics.

Pitaya guided the layering pressure: route registration, authenticated request/session context, protocol serialization, application handlers, backend service behavior, and persistence remain separated.

vibit adapts those lessons into its own WebSocket/Protobuf route model and application-owned service boundary. This slice does not add direct Nakama or Pitaya public API compatibility.

## Artifacts

- `proto/vibit/storage/v1/storage.proto`
- `runtime/internal/generated/proto/vibit/storage/v1/storage.pb.go`
- `runtime/internal/app/storage/routes.go`
- `runtime/internal/app/bootstrap/storage.go`
- `runtime/internal/app/bootstrap/storage_test.go`
- `runtime/internal/platform/protocol/protobuf/storage_bridge.go`
- `runtime/internal/platform/protocol/protobuf/storage_bridge_test.go`
- `runtime/cmd/vibit-server/main.go`
- `runtime/cmd/vibit-server/main_test.go`
- `decisions/ADR-0119-storage-objects-protocol-route-implementation.md`
- `changes/2026-05-22-implement-storage-objects-protocol-route/`
- `conversations/2026-05-22-storage-objects-protocol-route-implementation.md`
- `rules/check-rules.json`
- `tools/vibit`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `.arch/modules.yaml`
- `.arch/protocol.yaml`
- `modules/storage/module.yaml`
- `modules/storage/AGENTS.md`
- `modules/storage/AGENTS.zh-CN.md`

## Open Questions

- The local alpha request flow still needs a focused proof for storage object put/get/list/delete.
- Public permissions, ACLs, admin search, group/guild/party/room/match scopes, batch writes, JSON patch, merge semantics, TTL, script hooks, blob/S3 storage, and direct compatibility remain deferred.
- Live PostgreSQL verification is not required by this route implementation slice; the next proof may decide whether local flow coverage needs an opt-in persistent path.

## Follow-Up

- Implement `W-0212 Prove storage objects protocol route in local alpha request flow`.
- Preserve the current storage route Protobuf shape unless the proof reveals a concrete issue.
- Keep Nakama/Pitaya alignment explicit as capability and layering guidance, not direct API compatibility.

## Redaction Notes

No secrets, raw access tokens, generated credentials, digest bytes, HMAC input bytes, verifier keys, private account data, cookies, query tokens, WebSocket subprotocol token material, GitHub tokens, DSNs with credentials, or raw storage object values from a real user are recorded in this conversation log.
