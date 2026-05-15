# Conversation: Authentication Service Implementation Readiness Gate

Date: 2026-05-15
Participants: Maintainer, Agent
Related changes:

- `changes/2026-05-15-define-authentication-service-implementation-readiness-gate/`

Related artifacts:

- `docs/authentication-service-implementation-readiness-gate.md`
- `decisions/ADR-0044-authentication-service-implementation-readiness-gate.md`
- `.arch/work-items.yaml`

## Context

The maintainer asked the agent to continue the work queue for an extended period. Existing maintainer direction allows bounded professional decisions to proceed without unnecessary confirmation while true architecture forks still require discussion.

The next ready work item was `W-0095 Define authentication service implementation readiness gate`.

## Maintainer Narrative

The maintainer wants vibit to remain self-bootstrapping and controlled before implementation accelerates. The project should do necessary preparation instead of rushing into a demo-shaped authentication implementation.

The maintainer also wants vibit to keep referencing Nakama and Pitaya for game server capability coverage, while preserving vibit's core differentiator: agent-native maintainability through explicit boundaries, manifests, checks, and durable decisions.

## Agent Response Summary

The agent reviewed the runtime authentication boundary, service interface boundary, verifier algorithm boundary, secret configuration boundary, material generation boundary, and verifier digest computation/comparison boundary.

The agent concluded that the next safe step is to consolidate implementation readiness before code starts. The gate defines required prior boundaries, recommended package ownership, forbidden first-slice write areas, the first implementation queue, test classes, redaction expectations, Nakama/Pitaya mapping, and deferrals.

## Decisions

- Future authentication service implementation remains application-owned under `runtime/internal/app`.
- The recommended service package candidate is `runtime/internal/app/authentication`.
- The first code slice must be separately authorized.
- The recommended next implementation gate is local verifier key configuration loading.
- Protocol carriers, WebSocket behavior, repository changes, migrations, major dependencies, and production authentication behavior remain deferred.

## Artifacts

- `docs/authentication-service-implementation-readiness-gate.md`
- `docs/authentication-service-implementation-readiness-gate.zh-CN.md`
- `decisions/ADR-0044-authentication-service-implementation-readiness-gate.md`
- `changes/2026-05-15-define-authentication-service-implementation-readiness-gate/`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `.arch/work-items.yaml`
- `modules/authentication/module.yaml`

## Open Questions

- Exact first implementation helper names remain for the local verifier key configuration loading gate.
- Whether the first key loader reads process environment directly or accepts explicit in-memory configuration first remains for the next code gate.
- Protobuf authentication messages and WebSocket proof carriers remain deferred.

## Follow-Up

- Advance the local verifier key configuration loading gate as the next bounded work item.
- Keep authentication service behavior, token generation, credential generation, digest computation, verifier comparison, login execution, token validation, logout execution, cleanup, Protobuf messages, WebSocket proof carriers, authentication dependencies, repository changes, and migration changes behind later gates.

## Redaction Notes

No real secrets, tokens, credentials, generated material, device identifiers, account details, environment variable values, verifier key values, digest bytes, HMAC inputs, or private data are stored in this conversation log.
