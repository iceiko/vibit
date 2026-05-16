# Conversation: Authentication Service Behavior Skeleton

Date: 2026-05-16
Participants: Maintainer, Agent
Related changes:

- `changes/2026-05-16-implement-authentication-service-behavior-skeleton/`

Related artifacts:

- `runtime/internal/app/authentication/service.go`
- `runtime/internal/app/authentication/service_test.go`
- `.arch/work-items.yaml`

## Context

The authentication service behavior implementation gate was completed in `W-0106`. The next ready work item was `W-0107`, adding a skeleton-only service shape.

## Maintainer Narrative

The maintainer asked the agent to continue development.

The maintainer has established that `继续` means advancing one next-ready work item unless a real maintainer decision is required. The maintainer also wants routine technical choices handled by the agent while keeping major architectural or product branches explicit.

## Agent Response Summary

The agent advanced `W-0107` as a skeleton-only implementation slice.

The change adds an application-owned authentication service shape under `runtime/internal/app/authentication` with typed dependencies, request/result vocabulary, redacted internal failure classes, public error codes, and fail-closed behavior.

The service methods reserve `AuthenticateWithDeviceCredential`, `ValidateAccessToken`, `LogoutAccessToken`, and `RefreshAccessToken`, but they do not execute login, validate tokens, issue tokens, revoke tokens, refresh tokens, call repositories, expose protocol carriers, wire startup, add dependencies, or add production authentication behavior.

## Decisions

- Keep the first service code slice skeleton-only.
- Require a unit-of-work runner dependency at service construction, but do not call it until a later behavior work item authorizes repository orchestration.
- Return redacted `ServiceError` values and typed results instead of leaking proof material or internal verifier details.
- Open the next work item as a gate for device credential login service behavior rather than implementing login directly.

## Artifacts

- `runtime/internal/app/authentication/service.go`
- `runtime/internal/app/authentication/service_test.go`
- `changes/2026-05-16-implement-authentication-service-behavior-skeleton/`
- `conversations/2026-05-16-authentication-service-behavior-skeleton.md`

## Open Questions

- None for this skeleton-only change.

## Follow-Up

- Define the next bounded gate for real device credential login service behavior.
- Keep access-token validation, token issuance, repository orchestration, protocol carriers, startup wiring, dependencies, migrations, and production behavior behind later bounded work items.

## Redaction Notes

No secrets, tokens, generated credentials, digest bytes, HMAC input bytes, verifier keys, account details, or private data are recorded in this conversation log.
