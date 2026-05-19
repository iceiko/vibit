# Conversation: Access Token Payload Wrapper Route Protection Implementation

Date: 2026-05-17
Participants: Maintainer, Agent
Related changes:

- `changes/2026-05-17-implement-access-token-payload-wrapper-route-protection/`

Related artifacts:

- `proto/vibit/authentication/v1/authenticated_request.proto`
- `runtime/internal/generated/proto/vibit/authentication/v1/authenticated_request.pb.go`
- `runtime/internal/app/route_authentication.go`
- `runtime/internal/app/authentication/route_validator.go`
- `runtime/internal/platform/protocol/protobuf/frame_handler.go`
- `.arch/work-items.yaml`

## Context

The previous gate `W-0113` selected request-level access-token validation with a Protobuf payload wrapper candidate and explicitly deferred WebSocket handshake authentication, session persistence, startup wiring, repositories, migrations, dependencies, logout, refresh, cleanup, token rotation, and broader production behavior.

## Maintainer Narrative

The maintainer asked to continue progress for five hours.

## Agent Response Summary

The agent implemented the bounded `M-042/W-0114` slice. The change added the `vibit.authentication.v1.AuthenticatedRequest` Protobuf wrapper, generated Go Protobuf output through Buf, application route protection, authentication service validation handoff, Protobuf adapter wrapper handling, and focused tests.

Protected routes now require the wrapper and a validated player identity before domain dispatch. Metadata-only identity is rejected for protected routes. Public device credential login remains an explicit public route. The existing Protobuf envelope is unchanged, and WebSocket transport remains credential-neutral.

## Decisions

- Use `vibit.authentication.v1.AuthenticatedRequest` as the first request-level access-token proof wrapper.
- Keep the envelope route fields as the domain route, not as wrapper route metadata.
- Validate access-token proof before protected domain dispatch.
- Collapse malformed, invalid, and unavailable proof failures to route-level authentication error envelopes.
- Keep `SessionValidated` false until session persistence is separately ratified.
- Move the work queue to `M-043/W-0115`, a blocked next-direction confirmation gate.

## Deferrals

- WebSocket handshake authentication.
- Session persistence.
- Startup wiring.
- Repository interface changes.
- PostgreSQL adapter changes.
- Migrations.
- Dependencies.
- Logout, refresh, cleanup, and token rotation.
- Broader production authentication behavior.

## Artifacts

- `proto/vibit/authentication/v1/authenticated_request.proto`
- `runtime/internal/generated/proto/vibit/authentication/v1/authenticated_request.pb.go`
- `runtime/internal/app/route_authentication.go`
- `runtime/internal/app/authentication/route_validator.go`
- `runtime/internal/platform/protocol/protobuf/envelope.go`
- `runtime/internal/platform/protocol/protobuf/frame_handler.go`
- `runtime/internal/platform/protocol/protobuf/authenticated_request_test.go`
- `runtime/internal/platform/transport/ws/server_test.go`
- `changes/2026-05-17-implement-access-token-payload-wrapper-route-protection/`
- `.arch/work-items.yaml`

## Open Questions

- The next major milestone direction remains intentionally unselected until the maintainer confirms it.
- Session persistence and WebSocket handshake authentication remain future decisions.
- Startup composition for route protection remains deferred.

## Follow-Up

- Resolve `M-043/W-0115` by selecting the next milestone direction.
- Keep future work bounded before adding session persistence, WebSocket handshake authentication, startup wiring, logout, refresh, cleanup, token rotation, repository changes, migrations, dependencies, or broader game backend scope.

## Verification

Verification was recorded in the change spec and included Buf lint/generation, Go tests, generated/protocol/runtime/module/work/change/all repository checks, and `git diff --check`.

## Redaction Notes

No secrets, raw access tokens, generated credentials, digest bytes, HMAC input bytes, verifier keys, private account data, or GitHub tokens are recorded in this conversation log.
