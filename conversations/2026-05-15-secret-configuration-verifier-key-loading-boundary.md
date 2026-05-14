# Conversation: Secret Configuration And Verifier Key Loading Boundary

Date: 2026-05-15
Participants: Maintainer, Agent
Related changes:

- `changes/2026-05-15-define-secret-configuration-verifier-key-loading-boundary/`

Related artifacts:

- `docs/secret-configuration-verifier-key-loading-boundary.md`
- `decisions/ADR-0041-secret-configuration-verifier-key-loading-boundary.md`
- `.arch/work-items.yaml`

## Context

The maintainer asked the agent to continue the work queue for an extended period. Existing maintainer direction allows routine bounded work to proceed without unnecessary confirmation, while true architecture forks still require discussion.

The next ready work item was `W-0092 Define secret configuration and verifier key loading boundary`.

## Maintainer Narrative

The maintainer wants vibit to stay self-bootstrapping and controlled, with enough preparation that future agents can implement authentication without drifting into ad hoc code, hidden dependencies, or misplaced responsibilities.

The maintainer previously emphasized that Nakama and Pitaya should remain active references for game server capability planning, while vibit's differentiator is agent-native maintainability through explicit boundaries, manifests, checks, and durable decisions.

## Agent Response Summary

The agent reviewed the verifier algorithm/redaction boundary, application authentication service interface boundary, runtime authentication implementation boundary, authentication module guide, and current work queue.

The agent concluded that the next safe step is to ratify secret configuration and verifier key loading posture without adding code. The selected posture allows future local implementation to use process environment configuration or explicit runtime secret input, keeps KMS and cloud secret-manager integration behind later dependency and operations gates, requires four separated logical verifier keys, defines key identifier redaction rules, and requires fail-closed production behavior.

## Decisions

- Future verifier key loading is application-owned under `runtime/internal/app`.
- The first local implementation may use process environment configuration or explicit runtime secret input after a later code gate.
- KMS or external secret-manager integration is not required for the first local posture.
- Future KMS, cloud secret-manager, provider SDK, or operations secret-management integration requires a dependency adoption record and operations boundary.
- Four logical keys are required: credential lookup, credential verifier, token lookup, and token verifier.
- `verifier_key_id` is an internal key-set identifier, not a secret value, but it is not log-safe by default.
- Production behavior must fail closed when required key configuration is missing, malformed, too short, duplicated, or incomplete.
- Development and test keys must be explicit and must not introduce committed production-like secrets or default production keys.
- Runtime behavior remains deferred.

## Artifacts

- `docs/secret-configuration-verifier-key-loading-boundary.md`
- `docs/secret-configuration-verifier-key-loading-boundary.zh-CN.md`
- `decisions/ADR-0041-secret-configuration-verifier-key-loading-boundary.md`
- `changes/2026-05-15-define-secret-configuration-verifier-key-loading-boundary/`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `.arch/work-items.yaml`
- `modules/authentication/module.yaml`

## Open Questions

- Exact Go configuration struct names remain for a later implementation gate.
- Exact token and credential material generation helpers remain for a later gate.
- Exact verifier digest computation and constant-time comparison helpers remain for a later gate.
- Protobuf authentication messages and WebSocket proof carriers remain deferred.

## Follow-Up

- Advance a bounded token and credential material generation boundary, or the next explicitly ready preparation gate.
- Keep secret loading code, token generation, credential generation, verifier comparison, login execution, token validation, logout execution, cleanup, Protobuf messages, WebSocket proof carriers, authentication dependencies, repository changes, and migration changes behind later gates.

## Redaction Notes

No real secrets, tokens, credentials, device identifiers, account details, environment variable values, verifier key values, or private data are stored in this conversation log.
