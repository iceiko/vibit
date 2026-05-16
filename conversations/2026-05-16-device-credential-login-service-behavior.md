# Conversation: Device Credential Login Service Behavior

Date: 2026-05-16
Participants: Maintainer, Agent
Related changes:

- `changes/2026-05-16-implement-device-credential-login-service-behavior/`

Related artifacts:

- `runtime/internal/app/authentication/service.go`
- `runtime/internal/app/authentication/service_test.go`
- `docs/device-credential-login-service-behavior-gate.md`
- `decisions/ADR-0051-device-credential-login-service-behavior-gate.md`
- `.arch/work-items.yaml`

## Artifacts

- Change record: `changes/2026-05-16-implement-device-credential-login-service-behavior/`
- Runtime source: `runtime/internal/app/authentication/service.go`
- Runtime tests: `runtime/internal/app/authentication/service_test.go`
- Work registry: `.arch/work-items.yaml`
- Runtime manifest: `.arch/runtime.yaml`
- Conventions manifest: `.arch/conventions.yaml`
- Contracts manifest: `.arch/contracts.yaml`
- Authentication module manifest: `modules/authentication/module.yaml`
- Check rule catalog: `rules/check-rules.json`
- Check implementation: `tools/vibit`

## Context

`W-0108` defined the gate for future device credential login behavior. The maintainer asked the agent to continue, which advances one next-ready work item.

## Maintainer Narrative

The maintainer asked the agent to continue.

The standing workflow preference remains: `继续` advances one `next_ready` work item unless a real maintainer decision is required. Routine technical implementation details inside an already ratified direction should be handled by the agent.

## Agent Response Summary

The agent advanced `W-0109`.

The implementation keeps login behavior inside `runtime/internal/app/authentication/service.go`. It rejects missing or malformed proof before unit-of-work, uses existing verifier helpers, obtains repositories through unit-of-work capabilities, checks credential and player account state, generates opaque access-token material through an injected reader, stores token digests only, and returns raw access-token text only after token storage and unit-of-work success.

The implementation does not add access-token validation, logout, refresh, cleanup, protocol carriers, startup wiring, repository interface changes, migrations, generated files, dependencies, or broader production authentication behavior.

## Decisions

- Extend `ServiceDependencies` with service-local login dependencies instead of adding startup wiring.
- Use a local unit-of-work capability interface instead of expanding the global `tx.UnitOfWork` interface.
- Treat all credential lookup, verifier mismatch, invalid credential posture, and inactive player states as public invalid credential failures.
- Open the next work item as an access-token validation behavior gate rather than implementing validation immediately.

## Redaction Notes

No secrets, tokens, generated credential values, digest bytes, HMAC input bytes, verifier keys, account private data, or GitHub tokens are recorded in this conversation log.

## Open Questions

- Access-token validation behavior is intentionally deferred to `W-0110`.
- Protocol carriers, WebSocket handshake authentication, startup wiring, logout, refresh, cleanup, migrations, and broader production authentication behavior remain intentionally deferred.

## Follow-Up

- Continue with `W-0110`, the access-token validation service behavior gate, after `W-0109` is verified, committed, and pushed or reported as locally ahead if network push fails.
