# Conversation: Service Local Authentication Next Direction

Date: 2026-05-17
Participants: Maintainer, Agent
Related changes:

- `changes/2026-05-17-confirm-next-direction-after-service-local-authentication/`

Related artifacts:

- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/reference.yaml`

## Context

The work queue was blocked at `M-040/W-0112` after service-local device credential login and access-token validation were completed.

## Maintainer Narrative

The maintainer asked:

```text
继续推进5小时。
```

## Agent Response Summary

The agent interpreted the request as authorization to close the blocked direction-confirmation gate and use professional judgment to choose the next conservative direction. The selected direction was:

```text
expose_access_token_protocol_carrier_and_route_protection_gate
```

The choice starts with a gate-only milestone before any protocol carrier, route-protection, session, WebSocket handshake, startup wiring, generated output, repository, migration, dependency, logout, refresh, cleanup, or broader production authentication behavior is implemented.

## Decisions

- Close `M-040/W-0112`.
- Select `expose_access_token_protocol_carrier_and_route_protection_gate`.
- Create `M-041` with `W-0113` as the first next-ready gate work item.

## Artifacts

- `changes/2026-05-17-confirm-next-direction-after-service-local-authentication/`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/reference.yaml`

## Open Questions

- The later implementation gate still needs to decide exact `.proto` source content, generated output timing, adapter handoff, route policy API, and startup composition sequencing.

## Follow-Up

- Define the access-token protocol carrier and route-protection gate as a bounded documentation, ADR, manifest, and check-rule change.

## Redaction Notes

No secrets, raw access tokens, generated credentials, digest bytes, HMAC input bytes, verifier keys, private account data, or GitHub tokens are recorded in this conversation log.
