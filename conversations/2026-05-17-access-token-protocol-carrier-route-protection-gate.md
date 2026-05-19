# Conversation: Access Token Protocol Carrier Route Protection Gate

Date: 2026-05-17
Participants: Maintainer, Agent
Related changes:

- `changes/2026-05-17-define-access-token-protocol-carrier-route-protection-gate/`

Related artifacts:

- `docs/access-token-protocol-carrier-route-protection-gate.md`
- `docs/access-token-protocol-carrier-route-protection-gate.zh-CN.md`
- `decisions/ADR-0053-access-token-protocol-carrier-route-protection-gate.md`
- `.arch/work-items.yaml`

## Context

The previous direction gate selected access-token protocol carrier and route-protection gate work after service-local authentication was implemented.

## Maintainer Narrative

The maintainer authorized continued progress for five hours.

## Agent Response Summary

The agent defined a gate-only standard for future request-level access-token carrier and route-protection work. The gate selects an explicit Protobuf payload wrapper as the first future proof carrier posture and keeps the existing envelope unchanged. It requires application-owned route policy and validated player identity before protected domain dispatch.

No protocol sources, generated files, route-protection code, WebSocket handshake authentication, session persistence, startup wiring, repository changes, migrations, dependencies, logout, refresh, cleanup, or production behavior were added by the gate.

## Decisions

- Use request-level validation as the selected future posture for this milestone.
- Use a Protobuf payload wrapper candidate, `vibit.authentication.v1.AuthenticatedRequest`, as the first future carrier posture.
- Keep WebSocket transport credential-neutral.
- Keep existing Protobuf envelope fields as metadata only.
- Require application-owned route policy before protected domain dispatch.

## Artifacts

- `docs/access-token-protocol-carrier-route-protection-gate.md`
- `docs/access-token-protocol-carrier-route-protection-gate.zh-CN.md`
- `decisions/ADR-0053-access-token-protocol-carrier-route-protection-gate.md`
- `changes/2026-05-17-define-access-token-protocol-carrier-route-protection-gate/`

## Open Questions

- Exact wrapper `.proto` source content remains for the future implementation slice.
- Exact route-policy Go API remains for the future implementation slice.
- Startup wiring remains deferred behind a later composition work item.

## Follow-Up

- Implement the bounded wrapper and route-policy slice after the next work item authorizes `.proto`, generated output, protocol adapter tests, application route policy, and startup deferrals.

## Redaction Notes

No secrets, raw access tokens, generated credentials, digest bytes, HMAC input bytes, verifier keys, private account data, or GitHub tokens are recorded in this conversation log.
