# Conversation: Protocol Logout Route Implementation

Date: 2026-05-18
Participants: Maintainer, Agent
Related changes:

- `changes/2026-05-18-implement-protocol-logout-route/`

Related artifacts:

- `docs/protocol-logout-route-gate.md`
- `decisions/ADR-0079-protocol-logout-route-gate.md`
- `proto/vibit/authentication/v1/authentication.proto`
- `runtime/internal/generated/proto/vibit/authentication/v1/authentication.pb.go`
- `runtime/internal/platform/protocol/protobuf/authentication_bridge.go`
- `runtime/internal/app/bootstrap/authentication.go`
- `runtime/internal/app/route_authentication.go`
- `runtime/cmd/vibit-server/main.go`

## Context

The maintainer clarified that vibit should be planned as a Nakama/Pitaya-class game backend product, covering their common capability families while preserving vibit's agent-native standards.

The previous roadmap and gate work selected runtime lifecycle closure before broad product expansion. `W-0169` defined the client-facing protocol logout route gate, and `W-0170` implemented that bounded route.

## Maintainer Narrative

The maintainer asked to replan the project development route and method around the updated Nakama/Pitaya-class target, then continue development.

## Agent Response Summary

The agent implemented the bounded `runtime.authentication.LogoutAccessToken` protocol route. The route exposes existing application-service logout behavior through the WebSocket Protobuf pipeline without moving credential handling into WebSocket transport and without merging logout with socket close, runtime session revocation, reconnect, or protocol session carriers.

## Reference Review

Nakama is used as the product reference for explicit account/session lifecycle APIs and the distinction between logout/token invalidation and realtime socket disconnect.

Pitaya is used as the architecture reference for separating acceptor/transport concerns, session lifecycle, handler routing, and kick/disconnect style connection management.

vibit adapts those lessons by exposing logout as an application authentication route while leaving concrete socket close handoff to a later lifecycle gate.

## Decisions

- Add logout request and response messages to `proto/vibit/authentication/v1/authentication.proto`.
- Regenerate Go Protobuf output through Buf.
- Map logout request/result payloads in the Protobuf bridge.
- Register `runtime.authentication.LogoutAccessToken` through application bootstrap when the authentication service is composed.
- Treat logout as an explicit service-validated token lifecycle route.
- Reject `AuthenticatedRequest` wrapping for logout so the proof comes only from `LogoutAccessTokenRequest.access_token`.
- Add transaction bypass for logout because the authentication service owns its own unit-of-work.
- Preserve the existing Protobuf envelope and WebSocket transport credential neutrality.
- Keep socket close handoff, runtime session revocation, active connection invalidation, reconnect/epoch behavior, protocol session carriers, memory durable logout behavior, dependencies, and direct Nakama/Pitaya API compatibility deferred.

## Artifacts

- `changes/2026-05-18-implement-protocol-logout-route/`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `modules/authentication/module.yaml`
- `AGENTS.md`
- `AGENTS.zh-CN.md`
- `runtime/AGENTS.md`
- `runtime/AGENTS.zh-CN.md`
- `modules/authentication/AGENTS.md`
- `modules/authentication/AGENTS.zh-CN.md`
- `rules/check-rules.json`
- `tools/vibit`

## Verification

Verified:

- `buf generate`
- `node -c tools/vibit`
- `node tools/vibit check change implement-protocol-logout-route --json`
- `node tools/vibit check runtime --json`
- `cd runtime && go test ./...`

Runtime check passed with the existing `runtime.identity_boundary` warning in `runtime/internal/platform/persistence/postgres/authentication_repository.go`.

## Open Questions

- Whether a later close handoff implementation should close the current socket after logout remains deferred.
- Whether runtime session revocation should be a client route, admin action, or policy side effect remains deferred.
- Whether `token_record_id` should remain visible in protocol logout responses or move to audit/admin surfaces remains a future compatibility question.

## Follow-Up

The next ready work item is `M-099/W-0171 confirm_next_direction_after_protocol_logout_route_implementation`. The recommended direction is `define_transport_close_handoff_gate`, because current close policy can mark active connection records invalidated but cannot yet hand concrete close action to WebSocket transport.

## Redaction Notes

No raw access tokens, device credentials, generated secrets, digest bytes, HMAC input bytes, verifier keys, database secrets, player private data, or GitHub tokens are recorded in this conversation log.
