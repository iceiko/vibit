# Conversation: Device Credential Login Service Behavior Gate

Date: 2026-05-16
Participants: Maintainer, Agent
Related changes:

- `changes/2026-05-16-define-device-credential-login-service-behavior-gate/`

Related artifacts:

- `docs/device-credential-login-service-behavior-gate.md`
- `docs/device-credential-login-service-behavior-gate.zh-CN.md`
- `decisions/ADR-0051-device-credential-login-service-behavior-gate.md`
- `.arch/work-items.yaml`

## Context

The authentication service skeleton was completed in `W-0107`. The next ready work item was `W-0108`, defining the bounded gate for future device credential login behavior before any real login execution is added.

## Maintainer Narrative

The maintainer asked the agent to continue.

The maintainer's standing workflow preference is that `继续` advances one next-ready work item unless a real maintainer decision is required. Routine technical choices should be handled by the agent, while significant product or architecture branches should be made explicit.

## Agent Response Summary

The agent advanced `W-0108` as a gate-only standard.

The change defines the future device credential login sequence, proof input shape, future service dependencies, repository handoff, player account activity check, token issuance posture, public error collapse, redaction requirements, and required tests. It keeps the current service skeleton fail-closed and does not execute login, issue tokens, validate access tokens, expose protocol carriers, wire startup, change repositories or migrations, add dependencies, or add production authentication behavior.

## Decisions

- Treat first device credential proof as server-issued high-entropy Base64URL unpadded material, not a raw device id.
- Require proof pre-validation before any unit-of-work or repository call.
- Use existing authentication and player repositories through unit-of-work capabilities instead of changing repository interfaces in this gate.
- Generate access-token material only after credential proof and player account state are accepted.
- Return raw access-token text only after token storage and unit-of-work success.
- Collapse lookup miss, verifier mismatch, inactive credential, wrong key/version, and inactive player into the same public invalid-credential family.

## Artifacts

- `docs/device-credential-login-service-behavior-gate.md`
- `docs/device-credential-login-service-behavior-gate.zh-CN.md`
- `decisions/ADR-0051-device-credential-login-service-behavior-gate.md`
- `changes/2026-05-16-define-device-credential-login-service-behavior-gate/`

## Open Questions

- None for this gate-only change.

## Follow-Up

- The next work item may implement device credential login inside `runtime/internal/app/authentication/service.go` only if it preserves the gate sequence and deferrals.
- Access-token validation, logout, refresh, cleanup, protocol carriers, startup wiring, repository changes, migrations, dependencies, and production behavior remain separate work.

## Redaction Notes

No secrets, tokens, generated credentials, digest bytes, HMAC input bytes, verifier keys, account details, or private data are recorded in this conversation log.
