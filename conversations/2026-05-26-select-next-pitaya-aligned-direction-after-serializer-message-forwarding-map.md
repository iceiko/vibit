# Conversation: Select Next Pitaya-Aligned Direction After Serializer And Message Forwarding Map

Date: 2026-06-01

## Context

The maintainer asked to continue toward Pitaya. The repository next-ready item was `W-0264 Select next Pitaya-aligned direction after serializer and message forwarding map`, opened by `ADR-0171` after `node tools/vibit inspect pitaya-serializer-forwarding --json`.

## Maintainer Narrative

Continue advancing in bounded steps and keep moving toward Pitaya alignment.

## Agent Response Summary

The agent selected `define_pitaya_aligned_acceptor_connection_lifecycle_boundary_gate` as the next bounded Pitaya-aligned direction. This preserves the existing posture that Pitaya is deferred architecture vocabulary while choosing the next useful vocabulary surface for current vibit runtime facts: WebSocket acceptor, server-observed connection id and epoch, first-message binding, active connection registry state, close handoff, and presence lifecycle.

The selection completed W-0264, accepted `ADR-0172`, registered `runtime.next_pitaya_aligned_direction_after_serializer_message_forwarding_map`, and opened `M-193/W-0265 Define Pitaya-aligned acceptor and connection lifecycle boundary gate` as next-ready.

## Decisions

- Accept `ADR-0172`.
- Complete W-0264.
- Open W-0265 as next-ready.
- Keep acceptor behavior, connection lifecycle behavior, session binding behavior, kick/disconnect behavior, serializer behavior, message forwarding behavior, runtime behavior, protocol changes, generated output, persistence, dependencies, hosted surfaces, SDKs, and direct compatibility deferred.

## Artifacts

- `changes/2026-05-26-select-next-pitaya-aligned-direction-after-serializer-message-forwarding-map/`
- `decisions/ADR-0172-select-next-pitaya-aligned-direction-after-serializer-message-forwarding-map.md`
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

- Complete W-0265: define the Pitaya-aligned acceptor and connection lifecycle boundary gate.
- Do not add acceptor behavior, transport behavior changes, session binding behavior, kick/disconnect behavior, serializer behavior, message forwarding behavior, protocol changes, generated output, persistence, dependencies, metrics/tracing behavior, hosted surfaces, SDKs, or direct Nakama/Pitaya API compatibility without a later bounded work item.

## Redaction Notes

No secrets, credentials, raw access tokens, verifier material, DSNs with credentials, database payloads, transport payloads, or local ignored file contents are recorded.
