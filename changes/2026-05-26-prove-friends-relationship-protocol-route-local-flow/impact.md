# Impact Analysis

## Affected Modules

- `runtime`
- `friends`
- `workflow`
- `reference`

## Module Ownership Impact

The friends module continues to own storage-neutral friends relationship repository vocabulary and value types under `runtime/internal/modules/friends`.

The friends application service remains under `runtime/internal/app/friends` and is not changed by this proof. The PostgreSQL friends adapter remains under `runtime/internal/platform/persistence/postgres` and is not changed.

The proof adds a test-only in-memory friends relationship repository inside `runtime/internal/platform/protocol/protobuf/authenticated_gameplay_e2e_test.go` so the local alpha proof can exercise the route family without creating a new production memory social graph backend.

## Public Contract Impact

No public protocol contracts are added or changed.

The proof exercises the already implemented route family:

- command `friends.SendFriendRequest`
- command `friends.AcceptFriendRequest`
- command `friends.RejectFriendRequest`
- command `friends.RemoveFriend`
- command `friends.BlockPlayer`
- command `friends.UnblockPlayer`
- query `friends.ListFriendRelationships`
- query `friends.GetFriendRelationshipStatus`

The payload package remains `vibit.friends.v1`.

## Generated Output Impact

No generated files are added or changed.

The existing generated friends Protobuf output remains traced to `proto/vibit/friends/v1/friends.proto`.

## Data And Migration Impact

No migrations are added or changed.

The proof uses a test-only repository and does not change `friend_relationships` SQL shape, PostgreSQL adapter behavior, transaction behavior, startup migrations, or live database requirements.

## Runtime Impact

The local alpha scripts now run the focused friends route proof with the existing local alpha request-flow tests. They still execute inside the Go test runtime and do not start a hosted deployment or publish release artifacts.

The production runtime route registration and friends service behavior are unchanged.

## Authentication And Session Impact

The friends proof uses the existing local onboarding, device credential login, first-message connection binding, and authenticated request wrapper flow. It does not add handshake authentication, token carriers, session persistence changes, bound-session policy changes, or metadata-only identity trust.

## Nakama And Pitaya Alignment

Nakama informs the capability being proven: player-to-player friend request, accept/reject, list/status, remove, block, and unblock operations are common social graph primitives.

Pitaya informs the layering being proven: transport frame handling, Protobuf adaptation, authenticated route protection, session metadata, application handlers, service behavior, and repository handoff stay separated.

The proof uses vibit route names and payloads. It does not copy Nakama or Pitaya public APIs.

## Test Impact

Adds `TestFriendsRelationshipProtocolRouteLocalAlphaFlow` to the authenticated gameplay E2E test file and extends the fixture with friends route registration, multi-player deterministic fixture support, and a test-only friends repository.

The proof covers:

- local onboarding and login for two players;
- first-message connection binding for both players;
- authenticated friend request send and relationship status;
- authenticated request acceptance and relationship listing;
- authenticated friend removal;
- authenticated block and unblock;
- authenticated second friend request and rejection;
- secret redaction expectations for the final envelope.

## Compatibility Risks

The main risk is that the local proof could drift into a new production friends backend. This is bounded by keeping the repository implementation test-only, not changing production startup beyond the existing W-0241 route registration, and recording the deferral of production memory social graph behavior.
