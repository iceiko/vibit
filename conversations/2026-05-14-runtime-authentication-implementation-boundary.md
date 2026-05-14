# Conversation: Runtime Authentication Implementation Boundary

Date: 2026-05-14
Participants: Maintainer, Agent
Related changes:

- `changes/2026-05-14-define-runtime-authentication-implementation-boundary/`

Related artifacts:

- `docs/runtime-authentication-implementation-boundary.md`
- `docs/runtime-authentication-implementation-boundary.zh-CN.md`
- `decisions/ADR-0036-runtime-authentication-implementation-boundary.md`
- `.arch/work-items.yaml`

## Context

The maintainer asked the agent to continue through bounded work items unless a true decision or blocker required confirmation. The previous milestone completed the authentication PostgreSQL adapter while preserving runtime authentication deferral.

## Maintainer Narrative

The maintainer wants the project to move forward under professional judgment without stopping for routine technical sequencing, while still stopping for real architecture or product decisions.

The maintainer has also emphasized that vibit should remain self-bootstrapping and controlled, and that future work should avoid drifting away from the goal of an agent-native server framework.

## Agent Response Summary

The agent defined the runtime authentication implementation boundary as a planning artifact before adding runtime behavior.

The boundary assigns future authentication orchestration to the application layer, keeps `authentication.Repository` as the module-owned storage-neutral interface, keeps the PostgreSQL adapter persistence-only, preserves protocol and transport neutrality, and splits token generation, verifier comparison, login execution, token validation, logout execution, cleanup jobs, Protobuf messages, WebSocket carriers, generated authentication shapes, and dependencies into separate gates.

The agent referenced Nakama and Pitaya for capability and vocabulary alignment without copying their APIs.

## Decisions

- Define `ADR-0036 Runtime Authentication Implementation Boundary`.
- Add `docs/runtime-authentication-implementation-boundary.md` and its Simplified Chinese translation.
- Preserve metadata-only identity as non-authenticated.
- Preserve runtime authentication behavior deferral.
- Plan the next step as boundary-check and milestone-closeout work, not login implementation.

## Artifacts

- `changes/2026-05-14-define-runtime-authentication-implementation-boundary/`
- `docs/runtime-authentication-implementation-boundary.md`
- `docs/runtime-authentication-implementation-boundary.zh-CN.md`
- `decisions/ADR-0036-runtime-authentication-implementation-boundary.md`
- `conversations/2026-05-14-runtime-authentication-implementation-boundary.md`

## Open Questions

- Whether generated Go authentication contract shapes should appear before application service interfaces.
- Exact verifier algorithms and redaction tests for credential proof and access-token proof.
- Exact Protobuf message shape for login and explicit request proof payloads.
- Whether WebSocket proof carriers remain deferred indefinitely or become a later optional carrier.

## Follow-Up

- Advance the next boundary-check or closeout work item.
- Do not implement login, token generation, token validation, verifier comparison, logout execution, cleanup jobs, Protobuf authentication messages, WebSocket proof carriers, generated authentication shapes, or authentication dependencies in this planning change.

## Redaction Notes

No secrets, tokens, account details, or private data are stored in this conversation log.
