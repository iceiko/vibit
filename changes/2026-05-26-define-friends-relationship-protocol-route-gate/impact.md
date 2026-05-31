# Impact Analysis

## Affected Modules

- `friends`
- `runtime`
- `workflow`
- `reference`

## Module Ownership Impact

The friends module remains the semantic owner of storage-neutral friends relationship vocabulary and repository types. Application friends behavior remains under `runtime/internal/app/friends`. Future protocol bridge ownership is planned under `runtime/internal/platform/protocol/protobuf`, future route handler ownership under `runtime/internal/app/bootstrap`, and future Protobuf source ownership under `proto/vibit/friends/v1`.

No runtime ownership changes are implemented by this gate.

## Public Contract Impact

No public command, query, event, error, permission, or Protobuf contract is added by this gate.

The gate records candidate future protocol routes:

- `friends.SendFriendRequest`
- `friends.AcceptFriendRequest`
- `friends.RejectFriendRequest`
- `friends.RemoveFriend`
- `friends.BlockPlayer`
- `friends.UnblockPlayer`
- `friends.ListFriendRelationships`
- `friends.GetFriendRelationshipStatus`

## Data And Migration Impact

No migrations or data model changes are added. Existing friends relationship persistence remains the PostgreSQL `friend_relationships` table introduced by `W-0233`.

## Test Impact

No Go tests are required because this is a gate-only documentation and manifest change. Future implementation tests are recorded in the gate standard.

## Documentation Impact

Adds:

- `docs/friends-relationship-protocol-route-gate.md`
- `docs/friends-relationship-protocol-route-gate.zh-CN.md`
- `ADR-0148`
- conversation log
- change spec artifacts

Updates architecture manifests, friends module manifest/guides, repository docs, and continuation pointers.

## Compatibility Risks

This gate reduces compatibility risk by preventing ad hoc friends routes or copied Nakama/Pitaya API shapes before protocol contracts are ratified.

No wire compatibility changes occur because no `.proto` or generated output is added.
