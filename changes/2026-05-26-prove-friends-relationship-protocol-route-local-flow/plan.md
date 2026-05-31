# Plan

1. Extend the authenticated gameplay E2E fixture with friends route registration and multi-player deterministic fixture support.
2. Add a test-only in-memory friends relationship repository that implements the existing `runtime/internal/modules/friends.Repository` interface.
3. Add `TestFriendsRelationshipProtocolRouteLocalAlphaFlow` to exercise send/status/accept/list/remove/block/unblock/resend/reject over the existing `FrameHandler` path.
4. Update `examples/local-alpha-example-client.sh` and `examples/local-alpha-request-loop.sh` to include the friends proof.
5. Update examples documentation in English and Chinese with prerequisites, command shape, route list, redaction expectations, and Nakama/Pitaya alignment.
6. Record `ADR-0150`, the conversation log, change metadata, manifests, and check-rule updates.
7. Add `W-0243` as the next direction-selection work item after friends relationship route local proof.
8. Run focused Go tests, repository checks, full Go tests, local alpha scripts, and diff hygiene checks.

## Generated Artifacts

None.

## Handwritten Logic

- `runtime/internal/platform/protocol/protobuf/authenticated_gameplay_e2e_test.go`
- `examples/local-alpha-example-client.sh`
- `examples/local-alpha-request-loop.sh`
- `examples/README.md`
- `examples/README.zh-CN.md`
- `examples/local-alpha-client/README.md`
- `examples/local-alpha-client/README.zh-CN.md`
- manifest, ADR, conversation, and check-rule metadata

## Verification Commands

- `node -c tools/vibit`
- `node tools/vibit inspect next --json`
- `node tools/vibit inspect rule runtime.friends_relationship_protocol_route_local_proof`
- `node tools/vibit check change prove-friends-relationship-protocol-route-local-flow --json`
- `node tools/vibit check module friends --json`
- `node tools/vibit check work --json`
- `node tools/vibit check protocol --json`
- `node tools/vibit check generated --json`
- `node tools/vibit check runtime --json`
- `node tools/vibit check schemas --json`
- `node tools/vibit check memory --json`
- `node tools/vibit check all --json`
- `git diff --check`
- `cd runtime && go test ./internal/platform/protocol/protobuf -run 'TestAuthenticatedGameplayE2EUsesExistingOnboardingLoginBindingInventoryPresenceAndLogout|TestStorageObjectsProtocolRouteLocalAlphaFlow|TestFriendsRelationshipProtocolRouteLocalAlphaFlow|TestPresenceStatusLocalAlphaFlowReportsOfflineAfterCloseAndInvalidation|TestAuthenticatedGameplayFailurePathsLocalAlphaFlow' -v`
- `cd runtime && go test ./...`
- `examples/local-alpha-example-client.sh`
- `examples/local-alpha-request-loop.sh`

## Rollback Or Migration Notes

Reversal means removing the E2E proof test, local-alpha script updates, README updates, ADR-0150, change spec, conversation log, check-rule updates, and W-0242 manifest completion state, then reopening `M-170/W-0242`.

No database rollback is required because this change adds no migrations and does not change the `friend_relationships` table.
