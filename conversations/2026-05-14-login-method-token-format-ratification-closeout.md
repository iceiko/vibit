# Conversation: Login Method Token Format Ratification Closeout

Date: 2026-05-14
Participants: Maintainer, Agent
Related changes:

- `changes/2026-05-14-close-login-method-token-format-ratification-milestone/`

Related artifacts:

- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/reference.yaml`
- `docs/login-method-token-format-ratification.md`
- `docs/login-method-token-format-ratification.zh-CN.md`
- `decisions/ADR-0031-close-login-token-ratification-and-open-schema-gate.md`

## Context

The maintainer asked the agent to continue advancing bounded work items unless a true decision or blocker required confirmation. W-0073 followed W-0072, which added repository checks for the selected login/token posture.

## Maintainer Narrative

The maintainer wants the project to remain self-bootstrapping and controlled. The work should keep moving without stopping for routine technical details, but it must not cross long-term architectural or product decision boundaries implicitly.

## Agent Response Summary

The agent reviewed `M-013 Login Method And Token Format Ratification` against its completion criteria and closed the milestone.

The completed milestone ratified `device_credential_login`, opaque high-entropy access tokens, login-command token issuance, explicit request proof payloads, no refresh token in the first implementation, semantic authentication contracts, credential/token/session schema gates, and selected boundary checks.

The agent did not start runtime authentication implementation. Instead, it opened `M-014 Credential And Token Verifier Schema Ratification` with `W-0074 Define credential record schema boundary` as the next ready item.

## Decisions

- Mark `M-013` completed.
- Mark `W-0073` completed.
- Create `M-014 Credential And Token Verifier Schema Ratification`.
- Create `W-0074 Define credential record schema boundary` as the next ready work item.
- Preserve the rule that schema must be ratified before migrations, repositories, adapters, runtime lookup, handlers, routes, Protobuf changes, WebSocket changes, generated authentication shapes, or authentication implementation.
- Continue referencing Nakama and Pitaya for capability and vocabulary while rejecting direct public API copying.

## Artifacts

- `changes/2026-05-14-close-login-method-token-format-ratification-milestone/`
- `decisions/ADR-0031-close-login-token-ratification-and-open-schema-gate.md`
- `conversations/2026-05-14-login-method-token-format-ratification-closeout.md`

## Open Questions

- Exact credential record fields, uniqueness rules, verifier semantics, rotation behavior, and redaction model.
- Exact token verifier record fields, verifier algorithm, credential-token linkage, retention, cleanup, and revocation semantics.
- Whether generated authentication contract shapes should be introduced before or during runtime implementation.
- Whether request proof payloads become Protobuf messages or remain application-owned contract payloads in the first implementation milestone.

## Follow-Up

- Advance W-0074.
- Do not add credential migrations or runtime lookup until schema ratification explicitly authorizes the next step.

## Redaction Notes

No secrets, tokens, account details, or private data are stored in this conversation log.
