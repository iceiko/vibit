# Conversation: Pitaya-Aligned Serializer And Message Forwarding Boundary Gate

Date: 2026-06-01
Status: Accepted
Related change: `changes/2026-05-26-define-pitaya-aligned-serializer-message-forwarding-boundary-gate/`
Related decision: `ADR-0170`

## Context

The maintainer asked to continue pushing toward Pitaya alignment. `W-0261` selected `define_pitaya_aligned_serializer_message_forwarding_boundary_gate` after the route handler pipeline map.

## Maintainer Narrative

Continue toward Pitaya, but keep the work bounded. Serializer and message forwarding vocabulary should become a gate before any serializer, forwarding, backend targeting, or distributed behavior is implemented.

## Agent Response Summary

The agent treated W-0262 as a gate-only work item. It added the serializer/message forwarding standard and translation, accepted ADR-0170, registered `runtime.pitaya_aligned_serializer_message_forwarding_boundary_gate`, completed W-0262, and opened W-0263 as the next source-first map follow-up.

RED checks confirmed the rule and change artifacts were initially absent:

```text
node tools/vibit inspect rule runtime.pitaya_aligned_serializer_message_forwarding_boundary_gate
# Unknown rule_id: runtime.pitaya_aligned_serializer_message_forwarding_boundary_gate

node tools/vibit check change define-pitaya-aligned-serializer-message-forwarding-boundary-gate --json
# change directory does not exist
```

## Decisions

- `ADR-0170` defines the Pitaya-aligned serializer and message forwarding boundary gate.
- The allowed vocabulary is `serializer_boundary`, `serializer_format`, `encode_boundary`, `decode_boundary`, `message_forwarding`, `forwarding_target`, `forwarding_envelope`, and `delivery_handoff`.
- Current vibit behavior remains generated Protobuf bridge functions, protocol-adapter-owned envelope encoding, application-owned outbound messages, metadata-only target scopes, and single-process WebSocket delivery.

## Artifacts

- `docs/pitaya-aligned-serializer-message-forwarding-boundary-gate.md`
- `docs/pitaya-aligned-serializer-message-forwarding-boundary-gate.zh-CN.md`
- `decisions/ADR-0170-pitaya-aligned-serializer-message-forwarding-boundary-gate.md`
- `changes/2026-05-26-define-pitaya-aligned-serializer-message-forwarding-boundary-gate/`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/reference.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `tools/vibit`
- `rules/check-rules.json`

## Open Questions

None for W-0262. W-0263 must remain source-first and must not implement serializer or forwarding behavior.

## Follow-Up

- Advance `W-0263 Implement Pitaya-aligned serializer and message forwarding source-first map`.

## Redaction Notes

No ignored credential file contents, token values, DSNs with credentials, raw credentials, raw access tokens, verifier digests, lookup digests, node credentials, transport payloads, route payload contents, or forwarding payload contents were recorded.
