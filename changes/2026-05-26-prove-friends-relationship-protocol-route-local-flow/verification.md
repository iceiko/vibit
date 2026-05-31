# Verification

Verification started on 2026-05-31 for W-0242.

## TDD Evidence

Focused E2E coverage was added in:

- `runtime/internal/platform/protocol/protobuf/authenticated_gameplay_e2e_test.go`

The proof test is:

- `TestFriendsRelationshipProtocolRouteLocalAlphaFlow`

It exercises local onboarding, device credential login, first-message connection binding, authenticated friends relationship send/status/accept/list/remove/block/unblock/resend/reject behavior over the existing WebSocket/Protobuf `FrameHandler` path.

## Commands Run

- `cd runtime && go test ./internal/platform/protocol/protobuf -run TestFriendsRelationshipProtocolRouteLocalAlphaFlow -v`
  - Passed.

## Final Commands

Final verification completed on 2026-05-31:

- `node -c tools/vibit`
  - Passed.
- `node tools/vibit inspect rule runtime.friends_relationship_protocol_route_local_proof`
  - Passed.
- `node tools/vibit inspect next --json`
  - Passed; reported `M-171` as the active milestone and `W-0243` as the next-ready work item.
- `node tools/vibit check change prove-friends-relationship-protocol-route-local-flow --json`
  - Passed: 13 passed, 0 warnings, 0 failures.
- `node tools/vibit check module friends --json`
  - Passed: 23 passed, 0 warnings, 0 failures.
- `node tools/vibit check work --json`
  - Passed: 1470 passed, 0 warnings, 0 failures.
- `node tools/vibit check runtime --json`
  - Passed: 18894 passed, 1 warning, 0 failures.
  - The warning is the existing `runtime.identity_boundary` warning for `runtime/internal/platform/persistence/postgres/authentication_repository.go`.
- `node tools/vibit check protocol --json`
  - Passed: 200 passed, 0 warnings, 0 failures.
- `node tools/vibit check generated --json`
  - Passed: 231 passed, 0 warnings, 0 failures.
- `node tools/vibit check schemas --json`
  - Passed: 4480 passed, 0 warnings, 0 failures.
- `node tools/vibit check memory --json`
  - Passed: 4052 passed, 0 warnings, 0 failures.
- `node tools/vibit check all --json`
  - Passed: 297 subchecks passed, 1 warning, 0 failures.
- `cd runtime && go test ./internal/platform/protocol/protobuf -run 'TestAuthenticatedGameplayE2EUsesExistingOnboardingLoginBindingInventoryPresenceAndLogout|TestStorageObjectsProtocolRouteLocalAlphaFlow|TestFriendsRelationshipProtocolRouteLocalAlphaFlow|TestPresenceStatusLocalAlphaFlowReportsOfflineAfterCloseAndInvalidation|TestAuthenticatedGameplayFailurePathsLocalAlphaFlow' -v`
  - Passed.
- `examples/local-alpha-example-client.sh`
  - Passed.
- `examples/local-alpha-request-loop.sh`
  - Passed.
- `cd runtime && go test ./...`
  - Passed.
- `git diff --check`
  - Passed.

## Not Applicable

- `buf generate`, because this proof adds no Protobuf source and changes no generated output.
- Live PostgreSQL verification, because this slice proves the existing route family through a local test-only repository and does not change SQL mapping, migrations, or PostgreSQL adapters.
- Release artifact verification, because this change creates no binaries, packages, containers, checksums, provenance files, signing artifacts, install scripts, registry publications, SDK packages, hosted deployments, public announcements, or paid promotion.
