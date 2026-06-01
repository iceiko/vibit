# Conversation: Select Next Pitaya-Aligned Direction After Acceptor And Connection Lifecycle Map

Date: 2026-06-01

## Context

The maintainer asked to continue toward Pitaya. The repository next-ready item was `W-0267 Select next Pitaya-aligned direction after acceptor and connection lifecycle map`, opened by `ADR-0174` after `node tools/vibit inspect pitaya-acceptor-connection --json`.

## Maintainer Narrative

Continue advancing in bounded steps and keep moving toward Pitaya alignment.

## Agent Response Summary

The agent selected `define_pitaya_aligned_session_binding_kick_disconnect_session_data_boundary_gate` as the next bounded Pitaya-aligned direction. This preserves the existing posture that Pitaya is deferred architecture vocabulary while choosing the next useful vocabulary surface for current vibit runtime facts: authenticated sessions, first-message connection binding, active connection registry state, logout, close handoff, and server-owned presence lifecycle.

The selection completed W-0267, accepted `ADR-0175`, registered `runtime.next_pitaya_aligned_direction_after_acceptor_connection_lifecycle_map`, and opened `M-196/W-0268 Define Pitaya-aligned session binding, kick/disconnect, and session data boundary gate` as next-ready.

## Decisions

- Accept `ADR-0175`.
- Complete W-0267.
- Open W-0268 as next-ready.
- Keep session binding behavior, kick/disconnect behavior, session data behavior and persistence, acceptor behavior, TCP acceptors, WebSocket behavior changes, connection lifecycle behavior changes, route handler implementation, runtime behavior, protocol changes, generated output, persistence, dependencies, hosted surfaces, SDKs, and direct compatibility deferred.

## Artifacts

- `changes/2026-05-26-select-next-pitaya-aligned-direction-after-acceptor-connection-lifecycle-map/`
- `decisions/ADR-0175-select-next-pitaya-aligned-direction-after-acceptor-connection-lifecycle-map.md`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/reference.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/modules.yaml`
- `tools/vibit`
- `rules/check-rules.json`

## Open Questions

- None for this selection slice.

## Follow-Up

- Complete W-0268: define the Pitaya-aligned session binding, kick/disconnect, and session data boundary gate.
- Do not add session binding behavior, kick/disconnect behavior, session data behavior or persistence, acceptor behavior, transport behavior changes, connection lifecycle behavior changes, protocol changes, generated output, persistence, dependencies, metrics/tracing behavior, hosted surfaces, SDKs, or direct Nakama/Pitaya API compatibility without a later bounded work item.

## Redaction Notes

No secrets, credentials, raw access tokens, verifier material, DSNs with credentials, database payloads, transport payloads, or local ignored file contents are recorded.
