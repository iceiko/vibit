# Conversation: Application Authentication Service Interface Boundary

Date: 2026-05-15
Participants: Maintainer, Agent
Related changes:

- `changes/2026-05-15-define-application-authentication-service-interface-boundary/`

Related artifacts:

- `docs/application-authentication-service-interface-boundary.md`
- `decisions/ADR-0039-application-authentication-service-interface-boundary.md`
- `.arch/work-items.yaml`

## Context

The maintainer asked the agent to continue the work queue. Earlier maintainer direction established that routine professional sequencing should be advanced without stopping for unnecessary confirmation, while true architecture forks should be surfaced.

The next ready work item was `W-0090 Define application authentication service interface boundary`.

## Maintainer Narrative

The maintainer wants the project to remain self-bootstrapping and controlled, with enough preparation before implementation so future agents can extend the system without drifting away from the intended architecture.

The maintainer previously emphasized that Nakama and Pitaya should remain active references for comparable game server capability coverage, while vibit's differentiator is agent-native maintainability.

## Agent Response Summary

The agent reviewed the runtime authentication implementation boundary, generated authentication contract shape timing, generated shape output, authentication repository boundary, application unit-of-work code, and work queue.

The agent concluded that the next safe step is a boundary-only standard and ADR for future application authentication service interfaces. No Go service code, authentication behavior, token generation, verifier comparison, protocol messages, WebSocket carriers, dependencies, repository changes, or migration changes should be added in this step.

## Decisions

- Future authentication service interfaces are application-owned under `runtime/internal/app`.
- Generated authentication contract shapes inform service-level request/result vocabulary but remain metadata-only generated output.
- Future service behavior may use `authentication.Repository` only through the application unit-of-work boundary.
- Request identity handoff remains application-owned and must produce validated `RequestIdentity` before production-sensitive domain dispatch.
- Raw credential and raw access token material must be redacted from public errors, logs, audit-safe facts, change specs, and conversation logs.
- Runtime behavior remains deferred.

## Artifacts

- `docs/application-authentication-service-interface-boundary.md`
- `docs/application-authentication-service-interface-boundary.zh-CN.md`
- `decisions/ADR-0039-application-authentication-service-interface-boundary.md`
- `changes/2026-05-15-define-application-authentication-service-interface-boundary/`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/work-items.yaml`
- `modules/authentication/module.yaml`

## Open Questions

- The exact Go package path for future service code remains for a later implementation gate.
- Verifier algorithms and redaction tests still need a dedicated preparation gate before service code.
- Protobuf authentication messages and WebSocket proof carriers remain deferred.

## Follow-Up

- Advance the verifier algorithm and redaction test boundary.
- Keep service code behind a later bounded work item.
- Keep token generation, verifier comparison, login execution, token validation, logout execution, cleanup, Protobuf messages, WebSocket proof carriers, authentication dependencies, repository changes, and migration changes behind later gates.

## Redaction Notes

No secrets, tokens, account details, or private data are stored in this conversation log.
