# Plan

## Files To Create

- None expected beyond this change spec.

## Files To Edit

- `proto/vibit/authentication/v1/authentication.proto`
- `runtime/internal/platform/protocol/protobuf/authentication_bridge.go`
- `runtime/internal/platform/protocol/protobuf/authentication_bridge_test.go`
- `runtime/internal/platform/protocol/protobuf/authenticated_request_test.go`
- `runtime/internal/app/route_authentication.go`
- `runtime/internal/app/route_authentication_test.go`
- `runtime/internal/app/bootstrap/authentication.go`
- `runtime/internal/app/bootstrap/authentication_test.go`
- `runtime/internal/app/transactional_dispatch_test.go`
- `runtime/cmd/vibit-server/main.go`
- `runtime/cmd/vibit-server/main_test.go`
- architecture manifests, module manifests, AGENTS guides, `tools/vibit`, and `rules/check-rules.json`

## Generated Artifacts

- `runtime/internal/generated/proto/vibit/authentication/v1/authentication.pb.go`

Generate through:

```bash
buf generate
```

## Handwritten Logic

- Add `LogoutAccessTokenRoute()`.
- Classify logout as an explicit service-validated token lifecycle route that is allowed through route protection without an `AuthenticatedRequest` wrapper.
- Register and handle the logout route in authentication bootstrap.
- Map logout request/response payloads in the Protobuf bridge.
- Add logout to PostgreSQL startup route registration and transaction bypass.

## Tests

- Add protocol bridge tests for logout request and response mapping.
- Add bootstrap tests for registration, service call mapping, malformed payload, and redacted token errors.
- Add route policy tests for public/service-validated logout posture.
- Extend transactional dispatcher tests for both authentication lifecycle command routes.
- Add frame-handler test proving logout is not wrapped and does not call the request-token validator.
- Extend startup composition tests for logout route policy.

## Verification Commands

- `node -c tools/vibit`
- `buf generate`
- `cd runtime && go test ./...`
- `node tools/vibit check schemas --json`
- `node tools/vibit check generated --json`
- `node tools/vibit check protocol --json`
- `node tools/vibit check runtime --json`
- `node tools/vibit check work --json`
- `node tools/vibit check change implement-protocol-logout-route --json`
- `node tools/vibit check all --json`
- `git diff --check`

## Rollback Or Migration Notes

No database migration is added by this change. Rolling back this slice means removing the route registration, protocol bridge mapping, new Protobuf messages, and regenerated output while preserving the existing service-level logout behavior.
