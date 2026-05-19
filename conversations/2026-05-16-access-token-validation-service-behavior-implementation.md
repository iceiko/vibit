# Conversation: Access Token Validation Service Behavior Implementation

Date: 2026-05-16
Participants: Maintainer, Agent
Related changes:

- `changes/2026-05-16-implement-access-token-validation-service-behavior/`

Related artifacts:

- `runtime/internal/app/authentication/service.go`
- `runtime/internal/app/authentication/service_test.go`
- `.arch/work-items.yaml`

## Context

`W-0110` defined the access-token validation service behavior gate and left `ValidateAccessToken` fail-closed. The next ready work item was `W-0111`, implementing the service-local validation behavior.

## Maintainer Narrative

The maintainer asked the agent to continue.

The standing workflow preference remains: `继续` advances one `next_ready` work item unless a real maintainer decision is required. Routine technical implementation details inside an already ratified direction should be handled by the agent.

## Agent Response Summary

The agent advanced `W-0111` as the bounded service-local implementation slice.

The implementation validates opaque Base64URL unpadded access-token proof, computes token lookup and verifier digests, uses existing repositories through unit-of-work capabilities, checks token lifecycle posture and audience, requires active player account state, and returns application-owned request identity only after unit-of-work success. The implementation keeps `SessionValidated` false and does not expose protocol carriers, route protection, session persistence, startup wiring, repository changes, migrations, generated files, dependencies, or broader production behavior.

## Decisions

- Keep validation inside `runtime/internal/app/authentication`.
- Preserve the existing service method signature and existing repository interfaces.
- Treat missing or malformed proof as a pre-unit-of-work failure.
- Treat lookup miss, wrong token posture, wrong audience, verifier mismatch, and inactive player as public `AUTHENTICATION_TOKEN_INVALID`.
- Treat unit-of-work and repository construction failures as public `AUTHENTICATION_TOKEN_STORE_UNAVAILABLE`.
- Keep internal proof statuses specific enough for service tests while preserving public error collapse.
- Open the next workflow state as a blocked confirmation gate rather than implicitly choosing route protection, protocol carriers, session persistence, logout, refresh, or game-domain expansion.

## Artifacts

- `runtime/internal/app/authentication/service.go`
- `runtime/internal/app/authentication/service_test.go`
- `changes/2026-05-16-implement-access-token-validation-service-behavior/`

## Open Questions

- Which next major direction should follow service-local login and token validation: protocol carrier exposure, route protection, session persistence, logout/refresh/cleanup, runtime composition, or core game backend modules.

## Follow-Up

- The next step is a maintainer-confirmed direction gate before crossing protocol, session, route-protection, startup wiring, or broader production behavior boundaries.

## Redaction Notes

No secrets, raw access tokens, generated credentials, digest bytes, HMAC input bytes, verifier keys, private account data, or GitHub tokens are recorded in this conversation log.
