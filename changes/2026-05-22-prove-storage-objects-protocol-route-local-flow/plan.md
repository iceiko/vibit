# Plan

1. Extend the authenticated gameplay E2E fixture with storage route registration.
2. Add a test-only in-memory storage repository that implements the existing `runtime/internal/modules/storage.Repository` interface.
3. Add `TestStorageObjectsProtocolRouteLocalAlphaFlow` to exercise own-player put/get/list/delete over the existing `FrameHandler` path.
4. Update `examples/local-alpha-request-loop.sh` to run the existing gameplay proof and the storage proof.
5. Update `examples/README.md` and `examples/README.zh-CN.md` with prerequisites, command shape, redaction expectations, and Nakama/Pitaya alignment.
6. Record `ADR-0120`, the conversation log, change metadata, manifests, and check-rule updates.
7. Add `W-0213` as the next direction confirmation work item after storage object local proof.
8. Run focused Go tests, repository checks, full Go tests, local alpha request-loop script, and diff hygiene checks.

## Generated Artifacts

None.

## Handwritten Logic

- `runtime/internal/platform/protocol/protobuf/authenticated_gameplay_e2e_test.go`
- `examples/local-alpha-request-loop.sh`
- `examples/README.md`
- `examples/README.zh-CN.md`
- manifest, ADR, conversation, and check-rule metadata

## Verification Commands

- `node -c tools/vibit`
- `node tools/vibit inspect next --json`
- `node tools/vibit inspect rule runtime.storage_objects_protocol_route_local_proof`
- `node tools/vibit check change prove-storage-objects-protocol-route-local-flow --json`
- `node tools/vibit check module storage --json`
- `node tools/vibit check work --json`
- `node tools/vibit check protocol --json`
- `node tools/vibit check runtime --json`
- `node tools/vibit check all --json`
- `git diff --check`
- `cd runtime && go test ./internal/platform/protocol/protobuf -run 'TestAuthenticatedGameplayE2EUsesExistingOnboardingLoginBindingInventoryPresenceAndLogout|TestStorageObjectsProtocolRouteLocalAlphaFlow' -v`
- `cd runtime && go test ./...`
- `examples/local-alpha-request-loop.sh`

## Rollback Or Migration Notes

Reversal means removing the E2E proof test, local-alpha script update, README updates, ADR-0120, change spec, conversation log, check-rule updates, and W-0212 manifest completion state, then reopening `M-140/W-0212`.

No database rollback is required because this change adds no migrations and does not change the `storage_objects` table.
