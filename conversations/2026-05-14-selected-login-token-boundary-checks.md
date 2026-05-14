# Conversation: Selected Login Token Boundary Checks

Date: 2026-05-14
Participants: Maintainer, Agent
Related changes:

- `changes/2026-05-14-add-selected-login-token-boundary-checks/`

Related artifacts:

- `docs/selected-login-token-boundary-checks.md`
- `docs/selected-login-token-boundary-checks.zh-CN.md`
- `decisions/ADR-0030-selected-login-token-boundary-checks.md`
- `rules/check-rules.json`
- `tools/vibit`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `.arch/work-items.yaml`

## Context

The maintainer asked the agent to continue advancing ten work steps unless a true decision or blocker required confirmation. W-0072 followed W-0071, which had defined credential, token verifier, external identity, runtime session, and audit schema gates.

## Maintainer Narrative

The maintainer wants vibit to keep moving under professional agent judgment while preserving self-bootstrapping control. Routine technical work inside ratified boundaries should proceed without unnecessary interruptions. The project should continue to reference Nakama and Pitaya, but vibit's distinction is Agent-Native maintainability rather than copying an existing game-server API.

## Agent Response Summary

The agent added a selected login/token boundary standard and a new machine-readable runtime check rule: `runtime.selected_login_token_boundary`.

The rule protects the selected `device_credential_login` and opaque access-token posture from accidental implementation. It verifies status markers, contract blocked-until markers, generated output deferral, Protobuf source deferral, WebSocket carrier deferral, migration deferral, player account storage separation, forbidden auth dependencies, and stable JSON rule metadata.

The agent did not implement authentication, token behavior, credential storage, token storage, runtime sessions, external identity linking, migrations, repositories, adapters, Protobuf messages, generated authentication shapes, WebSocket routes, WebSocket handshake auth, or live service checks.

## Decisions

- Add `docs/selected-login-token-boundary-checks.md` and the Simplified Chinese translation.
- Record ADR-0030 for selected login/token boundary checks.
- Register `runtime.selected_login_token_boundary` in the rule catalog.
- Run the selected boundary check through `node tools/vibit check runtime --json` and `node tools/vibit check all --json`.
- Keep default verification local and static.
- Preserve the selected posture as semantic-contract-only until a later implementation milestone changes it explicitly.

## Artifacts

- `docs/selected-login-token-boundary-checks.md`
- `docs/selected-login-token-boundary-checks.zh-CN.md`
- `decisions/ADR-0030-selected-login-token-boundary-checks.md`
- `changes/2026-05-14-add-selected-login-token-boundary-checks/`
- `rules/check-rules.json`
- `tools/vibit`

## Open Questions

- Which future milestone will authorize credential schema ratification.
- Which future milestone will authorize token verifier schema ratification.
- Whether authentication generated contract shapes should be generated before runtime code or inside the first implementation milestone.
- Whether explicit request proof payloads become Protobuf messages, generated contract shapes, or an application-owned envelope extension in the first implementation milestone.

## Follow-Up

- Verify W-0072.
- Close M-013 if the milestone criteria are satisfied.
- Do not start production authentication implicitly while closing the milestone.

## Redaction Notes

No secrets, tokens, account details, or private data are stored in this conversation log.
