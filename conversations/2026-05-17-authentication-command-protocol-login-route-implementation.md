# Conversation: Authentication Command Protocol Login Route Implementation

Date: 2026-05-17
Participants: Maintainer, Agent
Related changes:

- `changes/2026-05-17-implement-authentication-command-protocol-login-route/`

Related artifacts:

- `proto/vibit/authentication/v1/authentication.proto`
- `runtime/internal/generated/proto/vibit/authentication/v1/authentication.pb.go`
- `runtime/internal/app/bootstrap/authentication.go`
- `runtime/internal/platform/protocol/protobuf/authentication_bridge.go`
- `runtime/cmd/vibit-server/main.go`

## Context

The authentication command protocol and login route gate authorized exposing the existing `AuthenticateWithDeviceCredential` service method through a public Protobuf command route.

## Maintainer Narrative

The maintainer asked in Chinese for the agent to recommend the next ten steps and then execute those ten steps. The maintainer selected `add_authentication_command_protocol_messages_and_login_route_registration` and emphasized using Nakama and Pitaya as important game server references.

## Agent Response Summary

The agent implemented the bounded W-0120 slice: Protobuf login command messages, generated output, protocol bridge behavior, application route handlers, PostgreSQL startup registration, transaction-wrapper bypass, focused tests, architecture manifests, AGENTS guides, change spec, conversation memory, and repository checks.

## Implementation Summary

The implementation:

- Adds `proto/vibit/authentication/v1/authentication.proto`.
- Generates `runtime/internal/generated/proto/vibit/authentication/v1/authentication.pb.go` through Buf.
- Maps login request and response payloads in `runtime/internal/platform/protocol/protobuf/authentication_bridge.go`.
- Registers the public login route in `runtime/internal/app/bootstrap/authentication.go`.
- Registers the route only in PostgreSQL startup composition where the authentication service exists.
- Bypasses the outer transactional dispatcher for `runtime.authentication.AuthenticateWithDeviceCredential`.
- Preserves protected-route validation behavior for non-public gameplay routes.

## Decisions

- Keep the existing Protobuf envelope unchanged.
- Keep WebSocket transport credential-neutral.
- Keep the first composed public login route limited to the PostgreSQL runtime path.
- Keep memory durable login unavailable until a future explicit work item.
- Let the existing authentication service own credential validation, token issuance, repository access, and its own unit-of-work.
- Use Nakama for the authenticate-before-gameplay capability sequence and Pitaya for transport/session/route/handler layering, without adopting direct public API compatibility.

## Artifacts

- `proto/vibit/authentication/v1/authentication.proto`
- `runtime/internal/generated/proto/vibit/authentication/v1/authentication.pb.go`
- `runtime/internal/app/bootstrap/authentication.go`
- `runtime/internal/app/bootstrap/authentication_test.go`
- `runtime/internal/platform/protocol/protobuf/authentication_bridge.go`
- `runtime/internal/platform/protocol/protobuf/authentication_bridge_test.go`
- `runtime/internal/app/transactional_dispatch.go`
- `runtime/cmd/vibit-server/main.go`
- `changes/2026-05-17-implement-authentication-command-protocol-login-route/`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `tools/vibit`
- `rules/check-rules.json`

## Deferrals Preserved

- Session persistence.
- WebSocket handshake authentication.
- HTTP `Authorization`, Bearer, cookie, query-string, and WebSocket subprotocol credential carriers.
- Repository interface changes.
- PostgreSQL adapter changes.
- Migrations.
- Dependencies.
- Logout, refresh, cleanup, token rotation, and token validation audit mutation.
- Memory durable authentication behavior.
- Direct Nakama/Pitaya API compatibility.

## Open Questions

- The next milestone direction is blocked at `M-049/W-0121` pending maintainer selection.
- Session persistence, WebSocket handshake authentication, logout/refresh/cleanup, token rotation, memory durable authentication behavior, operations posture, and broader game backend expansion remain future decisions.

## Follow-Up

- Run full runtime tests, repository checks, memory checks, generated/protocol checks, change checks, and diff checks.
- Record any remaining warnings or unavailable live verification explicitly.

## Redaction Notes

No secrets, raw access tokens, generated credentials, digest bytes, HMAC input bytes, verifier keys, private account data, or GitHub tokens are recorded in this conversation log.
